package model

import (
	"errors"

	"github.com/ashish19912009/zrms/services/authZ/pb"
)

type CheckAccess struct {
	AccountID   string            `json:"account_id"`
	FranchiseID string            `json:"franchise_id"`
	Resource    string            `json:"resource"`
	Action      string            `json:"action"`
	Context     map[string]string `json:"context,omitempty"`
}

type ResourceAction struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type BatchCheckAccess struct {
	AccountID   string            `json:"account_id"`
	FranchiseID string            `json:"franchise_id"`
	Resources   []ResourceAction  `json:"resource"`
	Context     map[string]string `json:"context,omitempty"`
}

func CheckAccessFromPbToModel(cA *pb.CheckAccessRequest) (*CheckAccess, error) {
	if cA.AccountId == "" || cA.FranchiseId == "" {
		return nil, errors.New("gRPC error: account id or franchise id can't be empty")
	}
	return &CheckAccess{
		AccountID:   cA.AccountId,
		FranchiseID: cA.FranchiseId,
		Resource:    cA.Resource,
		Action:      cA.Action,
		Context:     cA.Context,
	}, nil
}

func BatchCheckAccessFromPbToModel(cA *pb.BatchCheckAccessRequest) (*BatchCheckAccess, error) {
	if cA.AccountId == "" || cA.FranchiseId == "" {
		return nil, errors.New("gRPC error: account id or franchise id can't be empty")
	}
	resources := []ResourceAction{}
	for _, r := range cA.Resources {
		resources = append(resources, ResourceAction{
			Resource: r.Resource,
			Action:   r.Action,
		})
	}
	return &BatchCheckAccess{
		AccountID:   cA.AccountId,
		FranchiseID: cA.FranchiseId,
		Resources:   resources,
		Context:     cA.Context,
	}, nil
}
