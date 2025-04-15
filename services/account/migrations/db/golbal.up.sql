CREATE SCHEMA IF NOT EXISTS global;

CREATE TABLE IF NOT EXISTS global.document_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,               -- e.g. FSSAI License, GST Certificate, PAN Card
    description TEXT,
    is_mandatory BOOLEAN DEFAULT false,
    status TEXT DEFAULT 'active',            -- 'active', 'inactive'
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS outlet.franchise_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES outlet.franchises(id) ON DELETE CASCADE,
    document_type_id UUID NOT NULL REFERENCES global.document_types(id) ON DELETE RESTRICT,
    document_url TEXT NOT NULL,
    uploaded_by UUID,                                       -- FK to outlet.owner or outlet.team_accounts
    status TEXT DEFAULT 'pending',                          -- pending, verified, rejected
    remarks TEXT,                                           -- reviewer comments
    uploaded_at TIMESTAMPTZ DEFAULT now(),
    verified_at TIMESTAMPTZ
);