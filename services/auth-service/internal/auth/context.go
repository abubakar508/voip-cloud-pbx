package auth

import "github.com/abubakar508/voip-cloud-pbx/services/auth-service/internal/models"

type ContextUser struct {
	ID       string
	TenantID string
	Role     models.UserRole
	Email    string
}
