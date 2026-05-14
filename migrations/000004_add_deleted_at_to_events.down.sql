-- down: hapus kolom (rollback)
ALTER TABLE events DROP COLUMN deleted_at;