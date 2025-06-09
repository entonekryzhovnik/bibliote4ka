.PHONY: run build up down logs test

# 🚀 Запуск приложения
run:
	go run cmd/api/main.go

# 🔨 Сборка приложения
build:
	go build -o bin/api cmd/api/main.go

# 🐳 Запуск базы данных
up:
	docker-compose up -d

# 🛑 Остановка базы данных
down:
	docker-compose down

# 📝 Просмотр логов
logs:
	docker-compose logs -f

# 🧪 Тестирование API
test:
	@echo "🧪 Тестирование API..."
	@echo "1. Создание книги..."
	curl -s -X POST http://localhost:8080/books \
		-H "Content-Type: application/json" \
		-H "X-Admin-Secret: admin-secret" \
		-d '{"title":"Война и мир","author":"Лев Толстой","published":1869,"pages":1225}' | jq
	@echo "\n2. Получение списка книг..."
	curl -s http://localhost:8080/books | jq 