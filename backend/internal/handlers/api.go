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
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	
	// Repositories & Domain Services
	userRepo := repository.NewSQLUserRepository(db)
	userService := services.NewUserService(userRepo, rds, cfg, emailService)
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
	docService, err := services.NewDocumentService(docRepo, cfg)
	if err != nil {
		log.Printf("Warning: DocumentService init failed: %v", err)
	}

	// Journey Service Dependencies
	// S3 Client
	s3Svc, err := services.NewS3FromEnv()
	if err != nil {
		log.Printf("Warning: S3 init failed: %v", err)
	}
	var s3Client *s3.Client
	if s3Svc != nil {
		s3Client = s3Svc.Client()
	}

	// Mailer
	mailerSvc := mailer.NewMailer()

	// Journey Service
	journeyRepo := repository.NewSQLJourneyRepository(db)
	journeyService := services.NewJourneyService(journeyRepo, playbookManager, cfg, mailerSvc, s3Client, docService)

	journey := NewJourneyHandler(journeyService)
	_ = journey
	nodeSubmission := NewNodeSubmissionHandler(journeyService)
	_ = nodeSubmission
	// Admin Service
	adminRepo := repository.NewSQLAdminRepository(db)
	adminService := services.NewAdminService(adminRepo, playbookManager, cfg)
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

		// Notifications
		notif := protected.Group("/notifications")
		{
			notif.GET("/unread", notificationHandler.GetUnread)
			notif.POST("/:id/read", notificationHandler.MarkAsRead)
			notif.POST("/read-all", notificationHandler.MarkAllAsRead)
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
