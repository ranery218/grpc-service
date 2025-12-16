# Friend Service (auth + user + gateway)

## Требования
- Go ≥ 1.24 (grpc 1.77+)
- Docker + docker-compose (для Postgres)
- Goose для миграций: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- grpcurl (для gRPC тестов): `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

## Сервисы и директории
- Auth-service: `cmd/auth-service`, конфиг `config/auth.yaml`
- User-service: `cmd/user-service`, конфиг `config/user.yaml`
- Gateway: `cmd/gateway`, конфиг `config/gateway.yaml`
- Миграции: `migrations/auth`, `migrations/user`
- Proto: `proto/auth/v1/auth.proto`, `proto/user/v1/user.proto` (сгенерировано в `proto/gen/...`)

## Запуск Postgres
```
docker compose up -d postgres
```

## Прогон миграций
```
goose -dir migrations/auth postgres "postgres://app:app@localhost:5432/auth_db?sslmode=disable" up
goose -dir migrations/user  postgres "postgres://app:app@localhost:5432/user_db?sslmode=disable" up
```

## Конфиги (пример)
- `config/auth.yaml`
  - `server.grpc_addr: ":50051"`
  - `db.dsn: "postgres://app:app@localhost:5432/auth_db?sslmode=disable"`
  - `jwt.secret`, `jwt.access_ttl`, `jwt.refresh_ttl`, `jwt.iss`, `jwt.aud`
  - `clients.user_service_addr: "localhost:50052"`
- `config/user.yaml`
  - `server.grpc_addr: ":50052"`
  - `db.dsn: "postgres://app:app@localhost:5432/user_db?sslmode=disable"`
- `config/gateway.yaml`
  - `server.http_addr: ":8080"`
  - `backends.auth_addr: "localhost:50051"`
  - `backends.user_addr: "localhost:50052"`
  - `jwt.secret/iss/aud` (для проверки access токенов, если включено)

## Запуск сервисов (локально)
В разных терминалах:
```
go run ./cmd/user-service
go run ./cmd/auth-service
go run ./cmd/gateway
```

## HTTP через gateway (основное тестирование)
Шлюз на `:8080`, REST → gRPC. Публичные маршруты не требуют токен, остальные с `Authorization: Bearer <access_token>`.

- `POST /v1/auth/register`
  Body: `{"username":"alice","email":"a@ex.com","password":"pass12345"}`
  Ответ: `{"accessToken","accessExpiresAt","refreshToken","refreshExpiresAt"}`

- `POST /v1/auth/login`
  Body: `{"email":"a@ex.com","password":"pass12345"}`
  Ответ: `{"accessToken","accessExpiresAt","refreshToken","refreshExpiresAt"}`

- `POST /v1/auth/refresh`
  Body: `{"refresh_token":"..."}` → новый `access/refresh`.

- `POST /v1/auth/logout`
  Body: `{"refresh_token":"..."}` → 200 OK.

- `POST /v1/user/profile`
  Body: `{"user_id":"<uuid>","username":"alice"}` → создаёт профиль.

- `GET /v1/user/profile` (с user_id в query) → профиль.
- `PUT /v1/user/profile` → обновление профиля.
- `DELETE /v1/user/profile` → удаление профиля.
- Друзья: `POST/GET/DELETE /v1/user/friends` (согласно proto/user/v1/user.proto).

## gRPC тесты (опционально, grpcurl)
При необходимости можно дергать gRPC напрямую, указав `-proto proto/...` и include для googleapis. См. примеры в proto или генерируйте через grpcurl с путями к .proto.
