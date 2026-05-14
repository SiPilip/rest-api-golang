# 🧪 Testing di Go — Cheatsheet

## Konsep Dasar

| Jenis | Apa yang ditest | Butuh DB? |
|-------|----------------|:---------:|
| **Unit Test** | Fungsi individu (helper, utils) | ❌ |
| **Integration Test** | Endpoint lengkap (handler → DB) | ✅ |

## Aturan Go Testing

1. File test harus berakhiran `_test.go`
2. File test **di folder yang sama** dengan kode yang ditest
3. Nama fungsi test harus diawali `Test` + huruf kapital
4. Import package `"testing"`

```
helpers/
├── upload.go          ← kode asli
├── upload_test.go     ← test untuk upload.go
utils/
├── hash.go
├── hash_test.go
```

## Contoh Unit Test

### `helpers/upload_test.go`

```go
package helpers

import (
    "mime/multipart"
    "net/textproto"
    "testing"
)

func createFakeFile(filename string, size int64) *multipart.FileHeader {
    return &multipart.FileHeader{
        Filename: filename,
        Size:     size,
        Header:   make(textproto.MIMEHeader),
    }
}

func TestValidateImage_ValidJPG(t *testing.T) {
    file := createFakeFile("photo.jpg", 1024*1024)
    err := ValidateImage(file)
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}

func TestValidateImage_TooLarge(t *testing.T) {
    file := createFakeFile("huge.jpg", 5*1024*1024) // 5MB
    err := ValidateImage(file)
    if err == nil {
        t.Error("Expected error for file > 2MB, got nil")
    }
}

func TestValidateImage_InvalidExtension(t *testing.T) {
    file := createFakeFile("document.pdf", 1024)
    err := ValidateImage(file)
    if err == nil {
        t.Error("Expected error for PDF, got nil")
    }
}
```

### `utils/hash_test.go`

```go
package utils

import "testing"

func TestHashPassword_Success(t *testing.T) {
    hash, err := HashPassword("mypassword123")
    if err != nil {
        t.Fatalf("Expected no error, got: %v", err)
    }
    if hash == "mypassword123" {
        t.Error("Hash should not equal original password")
    }
}

func TestVerifyPassword_Correct(t *testing.T) {
    hash, _ := HashPassword("mypassword123")
    if !VerifyPassword("mypassword123", hash) {
        t.Error("Expected password to match hash")
    }
}

func TestVerifyPassword_Wrong(t *testing.T) {
    hash, _ := HashPassword("mypassword123")
    if VerifyPassword("wrongpassword", hash) {
        t.Error("Expected wrong password to NOT match")
    }
}
```

## Contoh Integration Test (Handler)

### `routes/users_test.go`

```go
package routes

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    return gin.Default()
}

func TestSignup_MissingFields(t *testing.T) {
    router := setupTestRouter()
    router.POST("/signup", signup)

    body := bytes.NewBuffer([]byte(`{}`))
    req, _ := http.NewRequest("POST", "/signup", body)
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("Expected 400, got %d", w.Code)
    }

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    if response["status"] != "error" {
        t.Errorf("Expected 'error', got '%s'", response["status"])
    }
}
```

## Perintah

```bash
# Jalankan semua test
go test ./... -v

# Test package tertentu
go test ./helpers/ -v
go test ./utils/ -v

# Cek test coverage
go test ./... -cover

# Test dengan detail coverage per fungsi
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out

# Test coverage visual di browser
go tool cover -html=coverage.out
```

## Testing Functions Cheatsheet

| Fungsi | Kegunaan |
|--------|----------|
| `t.Error(msg)` | Test gagal, lanjut test berikutnya |
| `t.Errorf(format, ...)` | Sama tapi dengan formatting |
| `t.Fatal(msg)` | Test gagal, **berhenti** |
| `t.Fatalf(format, ...)` | Sama tapi dengan formatting |
| `t.Log(msg)` | Print info (hanya muncul kalau `-v`) |
| `t.Skip(msg)` | Skip test ini |

## httptest Cheatsheet

```go
// Buat fake response recorder
w := httptest.NewRecorder()

// Buat request
req, _ := http.NewRequest("GET", "/events", nil)
req.Header.Set("Authorization", "token-here")

// Jalankan request
router.ServeHTTP(w, req)

// Inspect hasil
w.Code           // status code (200, 400, etc)
w.Body.String()  // response body sebagai string
w.Body.Bytes()   // response body sebagai []byte
```

## Tips

- Nama test: `TestNamaFungsi_Skenario` (contoh: `TestValidateImage_TooLarge`)
- Satu test = satu skenario, jangan campur
- Test **happy path** (sukses) DAN **error cases** (gagal)
- Untuk test yang butuh database, gunakan **test database** terpisah agar data production aman
