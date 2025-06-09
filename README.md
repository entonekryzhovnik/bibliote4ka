# �� Онлайн-библиотечка 🧸

Микросервисик для управления библиотекой с админским доступом и middleware для авторизации ✨

## 🛠️ Что нам понадобится 🧸

- Go 1.21 или новее
- Docker и Docker Compose
- PostgreSQL 15 (будет запущен в Docker)

## 🚀 Начинаем работу 🧸

1. Клонируем репозиторий:
```bash
git clone <repository-url>
cd <repository-name>
```

2. Запускаем базу данных:
```bash
make up
```

3. Запускаем приложение:
```bash
make run
```

Серверчик запустится на порту 8080 🎉 🧸

## 📡 API Endpoints 🧸

### 🌐 Публичные эндпоинты

- `GET /books` - список книжечек (с фильтрами) 📚
- `GET /books/{id}` - получить книжечку по ID 📖
- `POST /books/{id}/take` - взять книжечку 📥
- `POST /books/{id}/return` - вернуть книжечку 📤

### 🔐 Админские эндпоинты

Все админские эндпоинты требуют заголовок `X-Admin-Secret` 🎫

- `POST /books` - добавить книжечку ➕
- `PUT /books/{id}` - обновить книжечку ✏️
- `DELETE /books/{id}` - удалить книжечку 🗑️

## ⚙️ Переменные окружения 🧸

- `DB_HOST` - хостик базы данных (по умолчанию: localhost) 🏠
- `DB_PORT` - портик базы данных (по умолчанию: 5433) 🔌
- `DB_USER` - пользователь базы данных (по умолчанию: postgres) 👤
- `DB_PASSWORD` - пароль базы данных (по умолчанию: postgres) 🔑
- `DB_NAME` - имя базы данных (по умолчанию: library) 📚
- `ADMIN_SECRET` - секретный ключик для админского доступа (по умолчанию: admin-secret) 🎫
- `SERVER_PORT` - портик сервера (по умолчанию: 8080) 🌐

## 🛠️ Разработка 🧸

### 🐳 База данных

База данных запускается в Docker контейнере 🐳

Запускаем:
```bash
make up
```

Останавливаем:
```bash
make down
```

Смотрим логики:
```bash
make logs
```

Сбрасываем базу:
```bash
make down -v
make up
```

### 🧪 Тестирование

Запускаем тестики:
```bash
make test
```

## 📝 Примерчики запросиков 🧸

### Создаем книжечку (админ) 📚
```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -H "X-Admin-Secret: admin-secret" \
  -d '{
    "title": "Война и мир",
    "author": "Лев Толстой",
    "published": 1869,
    "pages": 1225
  }'
```

### Получаем список книжечек 📚
```bash
curl http://localhost:8080/books
```

### Берем книжечку 📚
```bash
curl -X POST http://localhost:8080/books/{id}/take \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }' 