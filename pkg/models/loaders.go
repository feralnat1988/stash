//go:generate go run -mod=vendor github.com/vektah/dataloaden SceneLoader int *github.com/stashapp/stash/pkg/models.Scene
//go:generate go run -mod=vendor github.com/vektah/dataloaden PerformerLoader int *github.com/stashappp/stash/pkg/models.Performer
//go:generate go run -mod=vendor github.com/vektah/dataloaden TagLoader int *github.com/stashappp/stash/pkg/models.Tag
//go:generate go run -mod=vendor github.com/vektah/dataloaden StudioLoader int *github.com/stashappp/stash/pkg/models.Studio

package models

import (
	"context"
	"errors"
	"time"
)

var ErrTxn = errors.New("transaction failure")

func NewSceneLoaderConfig(ctx context.Context, mgr TransactionManager) SceneLoaderConfig {
	return SceneLoaderConfig{
		Wait:     2 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []int) ([]*Scene, []error) {
			var scenes []*Scene
			var errors []error
			err := mgr.WithReadTxn(ctx, func(repo ReaderRepository) error {
				// The right implementation here would be a FindMany which return
				// all errors it finds, and optimizes the query to be a batch query:
				// either via sql IN, or via a temporary table join. At the time of
				// this writing, however, FindMany is just a wrapper around Find(),
				// so we mimic FindMany here to grab all errors and satisfy the dataloaden
				// fetch interface.
				r := repo.Scene()

				for _, k := range keys {
					s, e := r.Find(k)
					scenes = append(scenes, s)
					errors = append(errors, e)
				}

				return nil
			})

			if err != nil {
				// Txn failure, build error result set
				scenes = nil
				errors = nil
				for range keys {
					scenes = append(scenes, nil)
					errors = append(errors, ErrTxn)
				}
			}

			return scenes, errors
		},
	}
}

func NewPerformerLoaderConfig(ctx context.Context, mgr TransactionManager) PerformerLoaderConfig {
	return PerformerLoaderConfig{
		Wait:     2 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []int) ([]*Performer, []error) {
			var performers []*Performer
			var errors []error
			err := mgr.WithReadTxn(ctx, func(repo ReaderRepository) error {
				// The right implementation here would be a FindMany which return
				// all errors it finds, and optimizes the query to be a batch query:
				// either via sql IN, or via a temporary table join. At the time of
				// this writing, however, FindMany is just a wrapper around Find(),
				// so we mimic FindMany here to grab all errors and satisfy the dataloaden
				// fetch interface.
				r := repo.Performer()

				for _, k := range keys {
					p, e := r.Find(k)
					performers = append(performers, p)
					errors = append(errors, e)
				}

				return nil
			})

			if err != nil {
				// Txn failure, build error result set
				performers = nil
				errors = nil
				for range keys {
					performers = append(performers, nil)
					errors = append(errors, ErrTxn)
				}
			}

			return performers, errors
		},
	}
}

func NewTagLoaderConfig(ctx context.Context, mgr TransactionManager) TagLoaderConfig {
	return TagLoaderConfig{
		Wait:     2 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []int) ([]*Tag, []error) {
			var tags []*Tag
			var errors []error
			err := mgr.WithReadTxn(ctx, func(repo ReaderRepository) error {
				// The right implementation here would be a FindMany which return
				// all errors it finds, and optimizes the query to be a batch query:
				// either via sql IN, or via a temporary table join. At the time of
				// this writing, however, FindMany is just a wrapper around Find(),
				// so we mimic FindMany here to grab all errors and satisfy the dataloaden
				// fetch interface.
				r := repo.Tag()

				for _, k := range keys {
					s, e := r.Find(k)
					tags = append(tags, s)
					errors = append(errors, e)
				}

				return nil
			})

			if err != nil {
				// Txn failure, build error result set
				tags = nil
				errors = nil
				for range keys {
					tags = append(tags, nil)
					errors = append(errors, ErrTxn)
				}
			}

			return tags, errors
		},
	}
}

func NewStudioLoaderConfig(ctx context.Context, mgr TransactionManager) StudioLoaderConfig {
	return StudioLoaderConfig{
		Wait:     2 * time.Millisecond,
		MaxBatch: 100,
		Fetch: func(keys []int) ([]*Studio, []error) {
			var studios []*Studio
			var errors []error
			err := mgr.WithReadTxn(ctx, func(repo ReaderRepository) error {
				// The right implementation here would be a FindMany which return
				// all errors it finds, and optimizes the query to be a batch query:
				// either via sql IN, or via a temporary table join. At the time of
				// this writing, however, FindMany is just a wrapper around Find(),
				// so we mimic FindMany here to grab all errors and satisfy the dataloaden
				// fetch interface.
				r := repo.Studio()

				for _, k := range keys {
					s, e := r.Find(k)
					studios = append(studios, s)
					errors = append(errors, e)
				}

				return nil
			})

			if err != nil {
				// Txn failure, build error result set
				studios = nil
				errors = nil
				for range keys {
					studios = append(studios, nil)
					errors = append(errors, ErrTxn)
				}
			}

			return studios, errors
		},
	}
}
