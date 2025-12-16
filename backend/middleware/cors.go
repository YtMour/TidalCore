package middleware

import (
	"github.com/gin-gonic/gin"

	"tidalcore-backend/config"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigins := getAllowedOrigins()

		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		if !allowed && len(allowedOrigins) > 0 {
			c.Header("Access-Control-Allow-Origin", allowedOrigins[0])
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getAllowedOrigins() []string {
	cfg := config.Get()
	if cfg != nil && len(cfg.Server.AllowedOrigins) > 0 {
		return cfg.Server.AllowedOrigins
	}
	return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
}
