-- Revert reply threading and edit/delete support
DROP INDEX IF EXISTS idx_messages_parent_id;
ALTER TABLE messages DROP COLUMN IF EXISTS is_deleted;
ALTER TABLE messages DROP COLUMN IF EXISTS is_edited;
ALTER TABLE messages DROP COLUMN IF EXISTS parent_id;
