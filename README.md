# Subscription Service

REST-сервис для управления и агрегации данных об онлайн-подписках пользователей.

## Описание

Сервис позволяет выполнять CRUDL-операции над подписками, а также вычислять суммарную стоимость подписок за заданный период времени с возможностью фильтрации по пользователю и названию сервиса.

---

## Запуск проекта

### Клонирование репозитория

Клонируйте репозиторий:
```bash
git clone https://github.com/Tommych123/auth-service.git
cd auth-service
```
---
### Переменные окружения

Необходимо создать файл в корне проекта `.env` по примеру из `deploy/.env-example`:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscription
DB_SSLMODE=disable
SERVER_PORT=8080
```

---

### Сборка с помощью Docker Compose

```bash
cd deploy
docker-compose up --build
```

## API

### CRUDL подписок

- `POST /subscriptions/` — создать подписку  
- `GET /subscriptions/` — получить все подписки  
- `GET /subscriptions/{id}` — получить подписку по ID  
- `PUT /subscriptions/{id}` — обновить подписку  
- `DELETE /subscriptions/{id}` — удалить подписку  

Пример запроса:

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

---

### Агрегация стоимости

- `GET /total` — получить сумму подписок за период

#### Параметры запроса:
| Параметр | Тип | Обязательный | Описание |
|----------|-----|---------------|----------|
| `from` | `string` | + | Начало периода в формате `MM-YYYY` |
| `to` | `string` | + | Конец периода в формате `MM-YYYY` |
| `user_id` | `string` (UUID) | - | ID пользователя |
| `service_name` | `string` | - | Название сервиса |

---

## Swagger-документация

После запуска доступна по адресу:  
**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

---

## База данных

Используется PostgreSQL.  
Миграции автоматически применяются при старте приложения (`golang-migrate`).

Структура таблицы `subscriptions`:

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор |
| `service_name` | VARCHAR | Название сервиса |
| `price` | INTEGER | Стоимость (в рублях, без копеек) |
| `user_id` | UUID | ID пользователя |
| `start_date` | DATE | Дата начала подписки |
| `end_date` | DATE (NULLABLE) | Дата окончания подписки |

---

## Логирование

Используется библиотека [`uber-go/zap`](https://github.com/uber-go/zap) для логирования:
- уровни: `info`, `warn`, `error`
- логируются события API, ошибок и миграций

---

## Структура проекта

```
.
├── api/               # HTTP handlers
├── cmd/               # main.go entrypoint
├── internal/          # документация и утилиты
├── models/            # структуры и типы
├── pkg/db/            # PostgreSQL и миграции
├── repository/        # Работа с БД
├── service/           # Бизнес-логика
├── deploy/            # Dockerfile, скрипты запуска
└── migrations/        # SQL-модули миграции
```

---