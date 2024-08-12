package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sohamjaiswal/go-fiber-postgres/server/models"
	"github.com/sohamjaiswal/go-fiber-postgres/server/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}

	if err := c.BodyParser(&book); err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message":"request failed",
			},
		)
		return err
	}

	if err := r.DB.Create(&book).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message":"could not create book",
			},
		)
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":"book added",
		},
	)
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	if err := r.DB.Find(bookModels).Error; err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message":"could not get books",
			},
		)
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books fetched succesfully",
			"data": bookModels,
		},
	)

	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "id cannot be empty",
			},
		)
		return errors.New("cannot delete book without id")
	}
	if err := r.DB.Delete(bookModel, id).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "could not delete book",
			},
		)
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book deleted successfully",
		},
	)
	return nil
}

func (r *Repository) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModel := &models.Book{}
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id cannot be empty",
			},
		)
		return errors.New("cannot get book without id")
	}
	if err := r.DB.Where("id = ?", id).First(bookModel).Error; err != nil {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{
				"message": "could not get the book",
			},
		)
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book fetched successgully",
			"data": bookModel,
		},
	)
	return nil
}

func (r *Repository) SetupRoutes (app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create-book", r.CreateBook)
	api.Delete("/delete-book/:id", r.DeleteBook)
	api.Get("/get-book/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	mode := os.Getenv("MODE")	
	hostAddress := os.Getenv("HOST") + ":" + os.Getenv("PORT")

	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatalf("err occured while setting up db conn: %v", err)
	}

	if err = models.MigrateBooks(db); err != nil {
		log.Fatal("could not migrate books db")
	}

	r := Repository{
		DB: db,
	}

	log.Printf("Running on %v in %v mode", hostAddress, mode)

	app := fiber.New()
	r.SetupRoutes(app)
	log.Fatal((app.Listen(hostAddress)))
}