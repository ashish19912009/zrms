package model

// "id",
// "franchise_id",
// "employee_id",
// "login_id",
// "password_hash",
// "account_type",
// "name",
// "mobile_no",
// "email",
// "role_id",
// "status",

type User struct {
	AccountID   string `json:"account_id"`
	FranchiseID string `json:"franchise_id"`
	EmployeeID  string `json:"employee_id"`
	AccountType string `json:"account_type"`
	RoleID      string `json:"role_id"`
	Name        string `json:"name"`
	MobileNo    string `json:"mobile_no"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Status      string `json:"status"`
}

// type UserAdminPermission struct {
// 	Create bool `json:"create"`
// 	Read   bool `json:"read"`
// 	Update bool `json:"update"`
// 	Delete bool `json:"delete"`
// }

// type Permission struct {
// 	UserAccount UserAccountPermission `json:"user_account"`
// }

// type PermissionsArray []Permission

// func (p *PermissionsArray) Scan(src interface{}) error {
// 	rawArray, ok := src.([]byte)
// 	if !ok {
// 		return fmt.Errorf("permissions: failed to type assert []byte")
// 	}

// 	var rawStrings []string
// 	if err := json.Unmarshal(rawArray, &rawStrings); err != nil {
// 		return fmt.Errorf("permissions: error unmarshaling raw json array of strings: %w", err)
// 	}

// 	for _, raw := range rawStrings {
// 		var perm Permission
// 		if err := json.Unmarshal([]byte(raw), &perm); err != nil {
// 			return fmt.Errorf("permissions: failed to unmarshal element: %w", err)
// 		}
// 		*p = append(*p, perm)
// 	}
// 	return nil
// }

// type JSONStringArray []string

// func (j *JSONStringArray) Scan(value interface{}) error {
// 	bytes, ok := value.([]byte)
// 	if !ok {
// 		return fmt.Errorf("type assertion to []byte failed")
// 	}
// 	return json.Unmarshal(bytes, j)
// }

// type User struct {
// 	AccountID   string          `json:"account_id"`
// 	EmployeeID  string          `json:"employee_id"`
// 	AccountType string          `json:"account_type"`
// 	Name        string          `json:"name"`
// 	MobileNo    string          `json:"mobile_no"`
// 	Password    string          `json:"password"`
// 	Role        string          `json:"role"`
// 	Permissions JSONStringArray `json:"permissions"`
// 	Status      string          `json:"status"`
// }
