
INSERT INTO outlet.owners (id, name, gender, dob, mobile_no, email, address, aadhar_no, is_verified, status, created_at, updated_at)
VALUES
  (
    gen_random_uuid(), -- UUID for owner ID -- franchise_id linked to the above franchise
    'Ashish Kumar Rena', -- Owner Name
    'male', -- Gender
    '1990-01-01'::TIMESTAMPTZ, -- Date of Birth (example)
    '7088669991', -- Mobile number
    'ashish@example.com', -- Owner email (you can adjust this)
    'Address Line 1, City, Country', -- Owner address (fill in details)
    '123456789012', -- Aadhar number (make sure it's unique)
    true, -- Verified (true or false)
    'active', -- Status (active, inactive, suspended)
    now(), -- created_at
    now() -- updated_at
  );

INSERT INTO outlet.franchises (id, business_name, subdomain, logo_url, theme_settings, status, created_at, updated_at, franchise_owner_id)
VALUES
  (
    gen_random_uuid(), -- UUID for franchise ID
    'rail food', -- Business Name
    'railfood', -- Subdomain (you can customize this based on your needs)
    NULL, -- logo_url (add the logo URL if available)
    '{}'::jsonb, -- theme_settings (empty for now)
    'active', -- Status
    now(), -- created_at
    now(), -- updated_at
	'bc3aa99e-71aa-4063-a66b-6b616ace0adf'
  )
  

INSERT INTO outlet.roles (id, franchise_id, name, description, is_default, created_at, updated_at)
  VALUES
    (gen_random_uuid(), NULL, 'admin', 'Full control over the franchise, able to manage all settings', true, now(), now()),
    (gen_random_uuid(), NULL, 'manager', 'Can manage orders, staff, and settings related to the franchise', false, now(), now()),
    (gen_random_uuid(), NULL, 'packer', 'Responsible for packing orders for delivery', false, now(), now()),
    (gen_random_uuid(), NULL, 'delivery_agent', 'Delivers packed food to customers', false, now(), now()),
    (gen_random_uuid(), NULL, 'cook', 'Prepares food for the orders', false, now(), now())

INSERT INTO outlet.permissions (key, description, created_at)
  VALUES
    ('view_orders', 'View the list of orders', now()),
    ('edit_orders', 'Edit the details of an order', now()),
    ('manage_staff', 'Manage team members and their roles', now()),
    ('manage_menu', 'Add, update, or remove menu items', now()),
    ('view_reports', 'View sales and performance reports', now()),
    ('deliver_food', 'Deliver food to customers', now()),
    ('pack_food', 'Pack the food for delivery', now()),
    ('cook_food', 'Prepare the food for orders', now())

-- Admin has all permissions
INSERT INTO outlet.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM outlet.roles r
JOIN outlet.permissions p
ON r.name = 'admin';

-- Manager can manage orders and staff
INSERT INTO outlet.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM outlet.roles r
JOIN outlet.permissions p
ON r.name = 'manager' AND p.key IN ('view_orders', 'edit_orders', 'manage_staff', 'manage_menu', 'view_reports');

-- Delivery Agent can only deliver food
INSERT INTO outlet.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM outlet.roles r
JOIN outlet.permissions p
ON r.name = 'delivery_agent' AND p.key = 'deliver_food';

-- Cook can only prepare food
INSERT INTO outlet.role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM outlet.roles r
JOIN outlet.permissions p
ON r.name = 'cook' AND p.key = 'cook_food';
