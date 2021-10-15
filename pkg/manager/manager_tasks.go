package manager

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/remeh/sizedwaitgroup"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func isGallery(pathname string) bool {
	gExt := config.GetInstance().GetGalleryExtensions()
	return utils.MatchExtension(pathname, gExt)
}

func isVideo(pathname string) bool {
	vidExt := config.GetInstance().GetVideoExtensions()
	return utils.MatchExtension(pathname, vidExt)
}

func isImage(pathname string) bool {
	imgExt := config.GetInstance().GetImageExtensions()
	return utils.MatchExtension(pathname, imgExt)
}

func getScanPaths(inputPaths []string) []*models.StashConfig {
	if len(inputPaths) == 0 {
		return config.GetInstance().GetStashPaths()
	}

	var ret []*models.StashConfig
	for _, p := range inputPaths {
		s := getStashFromDirPath(p)
		if s == nil {
			logger.Warnf("%s is not in the configured stash paths", p)
			continue
		}

		// make a copy, changing the path
		ss := *s
		ss.Path = p
		ret = append(ret, &ss)
	}

	return ret
}

// ScanSubscribe subscribes to a notification that is triggered when a
// scan or clean is complete.
func (s *singleton) ScanSubscribe(ctx context.Context) <-chan bool {
	return s.scanSubs.subscribe(ctx)
}

func (s *singleton) Scan(ctx context.Context, input models.ScanMetadataInput) (int, error) {
	if err := s.validateFFMPEG(); err != nil {
		return 0, err
	}

	scanJob := ScanJob{
		txnManager:    s.TxnManager,
		input:         input,
		subscriptions: s.scanSubs,
	}

	return s.JobManager.Add(ctx, "Scanning...", &scanJob), nil
}

func (s *singleton) Import(ctx context.Context) (int, error) {
	config := config.GetInstance()
	metadataPath := config.GetMetadataPath()
	if metadataPath == "" {
		return 0, errors.New("metadata path must be set in config")
	}

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		task := ImportTask{
			txnManager:          s.TxnManager,
			BaseDir:             metadataPath,
			Reset:               true,
			DuplicateBehaviour:  models.ImportDuplicateEnumFail,
			MissingRefBehaviour: models.ImportMissingRefEnumFail,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}
		task.Start(ctx)
	})

	return s.JobManager.Add(ctx, "Importing...", j), nil
}

func (s *singleton) Export(ctx context.Context) (int, error) {
	config := config.GetInstance()
	metadataPath := config.GetMetadataPath()
	if metadataPath == "" {
		return 0, errors.New("metadata path must be set in config")
	}

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		var wg sync.WaitGroup
		wg.Add(1)
		task := ExportTask{
			txnManager:          s.TxnManager,
			full:                true,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}
		task.Start(&wg)
	})

	return s.JobManager.Add(ctx, "Exporting...", j), nil
}

func (s *singleton) RunSingleTask(ctx context.Context, t Task) int {
	var wg sync.WaitGroup
	wg.Add(1)

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		t.Start(ctx)
		wg.Done()
	})

	return s.JobManager.Add(ctx, t.GetDescription(), j)
}

func setGeneratePreviewOptionsInput(optionsInput *models.GeneratePreviewOptionsInput) {
	config := config.GetInstance()
	if optionsInput.PreviewSegments == nil {
		val := config.GetPreviewSegments()
		optionsInput.PreviewSegments = &val
	}

	if optionsInput.PreviewSegmentDuration == nil {
		val := config.GetPreviewSegmentDuration()
		optionsInput.PreviewSegmentDuration = &val
	}

	if optionsInput.PreviewExcludeStart == nil {
		val := config.GetPreviewExcludeStart()
		optionsInput.PreviewExcludeStart = &val
	}

	if optionsInput.PreviewExcludeEnd == nil {
		val := config.GetPreviewExcludeEnd()
		optionsInput.PreviewExcludeEnd = &val
	}

	if optionsInput.PreviewPreset == nil {
		val := config.GetPreviewPreset()
		optionsInput.PreviewPreset = &val
	}
}

func (s *singleton) Generate(ctx context.Context, input models.GenerateMetadataInput) (int, error) {
	if err := s.validateFFMPEG(); err != nil {
		return 0, err
	}
	if err := instance.Paths.Generated.EnsureTmpDir(); err != nil {
		logger.Warnf("could not generate temporary directory: %v", err)
	}

	sceneIDs, err := utils.StringSliceToIntSlice(input.SceneIDs)
	if err != nil {
		logger.Error(err.Error())
	}
	markerIDs, err := utils.StringSliceToIntSlice(input.MarkerIDs)
	if err != nil {
		logger.Error(err.Error())
	}

	// TODO - formalise this
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		var scenes []*models.Scene
		var err error
		var markers []*models.SceneMarker

		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			qb := r.Scene()
			if len(sceneIDs) > 0 {
				scenes, err = qb.FindMany(sceneIDs)
			} else {
				scenes, err = qb.All()
			}

			if err != nil {
				return err
			}

			if len(markerIDs) > 0 {
				markers, err = r.SceneMarker().FindMany(markerIDs)
				if err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		config := config.GetInstance()
		parallelTasks := config.GetParallelTasksWithAutoDetection()

		logger.Infof("Generate started with %d parallel tasks", parallelTasks)
		wg := sizedwaitgroup.New(parallelTasks)

		lenScenes := len(scenes)
		total := lenScenes + len(markers)
		progress.SetTotal(total)

		if job.IsCancelled(ctx) {
			logger.Info("Stopping due to user request")
			return
		}

		// TODO - consider removing this. Even though we're only waiting a maximum of
		// 90 seconds for this, it is all for a simple log message, and probably not worth
		// waiting for
		var totalsNeeded *totalsGenerate
		progress.ExecuteTask("Calculating content to generate...", func() {
			totalsNeeded = s.neededGenerate(scenes, input)

			if totalsNeeded == nil {
				logger.Infof("Taking too long to count content. Skipping...")
				logger.Infof("Generating content")
			} else {
				logger.Infof("Generating %d sprites %d previews %d image previews %d markers %d transcodes %d phashes", totalsNeeded.sprites, totalsNeeded.previews, totalsNeeded.imagePreviews, totalsNeeded.markers, totalsNeeded.transcodes, totalsNeeded.phashes)
			}
		})

		fileNamingAlgo := config.GetVideoFileNamingAlgorithm()

		overwrite := false
		if input.Overwrite != nil {
			overwrite = *input.Overwrite
		}

		generatePreviewOptions := input.PreviewOptions
		if generatePreviewOptions == nil {
			generatePreviewOptions = &models.GeneratePreviewOptionsInput{}
		}
		setGeneratePreviewOptionsInput(generatePreviewOptions)

		// Start measuring how long the generate has taken. (consider moving this up)
		start := time.Now()
		if err = instance.Paths.Generated.EnsureTmpDir(); err != nil {
			logger.Warnf("could not create temporary directory: %v", err)
		}

		for _, scene := range scenes {
			progress.Increment()
			if job.IsCancelled(ctx) {
				logger.Info("Stopping due to user request")
				wg.Wait()
				if err := instance.Paths.Generated.EmptyTmpDir(); err != nil {
					logger.Warnf("failure emptying temporary directory: %v", err)
				}
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping generate")
				continue
			}

			if utils.IsTrue(input.Sprites) {
				task := GenerateSpriteTask{
					Scene:               *scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				wg.Add()
				go progress.ExecuteTask(fmt.Sprintf("Generating sprites for %s", scene.Path), func() {
					task.Start()
					wg.Done()
				})
			}

			if utils.IsTrue(input.Previews) {
				task := GeneratePreviewTask{
					Scene:               *scene,
					ImagePreview:        utils.IsTrue(input.ImagePreviews),
					Options:             *generatePreviewOptions,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				wg.Add()
				go progress.ExecuteTask(fmt.Sprintf("Generating preview for %s", scene.Path), func() {
					task.Start()
					wg.Done()
				})
			}

			if utils.IsTrue(input.Markers) {
				wg.Add()
				task := GenerateMarkersTask{
					TxnManager:          s.TxnManager,
					Scene:               scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
					ImagePreview:        utils.IsTrue(input.MarkerImagePreviews),
					Screenshot:          utils.IsTrue(input.MarkerScreenshots),
				}
				go progress.ExecuteTask(fmt.Sprintf("Generating markers for %s", scene.Path), func() {
					task.Start()
					wg.Done()
				})
			}

			if utils.IsTrue(input.Transcodes) {
				wg.Add()
				task := GenerateTranscodeTask{
					Scene:               *scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				go progress.ExecuteTask(fmt.Sprintf("Generating transcode for %s", scene.Path), func() {
					task.Start()
					wg.Done()
				})
			}

			if utils.IsTrue(input.Phashes) {
				task := GeneratePhashTask{
					Scene:               *scene,
					fileNamingAlgorithm: fileNamingAlgo,
					txnManager:          s.TxnManager,
					Overwrite:           overwrite,
				}
				wg.Add()
				go progress.ExecuteTask(fmt.Sprintf("Generating phash for %s", scene.Path), func() {
					task.Start()
					wg.Done()
				})
			}
		}

		wg.Wait()

		for _, marker := range markers {
			progress.Increment()
			if job.IsCancelled(ctx) {
				logger.Info("Stopping due to user request")
				wg.Wait()
				if err := instance.Paths.Generated.EmptyTmpDir(); err != nil {
					logger.Warnf("failure emptying temporary directory: %v", err)
				}
				elapsed := time.Since(start)
				logger.Info(fmt.Sprintf("Generate finished (%s)", elapsed))
				return
			}

			if marker == nil {
				logger.Errorf("nil marker, skipping generate")
				continue
			}

			wg.Add()
			task := GenerateMarkersTask{
				TxnManager:          s.TxnManager,
				Marker:              marker,
				Overwrite:           overwrite,
				fileNamingAlgorithm: fileNamingAlgo,
			}
			go progress.ExecuteTask(fmt.Sprintf("Generating marker preview for marker ID %d", marker.ID), func() {
				task.Start()
				wg.Done()
			})
		}

		wg.Wait()

		if err = instance.Paths.Generated.EmptyTmpDir(); err != nil {
			logger.Warnf("failure emptying temporary directory: %v", err)
		}
		elapsed := time.Since(start)
		logger.Info(fmt.Sprintf("Generate finished (%s)", elapsed))
	})

	return s.JobManager.Add(ctx, "Generating...", j), nil
}

func (s *singleton) GenerateDefaultScreenshot(ctx context.Context, sceneId string) int {
	return s.generateScreenshot(ctx, sceneId, nil)
}

func (s *singleton) GenerateScreenshot(ctx context.Context, sceneId string, at float64) int {
	return s.generateScreenshot(ctx, sceneId, &at)
}

// generate default screenshot if at is nil
func (s *singleton) generateScreenshot(ctx context.Context, sceneId string, at *float64) int {
	if err := instance.Paths.Generated.EnsureTmpDir(); err != nil {
		logger.Warnf("failure generating screenshot: %v", err)
	}

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		sceneIdInt, err := strconv.Atoi(sceneId)
		if err != nil {
			logger.Errorf("Error parsing scene id %s: %s", sceneId, err.Error())
			return
		}

		var scene *models.Scene
		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			scene, err = r.Scene().Find(sceneIdInt)
			return err
		}); err != nil || scene == nil {
			logger.Errorf("failed to get scene for generate: %s", err.Error())
			return
		}

		task := GenerateScreenshotTask{
			txnManager:          s.TxnManager,
			Scene:               *scene,
			ScreenshotAt:        at,
			fileNamingAlgorithm: config.GetInstance().GetVideoFileNamingAlgorithm(),
		}

		task.Start()

		logger.Infof("Generate screenshot finished")
	})

	return s.JobManager.Add(ctx, fmt.Sprintf("Generating screenshot for scene id %s", sceneId), j)
}

func (s *singleton) AutoTag(ctx context.Context, input models.AutoTagMetadataInput) int {
	j := autoTagJob{
		txnManager: s.TxnManager,
		input:      input,
	}

	return s.JobManager.Add(ctx, "Auto-tagging...", &j)
}

func (s *singleton) Clean(ctx context.Context, input models.CleanMetadataInput) int {
	j := cleanJob{
		txnManager: s.TxnManager,
		input:      input,
		scanSubs:   s.scanSubs,
	}

	return s.JobManager.Add(ctx, "Cleaning...", &j)
}

func (s *singleton) MigrateHash(ctx context.Context) int {
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
		logger.Infof("Migrating generated files for %s naming hash", fileNamingAlgo.String())

		var scenes []*models.Scene
		if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
			var err error
			scenes, err = r.Scene().All()
			return err
		}); err != nil {
			logger.Errorf("failed to fetch list of scenes for migration: %s", err.Error())
			return
		}

		var wg sync.WaitGroup
		total := len(scenes)
		progress.SetTotal(total)

		for _, scene := range scenes {
			progress.Increment()
			if job.IsCancelled(ctx) {
				logger.Info("Stopping due to user request")
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping migrate")
				continue
			}

			wg.Add(1)

			task := MigrateHashTask{Scene: scene, fileNamingAlgorithm: fileNamingAlgo}
			go func() {
				task.Start()
				wg.Done()
			}()

			wg.Wait()
		}

		logger.Info("Finished migrating")
	})

	return s.JobManager.Add(ctx, "Migrating scene hashes...", j)
}

type totalsGenerate struct {
	sprites       int64
	previews      int64
	imagePreviews int64
	markers       int64
	transcodes    int64
	phashes       int64
}

func (s *singleton) neededGenerate(scenes []*models.Scene, input models.GenerateMetadataInput) *totalsGenerate {

	var totals totalsGenerate
	const timeout = 90 * time.Second

	// create a control channel through which to signal the counting loop when the timeout is reached
	chTimeout := make(chan struct{})

	//run the timeout function in a separate thread
	go func() {
		time.Sleep(timeout)
		chTimeout <- struct{}{}
	}()

	fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
	overwrite := false
	if input.Overwrite != nil {
		overwrite = *input.Overwrite
	}

	logger.Infof("Counting content to generate...")
	for _, scene := range scenes {
		if scene != nil {
			if utils.IsTrue(input.Sprites) {
				task := GenerateSpriteTask{
					Scene:               *scene,
					fileNamingAlgorithm: fileNamingAlgo,
				}

				if overwrite || task.required() {
					totals.sprites++
				}
			}

			if utils.IsTrue(input.Previews) {
				task := GeneratePreviewTask{
					Scene:               *scene,
					ImagePreview:        utils.IsTrue(input.ImagePreviews),
					fileNamingAlgorithm: fileNamingAlgo,
				}

				sceneHash := scene.GetHash(task.fileNamingAlgorithm)
				if overwrite || !task.doesVideoPreviewExist(sceneHash) {
					totals.previews++
				}

				if utils.IsTrue(input.ImagePreviews) && (overwrite || !task.doesImagePreviewExist(sceneHash)) {
					totals.imagePreviews++
				}
			}

			if utils.IsTrue(input.Markers) {
				task := GenerateMarkersTask{
					TxnManager:          s.TxnManager,
					Scene:               scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				totals.markers += int64(task.isMarkerNeeded())
			}

			if utils.IsTrue(input.Transcodes) {
				task := GenerateTranscodeTask{
					Scene:               *scene,
					Overwrite:           overwrite,
					fileNamingAlgorithm: fileNamingAlgo,
				}
				if task.isTranscodeNeeded() {
					totals.transcodes++
				}
			}

			if utils.IsTrue(input.Phashes) {
				task := GeneratePhashTask{
					Scene:               *scene,
					fileNamingAlgorithm: fileNamingAlgo,
				}

				if task.shouldGenerate() {
					totals.phashes++
				}
			}
		}
		//check for timeout
		select {
		case <-chTimeout:
			return nil
		default:
		}

	}
	return &totals
}

func (s *singleton) StashBoxBatchPerformerTag(ctx context.Context, input models.StashBoxBatchPerformerTagInput) int {
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) {
		logger.Infof("Initiating stash-box batch performer tag")

		boxes := config.GetInstance().GetStashBoxes()
		if input.Endpoint < 0 || input.Endpoint >= len(boxes) {
			logger.Error(fmt.Errorf("invalid stash_box_index %d", input.Endpoint))
			return
		}
		box := boxes[input.Endpoint]

		var tasks []StashBoxPerformerTagTask

		if len(input.PerformerIds) > 0 {
			if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
				performerQuery := r.Performer()

				for _, performerID := range input.PerformerIds {
					if id, err := strconv.Atoi(performerID); err == nil {
						performer, err := performerQuery.Find(id)
						if err == nil {
							tasks = append(tasks, StashBoxPerformerTagTask{
								txnManager:      s.TxnManager,
								performer:       performer,
								refresh:         input.Refresh,
								box:             box,
								excluded_fields: input.ExcludeFields,
							})
						} else {
							return err
						}
					}
				}
				return nil
			}); err != nil {
				logger.Error(err.Error())
			}
		} else if len(input.PerformerNames) > 0 {
			for i := range input.PerformerNames {
				if len(input.PerformerNames[i]) > 0 {
					tasks = append(tasks, StashBoxPerformerTagTask{
						txnManager:      s.TxnManager,
						name:            &input.PerformerNames[i],
						refresh:         input.Refresh,
						box:             box,
						excluded_fields: input.ExcludeFields,
					})
				}
			}
		} else {
			if err := s.TxnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
				performerQuery := r.Performer()
				var performers []*models.Performer
				var err error
				if input.Refresh {
					performers, err = performerQuery.FindByStashIDStatus(true, box.Endpoint)
				} else {
					performers, err = performerQuery.FindByStashIDStatus(false, box.Endpoint)
				}
				if err != nil {
					return fmt.Errorf("error querying performers: %v", err)
				}

				for _, performer := range performers {
					tasks = append(tasks, StashBoxPerformerTagTask{
						txnManager:      s.TxnManager,
						performer:       performer,
						refresh:         input.Refresh,
						box:             box,
						excluded_fields: input.ExcludeFields,
					})
				}
				return nil
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		if len(tasks) == 0 {
			return
		}

		progress.SetTotal(len(tasks))

		logger.Infof("Starting stash-box batch operation for %d performers", len(tasks))

		var wg sync.WaitGroup
		for _, task := range tasks {
			wg.Add(1)
			progress.ExecuteTask(task.Description(), func() {
				task.Start()
				wg.Done()
			})

			progress.Increment()
		}
	})

	return s.JobManager.Add(ctx, "Batch stash-box performer tag...", j)
}
