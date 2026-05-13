# 🎓 Rencana Belajar REST API Golang — Tier 1

## Pendekatan Belajar

> [!TIP]
> **Metode: Belajar Sambil Upgrade Proyek yang Sudah Ada**
> 
> Daripada buat proyek baru, kita akan **menerapkan setiap topik langsung ke REST API Event** yang sudah Anda buat. Ini lebih efektif karena:
> 1. Anda sudah paham konteks kodenya
> 2. Langsung lihat "before vs after" — jadi tahu kenapa topik itu penting
> 3. Di akhir Tier 1, proyek Anda akan jadi **production-ready**

## Urutan Belajar (Diurutkan dari yang Paling Berdampak)

Setiap sesi, kita akan:
- 📖 **Penjelasan singkat** — Apa dan kenapa topik ini penting
- 🛠️ **Praktek langsung** — Terapkan ke kode Anda
- ✅ **Verifikasi** — Test bareng untuk pastikan berjalan

---

### Sesi 1: 🔐 Environment Variables
**Durasi: ~15 menit**

**Kenapa ini duluan?** Karena saat ini DSN database dan JWT secret key Anda di-hardcode. Ini adalah security risk #1 yang harus diperbaiki sebelum belajar hal lain.

Yang akan dipelajari:
- Membuat file `.env`
- Menggunakan `godotenv` untuk load environment variables
- Memindahkan semua sensitive config (DB DSN, JWT secret, port) ke `.env`
- Membuat `.gitignore` agar `.env` tidak ikut di-push ke Git

---

### Sesi 2: 🚨 Error Handling Terpusat & Standard Response
**Durasi: ~25 menit**

**Kenapa?** Saat ini error response Anda tidak konsisten (kadang `"message"`, kadang `"error"`). API profesional harus punya format response yang seragam.

Yang akan dipelajari:
- Membuat custom error types
- Membuat standard response wrapper: `{"status": "success/error", "message": "...", "data": {...}}`
- Recovery middleware (menangkap panic agar server tidak crash)

---

### Sesi 3: 📝 Logging dengan `slog`
**Durasi: ~20 menit**

**Kenapa?** `fmt.Println` tidak cukup untuk production. Anda butuh log yang terstruktur dengan timestamp, level (INFO/WARN/ERROR), dan context.

Yang akan dipelajari:
- Setup `slog` (built-in Go, tidak perlu library external)
- Mengganti semua `fmt.Println` dengan structured logging
- Log setiap request masuk (request logger middleware)
- Log error dengan stack trace

---

### Sesi 4: ✅ Input Validation (Custom)
**Durasi: ~20 menit**

**Kenapa?** Saat ini Anda hanya pakai `binding:"required"`. Tidak ada validasi format email, panjang password, atau sanitasi input.

Yang akan dipelajari:
- Custom validator dengan `go-playground/validator`
- Validasi email format, password minimum length
- Sanitasi input untuk mencegah SQL injection & XSS
- Translate error message validator ke response yang user-friendly

---

### Sesi 5: 📄 Pagination & Filtering
**Durasi: ~25 menit**

**Kenapa?** `GetAllEvents()` Anda mengambil SEMUA data sekaligus. Kalau ada 1 juta event, server akan crash.

Yang akan dipelajari:
- Query params: `?page=1&limit=10&search=concert&sort=datetime`
- `LIMIT` dan `OFFSET` di SQL
- Response metadata: `{"data": [...], "meta": {"page": 1, "total": 100, "total_pages": 10}}`
- Default values untuk pagination

---

### Sesi 6: 🌐 CORS (Cross-Origin Resource Sharing)
**Durasi: ~10 menit**

**Kenapa?** Kalau nanti Anda buat frontend (React/Next.js), browser akan memblokir request ke API tanpa CORS.

Yang akan dipelajari:
- Apa itu CORS dan kenapa browser memblokir cross-origin request
- Setup `gin-contrib/cors` middleware
- Konfigurasi allowed origins, methods, headers

---

### Sesi 7: 📁 File Upload & Static Serving
**Durasi: ~30 menit**

**Kenapa?** Ini yang Anda sebut sendiri! Hampir semua aplikasi butuh upload file (foto profil, dokumen, thumbnail event).

Yang akan dipelajari:
- `multipart/form-data` parsing di Gin
- Validasi file: tipe (hanya gambar), ukuran maksimal
- Menyimpan file ke disk dengan nama unik (UUID)
- Serving static files (`server.Static()`)
- Menambahkan field `image_url` ke Event

---

### Sesi 8: 🛑 Graceful Shutdown
**Durasi: ~15 menit**

**Kenapa?** Saat ini kalau Anda Ctrl+C, server langsung mati. Koneksi database dan request yang sedang diproses bisa rusak.

Yang akan dipelajari:
- `os.Signal` untuk menangkap SIGINT/SIGTERM
- `server.Shutdown(ctx)` dengan timeout
- Menutup koneksi database dengan bersih
- Pattern `context.WithTimeout` untuk shutdown

---

## Ringkasan Timeline

```
Sesi 1: Environment Variables     ████░░░░░░ (~15 min)
Sesi 2: Error Handling            ██████░░░░ (~25 min)
Sesi 3: Logging (slog)            █████░░░░░ (~20 min)
Sesi 4: Input Validation          █████░░░░░ (~20 min)
Sesi 5: Pagination & Filtering    ██████░░░░ (~25 min)
Sesi 6: CORS                      ███░░░░░░░ (~10 min)
Sesi 7: File Upload               ███████░░░ (~30 min)
Sesi 8: Graceful Shutdown         ████░░░░░░ (~15 min)
────────────────────────────────────────────
Total estimasi:                    ~2.5 - 3 jam
```

> [!IMPORTANT]
> **Setiap sesi bisa dilakukan terpisah** — tidak harus selesai dalam satu waktu. Kita bisa belajar satu topik hari ini, lanjut besok, dst.

## Setelah Tier 1 Selesai ✨

Proyek REST API Event Anda akan berubah dari **"latihan belajar"** menjadi **"production-ready API"** dengan:
- ✅ Config aman (environment variables)
- ✅ Error handling konsisten
- ✅ Logging terstruktur
- ✅ Input tervalidasi
- ✅ Pagination untuk data besar
- ✅ CORS untuk frontend
- ✅ Upload file
- ✅ Shutdown yang aman

Kemudian baru kita masuk ke **Tier 2** atau langsung ke **Case Study** dengan fondasi yang kuat! 🚀
