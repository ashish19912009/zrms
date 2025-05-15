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
