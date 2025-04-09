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
    '8888888889',
    'ops.manager_2',
    'Suresh',
    'ops_manager',
    'EMP003',
    '$2a$12$KbLkDETK6qple93GvvS/q.UdQbxidA5dvqIKZzz.W8f1gxxqdXHka', -- bcrypt hashed password
    'admin',
    ARRAY['read', 'write'],
    'active',
    now(),
    now()
);
