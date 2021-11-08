// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// TagReaderWriter is an autogenerated mock type for the TagReaderWriter type
type TagReaderWriter struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *TagReaderWriter) All() ([]*models.Tag, error) {
	ret := _m.Called()

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func() []*models.Tag); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Count provides a mock function with given fields:
func (_m *TagReaderWriter) Count() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: newTag
func (_m *TagReaderWriter) Create(newTag models.Tag) (*models.Tag, error) {
	ret := _m.Called(newTag)

	var r0 *models.Tag
	if rf, ok := ret.Get(0).(func(models.Tag) *models.Tag); ok {
		r0 = rf(newTag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Tag) error); ok {
		r1 = rf(newTag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: id
func (_m *TagReaderWriter) Destroy(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DestroyImage provides a mock function with given fields: tagID
func (_m *TagReaderWriter) DestroyImage(tagID int) error {
	ret := _m.Called(tagID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(tagID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *TagReaderWriter) Find(id int) (*models.Tag, error) {
	ret := _m.Called(id)

	var r0 *models.Tag
	if rf, ok := ret.Get(0).(func(int) *models.Tag); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllAncestors provides a mock function with given fields: tagID, excludeIDs
func (_m *TagReaderWriter) FindAllAncestors(tagID int, excludeIDs []int) ([]*models.TagPath, error) {
	ret := _m.Called(tagID, excludeIDs)

	var r0 []*models.TagPath
	if rf, ok := ret.Get(0).(func(int, []int) []*models.TagPath); ok {
		r0 = rf(tagID, excludeIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.TagPath)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, []int) error); ok {
		r1 = rf(tagID, excludeIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllDescendants provides a mock function with given fields: tagID, excludeIDs
func (_m *TagReaderWriter) FindAllDescendants(tagID int, excludeIDs []int) ([]*models.TagPath, error) {
	ret := _m.Called(tagID, excludeIDs)

	var r0 []*models.TagPath
	if rf, ok := ret.Get(0).(func(int, []int) []*models.TagPath); ok {
		r0 = rf(tagID, excludeIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.TagPath)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, []int) error); ok {
		r1 = rf(tagID, excludeIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByChildTagID provides a mock function with given fields: childID
func (_m *TagReaderWriter) FindByChildTagID(childID int) ([]*models.Tag, error) {
	ret := _m.Called(childID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(childID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(childID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByGalleryID provides a mock function with given fields: galleryID
func (_m *TagReaderWriter) FindByGalleryID(galleryID int) ([]*models.Tag, error) {
	ret := _m.Called(galleryID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(galleryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(galleryID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByImageID provides a mock function with given fields: imageID
func (_m *TagReaderWriter) FindByImageID(imageID int) ([]*models.Tag, error) {
	ret := _m.Called(imageID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(imageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(imageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: name, nocase
func (_m *TagReaderWriter) FindByName(name string, nocase bool) (*models.Tag, error) {
	ret := _m.Called(name, nocase)

	var r0 *models.Tag
	if rf, ok := ret.Get(0).(func(string, bool) *models.Tag); ok {
		r0 = rf(name, nocase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(name, nocase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByNames provides a mock function with given fields: names, nocase
func (_m *TagReaderWriter) FindByNames(names []string, nocase bool) ([]*models.Tag, error) {
	ret := _m.Called(names, nocase)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func([]string, bool) []*models.Tag); ok {
		r0 = rf(names, nocase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, bool) error); ok {
		r1 = rf(names, nocase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByParentTagID provides a mock function with given fields: parentID
func (_m *TagReaderWriter) FindByParentTagID(parentID int) ([]*models.Tag, error) {
	ret := _m.Called(parentID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(parentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(parentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPerformerID provides a mock function with given fields: performerID
func (_m *TagReaderWriter) FindByPerformerID(performerID int) ([]*models.Tag, error) {
	ret := _m.Called(performerID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(performerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindBySceneID provides a mock function with given fields: sceneID
func (_m *TagReaderWriter) FindBySceneID(sceneID int) ([]*models.Tag, error) {
	ret := _m.Called(sceneID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(sceneID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(sceneID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindBySceneMarkerID provides a mock function with given fields: sceneMarkerID
func (_m *TagReaderWriter) FindBySceneMarkerID(sceneMarkerID int) ([]*models.Tag, error) {
	ret := _m.Called(sceneMarkerID)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(int) []*models.Tag); ok {
		r0 = rf(sceneMarkerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(sceneMarkerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields: ids
func (_m *TagReaderWriter) FindMany(ids []int) ([]*models.Tag, error) {
	ret := _m.Called(ids)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func([]int) []*models.Tag); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAliases provides a mock function with given fields: tagID
func (_m *TagReaderWriter) GetAliases(tagID int) ([]string, error) {
	ret := _m.Called(tagID)

	var r0 []string
	if rf, ok := ret.Get(0).(func(int) []string); ok {
		r0 = rf(tagID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(tagID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetImage provides a mock function with given fields: tagID
func (_m *TagReaderWriter) GetImage(tagID int) ([]byte, error) {
	ret := _m.Called(tagID)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(int) []byte); ok {
		r0 = rf(tagID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(tagID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Merge provides a mock function with given fields: source, destination
func (_m *TagReaderWriter) Merge(source []int, destination int) error {
	ret := _m.Called(source, destination)

	var r0 error
	if rf, ok := ret.Get(0).(func([]int, int) error); ok {
		r0 = rf(source, destination)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Query provides a mock function with given fields: tagFilter, findFilter
func (_m *TagReaderWriter) Query(tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, int, error) {
	ret := _m.Called(tagFilter, findFilter)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func(*models.TagFilterType, *models.FindFilterType) []*models.Tag); ok {
		r0 = rf(tagFilter, findFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(*models.TagFilterType, *models.FindFilterType) int); ok {
		r1 = rf(tagFilter, findFilter)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*models.TagFilterType, *models.FindFilterType) error); ok {
		r2 = rf(tagFilter, findFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryForAutoTag provides a mock function with given fields: words
func (_m *TagReaderWriter) QueryForAutoTag(words []string) ([]*models.Tag, error) {
	ret := _m.Called(words)

	var r0 []*models.Tag
	if rf, ok := ret.Get(0).(func([]string) []*models.Tag); ok {
		r0 = rf(words)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(words)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: updateTag
func (_m *TagReaderWriter) Update(updateTag models.TagPartial) (*models.Tag, error) {
	ret := _m.Called(updateTag)

	var r0 *models.Tag
	if rf, ok := ret.Get(0).(func(models.TagPartial) *models.Tag); ok {
		r0 = rf(updateTag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.TagPartial) error); ok {
		r1 = rf(updateTag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAliases provides a mock function with given fields: tagID, aliases
func (_m *TagReaderWriter) UpdateAliases(tagID int, aliases []string) error {
	ret := _m.Called(tagID, aliases)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []string) error); ok {
		r0 = rf(tagID, aliases)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateChildTags provides a mock function with given fields: tagID, parentIDs
func (_m *TagReaderWriter) UpdateChildTags(tagID int, parentIDs []int) error {
	ret := _m.Called(tagID, parentIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []int) error); ok {
		r0 = rf(tagID, parentIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateFull provides a mock function with given fields: updatedTag
func (_m *TagReaderWriter) UpdateFull(updatedTag models.Tag) (*models.Tag, error) {
	ret := _m.Called(updatedTag)

	var r0 *models.Tag
	if rf, ok := ret.Get(0).(func(models.Tag) *models.Tag); ok {
		r0 = rf(updatedTag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Tag) error); ok {
		r1 = rf(updatedTag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateImage provides a mock function with given fields: tagID, image
func (_m *TagReaderWriter) UpdateImage(tagID int, image []byte) error {
	ret := _m.Called(tagID, image)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []byte) error); ok {
		r0 = rf(tagID, image)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateParentTags provides a mock function with given fields: tagID, parentIDs
func (_m *TagReaderWriter) UpdateParentTags(tagID int, parentIDs []int) error {
	ret := _m.Called(tagID, parentIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []int) error); ok {
		r0 = rf(tagID, parentIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
