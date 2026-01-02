package config

import (
	"fmt"
	"log"
	"os"
)

// AppConfig holds env-driven configuration.
type AppConfig struct {
	RedisURL string

	Port            string
	Env             string
	JWTSecret       string
	JWTExpDays      int
	DatabaseURL     string
	UploadDir       string
	FileUploadMaxMB int
	SMTPHost        string
	SMTPPort        string
	SMTPUser        string
	SMTPPass        string
	SMTPFrom        string
	FrontendBase    string
	AdminEmail      string
	AdminPassword   string
	PlaybookPath    string
	S3Endpoint      string
	S3Bucket        string
	ServerURL       string
	IssuerURL       string
	OpenAIKey       string
}

// MustLoad loads configuration from environment variables.
func MustLoad() AppConfig {
	serverURL := get("SERVER_URL", "http://localhost:8080")
	cfg := AppConfig{
		RedisURL:        get("REDIS_URL", "redis://127.0.0.1:6379"),
		Port:            get("APP_PORT", "8080"),
		Env:             get("APP_ENV", "development"),
		JWTSecret:       get("JWT_SECRET", "change-me"),
		JWTExpDays:      atoi(get("JWT_EXP_DAYS", "180")),
		DatabaseURL:     get("DATABASE_URL", ""),
		UploadDir:       get("UPLOAD_DIR", "./uploads"),
		FileUploadMaxMB: atoi(get("FILE_UPLOAD_MAX_MB", "25")),
		SMTPHost:        get("SMTP_HOST", "localhost"),
		SMTPPort:        get("SMTP_PORT", "1025"),
		SMTPUser:        get("SMTP_USER", ""),
		SMTPPass:        get("SMTP_PASS", ""),
		SMTPFrom:        get("SMTP_FROM", "PhD Portal <no-reply@local>"),
		FrontendBase:    get("FRONTEND_BASE", "http://localhost:5173"),
		AdminEmail:      get("SUPERADMIN_EMAIL", get("ADMIN_EMAIL", "")),
		AdminPassword:   get("SUPERADMIN_PASSWORD", get("ADMIN_PASSWORD", "")),
		PlaybookPath:    get("PLAYBOOK_PATH", "../frontend/src/playbooks/playbook.json"),
		S3Endpoint:      get("S3_ENDPOINT", ""),
		S3Bucket:        get("S3_BUCKET", ""),
		ServerURL:       serverURL,
		IssuerURL:       get("ISSUER_URL", serverURL), // Default to server URL
		OpenAIKey:       get("OPENAI_API_KEY", ""),
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	// Warning only for AI
	if cfg.OpenAIKey == "" {
		log.Println("WARN: OPENAI_API_KEY is missing. AI features will be disabled.")
	}

	_ = os.MkdirAll(cfg.UploadDir, 0750)

	// Log important config values at startup
	log.Printf("Config loaded: Port=%s, Env=%s, FrontendBase=%s, ServerURL=%s", cfg.Port, cfg.Env, cfg.FrontendBase, cfg.ServerURL)

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
