// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// PerformerReaderWriter is an autogenerated mock type for the PerformerReaderWriter type
type PerformerReaderWriter struct {
	mock.Mock
}

// All provides a mock function with given fields: ctx
func (_m *PerformerReaderWriter) All(ctx context.Context) ([]*models.Performer, error) {
	ret := _m.Called(ctx)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context) []*models.Performer); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
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
func (_m *PerformerReaderWriter) Count(ctx context.Context) (int, error) {
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

// CountByTagID provides a mock function with given fields: ctx, tagID
func (_m *PerformerReaderWriter) CountByTagID(ctx context.Context, tagID int) (int, error) {
	ret := _m.Called(ctx, tagID)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, int) int); ok {
		r0 = rf(ctx, tagID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, tagID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, newPerformer
func (_m *PerformerReaderWriter) Create(ctx context.Context, newPerformer *models.Performer) error {
	ret := _m.Called(ctx, newPerformer)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Performer) error); ok {
		r0 = rf(ctx, newPerformer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Destroy provides a mock function with given fields: ctx, id
func (_m *PerformerReaderWriter) Destroy(ctx context.Context, id int) error {
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
func (_m *PerformerReaderWriter) Find(ctx context.Context, id int) (*models.Performer, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.Performer); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Performer)
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

// FindByGalleryID provides a mock function with given fields: ctx, galleryID
func (_m *PerformerReaderWriter) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Performer, error) {
	ret := _m.Called(ctx, galleryID)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Performer); ok {
		r0 = rf(ctx, galleryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
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

// FindByImageID provides a mock function with given fields: ctx, imageID
func (_m *PerformerReaderWriter) FindByImageID(ctx context.Context, imageID int) ([]*models.Performer, error) {
	ret := _m.Called(ctx, imageID)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Performer); ok {
		r0 = rf(ctx, imageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, imageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByNames provides a mock function with given fields: ctx, names, nocase
func (_m *PerformerReaderWriter) FindByNames(ctx context.Context, names []string, nocase bool) ([]*models.Performer, error) {
	ret := _m.Called(ctx, names, nocase)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, []string, bool) []*models.Performer); ok {
		r0 = rf(ctx, names, nocase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, bool) error); ok {
		r1 = rf(ctx, names, nocase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindBySceneID provides a mock function with given fields: ctx, sceneID
func (_m *PerformerReaderWriter) FindBySceneID(ctx context.Context, sceneID int) ([]*models.Performer, error) {
	ret := _m.Called(ctx, sceneID)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Performer); ok {
		r0 = rf(ctx, sceneID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, sceneID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByStashID provides a mock function with given fields: ctx, stashID
func (_m *PerformerReaderWriter) FindByStashID(ctx context.Context, stashID models.StashID) ([]*models.Performer, error) {
	ret := _m.Called(ctx, stashID)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, models.StashID) []*models.Performer); ok {
		r0 = rf(ctx, stashID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.StashID) error); ok {
		r1 = rf(ctx, stashID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByStashIDStatus provides a mock function with given fields: ctx, hasStashID, stashboxEndpoint
func (_m *PerformerReaderWriter) FindByStashIDStatus(ctx context.Context, hasStashID bool, stashboxEndpoint string) ([]*models.Performer, error) {
	ret := _m.Called(ctx, hasStashID, stashboxEndpoint)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, bool, string) []*models.Performer); ok {
		r0 = rf(ctx, hasStashID, stashboxEndpoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, bool, string) error); ok {
		r1 = rf(ctx, hasStashID, stashboxEndpoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields: ctx, ids
func (_m *PerformerReaderWriter) FindMany(ctx context.Context, ids []int) ([]*models.Performer, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, []int) []*models.Performer); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
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

// GetAliases provides a mock function with given fields: ctx, relatedID
func (_m *PerformerReaderWriter) GetAliases(ctx context.Context, relatedID int) ([]string, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, int) []string); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
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

// GetImage provides a mock function with given fields: ctx, performerID
func (_m *PerformerReaderWriter) GetImage(ctx context.Context, performerID int) ([]byte, error) {
	ret := _m.Called(ctx, performerID)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, int) []byte); ok {
		r0 = rf(ctx, performerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStashIDs provides a mock function with given fields: ctx, relatedID
func (_m *PerformerReaderWriter) GetStashIDs(ctx context.Context, relatedID int) ([]models.StashID, error) {
	ret := _m.Called(ctx, relatedID)

	var r0 []models.StashID
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.StashID); ok {
		r0 = rf(ctx, relatedID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.StashID)
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
func (_m *PerformerReaderWriter) GetTagIDs(ctx context.Context, relatedID int) ([]int, error) {
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

// HasImage provides a mock function with given fields: ctx, performerID
func (_m *PerformerReaderWriter) HasImage(ctx context.Context, performerID int) (bool, error) {
	ret := _m.Called(ctx, performerID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, performerID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx, performerFilter, findFilter
func (_m *PerformerReaderWriter) Query(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, int, error) {
	ret := _m.Called(ctx, performerFilter, findFilter)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, *models.PerformerFilterType, *models.FindFilterType) []*models.Performer); ok {
		r0 = rf(ctx, performerFilter, findFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, *models.PerformerFilterType, *models.FindFilterType) int); ok {
		r1 = rf(ctx, performerFilter, findFilter)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *models.PerformerFilterType, *models.FindFilterType) error); ok {
		r2 = rf(ctx, performerFilter, findFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryCount provides a mock function with given fields: ctx, performerFilter, findFilter
func (_m *PerformerReaderWriter) QueryCount(ctx context.Context, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) (int, error) {
	ret := _m.Called(ctx, performerFilter, findFilter)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, *models.PerformerFilterType, *models.FindFilterType) int); ok {
		r0 = rf(ctx, performerFilter, findFilter)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.PerformerFilterType, *models.FindFilterType) error); ok {
		r1 = rf(ctx, performerFilter, findFilter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryForAutoTag provides a mock function with given fields: ctx, words
func (_m *PerformerReaderWriter) QueryForAutoTag(ctx context.Context, words []string) ([]*models.Performer, error) {
	ret := _m.Called(ctx, words)

	var r0 []*models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*models.Performer); ok {
		r0 = rf(ctx, words)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, words)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, updatedPerformer
func (_m *PerformerReaderWriter) Update(ctx context.Context, updatedPerformer *models.Performer) error {
	ret := _m.Called(ctx, updatedPerformer)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Performer) error); ok {
		r0 = rf(ctx, updatedPerformer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateImage provides a mock function with given fields: ctx, performerID, image
func (_m *PerformerReaderWriter) UpdateImage(ctx context.Context, performerID int, image []byte) error {
	ret := _m.Called(ctx, performerID, image)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []byte) error); ok {
		r0 = rf(ctx, performerID, image)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePartial provides a mock function with given fields: ctx, id, updatedPerformer
func (_m *PerformerReaderWriter) UpdatePartial(ctx context.Context, id int, updatedPerformer models.PerformerPartial) (*models.Performer, error) {
	ret := _m.Called(ctx, id, updatedPerformer)

	var r0 *models.Performer
	if rf, ok := ret.Get(0).(func(context.Context, int, models.PerformerPartial) *models.Performer); ok {
		r0 = rf(ctx, id, updatedPerformer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Performer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, models.PerformerPartial) error); ok {
		r1 = rf(ctx, id, updatedPerformer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
