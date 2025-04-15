package model

import "time"

/**
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   franchise_id UUID NOT NULL REFERENCES outlet.franchises(id) ON DELETE CASCADE,
   employee_id VARCHAR(20),
   login_id TEXT NOT NULL UNIQUE,
   password_hash TEXT NOT NULL,
   account_type VARCHAR(50) NOT NULL, -- delivery_partner, food_packer, etc.
   name TEXT NOT NULL,
   mobile_no VARCHAR(15) NOT NULL,
   email TEXT,
   role_id UUID REFERENCES outlet.roles(id) ON DELETE SET NULL,
   status TEXT DEFAULT 'active',
   created_at TIMESTAMPTZ DEFAULT now(),
   updated_at TIMESTAMPTZ DEFAULT now(),
   deleted_at TIMESTAMPTZ,
*/

type FranchiseAccount struct {
	FranchiseID string     `json:"franchise_id"`
	EmployeeID  string     `json:"employee_id"` // Unique employee ID (for admins/delivery partners)
	LoginID     string     `json:"login_id"`
	AccountType string     `json:"account_type"`
	Name        string     `json:"name"`      // User's name
	MobileNo    string     `json:"mobile_no"` // Unique mobile number
	Email       string     `json:"email"`
	RoleID      string     `json:"role_id"` // Role: super_admin, admin, delivery_partner
	Status      string     `json:"status"`  // active, inactive, suspended
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"` // Nullable field for soft deletion
}

type FranchiseAccountResponse struct {
	ID          string     `json:"id"` // UUID for uniqueness
	FranchiseID string     `json:"franchise_id"`
	EmployeeID  string     `json:"employee_id"` // Unique employee ID (for admins/delivery partners)
	LoginID     string     `json:"login_id"`
	AccountType string     `json:"account_type"`
	Name        string     `json:"name"`      // User's name
	MobileNo    string     `json:"mobile_no"` // Unique mobile number
	Email       string     `json:"email"`
	RoleID      string     `json:"role_id"` // Role: super_admin, admin, delivery_partner
	Status      string     `json:"status"`  // active, inactive, suspended
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"` // Nullable field for soft deletion
}

/*

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    mobile_no VARCHAR(15) NOT NULL,
    name TEXT,
    role TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    employee_id VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for faster lookups
CREATE INDEX idx_accounts_mobile_no ON accounts(mobile_no);
CREATE INDEX idx_accounts_employee_id ON accounts(employee_id);
CREATE INDEX idx_accounts_deleted_at ON accounts(deleted_at);
*/
