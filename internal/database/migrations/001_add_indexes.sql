-- Create indexes for performance optimization

-- Add index for users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Add index for sessions table
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Add index for metrics table
CREATE INDEX idx_metrics_timestamp ON metrics(timestamp);
CREATE INDEX idx_metrics_type ON metrics(type);

-- Add composite indexes for common query patterns
CREATE INDEX idx_metrics_type_timestamp ON metrics(type, timestamp);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_age ON users(age);

-- Add composite indexes for common search patterns
CREATE INDEX IF NOT EXISTS idx_users_name_email ON users(name, email);
CREATE INDEX IF NOT EXISTS idx_users_age_name ON users(age, name);

-- Add indexes for sorting
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_updated_at ON users(updated_at DESC);

-- Down migration
-- Add rollback statements
DROP INDEX IF EXISTS idx_metrics_type_timestamp;
DROP INDEX IF EXISTS idx_metrics_type;
DROP INDEX IF EXISTS idx_metrics_timestamp;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;