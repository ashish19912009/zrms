package mapper

import (
	"errors"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/authZ/pb"
)

func CheckAccessFromPbToModel(cA *pb.CheckAccessResponse) *model.CheckAccessResponse {
	return &model.CheckAccessResponse{
		Allowed:       cA.Decision.Allowed,
		Reason:        cA.Decision.Reason,
		IssuedAt:      cA.Decision.IssuedAt,
		ExpiresAt:     cA.Decision.ExpiresAt,
		PolicyVersion: cA.Decision.PolicyVersion,
	}
}

func CheckAccessFromModelToPb(accountID, franchiseID, resource, action string) *pb.CheckAccessRequest {
	return &pb.CheckAccessRequest{
		AccountId:   accountID,
		FranchiseId: franchiseID,
		Resource:    resource,
		Action:      action,
	}
}

func BatchCheckAccessFromPbToModel(cA *pb.BatchCheckAccessResponse) (*model.BatchCheckAccessResponse, error) {
	pbResponse := &model.BatchCheckAccessResponse{
		Results: make([]*model.ResourceActionResult, len(cA.Results)),
	}

	for i, res := range cA.Results {
		// Create a new ResourceActionResult for each response
		result := &model.ResourceActionResult{}
		result.ResAct = &model.ResourceAction{
			Resource: res.ResAct.Resource,
			Action:   res.ResAct.Action,
		}
		if res != nil {
			result.Decision = &model.CheckAccessResponse{
				Allowed:       res.Decision.Allowed,
				Reason:        res.Decision.Reason,
				IssuedAt:      res.Decision.IssuedAt,
				ExpiresAt:     res.Decision.ExpiresAt,
				PolicyVersion: res.Decision.PolicyVersion,
			}
		} else {
			// Create a default denied decision for nil responses
			result.Decision = &model.CheckAccessResponse{
				Allowed:       false,
				Reason:        "no decision available",
				IssuedAt:      0,
				ExpiresAt:     0,
				PolicyVersion: "",
			}
		}
		pbResponse.Results[i] = result
	}
	return pbResponse, nil
}

func BatchCheckAccessFromModelToPb(req *model.BatchCheckAccess) (*pb.BatchCheckAccessRequest, error) {
	if req.AccountID == "" || req.FranchiseID == "" {
		return nil, errors.New("gRPC error: account id or franchise id can't be empty")
	}
	resources := []*pb.ResourceAction{}
	for _, r := range req.Resources {
		resources = append(resources, &pb.ResourceAction{
			Resource: r.Resource,
			Action:   r.Action,
		})
	}
	return &pb.BatchCheckAccessRequest{
		AccountId:   req.AccountID,
		FranchiseId: req.FranchiseID,
		Resources:   resources,
		Context:     req.Context,
	}, nil
}
