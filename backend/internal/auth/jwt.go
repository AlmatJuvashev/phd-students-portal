package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a signed JWT with role claims and 6-month expiry (configurable by days).
func GenerateJWT(sub string, role string, secret []byte, expDays int) (string, error) {
	claims := jwt.MapClaims{
		"sub":  sub,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * time.Duration(expDays)).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

// GenerateJWTWithTenant creates a signed JWT with tenant context for multitenancy.
func GenerateJWTWithTenant(sub string, roles []string, activeRole string, tenantID string, isSuperadmin bool, secret []byte, expDays int) (string, error) {
	claims := jwt.MapClaims{
		"sub":           sub,
		"roles":         roles,
		"active_role":   activeRole,
		"tenant_id":     tenantID,
		"is_superadmin": isSuperadmin,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(time.Hour * 24 * time.Duration(expDays)).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

