package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"

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
