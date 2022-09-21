package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
)

type PartialUpdater interface {
	UpdatePartial(ctx context.Context, id int, updatedGallery models.GalleryPartial) (*models.Gallery, error)
}

type ImageUpdater interface {
	GetImageIDs(ctx context.Context, galleryID int) ([]int, error)
	UpdateImages(ctx context.Context, galleryID int, imageIDs []int) error
}

func AddImage(ctx context.Context, qb ImageUpdater, galleryID int, imageID int) error {
	imageIDs, err := qb.GetImageIDs(ctx, galleryID)
	if err != nil {
		return err
	}

	imageIDs = intslice.IntAppendUnique(imageIDs, imageID)
	return qb.UpdateImages(ctx, galleryID, imageIDs)
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
