// TODO: document in README.md
package main

import (
	"context"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/hash/videophash"
)

func customusage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s [OPTIONS] VIDEOFILE...\n\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
}

func printphash(inputfile string, quiet *bool) error {
	ffmpegPath, ffprobePath := ffmpeg.GetPaths(nil)
	FFMPEG := ffmpeg.NewEncoder(ffmpegPath)
	FFMPEG.InitHWSupport(context.TODO())

	FFPROBE := ffmpeg.FFProbe(ffprobePath)
	ffvideoFile, err := FFPROBE.NewVideoFile(inputfile)
	if err != nil {
		return err
	}

	// All we need for videophash.Generate() is
	// videoFile.Path (from BaseFile)
	// videoFile.Duration
	// The rest of the struct isn't needed.
	vf := &file.VideoFile{
		BaseFile: &file.BaseFile{Path: inputfile},
		Duration: ffvideoFile.FileDuration,
	}

	phash, err := videophash.Generate(FFMPEG, vf)
	if err != nil {
		return err
	}

	if *quiet {
		fmt.Printf("%x\n", *phash)
	} else {
		fmt.Printf("%x %v\n", *phash, vf.Path)
	}
	return nil
}

func main() {
	flag.Usage = customusage
	quiet := flag.BoolP("quiet", "q", false, "print only the phash")
	help := flag.BoolP("help", "h", false, "print this help output")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Missing VIDEOFILE argument.\n")
		flag.Usage()
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Files will be processed sequentially! Consier using GNU Parallel.")
		fmt.Fprintf(os.Stderr, "Example: parallel %v ::: *.mp4\n", os.Args[0])
	}

	for _, item := range args {
		if err := printphash(item, quiet); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
