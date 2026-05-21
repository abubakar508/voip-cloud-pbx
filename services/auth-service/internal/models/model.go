package models

import "time"

type Tenant struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

type UserRole string

const (
	RoleAdmin       UserRole = "ADMIN"
	RoleTenantAdmin UserRole = "TENANT_ADMIN"
	RoleAgent       UserRole = "AGENT"
	RoleViewer      UserRole = "VIEWER"
)

type User struct {
	ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     string    `gorm:"type:uuid;not null;index"`
	Email        string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	DisplayName  string    `gorm:"type:varchar(255);not null"`
	Role         UserRole  `gorm:"type:varchar(32);not null;default:'AGENT'"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}
