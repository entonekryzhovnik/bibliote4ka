package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ksenia/kakashka/internal/models"
)

// BookRepository реализует интерфейс domain.BookRepository
type BookRepository struct {
	db *sql.DB
}

// NewBookRepository создает новый экземпляр BookRepository
func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Create создает новую книгу в базе данных
func (r *BookRepository) Create(book *models.Book) error {
	query := `
		INSERT INTO books (id, title, author, published, pages, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	book.ID = uuid.New().String()
	book.Status = "available"
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		book.ID,
		book.Title,
		book.Author,
		book.Published,
		book.Pages,
		book.Status,
		book.CreatedAt,
		book.UpdatedAt,
	)

	return err
}

// GetByID получает книгу по ID
func (r *BookRepository) GetByID(id string) (*models.Book, error) {
	query := `
		SELECT id, title, author, published, pages, status, taken_by, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	book := &models.Book{}
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Published,
		&book.Pages,
		&book.Status,
		&book.TakenBy,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("book not found")
	}

	return book, err
}

// List возвращает список книг с учетом фильтров
func (r *BookRepository) List(filters map[string]string) ([]*models.Book, error) {
	query := `
		SELECT id, title, author, published, pages, status, taken_by, created_at, updated_at
		FROM books
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if status, ok := filters["status"]; ok {
		query += " AND status = $" + string(rune('0'+argCount))
		args = append(args, status)
		argCount++
	}

	if author, ok := filters["author"]; ok {
		query += " AND author ILIKE $" + string(rune('0'+argCount))
		args = append(args, "%"+author+"%")
		argCount++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Published,
			&book.Pages,
			&book.Status,
			&book.TakenBy,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

// Update обновляет информацию о книге
	func (r *BookRepository) Update(book *models.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, published = $3, pages = $4, updated_at = $5
		WHERE id = $6
	`

	book.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		book.Title,
		book.Author,
		book.Published,
		book.Pages,
		book.UpdatedAt,
		book.ID,
	)

	return err
}

// Delete удаляет книгу
func (r *BookRepository) Delete(id string) error {
	query := `
		DELETE FROM books
		WHERE id = $1 AND status = 'available'
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("book not found or not available")
	}

	return nil
}

// TakeBook обновляет статус книги на "taken"
func (r *BookRepository) TakeBook(id string, email string) error {
	query := `
		UPDATE books
		SET status = 'taken', taken_by = $1, updated_at = $2
		WHERE id = $3 AND status = 'available'
	`

	result, err := r.db.Exec(query, email, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("book not found or not available")
	}

	return nil
}

// ReturnBook обновляет статус книги на "available"
func (r *BookRepository) ReturnBook(id string) error {
	query := `
		UPDATE books
		SET status = 'available', taken_by = NULL, updated_at = $1
		WHERE id = $2 AND status = 'taken'
	`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("book not found or not taken")
	}

	return nil
} 