# 📚 Полная документация Swagger API

## 🌍 Доступ к документации

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON документация**: http://localhost:8080/swagger/doc.json
- **Редирект**: http://localhost:8080/docs

## 📋 Группы эндпоинтов

### 1️⃣ Authentication (Аутентификация)

#### POST `/api/auth/register`

**Регистрация нового пользователя**

**Request:**

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password_123"
}
```

**Response (201):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2025-12-29T05:34:02Z"
}
```

**Ошибки:**

- `400` - Неверный ввод
- `429` - Превышен лимит запросов (10 за 24 часа)

---

#### POST `/api/auth/login`

**Вход пользователя**

**Request:**

```json
{
  "email": "john@example.com",
  "password": "secure_password_123"
}
```

**Response (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2025-12-29T05:34:02Z"
  }
}
```

**Ошибки:**

- `401` - Неверные учётные данные
- `429` - Превышен лимит запросов (20 за час)

---

#### POST `/api/auth/logout`

**Выход пользователя**

**Headers:**

```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**

```json
{
  "message": "Logout successful"
}
```

**Ошибки:**

- `401` - Отсутствует или неверный токен

---

### 2️⃣ Users (Пользователи)

#### GET `/api/users`

**Получить список всех пользователей**

**Query Parameters:**

- `page` (int, default: 1) - Номер страницы
- `limit` (int, default: 10) - Количество записей

**Headers:**

```
Authorization: Bearer <JWT_TOKEN>
```

**Response (200):**

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "created_at": "2025-12-29T05:34:02Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_records": 50
  }
}
```

---

#### GET `/api/users/{id}`

**Получить информацию о конкретном пользователе**

**Path Parameters:**

- `id` (UUID) - ID пользователя

**Response (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2025-12-29T05:34:02Z"
}
```

**Ошибки:**

- `404` - Пользователь не найден

---

#### POST `/api/users`

**Создать нового пользователя**

**Request:**

```json
{
  "username": "jane_doe",
  "email": "jane@example.com",
  "password": "secure_password_456"
}
```

**Response (201):**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "username": "jane_doe",
  "email": "jane@example.com",
  "created_at": "2025-12-29T05:35:00Z"
}
```

---

#### PUT `/api/users/{id}`

**Обновить пользователя**

**Path Parameters:**

- `id` (UUID) - ID пользователя

**Request:**

```json
{
  "username": "jane_updated",
  "email": "jane.updated@example.com"
}
```

**Response (200):**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "username": "jane_updated",
  "email": "jane.updated@example.com",
  "created_at": "2025-12-29T05:35:00Z"
}
```

---

#### DELETE `/api/users/{id}`

**Удалить пользователя**

**Path Parameters:**

- `id` (UUID) - ID пользователя

**Response (204):**

```
No Content
```

---

### 3️⃣ Organizations (Организации)

#### GET `/api/organizations`

**Получить список организаций**

**Query Parameters:**

- `page` (int, default: 1) - Номер страницы
- `limit` (int, default: 10) - Количество записей
- `search` (string) - Поиск по названию

**Response (200):**

```json
{
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "ООО Компания",
      "inn": "7728168971",
      "kpp": "772801001",
      "director": "Иван Петров",
      "created_at": "2025-12-29T05:34:02Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 3,
    "total_records": 25
  }
}
```

---

#### GET `/api/organizations/{id}`

**Получить информацию об организации**

**Response (200):**

```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "name": "ООО Компания",
  "inn": "7728168971",
  "kpp": "772801001",
  "director": "Иван Петров",
  "created_at": "2025-12-29T05:34:02Z"
}
```

---

#### POST `/api/organizations`

**Создать новую организацию**

**Request:**

```json
{
  "name": "ООО Новая Компания",
  "inn": "1234567890",
  "kpp": "123456789",
  "director": "Петр Иванов"
}
```

**Response (201):**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440003",
  "name": "ООО Новая Компания",
  "inn": "1234567890",
  "kpp": "123456789",
  "director": "Петр Иванов",
  "created_at": "2025-12-29T05:36:00Z"
}
```

---

#### PUT `/api/organizations/{id}`

**Обновить организацию**

**Request:**

```json
{
  "name": "ООО Обновленная Компания",
  "director": "Сергей Сидоров"
}
```

**Response (200):**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440003",
  "name": "ООО Обновленная Компания",
  "inn": "1234567890",
  "kpp": "123456789",
  "director": "Сергей Сидоров",
  "created_at": "2025-12-29T05:36:00Z"
}
```

---

#### DELETE `/api/organizations/{id}`

**Удалить организацию**

**Response (204):**

```
No Content
```

---

### 4️⃣ Documents (ЭСФ Документы)

#### GET `/api/documents`

**Получить список документов**

**Query Parameters:**

- `page` (int) - Номер страницы
- `limit` (int) - Количество записей
- `status` (string) - Фильтр по статусу: `draft|sent|received|processed`
- `organization_id` (UUID) - Фильтр по организации

**Response (200):**

```json
{
  "data": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440004",
      "number": "ЭСФ-001-2025",
      "status": "draft",
      "amount": 50000.0,
      "organization_id": "770e8400-e29b-41d4-a716-446655440002",
      "created_at": "2025-12-29T05:34:02Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 2,
    "total_records": 15
  }
}
```

---

#### GET `/api/documents/{id}`

**Получить документ по ID**

**Response (200):**

```json
{
  "id": "990e8400-e29b-41d4-a716-446655440004",
  "number": "ЭСФ-001-2025",
  "status": "draft",
  "amount": 50000.0,
  "organization_id": "770e8400-e29b-41d4-a716-446655440002",
  "created_at": "2025-12-29T05:34:02Z"
}
```

---

#### POST `/api/documents`

**Создать новый документ**

**Request:**

```json
{
  "number": "ЭСФ-002-2025",
  "status": "draft",
  "amount": 75000.0,
  "organization_id": "770e8400-e29b-41d4-a716-446655440002"
}
```

**Response (201):**

```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440005",
  "number": "ЭСФ-002-2025",
  "status": "draft",
  "amount": 75000.0,
  "organization_id": "770e8400-e29b-41d4-a716-446655440002",
  "created_at": "2025-12-29T05:37:00Z"
}
```

---

#### PUT `/api/documents/{id}`

**Обновить документ**

**Request:**

```json
{
  "status": "sent",
  "amount": 75500.0
}
```

**Response (200):**

```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440005",
  "number": "ЭСФ-002-2025",
  "status": "sent",
  "amount": 75500.0,
  "organization_id": "770e8400-e29b-41d4-a716-446655440002",
  "created_at": "2025-12-29T05:37:00Z"
}
```

---

#### DELETE `/api/documents/{id}`

**Удалить документ**

**Response (204):**

```
No Content
```

---

### 5️⃣ System (Система)

#### GET `/health`

**Проверка здоровья системы**

**Response (200):**

```json
{
  "status": "UP",
  "timestamp": "2025-12-29T05:40:00Z",
  "uptime": "1h 46m 30s",
  "components": [
    {
      "name": "PostgreSQL",
      "status": "UP",
      "response_time": "9.397ms",
      "message": "Database connected successfully",
      "last_checked": "2025-12-29T05:40:00Z"
    },
    {
      "name": "Redis",
      "status": "UP",
      "response_time": "186µs",
      "message": "Redis connected successfully",
      "last_checked": "2025-12-29T05:40:00Z"
    }
  ]
}
```

---

#### GET `/metrics`

**Prometheus метрики**

**Response (200):**

```
# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/users",status="200"} 145
http_requests_total{method="POST",path="/api/auth/login",status="200"} 89
...
```

---

## 🔐 Аутентификация

Все защищённые эндпоинты требуют JWT токен в заголовке:

```
Authorization: Bearer <JWT_TOKEN>
```

JWT токен получается при логине и содержит:

- User ID
- Email
- Username
- Время выпуска (iat)
- Время истечения (exp) - по умолчанию 24 часа

---

## 📊 Коды ответов

| Код | Описание                                    |
| --- | ------------------------------------------- |
| 200 | OK - Успешный запрос                        |
| 201 | Created - Ресурс создан                     |
| 204 | No Content - Успешно удалено                |
| 400 | Bad Request - Неверный запрос               |
| 401 | Unauthorized - Требуется аутентификация     |
| 403 | Forbidden - Доступ запрещён                 |
| 404 | Not Found - Ресурс не найден                |
| 429 | Too Many Requests - Превышен лимит запросов |
| 500 | Internal Server Error - Ошибка сервера      |
| 503 | Service Unavailable - Сервис недоступен     |

---

## ⚡ Rate Limiting

**Публичные эндпоинты (регистрация, логин):**

- 10 запросов за 24 часа на IP

**Аутентифицированные эндпоинты:**

- 1000 запросов за час на пользователя

**Health и Metrics:**

- 100 запросов за час

При превышении лимита возвращается ответ 429:

```json
{
  "error": "RATE_LIMIT_EXCEEDED",
  "message": "Too many requests",
  "reset": 1704067200
}
```

---

## 🚀 Быстрый старт тестирования

### 1. Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test12345!"
  }'
```

### 2. Логин

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test12345!"
  }'
```

### 3. Получить токен и использовать его

```bash
TOKEN="<token_from_login_response>"

curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Создать организацию

```bash
curl -X POST http://localhost:8080/api/organizations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "ООО Тестовая Компания",
    "inn": "7728168971",
    "kpp": "772801001",
    "director": "Тестовый Директор"
  }'
```

### 5. Создать документ

```bash
curl -X POST http://localhost:8080/api/documents \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "number": "ЭСФ-001-2025",
    "status": "draft",
    "amount": 100000,
    "organization_id": "<organization_id>"
  }'
```

---

## 📝 Примечания

- Все временные метки в UTC формате (ISO 8601)
- IDs - это UUID v4
- Все денежные суммы в основных единицах (рубли)
- Поиск чувствителен к регистру в некоторых полях
- Пагинация начинается с 1 (не с 0)

---

**Документация обновлена: 29 декабря 2025**
