package scraper

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
)

// autoTagScraperID is the scraper ID for the built-in AutoTag scraper
const (
	autoTagScraperID   = "builtin_autotag"
	autoTagScraperName = "Auto Tag"
)

type autotagScraper struct {
	txnManager   models.TransactionManager
	globalConfig GlobalConfig
}

func autotagMatchPerformers(path string, performerReader models.PerformerReader) ([]*models.ScrapedPerformer, error) {
	p, err := match.PathToPerformers(path, performerReader)
	if err != nil {
		return nil, fmt.Errorf("error matching performers: %w", err)
	}

	var ret []*models.ScrapedPerformer
	for _, pp := range p {
		id := strconv.Itoa(pp.ID)

		sp := &models.ScrapedPerformer{
			Name:     &pp.Name.String,
			StoredID: &id,
		}
		if pp.Gender.Valid {
			sp.Gender = &pp.Gender.String
		}

		ret = append(ret, sp)
	}

	return ret, nil
}

func autotagMatchStudio(path string, studioReader models.StudioReader) (*models.ScrapedStudio, error) {
	st, err := match.PathToStudios(path, studioReader)
	if err != nil {
		return nil, fmt.Errorf("error matching studios: %w", err)
	}

	if len(st) > 0 {
		id := strconv.Itoa(st[0].ID)
		return &models.ScrapedStudio{
			Name:     st[0].Name.String,
			StoredID: &id,
		}, nil
	}

	return nil, nil
}

func autotagMatchTags(path string, tagReader models.TagReader) ([]*models.ScrapedTag, error) {
	t, err := match.PathToTags(path, tagReader)
	if err != nil {
		return nil, fmt.Errorf("error matching tags: %w", err)
	}

	var ret []*models.ScrapedTag
	for _, tt := range t {
		id := strconv.Itoa(tt.ID)

		st := &models.ScrapedTag{
			Name:     tt.Name,
			StoredID: &id,
		}

		ret = append(ret, st)
	}

	return ret, nil
}

func (s autotagScraper) viaScene(ctx context.Context, _client *http.Client, scene *models.Scene) (*models.ScrapedScene, error) {
	var ret *models.ScrapedScene

	// populate performers, studio and tags based on scene path
	if err := s.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		path := scene.Path
		performers, err := autotagMatchPerformers(path, r.Performer())
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}
		studio, err := autotagMatchStudio(path, r.Studio())
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}

		tags, err := autotagMatchTags(path, r.Tag())
		if err != nil {
			return fmt.Errorf("autotag scraper viaScene: %w", err)
		}

		if len(performers) > 0 || studio != nil || len(tags) > 0 {
			ret = &models.ScrapedScene{
				Performers: performers,
				Studio:     studio,
				Tags:       tags,
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (s autotagScraper) viaGallery(ctx context.Context, _client *http.Client, gallery *models.Gallery) (*models.ScrapedGallery, error) {
	if !gallery.Path.Valid {
		// not valid for non-path-based galleries
		return nil, nil
	}

	var ret *models.ScrapedGallery

	// populate performers, studio and tags based on scene path
	if err := s.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		path := gallery.Path.String
		performers, err := autotagMatchPerformers(path, r.Performer())
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}
		studio, err := autotagMatchStudio(path, r.Studio())
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}

		tags, err := autotagMatchTags(path, r.Tag())
		if err != nil {
			return fmt.Errorf("autotag scraper viaGallery: %w", err)
		}

		if len(performers) > 0 || studio != nil || len(tags) > 0 {
			ret = &models.ScrapedGallery{
				Performers: performers,
				Studio:     studio,
				Tags:       tags,
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (s autotagScraper) supports(ty models.ScrapeContentType) bool {
	switch ty {
	case models.ScrapeContentTypeScene:
		return true
	case models.ScrapeContentTypeGallery:
		return true
	}

	return false
}

func (s autotagScraper) supportsURL(url string, ty models.ScrapeContentType) bool {
	return false
}

func (s autotagScraper) spec() models.Scraper {
	supportedScrapes := []models.ScrapeType{
		models.ScrapeTypeFragment,
	}

	return models.Scraper{
		ID:   autoTagScraperID,
		Name: autoTagScraperName,
		Scene: &models.ScraperSpec{
			SupportedScrapes: supportedScrapes,
		},
		Gallery: &models.ScraperSpec{
			SupportedScrapes: supportedScrapes,
		},
	}
}

func getAutoTagScraper(txnManager models.TransactionManager, globalConfig GlobalConfig) scraper {
	base := autotagScraper{
		txnManager:   txnManager,
		globalConfig: globalConfig,
	}

	return base
}
