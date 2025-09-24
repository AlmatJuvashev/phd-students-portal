package handlers

import (
	"phd-portal/backend/internal/services"
	"phd-portal/backend/internal/logging"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"phd-portal/backend/internal/config"
)

// BuildAPI wires routes and returns a *gin.Engine
func BuildAPI(r *gin.Engine, db *sqlx.DB, cfg config.AppConfig) *gin.Engine {
	r.Use(middleware.RequestLogger())
	api := r.Group("/api")

	api.GET("/me", func(c *gin.Context){
		// require auth, then return current user
		middleware.AuthRequired([]byte(cfg.JWTSecret))(c)
		if c.IsAborted() { return }
		// Hydrate from DB/Redis
		rds := dbpkg.NewRedis()
		middleware.HydrateUserFromClaims(c, db, rds)
		if val, ok := c.Get("current_user"); ok {
			c.JSON(200, val); return
		}
		c.JSON(401, gin.H{"error":"unauthenticated"})
	})


	// Auth routes (login, forgot/reset, logout is client-side token removal)
	auth := NewAuthHandler(db, cfg)
	api.POST("/auth/login", auth.Login)
	api.GET("/me", middleware.AuthRequired([]byte(cfg.JWTSecret)), me.Me)
	api.POST("/auth/forgot", auth.ForgotPassword) // sends email with reset link
	api.POST("/auth/reset", auth.ResetPassword)   // reset with token

	users := NewUsersHandler(db, cfg)
	me := NewMeHandler(db, cfg, services.NewRedis(cfg.RedisURL))
	api.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"ok":true}) })

	// Checklist + Documents + Comments
	check := NewChecklistHandler(db, cfg)
	api.GET("/checklist/modules", check.ListModules)
	api.GET("/checklist/steps", check.ListStepsByModule)
	api.GET("/students/:id/steps", check.ListStudentSteps)
	api.PATCH("/students/:id/steps/:stepId", check.UpdateStudentStep)

	docs := NewDocumentsHandler(db, cfg)
	api.POST("/students/:id/documents", docs.CreateDocument) // create doc metadata
	api.GET("/students/:id/documents", docs.ListDocuments) // list docs
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
	admin.Use(middleware.RequireRoles("admin","superadmin"))
	admin.GET("/users", users.ListUsers)
	admin.POST("/users", users.CreateUser)
	admin.PATCH("/users/:id/password", users.ResetPasswordForUser)
	admin.PATCH("/users/:id/active", users.SetActive)

	// Self-service password change
	api.PATCH("/me/password", middleware.AuthRequired([]byte(cfg.JWTSecret)), users.ChangeOwnPassword)

	// TODO: checklist, documents, comments handlers (skeletons for now)


	return r
}
