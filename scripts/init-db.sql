-- PostgreSQL Database Initialization Script for MysticFunds
-- This script creates all necessary databases for the microservices

-- Create databases for each service
CREATE DATABASE auth;
CREATE DATABASE wizard;
CREATE DATABASE mana;

-- Grant all privileges to mysticfunds user
GRANT ALL PRIVILEGES ON DATABASE auth TO mysticfunds;
GRANT ALL PRIVILEGES ON DATABASE wizard TO mysticfunds;
GRANT ALL PRIVILEGES ON DATABASE mana TO mysticfunds;

-- Set ownership
ALTER DATABASE auth OWNER TO mysticfunds;
ALTER DATABASE wizard OWNER TO mysticfunds;
ALTER DATABASE mana OWNER TO mysticfunds;