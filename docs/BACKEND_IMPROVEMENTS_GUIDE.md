# Backend Improvements Implementation Guide

> **Document Version:** 1.0  
> **Created:** January 3, 2026  
> **Status:** Implementation Ready

This guide provides detailed implementation steps for all recommended improvements to bring the Universal Education Portal backend to enterprise production standards.

---

## Table of Contents

1. [Critical Improvements](#critical-improvements)
   - [1.1 Structured Logging](#11-structured-logging)
   - [1.2 Database Connection Pooling](#12-database-connection-pooling)
   - [1.3 Health Check Endpoints](#13-health-check-endpoints)
   - [1.4 Secrets Management](#14-secrets-management)
2. [Preferrable Improvements](#preferrable-improvements)
   - [2.1 Request ID / Correlation ID](#21-request-id--correlation-id)
   - [2.2 Graceful Shutdown](#22-graceful-shutdown)
   - [2.3 OpenAPI Documentation](#23-openapi-documentation)
   - [2.4 API Versioning](#24-api-versioning)
3. [Implementation Timeline](#implementation-timeline)
4. [Testing Checklist](#testing-checklist)

---

## Critical Improvements

### 1.1 Structured Logging

**Priority:** 🔴 Critical  
**Estimated Time:** 4-6 hours  
**Impact:** Observability, Debugging, Production Monitoring

#### Why It Matters

- JSON-structured logs integrate with ELK Stack, Datadog, CloudWatch
- Enables log aggregation and searching across distributed systems
- Provides context (request_id, user_id, tenant_id) for debugging

#### Implementation Steps

**Step 1: Add zerolog dependency**

```bash
cd backend
go get github.com/rs/zerolog
```

**Step 2: Create new logging package**

Create file: `internal/logging/zerolog.go`

```go
package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// Config holds logging configuration
type Config struct {
	Level      string // debug, info, warn, error
	Pretty     bool   // Human-readable output (dev only)
	TimeFormat string
}

// Init initializes the global logger
func Init(cfg Config) {
	// Set global time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Output configuration
	var output io.Writer = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	}

	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = Logger
}

// WithContext creates a child logger with additional fields
func WithContext(fields map[string]interface{}) zerolog.Logger {
	ctx := Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx.Logger()
}

// Request logging helpers
func Info(msg string, fields ...map[string]interface{}) {
	event := Logger.Info()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Error(msg string, err error, fields ...map[string]interface{}) {
	event := Logger.Error().Err(err)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Warn(msg string, fields ...map[string]interface{}) {
	event := Logger.Warn()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Debug(msg string, fields ...map[string]interface{}) {
	event := Logger.Debug()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}
```

**Step 3: Create Gin logging middleware**

Create file: `internal/middleware/zerolog_middleware.go`

```go
package middleware

import (
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ZerologRequestLogger returns a gin middleware for structured request logging
func ZerologRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request context values
		requestID, _ := c.Get("request_id")
		tenantID := c.GetString("tenant_id")
		userID := c.GetString("userID")

		// Build log event
		var event *zerolog.Event
		status := c.Writer.Status()

		switch {
		case status >= 500:
			event = logging.Logger.Error()
		case status >= 400:
			event = logging.Logger.Warn()
		default:
			event = logging.Logger.Info()
		}

		if raw != "" {
			path = path + "?" + raw
		}

		event.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size())

		// Add context fields if present
		if requestID != nil {
			event.Str("request_id", requestID.(string))
		}
		if tenantID != "" {
			event.Str("tenant_id", tenantID)
		}
		if userID != "" {
			event.Str("user_id", userID)
		}

		// Add errors if any
		if len(c.Errors) > 0 {
			event.Strs("errors", c.Errors.Errors())
		}

		event.Msg("HTTP Request")
	}
}
```

**Step 4: Update config to include logging settings**

Update `internal/config/config.go`:

```go
type AppConfig struct {
	// ... existing fields ...

	// Logging
	LogLevel  string
	LogPretty bool
}

func MustLoad() AppConfig {
	cfg := AppConfig{
		// ... existing fields ...
		LogLevel:  get("LOG_LEVEL", "info"),
		LogPretty: get("LOG_PRETTY", "false") == "true",
	}
	// ...
}
```

**Step 5: Initialize in main.go**

Update `cmd/server/main.go`:

```go
import (
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	// Initialize structured logging FIRST
	logging.Init(logging.Config{
		Level:  cfg.LogLevel,
		Pretty: cfg.LogPretty || cfg.Env == "development",
	})

	logging.Info("Starting server", map[string]interface{}{
		"port": cfg.Port,
		"env":  cfg.Env,
	})

	// ... rest of initialization ...
}
```

**Step 6: Update .env.example**

```env
# Logging
LOG_LEVEL=info          # debug, info, warn, error
LOG_PRETTY=false        # true for human-readable (dev only)
```

#### Migration Strategy

1. Add new logging package alongside existing
2. Gradually replace `log.Printf` calls with `logging.Info/Error/Warn`
3. Remove old `internal/logging/logging.go` after migration complete

---

### 1.2 Database Connection Pooling

**Priority:** 🔴 Critical  
**Estimated Time:** 1-2 hours  
**Impact:** Performance, Reliability under load

#### Why It Matters

- Prevents connection exhaustion under high load
- Reduces latency by reusing connections
- Controls resource usage on database server

#### Implementation Steps

**Step 1: Update db/db.go**

Replace the current implementation:

```go
package db

import (
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PoolConfig holds connection pool settings
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultPoolConfig returns sensible defaults for most workloads
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,              // Max concurrent connections
		MaxIdleConns:    5,               // Keep 5 idle connections ready
		ConnMaxLifetime: 5 * time.Minute, // Recreate connections after 5 min
		ConnMaxIdleTime: 1 * time.Minute, // Close idle connections after 1 min
	}
}

// MustOpen connects to Postgres with connection pooling or exits on error.
func MustOpen(url string) *sqlx.DB {
	return MustOpenWithConfig(url, DefaultPoolConfig())
}

// MustOpenWithConfig connects with custom pool settings
func MustOpenWithConfig(url string, poolCfg PoolConfig) *sqlx.DB {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		logging.Error("Failed to connect to database", err, nil)
		panic(err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(poolCfg.MaxOpenConns)
	db.SetMaxIdleConns(poolCfg.MaxIdleConns)
	db.SetConnMaxLifetime(poolCfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(poolCfg.ConnMaxIdleTime)

	// Verify connection
	if err := db.Ping(); err != nil {
		logging.Error("Failed to ping database", err, nil)
		panic(err)
	}

	logging.Info("Database connected", map[string]interface{}{
		"max_open_conns":     poolCfg.MaxOpenConns,
		"max_idle_conns":     poolCfg.MaxIdleConns,
		"conn_max_lifetime":  poolCfg.ConnMaxLifetime.String(),
		"conn_max_idle_time": poolCfg.ConnMaxIdleTime.String(),
	})

	return db
}

// Stats returns current pool statistics
func Stats(db *sqlx.DB) map[string]interface{} {
	stats := db.Stats()
	return map[string]interface{}{
		"open_connections":   stats.OpenConnections,
		"in_use":             stats.InUse,
		"idle":               stats.Idle,
		"wait_count":         stats.WaitCount,
		"wait_duration":      stats.WaitDuration.String(),
		"max_idle_closed":    stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}
```

**Step 2: Add environment variables**

Add to `config/config.go`:

```go
type AppConfig struct {
	// ... existing ...

	// Database Pool
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime int // seconds
	DBConnMaxIdleTime int // seconds
}

func MustLoad() AppConfig {
	cfg := AppConfig{
		// ... existing ...
		DBMaxOpenConns:    atoi(get("DB_MAX_OPEN_CONNS", "25")),
		DBMaxIdleConns:    atoi(get("DB_MAX_IDLE_CONNS", "5")),
		DBConnMaxLifetime: atoi(get("DB_CONN_MAX_LIFETIME", "300")),
		DBConnMaxIdleTime: atoi(get("DB_CONN_MAX_IDLE_TIME", "60")),
	}
}
```

**Step 3: Update main.go**

```go
poolCfg := db.PoolConfig{
	MaxOpenConns:    cfg.DBMaxOpenConns,
	MaxIdleConns:    cfg.DBMaxIdleConns,
	ConnMaxLifetime: time.Duration(cfg.DBConnMaxLifetime) * time.Second,
	ConnMaxIdleTime: time.Duration(cfg.DBConnMaxIdleTime) * time.Second,
}
conn := db.MustOpenWithConfig(cfg.DatabaseURL, poolCfg)
```

**Step 4: Update .env.example**

```env
# Database Connection Pool
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300  # seconds
DB_CONN_MAX_IDLE_TIME=60  # seconds
```

#### Tuning Guidelines

| Workload               | MaxOpen | MaxIdle | Lifetime | IdleTime |
| ---------------------- | ------- | ------- | -------- | -------- |
| Low (< 100 req/s)      | 10      | 5       | 5 min    | 1 min    |
| Medium (100-500 req/s) | 25      | 10      | 5 min    | 1 min    |
| High (> 500 req/s)     | 50-100  | 20      | 5 min    | 30 sec   |

---

### 1.3 Health Check Endpoints

**Priority:** 🔴 Critical  
**Estimated Time:** 2-3 hours  
**Impact:** Kubernetes readiness, Load balancer health, Monitoring

#### Why It Matters

- Kubernetes uses probes to determine pod readiness
- Load balancers need health endpoints to route traffic
- Enables zero-downtime deployments

#### Implementation Steps

**Step 1: Create health handler**

Create file: `internal/handlers/health.go`

```go
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewHealthHandler(db *sqlx.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redis}
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Version   string                 `json:"version,omitempty"`
	Checks    map[string]CheckResult `json:"checks,omitempty"`
}

type CheckResult struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// Liveness - simple check that the service is running
// Used by: Kubernetes livenessProbe
// Endpoint: GET /health/live
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness - checks all dependencies are available
// Used by: Kubernetes readinessProbe, Load Balancers
// Endpoint: GET /health/ready
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]CheckResult)
	allHealthy := true

	// Check Database
	dbCheck := h.checkDatabase(ctx)
	checks["database"] = dbCheck
	if dbCheck.Status != "ok" {
		allHealthy = false
	}

	// Check Redis (if configured)
	if h.redis != nil {
		redisCheck := h.checkRedis(ctx)
		checks["redis"] = redisCheck
		if redisCheck.Status != "ok" {
			allHealthy = false
		}
	}

	status := HealthStatus{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	if allHealthy {
		status.Status = "ok"
		c.JSON(http.StatusOK, status)
	} else {
		status.Status = "degraded"
		c.JSON(http.StatusServiceUnavailable, status)
	}
}

// Detailed health with metrics (for monitoring dashboards)
// Endpoint: GET /health/detailed
func (h *HealthHandler) Detailed(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]CheckResult)

	// Database with pool stats
	dbCheck := h.checkDatabase(ctx)
	if dbCheck.Status == "ok" {
		stats := h.db.Stats()
		dbCheck.Message = ""
		checks["database"] = dbCheck
		checks["database_pool"] = CheckResult{
			Status: "ok",
			Message: formatPoolStats(stats),
		}
	} else {
		checks["database"] = dbCheck
	}

	// Redis
	if h.redis != nil {
		checks["redis"] = h.checkRedis(ctx)
	}

	allHealthy := true
	for _, check := range checks {
		if check.Status != "ok" {
			allHealthy = false
			break
		}
	}

	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}
	if !allHealthy {
		status.Status = "degraded"
	}

	c.JSON(http.StatusOK, status)
}

func (h *HealthHandler) checkDatabase(ctx context.Context) CheckResult {
	start := time.Now()
	err := h.db.PingContext(ctx)
	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			Status:   "error",
			Message:  err.Error(),
			Duration: duration.String(),
		}
	}
	return CheckResult{
		Status:   "ok",
		Duration: duration.String(),
	}
}

func (h *HealthHandler) checkRedis(ctx context.Context) CheckResult {
	start := time.Now()
	_, err := h.redis.Ping(ctx).Result()
	duration := time.Since(start)

	if err != nil {
		return CheckResult{
			Status:   "error",
			Message:  err.Error(),
			Duration: duration.String(),
		}
	}
	return CheckResult{
		Status:   "ok",
		Duration: duration.String(),
	}
}

func formatPoolStats(stats sql.DBStats) string {
	return fmt.Sprintf(
		"open=%d, in_use=%d, idle=%d, wait_count=%d",
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
	)
}
```

**Step 2: Register routes in api.go**

Add to `BuildAPI` function:

```go
// Health endpoints (no auth required)
healthHandler := NewHealthHandler(db, rds)
r.GET("/health/live", healthHandler.Liveness)
r.GET("/health/ready", healthHandler.Readiness)
r.GET("/health/detailed", healthHandler.Detailed)

// Legacy health endpoint
r.GET("/health", healthHandler.Readiness)
```

**Step 3: Create Kubernetes probe configuration**

Create file: `deploy/k8s/probes.yaml` (reference):

```yaml
# Example Kubernetes deployment probes
livenessProbe:
  httpGet:
    path: /health/live
    port: 8280
  initialDelaySeconds: 10
  periodSeconds: 15
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8280
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3

startupProbe:
  httpGet:
    path: /health/live
    port: 8280
  initialDelaySeconds: 0
  periodSeconds: 5
  timeoutSeconds: 5
  failureThreshold: 30 # Allow 150 seconds for startup
```

**Step 4: Update Dockerfile for health check**

Add to `Dockerfile`:

```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8280/health/live || exit 1
```

---

### 1.4 Secrets Management

**Priority:** 🔴 Critical (for production)  
**Estimated Time:** 4-8 hours (depends on infrastructure)  
**Impact:** Security, Compliance

#### Why It Matters

- Prevents credential leaks in version control
- Enables secret rotation without redeployment
- Required for SOC2/HIPAA/GDPR compliance

#### Implementation Options

##### Option A: Environment Variables from Orchestrator (Recommended for small teams)

**For Docker Compose:**

```yaml
# docker-compose.yml
services:
  backend:
    env_file:
      - .env.production # Not committed to git
    environment:
      - DATABASE_URL # Passed from host
```

**For Kubernetes:**

```yaml
# secrets.yaml (apply separately, never commit)
apiVersion: v1
kind: Secret
metadata:
  name: phd-portal-secrets
type: Opaque
data:
  DATABASE_URL: base64_encoded_value
  JWT_SECRET: base64_encoded_value
  S3_ACCESS_KEY: base64_encoded_value
  S3_SECRET_KEY: base64_encoded_value
```

```yaml
# deployment.yaml
spec:
  containers:
    - name: backend
      envFrom:
        - secretRef:
            name: phd-portal-secrets
```

##### Option B: HashiCorp Vault (Recommended for enterprise)

**Step 1: Add Vault client**

```bash
go get github.com/hashicorp/vault/api
```

**Step 2: Create secrets loader**

Create file: `internal/config/vault.go`

```go
package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Address   string
	Token     string
	Path      string
	MountPath string
}

func LoadSecretsFromVault(cfg VaultConfig) (map[string]string, error) {
	config := api.DefaultConfig()
	config.Address = cfg.Address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("vault client creation failed: %w", err)
	}

	client.SetToken(cfg.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	secret, err := client.KVv2(cfg.MountPath).Get(ctx, cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("vault secret read failed: %w", err)
	}

	secrets := make(map[string]string)
	for k, v := range secret.Data {
		if str, ok := v.(string); ok {
			secrets[k] = str
		}
	}

	return secrets, nil
}

// SetEnvFromVault loads secrets and sets as environment variables
func SetEnvFromVault() error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	vaultPath := os.Getenv("VAULT_SECRET_PATH")

	if vaultAddr == "" {
		return nil // Vault not configured, use env vars
	}

	secrets, err := LoadSecretsFromVault(VaultConfig{
		Address:   vaultAddr,
		Token:     vaultToken,
		Path:      vaultPath,
		MountPath: "secret",
	})
	if err != nil {
		return err
	}

	for k, v := range secrets {
		os.Setenv(k, v)
	}

	return nil
}
```

**Step 3: Initialize in main.go**

```go
func main() {
	// Load secrets from Vault (if configured)
	if err := config.SetEnvFromVault(); err != nil {
		log.Printf("Vault secret loading failed: %v", err)
	}

	_ = godotenv.Load() // Fallback to .env
	cfg := config.MustLoad()
	// ...
}
```

##### Option C: AWS Secrets Manager

```go
package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func LoadSecretsFromAWS(secretName string) (map[string]string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, err
	}

	var secrets map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}
```

#### .gitignore Update

Ensure these are in `.gitignore`:

```gitignore
# Secrets - NEVER COMMIT
.env
.env.production
.env.local
*.pem
*.key
secrets/
```

---

## Preferrable Improvements

### 2.1 Request ID / Correlation ID

**Priority:** 🟡 Medium  
**Estimated Time:** 1-2 hours  
**Impact:** Debugging, Distributed Tracing

#### Implementation

Create file: `internal/middleware/request_id.go`

```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// RequestID middleware generates or propagates request IDs
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for existing request ID (from load balancer/API gateway)
		requestID := c.GetHeader(RequestIDHeader)

		// Generate new if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context
		c.Set(RequestIDKey, requestID)

		// Add to response headers
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves request ID from context
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		return id.(string)
	}
	return ""
}
```

**Register in api.go:**

```go
func BuildAPI(r *gin.Engine, db *sqlx.DB, cfg config.AppConfig, ...) *gin.Engine {
	// Request ID should be first middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.ZerologRequestLogger())  // Logging uses request ID
	// ... other middleware
}
```

---

### 2.2 Graceful Shutdown

**Priority:** 🟡 Medium  
**Estimated Time:** 2-3 hours  
**Impact:** Zero-downtime deployments, Data integrity

#### Implementation

Update `cmd/server/main.go`:

```go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"
	// ... other imports
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	// Initialize logging
	logging.Init(logging.Config{
		Level:  cfg.LogLevel,
		Pretty: cfg.Env == "development",
	})

	conn := db.MustOpen(cfg.DatabaseURL)
	defer conn.Close()

	// ... initialization code ...

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      api,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logging.Info("Server starting", map[string]interface{}{
			"port": cfg.Port,
			"env":  cfg.Env,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Server failed to start", err, nil)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logging.Info("Shutdown signal received, initiating graceful shutdown...", nil)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop accepting new requests, wait for existing to complete
	if err := srv.Shutdown(ctx); err != nil {
		logging.Error("Server forced to shutdown", err, nil)
	}

	// Stop background workers
	workerCancel() // Cancel context passed to workers

	// Close database connections
	if err := conn.Close(); err != nil {
		logging.Error("Database connection close error", err, nil)
	}

	logging.Info("Server gracefully stopped", nil)
}
```

---

### 2.3 OpenAPI Documentation

**Priority:** 🟢 Low  
**Estimated Time:** 8-16 hours (initial setup + annotations)  
**Impact:** Developer experience, API discoverability

#### Implementation with swaggo/swag

**Step 1: Install swag**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

**Step 2: Add annotations to handlers**

Example for `handlers/auth.go`:

```go
// Login godoc
// @Summary      User login
// @Description  Authenticate user with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Tenant-Slug  header    string    false  "Tenant identifier"
// @Param        body           body      loginReq  true   "Login credentials"
// @Success      200  {object}  map[string]interface{}  "Login successful"
// @Failure      400  {object}  map[string]string       "Invalid request"
// @Failure      401  {object}  map[string]string       "Invalid credentials"
// @Failure      429  {object}  map[string]string       "Too many attempts"
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// ... existing code
}
```

**Step 3: Add main annotations**

Add to `cmd/server/main.go`:

```go
// @title           PhD Student Portal API
// @version         1.0
// @description     Universal Education Portal Backend API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@phd-portal.kz

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8280
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// ...
}
```

**Step 4: Generate and serve docs**

```bash
swag init -g cmd/server/main.go -o docs/swagger
```

Add to `api.go`:

```go
import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/AlmatJuvashev/phd-students-portal/backend/docs/swagger"
)

func BuildAPI(...) *gin.Engine {
	// ...

	// Swagger docs (disable in production if needed)
	if cfg.Env != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
```

**Step 5: Add to Makefile**

```makefile
swagger:
	swag init -g cmd/server/main.go -o docs/swagger
	@echo "Swagger docs generated. Access at /swagger/index.html"
```

---

### 2.4 API Versioning

**Priority:** 🟢 Low  
**Estimated Time:** 2-4 hours  
**Impact:** Backward compatibility, API evolution

#### Implementation Strategy

**Option A: URL Path Versioning (Recommended)**

```go
func BuildAPI(r *gin.Engine, ...) *gin.Engine {
	// Version 1 API
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", auth.Login)
		v1.POST("/auth/logout", auth.Logout)
		// ... all routes
	}

	// Legacy support (redirect to v1)
	r.GET("/api/*path", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/v1"+c.Param("path"))
	})
}
```

**Option B: Header-based Versioning**

```go
func APIVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetHeader("API-Version")
		if version == "" {
			version = "1" // Default
		}
		c.Set("api_version", version)
		c.Next()
	}
}
```

---

## Implementation Timeline

### Phase 1: Critical (Week 1-2)

| Task                     | Days | Owner   |
| ------------------------ | ---- | ------- |
| Structured Logging       | 2    | Backend |
| Database Connection Pool | 1    | Backend |
| Health Checks            | 1    | Backend |
| Secrets Management Setup | 2    | DevOps  |
| Testing & Validation     | 2    | QA      |

### Phase 2: Preferrable (Week 3-4)

| Task                  | Days | Owner   |
| --------------------- | ---- | ------- |
| Request ID Middleware | 1    | Backend |
| Graceful Shutdown     | 1    | Backend |
| API Versioning        | 1    | Backend |
| OpenAPI Documentation | 3-5  | Backend |
| Documentation Review  | 2    | Team    |

---

## Testing Checklist

### Structured Logging

- [ ] Logs appear in JSON format in production mode
- [ ] Logs appear in human-readable format in dev mode
- [ ] Request ID appears in all request logs
- [ ] Tenant ID appears in authenticated request logs
- [ ] Error logs include stack traces

### Database Connection Pool

- [ ] Application starts with default pool settings
- [ ] Pool stats visible in `/health/detailed`
- [ ] Load test shows no connection exhaustion
- [ ] Idle connections are cleaned up

### Health Checks

- [ ] `/health/live` returns 200 immediately after start
- [ ] `/health/ready` returns 503 when DB is down
- [ ] `/health/detailed` shows pool statistics
- [ ] Kubernetes probes work correctly

### Graceful Shutdown

- [ ] In-flight requests complete before shutdown
- [ ] New requests are rejected during shutdown
- [ ] Database connections are properly closed
- [ ] Background workers stop gracefully

### Request ID

- [ ] Request ID generated when not provided
- [ ] Request ID propagated from header
- [ ] Request ID appears in response header
- [ ] Request ID visible in logs

---

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [12-Factor App](https://12factor.net/)
- [Kubernetes Probes Best Practices](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [PostgreSQL Connection Pooling](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [HashiCorp Vault Go Client](https://github.com/hashicorp/vault/tree/main/api)
- [zerolog Documentation](https://github.com/rs/zerolog)
- [swaggo/swag](https://github.com/swaggo/swag)
