-- Disable foreign keys temporarily to avoid conflicts
SET session_replication_role = replica;

TRUNCATE TABLE cards RESTART IDENTITY CASCADE;
TRUNCATE TABLE columns RESTART IDENTITY CASCADE;
TRUNCATE TABLE boards RESTART IDENTITY CASCADE;

SET session_replication_role = DEFAULT;
