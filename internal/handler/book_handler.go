package handler

import (
	"net/http"

	"github.com/ksenia/kakashka/internal/models"
	"github.com/ksenia/kakashka/internal/service"
	"github.com/labstack/echo/v4"
)

// CustomError представляет собой кастомную ошибку
type CustomError struct {
	Message string
	Code    int
}

func (e *CustomError) Error() string {
	return e.Message
}

// BookHandler представляет собой обработчик HTTP-запросов для работы с книгами
type BookHandler struct {
	service *service.BookService
}

// NewBookHandler создает новый экземпляр BookHandler
func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

// RegisterRoutes регистрирует маршруты для работы с книгами
func (h *BookHandler) RegisterRoutes(e *echo.Echo) {
	// Настраиваем обработчик ошибок
	e.HTTPErrorHandler = h.errorHandler

	// Публичные маршруты
	e.GET("/books", h.List)
	e.GET("/books/:id", h.GetByID)
	e.POST("/books/:id/take", h.TakeBook)
	e.POST("/books/:id/return", h.ReturnBook)

	// Админские маршруты
	admin := e.Group("/books")
	admin.Use(h.adminMiddleware)
	admin.POST("", h.Create)
	admin.PUT("/:id", h.Update)
	admin.DELETE("/:id", h.Delete)
}

// errorHandler обрабатывает ошибки
func (h *BookHandler) errorHandler(err error, c echo.Context) {
	if customErr, ok := err.(*CustomError); ok {
		c.JSON(customErr.Code, echo.Map{"error": customErr.Message})
		return
	}

	// Обработка стандартных ошибок
	switch err.Error() {
	case "book not found":
		c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	case "invalid email format":
		c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal server error"})
	}
}

// List возвращает список книг
func (h *BookHandler) List(c echo.Context) error {
	filters := make(map[string]string)
	if status := c.QueryParam("status"); status != "" {
		filters["status"] = status
	}
	if author := c.QueryParam("author"); author != "" {
		filters["author"] = author
	}

	books, err := h.service.List(filters)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, books)
}

// GetByID возвращает книгу по ID
func (h *BookHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	book, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, book)
}

// Create создает новую книгу
func (h *BookHandler) Create(c echo.Context) error {
	var book models.Book
	if err := c.Bind(&book); err != nil {
		return &CustomError{Message: err.Error(), Code: http.StatusBadRequest}
	}

	if err := h.service.Create(&book); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, book)
}

// Update обновляет информацию о книге
func (h *BookHandler) Update(c echo.Context) error {
	id := c.Param("id")

	// Получаем текущую книгу
	book, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	// Обновляем только разрешенные поля
	var updateData struct {
		Title     string `json:"title"`
		Author    string `json:"author"`
		Published int    `json:"published"`
		Pages     int    `json:"pages"`
	}
	if err := c.Bind(&updateData); err != nil {
		return &CustomError{Message: err.Error(), Code: http.StatusBadRequest}
	}

	// Обновляем только разрешенные поля
	book.Title = updateData.Title
	book.Author = updateData.Author
	book.Published = updateData.Published
	book.Pages = updateData.Pages

	if err := h.service.Update(book); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, book)
}

// Delete удаляет книгу
func (h *BookHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	// Проверяем статус книги перед удалением
	book, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	if book.Status != "available" {
		return &CustomError{
			Message: "Cannot delete book: it is not available",
			Code:    http.StatusBadRequest,
		}
	}

	if err := h.service.Delete(id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// TakeBook позволяет взять книгу
func (h *BookHandler) TakeBook(c echo.Context) error {
	id := c.Param("id")

	var request struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&request); err != nil {
		return &CustomError{Message: err.Error(), Code: http.StatusBadRequest}
	}

	if err := h.service.TakeBook(id, request.Email); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// ReturnBook позволяет вернуть книгу
func (h *BookHandler) ReturnBook(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.ReturnBook(id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// adminMiddleware проверяет, является ли пользователь админом
func (h *BookHandler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-Admin-Secret") != "admin-secret" {
			return &CustomError{Message: "Unauthorized", Code: http.StatusUnauthorized}
		}
		return next(c)
	}
}
