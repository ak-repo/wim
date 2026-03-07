package middleware

import (
	"time"

	"github.com/ak-repo/wim/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Info("request",
			"method", c.Request.Method,
			"path", path,
			"status", statusCode,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)
	}
}

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Error("handler error",
					"error", e.Error(),
					"type", e.Type,
				)
			}
		}
	}
}

func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
