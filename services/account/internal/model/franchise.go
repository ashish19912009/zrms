package model

import "time"

type Franchise struct {
	BusinessName  string                 `json:"business_name"`
	LogoURL       string                 `json:"logo_url"`
	SubDomain     string                 `json:"sub_domain"`
	ThemeSettings map[string]interface{} `json:"theme_settings"`
	Status        string                 `json:"status"`
}

type FranchiseResponse struct {
	ID            string                 `json:"id"`
	BusinessName  string                 `json:"business_name"`
	LogoURL       string                 `json:"logo_url"`
	SubDomain     string                 `json:"sub_domain"`
	ThemeSettings map[string]interface{} `json:"theme_settings"`
	Status        string                 `json:"status"`
	CreatedAt     *time.Time             `json:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at"`
	DeletedAt     *time.Time             `json:"deleted_at,omitempty"` // Nullable field for soft deletion
}

type FranchiseDocument struct {
	DocumentName        string `json:"doc_name"`
	DocumentDescription string `json:"doc_desc"`
	IsMandate           string `json:"is_mandate"`
	DocumentURL         string `json:"document_url"`
	UploadedBy          string `json:"uploaded_by"`
}

type FranchiseDocumentResponse struct {
	ID                  string     `json:"id"`
	DocumentName        string     `json:"doc_name"`
	DocumentDescription string     `json:"doc_desc"`
	IsMandate           string     `json:"is_mandate"`
	DocumentURL         string     `json:"document_url"`
	UploadedBy          string     `json:"uploaded_by"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"` // Nullable field for soft deletion
}

type FranchiseAddress struct {
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Pincode     string `json:"pincode"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
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
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"` // Nullable field for soft deletion
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
