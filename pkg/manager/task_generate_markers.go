package manager

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type GenerateMarkersTask struct {
	Scene  models.Scene
	useMD5 bool
}

func (t *GenerateMarkersTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
	if len(sceneMarkers) == 0 {
		return
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.Scene.Path)
	if err != nil {
		logger.Errorf("error reading video file: %s", err.Error())
		return
	}

	sceneHash := t.Scene.GetHash(t.useMD5)

	// Make the folder for the scenes markers
	// use the existing folder if present, otherwise create one from the hash
	if !t.dirExists(t.Scene.OSHash) && !t.dirExists(t.Scene.Checksum) {
		markersFolder := filepath.Join(instance.Paths.Generated.Markers, sceneHash)
		utils.EnsureDir(markersFolder)
	}

	encoder := ffmpeg.NewEncoder(instance.FFMPEGPath)
	for i, sceneMarker := range sceneMarkers {
		index := i + 1
		logger.Progressf("[generator] <%s> scene marker %d of %d", sceneHash, index, len(sceneMarkers))

		seconds := int(sceneMarker.Seconds)

		videoExists := t.videoExists(sceneHash, seconds)
		imageExists := t.imageExists(sceneHash, seconds)

		baseFilename := strconv.Itoa(seconds)

		options := ffmpeg.SceneMarkerOptions{
			ScenePath: t.Scene.Path,
			Seconds:   seconds,
			Width:     640,
		}

		if !videoExists {
			videoFilename := baseFilename + ".mp4"
			videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneHash, seconds)

			options.OutputPath = instance.Paths.Generated.GetTmpPath(videoFilename) // tmp output in case the process ends abruptly
			if err := encoder.SceneMarkerVideo(*videoFile, options); err != nil {
				logger.Errorf("[generator] failed to generate marker video: %s", err)
			} else {
				_ = os.Rename(options.OutputPath, videoPath)
				logger.Debug("created marker video: ", videoPath)
			}
		}

		if !imageExists {
			imageFilename := baseFilename + ".webp"
			imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneHash, seconds)

			options.OutputPath = instance.Paths.Generated.GetTmpPath(imageFilename) // tmp output in case the process ends abruptly
			if err := encoder.SceneMarkerImage(*videoFile, options); err != nil {
				logger.Errorf("[generator] failed to generate marker image: %s", err)
			} else {
				_ = os.Rename(options.OutputPath, imagePath)
				logger.Debug("created marker image: ", imagePath)
			}
		}
	}
}

func (t *GenerateMarkersTask) isMarkerNeeded() int {
	markers := 0
	qb := models.NewSceneMarkerQueryBuilder()
	sceneMarkers, _ := qb.FindBySceneID(t.Scene.ID, nil)
	if len(sceneMarkers) == 0 {
		return 0
	}

	sceneHash := t.Scene.GetHash(t.useMD5)
	for _, sceneMarker := range sceneMarkers {
		seconds := int(sceneMarker.Seconds)

		if !t.markerExists(sceneHash, seconds) {
			markers++
		}
	}
	return markers
}

func (t *GenerateMarkersTask) markerExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneChecksum, seconds)
	imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneChecksum, seconds)
	videoExists, _ := utils.FileExists(videoPath)
	imageExists, _ := utils.FileExists(imagePath)

	return videoExists && imageExists
}

func (t *GenerateMarkersTask) dirExists(sceneChecksumNull sql.NullString) bool {
	if !sceneChecksumNull.Valid {
		return false
	}
	sceneChecksum := sceneChecksumNull.String

	markersFolder := filepath.Join(instance.Paths.Generated.Markers, sceneChecksum)
	dirExists, _ := utils.DirExists(markersFolder)

	return dirExists
}

func (t *GenerateMarkersTask) videoExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	videoPath := instance.Paths.SceneMarkers.GetStreamPath(sceneChecksum, seconds)
	videoExists, _ := utils.FileExists(videoPath)

	return videoExists
}

func (t *GenerateMarkersTask) imageExists(sceneChecksum string, seconds int) bool {
	if sceneChecksum == "" {
		return false
	}

	imagePath := instance.Paths.SceneMarkers.GetStreamPreviewImagePath(sceneChecksum, seconds)
	imageExists, _ := utils.FileExists(imagePath)

	return imageExists
}
