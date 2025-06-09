package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/ksenia/kakashka/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// BookService представляет собой сервисный слой для работы с книгами
type BookService struct {
	repo models.BookRepository
}

// NewBookService создает новый экземпляр BookService
func NewBookService(repo models.BookRepository) *BookService {
	return &BookService{repo: repo}
}

// Create создает новую книгу
func (s *BookService) Create(book *models.Book) error {
	if err := s.validateBook(book); err != nil {
		return err
	}
	return s.repo.Create(book)
}

// GetByID получает книгу по ID
func (s *BookService) GetByID(id string) (*models.Book, error) {
	return s.repo.GetByID(id)
}

// List возвращает список книг с учетом фильтров
func (s *BookService) List(filters map[string]string) ([]*models.Book, error) {
	return s.repo.List(filters)
}

// Update обновляет информацию о книге
func (s *BookService) Update(book *models.Book) error {
	if err := s.validateBook(book); err != nil {
		return err
	}
	return s.repo.Update(book)
}

// Delete удаляет книгу
func (s *BookService) Delete(id string) error {
	return s.repo.Delete(id)
}

// TakeBook позволяет взять книгу
func (s *BookService) TakeBook(id string, email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return s.repo.TakeBook(id, email)
}

// ReturnBook позволяет вернуть книгу
func (s *BookService) ReturnBook(id string) error {
	return s.repo.ReturnBook(id)
}

// validateBook проверяет валидность данных книги
func (s *BookService) validateBook(book *models.Book) error {
	if book.Title == "" {
		return errors.New("title is required")
	}
	if book.Author == "" {
		return errors.New("author is required")
	}
	if book.Published <= 0 {
		return errors.New("published year must be positive")
	}
	if book.Published > time.Now().Year() {
		return errors.New("published year cannot be in the future")
	}
	if book.Pages <= 0 {
		return errors.New("pages must be positive")
	}
	return nil
} 