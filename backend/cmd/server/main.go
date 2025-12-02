package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/seed"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/worker"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
)

// Entry point that wires configuration, DB, routes.
// Flags:
//
//	--seed : runs checklist seeding and exits
func main() {
	_ = godotenv.Load()

	seedFlag := flag.Bool("seed", false, "seed checklist and exit")
	bootstrapAdmin := flag.Bool("bootstrap-admin", false, "create/update superadmin from ADMIN_EMAIL/ADMIN_PASSWORD and exit")
	flag.Parse()

	cfg := config.MustLoad()
	conn := db.MustOpen(cfg.DatabaseURL)

	pbManager, err := playbook.EnsureActive(conn, cfg.PlaybookPath)
	if err != nil {
		log.Fatal(err)
	}

	if *seedFlag {
		if err := seed.Run(conn); err != nil {
			log.Fatal(err)
		}
		log.Println("Seed completed successfully")
		return
	}

	if *bootstrapAdmin {
		if gen, err := seed.EnsureSuperAdmin(conn, cfg); err != nil {
			log.Fatal(err)
		} else if gen != "" {
			log.Printf("Superadmin '%s' created with password: %s", cfg.AdminEmail, gen)
		} else {
			log.Printf("Superadmin ensured for '%s' (password unchanged)", cfg.AdminEmail)
		}
		return
	}

	// Ensure superadmin exists
	if gen, err := seed.EnsureSuperAdmin(conn, cfg); err != nil {
		log.Printf("superadmin ensure failed: %v", err)
	} else if gen != "" {
		log.Printf("Superadmin '%s' created with password: %s", cfg.AdminEmail, gen)
	}
	if err := seed.EnsureContacts(conn); err != nil {
		log.Printf("contact seed failed: %v", err)
	}

	r := gin.Default()
	// Do not trust any proxy headers by default; see Gin docs.
	_ = r.SetTrustedProxies(nil)

	api := handlers.BuildAPI(r, conn, cfg, pbManager)

	// Initialize S3 cleanup worker if S3 is configured
	s3Client, err := services.NewS3FromEnv()
	if err != nil {
		log.Printf("S3 initialization error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if s3Client != nil && s3Client.Client() != nil {
		cleanupWorker := worker.NewCleanupWorker(conn, s3Client.Client(), s3Client.Bucket())
		go cleanupWorker.Start(ctx)
		log.Println("S3 cleanup worker started")
	} else {
		log.Println("S3 not configured, cleanup worker disabled")
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logging.Info("API listening", "port", cfg.Port)
		if err := api.Run(":" + cfg.Port); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	cancel() // Stop cleanup worker
	log.Println("Server stopped")
}
