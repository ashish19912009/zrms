-- Drop the table and schema on rollback

DROP INDEX IF EXISTS idx_users_employee_id;
DROP INDEX IF EXISTS idx_users_mobile_no;
DROP TABLE IF EXISTS users.accounts;
DROP SCHEMA IF EXISTS users;