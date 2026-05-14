# 🎓 Rencana Belajar REST API Golang — Tier 3 (Advanced / Production-Grade)

## Recap

| Tier | Fokus | Status |
|------|-------|--------|
| **Tier 1** | API yang bekerja dengan benar & aman | ✅ 8/8 selesai |
| **Tier 2** | API yang profesional, teruji & maintainable | ✅ 8/8 selesai |
| **Tier 3** | API yang scalable, enterprise-ready & deployable | ⏳ |

> [!NOTE]
> **Tier 3** membawa API Anda ke level **enterprise**. Topik di sini adalah pembeda antara developer junior dan senior. Beberapa topik (Docker, CI/CD) Anda sudah pernah sentuh di proyek lain — kali ini kita fokus integrasinya dengan REST API ini.

---

## Urutan Belajar Tier 3

### Sesi 17: 🗄️ Caching dengan Redis
**Durasi: ~35 menit**

**Kenapa?** Setiap kali `GET /events` dipanggil, server query database. Kalau 1000 user hit endpoint yang sama dalam 1 detik — 1000 query identik ke database. **Redis** menyimpan hasil query di memori → response 10-100x lebih cepat.

Yang akan dipelajari:
- Install & koneksi Redis (`go-redis`)
- Cache response `GET /events` dan `GET /events/:id`
- Cache invalidation — hapus cache saat data berubah (create/update/delete)
- Set TTL (Time To Live) — cache expired otomatis setelah X detik
- Pattern: **Cache-Aside** (cek cache → miss → query DB → simpan ke cache)

---

### Sesi 18: 👥 Role-Based Access Control (RBAC)
**Durasi: ~30 menit**

**Kenapa?** Saat ini Anda hanya punya 1 jenis user. Di dunia nyata butuh: **Admin** (bisa semua), **User** (hanya event sendiri). Saat ini ownership check di-hardcode di setiap handler — tidak scalable.

Yang akan dipelajari:
- Menambah field `role` di tabel users (admin, user)
- Membuat authorization middleware (`RequireRole("admin")`)
- Admin bisa delete event siapapun, user hanya event sendiri
- Menyimpan role di JWT claims

---

### Sesi 19: 🔑 Refresh Token
**Durasi: ~30 menit**

**Kenapa?** JWT Anda sekarang expire dalam 2 jam. Setelah itu user harus login lagi. Di production, pattern yang lebih baik: **access token** (singkat, 15 menit) + **refresh token** (panjang, 7 hari).

Yang akan dipelajari:
- Generate access token (15 min) + refresh token (7 hari)
- Simpan refresh token di database
- Endpoint `POST /refresh` untuk minta access token baru
- Revoke refresh token saat logout
- Kenapa pattern ini lebih aman

---

### Sesi 20: 🏗️ Repository Pattern
**Durasi: ~30 menit**

**Kenapa?** Saat ini model langsung query database (`db.DB.Query`). Kalau nanti mau ganti dari MySQL ke PostgreSQL, harus ubah **semua** model. Repository pattern membuat **abstraksi layer** antara business logic dan database.

Yang akan dipelajari:
- Membuat interface `EventRepository` (contract)
- Membuat implementasi `MySQLEventRepository`
- Dependency injection — handler terima interface, bukan concrete type
- Kenapa ini mempermudah testing (bisa mock database)

---

### Sesi 21: 🏥 Health Check & Request ID
**Durasi: ~15 menit**

**Kenapa?** 
- **Health check** (`GET /health`) — load balancer & monitoring tools perlu tahu apakah server hidup
- **Request ID** — setiap request diberi UUID unik, sehingga kalau ada error di production, bisa trace log dari awal sampai akhir

Yang akan dipelajari:
- Endpoint `GET /health` yang cek status DB, Redis, dll
- Middleware request ID generator
- Menyisipkan request ID di setiap log entry

---

### Sesi 22: 📧 Email Sending
**Durasi: ~25 menit**

**Kenapa?** Fitur dasar: verifikasi email saat signup, notifikasi saat event baru, reset password.

Yang akan dipelajari:
- Setup SMTP (bisa pakai Mailtrap untuk testing)
- Membuat email service (`helpers/email.go`)
- Kirim email verifikasi saat signup
- HTML email template

---

### Sesi 23: ⚡ Background Jobs / Queue
**Durasi: ~30 menit**

**Kenapa?** Kirim email, proses gambar, generate report — semua ini **lambat** dan tidak boleh memblokir response ke user. Harus dijalankan di background.

Yang akan dipelajari:
- Goroutine + channel untuk simple background jobs
- Pattern worker pool
- Contoh: kirim email di background setelah register event
- Kenapa di production pakai message queue (Redis Queue, RabbitMQ)

---

### Sesi 24: 🐳 Docker & Deployment
**Durasi: ~30 menit**

**Kenapa?** Supaya API bisa di-deploy ke server manapun dengan satu perintah. Anda sudah pernah setup Docker di proyek lain — kali ini kita containerize REST API Go + MySQL.

Yang akan dipelajari:
- Membuat `Dockerfile` (multi-stage build untuk binary kecil)
- `docker-compose.yml` (Go API + MySQL + Redis)
- Environment variables di Docker
- Menjalankan migration di container

---

## Ringkasan Timeline

```
Sesi 17: Redis Caching          ████████░░ (~35 min)
Sesi 18: RBAC                   ███████░░░ (~30 min)
Sesi 19: Refresh Token          ███████░░░ (~30 min)
Sesi 20: Repository Pattern     ███████░░░ (~30 min)
Sesi 21: Health Check & Req ID  ████░░░░░░ (~15 min)
Sesi 22: Email Sending          ██████░░░░ (~25 min)
Sesi 23: Background Jobs        ███████░░░ (~30 min)
Sesi 24: Docker & Deployment    ███████░░░ (~30 min)
────────────────────────────────────────────────
Total estimasi:                  ~3.5 - 4.5 jam
```

> [!IMPORTANT]
> Sesi 17 (Redis) dan Sesi 18 (RBAC) paling berdampak untuk skill Anda. Sesi 20 (Repository Pattern) paling penting untuk interview karena menunjukkan pemahaman arsitektur.

## Setelah Tier 3 Selesai 🏆

Anda akan menguasai **seluruh stack** REST API production:

```
┌─────────────────────────────────────────┐
│           CLIENT (Browser/App)          │
├─────────────────────────────────────────┤
│  Rate Limit → CORS → Timeout → Logger  │  ← Middleware chain
│  Request ID → Auth → RBAC              │
├─────────────────────────────────────────┤
│        /api/v1/events (Handlers)        │  ← API Versioning
│        /api/v1/login                    │
│        /health                          │
├─────────────────────────────────────────┤
│      Repository (Interface Layer)       │  ← Abstraksi
├─────────────────────────────────────────┤
│  MySQL │ Redis Cache │ Email │ Queue    │  ← Infrastructure
├─────────────────────────────────────────┤
│          Docker + Deployment            │  ← Containerized
└─────────────────────────────────────────┘
```

Kemudian Anda 100% siap masuk **Case Study (Bank Soal API)** dengan fondasi yang sangat kuat! 🚀
