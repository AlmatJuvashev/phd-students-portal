package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/url"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/google/uuid"
)

type LTIService struct {
	repo repository.LTIRepository
	cfg  config.AppConfig
}

func NewLTIService(repo repository.LTIRepository, cfg config.AppConfig) *LTIService {
	return &LTIService{repo: repo, cfg: cfg}
}

// EnsureActiveKey checks if there is an active key, if not generates one.
func (s *LTIService) EnsureActiveKey(ctx context.Context) (*models.LTIKey, error) {
	key, err := s.repo.GetActiveKey(ctx)
	if err == nil && key != nil {
		return key, nil
	}
	// Generate new key
	return s.RotateKey(ctx)
}

func (s *LTIService) RotateKey(ctx context.Context) (*models.LTIKey, error) {
	// 1. Generate RSA 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	// 2. Encode to PEM
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	// 3. Store
	kid := uuid.New().String()
	ltiKey := models.LTIKey{
		KID:        kid,
		PrivateKey: string(privPEM),
		PublicKey:  string(pubPEM),
		Algorithm:  "RS256",
		Use:        "sig",
	}

	if err := s.repo.CreateKey(ctx, ltiKey); err != nil {
		return nil, err
	}
	
	// Re-fetch to get ID/dates
	return s.repo.GetActiveKey(ctx)
}

func (s *LTIService) GetJWKS(ctx context.Context) (*models.JWKS, error) {
	keys, err := s.repo.ListActiveKeys(ctx)
	if err != nil {
		return nil, err
	}

	jwks := &models.JWKS{Keys: []models.JWK{}}
	for _, k := range keys {
		block, _ := pem.Decode([]byte(k.PublicKey))
		if block == nil {
			continue
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			continue
		}
		rsaPub, ok := pubInterface.(*rsa.PublicKey)
		if !ok {
			continue
		}

		jwk := models.JWK{
			KID: k.KID,
			Kty: "RSA",
			Alg: "RS256",
			Use: "sig",
			N:   base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaPub.E)).Bytes()),
		}
		jwks.Keys = append(jwks.Keys, jwk)
	}
	return jwks, nil
}


func (s *LTIService) RegisterTool(ctx context.Context, p models.CreateToolParams) (*models.LTITool, error) {
	return s.repo.CreateTool(ctx, p)
}

func (s *LTIService) ListTools(ctx context.Context, tenantID string) ([]models.LTITool, error) {
	return s.repo.ListTools(ctx, tenantID)
}

// GenerateLoginInit constructs the OIDC Login Initiation URL to redirect the user TO the tool.
// This is Step 1 of the LTI 1.3 flow.
func (s *LTIService) GenerateLoginInit(ctx context.Context, toolID, userID string, targetLinkURI string) (string, error) {
	tool, err := s.repo.GetTool(ctx, toolID)
	if err != nil {
		return "", err
	}
	if tool == nil {
		return "", fmt.Errorf("tool not found")
	}

	// 1. Base URL
	u, err := url.Parse(tool.InitiateLoginURL)
	if err != nil {
		return "", fmt.Errorf("invalid tool login url: %w", err)
	}

	// 2. Query Params (OIDC Core)
	q := u.Query()
	q.Set("iss", s.cfg.IssuerURL)    // Our Platform Issuer
	q.Set("target_link_uri", targetLinkURI)
	q.Set("login_hint", userID)      // Current User ID
	q.Set("lti_message_hint", tool.DeploymentID) // Used to validate context on return
	q.Set("client_id", tool.ClientID) // The ID expected by the tool

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// ValidateLaunch handles the POST request from the tool (Step 2)
// This requires verifying the id_token signature against the Tool's public JWKS.
func (s *LTIService) ValidateLaunch(ctx context.Context, idToken string) (*models.LTITool, error) {
	// TODO: Phase 19.3 - Implement full JWT validation using go-jose or jwt-go and fetch JWKS
	// For now, returning nil to indicate logic is needed
	return nil, fmt.Errorf("launch validation not implemented")
}
