package constants

var DB = struct {
	Schema_Global             string
	Schema_Outlet             string
	Table_Franchise           string
	Table_Owner               string
	Table_Franchise_documents string

	Table_Franchise_Accounts string

	Table_Roles          string
	Table_Document_Types string
}{
	Schema_Global:             "global",
	Schema_Outlet:             "outlet",
	Table_Franchise:           "franchises",
	Table_Owner:               "owner",
	Table_Franchise_documents: "franchise_documents",

	Table_Franchise_Accounts: "team_accounts",

	Table_Roles:          "roles",
	Table_Document_Types: "document_types",
}
