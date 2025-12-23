-- Create test database if it doesn't exist
-- This script is run by Docker on first container initialization
-- or manually via `make create-test-db`

CREATE DATABASE phd_test;
