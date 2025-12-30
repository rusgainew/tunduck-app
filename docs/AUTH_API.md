# API Аутентификации и Регистрации

## Эндпоинты

### 1. Регистрация пользователя

**POST** `/api/auth/register`

**Тело запроса:**

```json
{
  "username": "ivan_petrov",
  "email": "ivan@example.com",
  "fullName": "Иван Петров",
  "phone": "+996555123456",
  "password": "SecurePassword123",
  "confirmPassword": "SecurePassword123"
}
```

**Валидация:**

- `username` - обязательно, 3-50 символов
- `email` - обязательно, валидный email
- `fullName` - обязательно, 2-100 символов
- `phone` - обязательно, 10-20 символов
- `password` - обязательно, минимум 6 символов
- `confirmPassword` - обязательно, должен совпадать с password

**Успешный ответ (201):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "ivan_petrov",
    "email": "ivan@example.com",
    "fullName": "Иван Петров",
    "phone": "+996555123456",
    "isActive": true
  }
}
```

**Ошибки:**

- `400` - некорректный запрос или валидация не прошла
- `400` - пользователь с таким логином/email уже существует

---

### 2. Вход в систему

**POST** `/api/auth/login`

**Тело запроса:**

```json
{
  "username": "ivan_petrov",
  "password": "SecurePassword123"
}
```

**Успешный ответ (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "ivan_petrov",
    "email": "ivan@example.com",
    "fullName": "Иван Петров",
    "phone": "+996555123456",
    "isActive": true
  }
}
```

**Ошибки:**

- `400` - некорректный запрос
- `401` - неверный логин или пароль
- `401` - учётная запись заблокирована

---

### 3. Получить текущего пользователя

**GET** `/api/auth/me`

**Заголовки:**

```
Authorization: Bearer <token>
```

**Успешный ответ (200):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "ivan_petrov",
  "email": "ivan@example.com",
  "fullName": "Иван Петров"
}
```

**Ошибки:**

- `401` - токен не предоставлен
- `401` - невалидный токен

---

## Структура базы данных

### Таблица `users`

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  full_name VARCHAR(100) NOT NULL,
  phone VARCHAR(20) NOT NULL,
  password VARCHAR(255) NOT NULL, -- bcrypt hash
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

---

## Безопасность

### Хеширование паролей

- Используется `bcrypt` с дефолтной сложностью (cost=10)
- Пароли никогда не возвращаются в JSON ответах (`json:"-"`)

### JWT токены

- Срок действия: 7 дней
- Секрет берётся из `JWT_SECRET` переменной окружения
- Claims содержат: `user_id`, `username`, `email`, `full_name`, `exp`

### Валидация

- Используется библиотека `validator/v10`
- Проверка уникальности username/email на уровне базы данных
- Проверка совпадения паролей перед регистрацией

---

## Примеры использования

### cURL

**Регистрация:**

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "email": "test@example.com",
    "fullName": "Test User",
    "phone": "+996555000000",
    "password": "password123",
    "confirmPassword": "password123"
  }'
```

**Логин:**

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_user",
    "password": "password123"
  }'
```

**Получить профиль:**

```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## Frontend интеграция

Обновите `ui-admin/lib/api.ts`:

```typescript
export interface RegisterData {
  username: string;
  email: string;
  fullName: string;
  phone: string;
  password: string;
  confirmPassword: string;
}

export interface LoginData {
  username: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    username: string;
    email: string;
    fullName: string;
    phone: string;
    isActive: boolean;
  };
}

export const authApi = {
  register: (data: RegisterData): Promise<AuthResponse> =>
    api.post("/api/auth/register", data),

  login: (data: LoginData): Promise<AuthResponse> =>
    api.post("/api/auth/login", data),

  getCurrentUser: (token: string) =>
    api.get("/api/auth/me", {
      headers: { Authorization: `Bearer ${token}` },
    }),
};
```
