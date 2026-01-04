package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
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
	tenantRepo := repository.NewSQLTenantRepository(db)
	userService := services.NewUserService(userRepo, tenantRepo, rds, cfg, emailService, s3Svc)
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
	// Analytics
	analyticsRepo := repository.NewSQLAnalyticsRepository(db)
	// Needs LMS and Attendance repos, which are defined later in the file usually.
	// I need to move their definition UP or move Analytics definition DOWN.
	// LMS and Attendance are defined around lines 205 and 211.
	// I will move Analytics instantiation DOWN after them.

	// Moving "analyticsService := ..." line to be AFTER lmsRepo and attendanceRepo are initialized.
	// So I will comment out/remove the lines here and insert them later.
	// But `analyticsHandler` is used? No, just created here.

	// Scheduler
	schedulerRepo := repository.NewSQLSchedulerRepository(db)
	resourceRepo := repository.NewSQLResourceRepository(db)
	resourceService := services.NewResourceService(resourceRepo)
	resourceHandler := NewResourceHandler(resourceService)

	curriculumRepo := repository.NewSQLCurriculumRepository(db)
	curriculumService := services.NewCurriculumService(curriculumRepo)
	curriculumHandler := NewCurriculumHandler(curriculumService)

	// Program Builder
	programBuilderService := services.NewProgramBuilderService(curriculumRepo)
	programBuilderHandler := NewProgramBuilderHandler(programBuilderService)

	// Mailer for Scheduler
	smtpMailer := mailer.NewMailer()

	schedulerService := services.NewSchedulerService(schedulerRepo, resourceRepo, curriculumRepo, userRepo, smtpMailer)
	schedulerHandler := NewSchedulerHandler(schedulerService)

	lmsRepo := repository.NewSQLLMSRepository(db)
	gradingRepo := repository.NewSQLGradingRepository(db)
	teacherService := services.NewTeacherService(schedulerRepo, lmsRepo, gradingRepo)
	teacherHandler := NewTeacherHandler(teacherService)

	studentService := services.NewStudentService(userRepo, journeyRepo, lmsRepo, schedulerRepo, curriculumRepo, gradingRepo, playbookManager)
	studentHandler := NewStudentHandler(studentService)

	// Attendance (Phase 17)
	attendanceRepo := repository.NewSQLAttendanceRepository(db)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	attendanceHandler := NewAttendanceHandler(attendanceService)
	// Initialize AnalyticsService (depends on AnalyticsRepo, LMSRepo, AttendanceRepo, UserRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, lmsRepo, attendanceRepo, userRepo)
	analyticsHandler := NewAnalyticsHandler(analyticsService)

	// AI Assistant
	aiService := services.NewAIService(cfg)
	aiHandler := NewAIHandler(aiService)

	// Transcript (Phase 17)
	transcriptRepo := repository.NewSQLTranscriptRepository(db)
	transcriptService := services.NewTranscriptService(transcriptRepo, schedulerRepo)
	transcriptHandler := NewTranscriptHandler(transcriptService)

	// Bulk Operations (Phase 17)
	// Bulk Operations (Phase 17)
	bulkService := services.NewBulkEnrollmentService(userService)
	bulkHandler := NewBulkHandler(bulkService)

	// LTI 1.3 (Phase 19)
	ltiRepo := repository.NewSQLLTIRepository(db)
	ltiService := services.NewLTIService(ltiRepo, cfg)
	ltiHandler := NewLTIHandler(ltiService, cfg)
	api.GET("/lti/login_init", ltiHandler.LoginInit)
	api.POST("/lti/launch", ltiHandler.Launch)
	api.GET("/.well-known/jwks.json", ltiHandler.GetJWKS)

	// ===========================================
	// RBAC (Phase 20)
	// ===========================================
	rbacRepo := repository.NewSQLRBACRepository(db)
	authzSvc := services.NewAuthzService(rbacRepo)
	rbac := middleware.NewRBACMiddleware(authzSvc)

	// ===========================================
	// External Audit (Phase 24)
	// ===========================================
	auditRepo := repository.NewSQLAuditRepository(db)
	auditService := services.NewAuditService(auditRepo, curriculumRepo)
	auditHandler := NewAuditHandler(auditService, curriculumService)

	// ===========================================
	// SUPERADMIN ROUTES (global platform admin)
	// ===========================================

	// Create Repos & Services for SuperAdmin/Tenant
	// tenantRepo moved up
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
		// Student Portal
		student := protected.Group("/student")
		student.Use(middleware.RequireRoles("student"))
		{
			student.GET("/dashboard", studentHandler.GetDashboard)
			student.GET("/courses", studentHandler.ListCourses)
			student.GET("/assignments", studentHandler.ListAssignments)
			student.GET("/grades", studentHandler.ListGrades)
			student.GET("/transcript", transcriptHandler.GetStudentTranscript)
		}

		// Calendar
		cal := protected.Group("/calendar")
		{
			cal.POST("/events", calendarHandler.CreateEvent)
			cal.GET("/events", calendarHandler.GetEvents)
			cal.PUT("/events/:id", calendarHandler.UpdateEvent)
			cal.DELETE("/events/:id", calendarHandler.DeleteEvent)
		}

		// Teacher / Faculty Dashboard
		teacher := protected.Group("/teacher")
		// Use new RBAC Middleware
		// Global check: Must have at least 'course.view' globally OR we can rely on specific endpoint checks
		// For dashboard, we check 'course.view' globally for now (as instructor/student have it)
		teacher.Use(rbac.RequirePermission("course.view", models.ContextGlobal, ""))
		{
			teacher.GET("/dashboard", teacherHandler.GetDashboardStats)
			teacher.GET("/courses", teacherHandler.GetMyCourses)
			// Context-Aware check: Must have 'course.view' for this SPECIFIC course ID
			teacher.GET("/courses/:id/roster", rbac.RequirePermission("course.view", models.ContextCourse, "id"), teacherHandler.GetCourseRoster)
			teacher.GET("/courses/:id/gradebook", teacherHandler.GetGradebook)
			teacher.GET("/submissions", teacherHandler.GetSubmissions)
			teacher.POST("/submissions/:id/annotations", teacherHandler.AddAnnotation)
			teacher.GET("/submissions/:id/annotations", teacherHandler.GetAnnotations)
			teacher.DELETE("/submissions/:id/annotations/:annId", teacherHandler.DeleteAnnotation)
			teacher.POST("/sessions/:session_id/attendance", attendanceHandler.BatchRecordAttendance)
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
			sched.GET("/offerings", schedulerHandler.ListOfferings)

			sched.GET("/sessions", schedulerHandler.ListSessions)
			sched.POST("/sessions", schedulerHandler.CreateSession)
			sched.POST("/optimize", schedulerHandler.AutoSchedule)
		}

		// Resources
		res := protected.Group("/resources")
		{
			res.GET("/buildings", resourceHandler.ListBuildings)
			res.GET("/rooms", resourceHandler.ListRooms)
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

		// Item Bank Module (Powered by Assessment Engine)
		assessmentRepo := repository.NewSQLAssessmentRepository(db)
		ibService := services.NewItemBankService(assessmentRepo)
		ibHandler := NewItemBankHandler(ibService)

		ib := protected.Group("/item-banks")
		{
			ib.GET("/banks", ibHandler.ListBanks)
			ib.POST("/banks", ibHandler.CreateBank)
			ib.PUT("/banks/:bankId", ibHandler.UpdateBank)
			ib.DELETE("/banks/:bankId", ibHandler.DeleteBank)
			ib.GET("/banks/:bankId/items", ibHandler.ListItems)
			ib.POST("/banks/:bankId/items", ibHandler.CreateItem)
			ib.PUT("/banks/:bankId/items/:itemId", ibHandler.UpdateItem)
			ib.DELETE("/banks/:bankId/items/:itemId", ibHandler.DeleteItem)
		}

		// Assessment Engine Module (Phase 25)
		assessmentService := services.NewAssessmentService(assessmentRepo)
		assessmentHandler := NewAssessmentHandler(assessmentService)

		ae := protected.Group("/assessments")
		{
			ae.POST("", assessmentHandler.CreateAssessment)
			ae.GET("", assessmentHandler.ListAssessments)
			ae.GET("/:id", assessmentHandler.GetAssessment)
			ae.PUT("/:id", assessmentHandler.UpdateAssessment)
			ae.DELETE("/:id", assessmentHandler.DeleteAssessment)
			ae.POST("/:id/attempts", assessmentHandler.StartAttempt)
			ae.GET("/:id/my-attempts", assessmentHandler.ListMyAttempts)
		}

		at := protected.Group("/attempts")
		{
			at.GET("/:id", assessmentHandler.GetAttemptDetails)
			at.POST("/:id/response", assessmentHandler.SubmitResponse)
			at.POST("/:id/complete", assessmentHandler.CompleteAttempt)
			at.POST("/:id/log", assessmentHandler.LogProctoringEvent)
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

		// Admin/Advisor Progress Monitoring & Review
		adm := protected.Group("/admin")
		// Note: We deliberately do NOT apply a blanket middleware here anymore.
		// Granular permissions are applied below.

		// 1. Monitor & Student Data (Accessible to Advisor)
		// Includes: Progress, Journey, Attachments, Reminders
		monitorGrp := adm.Group("")
		monitorGrp.Use(middleware.RequireRoles("admin", "advisor", "superadmin"))
		{
			monitorGrp.GET("/student-progress", adminHandler.StudentProgress)
			monitorGrp.GET("/monitor", adminHandler.MonitorStudents)
			monitorGrp.GET("/monitor/students", adminHandler.MonitorStudents)
			monitorGrp.GET("/monitor/analytics", adminHandler.MonitorAnalytics)
			monitorGrp.GET("/students/:id", adminHandler.GetStudentDetails)
			monitorGrp.GET("/students/:id/journey", adminHandler.StudentJourney)
			monitorGrp.GET("/students/:id/deadlines", adminHandler.GetStudentDeadlines)
			monitorGrp.GET("/students/:id/nodes/:nodeId/files", adminHandler.ListStudentNodeFiles)
			monitorGrp.PATCH("/students/:id/nodes/:nodeId/state", adminHandler.PatchStudentNodeState)

			// Review actions
			monitorGrp.POST("/attachments/:attachmentId/review", adminHandler.ReviewAttachment)
			monitorGrp.POST("/attachments/:attachmentId/presign", adminHandler.PresignReviewedDocumentUpload)
			monitorGrp.POST("/attachments/:attachmentId/attach-reviewed", adminHandler.AttachReviewedDocument)

			// Reminders
			monitorGrp.POST("/reminders", adminHandler.PostReminders)

			// Admin notifications (Advisors need to see alerts too)
			admNotif := monitorGrp.Group("/notifications")
			{
				admNotif.GET("", notificationHandler.GetNotifications)
				admNotif.GET("/unread", notificationHandler.GetUnread)
				admNotif.GET("/unread-count", notificationHandler.GetUnread)
				admNotif.POST("/:id/read", notificationHandler.MarkAsRead)
				admNotif.POST("/read-all", notificationHandler.MarkAllAsRead)
			}

			// Contacts Search (Advisors might need this)
			monitorGrp.GET("/contacts", contactsHandler.AdminList)
		}

		// 2. User Management (IT Admin only - NO Advisors)
		sysAdminGrp := adm.Group("")
		sysAdminGrp.Use(middleware.RequireRoles("admin", "it_admin", "superadmin"))
		{
			sysAdminGrp.GET("/users", users.ListUsers)
			sysAdminGrp.POST("/users", users.CreateUser)
			sysAdminGrp.PUT("/users/:id", users.UpdateUser)
			sysAdminGrp.PATCH("/users/:id/active", users.SetActive)
			sysAdminGrp.POST("/users/:id/reset-password", users.ResetPasswordForUser)
			sysAdminGrp.POST("/bulk/enroll", bulkHandler.BulkEnrollStudents)

			// Manage Contacts
			sysAdminGrp.POST("/contacts", contactsHandler.Create)
			sysAdminGrp.PUT("/contacts/:id", contactsHandler.Update)
			sysAdminGrp.DELETE("/contacts/:id", contactsHandler.Delete)

			// LTI Tool Management
			sysAdminGrp.POST("/lti/tools", ltiHandler.RegisterTool)
			sysAdminGrp.GET("/lti/tools", ltiHandler.ListTools)
		}

		// 3. Content / Dictionaries (Registrar / Content Mgr)
		// Advisors might need READ access, but Write should be restricted?
		// For simplicity, let's allow Admins/Registrars/ContentMgr/Superadmin modify.
		// Advisors: Read Only? Current handlers don't split R/W easily without more code.
		// Let's implement full access for Content Team.
		contentGrp := adm.Group("")
		contentGrp.Use(middleware.RequireRoles("admin", "registrar", "content_manager", "superadmin"))
		{
			// Admin dictionaries
			admDict := contentGrp.Group("/dictionaries")
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

			// Program Builder
			admBuilder := contentGrp.Group("/programs/:id/builder")
			{
				admBuilder.GET("/map", programBuilderHandler.GetJourneyMap)
				admBuilder.PUT("/map", programBuilderHandler.UpdateJourneyMap)
				admBuilder.GET("/nodes", programBuilderHandler.GetNodes)
				admBuilder.POST("/nodes", programBuilderHandler.CreateNode)
				admBuilder.PUT("/nodes/:nodeId", programBuilderHandler.UpdateNode)
			}

			// AI Assistant (Content Generation)
			contentGrp.POST("/ai/generate-course", aiHandler.GenerateCourseStructure)
			contentGrp.POST("/ai/generate-quiz", aiHandler.GenerateQuiz)
			contentGrp.POST("/ai/generate-survey", aiHandler.GenerateSurvey)
			contentGrp.POST("/ai/generate-assessment-items", aiHandler.GenerateAssessmentItems)
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
			adminChat.Use(rbac.RequirePermission("user.view", models.ContextGlobal, ""))
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
		an.Use(rbac.RequirePermission("user.view", models.ContextGlobal, "")) // Advisor/Admin access
		{
			an.GET("/stages", analyticsHandler.GetStageStats)
			an.GET("/overdue", analyticsHandler.GetOverdueStats)
			an.GET("/monitor", analyticsHandler.GetMonitorMetrics)
			an.GET("/at-risk", analyticsHandler.GetHighRiskStudents)
			an.POST("/batch-analysis", analyticsHandler.HandleBatchRiskAnalysis)
		}

		// Curriculum Management (Registrar/ContentMgr)
		curr := protected.Group("/curriculum")
		curr.Use(middleware.RequireRoles("admin", "registrar", "content_manager", "superadmin")) // "admin" kept for legacy compatibility
		{
			// Programs
			curr.GET("/programs", curriculumHandler.ListPrograms)
			curr.POST("/programs", curriculumHandler.CreateProgram)
			curr.GET("/programs/:id", curriculumHandler.GetProgram)
			curr.PUT("/programs/:id", curriculumHandler.UpdateProgram)
			curr.DELETE("/programs/:id", curriculumHandler.DeleteProgram)

			// Program Builder (Program Versions / Journey Map editor)
			currBuilder := curr.Group("/programs/:id/builder")
			{
				currBuilder.GET("/map", programBuilderHandler.GetJourneyMap)
				currBuilder.PUT("/map", programBuilderHandler.UpdateJourneyMap)
				currBuilder.GET("/nodes", programBuilderHandler.GetNodes)
				currBuilder.POST("/nodes", programBuilderHandler.CreateNode)
				currBuilder.PUT("/nodes/:nodeId", programBuilderHandler.UpdateNode)
			}

			// Courses
			curr.GET("/courses", curriculumHandler.ListCourses)
			curr.POST("/courses", curriculumHandler.CreateCourse)
			// curr.PUT("/courses/:id", curriculumHandler.UpdateCourse) // TODO: Add if implemented
			// curr.DELETE("/courses/:id", curriculumHandler.DeleteCourse) // TODO: Add if implemented
		}

		// ===========================================
		// External Audit (Phase 24)
		// Read-only access for external examiners and internal auditors
		// ===========================================
		audit := protected.Group("/audit")
		audit.Use(middleware.RequireRoles("external", "admin", "superadmin", "registrar"))
		{
			audit.GET("/programs", auditHandler.ListPrograms)
			audit.GET("/programs/:id", auditHandler.GetProgram)
			audit.GET("/courses", auditHandler.ListCourses)
			audit.GET("/outcomes", auditHandler.ListOutcomes)
			audit.GET("/changelog", auditHandler.ListChangeLog)
			audit.GET("/reports/program-summary", auditHandler.ProgramSummaryReport)
		}

		// Admin-only endpoints for managing outcomes
		outcomes := protected.Group("/outcomes")
		outcomes.Use(middleware.RequireRoles("admin", "superadmin", "registrar", "content_manager"))
		{
			outcomes.POST("", auditHandler.CreateOutcome)
			outcomes.PUT("/:id", auditHandler.UpdateOutcome)
			outcomes.DELETE("/:id", auditHandler.DeleteOutcome)
		}

		// Discussion Forums (Phase 26)
		forumRepo := repository.NewSQLForumRepository(db)
		forumService := services.NewForumService(forumRepo)
		forumHandler := NewForumHandler(forumService)

		// Course-level
		protected.GET("/courses/:id/forums", forumHandler.ListForums)
		protected.POST("/courses/:id/forums", forumHandler.CreateForum)

		// Forum-level
		forums := protected.Group("/forums")
		{
			forums.GET("/:id/topics", forumHandler.ListTopics)
			forums.POST("/:id/topics", forumHandler.CreateTopic)
		}

		// Topic-level
		topics := protected.Group("/topics")
		{
			topics.GET("/:id", forumHandler.GetTopic)
			topics.POST("/:id/posts", forumHandler.CreatePost)
		}

		// Rubric Grading (Phase 27)
		rubricRepo := repository.NewSQLRubricRepository(db)
		rubricService := services.NewRubricService(rubricRepo)
		rubricHandler := NewRubricHandler(rubricService)

		protected.POST("/courses/:id/rubrics", rubricHandler.CreateRubric)
		protected.GET("/courses/:id/rubrics", rubricHandler.ListRubrics)
		protected.GET("/rubrics/:id", rubricHandler.GetRubric)

		protected.POST("/submissions/:id/rubric_grade", rubricHandler.SubmitGrade)
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
