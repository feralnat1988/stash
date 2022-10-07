package models

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/file"
)

type FileQueryOptions struct {
	QueryOptions
	FileFilter *FileFilterType
}

type FileFilterType struct {
	And *FileFilterType `json:"AND"`
	Or  *FileFilterType `json:"OR"`
	Not *FileFilterType `json:"NOT"`

	// Filter by path
	Path *StringCriterionInput `json:"path"`
}

func PathsFileFilter(paths []string) *FileFilterType {
	if paths == nil {
		return nil
	}

	sep := string(filepath.Separator)

	var ret *FileFilterType
	var or *FileFilterType
	for _, p := range paths {
		newOr := &FileFilterType{}
		if or != nil {
			or.Or = newOr
		} else {
			ret = newOr
		}

		or = newOr

		if !strings.HasSuffix(p, sep) {
			p += sep
		}

		or.Path = &StringCriterionInput{
			Modifier: CriterionModifierEquals,
			Value:    p + "%",
		}
	}

	return ret
}

type FileQueryResult struct {
	// can't use QueryResult because id type is wrong

	IDs   []file.ID
	Count int

	finder     file.Finder
	files      []file.File
	resolveErr error
}

func NewFileQueryResult(finder file.Finder) *FileQueryResult {
	return &FileQueryResult{
		finder: finder,
	}
}

func (r *FileQueryResult) Resolve(ctx context.Context) ([]file.File, error) {
	// cache results
	if r.files == nil && r.resolveErr == nil {
		r.files, r.resolveErr = r.finder.Find(ctx, r.IDs...)
	}
	return r.files, r.resolveErr
}
