package model

import "time"

type Account struct {
	ID         string     `json:"id"`          // UUID for uniqueness
	MobileNo   string     `json:"mobile_no"`   // Unique mobile number
	EmployeeID string     `json:"employee_id"` // Unique employee ID (for admins/delivery partners)
	Name       string     `json:"name"`        // User's name
	Role       string     `json:"role"`        // Role: super_admin, admin, delivery_partner
	Status     string     `json:"status"`      // active, inactive, suspended
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"` // Nullable field for soft deletion
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
