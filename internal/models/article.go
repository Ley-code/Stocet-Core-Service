package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Article struct {
	ID                  uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title               string         `gorm:"column:title;not null" json:"title"`
	URL                 string         `gorm:"column:url;not null;uniqueIndex" json:"url"`
	Source              string         `gorm:"column:source;not null" json:"source"`
	PublishedDate       time.Time      `gorm:"column:published_date;not null" json:"published_date"`
	Content             string         `gorm:"column:content;not null" json:"content"`
	Summary             string         `gorm:"column:summary" json:"summary,omitempty"`
	Tags                pq.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	ImageURL            string         `gorm:"column:image_url" json:"image_url,omitempty"`
	Sentiment           string         `gorm:"column:sentiment" json:"sentiment,omitempty"`
	SentimentTickers    pq.StringArray `gorm:"column:sentiment_tickers;type:text[]" json:"sentiment_tickers,omitempty"`
	SentimentAnalyzedAt *time.Time     `gorm:"column:sentiment_analyzed_at" json:"sentiment_analyzed_at,omitempty"`
	CreatedAt           time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Article) TableName() string {
	return "articles"
}
