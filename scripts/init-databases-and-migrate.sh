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
    
    # Check for dirty state and fix if needed
    migrate -path "/migrations/$service" \
            -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=require" \
            version 2>/dev/null || {
        echo "Database in dirty state, forcing clean..."
        migrate -path "/migrations/$service" \
                -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=require" \
                force 1
    }
    
    # Run migrations
    migrate -path "/migrations/$service" \
            -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=require" \
            up
    
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

# Keep the container running (for web service)
echo "Migration service ready - databases initialized!"
tail -f /dev/null