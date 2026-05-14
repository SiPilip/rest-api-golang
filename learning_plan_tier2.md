# 🎓 Rencana Belajar REST API Golang — Tier 2 (Professional)

## Recap Tier 1 ✅

| # | Topik | Status |
|---|-------|--------|
| 1 | Environment Variables | ✅ |
| 2 | Error Handling Terpusat | ✅ |
| 3 | Logging (slog) | ✅ |
| 4 | Input Validation | ✅ |
| 5 | Pagination & Filtering | ✅ |
| 6 | CORS | ✅ |
| 7 | File Upload | ✅ |
| 8 | Graceful Shutdown | ✅ |

---

## Apa yang Berubah di Tier 2?

> [!NOTE]
> **Tier 1** fokus pada **"API yang bekerja dengan benar dan aman"**
> **Tier 2** fokus pada **"API yang profesional, teruji, dan maintainable"**
> 
> Setelah Tier 2, kode Anda layak untuk **portofolio** dan **review interview**.

## Urutan Belajar Tier 2

### Sesi 9: 🔄 Database Transaction
**Durasi: ~20 menit**

**Apa itu?** Transaction memastikan beberapa operasi database berjalan **semua atau tidak sama sekali** (atomik).

**Contoh masalah tanpa transaction:**
Saat delete event → harus hapus registrasi dulu, baru hapus event. Bagaimana kalau registrasi terhapus tapi event gagal dihapus? Data jadi **inkonsisten**.

Yang akan dipelajari:
- `db.Begin()`, `tx.Commit()`, `tx.Rollback()`
- Menerapkan transaction di delete event (hapus registrasi + event atomik)
- Pattern `defer tx.Rollback()` yang aman

---

### Sesi 10: 📦 Database Migration Tool
**Durasi: ~25 menit**

**Apa itu?** Saat ini schema database Anda di-hardcode di `db.go` pakai `CREATE TABLE IF NOT EXISTS`. Masalahnya: bagaimana kalau perlu **tambah kolom**, **ubah tipe data**, atau **rollback** perubahan schema?

Yang akan dipelajari:
- Install & setup `golang-migrate`
- Membuat migration files (up & down)
- Menjalankan migrasi dari CLI
- Memindahkan schema dari `db.go` ke migration files

---

### Sesi 11: 🗑️ Soft Delete
**Durasi: ~20 menit**

**Apa itu?** Saat ini `event.Delete()` menghapus data **permanen** dari database. Di dunia nyata, data jarang benar-benar dihapus — cukup ditandai `deleted_at`.

Yang akan dipelajari:
- Menambahkan kolom `deleted_at` (nullable TIMESTAMP)
- Mengubah `Delete()` menjadi soft delete (SET deleted_at = NOW())
- Mengubah semua query SELECT agar filter `WHERE deleted_at IS NULL`
- Membuat endpoint restore (opsional)

---

### Sesi 12: 🚦 Rate Limiting
**Durasi: ~15 menit**

**Apa itu?** Membatasi jumlah request per IP/user dalam periode waktu tertentu. Mencegah brute force login, DDoS, dan abuse API.

Yang akan dipelajari:
- Konsep token bucket / sliding window
- Implementasi rate limiter middleware
- Konfigurasi: max 100 requests per menit per IP
- Response `429 Too Many Requests`

---

### Sesi 13: ⏱️ Request Timeout / Context
**Durasi: ~15 menit**

**Apa itu?** Kalau database lambat atau query hang, request bisa menggantung selamanya. Timeout memastikan setiap request punya **batas waktu**.

Yang akan dipelajari:
- `context.WithTimeout` di database queries
- Timeout middleware untuk semua request
- Perbedaan `context.Context` vs `gin.Context`

---

### Sesi 14: 🧪 Unit Testing & Integration Testing
**Durasi: ~30 menit**

**Apa itu?** Menulis kode otomatis untuk memverifikasi bahwa API Anda bekerja dengan benar. Ini **wajib** di dunia profesional.

Yang akan dipelajari:
- Package `testing` dan `net/http/httptest`
- Menulis test untuk handler (integration test)
- Test helper functions (unit test)
- Menjalankan test dengan `go test`
- Test coverage

---

### Sesi 15: 🔢 API Versioning
**Durasi: ~15 menit**

**Apa itu?** Ketika API sudah dipakai client, Anda tidak bisa sembarangan ubah format response atau endpoint. Versioning memungkinkan Anda buat **versi baru** tanpa merusak yang lama.

Yang akan dipelajari:
- URL-based versioning: `/api/v1/events`, `/api/v2/events`
- Restructure routes dengan route groups
- Kenapa versioning penting untuk backward compatibility

---

### Sesi 16: 📖 Swagger / OpenAPI Documentation
**Durasi: ~25 menit**

**Apa itu?** Auto-generate dokumentasi API yang **interaktif** — developer lain bisa baca endpoint, parameter, response, dan langsung test dari browser.

Yang akan dipelajari:
- Install `swaggo/swag`
- Menulis annotation di handler functions
- Generate & serve Swagger UI
- Akses dokumentasi di `http://localhost:3010/swagger/index.html`

---

## Ringkasan Timeline

```
Sesi 9:  Database Transaction     █████░░░░░ (~20 min)
Sesi 10: Migration Tool           ██████░░░░ (~25 min)
Sesi 11: Soft Delete              █████░░░░░ (~20 min)
Sesi 12: Rate Limiting            ████░░░░░░ (~15 min)
Sesi 13: Request Timeout          ████░░░░░░ (~15 min)
Sesi 14: Unit Testing             ███████░░░ (~30 min)
Sesi 15: API Versioning           ████░░░░░░ (~15 min)
Sesi 16: Swagger Docs             ██████░░░░ (~25 min)
────────────────────────────────────────────────
Total estimasi:                    ~2.5 - 3 jam
```

> [!TIP]
> Sama seperti Tier 1 — setiap sesi bisa dilakukan terpisah. Semua diterapkan langsung ke proyek REST API Event yang sudah ada.

## Setelah Tier 2 Selesai ✨

API Anda akan memiliki kualitas **production-grade**:
- ✅ Data konsisten (transaction)
- ✅ Schema terkelola (migration)
- ✅ Data aman (soft delete)
- ✅ Terlindungi dari abuse (rate limiting)
- ✅ Tidak hang (timeout)
- ✅ Teruji otomatis (testing)
- ✅ Backward compatible (versioning)
- ✅ Terdokumentasi (Swagger)

Ini sudah layak jadi **portofolio** yang impress recruiter! 💼
