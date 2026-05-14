package helpers

import (
	"REST-API/db"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Simpan data ke cache
func SetCache(ctx context.Context, key string, data any, ttl time.Duration) error {
	if db.RedisClient == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return db.RedisClient.Set(ctx, key, jsonData, ttl).Err()
}

// Ambil data dari cache, return nil kalau tidak ada / expired
func GetCache(ctx context.Context, key string, dest any) error {
	if db.RedisClient == nil {
		return fmt.Errorf("redis not available")
	}

	val, err := db.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Hapus cache berdasarkan pattern
func DeleteCacheByPattern(ctx context.Context, pattern string) error {
	if db.RedisClient == nil {
		return nil
	}

	keys, err := db.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return db.RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

/*
SetCache — serialize data ke JSON, simpan di Redis dengan TTL (auto expired)
GetCache — ambil dari Redis, deserialize kembali ke struct
DeleteCacheByPattern — hapus semua cache yang cocok pattern (misal events:*)
*/