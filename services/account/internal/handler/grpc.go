package handler

import (
	"context"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/service"
	"github.com/ashish19912009/zrms/services/account/pb"
)

type GRPCHandler struct {
	pb.UnimplementedAccountServiceServer
	accountService service.AccountService
}

func NewGRPCHandler(accountService service.AccountService) *GRPCHandler {
	return &GRPCHandler{accountService: accountService}
}

func (h *GRPCHandler) CreateAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {

	acc := &model.Account{
		ID:         req.Id,
		MobileNo:   req.MobileNo,
		Name:       req.Name,
		Role:       req.Role,
		Status:     req.Status,
		EmployeeID: req.EmpId,
	}
	created, err := h.accountService.CreateAccount(ctx, acc)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
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
	return &pb.UpdateAccountResponse{
		Account: &pb.Account{
			Id:        updated.ID,
			MobileNo:  updated.MobileNo,
			Name:      updated.Name,
			Role:      updated.Role,
			Status:    updated.Status,
			EmpId:     updated.EmployeeID,
			CreatedAt: updated.CreatedAt.Format(time.RFC3339),
			UpdatedAt: updated.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *GRPCHandler) GetAccountByID(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := h.accountService.GetAccountByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:        account.ID,
			MobileNo:  account.MobileNo,
			Name:      account.Name,
			Role:      account.Role,
			Status:    account.Status,
			EmpId:     account.EmployeeID,
			CreatedAt: account.CreatedAt.Format(time.RFC3339),
			UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *GRPCHandler) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := h.accountService.ListAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var res []*pb.Account
	for _, acc := range accounts {
		res = append(res, &pb.Account{
			Id:        acc.ID,
			MobileNo:  acc.MobileNo,
			Name:      acc.Name,
			Role:      acc.Role,
			Status:    acc.Status,
			EmpId:     acc.EmployeeID,
			CreatedAt: acc.CreatedAt.Format(time.RFC3339),
		})
	}
	return &pb.GetAccountsResponse{Accounts: res}, nil
}
