package book

import (
	"github.com/sohamjaiswal/go-fiber-postgres/server/pkg/models"
	"gorm.io/gorm"
)

type CreateBookDTO struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateBook(createBookDto *CreateBookDTO) error {
	book := &models.Book{
		Author: &createBookDto.Author,
		Title: &createBookDto.Title,
		Publisher: &createBookDto.Publisher,
	}
	return r.DB.Create(&book).Error
}

func (r *Repository) GetBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.DB.Find(&books).Error
	return books, err
}

func (r *Repository) GetBookByID(id string) (models.Book, error) {
	var book models.Book
	err := r.DB.Where("id = ?", id).First(&book).Error
	return book, err
}

func (r *Repository) DeleteBook(id string) error {
	return r.DB.Delete(&models.Book{}, id).Error
}
