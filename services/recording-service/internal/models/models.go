package models

import "time"

type Recording struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CallID      string    `gorm:"type:uuid;not null;index"`
	TenantID    string    `gorm:"type:uuid;not null;index"`
	FilePath    string    `gorm:"type:text;not null"`
	DurationSec int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
}

type CallRecord struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CallID      string    `gorm:"type:uuid;not null;uniqueIndex"`
	TenantID    string    `gorm:"type:uuid;not null;index"`
	FromUser    string    `gorm:"type:varchar(128);not null"`
	ToUser      string    `gorm:"type:varchar(128);not null"`
	StartedAt   time.Time `gorm:"not null"`
	EndedAt     time.Time `gorm:"not null"`
	DurationSec int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
}
