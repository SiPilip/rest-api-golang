-- up: tambah kolom
ALTER TABLE events ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL;
