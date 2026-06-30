# Soal 4 - Payload Contract

Dokumen ini berisi kontrak payload API untuk backend `userapp`. Kontrak ini menjadi acuan utama bagi implementasi endpoint dan integrasi dengan frontend.

## 1. Standar Format Response

### 1.1 Response Sukses (Umum)

```json
{
  "status": "00",
  "message": "success",
  "data": { ... }
}
```

### 1.2 Response Pagination (List)

```json
{
  "status": "00",
  "message": "success",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "totalPage": 5
  }
}
```

### 1.3 Response Error (Umum)

```json
{
  "code": "BAD_REQUEST",
  "message": "Deskripsi error"
}
```

### 1.4 Response Error Validasi (400)

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Validation error",
  "details": [
    {
      "field": "title",
      "message": "title is required"
    }
  ]
}
```

---

## 2. Kontrak Endpoint Pocket

Semua endpoint di bawah ini memerlukan otentikasi melalui *header* `Authorization: Bearer <token>`.

### 2.1 Create Pocket Item
Membuat item pocket baru.

- **URL**: `/api/pockets`
- **Method**: `POST`
- **Auth Required**: Yes

**Request Body:**

```json
{
  "title": "React Performance Guide",
  "url": "https://example.com/react-performance",
  "description": "A guide about React rendering optimization",
  "contentType": "article",
  "tags": ["frontend", "react"]
}
```

**Validation Rule:**
- `title`: `required`, string, length 3-120.
- `url`: opsional jika `contentType` adalah `note`, wajib dan format URL jika tipe lainnya.
- `description`: opsional, maksimal 500 karakter.
- `contentType`: `required`, enum (`article`, `video`, `document`, `note`).
- `tags`: opsional, maksimal 10 tags, tiap tag maksimal 24 karakter.

**Response Success (201 Created):**

```json
{
  "message": "Pocket item created successfully",
  "data": {
    "id": "uuid",
    "title": "React Performance Guide",
    "url": "https://example.com/react-performance",
    "description": "A guide about React rendering optimization",
    "contentType": "article",
    "status": "unread",
    "isFavorite": false,
    "tags": ["frontend", "react"],
    "createdAt": "2026-06-30T10:00:00Z",
    "updatedAt": "2026-06-30T10:00:00Z"
  }
}
```

### 2.2 Update Pocket Item
Memperbarui keseluruhan detail item pocket.

- **URL**: `/api/pockets/:id`
- **Method**: `PUT`
- **Path Parameter**: `id` (UUID dari item pocket)
- **Auth Required**: Yes

**Request Body:** Sama seperti Create.

**Validation Rule:** Sama seperti Create.

**Response Success (200 OK):** Sama seperti Create, dengan `message`: `"Pocket item updated successfully"`.

### 2.3 List Pocket Items
Mengambil daftar pocket item dengan dukungan pencarian, filter, pengurutan, dan pagination.

- **URL**: `/api/pockets`
- **Method**: `GET`
- **Auth Required**: Yes

**Query Parameters:**
- `search` (string): Mencari berdasarkan title, description, atau url.
- `status` (string): Filter berdasarkan status (`unread`, `reading`, `read`, `archived`).
- `type` (string): Filter berdasarkan contentType (`article`, `video`, `document`, `note`).
- `favorite` (boolean): Filter `true` atau `false` untuk status favorite.
- `page` (int): Halaman (default: 1).
- `limit` (int): Jumlah item per halaman (default: 10).
- `sort` (string): Pengurutan, format `field:direction` (contoh: `createdAt:desc`, `title:asc`).

**Response Success (200 OK):**

```json
{
  "data": [
    {
      "id": "uuid",
      "title": "React Performance Guide",
      "url": "https://example.com/react-performance",
      "description": "A guide about React rendering optimization",
      "contentType": "article",
      "status": "unread",
      "isFavorite": false,
      "tags": ["frontend", "react"],
      "createdAt": "2026-06-30T10:00:00Z",
      "updatedAt": "2026-06-30T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "totalPage": 1
  }
}
```

### 2.4 Get Pocket Detail
Mendapatkan detail dari sebuah pocket item berdasarkan ID.

- **URL**: `/api/pockets/:id`
- **Method**: `GET`
- **Path Parameter**: `id` (UUID)
- **Auth Required**: Yes

**Response Success (200 OK):**
```json
{
  "data": {
    "id": "uuid",
    "title": "React Performance Guide",
    "url": "https://example.com/react-performance",
    "description": "A guide about React rendering optimization",
    "contentType": "article",
    "status": "unread",
    "isFavorite": false,
    "tags": ["frontend", "react"],
    "createdAt": "2026-06-30T10:00:00Z",
    "updatedAt": "2026-06-30T10:00:00Z"
  }
}
```

### 2.5 Delete / Archive Pocket Item
Menghapus atau mengarsipkan pocket item. Sesuai desain, proses ini menggunakan metode *soft delete* (archiving).

- **URL**: `/api/pockets/:id`
- **Method**: `DELETE`
- **Path Parameter**: `id` (UUID)
- **Auth Required**: Yes

**Response Success (200 OK):**
```json
{
  "message": "Pocket item archived successfully",
  "data": {
    "id": "uuid"
  }
}
```

### 2.6 Update Pocket Status
Memperbarui field status baca (reading status) pada pocket item tertentu.

- **URL**: `/api/pockets/:id/status`
- **Method**: `PATCH`
- **Path Parameter**: `id` (UUID)
- **Auth Required**: Yes

**Request Body:**
```json
{
  "status": "read"
}
```
**Validation Rule:** `status` harus salah satu dari `unread`, `reading`, `read`, atau `archived`.

**Response Success (200 OK):**
```json
{
  "message": "Status updated successfully",
  "data": {
    "id": "uuid",
    "status": "read",
    "updatedAt": "2026-06-30T10:05:00Z"
  }
}
```

### 2.7 Toggle Favorite
Mengubah status favorit dari pocket item.

- **URL**: `/api/pockets/:id/favorite`
- **Method**: `PATCH`
- **Path Parameter**: `id` (UUID)
- **Auth Required**: Yes

**Request Body:**
```json
{
  "isFavorite": true
}
```

**Response Success (200 OK):**
```json
{
  "message": "Favorite updated successfully",
  "data": {
    "id": "uuid",
    "isFavorite": true,
    "updatedAt": "2026-06-30T10:06:00Z"
  }
}
```

### 2.8 Summarize Pocket
(Opsional) Mendapatkan rangkuman konten (AI Summary) untuk pocket item.

- **URL**: `/api/pockets/:id/summarize`
- **Method**: `POST`
- **Path Parameter**: `id` (UUID)
- **Auth Required**: Yes

**Response Success (200 OK):**
```json
{
  "message": "Summary generated successfully",
  "data": {
    "summary": "Ini adalah teks hasil ekstraksi atau rangkuman AI..."
  }
}
```

### 2.9 Get Dashboard Summary
Mendapatkan ringkasan statistik (summary) dari item pocket pengguna.

- **URL**: `/api/pockets/summary`
- **Method**: `GET`
- **Auth Required**: Yes

**Response Success (200 OK):**
```json
{
  "message": "Dashboard summary retrieved successfully",
  "data": {
    "totalItems": 150,
    "unreadItems": 50,
    "readingItems": 20,
    "readItems": 70,
    "archivedItems": 10,
    "favoriteItems": 35
  }
}
```

---
## 3. Kontrak Endpoint Auth

Semua endpoint di bawah ini terkait otentikasi dan pendaftaran.

### 3.1 Register Account
Mendaftarkan pengguna baru dan menghubungkannya dengan tenant.

- **URL**: `/api/auth/register`
- **Method**: `POST`
- **Auth Required**: API Key (`x-api-key`)

**Request Body:**
```json
{
  "name": "User Name",
  "email": "user@example.com",
  "password": "password123",
  "tenant_code": "optional-uuid-here",
  "tenant_name": "Tenant Name"
}
```

**Validation Rule:**
- `name`: `required`, maksimal 255 karakter.
- `email`: `required`, valid format email, maksimal 255 karakter.
- `password`: `required`, minimal 6 karakter.
- `tenant_code`: opsional, harus berformat UUID jika ada.
- `tenant_name`: `required`, maksimal 255 karakter.

**Response Success (201 Created):**
```json
{
  "status": "00",
  "message": "success",
  "data": {
    "user": {
      "id": "uuid",
      "name": "User Name",
      "tenant_default": "uuid",
      "active_indicator": "Y",
      "email": "user@example.com",
      "created_at": "2026-06-30T10:00:00Z",
      "updated_at": "2026-06-30T10:00:00Z"
    },
    "tenant_code": "uuid",
    "tenant_name": "Tenant Name"
  }
}
```

### 3.2 Login
Melakukan proses login untuk mendapatkan token JWT.

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Auth Required**: API Key (`x-api-key`)

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "time_zone": "Asia/Jakarta"
}
```

**Response Success (200 OK):**
```json
{
  "message": "Login success",
  "data": {
    "token": "jwt-token-string",
    "tenant_code": "uuid",
    "user": {
      "id": "uuid",
      "name": "User Name",
      "email": "user@example.com",
      "avatarUrl": ""
    }
  }
}
```

---
## 4. Daftar Status Code
Berikut adalah HTTP Status Code standar yang digunakan dalam implementasi API ini:

| HTTP Status Code | Deskripsi | Format Response |
| :--- | :--- | :--- |
| `200 OK` | Operasi berhasil. | JSON dengan `data` dan (opsional) `message`. |
| `201 Created` | Data baru berhasil dibuat. | JSON dengan `data` entitas baru. |
| `400 Bad Request` | Payload request tidak sesuai, salah format, atau gagal validasi. | JSON dengan format `VALIDATION_ERROR` & `details`. |
| `401 Unauthorized` | Autentikasi gagal (token tidak ada atau kadaluwarsa). | JSON error (`UNAUTHORIZED`). |
| `404 Not Found` | Entitas yang dituju (berdasarkan ID) tidak ditemukan di tenant aktif. | JSON error (`POCKET_NOT_FOUND`). |
| `500 Internal Server Error` | Kesalahan sistem (database error, dll). | JSON error (`INTERNAL_SERVER_ERROR`). |
