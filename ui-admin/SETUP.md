# Tunduck Admin Panel - Инструкция по запуску

## Созданная функциональность

### ✅ Реализовано

1. **Аутентификация**

   - Страница входа с валидацией (/login)
   - Страница регистрации (/register)
   - JWT аутентификация
   - Автоматический выход при истечении токена
   - Защита роутов

2. **Дашборд** (/dashboard)

   - Статистика: пользователи, организации, документы
   - Последние пользователи
   - Список организаций
   - Карточки с быстрым доступом

3. **Управление пользователями** (/dashboard/users)

   - Таблица со всеми пользователями
   - Пагинация (10 на странице)
   - Поиск по имени, email, username
   - Отображение роли и статуса
   - Сортировка по дате создания

4. **Управление организациями** (/dashboard/organizations)

   - Карточки организаций в grid layout
   - Поиск по названию, описанию, БД
   - Удаление организаций
   - Отображение основной информации

5. **Управление документами** (/dashboard/documents)

   - Карточки документов ЭСФ
   - Поиск по ИНН, названию, комментарию
   - Фильтр по резидентности
   - Отображение суммы и валюты

6. **Инфраструктура**
   - TypeScript типы для всех сущностей
   - API клиент с перехватчиками
   - Zustand для управления состоянием
   - React Query для кеширования данных
   - Адаптивный дизайн (mobile-first)

## Быстрый старт

### 1. Запустите backend API

```bash
# В корне проекта
cd /home/rusgai_max/go/src/github.com/rusgainew/tunduck-app
docker-compose up -d
go run cmd/api/main.go
```

API будет доступен на `http://localhost:8080`

### 2. Запустите админ панель

```bash
# Перейдите в директорию ui-admin
cd ui-admin

# Установите зависимости (если еще не установлены)
pnpm install

# Создайте .env.local (уже создан)
# NEXT_PUBLIC_API_URL=http://localhost:8080

# Запустите dev сервер
pnpm dev
```

Админ панель будет доступна на `http://localhost:3000`

### 3. Создайте тестового пользователя

Зайдите на `http://localhost:3000/register` и создайте пользователя:

```
Username: admin
Email: admin@example.com
Full Name: Администратор
Phone: +996500123456
Password: admin123
```

### 4. Войдите в систему

После регистрации вы автоматически будете перенаправлены на дашборд.
Или зайдите на `http://localhost:3000/login`:

```
Username: admin
Password: admin123
```

## Структура API

Админ панель интегрирована со следующими endpoints:

### Auth

- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход
- `POST /api/auth/logout` - Выход
- `GET /api/auth/me` - Текущий пользователь

### Users

- `GET /api/users?page=1&limit=10` - Список пользователей
- `GET /api/users/:id` - Пользователь по ID

### Organizations

- `GET /api/esf-organizations` - Все организации
- `GET /api/esf-organizations/paginated` - С пагинацией
- `GET /api/esf-organizations/:id` - По ID
- `POST /api/esf-organizations` - Создать (требует JWT)
- `PUT /api/esf-organizations/:id` - Обновить (требует JWT)
- `DELETE /api/esf-organizations/:id` - Удалить (требует JWT)

### Documents

- `GET /api/esf-documents` - Все документы
- `GET /api/esf-documents/paginated` - С пагинацией
- `GET /api/esf-documents/:id` - По ID
- `POST /api/esf-documents` - Создать (требует JWT)
- `PUT /api/esf-documents/:id` - Обновить (требует JWT)
- `DELETE /api/esf-documents/:id` - Удалить (требует JWT)

## Следующие шаги для разработки

### В разработке (TODO)

1. **Формы создания/редактирования**

   - Модальные окна для организаций
   - Формы документов с валидацией всех полей
   - Управление пользователями (изменение ролей)

2. **Детальные страницы**

   - Детальный просмотр документа с таблицей записей
   - История изменений
   - Связанные данные

3. **Дополнительный функционал**

   - Фильтры по датам
   - Экспорт в Excel/CSV
   - Графики и аналитика
   - Уведомления
   - Темная тема
   - Мультиязычность

4. **Оптимизация**
   - Server-side rendering для SEO
   - Виртуализация длинных списков
   - Оптимизация изображений
   - Code splitting

## Технологический стек

- **Frontend**: Next.js 15, React 19, TypeScript
- **State**: Zustand (auth), React Query (data)
- **Styling**: Tailwind CSS
- **Forms**: React Hook Form + Zod
- **HTTP**: Axios
- **Icons**: Lucide React
- **Dates**: date-fns

## Полезные команды

```bash
# Разработка
pnpm dev          # Запуск dev сервера
pnpm build        # Production сборка
pnpm start        # Запуск production
pnpm lint         # Линтинг кода

# Тестирование API
curl http://localhost:8080/health
curl http://localhost:8080/api/users
```

## Примечания

- Все страницы адаптивные и работают на мобильных устройствах
- JWT токены хранятся в localStorage
- При 401 ошибке происходит автоматический logout
- Поиск работает в реальном времени на клиенте
- Пагинация серверная для оптимизации

## Проблемы и решения

**Проблема**: CORS ошибки при обращении к API  
**Решение**: Убедитесь, что в Go API настроены CORS headers

**Проблема**: Токен не сохраняется  
**Решение**: Проверьте localStorage в DevTools → Application

**Проблема**: Данные не загружаются  
**Решение**: Проверьте, что Backend API запущен на порту 8080

## Контакты

При возникновении вопросов или проблем создайте issue в репозитории.
