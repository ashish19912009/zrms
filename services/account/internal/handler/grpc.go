package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	"github.com/ashish19912009/zrms/services/account/internal/constants"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/mapper"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
	"github.com/ashish19912009/zrms/services/account/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	layer = constants.Handler
	erMsg = "unable to validate "
)

type GRPCHandler struct {
	pb.UnimplementedAccountServiceServer
	accountService service.AccountService
	adminService   service.AdminService
	client         client.AuthZClient
}

func NewGRPCHandler(
	accountService service.AccountService,
	adminService service.AdminService) *GRPCHandler {
	return &GRPCHandler{
		accountService: accountService,
		adminService:   adminService,
	}
}

// Helper to format *time.Time to RFC3339 string or empty if nil
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (h *GRPCHandler) CreateNewOwner(ctx context.Context, req *pb.AddFranchiseOwnerRequest) (*pb.AddResponse, error) {
	var method = constants.Methods.CreateNewOwner
	owner, err := mapper.AddFranchiseOwner_ProtoToModel(req.GetOwner())
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		errMsg := constants.MappingFromProtoToModel
		logger.Error(errMsg, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	// 💡 Run validations before calling service
	if err := validations.ValidateFranchiseOwner(owner); err != nil {
		return nil, err
	}
	createdOwner, err := h.adminService.CreateNewOwner(ctx, owner)
	if err != nil {
		eR := status.Errorf(codes.Internal, constants.FailedToCreateOwner, err)
		errMsg := constants.SomethinWentWrongOnNew + `new owner`
		logger.Error(errMsg, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	return mapper.Add_ModelToProto(createdOwner), nil
}

func (h *GRPCHandler) UpdateNewOwner(ctx context.Context, id string, req *pb.AddFranchiseOwnerRequest) (*pb.UpdateResponse, error) {
	var method = constants.Methods.UpdateOwner
	owner, err := mapper.AddFranchiseOwner_ProtoToModel(req.GetOwner())
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		errMsg := constants.MappingFromProtoToModel
		logger.Error(errMsg, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	// 💡 Run validations before calling service
	if err := validations.ValidateUUID(id); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateFranchiseOwner(owner); err != nil {
		logger.Error(erMsg+"owner fields", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	updateOwner, err := h.adminService.UpdateOwner(ctx, owner)
	if err != nil {
		eR := status.Errorf(codes.Internal, constants.FailedToCreateOwner, err)
		errMsg := constants.SomethinWentWrongOnUpdate + `existing owner`
		logger.Error(errMsg, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
	}
	return mapper.Update_ModelToProto(updateOwner), nil
}

func (h *GRPCHandler) CreateFranchise(ctx context.Context, req *pb.AddFranchiseRequest) (*pb.AddResponse, error) {
	var method = constants.Methods.CreateFranchise
	franchise, err := mapper.AddFranchise_ProtoToModel(req)
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		err := constants.MappingFromProtoToModel
		logger.Error(err, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	if err := validations.ValidateFranchise(franchise); err != nil {
		logger.Error(erMsg+" frachise fields", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	createdFranchise, err := h.adminService.CreateFranchise(ctx, franchise)
	if err != nil {
		eR := status.Errorf(codes.Internal, constants.FailedToCreateFranchsie, err)
		errMsg := constants.SomethinWentWrongOnNew + `new owner`
		logger.Error(errMsg, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	return mapper.Add_ModelToProto(createdFranchise), nil
}

func (h *GRPCHandler) UpdateFranchise(ctx context.Context, req *pb.UpdateFranchiseRequest) (*pb.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchise
	franchise, err := mapper.UpdateFranchise_ProtoToModel(req)
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		err := constants.MappingFromProtoToModel
		logger.Error(err, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	// 💡 Run validations before calling service
	if err := validations.ValidateUUID(req.Id); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateFranchise(franchise); err != nil {
		logger.Error(erMsg+" frachise fields", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	updatedFranchise, err := h.adminService.UpdateFranchise(ctx, franchise)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update franchise: %v", err)
	}
	return mapper.Update_ModelToProto(updatedFranchise), nil
}

func (h *GRPCHandler) UpdateFranchiseStatus(ctx context.Context, req *pb.UpdateFranchiseStatusRequest) (*pb.UpdateResponse, error) {
	// Convert proto to internal model
	var method = constants.Methods.UpdateFranchiseStatus
	// Convert proto to internal model
	franchise, err := mapper.AddFranchiseStatus_FromProtoToModel(req.GetId(), req.GetStatus())
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		err := constants.MappingFromProtoToModel
		logger.Error(err, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	// 💡 Run validations before calling service
	if err := validations.ValidateUUID(franchise.ID); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateStatus(franchise.Status); err != nil {
		logger.Error(erMsg+" frachise fields", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	updatedFranchise, err := h.adminService.UpdateFranchiseStatus(ctx, franchise)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update franchise: %v", err)
	}
	return mapper.Update_ModelToProto(updatedFranchise), nil
}

func (h *GRPCHandler) DeleteFranchise(ctx context.Context, req *pb.DeleteFranchiseRequest) (*pb.DeletedResponse, error) {
	var method = constants.Methods.DeleteFranchise
	franchise, err := mapper.DeleteFranchise_FromProtoToModel(req.GetId(), req.GetAdminId())
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		err := constants.MappingFromProtoToModel
		logger.Error(err, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	// 💡 Run validations before calling service
	if err := validations.ValidateUUID(franchise.ID); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(franchise.AdminID); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	updatedFranchise, err := h.adminService.DeleteFranchise(ctx, franchise)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update franchise: %v", err)
	}
	return mapper.Delete_ModelToProto(updatedFranchise), nil
}

func (h *GRPCHandler) GetAllFranchises(ctx context.Context, req *pb.GetFranchisesRequest) (*pb.GetFranchisesResponse, error) {
	var method = constants.Methods.GetAllFranchises
	page, limit, _, err := mapper.GetAllFranchises_ProtoToModel(req.GetPagination(), req.GetQuery())
	if err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.InvalidTimestamp, err)
		err := constants.MappingFromProtoToModel
		logger.Error(err, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}
	allFranchise, err := h.adminService.GetAllFranchises(ctx, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update franchise: %v", err)
	}
	return mapper.GetAllFranchises_ModelToProto(page, limit, allFranchise), nil
}

func (h *GRPCHandler) GetFranchiseByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetFranchiseByIDResponse, error) {
	var method = constants.Methods.GetFranchiseByID
	ID := mapper.GetFranchiseByID_FromProtoToModel(req.GetId())
	// 💡 Run validations before calling service
	if err := validations.ValidateUUID(ID); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(ID); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	franchise, err := h.accountService.GetFranchiseByID(ctx, ID)
	if err != nil {
		fmt.Printf("no Franchise Found")
	}
	// Critical nil check
	if franchise == nil {
		return &pb.GetFranchiseByIDResponse{}, nil
	}

	frByID, err := mapper.GetFranchiseByID_ModelToProto(franchise)
	if err != nil {
		return nil, err
	}
	return frByID, nil
}

func (h *GRPCHandler) GetFranchiseByBusinessName(ctx context.Context, req *pb.GetFranchiseByName) (*pb.GetFranchiseByIDResponse, error) {
	var method = constants.Methods.GetFranchiseByBusinessName

	if err := validations.ValidateNotEmpty(req.GetName()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	result, err := h.accountService.GetFranchiseByBusinessName(ctx, req.GetName())
	if err != nil {
		eR := status.Errorf(codes.Internal, "failed to fetch franchise: %v", err)
		logger.Error("failed to fetch franchise", eR, logger.BaseLogContext("layer", layer, "method", method))
		return nil, eR
	}

	frByName, err := mapper.GetFranchiseByID_ModelToProto(result)
	if err != nil {
		return nil, err
	}
	return frByName, nil
}

func (h *GRPCHandler) GetFranchiseOwnerByID(ctx context.Context, req *pb.GetFranchiseOwnerRequest) (*pb.GetFranchiseOwnerResponse, error) {
	var method = constants.Methods.GetFranchiseOwnerByID

	if err := validations.ValidateUUID(req.GetOwnerId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	result, err := h.accountService.GetFranchiseOwnerByID(ctx, req.GetOwnerId())
	if err != nil {
		eR := status.Errorf(codes.Internal, "failed to fetch franchise owner: %v", err)
		logger.Error("failed to fetch franchise owner", eR, logger.BaseLogContext("layer", layer, "method", method))
		return nil, eR
	}

	return mapper.FranchiseOwner_ModelToProto(result), nil
}

func (h *GRPCHandler) GetOwnerByAadharID(ctx context.Context, req *pb.AadharNoRequest) (*pb.BoolResponse, error) {
	var method = constants.Methods.CheckIfOwnerExistsByAadharID

	if err := validations.ValidateAadhaarNumber(req.GetAadharNo()); err != nil {
		logger.Error(erMsg+"Aadhar no", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	exists, err := h.accountService.CheckIfOwnerExistsByAadharID(ctx, req.GetAadharNo())
	if err != nil {
		eR := status.Errorf(codes.Internal, "failed to check owner existence: %v", err)
		logger.Error("failed to check if owner exists", eR, logger.BaseLogContext("layer", layer, "method", method))
		return nil, eR
	}

	return &pb.BoolResponse{Exist: exists}, nil
}

func (h *GRPCHandler) CreateFranchiseAccount(ctx context.Context, req *pb.AddFranchiseAccountRequest) (*pb.AddResponse, error) {
	var method = constants.Methods.CreateFranchiseAccount

	// Map incoming proto request to model
	account := mapper.AddFranchiseAccount_ProtoToModel(req)

	// Optional: Add validation here (if you have a validations.ValidateFranchiseAccount function)
	if err := validations.ValidateFranchiseAccounts(account); err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.ValidationFailed, err)
		logger.Error(constants.ValidationFailed, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}

	// Call service layer
	resp, err := h.accountService.CreateFranchiseAccount(ctx, account)
	if err != nil {
		eR := status.Errorf(codes.Internal, constants.SomethinWentWrongOnNew, err)
		logger.Error("failed to create franchise account", eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}

	// Return response
	return mapper.Add_ModelToProto(resp), nil
}

func (h *GRPCHandler) UpdateFranchiseAccountByID(ctx context.Context, req *pb.UpdateFranchiseAccountRequest) (*pb.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseAccount
	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetRoleId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	// Map proto to model
	account := mapper.UpdateFranchiseAccount_ProtoToModel(req)

	// Validate model
	if err := validations.ValidateFranchiseAccounts(account); err != nil {
		eR := status.Errorf(codes.InvalidArgument, constants.ValidationFailed, err)
		logger.Error(constants.ValidationFailed, eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}

	// Call service layer
	resp, err := h.accountService.UpdateFranchiseAccount(ctx, req.GetId(), account)
	if err != nil {
		eR := status.Errorf(codes.Internal, constants.SomethinWentWrongOnUpdate, err)
		logger.Error("failed to update franchise account", eR, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, eR
	}

	// Return proto response
	return &pb.UpdateResponse{
		Id:        resp.ID,
		UpdatedAt: resp.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h *GRPCHandler) GetFranchiseAccountByID(ctx context.Context, req *pb.GetFranchiseAccountByIDRequest) (*pb.GetFranchiseAccountByIDResponse, error) {
	var method = constants.Methods.GetFranchiseAccountByID

	if err := validations.ValidateUUID(req.GetAccountId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	resp, err := h.accountService.GetFranchiseAccountByID(ctx, req.GetAccountId())
	if err != nil {
		logger.Error("failed to get franchise account by id", err, logger.BaseLogContext("method", method))
		return nil, status.Errorf(codes.NotFound, "franchise account not found: %v", err)
	}

	if resp == nil {
		return nil, status.Errorf(codes.NotFound, "franchise account not found")
	}

	// Map model to proto response
	return &pb.GetFranchiseAccountByIDResponse{
		Id:          resp.ID,
		FranchiseId: resp.FranchiseID,
		Accounts: &pb.AccountInput{
			EmpId:       resp.EmployeeID,
			LoginId:     resp.LoginID,
			AccountType: resp.AccountType,
			Name:        resp.Name,
			MobileNo:    resp.MobileNo,
			EmailId:     resp.Email,
			Status:      resp.Status,
		},
		CreatedAt: formatTime(resp.CreatedAt),
		UpdatedAt: formatTime(resp.UpdatedAt),
	}, nil
}

func (h *GRPCHandler) GetAllFranchiseAccounts(ctx context.Context, req *pb.GetFranchiseAccountsRequest) (*pb.GetFranchiseAccountsResponse, error) {
	var method = constants.Methods.GetAllFranchiseAccounts

	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	franDetail := mapper.GetFranchisesRequest_ProtoToModel(req)

	accounts, err := h.accountService.GetAllFranchiseAccounts(ctx, franDetail)
	if err != nil {
		logger.Error("failed to get all franchise accounts", err, logger.BaseLogContext("method", method))
		return nil, status.Errorf(codes.Internal, "failed to fetch franchise accounts: %v", err)
	}

	pbAccounts := make([]*pb.GetFranchiseAccountByIDResponse, 0, len(accounts))
	for _, acc := range accounts {
		pbAcc := &pb.GetFranchiseAccountByIDResponse{
			Id:          acc.ID,
			FranchiseId: acc.FranchiseID,
			Accounts: &pb.AccountInput{
				EmpId:       acc.EmployeeID,
				LoginId:     acc.LoginID,
				AccountType: acc.AccountType,
				Name:        acc.Name,
				MobileNo:    acc.MobileNo,
				EmailId:     acc.Email,
				Status:      acc.Status,
			},
			CreatedAt: formatTime(acc.CreatedAt),
			UpdatedAt: formatTime(acc.UpdatedAt),
		}
		pbAccounts = append(pbAccounts, pbAcc)
	}

	return &pb.GetFranchiseAccountsResponse{
		Accounts: pbAccounts,
		Pagination: &pb.PaginationResponse{
			Page:  franDetail.GetPagination.Pagination.Page,
			Limit: franDetail.GetPagination.Pagination.Limit,
		},
	}, nil
}

func (h *GRPCHandler) AddFranchiseDocument(ctx context.Context, req *pb.AddFranchiseDocumentRequest) (*pb.AddResponse, error) {
	var method = constants.Methods.AddFranchiseDocument

	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	doc := mapper.AddFranchiseDocument_ProtoToModel(req)

	err := validations.ValidateFranchiseDocument(doc)
	if err != nil {
		logger.Error("validation failed", err, nil)
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	resp, err := h.accountService.AddFranchiseDocument(ctx, doc)
	if err != nil {
		logger.Error("failed to add franchise document", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnNew, err)
	}

	return &pb.AddResponse{
		Id:        resp.ID,
		CreatedAt: resp.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *GRPCHandler) UpdateFranchiseDocumentByID(ctx context.Context, req *pb.UpdateFranchiseDocumentRequest) (*pb.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseDocument

	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	doc := mapper.UpdateFranchiseDocument_ProtoToModel(req)

	if doc == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data")
	}

	if err := validations.ValidateFranchiseDocument(doc); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	resp, err := h.accountService.UpdateFranchiseDocument(ctx, req.Id, doc)
	if err != nil {
		logger.Error("failed to update franchise document", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnUpdate, err)
	}

	return &pb.UpdateResponse{
		Id:        resp.ID,
		UpdatedAt: resp.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h *GRPCHandler) CreateFranchiseAddress(ctx context.Context, req *pb.AddFranchiseAddressRequest) (*pb.AddResponse, error) {
	var method = constants.Methods.CreateFranchise
	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	addr := mapper.AddFranchiseAddress_ProtoToModel(req)
	if addr == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid franchise address request")
	}

	res, err := h.accountService.AddFranchiseAddress(ctx, addr)
	if err != nil {
		logger.Error("failed to add franchise address", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnNew, err)
	}

	return &pb.AddResponse{
		Id:        res.ID,
		CreatedAt: res.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *GRPCHandler) UpdateFranchiseAddress(ctx context.Context, req *pb.UpdateFranchiseAddressRequest) (*pb.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseAddress

	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing address ID")
	}
	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	addr := mapper.UpdateFranchiseAddress_ProtoToModel(req)
	if addr == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address data")
	}

	res, err := h.accountService.UpdateFranchiseAddress(ctx, req.Id, addr)
	if err != nil {
		logger.Error("failed to update franchise address", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnNew, err)
	}

	return &pb.UpdateResponse{
		Id:        res.ID,
		UpdatedAt: res.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h *GRPCHandler) GetFranchiseAddressByID(ctx context.Context, req *pb.GetFranchiseAddressRequest) (*pb.GetFranchiseAddressResponse, error) {
	var method = constants.Methods.GetFranchiseAddressByID

	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "address ID is required")
	}
	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	result, err := h.accountService.GetFranchiseAddressByID(ctx, req.Id)
	if err != nil {
		logger.Error("failed to fetch franchise address by ID", err, logger.BaseLogContext(
			"method", method,
			"layer", "handler",
			"id", req.Id,
		))
		return nil, status.Errorf(codes.Internal, "something went wrong: %s", err)
	}

	return mapper.FranchiseAddressModelToProto(result, req.FranchiseId), nil
}

func (h *GRPCHandler) CreateFranchiseRole(ctx context.Context, req *pb.AddFranchiseRoleRequest) (*pb.AddResponse, error) {
	var method = constants.Methods.CreateFranchiseRole

	if req == nil || req.FranchiseId == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "franchise_id and name are required")
	}
	if err := validations.ValidateUUID(req.GetFranchiseId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	role := mapper.FranchiseRoleProtoToModel(req)

	res, err := h.accountService.AddFranchiseRole(ctx, role)
	if err != nil {
		logger.Error("failed to create franchise role", err, logger.BaseLogContext(
			"method", method,
			"layer", "handler",
			"franchise_id", req.FranchiseId,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnNew, err)
	}

	return &pb.AddResponse{
		Id:        res.ID,
		CreatedAt: formatTime(&res.CreatedAt),
	}, nil
}

func (h *GRPCHandler) UpdateFranchiseRole(ctx context.Context, req *pb.UpdateFranchiseRoleRequest) (*pb.UpdateResponse, error) {
	var method = constants.Methods.UpdateFranchiseRole

	if req == nil || req.Id == "" || req.FRole == nil {
		return nil, status.Error(codes.InvalidArgument, "id and f_role are required")
	}
	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	role := mapper.MapFranchiseRoleInputToModel(req.FRole)

	updatedRes, err := h.accountService.UpdateFranchiseRole(ctx, req.Id, role)
	if err != nil {
		logger.Error("failed to update franchise role", err, logger.BaseLogContext(
			"method", method,
			"layer", "handler",
			"id", req.Id,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnUpdate, err)
	}

	return &pb.UpdateResponse{
		Id:        updatedRes.ID,
		UpdatedAt: formatTime(&updatedRes.UpdatedAt),
	}, nil
}

func (h *GRPCHandler) GetAllFranchiseRoles(ctx context.Context, req *pb.GetByIDRequest) (*pb.FranchiseRoleResponse, error) {
	var method = constants.Methods.GetAllFranchiseRoles

	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "franchise_id is required")
	}
	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}

	roles, err := h.accountService.GetAllFranchiseRoles(ctx, req.Id)
	if err != nil {
		logger.Error("failed to get franchise roles", err, logger.BaseLogContext(
			"method", method,
			"layer", "handler",
			"franchise_id", req.Id,
		))
		return nil, status.Errorf(codes.Internal, constants.SomethinWentWrongOnNew, err)
	}

	var franchiseRoles []*pb.AddFranchiseRoleRequest
	for _, r := range roles {
		franchiseRoles = append(franchiseRoles, &pb.AddFranchiseRoleRequest{
			FranchiseId: req.Id,
			Name:        r.Name,
			Description: r.Description,
			IsDefault:   r.IsDefault,
		})
	}

	return &pb.FranchiseRoleResponse{
		Id:    req.Id,
		FRole: franchiseRoles[0], // optional: for backward compatibility
		// consider defining a repeated field in proto if multiple roles are returned
		CreatedAt: formatTime(roles[0].CreatedAt),
		UpdatedAt: formatTime(roles[0].UpdatedAt),
	}, nil
}

func (h *GRPCHandler) AddPermissionsToRole(ctx context.Context, req *pb.AddRolePermission) (*pb.AddRolePermission, error) {
	var method = constants.Methods.AddPermissionsToRole
	if err := validations.ValidateUUID(req.GetRoleId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetPermissionId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	rolePerm := mapper.MapAddRolePermissionRequestToModel(req)
	created, err := h.accountService.AddPermissionsToRole(ctx, rolePerm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add permission to role: %v", err)
	}

	return &pb.AddRolePermission{
		RoleId:       created.RoleID,
		PermissionId: created.PermissionID,
	}, nil
}

func (h *GRPCHandler) UpdatePermissionsToRole(ctx context.Context, req *pb.AddRolePermission) (*pb.AddRolePermission, error) {
	var method = constants.Methods.UpdatePermissionsToRole
	if err := validations.ValidateUUID(req.GetRoleId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	if err := validations.ValidateUUID(req.GetPermissionId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	rolePerm := mapper.MapAddRolePermissionRequestToModel(req)
	updated, err := h.accountService.UpdatePermissionsToRole(ctx, rolePerm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update permission to role: %v", err)
	}

	return &pb.AddRolePermission{
		RoleId:       updated.RoleID,
		PermissionId: updated.PermissionID,
	}, nil
}

func (h *GRPCHandler) GetAllPermissionToRole(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetAllRolePermissionDetails, error) {
	var method = constants.Methods.UpdatePermissionsToRole
	if err := validations.ValidateUUID(req.GetId()); err != nil {
		logger.Error(erMsg+"UUID", err, logger.BaseLogContext(
			"layer", layer,
			"method", method,
		))
		return nil, err
	}
	perms, err := h.accountService.GetAllPermissionsToRole(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get role permissions: %v", err)
	}

	if len(perms) == 0 {
		return nil, status.Error(codes.NotFound, "no permissions found for the role")
	}

	// Map to repeated RolePermissionDetails
	rolePermissions := mapper.MapRolePermissionToProto(perms)

	return &pb.GetAllRolePermissionDetails{
		RoleP: rolePermissions,
	}, nil
}
