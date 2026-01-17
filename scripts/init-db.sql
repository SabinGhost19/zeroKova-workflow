-- Initialize database schemas for all services

-- Create schemas
CREATE SCHEMA IF NOT EXISTS orders;
CREATE SCHEMA IF NOT EXISTS inventory;
CREATE SCHEMA IF NOT EXISTS notifications;

-- Grant permissions
GRANT ALL ON SCHEMA orders TO postgres;
GRANT ALL ON SCHEMA inventory TO postgres;
GRANT ALL ON SCHEMA notifications TO postgres;
