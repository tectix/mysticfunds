#!/bin/bash
set -e

SERVICE=$1
if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service_name>"
    exit 1
fi

echo "Running migrations for $SERVICE service..."

# Install migrate tool if not present
if ! command -v migrate &> /dev/null; then
    echo "Installing migrate tool..."
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
    mv migrate /usr/local/bin/
fi

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
    echo "PostgreSQL is not ready yet. Waiting..."
    sleep 2
done

echo "PostgreSQL is ready!"

# Create database if it doesn't exist
echo "Creating database: $DB_NAME"
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME;"

echo "Database $DB_NAME ready"

# Check if migration directory exists
if [ ! -d "migrations/$SERVICE" ]; then
    echo "Warning: No migrations found for $SERVICE service"
    exit 0
fi

# Run migrations with error handling
echo "Running migrations for $SERVICE..."
migrate -path "migrations/$SERVICE" \
        -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require" \
        up 2>&1 | tee /tmp/migrate_output

# Check if migration failed due to dirty state
if grep -q "Dirty database version" /tmp/migrate_output; then
    echo "Database in dirty state, forcing clean and retrying..."
    migrate -path "migrations/$SERVICE" \
            -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require" \
            force 10
    
    # Retry migration
    migrate -path "migrations/$SERVICE" \
            -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require" \
            up
fi

if [ $? -eq 0 ]; then
    echo "$SERVICE migrations completed successfully"
else
    echo "$SERVICE migrations failed"
    exit 1
fi