package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/repository"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
)

/*
Responsibilities of service:
Validate input data (e.g., check if mobileNo is valid before inserting).

Interact with Repository to fetch/update accounts.

Use goroutines & WaitGroups where needed for parallel execution.

Handle errors properly.

Return structured responses.
*/
// mockery --name=AccountService --dir=services/account/internal/service --output=services/account/internal/service/mocks --outpkg=mocks

// AccountService defines business logic for accounts
type AccountService interface {
	GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error)
	GetFranchiseByBusinessName(ctx context.Context, b_name string) (*model.FranchiseResponse, error)
	GetFranchiseOwnerByID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error)
	CheckIfOwnerExistsByAadharID(ctx context.Context, id string) (bool, error)

	CreateFranchiseAccount(ctx context.Context, account *model.FranchiseAccount) (*model.AddResponse, error)
	UpdateFranchiseAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.UpdateResponse, error)
	GetFranchiseAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error)
	GetAllFranchiseAccounts(ctx context.Context, franDetail *model.GetFranchisesRequest) ([]model.FranchiseAccountResponse, error)

	AddFranchiseDocument(ctx context.Context, doc *model.FranchiseDocument) (*model.AddResponse, error)
	UpdateFranchiseDocument(ctx context.Context, id string, doc *model.FranchiseDocument) (*model.UpdateResponse, error)
	GetAllFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponseComplete, error)

	AddFranchiseAddress(ctx context.Context, addr *model.FranchiseAddress) (*model.AddResponse, error)
	UpdateFranchiseAddress(ctx context.Context, id string, addr *model.FranchiseAddress) (*model.UpdateResponse, error)
	GetFranchiseAddressByID(ctx context.Context, id string) (*model.FranchiseAddressResponse, error)

	AddFranchiseRole(ctx context.Context, role *model.FranchiseRole) (*model.AddResponse, error)
	UpdateFranchiseRole(ctx context.Context, id string, role *model.FranchiseRole) (*model.UpdateResponse, error)
	GetAllFranchiseRoles(ctx context.Context, id string) ([]model.FranchiseRoleResponse, error)

	AddPermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error)
	UpdatePermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error)
	GetAllPermissionsToRole(ctx context.Context, id string) ([]model.RoleToPermissionsComplete, error)
}

// accountService implements AccountService
type accountService struct {
	repo   repository.Repository
	client client.AuthZClient
}

// NewAccountService initializes service with a repository
func NewAccountService(repo repository.Repository, client client.AuthZClient) AccountService {
	return &accountService{
		repo:   repo,
		client: client,
	}
}

func (aS *accountService) GetFranchiseByID(ctx context.Context, id string) (*model.FranchiseResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}
	franchise, err := aS.repo.GetFranchiseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return franchise, nil
}

func (aS *accountService) GetFranchiseByBusinessName(ctx context.Context, b_name string) (*model.FranchiseResponse, error) {
	if err := validations.ValidateLength(b_name, 5, 100); err != nil {
		return nil, err
	}

	franchise, err := aS.repo.GetFranchiseByBusinessName(ctx, b_name)
	if err != nil {
		return nil, err
	}
	return franchise, nil
}

func (aS *accountService) GetFranchiseOwnerByID(ctx context.Context, id string) (*model.FranchiseOwnerResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}
	f_owner, err := aS.repo.GetFranchiseOwnerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return f_owner, nil
}

func (aS *accountService) CheckIfOwnerExistsByAadharID(ctx context.Context, aadharID string) (bool, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateAadhaarNumber(aadharID); err != nil {
		return false, err
	}
	f_owner, err := aS.repo.CheckIfOwnerExistsByAadharID(ctx, aadharID)
	if err != nil {
		return false, err
	}
	if f_owner == true {
		return false, nil
	}
	return true, nil
}

func (aS *accountService) CreateFranchiseAccount(ctx context.Context, account *model.FranchiseAccount) (*model.AddResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateFranchiseAccounts(account); err != nil {
		return nil, err
	}
	f_account, err := aS.repo.CreateFranchiseAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	return f_account, nil
}

func (aS *accountService) UpdateFranchiseAccount(ctx context.Context, id string, account *model.FranchiseAccount) (*model.UpdateResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}
	if err := validations.ValidateFranchiseAccounts(account); err != nil {
		return nil, err
	}
	f_owner, err := aS.repo.UpdateFranchiseAccount(ctx, id, account)
	if err != nil {
		return nil, err
	}
	return f_owner, nil
}
func (aS *accountService) GetFranchiseAccountByID(ctx context.Context, id string) (*model.FranchiseAccountResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}

	f_account, err := aS.repo.GetFranchiseAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return f_account, nil
}
func (aS *accountService) GetAllFranchiseAccounts(ctx context.Context, fran *model.GetFranchisesRequest) ([]model.FranchiseAccountResponse, error) {
	if err := validations.ValidateUUID(fran.FranchiseID); err != nil {
		return nil, err
	}
	franchiseExist, err := aS.repo.GetFranchiseByID(ctx, fran.FranchiseID)
	if err != nil {
		return nil, err
	}
	if franchiseExist == nil {
		return nil, errors.New("Frachise not found")
	}
	f_accounts, err := aS.repo.GetAllFranchiseAccounts(ctx, fran)
	if err != nil {
		return nil, err
	}
	return f_accounts, nil
}

func (aS *accountService) AddFranchiseDocument(ctx context.Context, doc *model.FranchiseDocument) (*model.AddResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateFranchiseDocument(doc); err != nil {
		return nil, err
	}
	franchiseExist, err := aS.repo.GetFranchiseByID(ctx, doc.FranchiseID)
	if err != nil {
		return nil, err
	}
	if franchiseExist == nil {
		return nil, errors.New("Frachise not found")
	}
	f_doc, err := aS.repo.AddFranchiseDocument(ctx, doc)
	if err != nil {
		return nil, err
	}
	return f_doc, nil
}

func (aS *accountService) UpdateFranchiseDocument(ctx context.Context, id string, doc *model.FranchiseDocument) (*model.UpdateResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}
	if err := validations.ValidateFranchiseDocument(doc); err != nil {
		return nil, err
	}
	f_doc, err := aS.repo.UpdateFranchiseDocument(ctx, id, doc)
	if err != nil {
		return nil, err
	}
	return f_doc, nil
}

func (aS *accountService) GetAllFranchiseDocuments(ctx context.Context, id string) ([]model.FranchiseDocumentResponseComplete, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}

	f_docs, err := aS.repo.GetAllFranchiseDocuments(ctx, id)
	if err != nil {
		return nil, err
	}
	return f_docs, nil
}

func (aS *accountService) AddFranchiseAddress(ctx context.Context, addr *model.FranchiseAddress) (*model.AddResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(addr.FranchiseID); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.AddressLine, 3, 200); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.City, 3, 100); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.Country, 3, 100); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.State, 3, 100); err != nil {
		return nil, err
	}
	if err := validations.ValidatePincode(addr.Pincode); err != nil {
		return nil, err
	}
	convertedLat, err := strconv.ParseFloat(addr.Latitude, 64)
	if err != nil {
		return nil, err
	}
	convertedLong, err := strconv.ParseFloat(addr.Longitude, 64)
	if err != nil {
		return nil, err
	}
	if err := validations.ValidateCoordinates(convertedLat, convertedLong); err != nil {
		return nil, err
	}

	f_addr, err := aS.repo.AddFranchiseAddress(ctx, addr)
	if err != nil {
		return nil, err
	}
	return f_addr, nil
}
func (aS *accountService) UpdateFranchiseAddress(ctx context.Context, id string, addr *model.FranchiseAddress) (*model.UpdateResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}
	if err := validations.ValidateUUID(addr.FranchiseID); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.AddressLine, 3, 200); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.City, 3, 100); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.Country, 3, 100); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(addr.State, 3, 100); err != nil {
		return nil, err
	}
	if err := validations.ValidatePincode(addr.Pincode); err != nil {
		return nil, err
	}
	convertedLat, err := strconv.ParseFloat(addr.Latitude, 64)
	if err != nil {
		return nil, err
	}
	convertedLong, err := strconv.ParseFloat(addr.Longitude, 64)
	if err != nil {
		return nil, err
	}
	if err := validations.ValidateCoordinates(convertedLat, convertedLong); err != nil {
		return nil, err
	}

	f_addr, err := aS.repo.UpdateFranchiseAddress(ctx, id, addr)
	if err != nil {
		return nil, err
	}
	return f_addr, nil
}
func (aS *accountService) GetFranchiseAddressByID(ctx context.Context, id string) (*model.FranchiseAddressResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}

	f_addr, err := aS.repo.GetFranchiseAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return f_addr, nil
}

func (aS *accountService) AddFranchiseRole(ctx context.Context, role *model.FranchiseRole) (*model.AddResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(role.FranchiseID); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(role.Name, 3, 200); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(role.Description, 3, 500); err != nil {
		return nil, err
	}

	f_role, err := aS.repo.AddFranchiseRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return f_role, nil
}
func (aS *accountService) UpdateFranchiseRole(ctx context.Context, id string, role *model.FranchiseRole) (*model.UpdateResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}
	if err := validations.ValidateUUID(role.FranchiseID); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(role.Name, 3, 200); err != nil {
		return nil, err
	}
	if err := validations.ValidateLength(role.Description, 3, 500); err != nil {
		return nil, err
	}

	f_role, err := aS.repo.UpdateFranchiseRole(ctx, id, role)
	if err != nil {
		return nil, err
	}
	return f_role, nil

}
func (aS *accountService) GetAllFranchiseRoles(ctx context.Context, id string) ([]model.FranchiseRoleResponse, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}

	f_roles, err := aS.repo.GetAllFranchiseRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	return f_roles, nil
}

func (aS *accountService) AddPermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(pRole.RoleID); err != nil {
		return nil, err
	}
	if err := validations.ValidateUUID(pRole.PermissionID); err != nil {
		return nil, err
	}

	p_role, err := aS.repo.AddPermissionsToRole(ctx, pRole)
	if err != nil {
		return nil, err
	}
	return p_role, nil
}
func (aS *accountService) UpdatePermissionsToRole(ctx context.Context, pRole *model.RoleToPermissions) (*model.RoleToPermissions, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(pRole.RoleID); err != nil {
		return nil, err
	}
	if err := validations.ValidateUUID(pRole.PermissionID); err != nil {
		return nil, err
	}

	p_role, err := aS.repo.UpdatePermissionsToRole(ctx, pRole)
	if err != nil {
		return nil, err
	}
	return p_role, nil
}
func (aS *accountService) GetAllPermissionsToRole(ctx context.Context, id string) ([]model.RoleToPermissionsComplete, error) {
	// 💡 Run validations before calling repo
	if err := validations.ValidateUUID(id); err != nil {
		return nil, err
	}

	p_roles, err := aS.repo.GetAllPermissionsToRole(ctx, id)
	if err != nil {
		return nil, err
	}
	return p_roles, nil
}
