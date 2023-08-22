package middleware

import "github.com/gin-gonic/gin"

var (
	defaultHeaders = map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, PATH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Origin, Authorization, Content-Type",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Max-Age":           "86400",
	}
)

func CORS(headers map[string]string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Override default headers
		if len(headers) != 0 {
			for key, value := range headers {
				defaultHeaders[key] = value
			}
		}

		// Set headers
		for key, value := range defaultHeaders {
			ctx.Writer.Header().Set(key, value)
		}

		// Handle preflight requests
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
