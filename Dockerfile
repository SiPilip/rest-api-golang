# ========== STAGE 1: BUILD ===============
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files dulu (untuk cache Layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary (statuc, tanpa CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# =========== Stage 2: RUN ================
FROM alpine:latest

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/server .

# Copy migraiton files (untuk dijalankan saat deploy)
COPY --from=builder /app/migrations ./migrations

# Buat folder uploads
RUN mkdir -p uploads

# Expose port
EXPOSE 3010

# Jalankan binary
CMD ["./server"]