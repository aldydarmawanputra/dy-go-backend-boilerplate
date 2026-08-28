# Postman collection

Import `dy-go-backend-boilerplate.postman_collection.json` ke Postman.

## Cara pakai
1. Set variable `base_url` kalau bukan `http://localhost:8080` (Collection → Variables).
2. **Register** → **Login**. Login otomatis nyimpen `access_token` + `refresh_token` ke collection variables, jadi request lain langsung keautentikasi (Bearer).
3. Request lain tinggal jalanin.

## Catatan
- **Admin-only** (`Create user`, `Delete user`, `Create API key`, `Broadcast`): butuh user dengan role `admin`. Kasih role admin lewat SQL (lihat README utama, bagian *Roles & authorization*), login ulang biar token bawa role admin.
- **Verify email / Reset password**: token dikirim via email. Pas dev pakai **Mailpit** (`http://localhost:8025`) buat ambil token, tempel ke variable `verify_token` / `reset_token`.
- **File upload**: pilih file di tab Body (form-data, key `file`).
- **API key**: `Create (admin)` otomatis nyimpen plaintext key ke variable `api_key`; dipakai `Whoami` lewat header `X-API-Key`.
- **Webhook**: kalau `PAYMENT_WEBHOOK_SECRET` diisi, isi header `X-Signature` = HMAC-SHA256(body, secret) hex.
- **WebSocket**: buat New → WebSocket Request ke `ws://localhost:8080/ws`.
