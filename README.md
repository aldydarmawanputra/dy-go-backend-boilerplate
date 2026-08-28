# dy-go-backend-boilerplate

Boilerplate REST API dengan **Go + Fiber (fasthttp)**, **PostgreSQL** (GORM + raw query), **Redis**, migrasi **dbmate**, **JWT auth**, CORS configurable, konfigurasi lewat **.env**, dan **Docker Compose**.

Go module: `go-backend-boilerplate`.

## Struktur

Layout **berbasis modul fitur** (`internal/modules/<fitur>`) — enak buat proyek kecil maupun besar. Nambah fitur = nambah folder modul, tanpa bongkar arsitektur.

```
cmd/api/                 # entrypoint (main.go)
internal/
  config/                # load .env -> struct Config (build DSN & Redis addr)
  database/              # koneksi GORM + pool, AutoMigrate opsional
  cache/                 # koneksi Redis (go-redis)
  server/                # bikin *fiber.App + wiring route (routes.go) + health
  middleware/            # JWT auth, recover, logger, cors, requestid
  shared/                # util lintas-fitur
    apperror/            #   error domain -> HTTP status + code
    response/            #   envelope JSON standar {success,data,error,meta}
    validator/           #   wrapper go-playground/validator
    hash/                #   bcrypt hash & compare
    jwtutil/             #   generate & parse JWT
  modules/
    user/                # model (users + user_details), dto, repository, service, handler
    auth/                # dto, service, handler (register/login)
    docs/                # dokumentasi API markdown (en/id), di-embed via go:embed
db/migrations/           # file migrasi dbmate (.sql)
```

Alur request: **handler** (parse + validate) -> **service** (business logic) -> **repository** (akses DB).

## Konfigurasi (.env)

Semua config lewat env (komponen terpisah, bukan satu URL):

| Group | Variable |
|-------|----------|
| App   | `APP_HOST`, `APP_PORT`, `APP_ENV` |
| DB    | `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE` |
| Redis | `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB` |
| JWT   | `JWT_SECRET`, `JWT_EXPIRE_HOURS` |
| CORS  | `CORS_ALLOW_ORIGINS`, `CORS_ALLOW_METHODS`, `CORS_ALLOW_HEADERS`, `CORS_ALLOW_CREDENTIALS` |
| Misc  | `AUTO_MIGRATE` |

App membangun DSN Postgres dari `DB_*`. dbmate butuh `DATABASE_URL`, yang otomatis dibangun oleh Makefile dari `DB_*` tersebut — jadi cukup atur komponennya sekali di `.env`.

## Prasyarat

- Go 1.25+
- Docker (buat Postgres + Redis) — atau install lokal
- [dbmate](https://github.com/amacneil/dbmate) untuk migrasi dari host:
  ```bash
  go install github.com/amacneil/dbmate/v2@latest
  ```
  (Alternatif: pakai service `migrate` di docker-compose, ga perlu install.)

## Menjalankan (lokal)

```bash
cp .env.example .env
go mod tidy
docker compose up -d postgres redis
make migrate-up
make run
```

Server jalan di `http://localhost:8080`.

## Menjalankan (full Docker)

```bash
docker compose up --build
```

Urutan: Postgres + Redis up -> service `migrate` jalan -> app start.

## Endpoint

| Method | Path                    | Auth | Keterangan                          |
|--------|-------------------------|------|-------------------------------------|
| GET    | `/health`               | -    | Health check (status db & redis)    |
| GET    | `/documentation`        | -    | Redirect ke `/documentation/en`     |
| GET    | `/documentation/en`     | -    | Dokumentasi API (English, markdown) |
| GET    | `/documentation/id`     | -    | Dokumentasi API (Indonesia, markdown)|
| POST   | `/api/v1/auth/register` | -    | Daftar user baru                    |
| POST   | `/api/v1/auth/login`    | -    | Login -> dapat JWT                  |
| GET    | `/api/v1/users`         | ✅   | List/search user (raw query)        |
| POST   | `/api/v1/users`         | ✅   | Buat user (+ detail opsional)       |
| GET    | `/api/v1/users/me`      | ✅   | Profil user yang login              |
| GET    | `/api/v1/users/:id`     | ✅   | Ambil satu user + detail            |
| PUT    | `/api/v1/users/:id`     | ✅   | Replace (name + detail)             |
| PATCH  | `/api/v1/users/:id`     | ✅   | Partial update                      |
| DELETE | `/api/v1/users/:id`     | ✅   | Hapus user (cascade ke detail)      |

Envelope response konsisten:

```json
{ "success": true, "data": { }, "meta": { } }
{ "success": false, "error": { "code": "UNAUTHORIZED", "message": "..." } }
```

## Contoh pakai (curl)

```bash
curl -s -X POST localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"aldy@example.com","name":"Aldy","password":"secret123"}'

curl -s -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"aldy@example.com","password":"secret123"}'

TOKEN=<paste access_token>

curl -s localhost:8080/api/v1/users/me -H "Authorization: Bearer $TOKEN"

curl -s -X PATCH localhost:8080/api/v1/users/<id> \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Aldy D","detail":{"city":"Jakarta","phone":"08123"}}'

curl -s -X DELETE localhost:8080/api/v1/users/<id> -H "Authorization: Bearer $TOKEN"
```

## Data model

- `users` — id, email, password_hash, name, timestamps.
- `user_details` — one-to-one ke `users` (FK `ON DELETE CASCADE`): phone, address, city, country, bio, avatar_url.

## Migrasi (dbmate)

```bash
make migrate-new name=create_posts
make migrate-up
make migrate-down
```

## Nambah modul baru (mis. `post`)

1. Buat `internal/modules/post/` (contek pola `user`: model, dto, repository, service, handler).
2. Wiring di `internal/server/routes.go`: repo -> service -> handler, lalu `v1.Group("/posts", ...)`.
3. Migrasi: `make migrate-new name=create_posts`.
4. (Opsional dev) daftarkan model di `internal/database/migrate.go` kalau pakai `AUTO_MIGRATE`.

## GORM vs raw query

- CRUD standar -> GORM query builder (`user/repository.go`: `Create`, `First`, `Save`).
- Query berat -> raw SQL via `db.Raw(...).Scan(...)` (`Search` di `user/repository.go`). Selalu pakai placeholder (`$1`, `$2`) — jangan string concat.

## Testing

```bash
make test
```

## Catatan

- Skema DB dikelola **dbmate** (source of truth). `AUTO_MIGRATE=true` hanya kenyamanan dev; default `false`.
- Redis bersifat opsional saat boot: kalau ga reachable, app tetap jalan (health check nandain `redis: down`).
- Ganti `JWT_SECRET` di production. Kalau `CORS_ALLOW_CREDENTIALS=true`, `CORS_ALLOW_ORIGINS` tidak boleh `*`.
