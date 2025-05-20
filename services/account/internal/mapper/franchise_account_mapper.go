package mapper

import (
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/pb"
)

func AddFranchiseAccount_ProtoToModel(req *pb.AddFranchiseAccountRequest) *model.FranchiseAccount {
	return &model.FranchiseAccount{
		FranchiseID: req.GetFranchiseId(),
		EmployeeID:  req.GetAccount().GetEmpId(),
		LoginID:     req.GetAccount().GetLoginId(),
		AccountType: req.GetAccount().GetAccountType(),
		Name:        req.GetAccount().GetName(),
		MobileNo:    req.GetAccount().GetMobileNo(),
		Email:       req.GetAccount().GetEmailId(),
		RoleID:      req.GetRoleId(),
		Status:      req.GetAccount().GetStatus(),
	}
}

func UpdateFranchiseAccount_ProtoToModel(req *pb.UpdateFranchiseAccountRequest) *model.FranchiseAccount {
	now := time.Now()

	return &model.FranchiseAccount{
		FranchiseID: req.GetFranchiseId(),
		EmployeeID:  req.GetAccount().GetEmpId(),
		LoginID:     req.GetAccount().GetLoginId(),
		AccountType: req.GetAccount().GetAccountType(),
		Name:        req.GetAccount().GetName(),
		MobileNo:    req.GetAccount().GetMobileNo(),
		Email:       req.GetAccount().GetEmailId(),
		Status:      req.GetAccount().GetStatus(),
		RoleID:      req.GetRoleId(),
		UpdatedAt:   &now,
	}
}

func GetFranchisesRequest_ProtoToModel(req *pb.GetFranchiseAccountsRequest) *model.GetFranchisesRequest {
	if req == nil {
		return nil
	}

	var pagination *model.Pagination
	if req.Pagination != nil {
		pagination = &model.Pagination{
			Page:  req.Pagination.Page,
			Limit: req.Pagination.Limit,
		}
	}

	var query *model.GetPaginationQuery
	query = &model.GetPaginationQuery{
		Pagination: pagination,
		Query:      req.Query,
	}

	return &model.GetFranchisesRequest{
		GetPagination: query,
		FranchiseID:   req.FranchiseId,
	}
}

func AddFranchiseDocument_ProtoToModel(req *pb.AddFranchiseDocumentRequest) *model.FranchiseDocument {
	if req == nil || req.FDoc == nil {
		return nil
	}

	return &model.FranchiseDocument{
		FranchiseID:    req.FranchiseId,
		DocumentTypeID: req.FDoc.DocName, // assuming doc_name maps to DocumentTypeID
		DocumentURL:    req.FDoc.DocumentUrl,
		UploadedBy:     req.FDoc.UploadedBy,
		Status:         req.FDoc.Status,
		Remark:         req.FDoc.Remarks,
	}
}

func UpdateFranchiseDocument_ProtoToModel(req *pb.UpdateFranchiseDocumentRequest) *model.FranchiseDocument {
	if req == nil || req.FDoc == nil {
		return nil
	}

	return &model.FranchiseDocument{
		FranchiseID:    req.FranchiseId,
		DocumentTypeID: req.FDoc.DocName, // Map doc_name as type ID (or change if you store actual doc_type_id separately)
		DocumentURL:    req.FDoc.DocumentUrl,
		UploadedBy:     req.FDoc.UploadedBy,
		Status:         req.FDoc.Status,
		Remark:         req.FDoc.Remarks,
		// You can map additional fields like doc_desc or is_mandate if your model supports it
	}
}

func AddFranchiseAddress_ProtoToModel(req *pb.AddFranchiseAddressRequest) *model.FranchiseAddress {
	if req == nil || req.Address == nil {
		return nil
	}

	return &model.FranchiseAddress{
		FranchiseID: req.FranchiseId,
		AddressLine: req.Address.AddressLine,
		City:        req.Address.City,
		State:       req.Address.State,
		Country:     req.Address.Country,
		Pincode:     req.Address.Pincode,
		Latitude:    req.Address.Latitude,
		Longitude:   req.Address.Longitude,
		IsVerified:  false, // default value, can be updated later
	}
}

func UpdateFranchiseAddress_ProtoToModel(req *pb.UpdateFranchiseAddressRequest) *model.FranchiseAddress {
	if req == nil || req.Address == nil {
		return nil
	}

	return &model.FranchiseAddress{
		FranchiseID: req.FranchiseId,
		AddressLine: req.Address.AddressLine,
		City:        req.Address.City,
		State:       req.Address.State,
		Country:     req.Address.Country,
		Pincode:     req.Address.Pincode,
		Latitude:    req.Address.Latitude,
		Longitude:   req.Address.Longitude,
		// IsVerified is not updated here, assumed to be handled separately
	}
}

func FranchiseAddressModelToProto(res *model.FranchiseAddressResponse, franchiseID string) *pb.GetFranchiseAddressResponse {
	if res == nil {
		return nil
	}

	return &pb.GetFranchiseAddressResponse{
		Id:          res.ID,
		FranchiseId: franchiseID,
		FAddress: &pb.FranchiseAddressInput{
			AddressLine: res.AddressLine,
			City:        res.City,
			State:       res.State,
			Country:     res.Country,
			Pincode:     res.Pincode,
			Latitude:    res.Latitude,
			Longitude:   res.Longitude,
		},
		CreatedAt: res.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: res.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func FranchiseRoleProtoToModel(req *pb.AddFranchiseRoleRequest) *model.FranchiseRole {
	return &model.FranchiseRole{
		FranchiseID: req.GetFranchiseId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		IsDefault:   req.GetIsDefault(),
	}
}

func MapFranchiseRoleInputToModel(input *pb.AddFranchiseRoleRequest) *model.FranchiseRole {
	if input == nil {
		return nil
	}

	return &model.FranchiseRole{
		FranchiseID: input.GetFranchiseId(),
		Name:        input.GetName(),
		Description: input.GetDescription(),
		IsDefault:   input.GetIsDefault(),
	}
}

func MapAddRolePermissionRequestToModel(req *pb.AddRolePermission) *model.RoleToPermissions {
	return &model.RoleToPermissions{
		RoleID:       req.GetRoleId(),
		PermissionID: req.GetPermissionId(),
	}
}

func MapRolePermissionToProto(p []model.RoleToPermissionsComplete) []*pb.RolePermissionDetails {
	result := make([]*pb.RolePermissionDetails, 0, len(p))
	for _, item := range p {
		result = append(result, &pb.RolePermissionDetails{
			FranchiseId:    item.FranchiseID,
			RoleName:       item.RoleName,
			RoleDesc:       item.Role_Description,
			IsDefault:      item.IsDefault,
			PermissionKey:  item.Permission_Key,
			PermissionDesc: item.Permission_Description,
			CreatedAt:      formatTime(item.CreatedAt),
			UpdatedAt:      formatTime(item.UpdatedAt),
		})
	}
	return result
}

// Helper to safely format time to RFC3339
func formatTime(t *time.Time) string {
	if t != nil {
		return t.Format(time.RFC3339)
	}
	return ""
}
