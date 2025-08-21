# Analytics Service

Высокопроизводительный Go-сервис для аналитики данных с параллельной обработкой, аутентификацией и REST API.

## Требования

* Go 1.21+
* PowerShell (для Windows)
* Git

## Установка и запуск

### Локальная разработка

1. Клонировать репозиторий:

```powershell
git clone <repository-url>
cd analytics-service
```

2. Установить зависимости:

```powershell
go mod tidy
```

3. Запустить тесты (сначала собери бинарник!!!):

```powershell
go test -v -cover ./...
```

4. Собрать и запустить сервис:

```powershell
go build -o analytics-service.exe ./cmd/server
.\analytics-service.exe
```

### Тестирование API

Для проверки работы всех эндпоинтов можно использовать PowerShell скрипт:

```powershell
powershell -ExecutionPolicy Bypass -File scripts\test_api.ps1
```

## API Endpoints

### 1. Health Check

```
GET /
```

Ответ:

```json
{
  "message": "Analytics Service is running"
}
```

### 2. Аутентификация

```
POST /auth
Content-Type: application/json
```

Пример тела запроса:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Ответ:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Валидация токена

```
GET /validate?token=<token>
```

Ответ:

```json
{
  "valid": true
}
```

### 4. Аналитика

```
POST /analytics
Content-Type: application/json
```

Пример тела запроса:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "StartDate": "01.01.2024",
  "FinishDate": "31.01.2024"
}
```

Пример ответа:

```json
{
  "items": [
    {
      "Name": "Товар 1",
      "Code": "1001",
      "Group": "Группа 1",
      "Sales": 500,
      "Loss": 500,
      "LossOfProfit": 100,
      "OSA": 100,
      "ABC": "A"
    }
  ],
  "total": 3
}
```

## Тестирование

### Unit тесты

```powershell
go test -v ./...
```

### Тесты с покрытием

```powershell
go test -v -cover ./...
```

### Интеграционные тесты API

```powershell
powershell -ExecutionPolicy Bypass -File scripts\test_api.ps1
```

## Конфигурация

### Переменные окружения

| Переменная        | Описание               | По умолчанию |
| ----------------- | ---------------------- | ------------ |
| PORT              | Порт сервера           | 8080         |
| AUTH\_SECRET\_KEY | Секретный ключ для JWT | secret       |
| WORKERS           | Количество worker'ов   | 4            |

Пример `.env` файла:

```env
PORT=8080
AUTH_SECRET_KEY=your-secret-key
WORKERS=8
```

## Структура проекта

```
analytics-service/
├── cmd/server/main.go          # Точка входа
├── internal/
│   ├── analytics/              # Бизнес-логика аналитики
│   ├── auth/                   # Аутентификация
│   ├── config/                 # Конфигурация
│   ├── handlers/               # HTTP обработчики
│   └── userdb/                 # Хранение токенов
├── routes/                     # Данные (LogPas.txt, *.json)
├── examples/                   # Примеры запросов
├── scripts/                    # Скрипты для тестирования
├── Makefile                    # Автоматизация сборки и тестов
├── go.mod                      # Зависимости
├── go.sum                      # Контроль версий зависимостей
└── README.md                   # Документация
```
