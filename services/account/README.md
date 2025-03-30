# Run app (using local env by default)
make run

# Run with a specific environment
make ENV=staging run

# Seed or reset DB
make seed
make reset-db

# Run migrations
make migrate-up
make migrate-down

# Generate mocks
make gen-mock
