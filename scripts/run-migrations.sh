#!/bin/sh

# Migration runner script for MysticFunds
set -e

echo "Waiting for PostgreSQL to be ready..."

# Wait for PostgreSQL to be ready
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
    echo "PostgreSQL is not ready yet. Waiting..."
    sleep 2
done

echo "PostgreSQL is ready. Running migrations..."

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
    
    # Run migrations
    migrate -path "/migrations/$service" \
            -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$db_name?sslmode=disable" \
            up
    
    if [ $? -eq 0 ]; then
        echo "✅ $service migrations completed successfully"
    else
        echo "❌ $service migrations failed"
        exit 1
    fi
}

# Run migrations for each service
run_migrations "auth" "auth"
run_migrations "wizard" "wizard" 
run_migrations "mana" "mana"

echo "🎉 All migrations completed successfully!"