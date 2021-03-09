package models

import (
	"database/sql"
	"path/filepath"
)

// Scene stores the metadata for a single video scene.
type Scene struct {
	ID          int                 `db:"id" json:"id"`
	Checksum    sql.NullString      `db:"checksum" json:"checksum"`
	OSHash      sql.NullString      `db:"oshash" json:"oshash"`
	Path        string              `db:"path" json:"path"`
	Title       sql.NullString      `db:"title" json:"title"`
	Details     sql.NullString      `db:"details" json:"details"`
	URL         sql.NullString      `db:"url" json:"url"`
	Date        SQLiteDate          `db:"date" json:"date"`
	Rating      sql.NullInt64       `db:"rating" json:"rating"`
	Organized   bool                `db:"organized" json:"organized"`
	OCounter    int                 `db:"o_counter" json:"o_counter"`
	Size        sql.NullString      `db:"size" json:"size"`
	Duration    sql.NullFloat64     `db:"duration" json:"duration"`
	VideoCodec  sql.NullString      `db:"video_codec" json:"video_codec"`
	Format      sql.NullString      `db:"format" json:"format_name"`
	AudioCodec  sql.NullString      `db:"audio_codec" json:"audio_codec"`
	Width       sql.NullInt64       `db:"width" json:"width"`
	Height      sql.NullInt64       `db:"height" json:"height"`
	Framerate   sql.NullFloat64     `db:"framerate" json:"framerate"`
	Bitrate     sql.NullInt64       `db:"bitrate" json:"bitrate"`
	StudioID    sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	FileModTime NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// ScenePartial represents part of a Scene object. It is used to update
// the database entry. Only non-nil fields will be updated.
type ScenePartial struct {
	ID          int                  `db:"id" json:"id"`
	Checksum    *sql.NullString      `db:"checksum" json:"checksum"`
	OSHash      *sql.NullString      `db:"oshash" json:"oshash"`
	Path        *string              `db:"path" json:"path"`
	Title       *sql.NullString      `db:"title" json:"title"`
	Details     *sql.NullString      `db:"details" json:"details"`
	URL         *sql.NullString      `db:"url" json:"url"`
	Date        *SQLiteDate          `db:"date" json:"date"`
	Rating      *sql.NullInt64       `db:"rating" json:"rating"`
	Organized   *bool                `db:"organized" json:"organized"`
	Size        *sql.NullString      `db:"size" json:"size"`
	Duration    *sql.NullFloat64     `db:"duration" json:"duration"`
	VideoCodec  *sql.NullString      `db:"video_codec" json:"video_codec"`
	Format      *sql.NullString      `db:"format" json:"format_name"`
	AudioCodec  *sql.NullString      `db:"audio_codec" json:"audio_codec"`
	Width       *sql.NullInt64       `db:"width" json:"width"`
	Height      *sql.NullInt64       `db:"height" json:"height"`
	Framerate   *sql.NullFloat64     `db:"framerate" json:"framerate"`
	Bitrate     *sql.NullInt64       `db:"bitrate" json:"bitrate"`
	StudioID    *sql.NullInt64       `db:"studio_id,omitempty" json:"studio_id"`
	MovieID     *sql.NullInt64       `db:"movie_id,omitempty" json:"movie_id"`
	FileModTime *NullSQLiteTimestamp `db:"file_mod_time" json:"file_mod_time"`
	CreatedAt   *SQLiteTimestamp     `db:"created_at" json:"created_at"`
	UpdatedAt   *SQLiteTimestamp     `db:"updated_at" json:"updated_at"`
}

// GetTitle returns the title of the scene. If the Title field is empty,
// then the base filename is returned.
func (s Scene) GetTitle() string {
	if s.Title.String != "" {
		return s.Title.String
	}

	return filepath.Base(s.Path)
}

// GetHash returns the hash of the scene, based on the hash algorithm provided. If
// hash algorithm is MD5, then Checksum is returned. Otherwise, OSHash is returned.
func (s Scene) GetHash(hashAlgorithm HashAlgorithm) string {
	if hashAlgorithm == HashAlgorithmMd5 {
		return s.Checksum.String
	} else if hashAlgorithm == HashAlgorithmOshash {
		return s.OSHash.String
	}

	panic("unknown hash algorithm")
}

func (s Scene) GetMinResolution() int64 {
	if s.Width.Int64 < s.Height.Int64 {
		return s.Width.Int64
	}

	return s.Height.Int64
}

// SceneFileType represents the file metadata for a scene.
type SceneFileType struct {
	Size       *string  `graphql:"size" json:"size"`
	Duration   *float64 `graphql:"duration" json:"duration"`
	VideoCodec *string  `graphql:"video_codec" json:"video_codec"`
	AudioCodec *string  `graphql:"audio_codec" json:"audio_codec"`
	Width      *int     `graphql:"width" json:"width"`
	Height     *int     `graphql:"height" json:"height"`
	Framerate  *float64 `graphql:"framerate" json:"framerate"`
	Bitrate    *int     `graphql:"bitrate" json:"bitrate"`
}

type Scenes []*Scene

func (s *Scenes) Append(o interface{}) {
	*s = append(*s, o.(*Scene))
}

func (s *Scenes) New() interface{} {
	return &Scene{}
}
