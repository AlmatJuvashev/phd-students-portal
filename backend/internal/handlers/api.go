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
	chatHandler := NewChatHandler(db, cfg)
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
	api.GET("/documents/:docId/comments", cmts.GetComments)
	api.POST("/documents/:docId/comments", cmts.CreateComment)

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

	// Notifications endpoints
	notifications := NewNotificationsHandler(db)
	admin.GET("/notifications", notifications.ListNotifications)
	admin.GET("/notifications/unread-count", notifications.GetUnreadCount)
	admin.PATCH("/notifications/:id/read", notifications.MarkAsRead)
	admin.POST("/notifications/read-all", notifications.MarkAllAsRead)

	// Monitor endpoints (admin/advisor access) - extend admin group
	admin.GET("/monitor/students", adminHandler.MonitorStudents)
	admin.GET("/students", adminHandler.StudentProgress)
	admin.GET("/students/:id", adminHandler.GetStudentDetails)
	admin.GET("/students/:id/journey", adminHandler.StudentJourney)
	admin.GET("/students/:id/nodes/:nodeId/files", adminHandler.ListStudentNodeFiles)
	admin.PATCH("/students/:id/nodes/:nodeId/state", adminHandler.PatchStudentNodeState)
	admin.PATCH("/attachments/:attachmentId/review", adminHandler.ReviewAttachment)
	admin.POST("/attachments/:attachmentId/reviewed-document", adminHandler.UploadReviewedDocument)
	admin.POST("/attachments/:attachmentId/presign", adminHandler.PresignReviewedDocumentUpload)
	admin.POST("/attachments/:attachmentId/attach-reviewed", adminHandler.AttachReviewedDocument)
	admin.POST("/reminders", adminHandler.PostReminders)

	// Self-service password change and profile update
	api.PATCH("/me", middleware.AuthRequired([]byte(cfg.JWTSecret)), users.UpdateMe)
	api.PATCH("/me/password", middleware.AuthRequired([]byte(cfg.JWTSecret)), users.ChangeOwnPassword)
	api.GET("/me/pending-email", middleware.AuthRequired([]byte(cfg.JWTSecret)), users.GetPendingEmailVerification)
	api.GET("/me/verify-email", users.VerifyEmailChange) // No auth required for better UX


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

	// Chat module (placeholder)
	chat := api.Group("/chat")
	chat.Use(middleware.AuthRequired([]byte(cfg.JWTSecret)))
	chatAdmin := chat.Group("")
	chatAdmin.Use(middleware.RequireRoles("admin", "superadmin"))
	chatAdmin.POST("/rooms", chatHandler.CreateRoom)
	chatAdmin.PATCH("/rooms/:roomId", chatHandler.UpdateRoom)
	chatAdmin.GET("/rooms/:roomId/members", chatHandler.ListMembers)
	chatAdmin.POST("/rooms/:roomId/members", chatHandler.AddMember)
	chatAdmin.POST("/rooms/:roomId/members/batch", chatHandler.AddRoomMembersBatch)
	chatAdmin.DELETE("/rooms/:roomId/members/batch", chatHandler.RemoveRoomMembersBatch)
	chatAdmin.DELETE("/rooms/:roomId/members/:userId", chatHandler.RemoveMember)

	chat.GET("/rooms", chatHandler.ListRooms)
	chat.GET("/rooms/:roomId/messages", chatHandler.ListMessages)
	chat.POST("/rooms/:roomId/messages", chatHandler.CreateMessage)

	// TODO: checklist, documents, comments handlers (skeletons for now)

	// Dictionary endpoints (admin only)
	dictHandler := NewDictionaryHandler(db)
	dictGroup := admin.Group("/dictionaries")
	{
		dictGroup.GET("/programs", dictHandler.ListPrograms)
		dictGroup.POST("/programs", dictHandler.CreateProgram)
		dictGroup.PUT("/programs/:id", dictHandler.UpdateProgram)
		dictGroup.DELETE("/programs/:id", dictHandler.DeleteProgram)

		dictGroup.GET("/specialties", dictHandler.ListSpecialties)
		dictGroup.POST("/specialties", dictHandler.CreateSpecialty)
		dictGroup.PUT("/specialties/:id", dictHandler.UpdateSpecialty)
		dictGroup.DELETE("/specialties/:id", dictHandler.DeleteSpecialty)

		dictGroup.GET("/cohorts", dictHandler.ListCohorts)
		dictGroup.POST("/cohorts", dictHandler.CreateCohort)
		dictGroup.PUT("/cohorts/:id", dictHandler.UpdateCohort)
		dictGroup.DELETE("/cohorts/:id", dictHandler.DeleteCohort)

		dictGroup.GET("/departments", dictHandler.ListDepartments)
		dictGroup.POST("/departments", dictHandler.CreateDepartment)
		dictGroup.PUT("/departments/:id", dictHandler.UpdateDepartment)
		dictGroup.DELETE("/departments/:id", dictHandler.DeleteDepartment)
	}

	return r
}
