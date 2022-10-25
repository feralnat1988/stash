package models

import (
	"context"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/file"
)

type Gallery struct {
	ID int `json:"id"`

	Title     string `json:"title"`
	URL       string `json:"url"`
	Date      *Date  `json:"date"`
	Details   string `json:"details"`
	Rating    *int   `json:"rating"`
	Rating100 *int   `json:"rating100"`
	Organized bool   `json:"organized"`
	StudioID  *int   `json:"studio_id"`

	// transient - not persisted
	Files RelatedFiles
	// transient - not persisted
	PrimaryFileID *file.ID
	// transient - path of primary file or folder
	Path string

	FolderID *file.FolderID `json:"folder_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	SceneIDs     RelatedIDs `json:"scene_ids"`
	TagIDs       RelatedIDs `json:"tag_ids"`
	PerformerIDs RelatedIDs `json:"performer_ids"`
}

func (g *Gallery) LoadFiles(ctx context.Context, l FileLoader) error {
	return g.Files.load(func() ([]file.File, error) {
		return l.GetFiles(ctx, g.ID)
	})
}

func (g *Gallery) LoadPrimaryFile(ctx context.Context, l file.Finder) error {
	return g.Files.loadPrimary(func() (file.File, error) {
		if g.PrimaryFileID == nil {
			return nil, nil
		}

		f, err := l.Find(ctx, *g.PrimaryFileID)
		if err != nil {
			return nil, err
		}

		if len(f) > 0 {
			return f[0], nil
		}
		return nil, nil
	})
}

func (g *Gallery) LoadSceneIDs(ctx context.Context, l SceneIDLoader) error {
	return g.SceneIDs.load(func() ([]int, error) {
		return l.GetSceneIDs(ctx, g.ID)
	})
}

func (g *Gallery) LoadPerformerIDs(ctx context.Context, l PerformerIDLoader) error {
	return g.PerformerIDs.load(func() ([]int, error) {
		return l.GetPerformerIDs(ctx, g.ID)
	})
}

func (g *Gallery) LoadTagIDs(ctx context.Context, l TagIDLoader) error {
	return g.TagIDs.load(func() ([]int, error) {
		return l.GetTagIDs(ctx, g.ID)
	})
}

func (g Gallery) PrimaryChecksum() string {
	// renamed from Checksum to prevent gqlgen from using it in the resolver
	if p := g.Files.Primary(); p != nil {
		v := p.Base().Fingerprints.Get(file.FingerprintTypeMD5)
		if v == nil {
			return ""
		}

		return v.(string)
	}
	return ""
}

// GalleryPartial represents part of a Gallery object. It is used to update
// the database entry. Only non-nil fields will be updated.
type GalleryPartial struct {
	// Path        OptionalString
	// Checksum    OptionalString
	// Zip         OptionalBool
	Title     OptionalString
	URL       OptionalString
	Date      OptionalDate
	Details   OptionalString
	Rating    OptionalInt
	Rating100 OptionalInt
	Organized OptionalBool
	StudioID  OptionalInt
	// FileModTime OptionalTime
	CreatedAt OptionalTime
	UpdatedAt OptionalTime

	SceneIDs      *UpdateIDs
	TagIDs        *UpdateIDs
	PerformerIDs  *UpdateIDs
	PrimaryFileID *file.ID
}

func NewGalleryPartial() GalleryPartial {
	updatedTime := time.Now()
	return GalleryPartial{
		UpdatedAt: NewOptionalTime(updatedTime),
	}
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (g Gallery) GetTitle() string {
	if g.Title != "" {
		return g.Title
	}

	return g.Path
}

// DisplayName returns a display name for the scene for logging purposes.
// It returns the path or title, or otherwise it returns the ID if both of these are empty.
func (g Gallery) DisplayName() string {
	if g.Path != "" {
		return g.Path
	}

	if g.Title != "" {
		return g.Title
	}

	return strconv.Itoa(g.ID)
}

const DefaultGthumbWidth int = 640

type Galleries []*Gallery

func (g *Galleries) Append(o interface{}) {
	*g = append(*g, o.(*Gallery))
}

func (g *Galleries) New() interface{} {
	return &Gallery{}
}
