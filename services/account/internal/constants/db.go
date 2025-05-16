package constants

var DB = struct {
	Schema_Global             string
	Schema_Outlet             string
	Table_Franchise           string
	Table_Owner               string
	Table_Franchise_documents string
	Table_Franchise_addresses string

	Table_Franchise_Accounts string

	Table_Roles            string
	Table_Document_Types   string
	Table_Role_Permissions string
}{
	Schema_Global:             "global",
	Schema_Outlet:             "outlet",
	Table_Franchise:           "franchises",
	Table_Owner:               "owner",
	Table_Franchise_documents: "franchise_documents",
	Table_Franchise_addresses: "franchise_addresses",

	Table_Franchise_Accounts: "team_accounts",

	Table_Roles:            "roles",
	Table_Document_Types:   "document_types",
	Table_Role_Permissions: "role_permissions",
}

var T_Fran = struct {
	UUID             string
	BusinessName     string
	LogoUrl          string
	Subdomain        string
	ThemeSettings    string
	Status           string
	CreatedAt        string
	UpdatedAt        string
	DeletedAt        string
	FranchiseOwnerID string
}{
	UUID:             "id",
	BusinessName:     "business_name",
	LogoUrl:          "logo_url",
	Subdomain:        "sub_domain",
	ThemeSettings:    "theme_settings",
	Status:           "status",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
	FranchiseOwnerID: "franchise_owner_id",
}

var T_Onr = struct {
	UUID       string
	Name       string
	Gender     string
	DOB        string
	MobileNo   string
	Email      string
	Address    string
	AadharNo   string
	IsVerified string
	Status     string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}{
	UUID:       "id",
	Name:       "name",
	Gender:     "gender",
	DOB:        "dob",
	MobileNo:   "mobile_no",
	Email:      "email",
	Address:    "address",
	AadharNo:   "aadhar_no",
	IsVerified: "is_verified",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

var Acc = struct {
	UUID        string
	FranchiseID string
	EmpID       string
	LoginID     string
	Password    string
	AccountType string
	Name        string
	MobileNo    string
	Email       string
	RoleID      string
	Status      string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}{
	UUID:        "id",
	FranchiseID: "franchise_id",
	EmpID:       "employee_id",
	LoginID:     "login_id",
	Password:    "password",
	AccountType: "account_type",
	Name:        "name",
	MobileNo:    "mobile_no",
	Email:       "email",
	RoleID:      "role_id",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

var F_doc = struct {
	UUID           string
	FranchiseID    string
	DocumentTypeID string
	DocumentURL    string
	UploadedBy     string
	Status         string
	Remark         string
	VerifiedAt     string
	CreatedAt      string
	UpdatedAt      string
	DeletedAt      string
}{
	UUID:           "id",
	FranchiseID:    "franchise_id",
	DocumentTypeID: "document_type_id",
	DocumentURL:    "document_url",
	UploadedBy:     "uploaded_by",
	Status:         "status",
	Remark:         "remark",
	VerifiedAt:     "verified_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

var F_addr = struct {
	UUID        string
	FranchiseID string
	AddressLine string
	City        string
	State       string
	Country     string
	Pincode     string
	Latitude    string
	Longitude   string
	IsVerified  string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}{
	UUID:        "id",
	FranchiseID: "franchise_id",
	AddressLine: "address_line",
	City:        "city",
	State:       "state",
	Country:     "country",
	Pincode:     "pincode",
	Latitude:    "latitude",
	Longitude:   "longitude",
	IsVerified:  "is_verified",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

var F_role = struct {
	UUID        string
	FranchiseID string
	Name        string
	Description string
	IsDefault   string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}{
	UUID:        "id",
	FranchiseID: "franchise_id",
	Name:        "name",
	Description: "description",
	IsDefault:   "is_default",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

var F_Role_Per = struct {
	RoleID       string
	PermissionID string
}{
	RoleID:       "role_id",
	PermissionID: "permission_id",
}
