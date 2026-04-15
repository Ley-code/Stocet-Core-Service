-- ============================================================================
-- Migration 000001 DOWN: Revert Existing Schema Optimizations
-- ============================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- 1. Drop foreign key constraints added in the UP migration
-- -------------------------------------------------------------------------
ALTER TABLE questions DROP CONSTRAINT IF EXISTS fk_questions_quiz;
ALTER TABLE quizzes DROP CONSTRAINT IF EXISTS fk_quizzes_lesson;
ALTER TABLE lessons DROP CONSTRAINT IF EXISTS fk_lessons_course;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_channel;
ALTER TABLE watchlists DROP CONSTRAINT IF EXISTS fk_watchlists_user;
ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS fk_subscriptions_user;
ALTER TABLE quiz_submissions DROP CONSTRAINT IF EXISTS fk_quiz_submissions_user;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS fk_messages_user;
ALTER TABLE alerts DROP CONSTRAINT IF EXISTS fk_alerts_user;

-- -------------------------------------------------------------------------
-- 2. Drop indexes
-- -------------------------------------------------------------------------
DROP INDEX IF EXISTS idx_questions_quiz_id;
DROP INDEX IF EXISTS idx_quizzes_lesson_id;
DROP INDEX IF EXISTS idx_lessons_course_id;
DROP INDEX IF EXISTS idx_messages_channel_id;
DROP INDEX IF EXISTS idx_watchlists_ticker;
DROP INDEX IF EXISTS idx_alerts_ticker;

-- -------------------------------------------------------------------------
-- 3. Revert quiz_submissions.lesson_id and quiz_id from UUID back to TEXT
-- -------------------------------------------------------------------------
ALTER TABLE quiz_submissions
    ALTER COLUMN lesson_id TYPE text USING lesson_id::text,
    ALTER COLUMN quiz_id TYPE text USING quiz_id::text;

-- -------------------------------------------------------------------------
-- 4. Revert JSONB columns back to arrays
-- -------------------------------------------------------------------------

-- quiz_submissions.answers: JSONB → integer[]
ALTER TABLE quiz_submissions
    ALTER COLUMN answers TYPE integer[]
    USING ARRAY(SELECT jsonb_array_elements_text(COALESCE(answers, '[]'::jsonb))::integer);

-- questions.options: JSONB → text[]
ALTER TABLE questions
    ALTER COLUMN options TYPE text[]
    USING ARRAY(SELECT jsonb_array_elements_text(COALESCE(options, '[]'::jsonb)));

-- -------------------------------------------------------------------------
-- 5. Revert user_id from UUID back to TEXT
-- -------------------------------------------------------------------------
ALTER TABLE watchlists ALTER COLUMN user_id TYPE text USING user_id::text;
ALTER TABLE subscriptions ALTER COLUMN user_id TYPE text USING user_id::text;
ALTER TABLE quiz_submissions ALTER COLUMN user_id TYPE text USING user_id::text;
ALTER TABLE messages ALTER COLUMN user_id TYPE text USING user_id::text;
ALTER TABLE alerts ALTER COLUMN user_id TYPE text USING user_id::text;

-- -------------------------------------------------------------------------
-- 6. Revert ENUM columns back to plain TEXT
-- -------------------------------------------------------------------------
ALTER TABLE users
    ALTER COLUMN role TYPE text USING role::text,
    ALTER COLUMN role SET DEFAULT 'user',
    ALTER COLUMN subscription_status TYPE text USING subscription_status::text,
    ALTER COLUMN subscription_status SET DEFAULT 'free';

-- -------------------------------------------------------------------------
-- 7. Drop ENUM types
-- -------------------------------------------------------------------------
DROP TYPE IF EXISTS subscription_status_enum;
DROP TYPE IF EXISTS user_role;

COMMIT;
