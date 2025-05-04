package model

import "github.com/ashish19912009/zrms/services/authZ/pb"

type CheckAccessResponse struct {
	Allowed       bool   `json:"allowed"`
	Reason        string `json:"reason"`
	IssuedAt      int64  `json:"issued_at"`
	ExpiresAt     int64  `json:"expires_at"`
	PolicyVersion string `json:"policy_version"`
}

type CheckBatchAccessResponse struct {
	Resource      string `json:"resource"`
	Action        string `json:"action"`
	Allowed       bool   `json:"allowed"`
	Reason        string `json:"reason"`
	IssuedAt      int64  `json:"issued_at"`
	ExpiresAt     int64  `json:"expires_at"`
	PolicyVersion string `json:"policy_version"`
}

func CheckAccessFromModelToPb(cR *CheckAccessResponse) *pb.CheckAccessResponse {
	return &pb.CheckAccessResponse{
		Decision: &pb.Decision{
			Allowed:       cR.Allowed,
			Reason:        cR.Reason,
			IssuedAt:      cR.IssuedAt,
			ExpiresAt:     cR.ExpiresAt,
			PolicyVersion: cR.PolicyVersion,
		},
	}
}

func BatchCheckAccessFromModelToPb(responses []*CheckBatchAccessResponse) *pb.BatchCheckAccessResponse {
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
			// result.ResAct = &pb.ResourceAction{
			// 	Resource: res.Resource,
			// 	Action:   res.Action,
			// }
		} else {
			// Create a default denied decision for nil responses
			result.Decision = &pb.Decision{
				Allowed:       false,
				Reason:        "no decision available",
				IssuedAt:      0,
				ExpiresAt:     0,
				PolicyVersion: "",
			}
			// result.ResAct = &pb.ResourceAction{
			// 	Resource: result.ResAct.Resource,
			// 	Action:   result.ResAct.Action,
			// }
		}

		pbResponse.Results[i] = result
	}

	return pbResponse
}
