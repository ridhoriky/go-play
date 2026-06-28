package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

func generateSlug(name string) string {
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := strings.Trim(re.ReplaceAllString(strings.ToLower(name), "-"), "-")
	if slug == "" {
		slug = "store"
	}
	return slug
}

func mapToResponse(s *entity.Store) *dto.StoreResponse {
	return &dto.StoreResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		StoreName:   s.StoreName,
		Slug:        s.Slug,
		Description: s.Description,
		LogoURL:     s.LogoURL,
		BannerURL:   s.BannerURL,
		IsVerified:  s.IsVerified,
		RatingAvg:   s.RatingAvg,
		TotalSales:  s.TotalSales,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func (s *storeService) CreateStore(ctx context.Context, userID string, req *dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	// Check if user already has a store
	existing, err := s.storeRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, dto.NewError(http.StatusBadRequest, "User already has a store")
	}

	slug := generateSlug(req.StoreName)
	exists, err := s.storeRepo.IsSlugExists(ctx, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		slug = fmt.Sprintf("%s-%s", slug, xid.New().String()[:4])
	}

	storeID, _ := uuid.NewV7()

	newStore := &entity.Store{
		ID:          storeID.String(),
		UserID:      userID,
		StoreName:   req.StoreName,
		Slug:        slug,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err = s.storeRepo.Create(ctx, newStore); err != nil {
		return nil, err
	}

	// Update user role to seller
	u, err := s.userRepo.GetByID(ctx, userID)
	if err == nil && u.Role != "admin" { // don't demote admin
		u.Role = "seller"
		_ = s.userRepo.Update(ctx, userID, u)
	}

	return mapToResponse(newStore), nil
}

func (s *storeService) GetMyStore(ctx context.Context, userID string) (*dto.StoreResponse, error) {
	store, err := s.storeRepo.GetByUserID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewError(http.StatusNotFound, "Store not found")
	}
	if err != nil {
		return nil, err
	}
	return mapToResponse(store), nil
}

func (s *storeService) UpdateStore(ctx context.Context, userID string, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	store, err := s.storeRepo.GetByUserID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewError(http.StatusNotFound, "Store not found")
	}
	if err != nil {
		return nil, err
	}

	store.StoreName = req.StoreName
	store.Description = req.Description
	store.LogoURL = req.LogoURL
	store.BannerURL = req.BannerURL

	if err = s.storeRepo.Update(ctx, store); err != nil {
		return nil, err
	}

	updated, err := s.storeRepo.GetByID(ctx, store.ID)
	if err != nil {
		return nil, err
	}
	return mapToResponse(updated), nil
}

func (s *storeService) GetStoreBySlug(ctx context.Context, slug string) (*dto.StoreResponse, error) {
	store, err := s.storeRepo.GetBySlug(ctx, slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewError(http.StatusNotFound, "Store not found")
	}
	if err != nil {
		return nil, err
	}
	return mapToResponse(store), nil
}

func (s *storeService) ListStores(ctx context.Context, params *dto.GetStoresQuery) (*dto.StoreListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	stores, total, err := s.storeRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	var res []dto.StoreResponse
	for i := range stores {
		st := stores[i]
		res = append(res, *mapToResponse(&st))
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.StoreListResponse{
		Data: res,
		Meta: dto.PaginationMeta{
			Page:       params.Page,
			Limit:      params.Limit,
			TotalPages: totalPages,
			Total:      total,
		},
	}, nil
}

func (s *storeService) GetStoreStats(ctx context.Context, storeID string) (*dto.StoreStats, error) {
	return s.storeRepo.GetStoreStats(ctx, storeID)
}
