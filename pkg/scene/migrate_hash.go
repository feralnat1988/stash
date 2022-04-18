package scene

import (
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models/paths"
)

func MigrateHash(p *paths.Paths, oldHash string, newHash string) {
	oldPath := filepath.Join(p.Generated.Markers, oldHash)
	newPath := filepath.Join(p.Generated.Markers, newHash)
	migrateSceneFiles(oldPath, newPath)

	scenePaths := p.Scene
	oldPath = scenePaths.GetThumbnailScreenshotPath(oldHash)
	newPath = scenePaths.GetThumbnailScreenshotPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetScreenshotPath(oldHash)
	newPath = scenePaths.GetScreenshotPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetVideoPreviewPath(oldHash)
	newPath = scenePaths.GetVideoPreviewPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetWebpPreviewPath(oldHash)
	newPath = scenePaths.GetWebpPreviewPath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetTranscodePath(oldHash)
	newPath = scenePaths.GetTranscodePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetSpriteVttFilePath(oldHash)
	newPath = scenePaths.GetSpriteVttFilePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetSpriteImageFilePath(oldHash)
	newPath = scenePaths.GetSpriteImageFilePath(newHash)
	migrateSceneFiles(oldPath, newPath)

	oldPath = scenePaths.GetInteractiveHeatmapPath(oldHash)
	newPath = scenePaths.GetInteractiveHeatmapPath(newHash)
	migrateSceneFiles(oldPath, newPath)
}

func migrateSceneFiles(oldName, newName string) {
	oldExists, err := fsutil.FileExists(oldName)
	if err != nil && !os.IsNotExist(err) {
		logger.Errorf("Error checking existence of %s: %s", oldName, err.Error())
		return
	}

	if oldExists {
		logger.Infof("renaming %s to %s", oldName, newName)
		if err := os.Rename(oldName, newName); err != nil {
			logger.Errorf("error renaming %s to %s: %s", oldName, newName, err.Error())
		}
	}
}
