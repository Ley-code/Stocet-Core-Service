-- ============================================================================
-- Migration 000001: Optimize Existing PostgreSQL Schema
-- ============================================================================
-- This migration:
--   1. Creates ENUMs for user role and subscription status
--   2. Converts role/subscription_status columns to use ENUMs
--   3. Changes user_id from TEXT to UUID with proper FOREIGN KEY constraints
--   4. Converts array columns (text[], integer[]) to JSONB
--   5. Adds missing indexes on frequent lookup columns
-- ============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- 1. Create ENUM types
-- -------------------------------------------------------------------------
CREATE TYPE user_role AS ENUM ('admin', 'user');
CREATE TYPE subscription_status_enum AS ENUM ('free', 'premium');

-- -------------------------------------------------------------------------
-- 2. Convert users.role and users.subscription_status to ENUMs
-- -------------------------------------------------------------------------
ALTER TABLE users
    ALTER COLUMN role DROP DEFAULT,
    ALTER COLUMN role TYPE user_role USING role::user_role,
    ALTER COLUMN role SET DEFAULT 'user'::user_role,
    
    ALTER COLUMN subscription_status DROP DEFAULT,
    ALTER COLUMN subscription_status TYPE subscription_status_enum USING subscription_status::subscription_status_enum,
    ALTER COLUMN subscription_status SET DEFAULT 'free'::subscription_status_enum;

-- -------------------------------------------------------------------------
-- 3. Convert user_id from TEXT to UUID and add FOREIGN KEY constraints
-- -------------------------------------------------------------------------

-- alerts
ALTER TABLE alerts
    ALTER COLUMN user_id TYPE uuid USING user_id::uuid;
ALTER TABLE alerts
    ADD CONSTRAINT fk_alerts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- messages
ALTER TABLE messages
    ALTER COLUMN user_id TYPE uuid USING user_id::uuid;
ALTER TABLE messages
    ADD CONSTRAINT fk_messages_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- quiz_submissions
ALTER TABLE quiz_submissions
    ALTER COLUMN user_id TYPE uuid USING user_id::uuid;
ALTER TABLE quiz_submissions
    ADD CONSTRAINT fk_quiz_submissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- subscriptions
ALTER TABLE subscriptions
    ALTER COLUMN user_id TYPE uuid USING user_id::uuid;
ALTER TABLE subscriptions
    ADD CONSTRAINT fk_subscriptions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- watchlists
ALTER TABLE watchlists
    ALTER COLUMN user_id TYPE uuid USING user_id::uuid;
ALTER TABLE watchlists
    ADD CONSTRAINT fk_watchlists_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- -------------------------------------------------------------------------
-- 4. Convert array columns to JSONB
-- -------------------------------------------------------------------------

-- questions.options: text[] → JSONB
ALTER TABLE questions
    ALTER COLUMN options TYPE jsonb USING to_jsonb(options);

-- quiz_submissions.answers: integer[] → JSONB
ALTER TABLE quiz_submissions
    ALTER COLUMN answers TYPE jsonb USING to_jsonb(answers);

-- Also convert quiz_submissions.lesson_id and quiz_id from TEXT to UUID
ALTER TABLE quiz_submissions
    ALTER COLUMN lesson_id TYPE uuid USING lesson_id::uuid,
    ALTER COLUMN quiz_id TYPE uuid USING quiz_id::uuid;

-- -------------------------------------------------------------------------
-- 5. Add missing indexes
-- -------------------------------------------------------------------------

-- alerts: ticker for price-alert lookups
CREATE INDEX IF NOT EXISTS idx_alerts_ticker ON alerts (ticker);

-- watchlists: ticker for watchlist lookups
CREATE INDEX IF NOT EXISTS idx_watchlists_ticker ON watchlists (ticker);

-- messages: channel_id for message retrieval
CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messages (channel_id);

-- lessons: course_id for lesson listing
CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons (course_id);

-- quizzes: lesson_id for quiz lookup
CREATE INDEX IF NOT EXISTS idx_quizzes_lesson_id ON quizzes (lesson_id);

-- questions: quiz_id for question listing
CREATE INDEX IF NOT EXISTS idx_questions_quiz_id ON questions (quiz_id);

-- messages: foreign key on channel_id
ALTER TABLE messages
    ADD CONSTRAINT fk_messages_channel FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;

-- lessons: foreign key on course_id
ALTER TABLE lessons
    ADD CONSTRAINT fk_lessons_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;

-- quizzes: foreign key on lesson_id
ALTER TABLE quizzes
    ADD CONSTRAINT fk_quizzes_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE;

-- questions: foreign key on quiz_id
ALTER TABLE questions
    ADD CONSTRAINT fk_questions_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE;

COMMIT;
