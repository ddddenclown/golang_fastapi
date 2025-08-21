# Analytics Service — README (код и только код)

Кратко и по делу — всё необходимое для запуска, тестирования и сборки сервиса.

---

## Описание

Go-сервис с REST API для:

* генерации/валидации JWT-токенов (`/auth`, `/validate`);
* вычислительной аналитики (`/analytics`) с параллельной обработкой данных.

---

## Требования

* Go 1.21+

---

## Быстрый старт (в корне проекта)

```bash
# перейти в корень проекта (там, где go.mod)
cd path/to/analytics-service

# установить зависимости
go mod download

# собрать бинарник
go build -o analytics-service ./cmd/server

# или сразу запустить
go run ./cmd/server
```

Windows (PowerShell):

```powershell
cd C:\path\to\analytics-service
go mod download
go build -o analytics-service.exe ./cmd/server
.\analytics-service.exe
```

Сервер по умолчанию слушает порт из переменной `PORT` (по умолчанию `8080`).

---

## Переменные окружения / .env

Поддерживаются (минимально):

```
PORT=8080
AUTH_SECRET_KEY=secret
WORKERS=4
```

Можно положить `.env` в корень (или экспортировать переменные в окружение перед запуском).

---

## API (контракт)

### Health

```
GET /
```

Ответ:

```json
{ "message": "Analytics Service is running" }
```

### Генерация токена

```
POST /auth
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Ответ:

```json
{ "token": "<jwt_token>" }
```

### Валидация токена

```
GET /validate?token=<jwt_token>
```

Ответ:

```json
{ "valid": true }
```

### Аналитика

```
POST /analytics
Content-Type: application/json

{
  "Token": "<jwt_token>",
  "StartDate": "01.01.2024",
  "FinishDate": "31.01.2024"
}
```

* **Примечание:** поле `Token`/регистрозависимость должно соответствовать тому, что ожидает реализация (см. код).
* Тело запроса можно подать из `examples/sample.json`.

Пример ответа:

```json
{
  "items": [ /*...*/ ],
  "total": 3
}
```

---

## Примеры запросов

Bash (Linux / Git Bash / WSL):

```bash
# POST /auth
curl -s -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# GET /validate
curl -s "http://localhost:8080/validate?token=YOUR_TOKEN"

# POST /analytics (from file)
curl -s -X POST http://localhost:8080/analytics \
  -H "Content-Type: application/json" \
  -d @examples/sample.json
```

PowerShell:

```powershell
# POST /auth
$auth = Invoke-RestMethod -Method Post -Uri http://localhost:8080/auth `
  -Body '{"email":"test@example.com","password":"password123"}' -ContentType 'application/json'
$token = $auth.token

# GET /validate
Invoke-RestMethod -Uri "http://localhost:8080/validate?token=$token"

# POST /analytics
Invoke-RestMethod -Method Post -Uri http://localhost:8080/analytics `
  -Body (Get-Content .\examples\sample.json -Raw) -ContentType 'application/json'
```

---

## Тесты

Запуск всех тестов (в корне проекта):

```bash
go test ./... -v
```

Покрытие:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Быстрый бенч (простая проверка)

Простой вариант (PowerShell) — замер N запросов к `/analytics`:

```powershell
$N=5; $body = Get-Content .\examples\sample.json -Raw
$time = Measure-Command { for ($i=0;$i -lt $N;$i++) { Invoke-RestMethod -Method Post -Uri http://localhost:8080/analytics -Body $body -ContentType 'application/json' | Out-Null } }
"Avg ms: $([math]::Round($time.TotalMilliseconds/$N,2))"
```

---

## Что отправлять проверяющему (минимум)

* `cmd/` и `internal/` — исходники Go
* `go.mod`, `go.sum`
* `examples/sample.json`
* `routes/` с `stock_dump.json`, `sales_dump.json`, `LogPass.txt`
* `scripts/test_api.ps1` (PowerShell) и/или `scripts/test_api.sh` (bash) — опционально
* `README.md`

> Docker/Kubernetes/monitoring/CI не обязательны — приложи только код.

---

## Структура (минимально важная)

```
cmd/server/main.go
internal/
  ├─ auth/
  ├─ analytics/
  ├─ config/
  ├─ handlers/
  └─ userdb/
examples/sample.json
routes/stock_dump.json
routes/sales_dump.json
routes/LogPass.txt
go.mod go.sum
scripts/test_api.ps1
```

---

## Отладка и полезные трюки

* Если видишь в PowerShell кривую кириллицу — визуальная проблема консоли; JSON UTF-8. Для копирования токена:

```powershell
$auth.token | Out-File token.txt -Encoding utf8
```

* Для быстрого перезапуска после правок: `go build` → запустить новый бинарник.
* Если endpoint возвращает пустые результаты — проверь:

  * токен валидный;
  * `examples/sample.json` имеет правильный регистр полей;
  * данные в `routes/` покрывают период дат из запроса.
