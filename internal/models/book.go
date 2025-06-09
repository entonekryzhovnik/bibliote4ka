package models

import "time"

// Book представляет собой модель книги в библиотеке
type Book struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Published int       `json:"published"`
	Pages     int       `json:"pages"`
	Status    string    `json:"status"`
	TakenBy   *string   `json:"taken_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookRepository определяет интерфейс для работы с книгами в базе данных
type BookRepository interface {
	Create(book *Book) error
	GetByID(id string) (*Book, error)
	List(filters map[string]string) ([]*Book, error)
	Update(book *Book) error
	Delete(id string) error
	TakeBook(id string, email string) error
	ReturnBook(id string) error
}
