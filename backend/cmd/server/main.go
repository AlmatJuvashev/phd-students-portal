package main

import (
	"flag"
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/seed"
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

	r := gin.Default()

	// Trust only localhost proxies (fix Gin warning)
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Printf("Warning: failed to set trusted proxies: %v", err)
	}

	api := handlers.BuildAPI(r, conn, cfg, pbManager)

	logging.Info("API listening", "port", cfg.Port)
	if err := api.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
