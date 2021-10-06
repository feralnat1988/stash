// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	models "github.com/stashapp/stash/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// MovieReaderWriter is an autogenerated mock type for the MovieReaderWriter type
type MovieReaderWriter struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *MovieReaderWriter) All() ([]*models.Movie, error) {
	ret := _m.Called()

	var r0 []*models.Movie
	if rf, ok := ret.Get(0).(func() []*models.Movie); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Movie)
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
func (_m *MovieReaderWriter) Count() (int, error) {
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

// CountByPerformerID provides a mock function with given fields: performerID
func (_m *MovieReaderWriter) CountByPerformerID(performerID int) (int, error) {
	ret := _m.Called(performerID)

	var r0 int
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(performerID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(performerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountByStudioID provides a mock function with given fields: studioID
func (_m *MovieReaderWriter) CountByStudioID(studioID int) (int, error) {
	ret := _m.Called(studioID)

	var r0 int
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(studioID)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(studioID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: newMovie
func (_m *MovieReaderWriter) Create(newMovie models.Movie) (*models.Movie, error) {
	ret := _m.Called(newMovie)

	var r0 *models.Movie
	if rf, ok := ret.Get(0).(func(models.Movie) *models.Movie); ok {
		r0 = rf(newMovie)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Movie)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Movie) error); ok {
		r1 = rf(newMovie)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: id
func (_m *MovieReaderWriter) Destroy(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DestroyImages provides a mock function with given fields: movieID
func (_m *MovieReaderWriter) DestroyImages(movieID int) error {
	ret := _m.Called(movieID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(movieID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *MovieReaderWriter) Find(id int) (*models.Movie, error) {
	ret := _m.Called(id)

	var r0 *models.Movie
	if rf, ok := ret.Get(0).(func(int) *models.Movie); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Movie)
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

// FindByName provides a mock function with given fields: name, nocase
func (_m *MovieReaderWriter) FindByName(name string, nocase bool) (*models.Movie, error) {
	ret := _m.Called(name, nocase)

	var r0 *models.Movie
	if rf, ok := ret.Get(0).(func(string, bool) *models.Movie); ok {
		r0 = rf(name, nocase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Movie)
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
func (_m *MovieReaderWriter) FindByNames(names []string, nocase bool) ([]*models.Movie, error) {
	ret := _m.Called(names, nocase)

	var r0 []*models.Movie
	if rf, ok := ret.Get(0).(func([]string, bool) []*models.Movie); ok {
		r0 = rf(names, nocase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Movie)
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

// FindByPerformerID provides a mock function with given fields: performerID
func (_m *MovieReaderWriter) FindByPerformerID(performerID int) ([]*models.Movie, error) {
	ret := _m.Called(performerID)

	var r0 []*models.Movie
	if rf, ok := ret.Get(0).(func(int) []*models.Movie); ok {
		r0 = rf(performerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Movie)
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

// FindByStudioID provides a mock function with given fields: studioID
func (_m *MovieReaderWriter) FindByStudioID(studioID int) ([]*models.Movie, error) {
	ret := _m.Called(studioID)

	var r0 []*models.Movie
	if rf, ok := ret.Get(0).(func(int) []*models.Movie); ok {
		r0 = rf(studioID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Movie)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(studioID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMany provides a mock function with given fields: ids
func (_m *MovieReaderWriter) FindMany(ids []int) ([]*models.Movie, error) {
	ret := _m.Called(ids)

	var r0 []*models.Movie
	if rf, ok := ret.Get(0).(func([]int) []*models.Movie); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Movie)
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

// GetBackImage provides a mock function with given fields: movieID
func (_m *MovieReaderWriter) GetBackImage(movieID int) ([]byte, error) {
	ret := _m.Called(movieID)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(int) []byte); ok {
		r0 = rf(movieID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(movieID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFrontImage provides a mock function with given fields: movieID
func (_m *MovieReaderWriter) GetFrontImage(movieID int) ([]byte, error) {
	ret := _m.Called(movieID)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(int) []byte); ok {
		r0 = rf(movieID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(movieID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: movieFilter, findFilter
func (_m *MovieReaderWriter) Query(movieFilter *models.MovieFilterType, findFilter *models.FindFilterType) ([]*models.Movie, int, error) {
	ret := _m.Called(movieFilter, findFilter)

	var r0 []*models.Movie
	if rf, ok := ret.Get(0).(func(*models.MovieFilterType, *models.FindFilterType) []*models.Movie); ok {
		r0 = rf(movieFilter, findFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Movie)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(*models.MovieFilterType, *models.FindFilterType) int); ok {
		r1 = rf(movieFilter, findFilter)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*models.MovieFilterType, *models.FindFilterType) error); ok {
		r2 = rf(movieFilter, findFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Update provides a mock function with given fields: updatedMovie
func (_m *MovieReaderWriter) Update(updatedMovie models.MoviePartial) (*models.Movie, error) {
	ret := _m.Called(updatedMovie)

	var r0 *models.Movie
	if rf, ok := ret.Get(0).(func(models.MoviePartial) *models.Movie); ok {
		r0 = rf(updatedMovie)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Movie)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.MoviePartial) error); ok {
		r1 = rf(updatedMovie)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFull provides a mock function with given fields: updatedMovie
func (_m *MovieReaderWriter) UpdateFull(updatedMovie models.Movie) (*models.Movie, error) {
	ret := _m.Called(updatedMovie)

	var r0 *models.Movie
	if rf, ok := ret.Get(0).(func(models.Movie) *models.Movie); ok {
		r0 = rf(updatedMovie)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Movie)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Movie) error); ok {
		r1 = rf(updatedMovie)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateImages provides a mock function with given fields: movieID, frontImage, backImage
func (_m *MovieReaderWriter) UpdateImages(movieID int, frontImage []byte, backImage []byte) error {
	ret := _m.Called(movieID, frontImage, backImage)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []byte, []byte) error); ok {
		r0 = rf(movieID, frontImage, backImage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
