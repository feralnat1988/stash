// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// ImageReaderWriter is an autogenerated mock type for the ImageReaderWriter type
type ImageReaderWriter struct {
	mock.Mock
}

// AddFileID provides a mock function with given fields: ctx, id, fileID
func (_m *ImageReaderWriter) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	ret := _m.Called(ctx, id, fileID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, models.FileID) error); ok {
		r0 = rf(ctx, id, fileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// All provides a mock function with given fields: ctx
func (_m *ImageReaderWriter) All(ctx context.Context) ([]*models.Image, error) {
	ret := _m.Called(ctx)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context) []*models.Image); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields: ctx
func (_m *ImageReaderWriter) Count(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByFileID provides a mock function with given fields: ctx, fileID
func (_m *ImageReaderWriter) CountByFileID(ctx context.Context, fileID models.FileID) (int, error) {
	ret := _m.Called(ctx, fileID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, models.FileID) int); ok {
		r0 = rf(ctx, fileID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FileID) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByGalleryID provides a mock function with given fields: ctx, galleryID
func (_m *ImageReaderWriter) CountByGalleryID(ctx context.Context, galleryID int) (int, error) {
	ret := _m.Called(ctx, galleryID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, galleryID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, galleryID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, newImage
func (_m *ImageReaderWriter) Create(ctx context.Context, newImage *models.ImageCreateInput) error {
	ret := _m.Called(ctx, newImage)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.ImageCreateInput) error); ok {
		r0 = rf(ctx, newImage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DecrementOCounter provides a mock function with given fields: ctx, id
func (_m *ImageReaderWriter) DecrementOCounter(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: ctx, id
func (_m *ImageReaderWriter) Destroy(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: ctx, id
func (_m *ImageReaderWriter) Find(ctx context.Context, id int) (*models.Image, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Image
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Image); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByChecksum provides a mock function with given fields: ctx, checksum
func (_m *ImageReaderWriter) FindByChecksum(ctx context.Context, checksum string) ([]*models.Image, error) {
	ret := _m.Called(ctx, checksum)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, string) []*models.Image); ok {
		r0 = rf(ctx, checksum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, checksum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByFileID provides a mock function with given fields: ctx, fileID
func (_m *ImageReaderWriter) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error) {
	ret := _m.Called(ctx, fileID)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, models.FileID) []*models.Image); ok {
		r0 = rf(ctx, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FileID) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByFingerprints provides a mock function with given fields: ctx, fp
func (_m *ImageReaderWriter) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Image, error) {
	ret := _m.Called(ctx, fp)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, []models.Fingerprint) []*models.Image); ok {
		r0 = rf(ctx, fp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []models.Fingerprint) error); ok {
		r1 = rf(ctx, fp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByFolderID provides a mock function with given fields: ctx, fileID
func (_m *ImageReaderWriter) FindByFolderID(ctx context.Context, fileID models.FolderID) ([]*models.Image, error) {
	ret := _m.Called(ctx, fileID)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, models.FolderID) []*models.Image); ok {
		r0 = rf(ctx, fileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FolderID) error); ok {
		r1 = rf(ctx, fileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByGalleryID provides a mock function with given fields: ctx, galleryID
func (_m *ImageReaderWriter) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Image, error) {
	ret := _m.Called(ctx, galleryID)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Image); ok {
		r0 = rf(ctx, galleryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, galleryID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByZipFileID provides a mock function with given fields: ctx, zipFileID
func (_m *ImageReaderWriter) FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]*models.Image, error) {
	ret := _m.Called(ctx, zipFileID)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, models.FileID) []*models.Image); ok {
		r0 = rf(ctx, zipFileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.FileID) error); ok {
		r1 = rf(ctx, zipFileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields: ctx, ids
func (_m *ImageReaderWriter) FindMany(ctx context.Context, ids []int) ([]*models.Image, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*models.Image
	if rf, ok := ret.Get(0).(func(context.Context, []int) []*models.Image); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFiles provides a mock function with given fields: ctx, relatedID
func (_m *ImageReaderWriter) GetFiles(ctx context.Context, relatedID int) ([]models.File, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []models.File
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.File); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.File)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGalleryIDs provides a mock function with given fields: ctx, relatedID
func (_m *ImageReaderWriter) GetGalleryIDs(ctx context.Context, relatedID int) ([]int, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, int) []int); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManyFileIDs provides a mock function with given fields: ctx, ids
func (_m *ImageReaderWriter) GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error) {
	ret := _m.Called(ctx, ids)

	var r0 [][]models.FileID
	if rf, ok := ret.Get(0).(func(context.Context, []int) [][]models.FileID); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]models.FileID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPerformerIDs provides a mock function with given fields: ctx, relatedID
func (_m *ImageReaderWriter) GetPerformerIDs(ctx context.Context, relatedID int) ([]int, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, int) []int); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTagIDs provides a mock function with given fields: ctx, relatedID
func (_m *ImageReaderWriter) GetTagIDs(ctx context.Context, relatedID int) ([]int, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context, int) []int); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, relatedID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncrementOCounter provides a mock function with given fields: ctx, id
func (_m *ImageReaderWriter) IncrementOCounter(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OCountByPerformerID provides a mock function with given fields: ctx, performerID
func (_m *ImageReaderWriter) OCountByPerformerID(ctx context.Context, performerID int) (int, error) {
	ret := _m.Called(ctx, performerID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, performerID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx, options
func (_m *ImageReaderWriter) Query(ctx context.Context, options models.ImageQueryOptions) (*models.ImageQueryResult, error) {
	ret := _m.Called(ctx, options)

	var r0 *models.ImageQueryResult
	if rf, ok := ret.Get(0).(func(context.Context, models.ImageQueryOptions) *models.ImageQueryResult); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ImageQueryResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.ImageQueryOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryCount provides a mock function with given fields: ctx, imageFilter, findFilter
func (_m *ImageReaderWriter) QueryCount(ctx context.Context, imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) (int, error) {
	ret := _m.Called(ctx, imageFilter, findFilter)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, *models.ImageFilterType, *models.FindFilterType) int); ok {
		r0 = rf(ctx, imageFilter, findFilter)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.ImageFilterType, *models.FindFilterType) error); ok {
		r1 = rf(ctx, imageFilter, findFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetOCounter provides a mock function with given fields: ctx, id
func (_m *ImageReaderWriter) ResetOCounter(ctx context.Context, id int) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields: ctx
func (_m *ImageReaderWriter) Size(ctx context.Context) (float64, error) {
	ret := _m.Called(ctx)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, updatedImage
func (_m *ImageReaderWriter) Update(ctx context.Context, updatedImage *models.Image) error {
	ret := _m.Called(ctx, updatedImage)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Image) error); ok {
		r0 = rf(ctx, updatedImage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePartial provides a mock function with given fields: ctx, id, partial
func (_m *ImageReaderWriter) UpdatePartial(ctx context.Context, id int, partial models.ImagePartial) (*models.Image, error) {
	ret := _m.Called(ctx, id, partial)

	var r0 *models.Image
	if rf, ok := ret.Get(0).(func(context.Context, int, models.ImagePartial) *models.Image); ok {
		r0 = rf(ctx, id, partial)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, models.ImagePartial) error); ok {
		r1 = rf(ctx, id, partial)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePerformers provides a mock function with given fields: ctx, imageID, performerIDs
func (_m *ImageReaderWriter) UpdatePerformers(ctx context.Context, imageID int, performerIDs []int) error {
	ret := _m.Called(ctx, imageID, performerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []int) error); ok {
		r0 = rf(ctx, imageID, performerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTags provides a mock function with given fields: ctx, imageID, tagIDs
func (_m *ImageReaderWriter) UpdateTags(ctx context.Context, imageID int, tagIDs []int) error {
	ret := _m.Called(ctx, imageID, tagIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []int) error); ok {
		r0 = rf(ctx, imageID, tagIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
