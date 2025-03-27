-- Make sure pgcrypto is enabled for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Seed Super Admin
INSERT INTO users.accounts (
  id, 
  mobile_no, 
  name, 
  role, 
  status, 
  employee_id, 
  created_at
) VALUES (
  gen_random_uuid(),              -- id (UUID)
  '7088669991',                   -- mobile_no
  'Super Admin',                  -- name
  'super_admin',                  -- role
  'active',                       -- status
  'EMP001',                       -- employee_id
  NOW()                           -- created_at
);