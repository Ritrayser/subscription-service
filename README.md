# Subscription Service

REST-сервис для агрегации данных об онлайн-подписках пользователей.

## Стек

- **Язык:** Go 1.24
- **HTTP-роутер:** chi v5
- **СУБД:** PostgreSQL 15
- **Драйвер БД:** pgx v5 + sqlx
- **Логирование:** slog (JSON)
- **Документация:** swaggo / swagger-ui
- **Контейнеризация:** Docker Compose

## Структура

```
├── cmd/server/main.go          # точка входа, роутинг
├── internal/
│   ├── config/config.go        # конфигурация из .env
│   ├── models/subscription.go  # модели и валидация
│   ├── repository/postgres.go  # слой работы с БД
│   ├── service/subscription.go # бизнес-логика
│   └── handlers/subscription.go# HTTP-обработчики
├── migrations/
│   ├── 000001_init.up.sql      # создание таблицы
│   └── 000001_init.down.sql    # откат
├── docs/                       # swagger-спецификация
├── Dockerfile
├── docker-compose.yml
└── .env
```

## Запуск

```bash
docker compose up --build
```

Сервис будет доступен на `http://localhost:8080`.

## Конфигурация

Переменные окружения (`.env`):

| Переменная     | Значение по умолчанию |
|----------------|-----------------------|
| DB_HOST        | postgres              |
| DB_PORT        | 5432                  |
| DB_USER        | postgres              |
| DB_PASSWORD    | postgres              |
| DB_NAME        | subscriptions         |
| DB_SSLMODE     | disable               |
| SERVER_PORT    | 8080                  |

## API

### Подписки

| Метод   | Путь                   | Описание                        |
|---------|------------------------|---------------------------------|
| POST    | `/subscriptions`       | Создать подписку                |
| GET     | `/subscriptions`       | Список подписок (с фильтрацией) |
| GET     | `/subscriptions/{id}`  | Получить подписку               |
| PUT     | `/subscriptions/{id}`  | Обновить подписку               |
| DELETE  | `/subscriptions/{id}`  | Удалить подписку                |

**Параметры GET /subscriptions:**

- `limit` — лимит записей (default 20)
- `offset` — смещение (default 0)
- `user_id` — фильтр по UUID пользователя
- `service_name` — фильтр по названию сервиса

**Тело POST /subscriptions:**

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

Поле `end_date` опционально.

### Стоимость за период

| Метод | Путь                      | Описание                              |
|-------|---------------------------|---------------------------------------|
| GET   | `/subscriptions/total`    | Суммарная стоимость подписок за период|

**Параметры:**

- `start` (обязательный) — начало периода, `MM-YYYY`
- `end` (обязательный) — конец периода, `MM-YYYY`
- `user_id` — фильтр по UUID пользователя
- `service_name` — фильтр по названию сервиса

**Ответ:**

```json
{
  "total": 2400
}
```

Стоимость рассчитывается как сумма `price × количество активных месяцев` для каждой подписки, попадающей в период.

## Swagger

После запуска документация доступна по адресу:

```
http://localhost:8080/swagger/
```
