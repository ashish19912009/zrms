INSERT INTO users.team_accounts (
    account_id,
    mobile_no,
    login_id,
    name,
    role,
    employee_id,
    password_hash,
    account_type,
    permissions,
    status,
    created_at,
    updated_at
)
VALUES
(
    gen_random_uuid(),
    '9999999999',
    'admin.zrms',
    'ZRMS Super Admin',
    'super_admin',
    'EMP001',
    '$2a$12$KbLkDETK6qple93GvvS/q.UdQbxidA5dvqIKZzz.W8f1gxxqdXHka', -- bcrypt hashed password
    'super_admin',
    ARRAY['read', 'write', 'delete'],
    'active',
    now(),
    now()
),
(
    gen_random_uuid(),
    '8888888888',
    'ops.manager',
    'Operations Manager',
    'ops_manager',
    'EMP002',
    '$2a$12$KbLkDETK6qple93GvvS/q.UdQbxidA5dvqIKZzz.W8f1gxxqdXHka', -- bcrypt hashed password
    'admin',
    ARRAY['read', 'write'],
    'active',
    now(),
    now()
);
