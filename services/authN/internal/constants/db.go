package constants

var DB = struct {
	Schema_Global            string
	Schema_Outlet            string
	Table_Direct_Permissions string
	Table_Permissions        string
	Table_Role               string

	Table_Franchise_Accounts string

	Table_Roles            string
	Table_Document_Types   string
	Table_Role_Permissions string
}{
	Schema_Global:            "global",
	Schema_Outlet:            "outlet",
	Table_Direct_Permissions: "direct_permissions",
	Table_Permissions:        "permissions",

	Table_Franchise_Accounts: "team_accounts",

	Table_Roles:            "roles",
	Table_Document_Types:   "document_types",
	Table_Role_Permissions: "role_permissions",
}
