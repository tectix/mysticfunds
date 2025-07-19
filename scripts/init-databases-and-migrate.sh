#!/bin/bash
set -e

echo "Initializing databases and running migrations..."

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
    echo "PostgreSQL is not ready yet. Waiting..."
    sleep 2
done

echo "PostgreSQL is ready!"

# Function to create database if it doesn't exist
create_database() {
    local db_name=$1
    echo "Creating database: $db_name"
    
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$db_name'" | grep -q 1 || \
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $db_name;"
    
    echo "Database $db_name ready"
}

# Function to run migrations for a service
run_migrations() {
    local service=$1
    local db_name=$2
    
    echo "Running migrations for $service service..."
    
    # Check if migration directory exists
    if [ ! -d "/migrations/$service" ]; then
        echo "Warning: No migrations found for $service service"
        return 0
    fi
    
    # Try to run migrations, if dirty state detected, force clean and retry
    migrate -path "/migrations/$service" \
            -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=require" \
            up 2>&1 | tee /tmp/migrate_output
    
    # Check if migration failed due to dirty state
    if grep -q "Dirty database version" /tmp/migrate_output; then
        echo "Database in dirty state, forcing clean and retrying..."
        migrate -path "/migrations/$service" \
                -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=require" \
                force 10
        
        # Retry migration
        migrate -path "/migrations/$service" \
                -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=require" \
                up
    fi
    
    if [ $? -eq 0 ]; then
        echo "$service migrations completed successfully"
    else
        echo "$service migrations failed"
        exit 1
    fi
}

# Create databases
create_database "auth"
create_database "wizard" 
create_database "mana"

# Run migrations for each service
run_migrations "auth" "auth"
run_migrations "wizard" "wizard" 
run_migrations "mana" "mana"

echo "All databases created and migrations completed successfully!"

# Exit successfully - job is complete
echo "Migration job completed successfully!"
exit 0