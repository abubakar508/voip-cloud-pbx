package models

import "time"

type AISummary struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CallID    string    `gorm:"type:uuid;not null;index"`
	TenantID  string    `gorm:"type:uuid;not null;index"`
	Summary   string    `gorm:"type:text;not null"`
	Sentiment string    `gorm:"type:varchar(32);not null"` // positive / neutral / negative
	Keywords  string    `gorm:"type:text;not null"`        // comma-separated keywords
	CreatedAt time.Time `gorm:"not null;default:now()"`
}
