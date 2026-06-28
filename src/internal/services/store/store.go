package store

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/repositories/store"
	"ne-project/src/internal/repositories/user"
)

type StoreServiceItf interface {
	CreateStore(ctx context.Context, userID string, req *dto.CreateStoreRequest) (*dto.StoreResponse, error)
	GetMyStore(ctx context.Context, userID string) (*dto.StoreResponse, error)
	UpdateStore(ctx context.Context, userID string, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error)
	GetStoreBySlug(ctx context.Context, slug string) (*dto.StoreResponse, error)
	ListStores(ctx context.Context, params *dto.GetStoresQuery) (*dto.StoreListResponse, error)
	GetStoreStats(ctx context.Context, storeID string) (*dto.StoreStats, error)
}

type storeService struct {
	storeRepo store.StoreRepositoryItf
	userRepo  user.UserRepositoryItf
}

func NewStoreService(storeRepo store.StoreRepositoryItf, userRepo user.UserRepositoryItf) StoreServiceItf {
	return &storeService{
		storeRepo: storeRepo,
		userRepo:  userRepo,
	}
}
