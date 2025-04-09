package models

import (
	"encoding/json"
	"fmt"
)

type JSONStringArray []string

func (j *JSONStringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type User struct {
	AccountID   string          `json:"account_id"`
	EmployeeID  string          `json:"employee_id"`
	AccountType string          `json:"account_type"`
	Name        string          `json:"name"`
	MobileNo    string          `json:"mobile_no"`
	Password    string          `json:"password"`
	Role        string          `json:"role"`
	Permissions JSONStringArray `json:"permissions"`
	Status      string          `json:"status"`
}
