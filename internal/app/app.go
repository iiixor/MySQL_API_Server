package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mysql-tui-editor/server/internal/api"
	"mysql-tui-editor/server/internal/config"
	"mysql-tui-editor/server/internal/executor"
	"mysql-tui-editor/server/internal/security"

	"github.com/gin-gonic/gin"
)

// App represents the application
type App struct {
	config    *config.Config
	executor  *executor.MySQLExecutor
	validator *security.Validator
	handler   *api.Handler
	server    *http.Server
}

// New creates a new application instance
func New(configPath string) (*App, error) {
	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create MySQL executor
	exec, err := executor.NewMySQLExecutor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL executor: %w", err)
	}

	// Create validator
	validator := security.NewValidator()

	// Create handler
	handler := api.NewHandler(exec, validator)

	app := &App{
		config:    cfg,
		executor:  exec,
		validator: validator,
		handler:   handler,
	}

	return app, nil
}

// Run starts the application
func (a *App) Run() error {
	// Setup Gin
	if a.config.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(api.RecoveryMiddleware())
	router.Use(api.LoggingMiddleware())
	router.Use(api.CORSMiddleware())

	// Rate limiting
	rateLimiter := api.NewRateLimiter(
		a.config.Security.RateLimitPerSecond,
		a.config.Security.RateLimitBurst,
	)
	router.Use(rateLimiter.RateLimitMiddleware())

	// Routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/execute", a.handler.ExecuteQuery)
		v1.GET("/health", a.handler.HealthCheck)
	}

	// Root health check
	router.GET("/health", a.handler.HealthCheck)

	// Create HTTP server
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.Server.Port),
		Handler:      router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("🚀 MySQL TUI Server starting on port %d...\n", a.config.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	a.waitForShutdown()

	return nil
}

// waitForShutdown waits for interrupt signal and performs graceful shutdown
func (a *App) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		fmt.Printf("❌ Server forced to shutdown: %v\n", err)
	}

	// Close MySQL connection
	if err := a.executor.Close(); err != nil {
		fmt.Printf("❌ Error closing MySQL connection: %v\n", err)
	}

	fmt.Println("✅ Server exited cleanly")
}

// Close closes all resources
func (a *App) Close() error {
	if a.executor != nil {
		return a.executor.Close()
	}
	return nil
}
