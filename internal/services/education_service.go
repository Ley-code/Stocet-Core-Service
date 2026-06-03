package services

import (
	"encoding/json"
	"errors"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type EducationService struct {
	repo *repositories.EducationRepository
}

func NewEducationService(repo *repositories.EducationRepository) *EducationService {
	return &EducationService{repo: repo}
}

// Course methods
func (s *EducationService) CreateCourse(course *models.Course) error {
	if course.Title == "" {
		return errors.New("course title is required")
	}
	return s.repo.CreateCourse(course)
}

func (s *EducationService) GetCourseByID(id string) (*models.Course, error) {
	return s.repo.GetCourseByID(id)
}

func (s *EducationService) ListCourses(limit, offset int) ([]models.Course, int64, error) {
	return s.repo.ListCourses(limit, offset)
}

func (s *EducationService) UpdateCourse(course *models.Course) error {
	if course.Title == "" {
		return errors.New("course title is required")
	}
	return s.repo.UpdateCourse(course)
}

func (s *EducationService) DeleteCourse(id string) error {
	return s.repo.DeleteCourse(id)
}

// Lesson methods
func (s *EducationService) CreateLesson(lesson *models.Lesson) error {
	if lesson.Title == "" {
		return errors.New("lesson title is required")
	}
	return s.repo.CreateLesson(lesson)
}

func (s *EducationService) GetLessonByID(id string) (*models.Lesson, error) {
	return s.repo.GetLessonByID(id)
}

func (s *EducationService) UpdateLesson(lesson *models.Lesson) error {
	if lesson.Title == "" {
		return errors.New("lesson title is required")
	}
	return s.repo.UpdateLesson(lesson)
}

func (s *EducationService) DeleteLesson(id string) error {
	return s.repo.DeleteLesson(id)
}

// Quiz submission
func (s *EducationService) SubmitQuiz(userID, lessonID string, answers []int) (*models.QuizSubmission, error) {
	quiz, err := s.repo.GetQuizByLessonID(lessonID)
	if err != nil {
		return nil, errors.New("quiz not found")
	}

	if len(answers) != len(quiz.Questions) {
		return nil, errors.New("number of answers does not match number of questions")
	}

	// Calculate score
	correct := 0
	for i, question := range quiz.Questions {
		if i < len(answers) && answers[i] == question.Correct {
			correct++
		}
	}

	score := float64(correct) / float64(len(quiz.Questions)) * 100

	answersBytes, err := json.Marshal(answers)
	if err != nil {
		return nil, errors.New("failed to marshal answers")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	lessonUUID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, errors.New("invalid lesson ID")
	}
	quizUUID, err := uuid.Parse(quiz.ID)
	if err != nil {
		return nil, errors.New("invalid quiz ID")
	}

	submission := &models.QuizSubmission{
		UserID:   userUUID,
		LessonID: lessonUUID,
		QuizID:   quizUUID,
		Answers:  datatypes.JSON(answersBytes),
		Score:    score,
	}

	if err := s.repo.CreateQuizSubmission(submission); err != nil {
		return nil, err
	}

	return submission, nil
}

// Glossary methods
func (s *EducationService) CreateGlossaryTerm(term *models.GlossaryTerm) error {
	if term.Term == "" || term.Definition == "" {
		return errors.New("term and definition are required")
	}
	return s.repo.CreateGlossaryTerm(term)
}

func (s *EducationService) ListGlossaryTerms(limit, offset int, search string) ([]models.GlossaryTerm, int64, error) {
	return s.repo.ListGlossaryTerms(limit, offset, search)
}

func (s *EducationService) GetGlossaryTermByID(id string) (*models.GlossaryTerm, error) {
	return s.repo.GetGlossaryTermByID(id)
}

func (s *EducationService) UpdateGlossaryTerm(term *models.GlossaryTerm) error {
	if term.Term == "" || term.Definition == "" {
		return errors.New("term and definition are required")
	}
	return s.repo.UpdateGlossaryTerm(term)
}

func (s *EducationService) DeleteGlossaryTerm(id string) error {
	return s.repo.DeleteGlossaryTerm(id)
}
