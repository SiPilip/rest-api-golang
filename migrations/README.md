# 📦 Database Migration Cheatsheet

## Setup

```bash
# Install golang-migrate CLI
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Connection String

```
mysql://root:@tcp(localhost:3306)/go_udemy
```

> Sesuaikan `user:password`, `host:port`, dan `dbname` dengan environment Anda.

## Perintah Utama

```bash
# Jalankan SEMUA migrasi yang belum dijalankan
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations up

# Jalankan N migrasi berikutnya saja
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations up 1

# Rollback 1 versi
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations down 1

# Rollback SEMUA migrasi (hati-hati!)
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations down

# Cek versi migrasi sekarang
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations version
```

## Buat Migration Baru

```bash
# Format: migrate create -ext sql -dir migrations -seq <nama_migration>
migrate create -ext sql -dir migrations -seq create_users_table
migrate create -ext sql -dir migrations -seq add_phone_to_users
```

Ini akan generate 2 file:

- `000005_add_phone_to_users.up.sql` → perubahan yang diterapkan
- `000005_add_phone_to_users.down.sql` → cara membatalkan perubahan

## Troubleshooting: Dirty Database

Kalau migrasi gagal di tengah jalan, database jadi "dirty". Solusi:

```bash
# 1. Perbaiki file .sql yang error
# 2. Force reset ke state sebelum migrasi yang gagal
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations force -1

# 3. Jalankan ulang
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations up
```

> **`force -1`** = reset ke "belum ada migrasi"
> **`force N`** = set ke versi N (tanpa menjalankan migrasi)

## Tabel `schema_migrations`

| Kolom     | Arti                                         |
| --------- | -------------------------------------------- |
| `version` | Nomor migrasi terakhir yang dijalankan       |
| `dirty`   | `0` = sukses, `1` = gagal (perlu fix manual) |

> Jangan edit tabel ini manual — biar CLI yang kelola.

## Tips Menulis Migration

```sql
-- UP: selalu gunakan IF NOT EXISTS / IF EXISTS untuk safety
CREATE TABLE IF NOT EXISTS users (...);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- DOWN: kebalikan dari UP
DROP TABLE IF EXISTS users;
ALTER TABLE users DROP COLUMN phone;
```
