package book

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	repo := NewRepository(db)
	bookHandler := NewHandler(repo)
	api := app.Group("/api")
	api.Post("/create-book", bookHandler.CreateBook)
	api.Delete("/delete-book/:id", bookHandler.DeleteBook)
	api.Get("/get-book/:id", bookHandler.GetBookByID)
	api.Get("/books", bookHandler.GetBooks)
}