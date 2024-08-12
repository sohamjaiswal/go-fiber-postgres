package book

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func (h *Handler) CreateBook(c *fiber.Ctx) error {
	book := CreateBookDTO{}
	if err := c.BodyParser(book); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"message": "request failed",
		})
	}

	if err := h.Repo.CreateBook(&book); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not create book",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book added",
	})
}

func (h *Handler) GetBooks(c *fiber.Ctx) error {
	books, err := h.Repo.GetBooks()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "could not get books",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    books,
	})
}

func (h *Handler) GetBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}

	book, err := h.Repo.GetBookByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "could not get the book",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book fetched successfully",
		"data": book,
	})
}

func (h *Handler) DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
	}

	if err := h.Repo.DeleteBook(id); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted successfully",
	})
}
