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

## gRPC тесты (без reflection, с указанием proto)
Учти include для googleapis:
```
export GAPI="$(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis"
```

- Register (auth):
```
grpcurl -plaintext \
  -import-path proto \
  -import-path "$GAPI" \
  -proto proto/auth/v1/auth.proto \
  -d '{"username":"alice","email":"a@ex.com","password":"pass12345"}' \
  localhost:50051 auth.v1.AuthService/Register
```
- Login:
```
grpcurl -plaintext -import-path proto -import-path "$GAPI" \
  -proto proto/auth/v1/auth.proto \
  -d '{"email":"a@ex.com","password":"pass12345"}' \
  localhost:50051 auth.v1.AuthService/Login
```
- Refresh/Logout аналогично (подставь refresh_token).

- User-service CreateProfile:
```
grpcurl -plaintext -import-path proto -import-path "$GAPI" \
  -proto proto/user/v1/user.proto \
  -d '{"user_id":"<uuid>","username":"alice"}' \
  localhost:50052 user.v1.UserService/CreateProfile
```

## HTTP через gateway
(если JWT-миддлварь включена, публичные пути `/v1/auth/register`, `/v1/auth/login` не требуют токена):
```
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"a@ex.com","password":"pass12345"}'
```
Остальные эндпоинты — с `Authorization: Bearer <access>`.

## Примечания
- В миграциях уникальные префиксы обязательны (Goose: версии должны отличаться).
- В DELETE для refresh-токенов убран LIMIT (Postgres).
- Для refresh-токена регэксп/генератор должны совпадать (если используешь base64 с паддингом — разреши `=` в валидации или убери паддинг через `base64.RawURLEncoding`).
- Для тестов grpcurl без proto можно включить gRPC reflection в dev, но по умолчанию выключено.
