CREATE SCHEMA IF NOT EXISTS outlet;

CREATE TABLE IF NOT EXISTS outlet.franchises (
    id UUID PRIMARY KEY,
    business_name TEXT NOT NULL,
    subdomain TEXT UNIQUE NOT NULL,
    logo_url TEXT,
    theme_settings JSONB DEFAULT '{}'::jsonb,
    status TEXT DEFAULT 'active', -- 'active', 'inactive', 'suspended'
    franchise_owner_id UUID NOT NULL, -- person who opened the franchise
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS outlet.owner (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    gender TEXT NOT NULL,
    dob TIMESTAMPTZ NOT NULL,
    mobile_no TEXT,
    email TEXT,
    address TEXT,
    aadhar_no TEXT UNIQUE NOT NULL,
    is_verified Boolean,
    status TEXT DEFAULT 'active', -- 'active', 'inactive', 'suspended'
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- New pivot table to link owners to franchises (many-to-many)
CREATE TABLE IF NOT EXISTS outlet.franchise_owners (
    id UUID PRIMARY KEY,
    franchise_id UUID NOT NULL REFERENCES outlet.franchises(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL REFERENCES outlet.owners(id) ON DELETE CASCADE,
    is_primary BOOLEAN DEFAULT FALSE, -- mark main owner if needed
    ownership_percentage NUMERIC(5, 2), -- optional: share in %
    role TEXT DEFAULT 'owner', -- 'owner', 'partner', 'manager', etc.
    joined_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(franchise_id, owner_id) -- prevent duplicates
);

ALTER TABLE outlet.franchises
ADD CONSTRAINT fk_franchise_owner
FOREIGN KEY (franchise_owner_id) REFERENCES outlet.owners(id) ON DELETE RESTRICT;

CREATE OR REPLACE FUNCTION outlet.ensure_creator_in_franchise_owners()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the owner exists in franchise_owners
    IF NOT EXISTS (
        SELECT 1 FROM outlet.franchise_owners
        WHERE franchise_id = NEW.id
          AND owner_id = NEW.franchise_owner_id
    ) THEN
        -- If not, raise an error
        RAISE EXCEPTION 'Franchise owner must be part of franchise_owners table';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_franchise_owner_exists
AFTER INSERT OR UPDATE ON outlet.franchises
FOR EACH ROW
EXECUTE FUNCTION outlet.ensure_creator_in_franchise_owners();

CREATE TABLE IF NOT EXISTS outlet.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES outlet.franchises(id) ON DELETE CASCADE,
    name TEXT NOT NULL, -- e.g., "manager", "packer", "delivery"
    description TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),

    UNIQUE(franchise_id, name)
);

CREATE TABLE IF NOT EXISTS outlet.permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT UNIQUE NOT NULL, -- e.g., "view_orders", "edit_menu", "manage_team"
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS outlet.role_permissions (
    role_id UUID NOT NULL REFERENCES outlet.roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES outlet.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY(role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS outlet.direct_permissions (
    account_id UUID NOT NULL REFERENCES outlet.team_accounts(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES outlet.permissions(id) ON DELETE CASCADE,
    is_granted BOOLEAN NOT NULL DEFAULT true, -- true: allow, false: deny
    PRIMARY KEY(account_id, permission_id)
);

CREATE TABLE IF NOT EXISTS outlet.team_accounts (
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

    UNIQUE(franchise_id, login_id),
);

CREATE TABLE IF NOT EXISTS outlet.franchise_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES outlet.franchises(id) ON DELETE CASCADE,
    address_line TEXT,
    city TEXT,
    state TEXT,
    country TEXT DEFAULT 'India',
    pincode VARCHAR(10),
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    is_verified Boolean,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_franchise_pincode ON outlet.franchise_addresses (pincode);
CREATE INDEX idx_franchise_id ON outlet.franchise_addresses (franchise_id);
CREATE INDEX idx_franchise_id ON outlet.team_accounts (franchise_id);
CREATE INDEX idx_users_mobile_no ON outlet.team_accounts (mobile_no);
CREATE INDEX idx_users_login_id ON outlet.team_accounts (login_id);
CREATE INDEX idx_users_account_type ON outlet.team_accounts (account_type);