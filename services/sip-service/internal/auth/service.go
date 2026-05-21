package auth

import (
	"errors"

	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Authenticate checks username + password against sip_accounts.
func (s *Service) Authenticate(username, password string) (*models.SipAccount, error) {
	var acc models.SipAccount
	if err := s.db.Where("username = ?", username).First(&acc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &acc, nil
}

func (s *Service) FindAccountByUsername(username string) (*models.SipAccount, error) {
	var acc models.SipAccount
	if err := s.db.Where("username = ?", username).First(&acc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	return &acc, nil
}
