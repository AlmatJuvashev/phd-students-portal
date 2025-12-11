package handlers

import (
	"log"
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
	// Clean FrontendBase in case env var has embedded quotes
	cleanedFrontendBase := strings.Trim(cfg.FrontendBase, "\"' \t")
	
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return false
			}
			// Log CORS check for debugging
			log.Printf("[CORS] Origin=%q, FrontendBase=%q, Cleaned=%q", origin, cfg.FrontendBase, cleanedFrontendBase)
			
			if origin == cleanedFrontendBase {
				return true
			}
			// allow any localhost/127.0.0.1 port for dev
			if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "https://localhost:") {
				return true
			}
			if strings.HasPrefix(origin, "http://127.0.0.1:") || strings.HasPrefix(origin, "https://127.0.0.1:") {
				return true
			}
			// Allow *.localhost for local multitenancy testing
			if strings.Contains(origin, ".localhost:") {
				return true
			}
			// Allow subdomain-based tenant URLs (e.g., kaznmu.phd-portal.kz)
			// Parse configured frontend base to extract the main domain
			if cleanedFrontendBase != "" {
				// Simple subdomain matching for production
				// e.g., if FrontendBase is "https://phd-portal.kz", allow "*.phd-portal.kz"
				mainDomain := strings.TrimPrefix(cleanedFrontendBase, "https://")
				mainDomain = strings.TrimPrefix(mainDomain, "http://")
				if strings.Contains(origin, "."+mainDomain) || strings.HasSuffix(origin, mainDomain) {
					return true
				}
			}
			// Allow all Vercel deployments (they use X-Tenant-Slug header for tenant resolution)
			if strings.HasSuffix(origin, ".vercel.app") {
				return true
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Tenant-Slug"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Add tenant middleware for all API routes
	// This resolves tenant from subdomain or X-Tenant-Slug header
	r.Use(middleware.TenantMiddleware(db))
	
	api := r.Group("/api")
	
	// Serve static files from uploads directory
	// This maps /uploads/... to the configured upload directory on disk
	r.Static("/uploads", cfg.UploadDir)

	// Shared Redis (for auth hydration)
	rds := services.NewRedis(cfg.RedisURL)

	// Debug endpoint to check CORS config (remove in production)
	api.GET("/debug/cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"frontend_base": cfg.FrontendBase,
			"origin":        c.Request.Header.Get("Origin"),
			"tenant_id":     middleware.GetTenantID(c),
			"tenant_slug":   middleware.GetTenantSlug(c),
		})
	})

	api.GET("/me", func(c *gin.Context) {
		// require auth, then return current user
		middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds)(c)
		if c.IsAborted() {
			return
		}
		if val, ok := c.Get("current_user"); ok {
			c.JSON(200, val)
			return
		}
		c.JSON(401, gin.H{"error": "unauthenticated"})
	})

	// Services
	emailService := services.NewEmailService()

	// Auth routes (login and password reset)
	auth := NewAuthHandler(db, cfg, emailService)
	api.POST("/auth/login", auth.Login)
	api.POST("/auth/forgot-password", auth.ForgotPassword)
	api.POST("/auth/reset-password", auth.ResetPassword)

	users := NewUsersHandler(db, cfg)
	journey := NewJourneyHandler(db, cfg, playbookManager)
	nodeSubmission := NewNodeSubmissionHandler(db, cfg, playbookManager)
	adminHandler := NewAdminHandler(db, cfg, playbookManager)
	chatHandler := NewChatHandler(db, cfg, emailService)
	contactsHandler := NewContactsHandler(db)
	meHandler := NewMeHandler(db, cfg, services.NewRedis(cfg.RedisURL))
	api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	api.GET("/contacts", contactsHandler.PublicList)

	// /me routes (require auth)
	meGroup := api.Group("/me")
	meGroup.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	{
		meGroup.GET("/tenants", meHandler.MyTenants)
		meGroup.GET("/tenant", meHandler.MyTenant)
	}

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

	// Monitor endpoints (admin/advisor access)
	monitor := api.Group("/admin") // Keep URL prefix /admin for compatibility but change middleware
	monitor.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	monitor.Use(middleware.RequireRoles("admin", "superadmin", "advisor"))

	// Notifications (admin/advisor view of student events)
	notifications := NewNotificationsHandler(db)
	monitor.GET("/notifications", notifications.ListNotifications)
	monitor.GET("/notifications/unread-count", notifications.GetUnreadCount)
	monitor.PATCH("/notifications/:id/read", notifications.MarkAsRead)
	monitor.POST("/notifications/read-all", notifications.MarkAllAsRead)

	monitor.GET("/monitor/students", adminHandler.MonitorStudents)
	monitor.GET("/students", adminHandler.StudentProgress)
	monitor.GET("/students/:id", adminHandler.GetStudentDetails)
	monitor.GET("/students/:id/journey", adminHandler.StudentJourney)
	monitor.GET("/students/:id/nodes/:nodeId/files", adminHandler.ListStudentNodeFiles)
	monitor.PATCH("/students/:id/nodes/:nodeId/state", adminHandler.PatchStudentNodeState)
	monitor.PATCH("/attachments/:attachmentId/review", adminHandler.ReviewAttachment)
	monitor.POST("/attachments/:attachmentId/reviewed-document", adminHandler.UploadReviewedDocument)
	monitor.POST("/attachments/:attachmentId/presign", adminHandler.PresignReviewedDocumentUpload)
	monitor.POST("/attachments/:attachmentId/attach-reviewed", adminHandler.AttachReviewedDocument)
	monitor.POST("/reminders", adminHandler.PostReminders)

	// Admin-only routes (strict)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	admin.Use(middleware.RequireRoles("admin", "superadmin"))
	admin.GET("/users", users.ListUsers)
	admin.POST("/users", users.CreateUser)
	admin.PUT("/users/:id", users.UpdateUser)
	admin.POST("/users/:id/reset-password", users.ResetPasswordForUser)
	admin.PATCH("/users/:id/active", users.SetActive)
	admin.GET("/student-progress", adminHandler.StudentProgress) // Duplicate? Keep for now if used elsewhere
	admin.GET("/contacts", contactsHandler.AdminList)
	admin.POST("/contacts", contactsHandler.Create)
	admin.PUT("/contacts/:id", contactsHandler.Update)
	admin.DELETE("/contacts/:id", contactsHandler.Delete)

	// Global Search
	searchHandler := NewSearchHandler(db, cfg)
	api.GET("/search", middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds), searchHandler.GlobalSearch)

	// Self-service password change and profile update
	api.PATCH("/me", middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds), users.UpdateMe)
	api.POST("/me/avatar/presign", middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds), users.PresignAvatarUpload)
	api.PATCH("/me/password", middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds), users.ChangeOwnPassword)
	api.GET("/me/pending-email", middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds), users.GetPendingEmailVerification)
	api.GET("/me/verify-email", users.VerifyEmailChange) // No auth required for better UX

	// Journey state (per-user)
	js := api.Group("/journey")
	js.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
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
	chat.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	chatAdmin := chat.Group("")
	chatAdmin.Use(middleware.RequireRoles("admin", "superadmin"))
	chatAdmin.GET("/rooms/all", chatHandler.ListAllRooms) // Admin-only: list ALL rooms for tenant
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
	chat.POST("/rooms/:roomId/upload", chatHandler.UploadFile)
	chat.GET("/rooms/:roomId/attachments/:filename", chatHandler.DownloadFile)
	chat.POST("/rooms/:roomId/read", chatHandler.MarkAsRead)
	chat.PATCH("/messages/:messageId", chatHandler.UpdateMessage)
	chat.DELETE("/messages/:messageId", chatHandler.DeleteMessage)

	// Calendar routes
	calendarService := services.NewCalendarService(db)
	calendarHandler := NewCalendarHandler(calendarService)

	events := api.Group("/events")
	events.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	events.GET("", calendarHandler.GetEvents)
	events.POST("", calendarHandler.CreateEvent)
	events.PUT("/:id", calendarHandler.UpdateEvent)
	events.DELETE("/:id", calendarHandler.DeleteEvent)

	// Notification routes (generic)
	notifService := services.NewNotificationService(db)
	notifHandler := NewNotificationHandler(notifService)

	notifs := api.Group("/notifications")
	notifs.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	notifs.GET("", notifHandler.GetUnread)
	notifs.PATCH("/:id/read", notifHandler.MarkAsRead)
	notifs.POST("/read-all", notifHandler.MarkAllAsRead)

	// Analytics routes (admin only)
	analyticsService := services.NewAnalyticsService(db)
	analyticsHandler := NewAnalyticsHandler(analyticsService)

	analytics := api.Group("/analytics")
	analytics.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	analytics.Use(middleware.RequireRoles("admin", "superadmin", "chair"))
	analytics.GET("/stages", analyticsHandler.GetStageStats)
	analytics.GET("/overdue", analyticsHandler.GetOverdueStats)

	// TODO: checklist, documents, comments handlers (skeletons for now)

	// Dictionary endpoints
	dictHandler := NewDictionaryHandler(db)

	// Public (Authenticated) Read Access
	dicts := api.Group("/dictionaries")
	dicts.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	{
		dicts.GET("/programs", dictHandler.ListPrograms)
		dicts.GET("/specialties", dictHandler.ListSpecialties)
		dicts.GET("/cohorts", dictHandler.ListCohorts)
		dicts.GET("/departments", dictHandler.ListDepartments)
	}

	// Admin Write Access
	dictAdminGroup := admin.Group("/dictionaries")
	{
		dictAdminGroup.POST("/programs", dictHandler.CreateProgram)
		dictAdminGroup.PUT("/programs/:id", dictHandler.UpdateProgram)
		dictAdminGroup.DELETE("/programs/:id", dictHandler.DeleteProgram)

		dictAdminGroup.POST("/specialties", dictHandler.CreateSpecialty)
		dictAdminGroup.PUT("/specialties/:id", dictHandler.UpdateSpecialty)
		dictAdminGroup.DELETE("/specialties/:id", dictHandler.DeleteSpecialty)

		dictAdminGroup.POST("/cohorts", dictHandler.CreateCohort)
		dictAdminGroup.PUT("/cohorts/:id", dictHandler.UpdateCohort)
		dictAdminGroup.DELETE("/cohorts/:id", dictHandler.DeleteCohort)

		dictAdminGroup.POST("/departments", dictHandler.CreateDepartment)
		dictAdminGroup.PUT("/departments/:id", dictHandler.UpdateDepartment)
		dictAdminGroup.DELETE("/departments/:id", dictHandler.DeleteDepartment)
	}

	// ===========================================
	// SUPERADMIN ROUTES (global platform admin)
	// ===========================================
	superadminTenantsHandler := NewSuperadminTenantsHandler(db, cfg)
	superadminAdminsHandler := NewSuperadminAdminsHandler(db, cfg)
	superadminLogsHandler := NewSuperadminLogsHandler(db, cfg)
	superadminSettingsHandler := NewSuperadminSettingsHandler(db, cfg)

	superadmin := api.Group("/superadmin")
	superadmin.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	superadmin.Use(middleware.RequireSuperadmin())
	{
		// Tenants management
		superadmin.GET("/tenants", superadminTenantsHandler.ListTenants)
		superadmin.POST("/tenants", superadminTenantsHandler.CreateTenant)
		superadmin.GET("/tenants/:id", superadminTenantsHandler.GetTenant)
		superadmin.PUT("/tenants/:id", superadminTenantsHandler.UpdateTenant)
		superadmin.DELETE("/tenants/:id", superadminTenantsHandler.DeleteTenant)
		superadmin.POST("/tenants/:id/logo", superadminTenantsHandler.UploadLogo)
		superadmin.PUT("/tenants/:id/services", superadminTenantsHandler.UpdateTenantServices)

		// Admins management
		superadmin.GET("/admins", superadminAdminsHandler.ListAdmins)
		superadmin.POST("/admins", superadminAdminsHandler.CreateAdmin)
		superadmin.GET("/admins/:id", superadminAdminsHandler.GetAdmin)
		superadmin.PUT("/admins/:id", superadminAdminsHandler.UpdateAdmin)
		superadmin.DELETE("/admins/:id", superadminAdminsHandler.DeleteAdmin)
		superadmin.POST("/admins/:id/reset-password", superadminAdminsHandler.ResetPassword)

		// Activity logs
		superadmin.GET("/logs", superadminLogsHandler.ListLogs)
		superadmin.GET("/logs/stats", superadminLogsHandler.GetLogStats)
		superadmin.GET("/logs/actions", superadminLogsHandler.GetActions)
		superadmin.GET("/logs/entity-types", superadminLogsHandler.GetEntityTypes)

		// Global settings
		superadmin.GET("/settings", superadminSettingsHandler.ListSettings)
		superadmin.GET("/settings/categories", superadminSettingsHandler.GetCategories)
		superadmin.GET("/settings/:key", superadminSettingsHandler.GetSetting)
		superadmin.PUT("/settings/:key", superadminSettingsHandler.UpdateSetting)
		superadmin.DELETE("/settings/:key", superadminSettingsHandler.DeleteSetting)
		superadmin.POST("/settings/bulk", superadminSettingsHandler.BulkUpdate)
	}

	return r
}
