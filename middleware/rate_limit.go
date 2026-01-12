// middleware/rate_limit.go
package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)


type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests
	window   time.Duration // time window
}

type visitor struct {
	count    int
	lastSeen time.Time
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	rate:     100,
	window:   1 * time.Minute,
}


func RateLimit(c *fiber.Ctx) error {
	ip := c.IP()

	limiter.mu.Lock()
	v, exists := limiter.visitors[ip]
	if !exists {
		limiter.visitors[ip] = &visitor{
			count:    1,
			lastSeen: time.Now(),
		}
		limiter.mu.Unlock()
		return c.Next()
	}


	if time.Since(v.lastSeen) > limiter.window {
		v.count = 1
		v.lastSeen = time.Now()
		limiter.mu.Unlock()
		return c.Next()
	}


	if v.count >= limiter.rate {
		limiter.mu.Unlock()
		return c.Status(429).JSON(fiber.Map{
			"error": "Too many requests. Please try again later.",
		})
	}

	v.count++
	v.lastSeen = time.Now()
	limiter.mu.Unlock()

	return c.Next()
}

// Cleanup old visitors periodically
func init() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limiter.mu.Lock()
			for ip, v := range limiter.visitors {
				if time.Since(v.lastSeen) > limiter.window {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}
