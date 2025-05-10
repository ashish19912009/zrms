package model

import (
	"time"
)

type ResourceAction struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	IsDefault   bool         `json:"is_default"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// RoleResponse represents the API response for role data
type RoleResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	PermissionCodes []string  `json:"permissions,omitempty"`
}

// ToResponse converts Role to RoleResponse
func (r *Role) ToResponse() RoleResponse {
	var permCodes []string
	for _, perm := range r.Permissions {
		permCodes = append(permCodes, perm.Code)
	}

	return RoleResponse{
		ID:              r.ID,
		Name:            r.Name,
		Description:     r.Description,
		IsDefault:       r.IsDefault,
		CreatedAt:       r.CreatedAt,
		PermissionCodes: permCodes,
	}
}

// RoleCreateInput represents data needed to create a new role
// type RoleCreateInput struct {
// 	Name        string `json:"name" validate:"required,alphanum"`
// 	Description string `json:"description" validate:"omitempty,max=500"`
// 	IsDefault   bool   `json:"isDefault"`
// }

// // RoleUpdateInput represents data needed to update a role
// type RoleUpdateInput struct {
// 	Name        *string `json:"name" validate:"omitempty,alphanum"`
// 	Description *string `json:"description" validate:"omitempty,max=500"`
// 	IsDefault   *bool   `json:"isDefault"`
// }

// // RoleResponse represents how we send role data in API responses
// type RoleResponse struct {
// 	ID          uuid.UUID `json:"id"`
// 	Name        string    `json:"name"`
// 	Description string    `json:"description"`
// 	IsDefault   bool      `json:"isDefault"`
// 	CreatedAt   time.Time `json:"createdAt"`
// 	PermissionCount int    `json:"permissionCount"`
// }
