package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
)

type AdminService interface {
	CreateFranchise(ctx context.Context, franchise *model.Franchise, f_owner *model.FranchiseOwner) (*model.FranchiseResponse, error)
	UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.FranchiseResponse, error)
	UpdateFranchiseStatus(ctx context.Context, id string, status string) (string, error)
	DeleteFranchise(ctx context.Context, id string) (bool, error)
	GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error)
}

type adminService struct {
	repo repository.AdminRepository // use interface, not pointer to struct
	// authorizer Authorizer (optionally inject this if using a real AuthZ client)
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{
		repo: repo,
	}
}

// ---- AuthZ stub (replace with real call to AuthZ service) ----
func authorize(ctx context.Context, action, resource string) error {
	// Replace with real logic: extract user/role from ctx and check permissions
	// e.g. call authZService.Check(ctx, subject, action, resource)
	return nil // return error if unauthorized
}

// ---------------------------------------------------------------

func (ad *adminService) CreateFranchise(ctx context.Context, franchise *model.Franchise, f_owner *model.FranchiseOwner) (*model.FranchiseResponse, error) {
	if err := authorize(ctx, "create", "franchise"); err != nil {
		return nil, err
	}

	// 💡 Run validations before calling repo
	if err := validations.ValidateFranchise(franchise); err != nil {
		return nil, err
	}
	if err := validations.ValidateFranchiseOwner(f_owner); err != nil {
		return nil, err
	}

	franchiseID, ownerID, err := ad.repo.CreateFranchise(ctx, franchise, f_owner)
	if err != nil {
		return nil, err
	}

	return &model.FranchiseResponse{
		FranchiseID: franchiseID,
		OwnerID:     ownerID,
	}, nil
}

func (ad *adminService) UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.FranchiseResponse, error) {
	if err := authorize(ctx, "update", "franchise"); err != nil {
		return nil, err
	}

	rowsAffected, err := ad.repo.UpdateFranchise(ctx, id, franchise)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("no rows updated")
	}

	return &model.FranchiseResponse{
		FranchiseID: id,
	}, nil
}

func (ad *adminService) UpdateFranchiseStatus(ctx context.Context, id string, status string) (string, error) {
	if err := authorize(ctx, "update_status", "franchise"); err != nil {
		return "", err
	}

	err := ad.repo.UpdateFranchiseStatus(ctx, id, status)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("franchise status updated to: %s", status), nil
}

func (ad *adminService) DeleteFranchise(ctx context.Context, id string) (bool, error) {
	if err := authorize(ctx, "delete", "franchise"); err != nil {
		return false, err
	}

	rowsAffected, err := ad.repo.DeleteFranchise(ctx, id)
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (ad *adminService) GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error) {
	if err := authorize(ctx, "read", "franchise"); err != nil {
		return nil, err
	}

	return ad.repo.GetAllFranchises(ctx, page, limit)
}
