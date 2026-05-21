package auth

import (
	"errors"
	"time"

	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/config"
	"github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
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

type refreshClaims struct {
	UserID   string          `json:"uid"`
	TenantID string          `json:"tid"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
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

func (s *Service) ParseAccessToken(tokenStr string) (*customClaims, error) {
	if tokenStr == "" {
		return nil, ErrInvalidCredentials
	}
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTAccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidCredentials
	}
	return claims, nil
}

func (s *Service) ParseRefreshToken(tokenStr string) (*refreshClaims, error) {
	if tokenStr == "" {
		return nil, ErrInvalidCredentials
	}
	token, err := jwt.ParseWithClaims(tokenStr, &refreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTRefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*refreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidCredentials
	}
	return claims, nil
}

func (s *Service) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := s.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	var user models.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(user, s.cfg.JWTAccessSecret, s.cfg.JWTAccessTTL)
}
