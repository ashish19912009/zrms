CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.team_accounts (
    account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mobile_no VARCHAR(15) NOT NULL,
    login_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    employee_id VARCHAR(50),
    password_hash TEXT NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    permissions TEXT[],
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_mobile_no ON users.team_accounts (mobile_no);
CREATE INDEX idx_users_login_id ON users.team_accounts (login_id);
CREATE INDEX idx_users_account_type ON users.team_accounts (account_type);