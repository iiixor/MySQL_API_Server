-- MySQL initialization script for MySQL TUI Editor
-- This script runs when the container starts for the first time

-- Ensure root has full privileges
FLUSH PRIVILEGES;

-- Set global variables for better performance with temporary databases
SET GLOBAL max_connections = 200;
SET GLOBAL connect_timeout = 10;
SET GLOBAL wait_timeout = 60;
SET GLOBAL interactive_timeout = 60;

-- Create a default database (optional, not used for student queries)
CREATE DATABASE IF NOT EXISTS mysql_tui_editor;
