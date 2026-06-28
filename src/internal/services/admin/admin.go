package admin

import (
	"context"

	"ne-project/src/internal/models/dto"
)

type AdminServiceItf interface {
	ListUsers(ctx context.Context, params *dto.AdminUserListParams) ([]dto.AdminUserResponse, int, error)
	GetUserByID(ctx context.Context, id string) (*dto.AdminUserResponse, error)
	UpdateUser(ctx context.Context, currentUserID string, targetUserID string, req dto.AdminUpdateUserRequest) error
	ListSellers(ctx context.Context, params *dto.AdminSellerListParams) ([]dto.AdminSellerResponse, int, error)
	VerifySeller(ctx context.Context, storeID string) error
	UnverifySeller(ctx context.Context, storeID string) error
}
