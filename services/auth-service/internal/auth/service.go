package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/config"
	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	db  *gorm.DB
	cfg *config.AppConfig
}

func NewService(db *gorm.DB, cfg *config.AppConfig) *Service {
	return &Service{db: db, cfg: cfg}
}

type customClaims struct {
	UserID   string          `json:"uid"`
	TenantID string          `json:"tid"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) RegisterTenantAdmin(req RegisterRequest) (*models.User, error) {
	// Check if user already exists
	var existing models.User
	if err := s.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("user already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tenant := &models.Tenant{
		Name: req.TenantName,
	}

	if err := s.db.Create(tenant).Error; err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		TenantID:     tenant.ID,
		Email:        req.Email,
		PasswordHash: string(hash),
		DisplayName:  req.DisplayName,
		Role:         models.RoleTenantAdmin,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(req LoginRequest) (*models.User, string, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	accessToken, err := s.generateToken(user, s.cfg.JWTAccessSecret, s.cfg.JWTAccessTTL)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := s.generateToken(user, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshTTL)
	if err != nil {
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

func (s *Service) generateToken(user models.User, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := customClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
