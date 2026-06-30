# Soal 2 - Technical Plan

## 1. Tujuan Technical Plan

Dokumen ini menerjemahkan hasil PRD analysis menjadi rencana teknis backend untuk Pocket App. Plan ini menjadi dasar untuk database design, payload contract, implementasi API, validasi, error handling, konfigurasi, dan testing.

Fokus MVP backend:

- Menyediakan API untuk auth mock, dashboard summary, dan pocket item management.
- Menjamin data pocket item selalu scoped ke tenant database milik user yang login.
- Menjalankan validasi business rule di backend, bukan hanya di frontend.
- Menyediakan response dan error contract yang stabil untuk integrasi frontend.
- Mengikuti struktur backend Go yang sudah ada di repo: domain, application, adapter, middleware, formatter, validator, dan config.

## 2. Tech Stack dan Alasan Pemilihan

| Komponen | Pilihan | Alasan |
| --- | --- | --- |
| Bahasa | Go | Cocok untuk backend API yang ringan, cepat, strongly typed, dan sudah digunakan di repo |
| HTTP framework | Fiber | Sudah tersedia di repo, routing sederhana, performa baik, middleware sudah ada |
| Database | PostgreSQL | Cocok untuk database-per-tenant, query filter/sort/pagination, dan indexing |
| Data access | sqlx + pgx | Sudah tersedia di repo, tetap dekat dengan SQL sehingga query search/filter lebih eksplisit |
| Query builder | Squirrel | Sudah tersedia di repo, membantu compose query dinamis untuk search/filter/sort |
| Validation | go-playground/validator | Sudah tersedia dan dipakai oleh package validator internal |
| Auth token | JWT | Sudah tersedia, cukup untuk mock auth dan bisa dikembangkan untuk production |
| Config | Viper + env/config.yaml | Sudah tersedia, mendukung konfigurasi lokal dan environment-based override |
| ID generator | google/uuid | Sudah tersedia, cocok untuk identifier user dan pocket item |
| Testing | Go testing + Testify + gomock | Sudah tersedia, mendukung unit test service dan repository mock |
| Migration/seed | Seedapp/scripts migration yang ada atau migration package serupa | Repo sudah punya pola migration/seed terpisah |

Catatan scope:

- Redis belum wajib untuk MVP Pocket karena belum ada kebutuhan cache/session store yang kuat.
- gRPC tidak dibutuhkan untuk API Pocket MVP kecuali activity log ingin diintegrasikan.
- OAuth, file upload, AI summary, full text extraction, dan collaboration tetap out of scope.

## 3. Microservices Architecture & Multitenancy

### 3.1. Daftar Microservices

Sistem di dalam repositori ini dibangun menggunakan pendekatan microservices, di mana setiap service memiliki tanggung jawab spesifik:

1. **`userapp` (Main App)**: Aplikasi utama (backend API) yang langsung berinteraksi dengan pengguna/frontend. Menangani autentikasi, manajemen pocket items, dan dashboard summary.
2. **`agentapp`**: Service berbasis Node.js yang bertanggung jawab untuk menghasilkan AI summary dari URL yang diberikan pengguna, serta menangani background processing dan task asinkron lainnya.
3. **`monitorapp`**: Service yang menangani pencatatan log aktivitas pengguna (user activity log), serta observabilitas seperti agregasi metrik dan monitoring kesehatan sistem.
4. **`seedapp`**: Service utilitas yang bertugas khusus untuk menjalankan database migration (schema update) dan seeding data awal, baik untuk main database maupun tenant database.

### 3.2. Konsep Multitenancy (Database-per-Tenant)

Sistem ini menerapkan arsitektur **Multitenant** dengan pendekatan **Database-per-Tenant**. Pendekatan ini mengisolasi data setiap pengguna secara fisik di level database.

**Cara Kerja Multitenant:**
1. **Main Database (Global)**: Menyimpan data terpusat seperti registrasi pengguna (users), kredensial login, dan *mapping* antara pengguna dengan tenant database mereka (`tenantCode`).
2. **Tenant Database (Isolated)**: Menyimpan business data yang spesifik milik satu pengguna, seperti `pocket_items`. Setiap pengguna memiliki database-nya sendiri.
3. **Tenant Resolver Workflow**:
   - Pengguna melakukan request HTTP dengan menyertakan JWT token. Token ini mengandung `tenantCode`.
   - Middleware (misal: Tenant Middleware) mengekstrak `tenantCode` dari request dan memanggil **Tenant Resolver**.
   - Tenant Resolver akan memilih pool koneksi database PostgreSQL yang sesuai untuk tenant tersebut dan meletakkannya di dalam `context` request.
   - Layanan (`application/service`) dan repositori (`adapter/repository/database`) selanjutnya hanya menjalankan query SQL menggunakan koneksi yang ada di `context` tersebut, memastikan data selalu scope ke tenant yang benar.

**Keuntungan Pendekatan Database-per-Tenant:**
- **Isolasi Data (Security)**: Mencegah risiko kebocoran data antar pengguna karena tidak perlu klausa `WHERE tenant_id = ?` di setiap query.
- **Skalabilitas**: Database pengguna yang *heavy-usage* dapat dialokasikan ke resource server yang lebih besar secara independen.

### 3.3. Layered Architecture (khususnya pada `userapp`)

Arsitektur internal untuk `userapp` mengikuti pola layered/clean-ish architecture:

```text
HTTP Request
  -> middleware (Auth & Tenant Context)
  -> tenant database resolver
  -> application/api (Handler)
  -> application/service (Use Case)
  -> domain (Repository Interface)
  -> adapter/repository/database (SQL Implementation)
  -> PostgreSQL (Tenant DB)
```

Tanggung jawab tiap layer:

| Layer | Tanggung Jawab |
| --- | --- |
| `adapter/rest` | Setup Fiber, global middleware, error mapper |
| `pkg/custommiddleware` | Auth, logging, request id, recover, global error handling |
| `application/api` | Parse request, bind DTO, call validator, call service, format response |
| `application/service` | Business use case, tenant-aware authorization, default value, orchestration |
| `domain/*` | Entity, enum, request/response model, repository interface, domain error |
| `adapter/repository/database` | SQL implementation, query builder, transaction, mapping DB row |
| `config` | Load config file/env and expose typed config |

Prinsip dependency:

- Handler bergantung ke service interface.
- Service bergantung ke repository interface.
- Repository implementation bergantung ke database client.
- Domain tidak bergantung ke Fiber, SQL, atau infrastructure package.

## 4. Struktur Direktori

Secara arsitektur, sistem ini dirancang agar setiap microservice nantinya dapat dipisah menjadi repositori mandiri (*1 service per repo*).

Pembagian service utamanya adalah:

```text
agentapp/        # Service untuk background processing & task asinkron
monitorapp/      # Service untuk observabilitas dan monitoring
seedapp/         # Service khusus untuk migrasi database dan seeding data
userapp/         # Aplikasi backend utama (REST API)
```

Di dalam masing-masing service (contohnya `userapp`), struktur direktori di dalamnya dijabarkan lebih mendetail dengan menggunakan kerangka *Clean Architecture* dan *Dependency Injection*:

```text
userapp/
├── cmd/                   # Command-line entry points
│   ├── cmd.go             # Setup environment/command awal
│   └── service.go         # Implementasi command untuk menjalankan API service
├── config/                # Configuration management
│   ├── model.go           # Struktur data dari konfigurasi aplikasi
│   └── config.yaml        # File konfigurasi default
├── internal/
│   ├── adapter/           # Infrastructure adapters
│   │   ├── repository/    # Data access layer
│   │   │   ├── database/  # Implementasi repositori ke database (PostgreSQL, dll)
│   │   │   └── cache/     # Implementasi repositori cache (Redis, dll)
│   │   └── rest/          # Konfigurasi REST (Fiber app, error mapper)
│   ├── application/       # Application layer
│   │   ├── api/           # HTTP API handler (controller/route handler)
│   │   ├── service/       # Business logic / Use cases
│   │   └── api.go         # Registrasi seluruh route API
│   ├── bootstrap/         # Dependency Injection / App bootstraping
│   │   ├── bootstrap.go   # Entry point perakitan dependencies
│   │   ├── adapter.go     # Register modul-modul infrastructure
│   │   └── application.go # Register layer application dan service
│   ├── domain/            # Domain layer (core business rules, murni tanpa infrastructure)
│   │   ├── auth/          # Domain untuk autentikasi
│   │   ├── shared/        # Entitas dan konstanta yang di-*share* secara umum
│   │   ├── pocket/        # (Modul Baru) Entitas, enum, repository interface Pocket
│   │   └── dashboard/     # (Modul Baru) Entitas Dashboard
│   └── pkg/               # Shared utility packages
│       ├── validator/     # Fungsi validasi request payload
│       ├── formatter/     # Fungsi untuk standarisasi format response sukses/error
│       └── custommiddleware/ # Middleware internal (logging, auth, error handling)
├── main.go                # Application entry point utama
├── Dockerfile             # Konfigurasi build *image* Docker
├── Makefile               # Task automation (build, run, test)
└── go.mod                 # Daftar dependensi modul Go
```

Untuk fitur Pocket, pengembangan domain logic, service, dan handler API akan terjadi sepenuhnya di dalam `userapp`. Namun, khusus untuk pengelolaan skema database (file migration) dan inisialisasi data awal (seed), prosesnya dikelola dan dieksekusi secara terpusat melalui service **`seedapp`**.

## 5. Module dan Domain Boundary

### 5.1 Auth Module

Scope:

- Login mock atau login existing user.
- Resolve tenant database yang boleh diakses user.
- Generate access token.
- Middleware auth membaca token dan mengisi context/locals user serta tenant context.
- Logout jika session/token store tersedia.

Boundary:

- Auth tidak mengatur pocket item.
- Auth menyediakan identity user dan tenant context untuk modul lain.

### 5.2 User Module

Scope:

- Menyimpan data user minimal di auth/main database.
- Menyimpan mapping user ke tenant database.
- Digunakan oleh auth dan tenant database resolver.

Boundary:

- User module tidak menghitung summary pocket.
- User module tidak melakukan search pocket.

### 5.3 Pocket Module

Scope:

- Create, list, detail, update, archive, update status, toggle favorite.
- Search/filter/sort/pagination.
- Validasi rule pocket item.
- Authorization item ownership.

Boundary:

- Pocket module berjalan menggunakan tenant database connection dari auth/tenant context.
- Pocket module tidak melakukan metadata extraction dari URL.
- Pocket module tidak menyimpan file.

### 5.4 Dashboard Module

Scope:

- Menghasilkan aggregate summary dari pocket item milik user.
- Menghasilkan recent items.

Boundary:

- Dashboard read-only.
- Dashboard tidak menyimpan state sendiri pada MVP.

### 5.5 Shared Module

Scope:

- Pagination model.
- Shared error type.
- Common response model jika dibutuhkan.
- Helper enum/identity jika sudah ada.

Boundary:

- Shared module tidak boleh menjadi tempat business logic spesifik pocket.

## 6. Layer Detail

### 6.1 Handler/API Layer

File:

- `internal/application/api/pocket.go`
- `internal/application/api/dashboard.go`

Tanggung jawab:

- Register route.
- Parse path/query/body.
- Ambil tenant context dari auth middleware.
- Jalankan validator untuk request body/query.
- Panggil service.
- Return response dengan `formatter.NewSuccessResponse`.

Contoh route:

```http
GET    /api/pockets
GET    /api/pockets/archived
GET    /api/pockets/:id
POST   /api/pockets
PUT    /api/pockets/:id
DELETE /api/pockets/:id
PATCH  /api/pockets/:id/status
PATCH  /api/pockets/:id/favorite
GET    /api/dashboard/summary
```

Handler tidak boleh:

- Menulis SQL.
- Mengandung business rule kompleks.
- Memilih database tenant sendiri tanpa melalui tenant resolver.

### 6.2 Service/Use Case Layer

File:

- `internal/application/service/pocket.go`
- `internal/application/service/dashboard.go`

Tanggung jawab:

- Menjalankan business rule.
- Normalize input: trim title, trim tags, remove empty tags.
- Set default: `status=unread`, `isFavorite=false`.
- Validate rule lintas field, misalnya URL wajib jika content type bukan note.
- Memastikan mutation berjalan pada tenant database milik user aktif.
- Memanggil repository.
- Mapping domain error yang jelas.

Use case Pocket:

| Use Case | Input | Output |
| --- | --- | --- |
| `ListPockets` | tenant DB, query | paginated pocket items |
| `ListArchivedPockets` | tenant DB, query | paginated archived items |
| `GetPocketDetail` | tenant DB, pocketID | pocket item |
| `CreatePocket` | tenant DB, request | created pocket item |
| `UpdatePocket` | tenant DB, pocketID, request | updated pocket item |
| `ArchivePocket` | tenant DB, pocketID | archived pocket item or nil |
| `UpdatePocketStatus` | tenant DB, pocketID, status | updated pocket item |
| `TogglePocketFavorite` | tenant DB, pocketID, isFavorite | updated pocket item |

### 6.3 Repository/Data Access Layer

File:

- `internal/domain/pocket/repository.go`
- `internal/adapter/repository/database/pocket.go`
- `internal/domain/dashboard/repository.go`
- `internal/adapter/repository/database/dashboard.go`

Tanggung jawab:

- Menjalankan query SQL.
- Compose search/filter/sort/pagination.
- Menjamin repository berjalan pada tenant database connection yang benar.
- Mapping row database ke domain model.
- Return domain error untuk not found.

Repository interface contoh:

```go
type Repository interface {
    Find(ctx context.Context, filter ListFilter) ([]PocketItem, PageMeta, error)
    FindArchived(ctx context.Context, filter ListFilter) ([]PocketItem, PageMeta, error)
    FindByID(ctx context.Context, id string) (*PocketItem, error)
    Create(ctx context.Context, item *PocketItem) error
    Update(ctx context.Context, item *PocketItem) error
    Archive(ctx context.Context, id string) error
    UpdateStatus(ctx context.Context, id string, status Status) (*PocketItem, error)
    UpdateFavorite(ctx context.Context, id string, isFavorite bool) (*PocketItem, error)
}
```

## 7. Database Design Plan

### 7.1 `pocket_items`

```sql
CREATE TABLE pocket_items (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  url TEXT NULL,
  description TEXT NULL,
  content_type VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'unread',
  is_favorite BOOLEAN NOT NULL DEFAULT FALSE,
  tags JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  archived_at TIMESTAMPTZ NULL
);
```

Constraint yang disarankan:

```sql
ALTER TABLE pocket_items
  ADD CONSTRAINT pocket_items_content_type_check
  CHECK (content_type IN ('article', 'video', 'document', 'note'));

ALTER TABLE pocket_items
  ADD CONSTRAINT pocket_items_status_check
  CHECK (status IN ('unread', 'reading', 'read', 'archived'));
```

Index yang disarankan:

```sql
CREATE INDEX idx_pocket_items_active_created_at ON pocket_items(created_at DESC) WHERE archived_at IS NULL;
CREATE INDEX idx_pocket_items_status_active ON pocket_items(status, created_at DESC) WHERE archived_at IS NULL;
CREATE INDEX idx_pocket_items_content_type_active ON pocket_items(content_type, created_at DESC) WHERE archived_at IS NULL;
CREATE INDEX idx_pocket_items_favorite_active ON pocket_items(is_favorite, created_at DESC) WHERE archived_at IS NULL;
CREATE INDEX idx_pocket_items_archived_at ON pocket_items(archived_at DESC) WHERE archived_at IS NOT NULL;
```

Search MVP:

- Gunakan `ILIKE` pada `title`, `url`, `description`, dan cast `tags::text`.
- Jika data membesar, tambahkan `tsvector`/GIN index.

### 7.2 Archive Model

Untuk MVP, gunakan dua indikator:

- `status = 'archived'`
- `archived_at IS NOT NULL`

Alasan:

- `status` selaras dengan PRD/glossary.
- `archived_at` berguna untuk audit dan sorting archive.

Rule konsistensi:

- List utama: `archived_at IS NULL`.
- Archive list: `archived_at IS NOT NULL`.
- Saat archive: set `status='archived'`, `archived_at=NOW()`, `updated_at=NOW()`.

## 8. Payload Contract Plan

### 8.1 Response Wrapper

Mengikuti formatter existing:

```json
{
  "status": "success",
  "message": "success",
  "data": {}
}
```

Error response:

```json
{
  "status": "bad_request",
  "message": "Validation error",
  "traceId": "request-id",
  "errorList": {
    "title": "title is required"
  }
}
```

Jika ingin lebih dekat dengan PRD, bisa ditambahkan `code` pada response error di fase berikutnya. Untuk MVP, gunakan mapping status formatter existing agar konsisten dengan repo.

### 8.2 Create/Update Pocket Request

```json
{
  "title": "React Performance Guide",
  "url": "https://example.com/react-performance",
  "description": "A guide about React rendering optimization",
  "contentType": "article",
  "tags": ["frontend", "react"]
}
```

Validation:

- `title`: required, trim, max 255.
- `url`: required jika `contentType` bukan `note`, harus valid URL.
- `description`: optional, max length disarankan 2000.
- `contentType`: required, enum.
- `tags`: optional, max item disarankan 20, no duplicate case-insensitive.

### 8.3 List Query

```http
GET /api/pockets?search=react&status=unread&contentType=article&favorite=true&page=1&limit=10&sort=createdAt:desc
```

Validation:

- `page`: integer, default 1, minimum 1.
- `limit`: integer, default 10, minimum 1, maximum 100.
- `status`: enum unread/reading/read.
- `contentType`: enum article/video/document/note.
- `favorite`: boolean.
- `sort`: whitelist `createdAt:desc`, `createdAt:asc`, `title:asc`, `title:desc`.

### 8.4 Paginated Response

```json
{
  "status": "success",
  "message": "success",
  "data": {
    "items": [],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 0,
      "totalPages": 0
    }
  }
}
```

## 9. Validation Strategy

Validation dibagi menjadi dua level.

### 9.1 Struct Validation

Menggunakan `go-playground/validator` via `internal/pkg/validator`.

Contoh rule:

- Required field.
- Email format.
- Max length.
- Enum melalui custom validation atau service-level validation.
- Basic URL format.

### 9.2 Business Validation

Dijalankan di service layer.

Rule:

- Title setelah trim tidak boleh kosong.
- URL wajib untuk article/video/document.
- URL optional untuk note.
- URL hanya boleh `http` atau `https`.
- Tag duplicate dicek case-insensitive.
- Search whitespace diabaikan.
- Status archived tidak boleh dipakai melalui update reading status endpoint.
- Favorite/status update pada archived item perlu ditentukan. Rekomendasi MVP: reject dengan `POCKET_ARCHIVED`.

## 10. Error Handling Strategy

Gunakan global error handler existing:

- Domain/service mengembalikan typed error.
- `adapter/rest/error_mapper.go` memetakan error ke formatter status dan HTTP status.
- Middleware menambahkan `traceId`.
- Response tidak mengekspos stack trace.

Domain error yang disarankan:

| Error | HTTP Status | Kapan Terjadi |
| --- | --- | --- |
| `ErrPocketNotFound` | 404 | Item tidak ditemukan di tenant database aktif |
| `ErrInvalidContentType` | 400 | Content type tidak valid |
| `ErrInvalidStatus` | 400 | Status tidak valid |
| `ErrURLRequired` | 400 | URL kosong untuk non-note |
| `ErrInvalidURL` | 400 | URL format/protocol tidak valid |
| `ErrDuplicateTag` | 400 | Tag duplicate dalam satu item |
| `ErrPocketArchived` | 409 | Mutation tidak valid pada item archived |
| `ErrUnauthorized` | 401 | Token/session tidak valid |

Prinsip:

- Akses tenant database yang tidak dimiliki user ditolak sebelum repository dijalankan.
- Validation error field-level dikembalikan dalam `errorList`.
- Internal error dicatat di log tetapi response tetap generic.

## 11. Configuration dan Environment Strategy

Gunakan pola config existing:

- `config.yaml` untuk default lokal.
- `.env` untuk override secret dan environment-specific value.
- `config/model.go` untuk typed config.

Config tambahan yang disarankan:

```yaml
pocket:
  default_page_limit: 10
  max_page_limit: 100
  max_title_length: 255
  max_description_length: 2000
  max_tags_per_item: 20
  allowed_url_schemes: http,https

auth:
  mock_email: "user@example.com"
  mock_password: "password123"
  access_token_ttl: "24h"

tenant:
  mode: "database_per_tenant"
  default_tenant_code: "pocket_demo_tenant"
```

Environment penting:

| Key | Keterangan |
| --- | --- |
| `APP_ENV` | development/test/production |
| `HTTP_PORT` | Port API |
| `DATABASE_HOST` | Host PostgreSQL |
| `DATABASE_PORT` | Port PostgreSQL |
| `DATABASE_NAME` | Nama database |
| `DATABASE_USER` | User database |
| `DATABASE_PASSWORD` | Password database |
| `TENANT_MODE` | Mode tenant, default `database_per_tenant` |
| `JWT_SECRET` | Secret token |
| `INTERNAL_API_KEY` | API key internal jika route tertentu butuh proteksi |

Security config:

- Secret tidak boleh di-commit untuk production.
- `.env.example` perlu berisi key tanpa secret nyata.
- JWT secret wajib berbeda per environment.

## 12. Authentication dan Authorization Plan

MVP auth:

- Login memakai credential mock atau user seeded.
- Generate JWT berisi `userID`, `email`, `tenantCode`, dan expiry.
- Auth middleware membaca `Authorization: Bearer <token>`.
- Middleware mengisi `user_id`, `email`, dan `tenant_code` ke `fiber.Ctx.Locals`.

Authorization:

- Handler mengambil tenant context dari middleware.
- Tenant resolver memilih database connection berdasarkan `tenant_code`.
- Service/repository hanya berjalan pada tenant database yang sudah dipilih.
- Query detail/update/archive cukup memakai `WHERE id = ?` di dalam tenant database aktif.
- Jika row tidak ditemukan, return `ErrPocketNotFound`.

## 13. Search, Filter, Sort Implementation Plan

Search:

- Trim keyword.
- Jika kosong, skip condition.
- Query:

```sql
title ILIKE :keyword
OR url ILIKE :keyword
OR description ILIKE :keyword
OR tags::text ILIKE :keyword
```

Filter:

- Tambahkan condition hanya jika parameter terisi.
- List utama selalu menambahkan `archived_at IS NULL`.
- Archive list selalu menambahkan `archived_at IS NOT NULL`.

Sort:

- Parse `field:direction`.
- Whitelist field dan direction.
- Default `created_at DESC`.
- Mapping API field ke DB column:

| API Field | DB Column |
| --- | --- |
| `createdAt` | `created_at` |
| `title` | `title` |

Pagination:

- `limit` dan `offset`.
- Query count terpisah untuk `total`.
- Response meta menghitung `totalPages`.

## 14. Logging dan Observability

Gunakan middleware existing:

- Request ID/trace ID.
- Request logging.
- Recover middleware.
- Error handler dengan traceId.

Yang perlu dicatat:

- Method, path, status code, duration.
- Trace ID.
- Error internal untuk debugging.
- Jangan log password, token, atau secret.

Untuk MVP, metrics/tracing terdistribusi belum wajib.

## 15. Testing Strategy

### 15.1 Unit Test Domain/Service

Target:

- Validasi title kosong.
- Validasi URL wajib untuk article/video/document.
- Note tanpa URL valid.
- Duplicate tag ditolak.
- Default status unread dan favorite false.
- Search whitespace dinormalisasi.
- Archived item tidak bisa muncul di list utama.
- User tidak bisa mengakses tenant database lain.

Tools:

- `testing`
- `testify`
- `gomock` untuk mock repository.

### 15.2 Handler/API Test

Target:

- Body parse error.
- Validation error format.
- Auth required.
- Response success create/list/detail.
- Query parameter validation.
- Error mapping ke status HTTP.

Gunakan Fiber app test dengan injected service mock.

### 15.3 Repository Test

Target:

- Query list dengan search/filter/sort/pagination.
- Query detail berjalan pada tenant database aktif.
- Archive update `status` dan `archived_at`.
- Count pagination benar.

Opsi:

- Integration test dengan test database PostgreSQL.
- Jika test DB belum tersedia, minimal repository unit test untuk query builder dan mapping.

### 15.4 End-to-End API Scenario

Minimal scenario:

1. Login berhasil.
2. Create pocket article valid.
3. Create note tanpa URL berhasil.
4. Create article tanpa URL gagal.
5. List pocket menampilkan item aktif.
6. Search/filter/sort bekerja.
7. Detail item berhasil.
8. Update item berhasil.
9. Toggle favorite berhasil.
10. Update status berhasil.
11. Archive item berhasil dan hilang dari list utama.
12. Archived item muncul di archive list.
13. Akses tanpa token gagal.
14. Akses tenant lain ditolak oleh tenant resolver.

## 16. Implementation Roadmap

### Phase 1 - Foundation

- Tambahkan domain `pocket` dan `dashboard`.
- Tambahkan model request/response, enum, dan error.
- Tambahkan repository interface.
- Tambahkan config Pocket.

### Phase 2 - Database

- Buat migration `pocket_items`.
- Tambahkan index.
- Tambahkan seed sample pocket item.

### Phase 3 - Repository

- Implement create, update, detail, archive.
- Implement list search/filter/sort/pagination.
- Implement dashboard summary query.

### Phase 4 - Service

- Implement business validation.
- Implement default values dan normalization.
- Implement authorization flow berbasis tenant database resolver.

### Phase 5 - API

- Register routes.
- Implement handler.
- Integrasi validator dan formatter.
- Integrasi auth middleware.

### Phase 6 - Testing

- Unit test service.
- Handler test.
- Repository/integration test.
- Manual API test dengan sample payload.

### Phase 7 - Documentation

- Tambahkan API contract detail.
- Tambahkan cara menjalankan migration/seed.
- Tambahkan contoh request/response.

## 17. Risiko Teknis dan Mitigasi

| Risiko | Dampak | Mitigasi |
| --- | --- | --- |
| Query search/filter terlalu kompleks di handler | Handler sulit dirawat | Pindahkan query composition ke repository/filter model |
| Auth mock terlalu jauh dari production | Sulit upgrade | Tetap gunakan JWT dengan user identity dan tenant context |
| Archive model tidak konsisten | Item muncul di list yang salah | Gunakan `archived_at` sebagai source of truth list |
| Error response berbeda dari frontend expectation | Integrasi frontend terganggu | Standarkan formatter dan dokumentasikan contoh response |
| Tag JSON sulit di-query | Search tag kurang optimal | MVP pakai `tags::text ILIKE`, future pakai GIN/relasi |
| Pagination tidak diterapkan | Dashboard/list bisa lambat | Implement page/limit sejak awal |
| Validasi hanya di frontend | Data kotor masuk DB | Validasi ulang di service dan DB constraint |

## 18. Definition of Ready untuk Implementasi

Backend siap diimplementasikan jika:

- API endpoint final untuk MVP disepakati.
- Nama query parameter dipilih: gunakan `contentType`, bukan `type`.
- Archive model disepakati: `archived_at` + `status='archived'`.
- Error response mengikuti formatter existing repo.
- Migration location disepakati: `userapp` atau `seedapp`.
- Mock auth credential atau seeded user tersedia.
- Mapping user ke tenant database tersedia.

## 19. Kesimpulan

Technical plan backend Pocket App akan memakai stack dan pola yang sudah tersedia di repo: Go, Fiber, PostgreSQL, sqlx/pgx, validator, JWT, Viper, layered architecture, domain repository interface, service use case, dan handler API.

Implementasi utama difokuskan pada modul `pocket` dan `dashboard`, dengan auth existing sebagai sumber identity dan tenant context. Dengan boundary dan layering ini, backend dapat mendukung requirement PRD secara realistis, mudah diuji, dan siap dikembangkan ke database design, payload contract detail, implementasi, serta testing pada soal berikutnya.
