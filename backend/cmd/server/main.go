package main

import (
	"flag"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/logging"
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/seed"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

// Entry point that wires configuration, DB, routes.
// Flags:
//
//	--seed : runs checklist seeding and exits
func main() {
	_ = godotenv.Load()

	seedFlag := flag.Bool("seed", false, "seed checklist and exit")
	flag.Parse()

	cfg := config.MustLoad()
	conn := db.MustOpen(cfg.DatabaseURL)

	if *seedFlag {
		if err := seed.Run(conn); err != nil {
			log.Fatal(err)
		}
		log.Println("Seed completed successfully")
		return
	}

	r := gin.Default()

	api := handlers.BuildAPI(r, conn, cfg)

	logging.Info("API listening", "port", cfg.Port)
	if err := api.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
