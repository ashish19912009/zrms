CREATE SCHEMA IF NOT EXISTS users;

-- Create the `accounts` table inside the `users` schema
CREATE TABLE IF NOT EXISTS users.accounts (
    id UUID PRIMARY KEY,
    mobile_no VARCHAR(15) NOT NULL UNIQUE,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    employee_id VARCHAR(50),
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_users_mobile_no ON users.accounts (mobile_no);
CREATE INDEX idx_users_employee_id ON users.accounts (employee_id);