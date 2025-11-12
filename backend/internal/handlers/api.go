package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// BuildAPI wires routes and returns a *gin.Engine
func BuildAPI(r *gin.Engine, db *sqlx.DB, cfg config.AppConfig, playbookManager *pb.Manager) *gin.Engine {
	r.Use(middleware.RequestLogger())
	// CORS for frontend dev and configured origin
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return false
			}
			if origin == cfg.FrontendBase {
				return true
			}
			// allow any localhost/127.0.0.1 port for dev
			if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "https://localhost:") {
				return true
			}
			if strings.HasPrefix(origin, "http://127.0.0.1:") || strings.HasPrefix(origin, "https://127.0.0.1:") {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	api := r.Group("/api")

	// Debug endpoint to check CORS config (remove in production)
	api.GET("/debug/cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"frontend_base": cfg.FrontendBase,
			"origin":        c.Request.Header.Get("Origin"),
		})
	})

	api.GET("/me", func(c *gin.Context) {
		// require auth, then return current user
		middleware.AuthRequired([]byte(cfg.JWTSecret))(c)
		if c.IsAborted() {
			return
		}
		// Hydrate from DB/Redis
		rds := services.NewRedis(cfg.RedisURL)
		middleware.HydrateUserFromClaims(c, db, rds)
		if val, ok := c.Get("current_user"); ok {
			c.JSON(200, val)
			return
		}
		c.JSON(401, gin.H{"error": "unauthenticated"})
	})

	// Auth routes (login only - password reset done by admin)
	auth := NewAuthHandler(db, cfg)
	api.POST("/auth/login", auth.Login)

    users := NewUsersHandler(db, cfg)
    journey := NewJourneyHandler(db, cfg, playbookManager)
    nodeSubmission := NewNodeSubmissionHandler(db, cfg, playbookManager)
    adminHandler := NewAdminHandler(db, cfg, playbookManager)
	_ = NewMeHandler(db, cfg, services.NewRedis(cfg.RedisURL)) // TODO: use me handler for routes
	api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// Checklist + Documents + Comments
	check := NewChecklistHandler(db, cfg)
	api.GET("/checklist/modules", check.ListModules)
	api.GET("/checklist/steps", check.ListStepsByModule)
	api.GET("/students/:id/steps", check.ListStudentSteps)
	api.PATCH("/students/:id/steps/:stepId", check.UpdateStudentStep)

	docs := NewDocumentsHandler(db, cfg)
	api.POST("/students/:id/documents", docs.CreateDocument)   // create doc metadata
	api.GET("/students/:id/documents", docs.ListDocuments)     // list docs
	api.POST("/documents/:docId/versions", docs.UploadVersion) // multipart upload
	api.GET("/documents/:docId", docs.GetDocument)
	api.GET("/documents/:docId/presign-get", docs.PresignGetLatest)
	api.GET("/documents/versions/:versionId/download", docs.DownloadVersion)
	api.POST("/documents/:docId/presign", docs.PresignUpload)

	cmts := NewCommentsHandler(db, cfg)
	api.GET("/documents/:docId/comments", cmts.ListComments)
	api.POST("/documents/:docId/comments", cmts.AddComment)
	api.PATCH("/comments/:id", cmts.UpdateComment)

	// Advisor inbox (pending submissions)
	api.GET("/advisor/inbox", check.AdvisorInbox)
	api.POST("/reviews/:id/steps/:stepId/approve", check.ApproveStep)
	api.POST("/reviews/:id/steps/:stepId/return", check.ReturnStep)

	// Admin-only routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired([]byte(cfg.JWTSecret)))
	admin.Use(middleware.RequireRoles("admin", "superadmin"))
    admin.GET("/users", users.ListUsers)
    admin.POST("/users", users.CreateUser)
    admin.PUT("/users/:id", users.UpdateUser)
    admin.POST("/users/:id/reset-password", users.ResetPasswordForUser)
    admin.PATCH("/users/:id/active", users.SetActive)
    admin.GET("/student-progress", adminHandler.StudentProgress)
    admin.GET("/monitor/students", adminHandler.MonitorStudents)
    admin.GET("/students/:id/journey", adminHandler.StudentJourney)
    admin.PATCH("/students/:id/nodes/:nodeId/state", adminHandler.PatchStudentNodeState)

	// Self-service password change
	api.PATCH("/me/password", middleware.AuthRequired([]byte(cfg.JWTSecret)), users.ChangeOwnPassword)

	// Journey state (per-user)
	js := api.Group("/journey")
	js.Use(middleware.AuthRequired([]byte(cfg.JWTSecret)))
	js.GET("/state", journey.GetState)
	js.PUT("/state", journey.SetState)
	js.POST("/reset", journey.Reset)
	js.GET("/profile", nodeSubmission.GetProfile)

	nodes := js.Group("/nodes")
	nodes.GET("/:nodeId/submission", nodeSubmission.GetSubmission)
	nodes.PUT("/:nodeId/submission", nodeSubmission.PutSubmission)
	nodes.POST("/:nodeId/uploads/presign", nodeSubmission.PresignUpload)
	nodes.POST("/:nodeId/uploads/attach", nodeSubmission.AttachUpload)
	nodes.PATCH("/:nodeId/state", nodeSubmission.PatchState)

	// TODO: checklist, documents, comments handlers (skeletons for now)

	return r
}
