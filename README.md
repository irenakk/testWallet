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

- Использование `SELECT FOR UPDATE` для блокировки строки на уровне БД
- Connection pool с настраиваемым количеством соединений (по умолчанию 100)
- Атомарные операции через транзакции PostgreSQL
- `ON CONFLICT` для безопасного создания кошелька при первом депозите