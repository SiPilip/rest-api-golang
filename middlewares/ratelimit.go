package middlewares

import (
	"REST-API/helpers"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Buat limiter per ip
type ipLimiter struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

var (
	limiters = make(map[string]*ipLimiter)
	mu sync.Mutex
)

// Bersihkan IP Lama setiap 3 menit
func init(){
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mu.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 1 * time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// getlimiter mengambil atau membuat lmiter untuk IP tertentu
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if l, exists := limiters[ip]; exists {
		l.lastSeen = time.Now()
		return l.limiter
	}

	// 10 Request per detik, dengan burst maksimal 20
	limiter := rate.NewLimiter(10, 20)
	limiters[ip] = &ipLimiter{
		limiter: limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			helpers.ErrorResponse(c, http.StatusTooManyRequests, "Too many request. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
rate.NewLimiter(10, 20) — artinya:
	10 = rate: 10 token per detik (10 request/detik diizinkan)
	20 = burst: boleh kirim 20 request sekaligus (untuk lonjakan singkat)
sync.Mutex — karena map diakses oleh banyak goroutine (banyak request bersamaan), harus dikunci supaya tidak race condition
Cleanup goroutine — IP yang tidak aktif selama 3 menit dihapus dari map agar tidak bocor memori
*/