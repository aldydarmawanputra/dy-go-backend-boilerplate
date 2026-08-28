# API Documentation

`go-backend-boilerplate` · REST · JSON · [Bahasa Indonesia](/documentation/id)

- **Base URL:** `http://localhost:8080`
- **API prefix:** `/api/v1`
- **Content-Type:** `application/json`

A starter REST API built with Go + Fiber, PostgreSQL (GORM + raw query), Redis, and JWT authentication.

---

## Authentication

Protected endpoints require a **JWT** obtained from the login endpoint. Send it as a Bearer token:

```
Authorization: Bearer <access_token>
```

Tokens expire after `JWT_EXPIRE_HOURS` (default 24h). A missing or invalid token returns `401 Unauthorized`.

---

## Response Format

Every response uses a consistent envelope.

**Success**

```json
{
  "success": true,
  "data": { },
  "meta": { }
}
```

`meta` is optional (e.g. pagination).

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

## Endpoints

### System

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | — | Service health, including DB & Redis status. |

### Auth

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | — | Create a new account. |
| POST | `/api/v1/auth/login` | — | Authenticate and receive a JWT. |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/users` | 🔒 | List / search users. Query: `q`, `limit`, `offset`. |
| POST | `/api/v1/users` | 🔒 | Create a user with optional detail. |
| GET | `/api/v1/users/me` | 🔒 | Profile of the authenticated user. |
| GET | `/api/v1/users/:id` | 🔒 | Fetch a single user with detail. |
| PUT | `/api/v1/users/:id` | 🔒 | Replace name and detail (full update). |
| PATCH | `/api/v1/users/:id` | 🔒 | Partial update (only provided fields). |
| DELETE | `/api/v1/users/:id` | 🔒 | Delete a user (cascades to detail). |

---

## Examples

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

**Partial update**

```json
PATCH /api/v1/users/<id>
Authorization: Bearer <token>
{ "name": "Aldy D", "detail": { "city": "Jakarta", "phone": "08123" } }
```

---

## Data Models

### User

| Field | Type | Notes |
|-------|------|-------|
| `id` | uuid | Primary key. |
| `email` | string | Unique. |
| `name` | string | |
| `detail` | UserDetail | One-to-one, optional. |
| `created_at` / `updated_at` | timestamp | |

### UserDetail

| Field | Type | Notes |
|-------|------|-------|
| `phone` | string | |
| `address` | string | |
| `city` | string | |
| `country` | string | |
| `bio` | string | |
| `avatar_url` | string | Must be a valid URL. |

---

## Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| `BAD_REQUEST` | 400 | Malformed body. |
| `UNAUTHORIZED` | 401 | Missing/invalid token or credentials. |
| `NOT_FOUND` | 404 | Resource does not exist. |
| `CONFLICT` | 409 | Duplicate (e.g. email already registered). |
| `VALIDATION_ERROR` | 422 | Field validation failed. |
| `INTERNAL` | 500 | Unexpected server error. |
