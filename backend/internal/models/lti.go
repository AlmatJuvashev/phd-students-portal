package models

import (
	"time"
)

// LTITool represents an external tool registered via LTI 1.3
type LTITool struct {
	ID               string    `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	Name             string    `json:"name" db:"name"`
	ClientID         string    `json:"client_id" db:"client_id"`
	InitiateLoginURL string    `json:"initiate_login_url" db:"initiate_login_url"`
	RedirectionURIs  []string  `json:"redirection_uris" db:"redirection_uris"`
	PublicJWKSURL    *string   `json:"public_jwks_url" db:"public_jwks_url"`
	DeploymentID     string    `json:"deployment_id" db:"deployment_id"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// LTIKey represents a stored RSA keypair for signing JWTs
type LTIKey struct {
	ID         string     `json:"id" db:"id"`
	KID        string     `json:"kid" db:"kid"`
	PrivateKey string     `json:"-" db:"private_key"`
	PublicKey  string     `json:"public_key" db:"public_key"`
	Algorithm  string     `json:"alg" db:"algorithm"`
	Use        string     `json:"use" db:"use"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at" db:"expires_at"`
}

type JWK struct {
	KID string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}


// CreateToolParams for registering a new tool
type CreateToolParams struct {
	TenantID         string   `json:"tenant_id" binding:"required"`
	Name             string   `json:"name" binding:"required"`
	ClientID         string   `json:"client_id" binding:"required"`
	InitiateLoginURL string   `json:"initiate_login_url" binding:"required,url"`
	RedirectionURIs  []string `json:"redirection_uris" binding:"required"`
	PublicJWKSURL    *string  `json:"public_jwks_url" binding:"omitempty,url"`
}
