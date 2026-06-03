ALTER TABLE lessons
    DROP COLUMN IF EXISTS prerequisite_lesson_id;

DROP INDEX IF EXISTS idx_course_certificates_user_id;
DROP TABLE IF EXISTS course_certificates;

DROP TABLE IF EXISTS user_education_stats;

DROP INDEX IF EXISTS idx_lesson_progress_user_course;
DROP TABLE IF EXISTS lesson_progress;
