# Dokumentasi API

`go-backend-boilerplate` · REST · JSON · [English](/documentation/en)

- **Base URL:** `http://localhost:8080`
- **Prefix API:** `/api/v1`
- **Content-Type:** `application/json`

REST API awalan yang dibangun dengan Go + Fiber, PostgreSQL (GORM + raw query), Redis, dan autentikasi JWT.

---

## Autentikasi

Endpoint terproteksi butuh **JWT** yang didapat dari endpoint login. Kirim sebagai Bearer token:

```
Authorization: Bearer <access_token>
```

Token kedaluwarsa setelah `JWT_EXPIRE_HOURS` (default 24 jam). Token yang tidak ada atau tidak valid mengembalikan `401 Unauthorized`.

---

## Format Response

Setiap response memakai envelope yang konsisten.

**Sukses**

```json
{
  "success": true,
  "data": { },
  "meta": { }
}
```

`meta` bersifat opsional (mis. paginasi).

**Error**

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "validation failed",
    "fields": { "email": "must be a valid email address" }
  }
}
```

---

## Endpoint

### Sistem

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| GET | `/health` | — | Status layanan, termasuk status DB & Redis. |

### Auth

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| POST | `/api/v1/auth/register` | — | Membuat akun baru. |
| POST | `/api/v1/auth/login` | — | Autentikasi dan menerima JWT. |

### Users

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| GET | `/api/v1/users` | 🔒 | List / cari user. Query: `q`, `limit`, `offset`. |
| POST | `/api/v1/users` | 🔒 | Membuat user dengan detail opsional. |
| GET | `/api/v1/users/me` | 🔒 | Profil user yang sedang login. |
| GET | `/api/v1/users/:id` | 🔒 | Mengambil satu user beserta detail. |
| PUT | `/api/v1/users/:id` | 🔒 | Mengganti name dan detail (update penuh). |
| PATCH | `/api/v1/users/:id` | 🔒 | Update sebagian (hanya field yang dikirim). |
| DELETE | `/api/v1/users/:id` | 🔒 | Menghapus user (cascade ke detail). |

---

## Contoh

**Register**

```json
POST /api/v1/auth/register
{
  "email": "aldy@example.com",
  "name": "Aldy",
  "password": "secret123"
}
```

**Login → token**

```json
POST /api/v1/auth/login
{ "email": "aldy@example.com", "password": "secret123" }

// 200
{ "success": true, "data": { "access_token": "eyJ...", "token_type": "Bearer" } }
```

**Update sebagian**

```json
PATCH /api/v1/users/<id>
Authorization: Bearer <token>
{ "name": "Aldy D", "detail": { "city": "Jakarta", "phone": "08123" } }
```

---

## Model Data

### User

| Field | Tipe | Catatan |
|-------|------|---------|
| `id` | uuid | Primary key. |
| `email` | string | Unik. |
| `name` | string | |
| `detail` | UserDetail | One-to-one, opsional. |
| `created_at` / `updated_at` | timestamp | |

### UserDetail

| Field | Tipe | Catatan |
|-------|------|---------|
| `phone` | string | |
| `address` | string | |
| `city` | string | |
| `country` | string | |
| `bio` | string | |
| `avatar_url` | string | Harus URL valid. |

---

## Kode Error

| Code | HTTP | Arti |
|------|------|------|
| `BAD_REQUEST` | 400 | Body tidak valid. |
| `UNAUTHORIZED` | 401 | Token/kredensial tidak ada atau salah. |
| `NOT_FOUND` | 404 | Resource tidak ditemukan. |
| `CONFLICT` | 409 | Duplikat (mis. email sudah terdaftar). |
| `VALIDATION_ERROR` | 422 | Validasi field gagal. |
| `INTERNAL` | 500 | Error server tak terduga. |
