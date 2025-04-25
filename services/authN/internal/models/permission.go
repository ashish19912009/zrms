package models

import "time"

type Permission struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionCreateRequest represents the input for creating a permission
type PermissionCreateRequest struct {
	Code        string `json:"code" validate:"required,alphanum"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// PermissionUpdateRequest represents the input for updating a permission
type PermissionUpdateRequest struct {
	Code        *string `json:"code,omitempty" validate:"omitempty,alphanum"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// PermissionResponse represents the API response for permission data
type PermissionResponse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse converts Permission to PermissionResponse
func (p *Permission) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Code:        p.Code,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}
