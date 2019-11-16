package manager

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type AutoTagPerformerTask struct {
	performer *models.Performer
}

func (t *AutoTagPerformerTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagPerformer()
}

func getQueryRegex(name string) string {
	const separator = `[.\-_ ]`
	ret := strings.Replace(name, " ", separator+"*", -1)
	ret = "(?:^|" + separator + "+)" + ret + "(?:$|" + separator + "+)"
	return ret
}

func (t *AutoTagPerformerTask) autoTagPerformer() {
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	regex := getQueryRegex(t.performer.Name.String)
	filter := models.FindFilterType{
		Q: &regex,
	}

	logger.Infof("Using regex '%s' to search for scenes", regex)
	scenes, _ := qb.QueryByPathRegex(&filter)

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, scene := range scenes {
		added, err := jqb.AddPerformerScene(scene.ID, t.performer.ID, tx)

		if err != nil {
			logger.Infof("Error adding performer '%s' to scene '%s': %s", t.performer.Name.String, scene.GetTitle(), err.Error())
			tx.Rollback()
			return
		}

		if added {
			logger.Infof("Added performer '%s' to scene '%s'", t.performer.Name.String, scene.GetTitle())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding performer to scene: %s", err.Error())
		return
	}
}

type AutoTagStudioTask struct {
	studio *models.Studio
}

func (t *AutoTagStudioTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagStudio()
}

func (t *AutoTagStudioTask) autoTagStudio() {
	qb := models.NewSceneQueryBuilder()

	regex := getQueryRegex(t.studio.Name.String)
	filter := models.FindFilterType{
		Q: &regex,
	}

	scenes, _ := qb.QueryByPathRegex(&filter)

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, scene := range scenes {
		if scene.StudioID.Int64 == int64(t.studio.ID) {
			// don't modify
			continue
		}

		logger.Infof("Adding studio '%s' to scene '%s'", t.studio.Name.String, scene.GetTitle())

		// set the studio id
		studioID := sql.NullInt64{Int64: int64(t.studio.ID), Valid: true}
		scenePartial := models.ScenePartial{
			ID:       scene.ID,
			StudioID: &studioID,
		}

		_, err := qb.Update(scenePartial, tx)

		if err != nil {
			logger.Infof("Error adding studio to scene: %s", err.Error())
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding studio to scene: %s", err.Error())
		return
	}
}

type AutoTagTagTask struct {
	tag *models.Tag
}

func (t *AutoTagTagTask) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	t.autoTagTag()
}

func (t *AutoTagTagTask) autoTagTag() {
	qb := models.NewSceneQueryBuilder()
	jqb := models.NewJoinsQueryBuilder()

	regex := getQueryRegex(t.tag.Name)
	filter := models.FindFilterType{
		Q: &regex,
	}

	scenes, _ := qb.QueryByPathRegex(&filter)

	ctx := context.TODO()
	tx := database.DB.MustBeginTx(ctx, nil)

	for _, scene := range scenes {
		added, err := jqb.AddSceneTag(scene.ID, t.tag.ID, tx)

		if err != nil {
			logger.Infof("Error adding tag '%s' to scene '%s': %s", t.tag.Name, scene.GetTitle(), err.Error())
			tx.Rollback()
			return
		}

		if added {
			logger.Infof("Added tag '%s' to scene '%s'", t.tag.Name, scene.GetTitle())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Infof("Error adding tag to scene: %s", err.Error())
		return
	}
}
