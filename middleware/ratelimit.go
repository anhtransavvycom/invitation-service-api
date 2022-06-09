package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
	"log"
	"time"
)

func RequestLimit(rate int) func(c *gin.Context) {
	limit := ratelimit.New(rate)
	prev := time.Now()

	log.SetPrefix("[GIN] ")
	log.SetOutput(gin.DefaultWriter)

	return func(c *gin.Context) {
		now := limit.Take()
		log.Print(now.Sub(prev))
		prev = now
		c.Next()
	}
}
