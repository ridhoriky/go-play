package store

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type StoreRepositoryItf interface {
	Create(ctx context.Context, store *entity.Store) error
	GetByID(ctx context.Context, id string) (*entity.Store, error)
	GetByUserID(ctx context.Context, userID string) (*entity.Store, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Store, error)
	Update(ctx context.Context, store *entity.Store) error
	SoftDelete(ctx context.Context, id string) error
	List(ctx context.Context, params *dto.GetStoresQuery) ([]entity.Store, int, error)
	IsSlugExists(ctx context.Context, slug string) (bool, error)
	GetStoreStats(ctx context.Context, storeID string) (*dto.StoreStats, error)
}

type storeRepository struct {
	db *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) StoreRepositoryItf {
	return &storeRepository{
		db: db,
	}
}
