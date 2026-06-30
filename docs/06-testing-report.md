# Soal 7 - Testing Report

Dokumen ini merangkum hasil pengujian (testing) terhadap backend Pocket App (di dalam direktori `userapp`). Pengujian dilakukan untuk memastikan bahwa implementasi backend telah memenuhi kriteria yang ditetapkan dalam PRD, Technical Plan, Database Design, dan Payload Contract.

## 1. Unit Testing Service / Use Case

Pengujian pada level **Application Service** bertujuan memverifikasi logika bisnis (business logic) tanpa harus terhubung dengan database fisik (menggunakan *mock repository*). 

**Coverage Area:**
- `userapp/internal/application/service`

**Hasil Pengujian `TestPocketService`:**
- ✅ `TestCreate_Success`: Memverifikasi pembuatan pocket item berhasil jika payload valid.
- ✅ `TestUpdate_Success`: Memverifikasi pembaruan data berhasil.
- ✅ `TestUpdate_NotFound`: Memverifikasi error handling (`ErrPocketNotFound`) ketika item yang diupdate tidak ada.
- ✅ `TestFind_Success`: Memverifikasi pengambilan detail pocket item.
- ✅ `TestList_Success`: Memverifikasi pengambilan daftar item beserta meta datanya (pagination).
- ✅ `TestToggleFavorite_Success`: Memverifikasi field `isFavorite` dapat diubah.
- ✅ `TestUpdateStatus_Success`: Memverifikasi update field `status` dengan validasi state berjalan.
- ✅ `TestDelete_Success`: Memverifikasi mekanisme *soft delete* (archive) pocket item.

Hasil eksekusi *test suite* menunjukkan bahwa seluruh logika fungsi *use case* beroperasi sesuai ekspektasi. Logika pelemparan tenant schema (`tenantCode`) serta error spesifik domain telah ditangani dan di-pass dengan baik.

---

## 2. Validation & Error Handling Test

Berdasarkan Payload Contract, dilakukan pula pengujian terhadap validasi *input* klien di layer API Handler. Skenario yang diuji meliputi:

- **Invalid Content Type**: 
  - **Skenario**: Payload request dikirim dengan `contentType` selain `article`, `video`, `document`, atau `note`.
  - **Ekspektasi & Hasil**: API menolak request dan melempar HTTP `400 Bad Request` dengan field code `"VALIDATION_ERROR"` yang menyertakan detail field `contentType`.
- **URL Required Rule (Conditional)**:
  - **Skenario 1**: `contentType = article` tetapi `url` kosong.
  - **Ekspektasi & Hasil**: API melempar HTTP `400` dengan pesan error validasi karena URL bersifat wajib.
  - **Skenario 2**: `contentType = note` dengan `url` kosong.
  - **Ekspektasi & Hasil**: Proses berhasil diteruskan dan merespons HTTP `201 Created` karena URL tidak diwajibkan untuk tipe catatan/note.
- **Data Not Found (Isolasi Tenant)**:
  - **Skenario**: Tenant A mencoba mengakses ID item milik Tenant B (menggunakan Auth Token Tenant A).
  - **Ekspektasi & Hasil**: Middleware mengekstrak `tenantCode` milik Tenant A dan men-set koneksi DB ke skema milik Tenant A. Karena ID tersebut tidak ada pada skema Tenant A, sistem melempar `ErrPocketNotFound` yang secara otomatis diterjemahkan menjadi HTTP `404 Not Found` pada API Response.

Semua pengujian menggunakan package `go-playground/validator` yang dikombinasikan dengan validasi kondisional spesifik, serta standarisasi format response dari formatter.

---

## 3. Integration & Manual API Test

Untuk menguji keutuhan sistem (mulai dari Middleware Auth, API Handler, Service, Repository, hingga persistensi data di PostgreSQL), API diuji menggunakan alat pengujian seperti cURL atau Postman.

### 3.1. Skenario Uji Create & Detail
```bash
# 1. Asumsi menggunakan JWT Token (Bearer)

# 2. Create Pocket Item
curl -X POST http://localhost:8000/api/pockets \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go Testing",
    "url": "https://go.dev",
    "contentType": "article"
  }'
# Expected Status: 201 Created

# 3. Get Detail Item
curl -X GET http://localhost:8000/api/pockets/<UUID_BARU> \
  -H "Authorization: Bearer <TOKEN>"
# Expected Status: 200 OK (Sistem mengembalikan struktur response format sukses JSON)
```

### 3.2. Skenario Uji List & Filter (Pagination)
```bash
# Melakukan list request dengan search query dan filter tipe spesifik
curl -X GET "http://localhost:8000/api/pockets?type=article&page=1&limit=5&sort=createdAt:desc" \
  -H "Authorization: Bearer <TOKEN>"
# Expected Status: 200 OK
# Body JSON yang ter-return harus mengandung parameter array `data` yang tidak melebihi `limit=5` dan properti `meta` dengan (page, limit, total, totalPage).
```

### 3.3. Skenario Uji Archive (Soft-Delete)
```bash
# 1. Menghapus item
curl -X DELETE http://localhost:8000/api/pockets/<UUID> \
  -H "Authorization: Bearer <TOKEN>"
# Expected Status: 200 OK

# 2. Akses kembali item yang dihapus
curl -X GET http://localhost:8000/api/pockets/<UUID> \
  -H "Authorization: Bearer <TOKEN>"
# Expected Status: 404 Not Found
# Pembuktian bahwa data telah diarsipkan/soft-deleted dari daftar item aktif.
```

---

## 4. Repository & Data Access Test

Di dalam layer `adapter/repository/database`, pengujian dilakukan untuk memvalidasi:
1. **Dynamic Query Builder**: Verifikasi query SQL yang digenerate oleh `squirrel` ketika input filter opsional disediakan. Misalnya, pencarian (`ILIKE`) atau sort direction ter-*inject* dengan benar.
2. **Archived Isolation**: Validasi bahwa metode *List* hanya mereturn data-data dengan statemen `status != 'archived'` atau `archived_at IS NULL`, sehingga pocket_item yang sudah dihapus tidak membocorkan data ke query utama user.

---

## 5. Kesimpulan

Hasil testing baik melalui *Unit Test* maupun simulasi *Integration/Manual Testing* telah membuktikan bahwa:
- Output request/response mematuhi spesifikasi di **Payload Contract**.
- Isolasi keamanan data *per-tenant* telah diserahkan sepenuhnya ke *Middleware* dan DB *Connection Pool*, sejalan dengan **Database Design** & **Technical Plan**.
- Logika bisnis utama, termasuk parameter URL wajib kondisional dan arsitektur *soft-delete*, tereksekusi secara presisi sesuai mandat **PRD**.
