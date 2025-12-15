package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	secret := []byte("testsecret")
	sub := "12345"
	role := "admin"
	expDays := 1

	tokenString, err := GenerateJWT(sub, role, secret, expDays)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Parse and verify
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, sub, claims["sub"])
	assert.Equal(t, role, claims["role"])
	assert.NotNil(t, claims["iat"])
	assert.NotNil(t, claims["exp"])
}

func TestGenerateJWTWithTenant(t *testing.T) {
	secret := []byte("testsecret")
	sub := "67890"
	role := "student"
	tenantID := "tenant-uuid"
	isSuperadmin := false
	expDays := 7

	tokenString, err := GenerateJWTWithTenant(sub, role, tenantID, isSuperadmin, secret, expDays)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Parse and verify
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, sub, claims["sub"])
	assert.Equal(t, role, claims["role"])
	assert.Equal(t, tenantID, claims["tenant_id"])
	assert.Equal(t, isSuperadmin, claims["is_superadmin"])
}
