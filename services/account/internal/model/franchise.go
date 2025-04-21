package model

import "time"

type Franchise struct {
	BusinessName       string                 `json:"business_name"`
	LogoURL            string                 `json:"logo_url"`
	SubDomain          string                 `json:"sub_domain"`
	ThemeSettings      map[string]interface{} `json:"theme_settings"`
	Status             string                 `json:"status"`
	Franchise_Owner_id string                 `json:"franchise_owner_id"`
}

type FranchiseResponse struct {
	ID                 string                 `json:"id"`
	BusinessName       string                 `json:"business_name"`
	LogoURL            string                 `json:"logo_url"`
	SubDomain          string                 `json:"sub_domain"`
	ThemeSettings      map[string]interface{} `json:"theme_settings"`
	Status             string                 `json:"status"`
	Franchise_Owner_id string                 `json:"franchise_owner_id"`
	CreatedAt          *time.Time             `json:"created_at"`
	UpdatedAt          *time.Time             `json:"updated_at"`
	DeletedAt          *time.Time             `json:"deleted_at,omitempty"` // Nullable field for soft deletion
}

type CommonReturn struct {
	ID string `json:"id"`
}

type FranchiseDocument struct {
	FranchiseID    string `json:"franchise_id"`
	DocumentTypeID string `json:"document_type_id"`
	DocumentURL    string `json:"document_url"`
	UploadedBy     string `json:"uploaded_by"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
	VerifiedAt     string `json:"verified_id"`
}

type FranchiseDocumentResponse struct {
	ID             string `json:"id"`
	FranchiseID    string `json:"franchise_id"`
	DocumentTypeID string `json:"document_type_id"`
	DocumentURL    string `json:"document_url"`
	UploadedBy     string `json:"uploaded_by"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
	VerifiedAt     string `json:"verified_id"`
}

type FranchiseDocumentResponseComplete struct {
	ID                  string     `json:"id"`
	DocumentName        string     `json:"doc_name"`
	DocumentDescription string     `json:"doc_desc"`
	IsMandate           string     `json:"is_mandate"`
	DocumentURL         string     `json:"document_url"`
	UploadedBy          string     `json:"uploaded_by"`
	Status              string     `json:"status"`
	Remark              string     `json:"remark"`
	VerifiedAt          string     `json:"verified_id"`
	UploadedAt          *time.Time `json:"uploaded_at"`
}

type FranchiseAddress struct {
	FranchiseID string `json:"franchise_id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Pincode     string `json:"pincode"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	IsVerified  string `json:"is_verified"`
}

type FranchiseAddressResponse struct {
	ID          string     `json:"id"`
	AddressLine string     `json:"address_line"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	Country     string     `json:"country"`
	Pincode     string     `json:"pincode"`
	Latitude    string     `json:"latitude"`
	Longitude   string     `json:"longitude"`
	IsVerified  string     `json:"is_verified"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type FranchiseOwner struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Dob        string `json:"dob"`
	MobileNo   string `json:"mobile_no"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	AadharNo   string `json:"aadhar_no"`
	IsVerified string `json:"is_verified"`
	Status     string `json:"status"`
}

type FranchiseOwnerResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Gender     string     `json:"gender"`
	Dob        string     `json:"dob"`
	MobileNo   string     `json:"mobile_no"`
	Email      string     `json:"email"`
	Address    string     `json:"address"`
	AadharNo   string     `json:"aadhar_no"`
	IsVerified string     `json:"is_verified"`
	Status     string     `json:"status"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type FranchiseRole struct {
	FranchiseID string `json:"franchise_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
}

type FranchiseRoleResponse struct {
	ID          string     `json:"id"`
	FranchiseID string     `json:"franchise_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsDefault   bool       `json:"is_default"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type RoleToPermissions struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

type RoleToPermissionsComplete struct {
	FranchiseID            string     `json:"franchise_id"`
	RoleName               string     `json:"name"`
	Role_Description       string     `json:"description"`
	IsDefault              string     `json:"is_default"`
	Permission_Key         string     `json:"p.key"`
	Permission_Description string     `json:"p.description"`
	CreatedAt              *time.Time `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
}

// func (fo *FranchiseOwner) ToResponse() *FranchiseOwnerResponse{
// 	return &FranchiseOwnerResponse{
// 		ID: fo.ID,
// 	}
// }
