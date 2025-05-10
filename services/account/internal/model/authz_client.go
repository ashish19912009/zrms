package model

type CheckAccess struct {
	AccountID   string            `json:"account_id"`
	FranchiseID string            `json:"franchise_id"`
	Resource    string            `json:"resource"`
	Action      string            `json:"action"`
	Context     map[string]string `json:"context,omitempty"`
}

type CheckAccessResponse struct {
	Allowed       bool   `json:"allowed"`
	Reason        string `json:"reason"`
	IssuedAt      int64  `json:"issued_at"`
	ExpiresAt     int64  `json:"expires_at"`
	PolicyVersion string `json:"policy_version"`
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

type CheckBatchAccessResponse struct {
	Resource      string `json:"resource"`
	Action        string `json:"action"`
	Allowed       bool   `json:"allowed"`
	Reason        string `json:"reason"`
	IssuedAt      int64  `json:"issued_at"`
	ExpiresAt     int64  `json:"expires_at"`
	PolicyVersion string `json:"policy_version"`
}
