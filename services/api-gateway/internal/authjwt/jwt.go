package authjwt

import (
	"errors"
	"strings"
	"time"

	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"uid"`
	TenantID string `json:"tid"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var ErrInvalidToken = errors.New("invalid token")

type Validator struct {
	secret []byte
}

func NewValidator(cfg *config.AppConfig) *Validator {
	return &Validator{
		secret: []byte(cfg.JWTAccessSecret),
	}
}

func (v *Validator) ParseFromHeader(authHeader string) (*Claims, error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrInvalidToken
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return v.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
