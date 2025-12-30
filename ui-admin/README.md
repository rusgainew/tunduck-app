# Tunduck Admin Panel

Админ панель для управления системой Tunduck ESF.

## Технологии

- **Next.js 15** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **React Query** - Data fetching и кеширование
- **Zustand** - State management
- **React Hook Form** - Формы с валидацией
- **Zod** - Schema validation
- **Axios** - HTTP client
- **Lucide React** - Icons
- **date-fns** - Работа с датами

## Возможности

### Аутентификация

- ✅ Вход в систему
- ✅ Регистрация
- ✅ Выход из системы
- ✅ JWT токены
- ✅ Автоматический редирект

### Управление пользователями

- ✅ Просмотр списка пользователей
- ✅ Пагинация
- ✅ Поиск по имени, email, username
- ✅ Отображение ролей и статусов
- ⏳ Создание/редактирование (в разработке)

### Управление организациями

- ✅ Просмотр списка организаций ЭСФ
- ✅ Поиск по названию, описанию, БД
- ✅ Удаление организаций
- ⏳ Создание/редактирование (в разработке)

### Управление документами

- ✅ Просмотр списка документов ЭСФ
- ✅ Поиск по ИНН, названию, комментарию
- ✅ Отображение основной информации
- ⏳ Детальный просмотр (в разработке)
- ⏳ Создание/редактирование (в разработке)

### Дашборд

- ✅ Статистика по пользователям, организациям, документам
- ✅ Последние пользователи
- ✅ Список организаций
- ✅ Быстрый доступ к разделам

## Установка

```bash
# Установите зависимости
pnpm install

# Создайте .env.local файл
cp .env.example .env.local

# Отредактируйте .env.local и укажите URL API
# NEXT_PUBLIC_API_URL=http://localhost:8080

# Запустите dev сервер
pnpm dev
```

Откройте [http://localhost:3000](http://localhost:3000) в браузере.

## Структура проекта

```
ui-admin/
├── app/                      # Next.js App Router
│   ├── dashboard/           # Защищенные страницы дашборда
│   │   ├── layout.tsx      # Layout с навигацией
│   │   ├── page.tsx        # Главная страница дашборда
│   │   ├── users/          # Управление пользователями
│   │   ├── organizations/  # Управление организациями
│   │   └── documents/      # Управление документами
│   ├── login/              # Страница входа
│   ├── register/           # Страница регистрации
│   ├── layout.tsx          # Root layout
│   ├── page.tsx            # Home page (редирект)
│   ├── providers.tsx       # React Query provider
│   └── globals.css         # Global styles
├── lib/                     # Утилиты и хелперы
│   ├── api.ts              # API функции
│   ├── api-client.ts       # Axios конфигурация
│   ├── store.ts            # Zustand store (auth)
│   └── types.ts            # TypeScript типы
└── public/                  # Статичные файлы
```

## API Endpoints

Админ панель использует следующие API endpoints:

### Auth

- `POST /api/auth/login` - Вход
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/logout` - Выход
- `GET /api/auth/me` - Текущий пользователь

### Users

- `GET /api/users` - Список пользователей (с пагинацией)
- `GET /api/users/:id` - Пользователь по ID

### Organizations

- `GET /api/esf-organizations` - Список организаций
- `GET /api/esf-organizations/paginated` - С пагинацией
- `GET /api/esf-organizations/:id` - По ID
- `POST /api/esf-organizations` - Создать
- `PUT /api/esf-organizations/:id` - Обновить
- `DELETE /api/esf-organizations/:id` - Удалить

### Documents

- `GET /api/esf-documents` - Список документов
- `GET /api/esf-documents/paginated` - С пагинацией
- `GET /api/esf-documents/:id` - По ID
- `POST /api/esf-documents` - Создать
- `PUT /api/esf-documents/:id` - Обновить
- `DELETE /api/esf-documents/:id` - Удалить

## Скрипты

```bash
# Запуск в режиме разработки
pnpm dev

# Сборка для production
pnpm build

# Запуск production сборки
pnpm start

# Линтинг
pnpm lint
```

## Переменные окружения

Создайте `.env.local` файл:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Будущие улучшения

- [ ] Формы создания/редактирования для всех сущностей
- [ ] Детальные страницы для документов
- [ ] Управление ролями и правами
- [ ] Фильтры и сортировка
- [ ] Экспорт данных
- [ ] Dark mode
- [ ] Уведомления
- [ ] Загрузка файлов
- [ ] Графики и аналитика
- [ ] Поддержка нескольких языков

## Лицензия

MIT
