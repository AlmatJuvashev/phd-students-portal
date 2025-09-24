package main

import (
	"phd-portal/backend/internal/logging"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"phd-portal/backend/internal/config"
	"phd-portal/backend/internal/db"
	"phd-portal/backend/internal/handlers"
	"phd-portal/backend/internal/seed"

	"github.com/gin-gonic/gin"
)

// Entry point that wires configuration, DB, routes.
// Flags:
//   --seed : runs checklist seeding and exits
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
