#!/bin/sh

# Railway startup script for MysticFunds
echo "Starting MysticFunds on Railway..."

# Wait for database to be ready
echo "Waiting for database connection..."
until pg_isready -h $PGHOST -p $PGPORT -U $PGUSER; do
  echo "Database not ready, waiting..."
  sleep 2
done

echo "Database ready!"

# Run database migrations
echo "Running database migrations..."

# Create databases if they don't exist
echo "Creating databases..."
PGPASSWORD=$PGPASSWORD psql -h $PGHOST -p $PGPORT -U $PGUSER -d $PGDATABASE -c "CREATE DATABASE auth;" 2>/dev/null || echo "auth database exists"
PGPASSWORD=$PGPASSWORD psql -h $PGHOST -p $PGPORT -U $PGUSER -d $PGDATABASE -c "CREATE DATABASE wizard;" 2>/dev/null || echo "wizard database exists"
PGPASSWORD=$PGPASSWORD psql -h $PGHOST -p $PGPORT -U $PGUSER -d $PGDATABASE -c "CREATE DATABASE mana;" 2>/dev/null || echo "mana database exists"

# Run migrations for each service
echo "Running auth migrations..."
# Use environment variables for database connection
export AUTH_DB_HOST=$PGHOST
export AUTH_DB_PORT=$PGPORT
export AUTH_DB_USER=$PGUSER
export AUTH_DB_PASSWORD=$PGPASSWORD
export AUTH_DB_NAME=auth

export WIZARD_DB_HOST=$PGHOST
export WIZARD_DB_PORT=$PGPORT
export WIZARD_DB_USER=$PGUSER
export WIZARD_DB_PASSWORD=$PGPASSWORD
export WIZARD_DB_NAME=wizard

export MANA_DB_HOST=$PGHOST
export MANA_DB_PORT=$PGPORT
export MANA_DB_USER=$PGUSER
export MANA_DB_PASSWORD=$PGPASSWORD
export MANA_DB_NAME=mana

# Start services in background
echo "Starting auth service..."
PORT=50051 ./bin/auth &
AUTH_PID=$!

echo "Starting wizard service..."
PORT=50052 ./bin/wizard &
WIZARD_PID=$!

echo "Starting mana service..."
PORT=50053 ./bin/mana &
MANA_PID=$!

# Wait a moment for services to start
sleep 5

echo "Starting API gateway..."
PORT=$PORT ./bin/api-gateway &
GATEWAY_PID=$!

# Function to handle shutdown
cleanup() {
    echo "Shutting down services..."
    kill $AUTH_PID $WIZARD_PID $MANA_PID $GATEWAY_PID 2>/dev/null
    wait
    exit 0
}

# Trap signals
trap cleanup SIGTERM SIGINT

# Wait for all processes
wait