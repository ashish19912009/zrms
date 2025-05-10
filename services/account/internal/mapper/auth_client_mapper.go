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

func BatchCheckAccessFromPbToModel(cA *pb.BatchCheckAccessRequest) (*model.BatchCheckAccess, error) {
	if cA.AccountId == "" || cA.FranchiseId == "" {
		return nil, errors.New("gRPC error: account id or franchise id can't be empty")
	}
	resources := []model.ResourceAction{}
	for _, r := range cA.Resources {
		resources = append(resources, model.ResourceAction{
			Resource: r.Resource,
			Action:   r.Action,
		})
	}
	return &model.BatchCheckAccess{
		AccountID:   cA.AccountId,
		FranchiseID: cA.FranchiseId,
		Resources:   resources,
		Context:     cA.Context,
	}, nil
}

func BatchCheckAccessFromModelToPb(responses []*model.CheckBatchAccessResponse) *pb.BatchCheckAccessResponse {
	pbResponse := &pb.BatchCheckAccessResponse{
		Results: make([]*pb.ResourceActionResult, len(responses)),
	}

	for i, res := range responses {
		// Create a new ResourceActionResult for each response
		result := &pb.ResourceActionResult{}

		if res != nil {
			// Create and populate the Decision if response exists
			result.Decision = &pb.Decision{
				Allowed:       res.Allowed,
				Reason:        res.Reason,
				IssuedAt:      res.IssuedAt,
				ExpiresAt:     res.ExpiresAt,
				PolicyVersion: res.PolicyVersion,
			}
			result.ResAct = &pb.ResourceAction{
				Resource: res.Resource,
				Action:   res.Action,
			}
		} else {
			// Create a default denied decision for nil responses
			result.Decision = &pb.Decision{
				Allowed:       false,
				Reason:        "no decision available",
				IssuedAt:      0,
				ExpiresAt:     0,
				PolicyVersion: "",
			}
			result.ResAct = &pb.ResourceAction{
				Resource: result.ResAct.Resource,
				Action:   result.ResAct.Action,
			}
		}
		pbResponse.Results[i] = result
	}

	return pbResponse
}
