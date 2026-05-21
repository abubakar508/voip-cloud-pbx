package models

import "time"

type SipAccount struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     string    `gorm:"type:uuid;not null;index"`
	Extension    string    `gorm:"type:varchar(64);not null;index"`
	Username     string    `gorm:"type:varchar(128);not null;uniqueIndex"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	DisplayName  string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}
