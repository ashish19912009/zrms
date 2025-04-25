package service

import (
	"context"
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
)

type AdminService interface {
	CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error)
	UpdateNewOwner(ctx context.Context, id string, owner *model.FranchiseOwner) (*model.UpdateResponse, error)

	CreateFranchise(ctx context.Context, franchise *model.Franchise) (*model.AddResponse, error)
	UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.UpdateResponse, error)
	UpdateFranchiseStatus(ctx context.Context, f_status *model.FranchiseStatusRequest) (*model.UpdateResponse, error)
	DeleteFranchise(ctx context.Context, f_del *model.DeleteFranchiseRequest) (*model.DeletedResponse, error)
	GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error)
}

type adminService struct {
	a_repo repository.AdminRepository // use interface, not pointer to struct+
	repo   repository.Repository
	// authorizer Authorizer (optionally inject this if using a real AuthZ client)
}

func NewAdminService(admin_repo repository.AdminRepository, repo repository.Repository) AdminService {
	return &adminService{
		a_repo: admin_repo,
		repo:   repo,
	}
}

// ---- AuthZ stub (replace with real call to AuthZ service) ----
func authorize(ctx context.Context, action, resource string) error {
	// Replace with real logic: extract user/role from ctx and check permissions
	// e.g. call authZService.Check(ctx, subject, action, resource)
	return nil // return error if unauthorized
}

// ---------------------------------------------------------------

func (ad *adminService) CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error) {
	if err := authorize(ctx, "create", "franchise"); err != nil {
		return nil, err
	}

	// check if business name already registered or not
	ownerExist, err := ad.repo.CheckIfOwnerExistsByAadharID(ctx, owner.AadharNo)
	if err != nil {
		return nil, err
	}

	var newFranchiseOwner *model.AddResponse
	if ownerExist != nil {
		return nil, errors.New(constants.FranchiseOwnerExist)
	} else {
		newFranchiseOwner, err = ad.a_repo.CreateNewOwner(ctx, owner)
		if err != nil {
			return nil, err
		}
	}
	return newFranchiseOwner, nil
}

func (ad *adminService) UpdateNewOwner(ctx context.Context, id string, owner *model.FranchiseOwner) (*model.UpdateResponse, error) {
	if err := authorize(ctx, "create", "franchise"); err != nil {
		return nil, err
	}

	// 💡 Run validations before calling repo
	franchiseExist, err := ad.repo.GetFranchiseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if franchiseExist == nil {
		return nil, errors.New("No owner to update")
	}

	newFranchiseOwner, err := ad.a_repo.UpdateNewOwner(ctx, id, owner)
	if err != nil {
		return nil, err
	}
	return newFranchiseOwner, nil
}

func (ad *adminService) CreateFranchise(ctx context.Context, franchise *model.Franchise) (*model.AddResponse, error) {
	if err := authorize(ctx, "create", "franchise"); err != nil {
		return nil, err
	}

	// 💡 Run validations before calling repo
	// check if business name already registered or not
	franchiseExists, err := ad.repo.GetFranchiseByBusinessName(ctx, franchise.BusinessName)
	if err != nil {
		return nil, err
	}

	var newFranchise *model.AddResponse
	if franchiseExists != nil && franchiseExists.BusinessName == franchise.BusinessName && franchiseExists.Franchise_Owner_id == franchise.Franchise_Owner_id {
		return nil, errors.New(constants.BusinessAlreadyExist)
	} else {
		newFranchise, err = ad.a_repo.CreateFranchise(ctx, franchise)
		if err != nil {
			return nil, err
		}
	}
	return newFranchise, nil
}

func (ad *adminService) UpdateFranchise(ctx context.Context, id string, franchise *model.Franchise) (*model.UpdateResponse, error) {
	if err := authorize(ctx, "update", "franchise"); err != nil {
		return nil, err
	}

	// 💡 Run validations before calling repo
	if err := validations.ValidateFranchise(franchise); err != nil {
		return nil, err
	}
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}

	rowsAffected, err := ad.a_repo.UpdateFranchise(ctx, id, franchise)
	if err != nil {
		return nil, err
	}

	if rowsAffected == nil {
		return nil, errors.New("no rows updated")
	}

	return rowsAffected, nil
}

func (ad *adminService) UpdateFranchiseStatus(ctx context.Context, f_status *model.FranchiseStatusRequest) (*model.UpdateResponse, error) {
	if err := authorize(ctx, "update_status", "franchise"); err != nil {
		return nil, err
	}

	updated_completed, err := ad.a_repo.UpdateFranchiseStatus(ctx, f_status.ID, f_status.Status)
	if err != nil {
		return nil, err
	}

	if updated_completed != nil {
		return &model.UpdateResponse{
			ID:        updated_completed.ID,
			UpdatedAt: updated_completed.UpdatedAt,
		}, nil
	}
	return nil, nil
}

func (ad *adminService) DeleteFranchise(ctx context.Context, f_del *model.DeleteFranchiseRequest) (*model.DeletedResponse, error) {
	// check if the admin is authorized to delete user or not
	if err := authorize(ctx, "delete", "franchise"); err != nil {
		return nil, err
	}
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(f_del.ID); err != nil {
		return nil, err
	}

	rowsAffected, err := ad.a_repo.DeleteFranchise(ctx, f_del.ID)
	if err != nil {
		return nil, err
	}

	return &model.DeletedResponse{
		ID:        rowsAffected.ID,
		DeletedAt: rowsAffected.DeletedAt,
	}, nil
}

func (ad *adminService) GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error) {
	if err := authorize(ctx, "read", "franchise"); err != nil {
		return nil, err
	}
	return ad.a_repo.GetAllFranchises(ctx, page, limit)
}
