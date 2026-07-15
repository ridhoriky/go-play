package admin

import (
	"context"
	"errors"

	"ne-project/src/internal/models/dto"
	adminRepo "ne-project/src/internal/repositories/admin"
)

type AdminService struct {
	adminRepo adminRepo.AdminRepositoryItf
}

func NewAdminService(adminRepo adminRepo.AdminRepositoryItf) AdminServiceItf {
	return &AdminService{adminRepo: adminRepo}
}

func (s *AdminService) ListUsers(ctx context.Context, params *dto.AdminUserListParams) ([]dto.AdminUserResponse, int, error) {
	return s.adminRepo.ListUsers(ctx, params)
}

func (s *AdminService) GetUserByID(ctx context.Context, id string) (*dto.AdminUserResponse, error) {
	user, err := s.adminRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminService) UpdateUser(ctx context.Context, currentUserID string, targetUserID string, req dto.AdminUpdateUserRequest) error {
	user, err := s.adminRepo.GetUserByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if currentUserID == targetUserID {
		if req.Role != nil && *req.Role != user.Role {
			return errors.New("cannot change own role")
		}
		if req.IsActive != nil && !*req.IsActive && user.IsActive {
			return errors.New("cannot deactivate own account")
		}
	}

	if req.Role != nil {
		role := *req.Role
		if role != "admin" && role != "seller" && role != "user" {
			return errors.New("invalid role")
		}
	}

	return s.adminRepo.UpdateUser(ctx, targetUserID, req)
}

func (s *AdminService) ListSellers(ctx context.Context, params *dto.AdminSellerListParams) ([]dto.AdminSellerResponse, int, error) {
	return s.adminRepo.ListSellers(ctx, params)
}

func (s *AdminService) VerifySeller(ctx context.Context, storeID string) error {
	return s.adminRepo.VerifySeller(ctx, storeID)
}

func (s *AdminService) UnverifySeller(ctx context.Context, storeID string) error {
	return s.adminRepo.UnverifySeller(ctx, storeID)
}
