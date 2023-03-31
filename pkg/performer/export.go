package performer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/json"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type ImageAliasStashIDGetter interface {
	GetImage(ctx context.Context, performerID int) ([]byte, error)
	models.AliasLoader
	models.StashIDLoader
}

// ToJSON converts a Performer object into its JSON equivalent.
func ToJSON(ctx context.Context, reader ImageAliasStashIDGetter, performer *models.Performer) (*jsonschema.Performer, error) {
	newPerformerJSON := jsonschema.Performer{
		Name:           performer.Name,
		Disambiguation: performer.Disambiguation,
		Gender:         performer.Gender.String(),
		URL:            performer.URL,
		Ethnicity:      performer.Ethnicity,
		Country:        performer.Country,
		EyeColor:       performer.EyeColor,
		Measurements:   performer.Measurements,
		FakeTits:       performer.FakeTits,
		Circumcised:    performer.Circumcised.String(),
		CareerLength:   performer.CareerLength,
		Tattoos:        performer.Tattoos,
		Piercings:      performer.Piercings,
		Twitter:        performer.Twitter,
		Instagram:      performer.Instagram,
		Favorite:       performer.Favorite,
		Details:        performer.Details,
		HairColor:      performer.HairColor,
		IgnoreAutoTag:  performer.IgnoreAutoTag,
		CreatedAt:      json.JSONTime{Time: performer.CreatedAt},
		UpdatedAt:      json.JSONTime{Time: performer.UpdatedAt},
	}

	if performer.Birthdate != nil {
		newPerformerJSON.Birthdate = performer.Birthdate.String()
	}
	if performer.Rating != nil {
		newPerformerJSON.Rating = *performer.Rating
	}
	if performer.DeathDate != nil {
		newPerformerJSON.DeathDate = performer.DeathDate.String()
	}

	if performer.Height != nil {
		newPerformerJSON.Height = strconv.Itoa(*performer.Height)
	}

	if performer.Weight != nil {
		newPerformerJSON.Weight = *performer.Weight
	}

	if performer.PenisLength != nil {
		newPerformerJSON.PenisLength = *performer.PenisLength
	}

	if err := performer.LoadAliases(ctx, reader); err != nil {
		return nil, fmt.Errorf("loading performer aliases: %w", err)
	}

	newPerformerJSON.Aliases = performer.Aliases.List()

	if err := performer.LoadStashIDs(ctx, reader); err != nil {
		return nil, fmt.Errorf("loading performer stash ids: %w", err)
	}

	newPerformerJSON.StashIDs = performer.StashIDs.List()

	image, err := reader.GetImage(ctx, performer.ID)
	if err != nil {
		logger.Errorf("Error getting performer image: %v", err)
	}

	if len(image) > 0 {
		newPerformerJSON.Image = utils.GetBase64StringFromData(image)
	}

	return &newPerformerJSON, nil
}

func GetIDs(performers []*models.Performer) []int {
	var results []int
	for _, performer := range performers {
		results = append(results, performer.ID)
	}

	return results
}

func GetNames(performers []*models.Performer) []string {
	var results []string
	for _, performer := range performers {
		if performer.Name != "" {
			results = append(results, performer.Name)
		}
	}

	return results
}
