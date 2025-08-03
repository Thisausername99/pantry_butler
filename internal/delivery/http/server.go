package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/thisausername99/pantry_butler/internal/delivery/graphql"
	"github.com/thisausername99/pantry_butler/internal/usecase"
)

// Server represents the HTTP server
type Server struct {
	router  *gin.Engine
	server  *http.Server
	logger  *zap.Logger
	useCase *usecase.Usecase
}

// NewServer creates a new HTTP server
func NewServer(logger *zap.Logger, useCase *usecase.Usecase) *Server {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	server := &Server{
		router:  router,
		logger:  logger,
		useCase: useCase,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

// setupMiddleware configures all middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware (handles panics)
	s.router.Use(RecoveryMiddleware(s.logger))

	// Request logging
	s.router.Use(RequestLogger(s.logger))

	// CORS middleware
	s.router.Use(CORSMiddleware())

	// Request context (adds request ID, user info, etc.)
	s.router.Use(RequestContext(s.logger))

	// Rate limiting
	s.router.Use(RateLimitMiddleware(s.logger))

	// Health check middleware
	s.router.Use(HealthCheckMiddleware())

	// GraphQL context middleware
	s.router.Use(GraphQLContextMiddleware())

	// Authentication middleware (optional - can be enabled per route)
	// authConfig := config.NewAuthConfig(s.useCase, s.logger)
	// s.router.Use(AuthMiddleware(authConfig))
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "pantry-butler",
			"version":   "1.0.0",
		})
	})

	// API info endpoint
	s.router.GET("/api/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Pantry Butler API",
			"version":     "1.0.0",
			"description": "GraphQL API for managing pantry items and recipes",
			"endpoints": gin.H{
				"graphql":    "/query",
				"playground": "/",
				"health":     "/health",
			},
		})
	})

	// GraphQL playground (GET requests)
	s.router.GET("/", playgroundHandler())

	// GraphQL endpoint (POST requests)
	s.router.POST("/query", graphqlHandler(s.useCase, s.logger))

	// GraphQL endpoint (GET requests for queries)
	s.router.GET("/query", graphqlHandler(s.useCase, s.logger))

	// 404 handler
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Endpoint not found",
			"code":  "NOT_FOUND",
			"path":  c.Request.URL.Path,
		})
	})
}

// playgroundHandler creates the GraphQL playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// graphqlHandler creates the GraphQL handler
func graphqlHandler(useCase *usecase.Usecase, logger *zap.Logger) gin.HandlerFunc {
	// Create GraphQL schema
	schema := graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &graphql.Resolver{UseCase: *useCase},
	})

	// Create GraphQL handler
	h := handler.NewDefaultServer(schema)

	// Add panic recovery
	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.Error("GraphQL panic recovered",
			zap.Any("panic", err),
			zap.String("requestID", getRequestIDFromContext(ctx)),
		)
		return fmt.Errorf("internal server error")
	})

	return func(c *gin.Context) {
		// Log GraphQL request
		logger.Info("GraphQL request",
			zap.String("requestID", getRequestIDFromGin(c)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("userID", c.GetHeader("X-User-ID")),
		)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	if port == "" {
		port = "8080"
	}

	s.server = &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info("Starting HTTP server",
		zap.String("port", port),
		zap.String("address", s.server.Addr),
	)

	return s.server.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		s.logger.Info("Stopping HTTP server")
		return s.server.Shutdown(ctx)
	}
	return nil
}

// GetRouter returns the Gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Helper function to get request ID from context
func getRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value("requestID").(string); ok {
		return id
	}
	return "unknown"
}

// Helper function to get request ID from Gin context
func getRequestIDFromGin(c *gin.Context) string {
	if ctx := c.Request.Context(); ctx != nil {
		if id, ok := ctx.Value("requestID").(string); ok {
			return id
		}
	}
	return "unknown"
}
