package studio

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	studioNameErr      = "studioNameErr"
	existingStudioName = "existingTagName"

	existingStudioID = 100

	existingParentStudioName = "existingParentStudioName"
	existingParentStudioErr  = "existingParentStudioErr"
	missingParentStudioName  = "existingParentStudioName"
)

var testCtx = context.Background()

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Studio{
			Name: studioName,
		},
	}

	assert.Equal(t, studioName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Studio{
			Name:          studioName,
			Image:         invalidImage,
			IgnoreAutoTag: autoTagIgnored,
		},
	}

	err := i.PreImport(testCtx)

	assert.NotNil(t, err)

	i.Input.Image = image

	err = i.PreImport(testCtx)

	assert.Nil(t, err)

	i.Input = *createFullJSONStudio(studioName, image, []string{"alias"})
	i.Input.ParentStudio = ""

	err = i.PreImport(testCtx)

	assert.Nil(t, err)
	expectedStudio := createFullStudio(0, 0)
	expectedStudio.ParentID = nil
	assert.Equal(t, expectedStudio, i.studio)
}

func TestImporterPreImportWithParent(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Studio,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: existingParentStudioName,
		},
	}

	db.Studio.On("FindByName", testCtx, existingParentStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	db.Studio.On("FindByName", testCtx, existingParentStudioErr, false).Return(nil, errors.New("FindByName error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingStudioID, *i.studio.ParentID)

	i.Input.ParentStudio = existingParentStudioErr
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.Studio.AssertExpectations(t)
}

func TestImporterPreImportWithMissingParent(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Studio,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: missingParentStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Studio.On("FindByName", testCtx, missingParentStudioName, false).Return(nil, nil).Times(3)
	db.Studio.On("Create", testCtx, mock.AnythingOfType("*models.Studio")).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Studio)
		s.ID = existingStudioID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingStudioID, *i.studio.ParentID)

	db.Studio.AssertExpectations(t)
}

func TestImporterPreImportWithMissingParentCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Studio,
		Input: jsonschema.Studio{
			Name:         studioName,
			Image:        image,
			ParentStudio: missingParentStudioName,
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Studio.On("FindByName", testCtx, missingParentStudioName, false).Return(nil, nil).Once()
	db.Studio.On("Create", testCtx, mock.AnythingOfType("*models.Studio")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Studio,
		Input: jsonschema.Studio{
			Aliases: []string{"alias"},
		},
		imageData: imageBytes,
	}

	updateStudioImageErr := errors.New("UpdateImage error")

	db.Studio.On("UpdateImage", testCtx, studioID, imageBytes).Return(nil).Once()
	db.Studio.On("UpdateImage", testCtx, errImageID, imageBytes).Return(updateStudioImageErr).Once()

	err := i.PostImport(testCtx, studioID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errImageID)
	assert.NotNil(t, err)

	db.Studio.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Studio,
		Input: jsonschema.Studio{
			Name: studioName,
		},
	}

	errFindByName := errors.New("FindByName error")
	db.Studio.On("FindByName", testCtx, studioName, false).Return(nil, nil).Once()
	db.Studio.On("FindByName", testCtx, existingStudioName, false).Return(&models.Studio{
		ID: existingStudioID,
	}, nil).Once()
	db.Studio.On("FindByName", testCtx, studioNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingStudioName
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingStudioID, *id)
	assert.Nil(t, err)

	i.Input.Name = studioNameErr
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.Studio.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	db := mocks.NewDatabase()

	studio := models.Studio{
		Name: studioName,
	}

	studioErr := models.Studio{
		Name: studioNameErr,
	}

	i := Importer{
		ReaderWriter: db.Studio,
		studio:       studio,
	}

	errCreate := errors.New("Create error")
	db.Studio.On("Create", testCtx, &studio).Run(func(args mock.Arguments) {
		s := args.Get(1).(*models.Studio)
		s.ID = studioID
	}).Return(nil).Once()
	db.Studio.On("Create", testCtx, &studioErr).Return(errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, studioID, *id)
	assert.Nil(t, err)

	i.studio = studioErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.Studio.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	db := mocks.NewDatabase()

	studio := models.Studio{
		Name: studioName,
	}

	studioErr := models.Studio{
		Name: studioNameErr,
	}

	i := Importer{
		ReaderWriter: db.Studio,
		studio:       studio,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	studio.ID = studioID
	db.Studio.On("Update", testCtx, &studio).Return(nil).Once()

	err := i.Update(testCtx, studioID)
	assert.Nil(t, err)

	i.studio = studioErr

	// need to set id separately
	studioErr.ID = errImageID
	db.Studio.On("Update", testCtx, &studioErr).Return(errUpdate).Once()

	err = i.Update(testCtx, errImageID)
	assert.NotNil(t, err)

	db.Studio.AssertExpectations(t)
}
