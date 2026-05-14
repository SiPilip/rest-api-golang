# 📚 REST API Golang — Learning Cheatsheet

## Project Structure

```
REST-API/
├── main.go                  # Entry point, graceful shutdown, swagger
├── .env                     # Environment variables (JANGAN push ke git!)
├── .env.example             # Template .env (aman di-push)
├── .gitignore               # .env, uploads/
├── db/
│   └── db.go                # Database connection & pooling
├── models/
│   ├── event.go             # Event CRUD, transaction, soft delete
│   └── user.go              # User auth, password hashing
├── routes/
│   ├── routes.go            # Route registration, versioning (/api/v1)
│   ├── event.go             # Event handlers + swagger annotations
│   ├── users.go             # Auth handlers (signup, login)
│   └── register.go          # Event registration handlers
├── middlewares/
│   ├── auth.go              # JWT authentication
│   ├── cors.go              # CORS policy
│   ├── logger.go            # Request logging (slog)
│   ├── ratelimit.go         # Rate limiting per IP
│   └── timeout.go           # Request timeout
├── helpers/
│   ├── response.go          # Standard response + validation errors
│   └── upload.go            # File upload validation & save
├── utils/
│   ├── jwt.go               # JWT generate & verify
│   └── hash.go              # bcrypt hash & verify
├── migrations/              # Database migration files
│   ├── README.md            # Migration cheatsheet
│   └── 000001_*.sql ...     # Migration files (up & down)
├── uploads/                 # Uploaded files (gitignored)
├── docs/                    # Swagger generated docs
└── tests/
    └── README.md            # Testing cheatsheet
```

---

## Dependencies

### Go Get

```bash
# Tier 1
go get github.com/joho/godotenv               # Environment variables
go get github.com/gin-contrib/cors             # CORS middleware
go get github.com/go-playground/validator/v10   # Input validation (sudah include di Gin)

# Tier 2
go get golang.org/x/time/rate                  # Rate limiting
go get github.com/swaggo/swag@latest           # Swagger core
go get github.com/swaggo/gin-swagger           # Swagger Gin integration
go get github.com/swaggo/files                 # Swagger UI files

# Sudah ada dari sebelumnya
go get github.com/gin-gonic/gin                # Web framework
go get github.com/go-sql-driver/mysql          # MySQL driver
go get github.com/golang-jwt/jwt/v5            # JWT
go get golang.org/x/crypto                     # bcrypt
```

### Go Install (CLI tools)

```bash
# Swagger CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Database migration CLI
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

## Tier 1 — Essential (Production-Ready)

### 1. Environment Variables

```go
// .env
DB_DSN=root:@tcp(localhost:3306)/go_udemy?parseTime=true
JWT_SECRET=your-secret-key
SERVER_PORT=3010

// main.go
godotenv.Load()
os.Getenv("DB_DSN")
```

### 2. Standard Response

```go
// helpers/response.go
type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

helpers.SuccessResponse(ctx, 200, "OK", data)
helpers.ErrorResponse(ctx, 400, "Error message")
helpers.ValidationErrorResponse(ctx, err)
```

### 3. Logging (slog) — Built-in Go

```go
import "log/slog"

// Setup
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
slog.SetDefault(logger)

// Usage
slog.Info("Message", "key", value)
slog.Error("Failed", "error", err, "userId", id)
```

### 4. Input Validation

```go
// Struct tags
Email    string `binding:"required,email"`
Password string `json:"-" binding:"required,min=6,max=255"`  // json:"-" = hide from response
Name     string `binding:"required,min=3,max=80"`

// Handler
err := ctx.ShouldBind(&data)   // Auto-detect JSON / FormData
err := ctx.ShouldBindJSON(&data) // JSON only
```

### 5. Pagination

```go
// Query params: ?page=1&limit=10&search=keyword
page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

// SQL
"SELECT ... LIMIT ? OFFSET ?"  // offset = (page - 1) * limit
"SELECT COUNT(*) FROM ..."     // untuk total & total_pages
```

### 6. CORS

```go
// middlewares/cors.go
cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
})
// CORS hanya berlaku di browser, BUKAN Postman/curl
```

### 7. File Upload

```go
// Handler — pakai FormData
file, err := ctx.FormFile("image")
helpers.ValidateImage(file)                    // cek tipe & ukuran
filePath, _ := helpers.SaveUploadedFile(file)  // generate nama unik
ctx.SaveUploadedFile(file, filePath)           // simpan ke disk

// Serve static
server.Static("/uploads", "./uploads")

// Struct tag untuk support JSON + FormData
Name string `json:"name" form:"name" binding:"required"`
```

### 8. Graceful Shutdown

```go
// Jalankan server di goroutine
srv := &http.Server{Addr: ":3010", Handler: server}
go func() { srv.ListenAndServe() }()

// Tunggu sinyal OS (Ctrl+C)
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// Shutdown dengan timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)
db.DB.Close()
```

---

## Tier 2 — Professional

### 9. Database Transaction

```go
tx, err := db.DB.Begin()
defer tx.Rollback()            // auto-rollback kalau belum commit

_, err = tx.Exec("DELETE ...", id)
if err != nil { return err }

_, err = tx.Exec("DELETE ...", id)
if err != nil { return err }

return tx.Commit()             // semua berhasil → simpan
```

### 10. Database Migration

```bash
# Buat migration baru
migrate create -ext sql -dir migrations -seq nama_migration

# Jalankan semua
migrate -database "mysql://root:@tcp(localhost:3306)/dbname" -path migrations up

# Rollback 1
migrate -database "..." -path migrations down 1

# Fix dirty database
migrate -database "..." -path migrations force -1
```

### 11. Soft Delete

```go
// Struct
DeletedAt *time.Time `json:"deleted_at,omitempty"`  // pointer untuk NULL

// Delete → UPDATE bukan DELETE
"UPDATE events SET deleted_at = NOW() WHERE id = ?"

// Semua SELECT harus filter
"SELECT ... WHERE deleted_at IS NULL"

// Restore
"UPDATE events SET deleted_at = NULL WHERE id = ?"
```

### 12. Rate Limiting

```go
// golang.org/x/time/rate
limiter := rate.NewLimiter(10, 20) // 10 req/detik, burst 20
if !limiter.Allow() {
    // 429 Too Many Requests
}

// Simpan limiter per IP di map + sync.Mutex
// Cleanup IP lama dengan goroutine background
```

### 13. Request Timeout

```go
// Middleware
ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
defer cancel()
c.Request = c.Request.WithContext(ctx)

// Model — pakai Context version
db.DB.QueryContext(ctx, "SELECT ...")
db.DB.ExecContext(ctx, "INSERT ...")
tx.ExecContext(ctx, "UPDATE ...")

// Handler — skip response kalau timeout
if c.Request.Context().Err() != nil {
    return // biar timeout middleware yang response
}
```

### 14. Testing

```bash
go test ./... -v          # semua test
go test ./helpers/ -v     # test package tertentu
go test ./... -cover      # test coverage
```

```go
// File: xxx_test.go (di folder yang sama)
func TestNama_Skenario(t *testing.T) {
    t.Error("gagal tapi lanjut")
    t.Fatal("gagal dan berhenti")
}

// HTTP test
w := httptest.NewRecorder()
req, _ := http.NewRequest("GET", "/events", nil)
router.ServeHTTP(w, req)
// w.Code, w.Body.String()
```

### 15. API Versioning

```go
v1 := server.Group("/api/v1")
{
    v1.GET("/events", getEvents)
    v1.POST("/login", login)
    // ...
}
// Endpoint: /api/v1/events, /api/v1/login, etc.
```

### 16. Swagger

```bash
# Generate docs
swag init

# Setiap ubah annotation → swag init ulang
```

```go
// Annotation di atas handler
// @Summary      Get all events
// @Tags         Events
// @Param        page query int false "Page"
// @Success      200  {object} helpers.Response
// @Security     BearerAuth
// @Router       /events [get]
func getEvents(ctx *gin.Context) {

// Serve UI
server.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// Akses: http://localhost:3010/api/docs/index.html
```

---

## Middleware Order (penting!)

```go
server.Use(middlewares.RateLimiter())                      // 1. Tolak spam dulu
server.Use(middlewares.CORSMiddleware())                   // 2. Handle preflight
server.Use(middlewares.TimeoutMiddleware(5 * time.Second))  // 3. Set batas waktu
server.Use(middlewares.RequestLogger())                     // 4. Log request
// ... lalu routes + auth middleware
```

---

## Quick Commands

```bash
# Run server
go run .

# Install dependency
go get package@latest

# Generate swagger
swag init

# Run migrations
migrate -database "mysql://root:@tcp(localhost:3306)/go_udemy" -path migrations up

# Run tests
go test ./... -v

# Build production binary
go build -o server .
```
