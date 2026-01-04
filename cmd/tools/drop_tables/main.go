package main

import (
	"log"

	"github.com/qsrrvln/mountapi/internal/config"
	"github.com/qsrrvln/mountapi/internal/database"
	"github.com/qsrrvln/mountapi/internal/models"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Dropping tables...")
	err = db.Migrator().DropTable(&models.Post{}, &models.Route{}, &models.Mountain{})
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
	log.Println("Tables dropped successfully.")
}
