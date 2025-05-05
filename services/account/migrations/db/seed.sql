-- Sample owner
INSERT INTO outlet.owners (
    id, name, gender, dob, mobile_no, email, address, aadhar_no, is_verified
) VALUES (
    '11111111-1111-1111-1111-111111111111', 'Amit Sharma', 'Male', '1985-04-15',
    '9876543210', 'amit@example.com', '123 Main Street, Delhi', '123456789012', true
);

-- Sample franchise created by owner
INSERT INTO outlet.franchises (
    id, business_name, subdomain, logo_url, theme_settings, status, franchise_owner_id
) VALUES (
    '22222222-2222-2222-2222-222222222222', 'Zippy Snacks', 'zippy-delhi', NULL, '{}', 'active',
    '11111111-1111-1111-1111-111111111111'
);

-- Owner-franchise mapping (can support co-owners later)
INSERT INTO outlet.franchise_owners (
    id, franchise_id, owner_id, is_primary, ownership_percentage, role
) VALUES (
    '33333333-3333-3333-3333-333333333333', '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111', true, 100.0, 'owner'
);

-- Franchise address
INSERT INTO outlet.franchise_addresses (
    franchise_id, address_line, city, state, pincode, is_verified, latitude, longitude
) VALUES (
    '22222222-2222-2222-2222-222222222222',
    'Shop No. 5, Near Metro Station', 'New Delhi', 'Delhi', '110001', true, 28.6139, 77.2090
);

-- Roles
INSERT INTO outlet.roles (
    id, franchise_id, name, description, is_default
) VALUES 
    ('44444444-4444-4444-4444-444444444444', '22222222-2222-2222-2222-222222222222', 'Manager', 'Can manage everything', true),
    ('55555555-5555-5555-5555-555555555555', '22222222-2222-2222-2222-222222222222', 'Delivery', 'Handles delivery tasks', false);

-- Permissions
INSERT INTO outlet.permissions (id, resource, action, key, description) VALUES
    ('aaa11111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'order', 'view', 'order:view', 'View orders'),
    ('aaa11112-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'order', 'edit', 'order:edit', 'Edit orders'),
    ('aaa11113-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'menu', 'view', 'menu:view', 'View menu'),
    ('aaa11114-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'account', 'manage', 'account:manage', 'Manage team accounts');

-- Role-Permission mapping
INSERT INTO outlet.role_permissions (role_id, permission_id) VALUES
    ('44444444-4444-4444-4444-444444444444', 'aaa11111-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('44444444-4444-4444-4444-444444444444', 'aaa11112-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('44444444-4444-4444-4444-444444444444', 'aaa11114-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('55555555-5555-5555-5555-555555555555', 'aaa11111-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

-- Team accounts
INSERT INTO outlet.team_accounts (
    id, franchise_id, employee_id, login_id, password_hash, account_type, name, mobile_no, email, role_id
) VALUES (
    '66666666-6666-6666-6666-666666666666', '22222222-2222-2222-2222-222222222222',
    'EMP001', 'manager@zippy.com', '$2a$10$examplehashhere1234567890123456789012', 'manager',
    'Rohit Manager', '9123456780', 'rohit@zippy.com', '44444444-4444-4444-4444-444444444444'
);

-- Direct permission (example of override)
INSERT INTO outlet.direct_permissions (account_id, permission_id, is_granted) VALUES
    ('66666666-6666-6666-6666-666666666666', 'aaa11113-aaaa-aaaa-aaaa-aaaaaaaaaaaa', true);
