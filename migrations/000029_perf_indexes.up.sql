CREATE INDEX IF NOT EXISTS idx_alerts_user_created ON alerts (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_channel_created ON messages (channel_id, created_at DESC);
