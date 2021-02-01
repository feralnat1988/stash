package studio

import (
	"fmt"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

// ToJSON converts a Studio object into its JSON equivalent.
func ToJSON(reader models.StudioReader, studio *models.Studio) (*jsonschema.Studio, error) {
	newStudioJSON := jsonschema.Studio{
		CreatedAt: models.JSONTime{Time: studio.CreatedAt.Timestamp},
		UpdatedAt: models.JSONTime{Time: studio.UpdatedAt.Timestamp},
	}

	if studio.Name.Valid {
		newStudioJSON.Name = studio.Name.String
	}

	if studio.URL.Valid {
		newStudioJSON.URL = studio.URL.String
	}

	if studio.ParentID.Valid {
		parent, err := reader.Find(int(studio.ParentID.Int64))
		if err != nil {
			return nil, fmt.Errorf("error getting parent studio: %s", err.Error())
		}

		if parent != nil {
			newStudioJSON.ParentStudio = parent.Name.String
		}
	}

	image, err := reader.GetImage(studio.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting studio image: %s", err.Error())
	}

	if len(image) > 0 {
		newStudioJSON.Image = utils.GetBase64StringFromData(image)
	}

	return &newStudioJSON, nil
}
