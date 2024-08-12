package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sohamjaiswal/go-fiber-postgres/server/pkg/api/book"
	"github.com/sohamjaiswal/go-fiber-postgres/server/pkg/models"
	"github.com/sohamjaiswal/go-fiber-postgres/server/pkg/storage"
)

func main() {
	mode := os.Getenv("MODE")	
	hostAddress := os.Getenv("HOST") + ":" + os.Getenv("PORT")

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalf("Error occurred while setting up DB connection: %v", err)
	}

	if err = models.MigrateBooks(db); err != nil {
		log.Fatal("Could not migrate books DB")
	}

	app := fiber.New()
	book.RegisterRoutes(app, db)

	log.Printf("Running on %v in %v mode", hostAddress, mode)
	log.Fatal(app.Listen(hostAddress))
}
