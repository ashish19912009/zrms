
->	migrate create -ext sql -dir migrations/db -seq <name>

-> 	It will create two file 000001_create_account_table.up.sql and 000001_create_account_table.down.sql.

->	Then update both the sql file with sql query.

->	migrate -path=migrations/db -database "postgres://admin:12345@localhost:5432/zrms?sslmode=disable" -verbose up

 protoc --go_out=/Users/charuhashirena/Documents/Projects/zrms/services/account --go-grpc_out=/Users/charuhashirena/Documents/Projects/zrms/services/account services/account/proto/account.proto





 


STEP 1

docker pull postgres

STEP 2

docker ps -a


STEP 3

docker network create zippy-net

STEP 4

docker run --name pg_container --network zippy-net \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=12345 \
  -e POSTGRES_DB=zrms \
  -p 5432:5432 \
  -d postgres:15

STEP 5

docker run --rm \
--network zippy-net \
-v /Users/charuhashirena/Documents/Projects/zrms/services/account/migrations/db:/migrations \
migrate/migrate \
-path=/migrations \
-database "postgres://admin:12345@pg_container:5432/zrms?sslmode=disable" \
-verbose up




✅ 1. Verify Migration File Exists Inside the Mounted Volume
On your host machine:

bash
Copy
Edit
ls /Users/charuhashirena/Documents/Projects/zrms/services/account/migrations/db
Make sure there's a file like:

pgsql
Copy
Edit
000001_create_users_team_accounts.up.sql
Inside it, you should see SQL like:

sql
Copy
Edit
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.team_accounts (
  account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ...
);
✅ 2. Add Debugging to See Logs
To print migration logs to stdout (to debug), add -verbose flag:

bash
Copy
Edit
docker run --rm \
--network zippy-net \
-v /Users/charuhashirena/Documents/Projects/zrms/services/account/migrations/db:/migrations \
migrate/migrate \
-path=/migrations \
-database "postgres://admin:12345@pg_container:5432/zrms?sslmode=disable" \
-verbose up
This will show what SQL is being executed or whether it's being skipped due to already being applied.

✅ 3. Check schema_migrations Table
If the migration has run once successfully, golang-migrate won’t re-run it unless you:

Add a new file with a higher version number, or

Run:

bash
Copy
Edit
migrate ... force 0
To see what's applied:

bash
Copy
Edit
docker run --rm \
--network zippy-net \
-v /Users/charuhashirena/Documents/Projects/zrms/services/account/migrations/db:/migrations \
migrate/migrate \
-path=/migrations \
-database "postgres://admin:12345@pg_container:5432/zrms?sslmode=disable" \
version
✅ 4. Shell into the Container and Connect to Postgres
Just to verify from inside Docker that the table doesn’t exist:

bash
Copy
Edit
docker exec -it pg_container psql -U admin -d zrms
Then run:

sql
Copy
Edit
\dt users.*
To list all tables in the users schema.