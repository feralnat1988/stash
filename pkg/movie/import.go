package movie

import (
	"context"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type ImporterReaderWriter interface {
	models.MovieCreatorUpdater
	FindByName(ctx context.Context, name string, nocase bool) (*models.Movie, error)
}

type Importer struct {
	ReaderWriter        ImporterReaderWriter
	StudioWriter        models.StudioFinderCreator
	TagWriter           models.TagFinderCreator
	Input               jsonschema.Movie
	MissingRefBehaviour models.ImportMissingRefEnum

	movie          models.Movie
	frontImageData []byte
	backImageData  []byte
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.movie = i.movieJSONToMovie(i.Input)

	if err := i.populateStudio(ctx); err != nil {
		return err
	}

	if err := i.populateTags(ctx); err != nil {
		return err
	}

	var err error
	if len(i.Input.FrontImage) > 0 {
		i.frontImageData, err = utils.ProcessBase64Image(i.Input.FrontImage)
		if err != nil {
			return fmt.Errorf("invalid front_image: %v", err)
		}
	}
	if len(i.Input.BackImage) > 0 {
		i.backImageData, err = utils.ProcessBase64Image(i.Input.BackImage)
		if err != nil {
			return fmt.Errorf("invalid back_image: %v", err)
		}
	}

	return nil
}

func (i *Importer) populateTags(ctx context.Context) error {
	if len(i.Input.Tags) > 0 {

		tags, err := importTags(ctx, i.TagWriter, i.Input.Tags, i.MissingRefBehaviour)
		if err != nil {
			return err
		}

		for _, p := range tags {
			i.movie.TagIDs.Add(p.ID)
		}
	}

	return nil
}

func importTags(ctx context.Context, tagWriter models.TagFinderCreator, names []string, missingRefBehaviour models.ImportMissingRefEnum) ([]*models.Tag, error) {
	tags, err := tagWriter.FindByNames(ctx, names, false)
	if err != nil {
		return nil, err
	}

	var pluckedNames []string
	for _, tag := range tags {
		pluckedNames = append(pluckedNames, tag.Name)
	}

	missingTags := sliceutil.Filter(names, func(name string) bool {
		return !sliceutil.Contains(pluckedNames, name)
	})

	if len(missingTags) > 0 {
		if missingRefBehaviour == models.ImportMissingRefEnumFail {
			return nil, fmt.Errorf("tags [%s] not found", strings.Join(missingTags, ", "))
		}

		if missingRefBehaviour == models.ImportMissingRefEnumCreate {
			createdTags, err := createTags(ctx, tagWriter, missingTags)
			if err != nil {
				return nil, fmt.Errorf("error creating tags: %v", err)
			}

			tags = append(tags, createdTags...)
		}

		// ignore if MissingRefBehaviour set to Ignore
	}

	return tags, nil
}

func createTags(ctx context.Context, tagWriter models.TagFinderCreator, names []string) ([]*models.Tag, error) {
	var ret []*models.Tag
	for _, name := range names {
		newTag := models.NewTag()
		newTag.Name = name

		err := tagWriter.Create(ctx, &newTag)
		if err != nil {
			return nil, err
		}

		ret = append(ret, &newTag)
	}

	return ret, nil
}

func (i *Importer) movieJSONToMovie(movieJSON jsonschema.Movie) models.Movie {
	newMovie := models.Movie{
		Name:      movieJSON.Name,
		Aliases:   movieJSON.Aliases,
		Director:  movieJSON.Director,
		Synopsis:  movieJSON.Synopsis,
		CreatedAt: movieJSON.CreatedAt.GetTime(),
		UpdatedAt: movieJSON.UpdatedAt.GetTime(),

		TagIDs: models.NewRelatedIDs([]int{}),
	}

	if len(movieJSON.URLs) > 0 {
		newMovie.URLs = models.NewRelatedStrings(movieJSON.URLs)
	} else if movieJSON.URL != "" {
		newMovie.URLs = models.NewRelatedStrings([]string{movieJSON.URL})
	}
	if movieJSON.Date != "" {
		d, err := models.ParseDate(movieJSON.Date)
		if err == nil {
			newMovie.Date = &d
		}
	}
	if movieJSON.Rating != 0 {
		newMovie.Rating = &movieJSON.Rating
	}

	if movieJSON.Duration != 0 {
		newMovie.Duration = &movieJSON.Duration
	}

	return newMovie
}

func (i *Importer) populateStudio(ctx context.Context) error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(ctx, i.Input.Studio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %v", err)
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("movie studio '%s' not found", i.Input.Studio)
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				studioID, err := i.createStudio(ctx, i.Input.Studio)
				if err != nil {
					return err
				}
				i.movie.StudioID = &studioID
			}
		} else {
			i.movie.StudioID = &studio.ID
		}
	}

	return nil
}

func (i *Importer) createStudio(ctx context.Context, name string) (int, error) {
	newStudio := models.NewStudio()
	newStudio.Name = name

	err := i.StudioWriter.Create(ctx, &newStudio)
	if err != nil {
		return 0, err
	}

	return newStudio.ID, nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	if len(i.frontImageData) > 0 {
		if err := i.ReaderWriter.UpdateFrontImage(ctx, id, i.frontImageData); err != nil {
			return fmt.Errorf("error setting movie front image: %v", err)
		}
	}

	if len(i.backImageData) > 0 {
		if err := i.ReaderWriter.UpdateBackImage(ctx, id, i.backImageData); err != nil {
			return fmt.Errorf("error setting movie back image: %v", err)
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	const nocase = false
	existing, err := i.ReaderWriter.FindByName(ctx, i.Name(), nocase)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := existing.ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	err := i.ReaderWriter.Create(ctx, &i.movie)
	if err != nil {
		return nil, fmt.Errorf("error creating movie: %v", err)
	}

	id := i.movie.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	movie := i.movie
	movie.ID = id
	err := i.ReaderWriter.Update(ctx, &movie)
	if err != nil {
		return fmt.Errorf("error updating existing movie: %v", err)
	}

	return nil
}
