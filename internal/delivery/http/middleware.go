package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestContext adds request information to context
func RequestContext(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Extract request information
		userID := c.GetHeader("X-User-ID")
		authToken := c.GetHeader("Authorization")
		userAgent := c.GetHeader("User-Agent")

		// Add request information to context
		ctx := context.WithValue(c.Request.Context(), "requestID", requestID)
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "authToken", authToken)
		ctx = context.WithValue(ctx, "userAgent", userAgent)
		ctx = context.WithValue(ctx, "remoteAddr", c.ClientIP())

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		// Log request start
		logger.Info("Request started",
			zap.String("requestID", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", userID),
			zap.String("userAgent", userAgent),
			zap.String("remoteAddr", c.ClientIP()),
		)

		c.Next()
	}
}

// RequestLogger logs request completion
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract request ID from context
		requestID := ""
		if ctx := param.Request.Context(); ctx != nil {
			if id, ok := ctx.Value("requestID").(string); ok {
				requestID = id
			}
		}

		// Log request completion
		logger.Info("Request completed",
			zap.String("requestID", requestID),
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("clientIP", param.ClientIP),
			zap.String("error", param.ErrorMessage),
		)

		return ""
	})
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-ID")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthMiddleware validates authentication
func AuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for GraphQL playground
		if c.Request.URL.Path == "/" {
			c.Next()
			return
		}

		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			logger.Warn("Missing authorization token",
				zap.String("requestID", getRequestID(c)),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		// TODO: Add actual token validation logic here
		// For now, just check if token exists
		if authToken == "" {
			logger.Warn("Invalid authorization token",
				zap.String("requestID", getRequestID(c)),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Add user info to context
		ctx := context.WithValue(c.Request.Context(), "authenticated", true)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware(logger *zap.Logger) gin.HandlerFunc {
	// Simple in-memory rate limiter (use Redis for production)
	clients := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old requests (older than 1 minute)
		if times, exists := clients[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					validTimes = append(validTimes, t)
				}
			}
			clients[clientIP] = validTimes
		}

		// Check rate limit (max 100 requests per minute)
		if times, exists := clients[clientIP]; exists && len(times) >= 100 {
			logger.Warn("Rate limit exceeded",
				zap.String("requestID", getRequestID(c)),
				zap.String("clientIP", clientIP),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		// Add current request
		clients[clientIP] = append(clients[clientIP], now)

		c.Next()
	}
}

// RecoveryMiddleware handles panics
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := getRequestID(c)

		logger.Error("Panic recovered",
			zap.String("requestID", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("panic", recovered),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal server error",
			"code":      "INTERNAL_ERROR",
			"requestID": requestID,
		})
	})
}

// HealthCheckMiddleware adds health check endpoint
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"service":   "pantry-butler",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GraphQLContextMiddleware adds GraphQL-specific context
func GraphQLContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add GraphQL-specific context for resolvers
		ctx := c.Request.Context()

		// Add request information that will be available in GraphQL resolvers
		ctx = context.WithValue(ctx, "ginContext", c)
		ctx = context.WithValue(ctx, "requestMethod", c.Request.Method)
		ctx = context.WithValue(ctx, "requestPath", c.Request.URL.Path)
		ctx = context.WithValue(ctx, "requestQuery", c.Request.URL.RawQuery)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Helper function to get request ID from context
func getRequestID(c *gin.Context) string {
	if ctx := c.Request.Context(); ctx != nil {
		if id, ok := ctx.Value("requestID").(string); ok {
			return id
		}
	}
	return "unknown"
}

// ExtractRequestInfo extracts comprehensive request information
func ExtractRequestInfo(c *gin.Context) map[string]interface{} {
	info := make(map[string]interface{})

	// Basic request info
	info["method"] = c.Request.Method
	info["path"] = c.Request.URL.Path
	info["query"] = c.Request.URL.RawQuery
	info["clientIP"] = c.ClientIP()
	info["userAgent"] = c.GetHeader("User-Agent")

	// Headers
	info["userID"] = c.GetHeader("X-User-ID")
	info["authToken"] = c.GetHeader("Authorization")
	info["contentType"] = c.GetHeader("Content-Type")

	// Context values
	if ctx := c.Request.Context(); ctx != nil {
		if requestID, ok := ctx.Value("requestID").(string); ok {
			info["requestID"] = requestID
		}
		if userID, ok := ctx.Value("userID").(string); ok {
			info["contextUserID"] = userID
		}
		if authToken, ok := ctx.Value("authToken").(string); ok {
			info["contextAuthToken"] = authToken
		}
	}

	return info
}
