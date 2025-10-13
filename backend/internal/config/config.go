package config

import (
	"fmt"
	"log"
	"os"
)

// AppConfig holds env-driven configuration.
type AppConfig struct {
	RedisURL string

	Port          string
	Env           string
	JWTSecret     string
	JWTExpDays    int
	DatabaseURL   string
	UploadDir     string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string
	SMTPFrom      string
	FrontendBase  string
	AdminEmail    string
	AdminPassword string
	PlaybookPath  string
}

// MustLoad loads configuration from environment variables.
func MustLoad() AppConfig {
	cfg := AppConfig{
		RedisURL:      get("REDIS_URL", "redis://localhost:6379"),
		Port:          get("APP_PORT", "8080"),
		Env:           get("APP_ENV", "development"),
		JWTSecret:     get("JWT_SECRET", "change-me"),
		JWTExpDays:    atoi(get("JWT_EXP_DAYS", "180")),
		DatabaseURL:   get("DATABASE_URL", ""),
		UploadDir:     get("UPLOAD_DIR", "./uploads"),
		SMTPHost:      get("SMTP_HOST", "localhost"),
		SMTPPort:      get("SMTP_PORT", "1025"),
		SMTPUser:      get("SMTP_USER", ""),
		SMTPPass:      get("SMTP_PASS", ""),
		SMTPFrom:      get("SMTP_FROM", "PhD Portal <no-reply@local>"),
		FrontendBase:  get("FRONTEND_BASE", "http://localhost:5173"),
		AdminEmail:    get("ADMIN_EMAIL", ""),
		AdminPassword: get("ADMIN_PASSWORD", ""),
		PlaybookPath:  get("PLAYBOOK_PATH", "../frontend/src/playbooks/playbook.json"),
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	_ = os.MkdirAll(cfg.UploadDir, 0755)
	
	// Log important config values at startup
	log.Printf("Config loaded: Port=%s, Env=%s, FrontendBase=%s", cfg.Port, cfg.Env, cfg.FrontendBase)
	
	return cfg
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoi(s string) int {
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}
