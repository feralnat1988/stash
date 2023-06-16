package gallery

import (
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
)

type ImageUpdater interface {
	GetImageIDs(ctx context.Context, galleryID int) ([]int, error)
	AddImages(ctx context.Context, galleryID int, imageIDs ...int) error
	RemoveImages(ctx context.Context, galleryID int, imageIDs ...int) error
}

func (s *Service) Updated(ctx context.Context, galleryID int) error {
	_, err := s.Repository.UpdatePartial(ctx, galleryID, models.GalleryPartial{
		UpdatedAt: models.NewOptionalTime(time.Now()),
	})
	return err
}

// AddImages adds images to the provided gallery.
// It returns an error if the gallery does not support adding images, or if
// the operation fails.
func (s *Service) AddImages(ctx context.Context, g *models.Gallery, toAdd ...int) error {
	if err := validateContentChange(g); err != nil {
		return err
	}

	if err := s.Repository.AddImages(ctx, g.ID, toAdd...); err != nil {
		return fmt.Errorf("failed to add images to gallery: %w", err)
	}

	// #3759 - update the gallery's UpdatedAt timestamp
	return s.Updated(ctx, g.ID)
}

// RemoveImages removes images from the provided gallery.
// It does not validate if the images are part of the gallery.
// It returns an error if the gallery does not support removing images, or if
// the operation fails.
func (s *Service) RemoveImages(ctx context.Context, g *models.Gallery, toRemove ...int) error {
	if err := validateContentChange(g); err != nil {
		return err
	}

	if err := s.Repository.RemoveImages(ctx, g.ID, toRemove...); err != nil {
		return fmt.Errorf("failed to remove images from gallery: %w", err)
	}

	// #3759 - update the gallery's UpdatedAt timestamp
	return s.Updated(ctx, g.ID)
}

func AddPerformer(ctx context.Context, qb PartialUpdater, o *models.Gallery, performerID int) error {
	_, err := qb.UpdatePartial(ctx, o.ID, models.GalleryPartial{
		PerformerIDs: &models.UpdateIDs{
			IDs:  []int{performerID},
			Mode: models.RelationshipUpdateModeAdd,
		},
	})
	return err
}

func AddTag(ctx context.Context, qb PartialUpdater, o *models.Gallery, tagID int) error {
	_, err := qb.UpdatePartial(ctx, o.ID, models.GalleryPartial{
		TagIDs: &models.UpdateIDs{
			IDs:  []int{tagID},
			Mode: models.RelationshipUpdateModeAdd,
		},
	})
	return err
}
