package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Course struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Duration    string         `json:"duration"`
	Lessons     []Lesson       `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Lesson struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CourseID  string         `gorm:"type:uuid;not null;index" json:"course_id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	VideoURL  string         `gorm:"type:text" json:"video_url"`
	Order     int            `json:"order"`
	Quiz      *Quiz          `gorm:"foreignKey:LessonID" json:"quiz,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Quiz struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	LessonID  string     `gorm:"type:uuid;not null;index" json:"lesson_id"`
	Questions []Question `gorm:"foreignKey:QuizID" json:"questions"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Question struct {
	ID       string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	QuizID   string         `gorm:"type:uuid;not null;index" json:"quiz_id"`
	Question string         `gorm:"not null" json:"question"`
	Options  datatypes.JSON `gorm:"type:jsonb;serializer:json" json:"options"`
	Correct  int            `gorm:"not null" json:"correct"`
}

type GlossaryTerm struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Term       string         `gorm:"uniqueIndex;not null" json:"term"`
	Definition string         `gorm:"type:text;not null" json:"definition"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type QuizSubmission struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	LessonID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"lesson_id"`
	QuizID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"quiz_id"`
	Answers   datatypes.JSON `gorm:"type:jsonb;serializer:json" json:"answers"`
	Score     float64        `json:"score"`
	CreatedAt time.Time      `json:"created_at"`
}

// TableName methods
func (Course) TableName() string {
	return "courses"
}

func (Lesson) TableName() string {
	return "lessons"
}

func (Quiz) TableName() string {
	return "quizzes"
}

func (Question) TableName() string {
	return "questions"
}

func (GlossaryTerm) TableName() string {
	return "glossary_terms"
}

func (QuizSubmission) TableName() string {
	return "quiz_submissions"
}
