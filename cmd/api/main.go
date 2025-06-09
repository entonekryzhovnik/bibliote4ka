package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ksenia/kakashka/internal/config"
	"github.com/ksenia/kakashka/internal/handler"
	"github.com/ksenia/kakashka/internal/repository/postgres"
	"github.com/ksenia/kakashka/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.New()

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение к базе данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализируем слои приложения
	bookRepo := postgres.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	// Создаем экземпляр Echo
	e := echo.New()

	// Добавляем middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрируем маршруты
	bookHandler.RegisterRoutes(e)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server is starting on port %s", cfg.ServerPort)
		if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Даем 10 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
}
