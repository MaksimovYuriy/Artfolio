# Artfolio

Artfolio — публичное портфолио художницы с закрытой административной панелью. Владелец может редактировать профиль и социальные ссылки, загружать работы, менять порядок и управлять публикацией.

Проект состоит из:

- backend на Go с REST API `/api/v1`;
- frontend на React, TypeScript, Vite и Material UI;
- PostgreSQL;
- Caddy для HTTPS, маршрутизации API и раздачи изображений;
- Docker Compose для production-запуска.

## Первый запуск через Docker Compose

Требуются Docker и Docker Compose.

1. Создайте файл окружения:

   ```bash
   cp .env.example .env
   ```

2. Задайте в `.env` надёжный `DB_PASSWORD`. Для production также проверьте остальные параметры и домен в `Caddyfile`.

3. Запустите приложение:

   ```bash
   docker compose up -d
   ```

   Сервис `migrate` автоматически применит миграции перед запуском backend.

4. Создайте административный ключ:

   ```bash
   docker compose run --rm key
   ```

   Ключ выводится один раз. Сохраните его в менеджере паролей.

5. Откройте `/admin`, войдите с созданным ключом и заполните профиль, социальные ссылки и работы.

Файлы `.env` и административные ключи нельзя добавлять в Git или публиковать в логах.

## Переменные окружения

Полный шаблон находится в `.env.example`.

| Переменная | Назначение |
| --- | --- |
| `APP_ENV` | Режим логирования: `dev`, `test` или `prod` |
| `HTTP_*` | Адрес, порт и тайм-ауты backend |
| `DB_*` | Подключение к PostgreSQL |
| `STORAGE_PATH` | Каталог загруженных изображений |
| `STORAGE_PUBLIC_URL` | Публичный префикс URL изображений |
| `STORAGE_MAX_FILE_SIZE` | Максимальный размер файла в байтах |
| `STORAGE_MAX_PIXELS` | Максимальное количество пикселей изображения |
| `GOOSE_*` | Параметры CLI Goose для ручной работы с миграциями |

## Логирование

Backend пишет структурированные логи в стандартный вывод контейнера. В `prod` используется JSON,
в `dev` — удобный для чтения текстовый формат. Каждая HTTP-запись содержит request ID, метод,
маршрут, статус, размер ответа и длительность. Для авторизованных запросов также записываются
внутренние `actor_id` и `session_id`, позволяющие связать действия пользователя и отдельного входа.
Заголовок `X-Request-ID` возвращается клиенту и может использоваться для поиска конкретного запроса
в логах.

Внутренние причины серверных ошибок записываются в лог, но не возвращаются клиенту. Тела запросов,
административные ключи и session cookie не логируются. Docker Compose ограничивает локальные файлы
логов тремя файлами по 10 МБ для каждого сервиса.

Просмотр логов backend:

```bash
docker compose logs -f backend
```

В Docker Compose backend хранит изображения в volume `artwork_media`, а PostgreSQL — в `postgres_data`.

## Локальная разработка

Backend требует Go 1.26 и доступный PostgreSQL. Переменные `DB_*` и `STORAGE_*` можно взять из `.env.example`, изменив адреса и пути под локальное окружение.

```bash
cd backend
go run ./cmd/migrate
go run ./cmd/api
```

Frontend требует Node.js 22:

```bash
cd frontend
npm ci
npm run dev
```

Для локальной разработки запросы `/api` и `/media` должны проксироваться на backend либо обслуживаться через локально настроенный reverse proxy.

## Проверки

Unit-тесты backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm run lint
npm run build
```

### Интеграционные тесты backend

Интеграционные тесты из `backend/integration-tests` используют отдельный PostgreSQL на порту `55432`:

```bash
docker compose -f compose.integration.yml up -d --wait
cd backend
go test -tags=integration ./integration-tests
cd ..
docker compose -f compose.integration.yml down
```

Для внешней тестовой БД передайте `ARTFOLIO_TEST_DATABASE_URL`.

## Основные маршруты

- `/` — публичное портфолио;
- `/admin` — административная панель;
- `/api/v1` — REST API;
- `/media` — загруженные изображения.
