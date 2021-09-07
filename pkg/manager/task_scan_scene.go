package manager

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/utils"
)

func (t *ScanTask) scanScene() *models.Scene {
	logError := func(err error) *models.Scene {
		logger.Error(err.Error())
		return nil
	}

	var retScene *models.Scene
	var s *models.Scene

	if err := t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		var err error
		s, err = r.Scene().FindByPath(t.FilePath)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return nil
	}

	fileModTime, err := t.getFileModTime()
	if err != nil {
		return logError(err)
	}
	interactive := t.getInteractive()

	if s != nil {
		// if file mod time is not set, set it now
		if !s.FileModTime.Valid {
			logger.Infof("setting file modification time on %s", t.FilePath)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Scene()
				if _, err := scene.UpdateFileModTime(qb, s.ID, models.NullSQLiteTimestamp{
					Timestamp: fileModTime,
					Valid:     true,
				}); err != nil {
					return err
				}

				// update our copy of the scene
				var err error
				s, err = qb.Find(s.ID)
				return err
			}); err != nil {
				return logError(err)
			}
		}

		// if the mod time of the file is different than that of the associated
		// scene, then recalculate the checksum and regenerate the thumbnail
		modified := t.isFileModified(fileModTime, s.FileModTime)
		config := config.GetInstance()
		if modified || !s.Size.Valid {
			oldHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
			s, err = t.rescanScene(s, fileModTime)
			if err != nil {
				return logError(err)
			}

			// Migrate any generated files if the hash has changed
			newHash := s.GetHash(config.GetVideoFileNamingAlgorithm())
			if newHash != oldHash {
				MigrateHash(oldHash, newHash)
			}
		}

		// We already have this item in the database
		// check for thumbnails,screenshots
		t.makeScreenshots(nil, s.GetHash(t.fileNamingAlgorithm))

		// check for container
		if !s.Format.Valid {
			videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
			if err != nil {
				return logError(err)
			}
			container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)
			logger.Infof("Adding container %s to file %s", container, t.FilePath)

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := scene.UpdateFormat(r.Scene(), s.ID, string(container))
				return err
			}); err != nil {
				return logError(err)
			}
		}

		// check if oshash is set
		if !s.OSHash.Valid {
			logger.Infof("Calculating oshash for existing file %s ...", t.FilePath)
			oshash, err := utils.OSHashFromFilePath(t.FilePath)
			if err != nil {
				return nil
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Scene()
				// check if oshash clashes with existing scene
				dupe, _ := qb.FindByOSHash(oshash)
				if dupe != nil {
					return fmt.Errorf("OSHash for file %s is the same as that of %s", t.FilePath, dupe.Path)
				}

				_, err := scene.UpdateOSHash(qb, s.ID, oshash)
				return err
			}); err != nil {
				return logError(err)
			}
		}

		// check if MD5 is set, if calculateMD5 is true
		if t.calculateMD5 && !s.Checksum.Valid {
			checksum, err := t.calculateChecksum()
			if err != nil {
				return logError(err)
			}

			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Scene()
				// check if checksum clashes with existing scene
				dupe, _ := qb.FindByChecksum(checksum)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", t.FilePath, dupe.Path)
				}

				_, err := scene.UpdateChecksum(qb, s.ID, checksum)
				return err
			}); err != nil {
				return logError(err)
			}
		}

		if s.Interactive != interactive {
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				qb := r.Scene()
				scenePartial := models.ScenePartial{
					ID:          s.ID,
					Interactive: &interactive,
				}
				_, err := qb.Update(scenePartial)
				return err
			}); err != nil {
				return logError(err)
			}
		}

		return nil
	}

	// Ignore directories.
	if isDir, _ := utils.DirExists(t.FilePath); isDir {
		return nil
	}

	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)

	// Override title to be filename if UseFileMetadata is false
	if !t.UseFileMetadata {
		videoFile.SetTitleFromPath(t.StripFileExtension)
	}

	var checksum string

	logger.Infof("%s not found. Calculating oshash...", t.FilePath)
	oshash, err := utils.OSHashFromFilePath(t.FilePath)
	if err != nil {
		return logError(err)
	}

	if t.fileNamingAlgorithm == models.HashAlgorithmMd5 || t.calculateMD5 {
		checksum, err = t.calculateChecksum()
		if err != nil {
			return logError(err)
		}
	}

	// check for scene by checksum and oshash - MD5 should be
	// redundant, but check both
	t.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		qb := r.Scene()
		if checksum != "" {
			s, _ = qb.FindByChecksum(checksum)
		}

		if s == nil {
			s, _ = qb.FindByOSHash(oshash)
		}

		return nil
	})

	sceneHash := oshash

	if t.fileNamingAlgorithm == models.HashAlgorithmMd5 {
		sceneHash = checksum
	}

	t.makeScreenshots(videoFile, sceneHash)

	if s != nil {
		exists, _ := utils.FileExists(s.Path)
		if !t.CaseSensitiveFs {
			// #1426 - if file exists but is a case-insensitive match for the
			// original filename, then treat it as a move
			if exists && strings.EqualFold(t.FilePath, s.Path) {
				exists = false
			}
		}

		if exists {
			logger.Infof("%s already exists. Duplicate of %s", t.FilePath, s.Path)
		} else {
			logger.Infof("%s already exists. Updating path...", t.FilePath)
			scenePartial := models.ScenePartial{
				ID:          s.ID,
				Path:        &t.FilePath,
				Interactive: &interactive,
			}
			if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
				_, err := r.Scene().Update(scenePartial)
				return err
			}); err != nil {
				return logError(err)
			}

			GetInstance().PluginCache.ExecutePostHooks(t.ctx, s.ID, plugin.SceneUpdatePost, nil, nil)
		}
	} else {
		logger.Infof("%s doesn't exist. Creating new item...", t.FilePath)
		currentTime := time.Now()
		newScene := models.Scene{
			Checksum:   sql.NullString{String: checksum, Valid: checksum != ""},
			OSHash:     sql.NullString{String: oshash, Valid: oshash != ""},
			Path:       t.FilePath,
			Title:      sql.NullString{String: videoFile.Title, Valid: true},
			Duration:   sql.NullFloat64{Float64: videoFile.Duration, Valid: true},
			VideoCodec: sql.NullString{String: videoFile.VideoCodec, Valid: true},
			AudioCodec: sql.NullString{String: videoFile.AudioCodec, Valid: true},
			Format:     sql.NullString{String: string(container), Valid: true},
			Width:      sql.NullInt64{Int64: int64(videoFile.Width), Valid: true},
			Height:     sql.NullInt64{Int64: int64(videoFile.Height), Valid: true},
			Framerate:  sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true},
			Bitrate:    sql.NullInt64{Int64: videoFile.Bitrate, Valid: true},
			Size:       sql.NullString{String: strconv.FormatInt(videoFile.Size, 10), Valid: true},
			FileModTime: models.NullSQLiteTimestamp{
				Timestamp: fileModTime,
				Valid:     true,
			},
			CreatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			UpdatedAt:   models.SQLiteTimestamp{Timestamp: currentTime},
			Interactive: interactive,
		}

		if t.UseFileMetadata {
			newScene.Details = sql.NullString{String: videoFile.Comment, Valid: true}
			newScene.Date = models.SQLiteDate{String: videoFile.CreationTime.Format("2006-01-02")}
		}

		if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			var err error
			retScene, err = r.Scene().Create(newScene)
			return err
		}); err != nil {
			return logError(err)
		}

		GetInstance().PluginCache.ExecutePostHooks(t.ctx, retScene.ID, plugin.SceneCreatePost, nil, nil)
	}

	return retScene
}

func (t *ScanTask) rescanScene(s *models.Scene, fileModTime time.Time) (*models.Scene, error) {
	logger.Infof("%s has been updated: rescanning", t.FilePath)

	// update the oshash/checksum and the modification time
	logger.Infof("Calculating oshash for existing file %s ...", t.FilePath)
	oshash, err := utils.OSHashFromFilePath(t.FilePath)
	if err != nil {
		return nil, err
	}

	var checksum *sql.NullString
	if t.calculateMD5 {
		cs, err := t.calculateChecksum()
		if err != nil {
			return nil, err
		}

		checksum = &sql.NullString{
			String: cs,
			Valid:  true,
		}
	}

	// regenerate the file details as well
	videoFile, err := ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)
	if err != nil {
		return nil, err
	}
	container := ffmpeg.MatchContainer(videoFile.Container, t.FilePath)

	currentTime := time.Now()
	scenePartial := models.ScenePartial{
		ID:       s.ID,
		Checksum: checksum,
		OSHash: &sql.NullString{
			String: oshash,
			Valid:  true,
		},
		Duration:   &sql.NullFloat64{Float64: videoFile.Duration, Valid: true},
		VideoCodec: &sql.NullString{String: videoFile.VideoCodec, Valid: true},
		AudioCodec: &sql.NullString{String: videoFile.AudioCodec, Valid: true},
		Format:     &sql.NullString{String: string(container), Valid: true},
		Width:      &sql.NullInt64{Int64: int64(videoFile.Width), Valid: true},
		Height:     &sql.NullInt64{Int64: int64(videoFile.Height), Valid: true},
		Framerate:  &sql.NullFloat64{Float64: videoFile.FrameRate, Valid: true},
		Bitrate:    &sql.NullInt64{Int64: videoFile.Bitrate, Valid: true},
		Size:       &sql.NullString{String: strconv.FormatInt(videoFile.Size, 10), Valid: true},
		FileModTime: &models.NullSQLiteTimestamp{
			Timestamp: fileModTime,
			Valid:     true,
		},
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: currentTime},
	}

	var ret *models.Scene
	if err := t.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		var err error
		ret, err = r.Scene().Update(scenePartial)
		return err
	}); err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	GetInstance().PluginCache.ExecutePostHooks(t.ctx, ret.ID, plugin.SceneUpdatePost, nil, nil)

	// leave the generated files as is - the scene file may have been moved
	// elsewhere

	return ret, nil
}
func (t *ScanTask) makeScreenshots(probeResult *ffmpeg.VideoFile, checksum string) {
	thumbPath := instance.Paths.Scene.GetThumbnailScreenshotPath(checksum)
	normalPath := instance.Paths.Scene.GetScreenshotPath(checksum)

	thumbExists, _ := utils.FileExists(thumbPath)
	normalExists, _ := utils.FileExists(normalPath)

	if thumbExists && normalExists {
		return
	}

	if probeResult == nil {
		var err error
		probeResult, err = ffmpeg.NewVideoFile(instance.FFProbePath, t.FilePath, t.StripFileExtension)

		if err != nil {
			logger.Error(err.Error())
			return
		}
		logger.Infof("Regenerating images for %s", t.FilePath)
	}

	at := float64(probeResult.Duration) * 0.2

	if !thumbExists {
		logger.Debugf("Creating thumbnail for %s", t.FilePath)
		makeScreenshot(*probeResult, thumbPath, 5, 320, at)
	}

	if !normalExists {
		logger.Debugf("Creating screenshot for %s", t.FilePath)
		makeScreenshot(*probeResult, normalPath, 2, probeResult.Width, at)
	}
}

func (t *ScanTask) getInteractive() bool {
	_, err := os.Stat(utils.GetFunscriptPath(t.FilePath))
	return err == nil
}
