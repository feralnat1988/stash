package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	MimeHLS    string = "application/vnd.apple.mpegurl"
	MimeMpegTS string = "video/MP2T"

	segmentLength = 2

	maxSegmentWait  = 15 * time.Second
	monitorInterval = 200 * time.Millisecond

	// segment gap before counting a request as a seek and
	// restarting the transcode process at the requested segment
	maxSegmentGap = 5

	// maximum number of segments to generate
	// ahead of the currently streaming segment
	maxSegmentBuffer = 15

	// maximum idle time between segment requests before
	// stopping transcode and deleting cache folder
	maxIdleTime = 30 * time.Second
)

type StreamType struct {
	Name          string
	SegmentType   *SegmentType
	ServeManifest func(sm *StreamManager, w http.ResponseWriter, r *http.Request, vf *file.VideoFile, resolution string)
	Args          func(segment int, videoFilter VideoFilter, videoOnly bool, outputDir string) Args
}

var (
	StreamTypeHLS = &StreamType{
		Name:          "hls",
		SegmentType:   SegmentTypeTS,
		ServeManifest: serveHLSManifest,
		Args: func(segment int, videoFilter VideoFilter, videoOnly bool, outputDir string) (args Args) {
			args = append(args,
				"-c:v", "libx264",
				"-pix_fmt", "yuv420p",
				"-preset", "veryfast",
				"-crf", "25",
				"-flags", "+cgop",
				"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentLength),
				"-sc_threshold", "0",
			)
			args = args.VideoFilter(videoFilter)
			if videoOnly {
				args = append(args, "-an")
			} else {
				args = append(args,
					"-c:a", "aac",
					"-ac", "2",
				)
			}
			args = append(args,
				"-sn",
				"-copyts",
				"-avoid_negative_ts", "disabled",
				"-f", "hls",
				"-start_number", fmt.Sprint(segment),
				"-hls_time", "2",
				"-hls_segment_type", "mpegts",
				"-hls_playlist_type", "vod",
				"-hls_segment_filename", filepath.Join(outputDir, ".%d.ts"),
				filepath.Join(outputDir, "manifest.m3u8"),
			)
			return
		},
	}
	StreamTypeHLSCopy = &StreamType{
		Name:          "hls-copy",
		SegmentType:   SegmentTypeTS,
		ServeManifest: serveHLSManifest,
		Args: func(segment int, videoFilter VideoFilter, videoOnly bool, outputDir string) (args Args) {
			args = append(args,
				"-c:v", "copy",
			)
			if videoOnly {
				args = append(args, "-an")
			} else {
				args = append(args,
					"-c:a", "aac",
					"-ac", "2",
				)
			}
			args = append(args,
				"-sn",
				"-copyts",
				"-avoid_negative_ts", "disabled",
				"-f", "hls",
				"-start_number", fmt.Sprint(segment),
				"-hls_time", "2",
				"-hls_segment_type", "mpegts",
				"-hls_playlist_type", "vod",
				"-hls_segment_filename", filepath.Join(outputDir, ".%d.ts"),
				filepath.Join(outputDir, "manifest.m3u8"),
			)
			return
		},
	}
)

type SegmentType struct {
	Format       string
	MimeType     string
	MakeFilename func(segment int) string
	ParseSegment func(str string) (int, error)
}

var (
	SegmentTypeTS = &SegmentType{
		Format:   "%d.ts",
		MimeType: MimeMpegTS,
		MakeFilename: func(segment int) string {
			return fmt.Sprintf("%d.ts", segment)
		},
		ParseSegment: func(str string) (int, error) {
			segment, err := strconv.Atoi(str)
			if err != nil || segment < 0 {
				err = ErrInvalidSegment
			}
			return segment, err
		},
	}
)

var ErrInvalidSegment = errors.New("invalid segment")

type StreamOptions struct {
	StreamType *StreamType
	VideoFile  *file.VideoFile
	Resolution string
	Hash       string
	Segment    string
}

type transcodeProcess struct {
	cmd         *exec.Cmd
	context     context.Context
	cancel      context.CancelFunc
	cancelled   bool
	outputDir   string
	segmentType *SegmentType
	segment     int
}

type waitingSegment struct {
	segmentType *SegmentType
	idx         int
	file        string
	path        string
	accessed    time.Time
	available   chan error
	done        atomic.Bool
}

type runningStream struct {
	dir              string
	streamType       *StreamType
	vf               *file.VideoFile
	maxTranscodeSize int
	outputDir        string

	waitingSegments []*waitingSegment
	tp              *transcodeProcess
	lastAccessed    time.Time
	lastSegment     int
}

func (t StreamType) String() string {
	return t.Name
}

func (t StreamType) FileDir(hash string, maxTranscodeSize int) string {
	if maxTranscodeSize == 0 {
		return fmt.Sprintf("%s_%s", hash, t)
	} else {
		return fmt.Sprintf("%s_%s_%d", hash, t, maxTranscodeSize)
	}
}

func (s *runningStream) makeStreamArgs(segment int) Args {
	args := Args{"-hide_banner"}
	args = args.LogLevel(LogLevelError)

	if segment > 0 {
		args = args.Seek(float64(segment * segmentLength))
	}

	args = args.Input(s.vf.Path)

	videoOnly := ProbeAudioCodec(s.vf.AudioCodec) == MissingUnsupported

	var videoFilter VideoFilter
	videoFilter = videoFilter.ScaleMax(s.vf.Width, s.vf.Height, s.maxTranscodeSize)

	args = append(args, s.streamType.Args(segment, videoFilter, videoOnly, s.outputDir)...)

	return args
}

// checkSegments renames temp segments that have been completely generated.
// existing segments are not replaced - if a segment is generated
// multiple times, then only the first one is kept.
func (tp *transcodeProcess) checkSegments() {
	doSegment := func(filename string) {
		if filename != "" {
			oldPath := filepath.Join(tp.outputDir, filename)
			newPath := filepath.Join(tp.outputDir, filename[1:])
			if !segmentExists(newPath) {
				_ = os.Rename(oldPath, newPath)
			} else {
				os.Remove(oldPath)
			}
		}
	}

	processState := tp.cmd.ProcessState
	var lastFilename string
	for i := tp.segment; ; i++ {
		filename := fmt.Sprintf("."+tp.segmentType.Format, i)
		if segmentExists(filepath.Join(tp.outputDir, filename)) {
			// this segment exists so the previous segment is valid
			doSegment(lastFilename)
		} else {
			// if the transcode process has exited then
			// we need to do something with the last segment
			if processState != nil {
				if processState.Success() {
					// if the process exited successfully then
					// count the last segment as valid
					doSegment(lastFilename)
				} else if lastFilename != "" {
					// if the process exited unsuccessfully then just delete
					// the last segment, it's probably incomplete
					os.Remove(filepath.Join(tp.outputDir, lastFilename))
				}
			}
			break
		}

		lastFilename = filename
		tp.segment = i
	}
}

func lastSegment(vf *file.VideoFile) int {
	return int(math.Ceil(vf.Duration/segmentLength)) - 1
}

func segmentExists(path string) bool {
	exists, _ := fsutil.FileExists(path)
	return exists
}

// serveHLSManifest serves a generated HLS playlist. The URLs for the segments
// are of the form {r.URL}/%d.ts{?urlQuery} where %d is the segment index.
func serveHLSManifest(sm *StreamManager, w http.ResponseWriter, r *http.Request, vf *file.VideoFile, resolution string) {
	if sm.cacheDir == "" {
		logger.Error("[transcode] cannot live transcode with HLS because cache dir is unset")
		http.Error(w, "cannot live transcode with HLS because cache dir is unset", http.StatusServiceUnavailable)
		return
	}

	probeResult, err := sm.ffprobe.NewVideoFile(vf.Path)
	if err != nil {
		logger.Warnf("[transcode] error generating HLS manifest: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	baseUrl := *r.URL
	baseUrl.RawQuery = ""
	baseURL := baseUrl.String()

	var urlQuery string
	if resolution != "" {
		urlQuery = fmt.Sprintf("?resolution=%s", resolution)
	}

	var buf bytes.Buffer

	fmt.Fprint(&buf, "#EXTM3U\n")

	fmt.Fprint(&buf, "#EXT-X-VERSION:3\n")
	fmt.Fprint(&buf, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprintf(&buf, "#EXT-X-TARGETDURATION:%d\n", segmentLength)
	fmt.Fprint(&buf, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := probeResult.FileDuration
	segment := 0

	for leftover > 0 {
		thisLength := float64(segmentLength)
		if leftover < thisLength {
			thisLength = leftover
		}

		fmt.Fprintf(&buf, "#EXTINF:%f,\n", thisLength)
		fmt.Fprintf(&buf, "%s/%d.ts%s\n", baseURL, segment, urlQuery)

		leftover -= thisLength
		segment++
	}

	fmt.Fprint(&buf, "#EXT-X-ENDLIST\n")

	w.Header().Set("Content-Type", MimeHLS)
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(buf.Bytes()))
}

func (sm *StreamManager) ServeManifest(w http.ResponseWriter, r *http.Request, streamType *StreamType, vf *file.VideoFile, resolution string) {
	streamType.ServeManifest(sm, w, r, vf, resolution)
}

func (sm *StreamManager) serveWaitingSegment(w http.ResponseWriter, r *http.Request, segment *waitingSegment) {
	select {
	case <-r.Context().Done():
		break
	case err := <-segment.available:
		if err == nil {
			logger.Tracef("[transcode] streaming segment file %s", segment.file)
			w.Header().Set("Content-Type", segment.segmentType.MimeType)
			// Prevent caching as segments are generated on the fly
			w.Header().Add("Cache-Control", "no-cache")
			http.ServeFile(w, r, segment.path)
		} else if !errors.Is(err, context.Canceled) {
			logger.Errorf("[transcode] %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	segment.done.Store(true)
}

func (sm *StreamManager) ServeSegment(w http.ResponseWriter, r *http.Request, options StreamOptions) {
	if sm.cacheDir == "" {
		logger.Error("[transcode] cannot live transcode files because cache dir is unset")
		http.Error(w, "cannot live transcode files because cache dir is unset", http.StatusServiceUnavailable)
		return
	}

	if options.Hash == "" {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		return
	}

	streamType := options.StreamType

	segment, err := streamType.SegmentType.ParseSegment(options.Segment)
	// error if segment is past the end of the video
	if err != nil || segment > lastSegment(options.VideoFile) {
		http.Error(w, "invalid segment", http.StatusBadRequest)
		return
	}

	maxTranscodeSize := sm.config.GetMaxStreamingTranscodeSize().GetMaxResolution()
	if options.Resolution != "" {
		maxTranscodeSize = models.StreamingResolutionEnum(options.Resolution).GetMaxResolution()
	}

	dir := options.StreamType.FileDir(options.Hash, maxTranscodeSize)
	outputDir := filepath.Join(sm.cacheDir, dir)

	name := streamType.SegmentType.MakeFilename(segment)
	file := filepath.Join(dir, name)

	sm.streamsMutex.Lock()

	stream := sm.runningStreams[dir]
	if stream == nil {
		stream = &runningStream{
			dir:              dir,
			streamType:       options.StreamType,
			vf:               options.VideoFile,
			maxTranscodeSize: maxTranscodeSize,
			outputDir:        outputDir,

			// initialize to cap 10 to avoid reallocations
			waitingSegments: make([]*waitingSegment, 0, 10),
		}
		sm.runningStreams[dir] = stream
	}

	now := time.Now()
	stream.lastAccessed = now
	if segment != -1 {
		stream.lastSegment = segment
	}

	waitingSegment := &waitingSegment{
		segmentType: streamType.SegmentType,
		idx:         segment,
		file:        file,
		path:        filepath.Join(sm.cacheDir, file),
		accessed:    now,
		available:   make(chan error, 1),
	}
	stream.waitingSegments = append(stream.waitingSegments, waitingSegment)

	sm.streamsMutex.Unlock()

	sm.serveWaitingSegment(w, r, waitingSegment)
}

// assume lock is held
func (sm *StreamManager) startTranscode(stream *runningStream, segment int, done chan<- error) {
	// generate segment 0 if init segment requested
	if segment == -1 {
		segment = 0
	}

	logger.Debugf("[transcode] starting transcode for %s at segment #%d", stream.dir, segment)

	if err := os.MkdirAll(stream.outputDir, os.ModePerm); err != nil {
		done <- err
		return
	}

	lockCtx := sm.lockManager.ReadLock(sm.context, stream.vf.Path)

	args := stream.makeStreamArgs(segment)
	cmd := sm.encoder.Command(lockCtx, args)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Errorf("[transcode] ffmpeg stderr not available: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Errorf("[transcode] ffmpeg stdout not available: %v", err)
	}

	logger.Tracef("[transcode] running %s", cmd)
	if err := cmd.Start(); err != nil {
		lockCtx.Cancel()
		done <- fmt.Errorf("error starting transcode process: %w", err)
		return
	}

	tp := &transcodeProcess{
		cmd:         cmd,
		context:     lockCtx,
		cancel:      lockCtx.Cancel,
		outputDir:   stream.outputDir,
		segmentType: stream.streamType.SegmentType,
		segment:     segment,
	}
	stream.tp = tp

	go func() {
		errStr, _ := io.ReadAll(stderr)
		outStr, _ := io.ReadAll(stdout)

		errCmd := cmd.Wait()

		var err error

		// don't log error if cancelled
		if !tp.cancelled {
			e := string(errStr)
			if e == "" {
				e = string(outStr)
			}
			if e != "" {
				err = errors.New(e)
			} else {
				err = errCmd
			}

			if err != nil {
				err = fmt.Errorf("[transcode] ffmpeg error when running command <%s>: %w", strings.Join(cmd.Args, " "), err)
			}
		}

		sm.streamsMutex.Lock()

		// make sure that cancel is called to prevent memory leaks
		tp.cancel()

		// clear remaining segments after ffmpeg exit
		tp.checkSegments()

		if stream.tp == tp {
			stream.tp = nil
		}

		sm.streamsMutex.Unlock()

		done <- err
	}()
}

// assume lock is held
func (sm *StreamManager) stopTranscode(stream *runningStream) {
	tp := stream.tp
	if tp != nil {
		tp.cancel()
		tp.cancelled = true
	}
}

func (sm *StreamManager) checkTranscode(stream *runningStream, now time.Time) {
	if len(stream.waitingSegments) == 0 && stream.lastAccessed.Add(maxIdleTime).Before(now) {
		// Stream expired. Cancel the transcode process and delete the files
		logger.Debugf("[transcode] stream for %s not accessed recently. Cancelling transcode and removing files", stream.dir)

		sm.stopTranscode(stream)
		sm.removeTranscodeFiles(stream)

		delete(sm.runningStreams, stream.dir)
		return
	}

	if stream.tp != nil {
		segmentType := stream.streamType.SegmentType
		segment := stream.lastSegment
		// if all segments up to maxSegmentBuffer exist, stop transcode
		for i := segment; i < segment+maxSegmentBuffer; i++ {
			if !segmentExists(filepath.Join(stream.outputDir, segmentType.MakeFilename(i))) {
				return
			}
		}

		logger.Debugf("[transcode] stopping transcode for %s, buffer is full", stream.dir)
		sm.stopTranscode(stream)
	}
}

func (s *waitingSegment) checkAvailable(now time.Time) bool {
	if segmentExists(s.path) {
		s.available <- nil
		return true
	} else if s.accessed.Add(maxSegmentWait).Before(now) {
		s.available <- fmt.Errorf("timed out waiting for segment file %s to be generated", s.file)
		return true
	}
	return false
}

// ensureTranscode will start a new transcode process if the transcode
// is more than maxSegmentGap behind the requested segment
func (sm *StreamManager) ensureTranscode(stream *runningStream, segment *waitingSegment) bool {
	segmentIdx := segment.idx
	tp := stream.tp
	if tp == nil {
		sm.startTranscode(stream, segmentIdx, segment.available)
		return true
	} else if segmentIdx < tp.segment || tp.segment+maxSegmentGap < segmentIdx {
		// only stop the transcode process here - it will be restarted only
		// after the old process exits as stream.tp will then be nil.
		sm.stopTranscode(stream)
		return true
	}
	return false
}

// runs every monitorInterval
func (sm *StreamManager) monitorStreams() {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	now := time.Now()

	for _, stream := range sm.runningStreams {
		if stream.tp != nil {
			stream.tp.checkSegments()
		}

		transcodeStarted := false
		temp := stream.waitingSegments[:0]
		for _, segment := range stream.waitingSegments {
			remove := false
			if segment.done.Load() || segment.checkAvailable(now) {
				remove = true
			} else if !transcodeStarted {
				transcodeStarted = sm.ensureTranscode(stream, segment)
			}
			if !remove {
				temp = append(temp, segment)
			}
		}
		stream.waitingSegments = temp

		if !transcodeStarted {
			sm.checkTranscode(stream, now)
		}
	}
}

// assume lock is held
func (sm *StreamManager) removeTranscodeFiles(stream *runningStream) {
	path := stream.outputDir
	if err := os.RemoveAll(path); err != nil {
		logger.Warnf("[transcode] error removing segment directory %s: %v", path, err)
	}
}

// stopAndRemoveAll stops all current streams and removes all cache files
func (sm *StreamManager) stopAndRemoveAll() {
	sm.streamsMutex.Lock()
	defer sm.streamsMutex.Unlock()

	for _, stream := range sm.runningStreams {
		for _, segment := range stream.waitingSegments {
			if len(segment.available) == 0 {
				segment.available <- context.Canceled
			}
		}
		sm.stopTranscode(stream)
		sm.removeTranscodeFiles(stream)
	}

	// ensure nothing else can use the map
	sm.runningStreams = nil
}
