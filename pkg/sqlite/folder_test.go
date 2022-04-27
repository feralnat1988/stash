//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stretchr/testify/assert"
)

var (
	invalidFolderID = file.FolderID(invalidID)
	invalidFileID   = file.ID(invalidID)
)

func Test_folderQueryBuilder_Create(t *testing.T) {
	var (
		path        = "path"
		fileModTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt   = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	)

	tests := []struct {
		name      string
		newObject file.Folder
		wantErr   bool
	}{
		{
			"full",
			file.Folder{
				DirEntry: file.DirEntry{
					ZipFileID:    &fileIDs[fileIdxZip],
					ModTime:      fileModTime,
					MissingSince: &updatedAt,
					LastScanned:  createdAt,
				},
				Path:      path,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"invalid parent folder id",
			file.Folder{
				Path:           path,
				ParentFolderID: &invalidFolderID,
			},
			true,
		},
		{
			"invalid zip file id",
			file.Folder{
				DirEntry: file.DirEntry{
					ZipFileID: &invalidFileID,
				},
				Path: path,
			},
			true,
		},
	}

	qb := db.Folder

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			s := tt.newObject
			if err := qb.Create(ctx, &s); (err != nil) != tt.wantErr {
				t.Errorf("folderQueryBuilder.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.ID)
				return
			}

			assert.NotZero(s.ID)

			copy := tt.newObject
			copy.ID = s.ID

			assert.Equal(copy, s)

			// ensure can find the folder
			found, err := qb.FindByPath(ctx, path)
			if err != nil {
				t.Errorf("folderQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(copy, *found)
		})
	}
}

func Test_folderQueryBuilder_Update(t *testing.T) {
	var (
		path        = "path"
		fileModTime = time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
		createdAt   = time.Date(2001, 1, 2, 3, 4, 5, 6, time.UTC)
		updatedAt   = time.Date(2002, 1, 2, 3, 4, 5, 6, time.UTC)
	)

	tests := []struct {
		name          string
		updatedObject *file.Folder
		wantErr       bool
	}{
		{
			"full",
			&file.Folder{
				ID: folderIDs[folderIdxIsMissing],
				DirEntry: file.DirEntry{
					ZipFileID:    &fileIDs[fileIdxZip],
					ModTime:      fileModTime,
					MissingSince: &updatedAt,
					LastScanned:  createdAt,
				},
				Path:      path,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"clear zip",
			&file.Folder{
				ID:   folderIDs[folderIdxInZip],
				Path: path,
			},
			false,
		},
		{
			"clear missing since",
			&file.Folder{
				ID:   folderIDs[folderIdxIsMissing],
				Path: path,
			},
			false,
		},
		{
			"clear folder",
			&file.Folder{
				ID:   folderIDs[folderIdxWithParentFolder],
				Path: path,
			},
			false,
		},
		{
			"invalid parent folder id",
			&file.Folder{
				ID:             folderIDs[folderIdxIsMissing],
				Path:           path,
				ParentFolderID: &invalidFolderID,
			},
			true,
		},
		{
			"invalid zip file id",
			&file.Folder{
				ID: folderIDs[folderIdxIsMissing],
				DirEntry: file.DirEntry{
					ZipFileID: &invalidFileID,
				},
				Path: path,
			},
			true,
		},
	}

	qb := db.Folder
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := *tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("folderQueryBuilder.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.FindByPath(ctx, path)
			if err != nil {
				t.Errorf("folderQueryBuilder.Find() error = %v", err)
			}

			assert.Equal(copy, *s)

			return
		})
	}
}

func makeFolderWithID(index int) *file.Folder {
	ret := makeFolder(index)
	ret.ID = folderIDs[index]

	return &ret
}

func Test_folderQueryBuilder_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		return getFolderPath(index)
	}

	tests := []struct {
		name    string
		path    string
		want    *file.Folder
		wantErr bool
	}{
		{
			"valid",
			getPath(folderIdxWithFiles),
			makeFolderWithID(folderIdxWithFiles),
			false,
		},
		{
			"invalid",
			"invalid path",
			nil,
			false,
		},
	}

	qb := db.Folder

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			got, err := qb.FindByPath(ctx, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("folderQueryBuilder.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("folderQueryBuilder.FindByPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
