package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Tipe file yang diizinkan
var allowedExtensions = map[string]bool{
	".jpg": true,
	".jpeg": true,
	".png": true,
	".webp": true,
}

const maxFileSize = 2 * 1024 * 1024 //2MB

func ValidateImage(file *multipart.FileHeader) error {
	// Cek ukuran file
	if file.Size > maxFileSize {
		return fmt.Errorf("file size exceeds 2mb limit")

	}

	// Cek ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type %s is not allowed. Use: jpg, jpeg, png, webp", ext)
	}
	return nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) (string,error) {
	// Buat folder uploads/ kalau belum ada
	uploadDir := "uploads"
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Generate nama unik: timestamp + ekstensi asli
	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	return filepath.ToSlash(filePath), nil
	
}