# Simple Wallet API

REST API для управления кошельками с поддержкой депозитов и снятий средств.

## Технологии

- Go 1.24
- PostgreSQL 16
- Docker & Docker Compose

## Запуск

### С Docker Compose (рекомендуется)

```bash
# Собрать и запустить все сервисы
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build
```

## API Endpoints

### POST /api/v1/wallet
Выполнение операции (депозит/снятие)
```json
{
    "walletId": "550e8400-e29b-41d4-a716-446655440000",
    "operationType": "DEPOSIT",
    "amount": 1000
}
```

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "balance": 1000
}
```

### GET /api/v1/wallets/{walletId}
Получение баланса кошелька

**Response (200 OK):**
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "balance": 1000
}
```

## Конкурентность

Приложение оптимизировано для обработки 1000 RPS на один кошелёк:

- `SELECT FOR UPDATE` для блокировки строки на уровне БД
- Connection pool с настраиваемым количеством соединений
- Транзакции PostgreSQL
- `ON CONFLICT`: безопасное создание кошелька при первом пополнении

## Конфигурация

Переменные окружения (config.env):

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| DATABASE_HOST | Хост БД | postgres |
| DATABASE_PORT | Порт БД | 5432 |
| DATABASE_USER | Пользователь БД | wallet_user |
| DATABASE_PASSWORD | Пароль БД | wallet_password |
| DATABASE_NAME | Имя БД | wallet_db |
| DATABASE_SSL_MODE | SSL режим | disable |
| DATABASE_MAX_CONNS | Макс. соединений | 100 |
| DATABASE_MIN_CONNS | Мин. соединений | 10 |
| SERVER_HOST | Хост сервера | 0.0.0.0 |
| SERVER_PORT | Порт сервера | 8080 |

## Тестирование

```bash
# Unit тесты
go test -v ./...

# Интеграционные тесты (требуется запущенная БД)
go test -v -tags=integration ./...

# С Makefile
make test
make test-integration
```

## Makefile команды

```bash
make build          # Собрать бинарник
make run            # Запустить локально
make test           # Запустить unit тесты
make test-integration # Запустить интеграционные тесты
make docker-up      # Запустить в Docker
make docker-down    # Остановить Docker
make docker-clean   # Удалить volumes
make logs           # Просмотр логов
```

## Структура проекта

```
.
├── cmd/server/          # Точка входа
├── internal/
│   ├── config/          # Конфигурация
│   ├── handler/         # HTTP обработчики
│   ├── model/           # Модели данных
│   ├── repository/      # Слой работы с БД
│   └── service/         # Бизнес-логика
├── migrations/          # SQL миграции
├── config.env           # Переменные окружения
├── Dockerfile
└── docker-compose.yml
```