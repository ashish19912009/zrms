package mapper

import (
	"errors"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/pb"
)

func AddFranchiseOwner_ProtoToModel(pbOwner *pb.FranchiseOwnerInput) (*model.FranchiseOwner, error) {
	return &model.FranchiseOwner{
		Name:       pbOwner.Name,
		Gender:     pbOwner.Gender,
		Dob:        pbOwner.Dob,
		MobileNo:   pbOwner.MobileNo,
		Email:      pbOwner.Email,
		Address:    pbOwner.Address,
		AadharNo:   pbOwner.AadharNo,
		IsVerified: pbOwner.IsVerified,
	}, nil
}

func Update_ModelToProto(owner *model.UpdateResponse) *pb.UpdateResponse {
	return &pb.UpdateResponse{
		Id:        owner.ID,
		UpdatedAt: owner.UpdatedAt.Format(time.RFC3339),
	}
}

func Delete_ModelToProto(owner *model.DeletedResponse) *pb.DeletedResponse {
	return &pb.DeletedResponse{
		Id:        owner.ID,
		DeletedAt: owner.DeletedAt.Format(time.RFC3339),
	}
}

func AddFranchiseStatus_FromProtoToModel(id string, status string) (*model.FranchiseStatusRequest, error) {
	if id == "" || status == "" {
		return nil, errors.New("Id or status can't be empty")
	}
	return &model.FranchiseStatusRequest{
		ID:     id,
		Status: status,
	}, nil
}

func DeleteFranchise_FromProtoToModel(id string, admin_id string) (*model.DeleteFranchiseRequest, error) {
	if id == "" || admin_id == "" {
		return nil, errors.New("Id or admin id can't be empty")
	}
	return &model.DeleteFranchiseRequest{
		ID:      id,
		AdminID: admin_id,
	}, nil
}

func AddFranchise_ProtoToModel(pbFranchise *pb.FranchiseInput) (*model.Franchise, error) {
	return &model.Franchise{
		BusinessName:       pbFranchise.BusinessName,
		LogoURL:            pbFranchise.LogoUrl,
		SubDomain:          pbFranchise.Subdomain,
		ThemeSettings:      pbFranchise.ThemeSettings.AsMap(),
		Status:             pbFranchise.Status,
		Franchise_Owner_id: pbFranchise.FranchiseOwnerId,
	}, nil
}

func GetAllFranchises_ProtoToModel(pagination *pb.PaginationRequest, query string) (int32, int32, string, error) {
	if pagination.Page == 0 || pagination.Limit == 0 {
		return 0, 0, "", errors.New("page and limit can't be empty")
	}
	return pagination.Page, pagination.Limit, query, nil
}

func GetAllFranchises_ModelToProto(page, limit int32, allFranchise []model.FranchiseResponse) *pb.GetFranchisesResponse {
	var protoFranchise []*pb.FranchiseByIDInput
	for _, franchise := range allFranchise {
		var f_details = &pb.FranchiseInput{
			BusinessName: franchise.BusinessName,
			LogoUrl:      franchise.LogoURL,
			Subdomain:    franchise.SubDomain,
			Status:       franchise.Status,
		}
		protoFranchise = append(protoFranchise, &pb.FranchiseByIDInput{
			Id:               franchise.ID,
			FranchiseDetails: f_details,
			CreatedAt:        franchise.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        franchise.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &pb.GetFranchisesResponse{
		Franchises: protoFranchise,
		Pagination: &pb.PaginationResponse{
			Page:  page,
			Limit: limit,
		},
	}
}

func Add_ModelToProto(franchise *model.AddResponse) *pb.AddResponse {
	return &pb.AddResponse{
		Id:        franchise.ID,
		CreatedAt: franchise.CreatedAt.Format(time.RFC3339),
	}
}

func UpdateFranchiseToProto(franchise *model.UpdateResponse) *pb.UpdateResponse {
	return &pb.UpdateResponse{
		Id:        franchise.ID,
		UpdatedAt: franchise.UpdatedAt.Format(time.RFC3339),
	}
}
