-- Table for user activity logs
CREATE TABLE IF NOT EXISTS user_activity_logs (
    id UUID PRIMARY KEY,
    object_name VARCHAR(100) NOT NULL,
    record_id VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    changed_by VARCHAR(255) NOT NULL,
    request JSONB,
    response JSONB,
    change_stamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for typical audit log lookup patterns
CREATE INDEX IF NOT EXISTS idx_user_activity_logs_change_stamp ON user_activity_logs(change_stamp);
CREATE INDEX IF NOT EXISTS idx_user_activity_logs_changed_by ON user_activity_logs(changed_by);
CREATE INDEX IF NOT EXISTS idx_user_activity_logs_object_name_record_id ON user_activity_logs(object_name, record_id);
