package repositories

import (
	"errors"
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm/clause"
	"gorm.io/gorm"
)

type EducationRepository struct{}

func NewEducationRepository() *EducationRepository {
	return &EducationRepository{}
}

// Course methods
func (r *EducationRepository) CreateCourse(course *models.Course) error {
	return database.DB.Create(course).Error
}

func (r *EducationRepository) GetCourseByID(id string) (*models.Course, error) {
	var course models.Course
	err := database.DB.Preload("Lessons.Quiz.Questions").Where("id = ?", id).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *EducationRepository) ListCourses(limit, offset int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	if err := database.DB.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "course_id", "title", "order", "video_url", "created_at", "updated_at")
	}).Limit(limit).Offset(offset).Find(&courses).Error
	return courses, total, err
}

func (r *EducationRepository) UpdateCourse(course *models.Course) error {
	return database.DB.Save(course).Error
}

func (r *EducationRepository) DeleteCourse(id string) error {
	return database.DB.Delete(&models.Course{}, "id = ?", id).Error
}

// Lesson methods
func (r *EducationRepository) CreateLesson(lesson *models.Lesson) error {
	return database.DB.Create(lesson).Error
}

func (r *EducationRepository) GetLessonByID(id string) (*models.Lesson, error) {
	var lesson models.Lesson
	err := database.DB.Preload("Quiz.Questions").Where("id = ?", id).First(&lesson).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *EducationRepository) UpdateLesson(lesson *models.Lesson) error {
	return database.DB.Save(lesson).Error
}

func (r *EducationRepository) DeleteLesson(id string) error {
	return database.DB.Delete(&models.Lesson{}, "id = ?", id).Error
}

// Quiz methods
func (r *EducationRepository) CreateQuiz(quiz *models.Quiz) error {
	return database.DB.Create(quiz).Error
}

func (r *EducationRepository) GetQuizByLessonID(lessonID string) (*models.Quiz, error) {
	var quiz models.Quiz
	err := database.DB.Preload("Questions").Where("lesson_id = ?", lessonID).First(&quiz).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

// Quiz Submission methods
func (r *EducationRepository) CreateQuizSubmission(submission *models.QuizSubmission) error {
	return database.DB.Create(submission).Error
}

func (r *EducationRepository) GetQuizSubmissionsByUser(userID string, limit, offset int) ([]models.QuizSubmission, int64, error) {
	var submissions []models.QuizSubmission
	var total int64

	if err := database.DB.Model(&models.QuizSubmission{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at DESC").Find(&submissions).Error
	return submissions, total, err
}

// Glossary methods
func (r *EducationRepository) CreateGlossaryTerm(term *models.GlossaryTerm) error {
	return database.DB.Create(term).Error
}

func (r *EducationRepository) ListGlossaryTerms(limit, offset int, search string) ([]models.GlossaryTerm, int64, error) {
	var terms []models.GlossaryTerm
	var total int64

	query := database.DB.Model(&models.GlossaryTerm{})
	if search != "" {
		query = query.Where("term ILIKE ? OR definition ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("term ASC").Find(&terms).Error
	return terms, total, err
}

func (r *EducationRepository) GetGlossaryTermByID(id string) (*models.GlossaryTerm, error) {
	var term models.GlossaryTerm
	err := database.DB.Where("id = ?", id).First(&term).Error
	if err != nil {
		return nil, err
	}
	return &term, nil
}

func (r *EducationRepository) UpdateGlossaryTerm(term *models.GlossaryTerm) error {
	return database.DB.Save(term).Error
}

func (r *EducationRepository) DeleteGlossaryTerm(id string) error {
	return database.DB.Delete(&models.GlossaryTerm{}, "id = ?", id).Error
}

func (r *EducationRepository) GetCourseProgress(userID, courseID string) ([]string, int64, error) {
	var completedLessonIDs []string
	if err := database.DB.Model(&models.LessonProgress{}).
		Where("user_id = ? AND course_id = ? AND completed = true", userID, courseID).
		Pluck("lesson_id", &completedLessonIDs).Error; err != nil {
		return nil, 0, err
	}

	var totalLessons int64
	if err := database.DB.Model(&models.Lesson{}).
		Where("course_id = ?", courseID).
		Count(&totalLessons).Error; err != nil {
		return nil, 0, err
	}

	return completedLessonIDs, totalLessons, nil
}

func (r *EducationRepository) MarkLessonComplete(userID, lessonID, courseID string) (*models.UserEducationStats, bool, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, false, err
	}
	lessonUUID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, false, err
	}
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return nil, false, err
	}

	var updatedStats models.UserEducationStats
	awardedXP := false

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var lesson models.Lesson
		if err := tx.Select("id", "course_id").Where("id = ?", lessonID).First(&lesson).Error; err != nil {
			return err
		}
		if lesson.CourseID != courseID {
			return errors.New("lesson does not belong to course")
		}

		var existing models.LessonProgress
		err := tx.Where("user_id = ? AND lesson_id = ?", userUUID, lessonUUID).First(&existing).Error
		if err == nil && existing.Completed {
			var stats models.UserEducationStats
			if err := tx.Where("user_id = ?", userUUID).First(&stats).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					stats = models.UserEducationStats{UserID: userUUID}
				} else {
					return err
				}
			}
			updatedStats = stats
			return nil
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		now := time.Now().UTC()
		progress := models.LessonProgress{
			UserID:      userUUID,
			LessonID:    lessonUUID,
			CourseID:    courseUUID,
			Completed:   true,
			CompletedAt: &now,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "lesson_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"completed":    true,
				"completed_at": now,
				"course_id":    courseUUID,
			}),
		}).Create(&progress).Error; err != nil {
			return err
		}

		awardedXP = true
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		yesterday := today.AddDate(0, 0, -1)

		var stats models.UserEducationStats
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userUUID).First(&stats).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stats = models.UserEducationStats{
				UserID:           userUUID,
				TotalXP:          10,
				CurrentStreak:    1,
				LongestStreak:    1,
				LastActivityDate: &today,
			}
			if err := tx.Create(&stats).Error; err != nil {
				return err
			}
			updatedStats = stats
			return nil
		}
		if err != nil {
			return err
		}

		stats.TotalXP += 10
		if stats.LastActivityDate == nil {
			stats.CurrentStreak = 1
		} else {
			last := time.Date(stats.LastActivityDate.Year(), stats.LastActivityDate.Month(), stats.LastActivityDate.Day(), 0, 0, 0, 0, time.UTC)
			if last.Equal(today) {
				// Same day completion; keep current streak unchanged.
			} else if last.Equal(yesterday) {
				stats.CurrentStreak++
			} else {
				stats.CurrentStreak = 1
			}
		}
		if stats.CurrentStreak > stats.LongestStreak {
			stats.LongestStreak = stats.CurrentStreak
		}
		stats.LastActivityDate = &today

		if err := tx.Save(&stats).Error; err != nil {
			return err
		}
		updatedStats = stats
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return &updatedStats, awardedXP, nil
}

func (r *EducationRepository) GetUserStats(userID string) (*models.UserEducationStats, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	var stats models.UserEducationStats
	err = database.DB.Where("user_id = ?", userUUID).First(&stats).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stats = models.UserEducationStats{UserID: userUUID}
		return &stats, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *EducationRepository) IssueCertificate(userID, courseID string) (*models.CourseCertificate, bool, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, false, err
	}
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return nil, false, err
	}

	var certificate models.CourseCertificate
	err = database.DB.Where("user_id = ? AND course_id = ?", userUUID, courseUUID).First(&certificate).Error
	if err == nil {
		return &certificate, false, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}

	var totalLessons int64
	if err := database.DB.Model(&models.Lesson{}).Where("course_id = ?", courseUUID).Count(&totalLessons).Error; err != nil {
		return nil, false, err
	}
	if totalLessons == 0 {
		return nil, false, nil
	}

	var completedLessons int64
	if err := database.DB.Model(&models.LessonProgress{}).
		Where("user_id = ? AND course_id = ? AND completed = true", userUUID, courseUUID).
		Count(&completedLessons).Error; err != nil {
		return nil, false, err
	}
	if completedLessons < totalLessons {
		return nil, false, nil
	}

	certificate = models.CourseCertificate{
		UserID:   userUUID,
		CourseID: courseUUID,
		IssuedAt: time.Now().UTC(),
	}
	if err := database.DB.Create(&certificate).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			var existing models.CourseCertificate
			if lookupErr := database.DB.Where("user_id = ? AND course_id = ?", userUUID, courseUUID).First(&existing).Error; lookupErr == nil {
				return &existing, false, nil
			}
		}
		return nil, false, err
	}

	return &certificate, true, nil
}

func (r *EducationRepository) GetCertificates(userID string) ([]models.CourseCertificate, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	var certificates []models.CourseCertificate
	if err := database.DB.Where("user_id = ?", userUUID).Order("issued_at DESC").Find(&certificates).Error; err != nil {
		return nil, err
	}
	return certificates, nil
}

func (r *EducationRepository) ReorderLessons(courseID string, orderedLessonIDs []string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for idx, lessonID := range orderedLessonIDs {
			if err := tx.Model(&models.Lesson{}).
				Where("id = ? AND course_id = ?", lessonID, courseID).
				Update("order", idx+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
