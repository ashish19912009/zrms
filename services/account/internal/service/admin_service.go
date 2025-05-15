package service

import (
	"context"
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminService interface {
	CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error)
	UpdateOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.UpdateResponse, error)

	CreateFranchise(ctx context.Context, franchise *model.Franchise) (*model.AddResponse, error)
	UpdateFranchise(ctx context.Context, franchise *model.Franchise) (*model.UpdateResponse, error)
	UpdateFranchiseStatus(ctx context.Context, f_status *model.FranchiseStatusRequest) (*model.UpdateResponse, error)
	DeleteFranchise(ctx context.Context, f_del *model.DeleteFranchiseRequest) (*model.DeletedResponse, error)
	GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error)
}

type adminService struct {
	a_repo repository.AdminRepository // use interface, not pointer to struct+
	repo   repository.Repository
	client client.AuthZClient
	// authorizer Authorizer (optionally inject this if using a real AuthZ client)
}

func NewAdminService(admin_repo repository.AdminRepository, repo repository.Repository, client client.AuthZClient) (AdminService, error) {
	if admin_repo == nil || repo == nil || client == nil {
		return nil, errors.New("fields required")
	}
	return &adminService{
		a_repo: admin_repo,
		repo:   repo,
		client: client,
	}, nil
}

func (ad *adminService) CreateNewOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.AddResponse, error) {
	// check if owner already exists
	ownerExist, err := ad.repo.CheckIfOwnerExistsByAadharID(ctx, owner.AadharNo)
	if err != nil {
		return nil, err
	}
	if ownerExist != nil {
		return nil, errors.New(constants.FranchiseOwnerExist)
	}
	newFranchiseOwner, err := ad.a_repo.CreateNewOwner(ctx, owner)
	if err != nil {
		return nil, err
	}
	return newFranchiseOwner, nil
}

func (ad *adminService) UpdateOwner(ctx context.Context, owner *model.FranchiseOwner) (*model.UpdateResponse, error) {
	// 💡 Run validations before calling repo
	franchiseExist, err := ad.repo.GetFranchiseByID(ctx, owner.ID)
	if err != nil {
		return nil, err
	}
	if franchiseExist == nil {
		return nil, errors.New("No owner to update")
	}

	updateFOwner, err := ad.a_repo.UpdateNewOwner(ctx, owner)
	if err != nil {
		return nil, err
	}
	return updateFOwner, nil
}

func (ad *adminService) CreateFranchise(ctx context.Context, franchise *model.Franchise) (*model.AddResponse, error) {
	// check if business name already registered or not
	franchiseExists, err := ad.repo.GetFranchiseByBusinessName(ctx, franchise.BusinessName)
	if err != nil {
		return nil, err
	}
	if franchiseExists != nil && franchiseExists.BusinessName == franchise.BusinessName && franchiseExists.Franchise_Owner_id == franchise.Franchise_Owner_id {
		return nil, status.Error(codes.PermissionDenied, constants.BusinessAlreadyExist)
	}
	var newFranchise *model.AddResponse
	uuid := uuid.New().String()
	franchise.ID = uuid
	newFranchise, err = ad.a_repo.CreateFranchise(ctx, franchise)
	if err != nil {
		return nil, err
	}
	return newFranchise, nil
}

func (ad *adminService) UpdateFranchise(ctx context.Context, franchise *model.Franchise) (*model.UpdateResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateFranchise(franchise); err != nil {
		return nil, err
	}

	if err := validations.ValidateUUID(franchise.ID); err != nil {
		return nil, err
	}

	rowsAffected, err := ad.a_repo.UpdateFranchise(ctx, franchise)
	if err != nil {
		return nil, err
	}

	if rowsAffected == nil {
		return nil, errors.New("no rows updated")
	}

	return rowsAffected, nil
}

func (ad *adminService) UpdateFranchiseStatus(ctx context.Context, f_status *model.FranchiseStatusRequest) (*model.UpdateResponse, error) {
	if err := validations.ValidateUUID(f_status.ID); err != nil {
		return nil, err
	}
	franchiseExist, err := ad.repo.GetFranchiseByID(ctx, f_status.ID)
	if err != nil {
		return nil, err
	}
	if franchiseExist == nil {
		return nil, errors.New("No Franchise found")
	}

	updated_completed, err := ad.a_repo.UpdateFranchiseStatus(ctx, f_status.ID, f_status.Status)
	if err != nil {
		return nil, err
	}
	return updated_completed, nil
}

func (ad *adminService) DeleteFranchise(ctx context.Context, f_del *model.DeleteFranchiseRequest) (*model.DeletedResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(f_del.ID); err != nil {
		return nil, err
	}
	franchiseExist, err := ad.repo.GetFranchiseByID(ctx, f_del.ID)
	if err != nil {
		return nil, err
	}
	if franchiseExist == nil {
		return nil, errors.New("No Franchise found")
	}
	rowsAffected, err := ad.a_repo.DeleteFranchise(ctx, f_del.ID)
	if err != nil {
		return nil, err
	}
	return rowsAffected, nil
}

func (ad *adminService) GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error) {
	return ad.a_repo.GetAllFranchises(ctx, page, limit)
}
