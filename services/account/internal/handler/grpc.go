package handler

import (
	"context"

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
}

func NewGRPCHandler(accountService service.AccountService, adminService service.AdminService) *GRPCHandler {
	return &GRPCHandler{
		accountService: accountService,
		adminService:   adminService,
	}
}

// UpdateFranchiseStatus(ctx context.Context, id string, status string) (bool, error)
// DeleteFranchise(ctx context.Context, id string) (bool, error)
// GetAllFranchises(ctx context.Context, page int32, limit int32) ([]model.FranchiseResponse, error)

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
	updateOwner, err := h.adminService.UpdateNewOwner(ctx, id, owner)
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
	// Convert proto to internal model
	franchise, err := mapper.AddFranchise_ProtoToModel(req.GetFranchiseDetails())
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
	// Convert proto to internal model
	var method = constants.Methods.UpdateFranchise
	// Convert proto to internal model
	franchise, err := mapper.AddFranchise_ProtoToModel(req.GetFranchiseDetails())
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
	updatedFranchise, err := h.adminService.UpdateFranchise(ctx, req.Id, franchise)
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
	// Convert proto to internal model
	var method = constants.Methods.DeleteFranchise
	// Convert proto to internal model
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
	// Convert proto to internal model
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
