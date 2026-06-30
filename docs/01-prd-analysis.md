# Soal 1 - PRD Analysis & Backend Requirement

## 1. Ringkasan Pemahaman PRD dari Sudut Pandang Backend

Pocket App adalah aplikasi personal pocket management untuk menyimpan, mengelola, mencari, membaca ulang, dan mengarsipkan item seperti artikel, link, video, dokumen referensi, dan catatan pribadi.

Dari sisi backend, PRD ini dapat dipahami sebagai sistem CRUD personal yang berpusat pada entitas `PocketItem`, dengan autentikasi sederhana, kepemilikan data per user, pencarian/filter/sort, soft delete melalui archive, serta validasi data yang konsisten antara frontend dan backend.

Core flow produk adalah:

```text
Save -> Organize -> Retrieve -> Act
```

Implikasi backend dari flow tersebut:

| Flow | Makna Produk | Kebutuhan Backend |
| --- | --- | --- |
| Save | User menyimpan item baru | API create item, validasi form, default value |
| Organize | User memberi type, tag, status, favorite | Penyimpanan metadata, update status/favorite/tag |
| Retrieve | User mencari dan memfilter item | API list dengan query search/filter/sort/pagination |
| Act | User membuka detail, edit, archive | API detail, update, archive/soft delete |

Walaupun PRD ditulis untuk take-home test frontend dan menyebut mock API, requirement yang tersedia sudah cukup untuk diturunkan menjadi kontrak backend nyata atau mock backend yang menyerupai service produksi.

## 2. Domain Backend yang Muncul

### 2.1 Auth Domain

Auth domain menangani login, logout, session/token, dan protected access.

Tanggung jawab backend:

- Memvalidasi email dan password.
- Menghasilkan token/session mock setelah login berhasil.
- Menolak request tanpa token/session valid.
- Menyediakan identitas user aktif untuk authorization.

Catatan:

- PRD menyebut mock authentication, sehingga implementasi MVP boleh menggunakan credential statis atau seeded user.
- Jika berkembang menjadi production, auth perlu diganti ke mekanisme nyata seperti password hashing, refresh token, session expiry, dan rate limiting.

### 2.2 User Domain

User domain merepresentasikan pemilik data.

Data minimal:

| Field | Keterangan |
| --- | --- |
| `id` | Identifier unik user |
| `name` | Nama user untuk tampilan |
| `email` | Email user untuk login |
| `createdAt` | Waktu user dibuat |

Business rule utama:

- Semua pocket item harus terkait dengan satu user.
- User hanya boleh melihat dan memodifikasi item miliknya sendiri.

### 2.3 Pocket Item Domain

Pocket Item adalah domain utama aplikasi.

Data minimal:

| Field | Keterangan |
| --- | --- |
| `id` | Identifier unik pocket item |
| `userId` | Pemilik item |
| `title` | Judul item, wajib |
| `url` | URL item, wajib untuk article/video/document |
| `description` | Deskripsi opsional |
| `contentType` | Jenis konten: article, video, document, note |
| `status` | Status baca: unread, reading, read, archived |
| `isFavorite` | Penanda item penting |
| `tags` | Daftar tag manual |
| `createdAt` | Waktu item dibuat |
| `updatedAt` | Waktu item terakhir diubah |
| `archivedAt` | Waktu item diarsipkan, opsional tetapi disarankan |

Catatan desain backend:

- `archivedAt` tidak eksplisit di sample PRD, tetapi berguna untuk membedakan item aktif dan archived secara audit-friendly.
- `status = archived` disebut di glossary, tetapi business rule juga menyebut archive/soft delete. Backend perlu memilih satu model konsisten: status archived, field `archivedAt`, atau keduanya.

### 2.4 Tag Domain

Untuk MVP, tag dapat disimpan sebagai array string pada pocket item.

Alternatif backend:

| Opsi | Kelebihan | Kekurangan |
| --- | --- | --- |
| Array/string JSON pada item | Simpel untuk MVP dan mock API | Sulit query/index jika data besar |
| Tabel relasi `pocket_item_tags` | Lebih normal dan scalable | Implementasi lebih kompleks |

Rekomendasi MVP:

- Simpan tag sebagai array string atau tabel sederhana tergantung database.
- Tetap validasi duplicate tag per item.
- Normalisasi tag dengan trim dan case-insensitive comparison.

### 2.5 Dashboard/Summary Domain

Dashboard membutuhkan agregasi data milik user.

Data summary:

| Field | Keterangan |
| --- | --- |
| `totalItems` | Jumlah item aktif |
| `totalUnread` | Jumlah item aktif berstatus unread |
| `totalReading` | Jumlah item aktif berstatus reading |
| `totalRead` | Jumlah item aktif berstatus read |
| `totalFavorite` | Jumlah item aktif yang favorite |
| `recentItems` | Beberapa item terbaru |

Backend dapat menghitung summary secara real-time dari tabel pocket item. Untuk skala MVP, belum perlu materialized view/cache.

## 3. Fitur yang Membutuhkan API

### 3.1 Authentication API

| Method | Endpoint | Kebutuhan |
| --- | --- | --- |
| `POST` | `/api/auth/login` | Login dengan email/password mock |
| `POST` | `/api/auth/logout` | Logout/invalidate session jika backend menyimpan session |
| `GET` | `/api/auth/me` | Mengambil user aktif dari token/session |

PRD eksplisit menyebut login API, tetapi `logout` dan `me` disarankan agar protected route dan refresh page bisa bekerja konsisten.

### 3.2 Pocket Item API

| Method | Endpoint | Kebutuhan |
| --- | --- | --- |
| `GET` | `/api/pockets` | List item aktif dengan search/filter/sort/pagination |
| `GET` | `/api/pockets/:id` | Detail pocket item |
| `POST` | `/api/pockets` | Create pocket item |
| `PUT` | `/api/pockets/:id` | Update pocket item |
| `DELETE` | `/api/pockets/:id` | Archive/soft delete item |
| `PATCH` | `/api/pockets/:id/status` | Update reading status |
| `PATCH` | `/api/pockets/:id/favorite` | Toggle/set favorite |

### 3.3 Archive API

PRD menyebut archived item tampil di archive page, tetapi API contract hanya menampilkan `DELETE /api/pockets/:id`.

Backend perlu endpoint tambahan:

| Method | Endpoint | Kebutuhan |
| --- | --- | --- |
| `GET` | `/api/pockets/archived` | List item yang sudah archived |
| `PATCH` | `/api/pockets/:id/restore` | Restore item archived, jika fitur restore diinginkan |

Jika restore belum masuk MVP, endpoint restore dapat ditunda.

### 3.4 Dashboard API

| Method | Endpoint | Kebutuhan |
| --- | --- | --- |
| `GET` | `/api/dashboard/summary` | Total item, unread, reading, read, favorite, recently added |

Alternatif:

- Frontend menghitung summary dari `/api/pockets`.
- Backend menyediakan summary khusus agar lebih efisien dan kontrak dashboard jelas.

Rekomendasi backend:

- Sediakan `/api/dashboard/summary`, terutama jika list memakai pagination.

## 4. Data yang Perlu Disimpan

### 4.1 Tabel/Collection `users`

| Field | Type | Constraint |
| --- | --- | --- |
| `id` | string/uuid | primary key |
| `name` | string | required |
| `email` | string | required, unique |
| `passwordHash` | string | optional untuk mock, required untuk real auth |
| `createdAt` | datetime | required |
| `updatedAt` | datetime | required |

Untuk mock auth, `passwordHash` dapat diganti credential hardcoded, tetapi model data sebaiknya tetap siap untuk real backend.

### 4.2 Tabel/Collection `pocket_items`

| Field | Type | Constraint |
| --- | --- | --- |
| `id` | string/uuid | primary key |
| `userId` | string/uuid | required, indexed |
| `title` | string | required, trimmed |
| `url` | string/null | required jika type bukan note |
| `description` | string/null | optional |
| `contentType` | enum | article/video/document/note |
| `status` | enum | unread/reading/read/archived |
| `isFavorite` | boolean | default false |
| `tags` | string array/json | optional, no duplicate per item |
| `createdAt` | datetime | required |
| `updatedAt` | datetime | required |
| `archivedAt` | datetime/null | optional |

Index yang disarankan:

- `userId`
- `userId, createdAt`
- `userId, status`
- `userId, contentType`
- `userId, isFavorite`
- Full-text/search index untuk `title`, `url`, `description`, dan `tags` jika database mendukung.

### 4.3 Session/Token Storage

Untuk mock API, token bisa stateless dan sederhana. Untuk backend nyata, perlu salah satu:

- JWT dengan expiry.
- Session table/cache.
- Refresh token storage.

Minimal payload token:

```json
{
  "userId": "user_001",
  "email": "user@example.com"
}
```

## 5. Business Rule yang Perlu Divalidasi Backend

Backend tidak boleh hanya mengandalkan validasi frontend. Semua rule penting dari PRD perlu divalidasi ulang di API.

| Rule | Validasi Backend |
| --- | --- |
| User harus login | Semua endpoint pocket/dashboard wajib auth |
| User hanya melihat item miliknya | Query selalu dibatasi `userId` dari token |
| Title wajib | Reject title kosong atau hanya spasi |
| URL wajib untuk article/video/document | Reject jika type bukan note dan URL kosong |
| URL tidak wajib untuk note | Izinkan null/empty URL untuk note |
| URL harus valid | Validasi format URL dan protocol yang diizinkan |
| Tag opsional | Izinkan array kosong |
| Tag tidak boleh duplicate | Deduplicate atau reject duplicate case-insensitive |
| Status default item baru unread | Set default di backend |
| Favorite default false | Set default di backend |
| Archived item tidak tampil di list utama | Filter default exclude archived |
| Archived item tampil di archive page | Endpoint archive hanya return archived |
| Delete adalah archive/soft delete | Jangan hard delete pada MVP |
| Search mencakup title, URL, description, tags | Implement query di field tersebut |
| Search dan filter bisa dikombinasi | Query builder harus mendukung kombinasi parameter |
| Sort default createdAt desc | Default list tanpa sort |
| Prevent double submit | Backend sebaiknya idempotent atau aman dari duplicate cepat |

Catatan double submit:

- Frontend men-disable submit saat loading.
- Backend tetap perlu aman jika request ganda masuk, misalnya dengan optional idempotency key atau duplicate URL policy jika nanti ditentukan.

## 6. Search, Filter, Sort, dan Pagination

### 6.1 Search

Search dilakukan terhadap:

- `title`
- `url`
- `description`
- `tags`

Perilaku yang disarankan:

- Trim keyword.
- Jika keyword hanya spasi, anggap tidak ada search.
- Search case-insensitive.
- Untuk MVP, partial match sudah cukup.
- Untuk database relasional kecil, `LIKE/ILIKE` cukup.
- Untuk skala lebih besar, gunakan full-text index.

### 6.2 Filter

Parameter filter:

| Parameter | Nilai |
| --- | --- |
| `status` | unread, reading, read |
| `type` atau `contentType` | article, video, document, note |
| `favorite` | true/false |

Catatan:

- `archived` sebaiknya tidak dicampur dengan list utama kecuali ada parameter eksplisit seperti `archived=true`.
- Nama query parameter perlu distandarkan: PRD memakai `type` di API example, sedangkan data model memakai `contentType`.

### 6.3 Sort

Sort option dari PRD:

| Sort | Makna |
| --- | --- |
| `createdAt:desc` | Newest first, default |
| `createdAt:asc` | Oldest first |
| `title:asc` | Title A-Z |
| `title:desc` | Title Z-A |

Backend perlu whitelist field dan direction agar aman dari query injection.

### 6.4 Pagination

PRD menyebut pagination/load more sebagai could-have, namun API contract sudah memakai `page` dan `limit`.

Rekomendasi:

- Implement `page` dan `limit` sejak awal.
- Default `page=1`, `limit=10`.
- Batasi `limit` maksimum, misalnya 100.
- Response menyertakan `total`, `page`, `limit`, dan `totalPages`.

## 7. Error Handling dan Response Contract

PRD menyarankan wrapper:

```ts
type ApiResponse<T> = {
  data: T;
  message?: string;
};
```

Untuk error:

```ts
type ApiErrorResponse = {
  code: string;
  message: string;
  details?: {
    field?: string;
    message: string;
  }[];
};
```

Mapping error backend:

| Kondisi | HTTP Status | Code |
| --- | --- | --- |
| Input tidak valid | 400/422 | `VALIDATION_ERROR` |
| Login gagal | 401 | `INVALID_CREDENTIALS` |
| Token tidak ada/expired | 401 | `UNAUTHORIZED` |
| Akses data user lain | 403 atau 404 | `FORBIDDEN` atau `POCKET_NOT_FOUND` |
| Item tidak ditemukan | 404 | `POCKET_NOT_FOUND` |
| Konflik data | 409 | `CONFLICT` |
| Error server | 500 | `INTERNAL_ERROR` |

Rekomendasi security:

- Untuk item milik user lain, lebih aman return 404 agar tidak membocorkan keberadaan resource.
- Error response tidak boleh menampilkan stack trace.

## 8. Security dan Access Control

Kebutuhan security dari PRD yang berdampak ke backend:

- Protected endpoint wajib memeriksa token/session.
- Authorization per item wajib berdasarkan `userId`.
- Client validation tidak menggantikan server validation.
- User input tidak boleh dirender sebagai raw HTML.
- Error tidak boleh membocorkan stack trace.
- Token localStorage boleh untuk test, tetapi tradeoff harus dicatat.

Tambahan backend yang disarankan:

- Validasi allowed URL protocol, minimal `http` dan `https`.
- Normalize URL jika product memutuskan URL tanpa protocol boleh otomatis ditambahkan.
- Rate limit login jika auth nyata.
- Audit field `createdAt`, `updatedAt`, `archivedAt`.
- CORS hanya untuk origin frontend yang diizinkan jika API terpisah.

## 9. Risiko dan Ambiguitas Requirement

| Area | Ambiguitas/Risiko | Dampak Backend | Rekomendasi |
| --- | --- | --- | --- |
| Auth | Mock auth atau real backend belum final | Desain security bisa terlalu ringan atau terlalu berat | MVP pakai mock token, dokumentasikan upgrade path |
| API | PRD menyebut API/mock API | Bisa terjadi beda ekspektasi antara FE dan BE | Tetapkan contract OpenAPI atau README API |
| Delete | Delete disebut archive/soft delete, tetapi endpoint memakai `DELETE` | Salah implementasi hard delete | Implement `DELETE` sebagai soft delete |
| Archive model | Status punya `archived`, rule juga butuh archive page | Bisa double state antara `status` dan `archivedAt` | Pilih model konsisten dan dokumentasikan |
| Search | Client-side atau server-side masih open question | Pagination dan result count bisa berbeda | Untuk backend nyata, search server-side |
| Pagination | Disebut could-have, tetapi API memakai page/limit | FE bisa asumsi load all | Implement pagination default |
| URL tanpa protocol | PRD menyebut validasi atau normalisasi sesuai keputusan | Behavior create/edit tidak konsisten | Pilih reject atau normalize; rekomendasi reject untuk MVP |
| Metadata | Auto metadata extraction out of scope/open question | Backend tidak perlu crawler | Semua metadata manual |
| Duplicate URL | Belum ditentukan | User bisa menyimpan URL sama berkali-kali | MVP izinkan, future enhancement bisa deteksi duplicate |
| Tag casing | Duplicate tag tidak dijelaskan case-sensitive atau tidak | Data bisa kotor: `React` dan `react` | Validasi duplicate secara case-insensitive |
| Favorite/status race | Update cepat bisa konflik state | UI perlu rollback/refetch | Endpoint mutation return item terbaru |
| Archive restore | Archive page disebut, restore tidak disebut | Archive page hanya read-only | Tunda restore sampai ada requirement |

## 10. Asumsi Teknis Backend

Asumsi yang dibuat untuk menurunkan requirement:

- Aplikasi bersifat personal, bukan team collaboration.
- Satu pocket item hanya dimiliki oleh satu user.
- File upload, AI summary, auto tagging, full text extraction, dan public sharing tidak masuk MVP.
- Metadata title, description, type, dan tag diinput manual oleh user.
- Delete diperlakukan sebagai archive/soft delete.
- API list default hanya mengembalikan item aktif/non-archived.
- Backend tetap melakukan validasi walaupun frontend juga melakukan validasi.
- Search/filter/sort dilakukan server-side agar konsisten dengan pagination.
- ID dapat menggunakan UUID/string.
- Timezone penyimpanan menggunakan UTC.
- Response mutation mengembalikan data terbaru agar frontend mudah menjaga konsistensi state.

## 11. Prioritas Backend MVP

### Must Have

- Auth mock login dan protected endpoint.
- CRUD pocket item.
- Soft delete/archive.
- List item aktif.
- Detail item.
- Search by keyword.
- Validasi create/edit.
- Authorization per user.
- Error response standar.

### Should Have

- Dashboard summary endpoint.
- Filter status/content type/favorite.
- Sort newest/oldest/title.
- Update status endpoint.
- Toggle favorite endpoint.
- Archive list endpoint.
- Pagination.

### Could Have

- Restore archived item.
- Idempotency key untuk create.
- Full-text search index.
- Dedicated tag table.
- Duplicate URL detection.

## 12. Rekomendasi API Contract Ringkas

```http
POST   /api/auth/login
POST   /api/auth/logout
GET    /api/auth/me

GET    /api/dashboard/summary

GET    /api/pockets
GET    /api/pockets/archived
GET    /api/pockets/:id
POST   /api/pockets
PUT    /api/pockets/:id
DELETE /api/pockets/:id
PATCH  /api/pockets/:id/status
PATCH  /api/pockets/:id/favorite
```

Contoh query list:

```http
GET /api/pockets?search=react&status=unread&contentType=article&favorite=true&page=1&limit=10&sort=createdAt:desc
```

## 13. Kesimpulan

Dari sudut pandang backend, PRD Pocket App bukan hanya membutuhkan mock data, tetapi membutuhkan kontrak API dan model domain yang jelas agar frontend dapat menguji flow login, protected route, CRUD, search/filter/sort, favorite, status update, archive, loading/error state, dan validasi secara realistis.

Backend MVP sebaiknya difokuskan pada:

- Konsistensi data user-owned pocket item.
- Validasi rule inti di server.
- API contract yang stabil untuk integrasi frontend.
- Soft delete yang aman.
- Error handling yang predictable.
- Search/filter/sort yang konsisten dengan pagination.

Dengan pendekatan tersebut, kebutuhan frontend dalam PRD dapat terpenuhi tanpa memperbesar scope ke fitur non-MVP seperti OAuth nyata, file upload, AI summary, auto metadata extraction, collaboration, atau public sharing.
