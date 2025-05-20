package mapper

import (
	"errors"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/dbutils"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/pb"
	"google.golang.org/protobuf/types/known/structpb"
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

func GetFranchiseByID_FromProtoToModel(id string) string {
	return id
}

func GetFranchiseByID_ModelToProto(franc *model.FranchiseResponse) (*pb.GetFranchiseByIDResponse, error) {
	if franc == nil {
		return &pb.GetFranchiseByIDResponse{}, nil
	}

	// Ensure ThemeSettings is never nil
	if franc.ThemeSettings == nil {
		franc.ThemeSettings = make(map[string]interface{})
	}

	settingsStruct, err := structpb.NewStruct(franc.ThemeSettings)
	if err != nil {
		return nil, err
	}
	frachise := &pb.FranchiseInput{
		BusinessName:     franc.BusinessName,
		LogoUrl:          franc.LogoURL,
		Subdomain:        franc.SubDomain,
		ThemeSettings:    settingsStruct,
		Status:           franc.Status,
		FranchiseOwnerId: franc.Franchise_Owner_id,
	}
	frByID := &pb.FranchiseByIDInput{
		Id:               franc.ID,
		FranchiseDetails: frachise,
		CreatedAt:        franc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        franc.UpdatedAt.Format(time.RFC3339),
	}
	return &pb.GetFranchiseByIDResponse{
		Franchise: frByID,
	}, nil
}

func AddFranchise_ProtoToModel(pbFranchise *pb.AddFranchiseRequest) (*model.Franchise, error) {
	frach := pbFranchise.GetFranchiseDetails()
	jsonBytes, err := dbutils.ConvertStringMapToJson(frach.ThemeSettings.AsMap())
	if err != nil {
		logger.Error("something went wrong while converting map string to JSON", err, nil)
	}
	return &model.Franchise{
		BusinessName:       frach.BusinessName,
		LogoURL:            frach.LogoUrl,
		SubDomain:          frach.Subdomain,
		ThemeSettings:      string(jsonBytes),
		Status:             frach.Status,
		Franchise_Owner_id: frach.FranchiseOwnerId,
	}, nil
}

func FranchiseOwner_ModelToProto(res *model.FranchiseOwnerResponse) *pb.GetFranchiseOwnerResponse {
	if res == nil {
		return nil
	}

	fOwner := &pb.FranchiseOwnerInput{
		Name:       res.Name,
		Gender:     res.Gender,
		Dob:        res.Dob,
		MobileNo:   res.MobileNo,
		Email:      res.Email,
		Address:    res.Address,
		AadharNo:   res.AadharNo,
		IsVerified: res.IsVerified,
		Status:     res.Status,
	}

	var createdAtStr, updatedAtStr string
	if res.CreatedAt != nil {
		createdAtStr = res.CreatedAt.Format(time.RFC3339)
	}
	if res.UpdatedAt != nil {
		updatedAtStr = res.UpdatedAt.Format(time.RFC3339)
	}

	return &pb.GetFranchiseOwnerResponse{
		Id:        res.ID,
		FOwner:    fOwner,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func UpdateFranchise_ProtoToModel(pbFranchise *pb.UpdateFranchiseRequest) (*model.Franchise, error) {
	frach := pbFranchise.GetFranchiseDetails()
	jsonBytes, err := dbutils.ConvertStringMapToJson(frach.ThemeSettings.AsMap())
	if err != nil {
		logger.Error("something went wrong while converting map string to JSON", err, nil)
	}
	return &model.Franchise{
		ID:                 pbFranchise.Id,
		BusinessName:       frach.BusinessName,
		LogoURL:            frach.LogoUrl,
		SubDomain:          frach.Subdomain,
		ThemeSettings:      string(jsonBytes),
		Status:             frach.Status,
		Franchise_Owner_id: frach.FranchiseOwnerId,
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
	if franchise == nil {
		return nil
	}
	var createdAt string
	if !franchise.CreatedAt.IsZero() {
		createdAt = franchise.CreatedAt.Format(time.RFC3339)
	}
	return &pb.AddResponse{
		Id:        franchise.ID,
		CreatedAt: createdAt,
	}
}

func UpdateFranchiseToProto(franchise *model.UpdateResponse) *pb.UpdateResponse {
	return &pb.UpdateResponse{
		Id:        franchise.ID,
		UpdatedAt: franchise.UpdatedAt.Format(time.RFC3339),
	}
}
