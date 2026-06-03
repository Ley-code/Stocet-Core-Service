CREATE TABLE IF NOT EXISTS lesson_progress (
    user_id UUID NOT NULL,
    lesson_id UUID NOT NULL,
    course_id UUID NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ NULL,
    PRIMARY KEY (user_id, lesson_id)
);

CREATE INDEX IF NOT EXISTS idx_lesson_progress_user_course
    ON lesson_progress (user_id, course_id);

CREATE TABLE IF NOT EXISTS user_education_stats (
    user_id UUID PRIMARY KEY,
    total_xp INT NOT NULL DEFAULT 0,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    last_activity_date DATE NULL
);

CREATE TABLE IF NOT EXISTS course_certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_course_certificates_user_id
    ON course_certificates (user_id);

ALTER TABLE lessons
    ADD COLUMN IF NOT EXISTS prerequisite_lesson_id UUID REFERENCES lessons(id);
