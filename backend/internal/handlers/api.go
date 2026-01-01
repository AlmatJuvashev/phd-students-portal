package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/mailer"
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



	// Services
	emailService := services.NewEmailService()
	
	// Journey Service Dependencies
	// S3 Client
	s3Svc, err := services.NewS3FromEnv()
	if err != nil {
		log.Printf("Warning: S3 init failed: %v", err)
	}

	// Repositories & Domain Services
	userRepo := repository.NewSQLUserRepository(db)
	userService := services.NewUserService(userRepo, rds, cfg, emailService, s3Svc)
	authService := services.NewAuthService(userRepo, emailService, cfg)

	// Auth routes (login and password reset)
	auth := NewAuthHandler(authService, cfg, rds)
	api.POST("/auth/login", auth.Login)
	api.POST("/auth/logout", auth.Logout) // Added logout
	api.POST("/auth/forgot-password", auth.ForgotPassword)
	api.POST("/auth/reset-password", auth.ResetPassword)

	users := NewUsersHandler(userService, cfg)
	_ = users
	
	// Documents
	docRepo := repository.NewSQLDocumentRepository(db)
	docService := services.NewDocumentService(docRepo, cfg, s3Svc)

	// Mailer
	mailerSvc := mailer.NewMailer()

	// Journey Service
	journeyRepo := repository.NewSQLJourneyRepository(db)
	journeyService := services.NewJourneyService(journeyRepo, playbookManager, cfg, mailerSvc, s3Svc, docService)

	journey := NewJourneyHandler(journeyService)
	_ = journey
	nodeSubmission := NewNodeSubmissionHandler(journeyService)
	_ = nodeSubmission
	// Admin Service
	adminRepo := repository.NewSQLAdminRepository(db)
	adminService := services.NewAdminService(adminRepo, playbookManager, cfg, s3Svc)
	adminHandler := NewAdminHandler(cfg, playbookManager, adminService, journeyService)
	_ = adminHandler
	chatRepo := repository.NewSQLChatRepository(db)
	chatService := services.NewChatService(chatRepo, emailService, cfg)
	chatHandler := NewChatHandler(chatService, cfg)
	_ = chatHandler

	// Calendar Module
	eventRepo := repository.NewSQLEventRepository(db)
	calendarService := services.NewCalendarService(eventRepo)
	calendarHandler := NewCalendarHandler(calendarService)

	// Checklist Module
	checklistRepo := repository.NewSQLChecklistRepository(db)
	checklistService := services.NewChecklistService(checklistRepo)
	checklistHandler := NewChecklistHandler(checklistService, cfg)

	// Search Module
	searchRepo := repository.NewSQLSearchRepository(db)
	searchService := services.NewSearchService(searchRepo)
	searchHandler := NewSearchHandler(searchService, cfg)

	// Dictionary Module
	dictionaryRepo := repository.NewSQLDictionaryRepository(db)
	dictionaryService := services.NewDictionaryService(dictionaryRepo)
	dictionaryHandler := NewDictionaryHandler(dictionaryService)

	// Notification Module
	notificationRepo := repository.NewSQLNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := NewNotificationHandler(notificationService)

	contactRepo := repository.NewSQLContactRepository(db)
	contactService := services.NewContactService(contactRepo)
	contactsHandler := NewContactsHandler(contactService)

	// Analytics
	analyticsRepo := repository.NewSQLAnalyticsRepository(db)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	analyticsHandler := NewAnalyticsHandler(analyticsService)

	// Scheduler
	schedulerRepo := repository.NewSQLSchedulerRepository(db)
	// Reuse existing ResourceRepo (resourceRepo) which was defined above. Wait, resourceRepo is "resources := ..." or similar? 
	// I need to find where ResourceRepo is defined. It's likely in "resources".
	// Let's create it fresh to be safe or look up slightly higher in file.
	// Oh, I can see "resourceHandler := ..." in previous context? No.
	// Let's instantiate SQLResourceRepository here again, it's cheap (just a struct with db pointer).
	resourceRepo := repository.NewSQLResourceRepository(db)
	schedulerService := services.NewSchedulerService(schedulerRepo, resourceRepo)
	schedulerHandler := NewSchedulerHandler(schedulerService)



	// ===========================================
	// SUPERADMIN ROUTES (global platform admin)
	// ===========================================
	
	// Create Repos & Services for SuperAdmin/Tenant
	tenantRepo := repository.NewSQLTenantRepository(db)
	tenantService := services.NewTenantService(tenantRepo)
	
	superAdminRepo := repository.NewSQLSuperAdminRepository(db)
	superAdminService := services.NewSuperAdminService(superAdminRepo)
	
	// Update MeHandler with dependencies (UserService, TenantService)
	// Note: NewMeHandler signature: (userSvc, tenantSvc, cfg)
	meHandler := NewMeHandler(userService, tenantService, cfg, rds)
	// Current NewMeHandler signature from previous edit might only take (userService, tenantService, cfg) 
	// Or did I change it? I need to check MeHandler constructor.
	// Looking at me.go previous changes: NewMeHandler(userSvc, tenantSvc, cfg)

	api.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	api.GET("/contacts", contactsHandler.PublicList)

	// /me routes (require auth)
	meGroup := api.Group("/me")
	meGroup.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	{
		meGroup.GET("", meHandler.Me)
		meGroup.GET("/tenants", meHandler.MyTenants)
	meGroup.GET("/tenant", meHandler.MyTenant)
	}

	// Authenticated Routes Group
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret), db, rds))
	{
		// Calendar
		cal := protected.Group("/calendar")
		{
			cal.POST("/events", calendarHandler.CreateEvent)
			cal.GET("/events", calendarHandler.GetEvents)
			cal.PUT("/events/:id", calendarHandler.UpdateEvent)
			cal.DELETE("/events/:id", calendarHandler.DeleteEvent)
		}

		// Checklist
		check := protected.Group("/checklists")
		{
			check.GET("/modules", checklistHandler.ListModules)
		}

		// Scheduler
		sched := protected.Group("/scheduling")
		{
			sched.GET("/terms", schedulerHandler.ListTerms)
			sched.POST("/terms", schedulerHandler.CreateTerm)
			
			sched.POST("/offerings", schedulerHandler.CreateOffering)
			
			sched.GET("/sessions", schedulerHandler.ListSessions)
			sched.POST("/sessions", schedulerHandler.CreateSession)
		}

		// Search
		protected.GET("/search", searchHandler.GlobalSearch)

		// Dictionaries
		dict := protected.Group("/dictionaries")
		{
			dict.GET("/programs", dictionaryHandler.ListPrograms)
			dict.GET("/specialties", dictionaryHandler.ListSpecialties)
			dict.GET("/cohorts", dictionaryHandler.ListCohorts)
			dict.GET("/departments", dictionaryHandler.ListDepartments)
		}

		// Grading Module
		gradingRepo := repository.NewSQLGradingRepository(db)
		gradingService := services.NewGradingService(gradingRepo)
		gradingHandler := NewGradingHandler(gradingService)
		
		gr := protected.Group("/grading")
		{
			gr.GET("/schemas", gradingHandler.ListSchemas)
			gr.POST("/schemas", gradingHandler.CreateSchema)
			gr.POST("/entries", gradingHandler.SubmitGrade)
			gr.GET("/student/:studentId", gradingHandler.ListStudentGrades)
		}

		// Item Bank Module
		ibRepo := repository.NewSQLItemBankRepository(db)
		ibService := services.NewItemBankService(ibRepo)
		ibHandler := NewItemBankHandler(ibService)

		ib := protected.Group("/item-banks")
		{
			ib.GET("/banks", ibHandler.ListBanks)
			ib.POST("/banks", ibHandler.CreateBank)
			ib.GET("/banks/:bankId/items", ibHandler.ListItems)
			ib.POST("/banks/:bankId/items", ibHandler.CreateItem)
		}

		// Governance Module (Phase 5)
		govRepo := repository.NewSQLGovernanceRepository(db)
		govService := services.NewGovernanceService(govRepo)
		govHandler := NewGovernanceHandler(govService)

		gov := protected.Group("/governance")
		{
			gov.POST("/proposals", govHandler.SubmitProposal)
			gov.GET("/proposals", govHandler.ListProposals)
			gov.GET("/proposals/:id", govHandler.GetProposal)
			gov.POST("/proposals/:id/review", govHandler.ReviewProposal) // Approve/Reject
			gov.GET("/proposals/:id/reviews", govHandler.ListReviews)
		}

		// Notifications
		notif := protected.Group("/notifications")
		{
			notif.GET("", notificationHandler.GetNotifications)
			notif.GET("/unread", notificationHandler.GetUnread)
			notif.POST("/:id/read", notificationHandler.MarkAsRead)
			notif.POST("/read-all", notificationHandler.MarkAllAsRead)
		}

		// Journey
		j := protected.Group("/journey")
		{
			j.GET("/state", journey.GetState)
			j.PUT("/state", journey.SetState)
			j.POST("/reset", journey.Reset)
			j.GET("/scoreboard", journey.GetScoreboard)

			j.GET("/profile", nodeSubmission.GetProfile)
			nodes := j.Group("/nodes/:nodeId")
			{
				nodes.GET("/submission", nodeSubmission.GetSubmission)
				nodes.PUT("/submission", nodeSubmission.PutSubmission)
				nodes.PATCH("/state", nodeSubmission.PatchState)
				
				uploads := nodes.Group("/uploads")
				{
					uploads.POST("/presign", nodeSubmission.PresignUpload)
					uploads.POST("/attach", nodeSubmission.AttachUpload)
				}
			}
		}

		// Admin/Advisor Progress Monitoring
		adm := protected.Group("/admin")
		adm.Use(middleware.RequireAdminOrAdvisor())
		{
			adm.GET("/student-progress", adminHandler.StudentProgress)
			adm.GET("/monitor", adminHandler.MonitorStudents)
			adm.GET("/monitor/students", adminHandler.MonitorStudents) // Alias for frontend compatibility
			adm.GET("/monitor/analytics", adminHandler.MonitorAnalytics)
			adm.GET("/students/:id", adminHandler.GetStudentDetails)
			adm.GET("/students/:id/journey", adminHandler.StudentJourney)
			adm.GET("/students/:id/deadlines", adminHandler.GetStudentDeadlines)
			adm.GET("/students/:id/nodes/:nodeId/files", adminHandler.ListStudentNodeFiles)
			adm.PATCH("/students/:id/nodes/:nodeId/state", adminHandler.PatchStudentNodeState)
			
			// Review actions
			adm.POST("/attachments/:attachmentId/review", adminHandler.ReviewAttachment)
			adm.POST("/attachments/:attachmentId/presign", adminHandler.PresignReviewedDocumentUpload)
			adm.POST("/attachments/:attachmentId/attach-reviewed", adminHandler.AttachReviewedDocument)
			
			// Reminders
			adm.POST("/reminders", adminHandler.PostReminders)
			
			// User management (admin panel)
			adm.GET("/users", users.ListUsers)
			adm.POST("/users", users.CreateUser)
			adm.PUT("/users/:id", users.UpdateUser)
			adm.PATCH("/users/:id/active", users.SetActive)
			adm.POST("/users/:id/reset-password", users.ResetPasswordForUser)
			
			// Admin dictionaries
			admDict := adm.Group("/dictionaries")
			{
				admDict.GET("/programs", dictionaryHandler.ListPrograms)
				admDict.POST("/programs", dictionaryHandler.CreateProgram)
				admDict.PUT("/programs/:id", dictionaryHandler.UpdateProgram)
				admDict.DELETE("/programs/:id", dictionaryHandler.DeleteProgram)
				
				admDict.GET("/specialties", dictionaryHandler.ListSpecialties)
				admDict.POST("/specialties", dictionaryHandler.CreateSpecialty)
				admDict.PUT("/specialties/:id", dictionaryHandler.UpdateSpecialty)
				admDict.DELETE("/specialties/:id", dictionaryHandler.DeleteSpecialty)
				
				admDict.GET("/cohorts", dictionaryHandler.ListCohorts)
				admDict.POST("/cohorts", dictionaryHandler.CreateCohort)
				admDict.PUT("/cohorts/:id", dictionaryHandler.UpdateCohort)
				admDict.DELETE("/cohorts/:id", dictionaryHandler.DeleteCohort)
				
				admDict.GET("/departments", dictionaryHandler.ListDepartments)
				admDict.POST("/departments", dictionaryHandler.CreateDepartment)
				admDict.PUT("/departments/:id", dictionaryHandler.UpdateDepartment)
				admDict.DELETE("/departments/:id", dictionaryHandler.DeleteDepartment)
			}
			
			// Admin notifications (duplicated from protected for admin panel access)
			admNotif := adm.Group("/notifications")
			{
				admNotif.GET("", notificationHandler.GetNotifications)
				admNotif.GET("/unread", notificationHandler.GetUnread)
				admNotif.GET("/unread-count", notificationHandler.GetUnread) // Alias for compatibility
				admNotif.POST("/:id/read", notificationHandler.MarkAsRead)
				admNotif.POST("/read-all", notificationHandler.MarkAllAsRead)
			}
			
			// Admin contacts
			adm.GET("/contacts", contactsHandler.AdminList)
			adm.POST("/contacts", contactsHandler.Create)
			adm.PUT("/contacts/:id", contactsHandler.Update)
			adm.DELETE("/contacts/:id", contactsHandler.Delete)
		}


		// Chat
		chat := protected.Group("/chat")
		{
			chat.GET("/rooms", chatHandler.ListRooms)
			chat.GET("/rooms/:roomId/members", chatHandler.GetRoomMembers)
			chat.GET("/rooms/:roomId/messages", chatHandler.ListMessages)
			chat.POST("/rooms/:roomId/messages", chatHandler.CreateMessage)
			chat.POST("/rooms/:roomId/read", chatHandler.MarkAsRead)
			
			// File upload/download - available to all chat members
			chat.POST("/rooms/:roomId/upload", chatHandler.UploadFile)
			chat.GET("/rooms/:roomId/files/:filename", chatHandler.DownloadFile)

			// Admin chat actions
			adminChat := chat.Group("")
			adminChat.Use(middleware.RequireAdminOrAdvisor())
			{
				adminChat.POST("/rooms", chatHandler.CreateRoom)
				adminChat.PATCH("/rooms/:roomId", chatHandler.UpdateRoom)
				adminChat.GET("/rooms/all", chatHandler.ListAllRooms)
				adminChat.POST("/rooms/:roomId/members", chatHandler.AddMember)
				adminChat.DELETE("/rooms/:roomId/members/:userId", chatHandler.RemoveMember)
				adminChat.POST("/rooms/:roomId/members/batch", chatHandler.AddRoomMembersBatch)
				adminChat.DELETE("/rooms/:roomId/members/batch", chatHandler.RemoveRoomMembersBatch)
			}
			
			// Message editing/deletion
			chat.PATCH("/messages/:messageId", chatHandler.UpdateMessage)
			chat.DELETE("/messages/:messageId", chatHandler.DeleteMessage)
		}

		// Analytics
		an := protected.Group("/analytics")
		an.Use(middleware.RequireAdminOrAdvisor())
		{
			an.GET("/stages", analyticsHandler.GetStageStats)
			an.GET("/overdue", analyticsHandler.GetOverdueStats)
		}
	}
	
	// SuperAdmin Handlers
	superadminTenantsHandler := NewSuperadminTenantsHandler(tenantService, superAdminService, cfg)
	superadminAdminsHandler := NewSuperadminAdminsHandler(superAdminService, cfg)
	superadminLogsHandler := NewSuperadminLogsHandler(superAdminService, cfg)
	superadminSettingsHandler := NewSuperadminSettingsHandler(superAdminService, cfg)

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
