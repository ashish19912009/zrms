package handler

import (
	"context"
	"log"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/helper"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	pb.UnimplementedAccountServiceServer
	accountService service.AccountService
}

func NewGRPCHandler(accountService service.AccountService) *GRPCHandler {
	return &GRPCHandler{accountService: accountService}
}

func formatDatesToString(acc *model.Account) (createdAt, updatedAt string) {
	if acc.CreatedAt != nil {
		createdAt = helper.FormatTimePtr(acc.CreatedAt)
	}
	if acc.UpdatedAt != nil {
		updatedAt = helper.FormatTimePtr(acc.UpdatedAt)
	}
	return
}

func (h *GRPCHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	if h.accountService == nil {
		log.Fatal("❌ accountService is nil")
	}
	if req.MobileNo == "" || req.EmpId == "" || req.Name == "" || req.Role == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "required fields missing")
	}
	acc := &model.Account{
		ID:         req.Id,
		MobileNo:   req.MobileNo,
		Name:       req.Name,
		Role:       req.Role,
		Status:     req.Status,
		EmployeeID: req.EmpId,
		CreatedAt:  helper.NowPtr(),
	}
	created, err := h.accountService.CreateAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	return &pb.CreateAccountResponse{
		Account: &pb.Account{
			Id:        created.ID,
			MobileNo:  created.MobileNo,
			Name:      created.Name,
			Role:      created.Role,
			Status:    created.Status,
			EmpId:     created.EmployeeID,
			CreatedAt: created.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *GRPCHandler) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "account ID is required")
	}
	if req.MobileNo == "" && req.Name == "" && req.Role == "" && req.Status == "" && req.EmpId == "" {
		return nil, status.Error(codes.InvalidArgument, "no fields to update")
	}
	acc := &model.Account{
		ID:         req.Id,
		MobileNo:   req.MobileNo,
		Name:       req.Name,
		Role:       req.Role,
		Status:     req.Status,
		EmployeeID: req.EmpId,
	}
	updated, err := h.accountService.UpdateAccount(ctx, acc)
	if err != nil {
		return nil, err
	}

	var created_at, updated_at = formatDatesToString(updated)

	return &pb.UpdateAccountResponse{
		Account: &pb.Account{
			Id:        updated.ID,
			MobileNo:  updated.MobileNo,
			Name:      updated.Name,
			Role:      updated.Role,
			Status:    updated.Status,
			EmpId:     updated.EmployeeID,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		},
	}, nil
}

func (h *GRPCHandler) GetAccountByID(ctx context.Context, req *pb.GetAccountByIDRequest) (*pb.GetAccountByIDResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "account ID is required")
	}

	account, err := h.accountService.GetAccountByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var created_at, updated_at = formatDatesToString(account)

	return &pb.GetAccountByIDResponse{
		Account: &pb.Account{
			Id:        account.ID,
			MobileNo:  account.MobileNo,
			Name:      account.Name,
			Role:      account.Role,
			Status:    account.Status,
			EmpId:     account.EmployeeID,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		},
	}, nil
}

func (h *GRPCHandler) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	if req.Take == 0 {
		return nil, status.Error(codes.InvalidArgument, "take must be greater than zero")
	}

	accounts, err := h.accountService.ListAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var res []*pb.Account
	var created_at, updated_at string
	for _, acc := range accounts {
		created_at, updated_at = formatDatesToString(acc)
		res = append(res, &pb.Account{
			Id:        acc.ID,
			MobileNo:  acc.MobileNo,
			Name:      acc.Name,
			Role:      acc.Role,
			Status:    acc.Status,
			EmpId:     acc.EmployeeID,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		})
	}
	return &pb.GetAccountsResponse{Accounts: res}, nil
}
