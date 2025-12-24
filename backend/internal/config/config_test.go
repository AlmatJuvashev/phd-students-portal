package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustLoad_Defaults(t *testing.T) {
	// Save all config-related environment variables
	envVars := []string{
		"DATABASE_URL", "APP_PORT", "APP_ENV", "JWT_SECRET", "JWT_EXP_DAYS",
		"UPLOAD_DIR", "FILE_UPLOAD_MAX_MB", "SMTP_HOST", "SMTP_PORT",
		"SMTP_USER", "SMTP_PASS", "SMTP_FROM", "FRONTEND_BASE", "SERVER_URL",
		"REDIS_URL", "S3_ENDPOINT", "S3_REGION", "S3_BUCKET", "S3_ACCESS_KEY", "S3_SECRET_KEY",
	}
	
	originalValues := make(map[string]string)
	for _, key := range envVars {
		originalValues[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	
	defer func() {
		for key, val := range originalValues {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	// Set only the required DATABASE_URL
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")

	cfg := MustLoad()

	// Check defaults are applied
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "change-me", cfg.JWTSecret)
	assert.Equal(t, 180, cfg.JWTExpDays)
	assert.Equal(t, "./uploads", cfg.UploadDir)
	assert.Equal(t, 25, cfg.FileUploadMaxMB)
	assert.Equal(t, "localhost", cfg.SMTPHost)
	assert.Equal(t, "1025", cfg.SMTPPort)
	assert.Equal(t, "http://localhost:5173", cfg.FrontendBase)
}

func TestMustLoad_OverrideWithEnv(t *testing.T) {
	// Save and restore environment
	originalVars := map[string]string{
		"DATABASE_URL": os.Getenv("DATABASE_URL"),
		"APP_PORT":     os.Getenv("APP_PORT"),
		"APP_ENV":      os.Getenv("APP_ENV"),
		"JWT_SECRET":   os.Getenv("JWT_SECRET"),
	}
	defer func() {
		for k, v := range originalVars {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// Set custom env vars
	os.Setenv("DATABASE_URL", "postgres://prod:prod@db.example.com:5432/prod")
	os.Setenv("APP_PORT", "3000")
	os.Setenv("APP_ENV", "production")
	os.Setenv("JWT_SECRET", "super-secret-key")

	cfg := MustLoad()

	assert.Equal(t, "postgres://prod:prod@db.example.com:5432/prod", cfg.DatabaseURL)
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, "super-secret-key", cfg.JWTSecret)
}

func TestAtoi(t *testing.T) {
	// Test the atoi helper via JWTExpDays
	originalVars := map[string]string{
		"DATABASE_URL": os.Getenv("DATABASE_URL"),
		"JWT_EXP_DAYS": os.Getenv("JWT_EXP_DAYS"),
	}
	defer func() {
		for k, v := range originalVars {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("JWT_EXP_DAYS", "30")

	cfg := MustLoad()
	assert.Equal(t, 30, cfg.JWTExpDays)
}
