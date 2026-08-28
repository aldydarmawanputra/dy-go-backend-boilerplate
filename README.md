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
    model/               #   Base model (UUID google + timestamps) buat di-embed model
    apperror/            #   error domain -> HTTP status + code
    response/            #   envelope JSON standar {success,data,error,meta}
    validator/           #   wrapper go-playground/validator
    redact/              #   tipe Secret (auto [REDACTED] di JSON/log) + mask helper
    sanitize/            #   trim/normalisasi input (zero-trust)
    pagination/          #   params + meta (limit/offset/total)
    logging/             #   setup slog (JSON di prod, text di dev)
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

## Setelah clone (init sekali)

Pastikan ada server PostgreSQL jalan (lokal atau `docker compose up -d postgres redis`), lalu:

```bash
make setup
```

`make setup` otomatis: bikin `.env` dari `.env.example` (kalau belum ada) → `go mod download` → pasang `dbmate` (kalau belum ada) → **bikin database sesuai `.env` + jalanin semua migrasi** (`dbmate up`). Sesuaikan kredensial DB di `.env` dulu sebelum jalanin.

## Menjalankan (lokal)

```bash
make run
```

Server jalan di `http://localhost:8080`. Saat start, app juga otomatis **membuat database** (sesuai `DB_*` di `.env`) kalau belum ada — lihat `internal/database/ensure.go`. Skema tabel tetap dikelola dbmate (`make setup` / `make migrate-up`).

## Menjalankan (full Docker)

```bash
docker compose up --build
```

Urutan: Postgres + Redis up -> service `migrate` jalan -> app start.

Config service `app` & `migrate` dibaca dari **`.env.docker`** (via `env_file`), bukan `environment` inline. Hostname di situ pakai nama service compose (`DB_HOST=postgres`, `REDIS_HOST=redis`) — beda dari `.env` lokal yang pakai `localhost`. `.env` asli & semua `.env*` tidak ikut ke image (lihat `.dockerignore`).

## Endpoint

| Method | Path                    | Auth | Keterangan                          |
|--------|-------------------------|------|-------------------------------------|
| GET    | `/health`               | -    | Health check (status db & redis)    |
| GET    | `/documentation`        | -    | Redirect ke `/documentation/en`     |
| GET    | `/documentation/en`     | -    | Dokumentasi API (English, markdown) |
| GET    | `/documentation/id`     | -    | Dokumentasi API (Indonesia, markdown)|
| POST   | `/api/v1/auth/register` | -    | Daftar user baru                    |
| POST   | `/api/v1/auth/login`    | -    | Login -> access + refresh token     |
| POST   | `/api/v1/auth/refresh`  | -    | Tukar refresh token -> pasangan baru (rotation) |
| POST   | `/api/v1/auth/logout`   | -    | Revoke refresh token                |
| GET    | `/api/v1/users`         | ✅   | List/search user (raw query)        |
| POST   | `/api/v1/users`         | admin | Buat user (+ detail opsional)     |
| GET    | `/api/v1/users/search`  | ✅   | Full-text search (Postgres tsvector) `?q=` |
| GET    | `/api/v1/users/me`      | ✅   | Profil user yang login              |
| GET    | `/api/v1/users/:id`     | ✅   | Ambil satu user + detail            |
| PUT    | `/api/v1/users/:id`     | ✅   | Replace (name + detail)             |
| PATCH  | `/api/v1/users/:id`     | ✅   | Partial update                      |
| DELETE | `/api/v1/users/:id`     | admin | Hapus user (cascade ke detail)    |
| POST   | `/api/v1/files`         | ✅   | Upload file (multipart `file`) → `{key,url}` |

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

Semua ID tabel entity pakai **UUID (google/uuid)**, di-generate otomatis lewat GORM hook `BeforeCreate` di `internal/shared/model`. Bikin model baru cukup embed `model.Base` — ID + timestamps langsung ke-handle.

### Roles & authorization

- `roles` — id integer + name, di-seed `1=admin`, `2=user` (migrasi).
- `user_roles` — join many-to-many `users`↔`roles` (FK `ON DELETE CASCADE`). Satu user bisa punya >1 role.
- User baru (register) otomatis dapat role `user`.
- JWT membawa claim `roles` (array). Login memuat role user ke token.
- Middleware `middleware.RequireRole("admin")` menjaga endpoint (fail-closed: tanpa role yang cocok → `403`). Contoh: `POST` & `DELETE /users` khusus admin.
- Jadikan seorang user admin (sementara, via SQL):
  ```sql
  INSERT INTO user_roles (user_id, role_id)
  SELECT '<user-uuid>', id FROM roles WHERE name = 'admin'
  ON CONFLICT DO NOTHING;
  ```

## i18n (multi-bahasa)

Pesan error diterjemah otomatis berdasarkan **locale request**: dari `?lang=id` atau header `Accept-Language`. Didukung `en` (default) & `id`. Katalog di `internal/shared/i18n`. Client sebaiknya tetap pakai `error.code` yang stabil; `error.message` yang human-readable ngikut bahasa.

```bash
curl -s localhost:8080/api/v1/users/xxx -H "Accept-Language: id" ...
# {"success":false,"error":{"code":"UNAUTHORIZED","message":"tidak terautentikasi"}}
```

Nambah bahasa = tambah entri di `catalog`. Nambah pesan = tambah key baru.

## Payment (kerangka)

Abstraksi `payment.Gateway` (interface `Charge` + `ParseWebhook`) dengan driver `PAYMENT_PROVIDER` — sekarang cuma **`stub`** (bikin charge palsu, tanpa network). Integrasi Midtrans/Xendit/Stripe tinggal implement `Gateway` dan daftarin di `payment.New()`.

- `POST /api/v1/payments` (login) → bikin charge, simpan ke tabel `payments`, balikin `payment_url`.
- `POST /api/v1/payments/webhook` (public) → update status by `order_id`. **Di provider beneran, verifikasi signature dulu di sini.**

## Realtime (WebSocket)

Kerangka pub/sub general di `internal/realtime`: satu `Hub` (broadcast ke semua client, register/unregister, drop client lambat). Endpoint:

- `GET /ws` — koneksi WebSocket (client subscribe).
- `POST /api/v1/broadcast` (admin) — push pesan ke semua client `{ "message": "..." }` → server→clients.

Ini sengaja generik — tinggal extend jadi rooms/topic/auth-per-koneksi sesuai fitur (chat, notif live, dsb).

## Email (SMTP)

`mailer.Mailer` pakai SMTP standar (net/smtp, zero-dep). Default nunjuk ke **Mailpit** (`docker compose up mailpit`) — email nyangkut di UI `http://localhost:8025`, ga keluar ke inbox beneran, **gratis** buat dev. Register ngirim welcome email (best-effort, ga nge-block).

Ganti ke email beneran tinggal edit `.env` (contoh Gmail):

```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=you@gmail.com
SMTP_PASSWORD=<app-password>
SMTP_FROM=you@gmail.com
```

net/smtp otomatis STARTTLS di `:587`. Kalau `SMTP_USER` kosong → tanpa auth (Mailpit).

## File storage

Abstraksi `storage.Storage` dengan driver dipilih lewat `STORAGE_DRIVER`:

- `local` (default) — simpan ke `STORAGE_LOCAL_PATH` (`./storage`), disajikan di `GET /storage/*`. URL publik dari `STORAGE_PUBLIC_BASE_URL`.
- `r2` — Cloudflare R2 (S3-compatible via aws-sdk-go-v2). Isi `R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `R2_BUCKET`, `R2_PUBLIC_BASE_URL`.
- `supabase` — Supabase Storage (REST). Isi `SUPABASE_URL`, `SUPABASE_SERVICE_KEY`, `SUPABASE_BUCKET`, `SUPABASE_PUBLIC_BASE_URL`.

Upload: `POST /api/v1/files` (multipart, field `file`) → `{ "key": "...", "url": "..." }`. Key di-generate UUID + ekstensi; driver local aman dari path-traversal (`filepath.Clean`).

Kalau file-nya **gambar** (jpeg/png), otomatis dibuatin **thumbnail** (maks 300×300, jaga rasio, high-quality scaling via `x/image`) dan disimpan di `uploads/thumbs/` — response nambah `thumbnail_url`. Lihat `internal/imageproc`.

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
- **OpenTelemetry** tracing opsional: set `OTEL_ENABLED=true` + `OTEL_EXPORTER_OTLP_ENDPOINT` (mis. `localhost:4318`). Kalau `false` (default), no-op — zero overhead. Otomatis instrument request Fiber via `otelfiber`.
- Ganti `JWT_SECRET` di production. Kalau `CORS_ALLOW_CREDENTIALS=true`, `CORS_ALLOW_ORIGINS` tidak boleh `*`.

### Keamanan field sensitif

- Password disimpan sebagai **bcrypt hash** (`internal/shared/hash`), tidak pernah plaintext.
- `PasswordHash` di-tag `json:"-"` → tidak pernah muncul di response.
- Field password di DTO pakai tipe `redact.Secret` → otomatis jadi `[REDACTED]` kalau ke-marshal JSON atau ke-log; nilai asli cuma dibuka via `.Reveal()` saat hashing.
- GORM SQL log pakai `ParameterizedQueries: true` → nilai argumen (mis. `password_hash`, email) tidak di-expand di log, cuma muncul `$1, $2`.
- Request logger (Fiber) tidak mencatat body, jadi password di body tidak masuk access log.
- Buat mask nilai saat logging manual: `redact.String(x)`, `redact.Email(x)`, `redact.Tail(x, n)`.

### Proteksi lain

- **Rate limit**: global (`RATE_LIMIT_MAX`/`RATE_LIMIT_WINDOW_SEC`) + limiter ketat khusus `/auth` (`AUTH_RATE_LIMIT_MAX`) → tahan brute-force. Balikan `429 TOO_MANY_REQUESTS`.
- **Timeout & body limit**: `READ/WRITE/IDLE_TIMEOUT_SEC` + `BODY_LIMIT_BYTES` di `fiber.Config` → tahan slow-loris & payload gede.
- **Security headers**: middleware `helmet`.
- **IDOR**: `RequireSelfOrAdmin("id")` di `GET/PUT/PATCH /users/:id` → non-admin cuma bisa akses dirinya sendiri.
- **Timing/user-enumeration**: login melakukan `bcrypt` dummy-compare saat email tak ditemukan → waktu respons seragam.
- **SQL injection**: semua query pakai placeholder parameter (GORM builder & `Raw($1,$2)`), tak ada string-concat input.
- **Sanitasi input (zero-trust)**: `sanitize.String` (trim + buang control char) & `sanitize.Email` (lowercase+trim) di service sebelum simpan.
