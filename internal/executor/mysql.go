package executor

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"mysql-tui-editor/server/internal/config"
	"mysql-tui-editor/server/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLExecutor handles MySQL query execution
type MySQLExecutor struct {
	db           *sql.DB
	queryTimeout time.Duration
	dbPrefix     string
}

// NewMySQLExecutor creates a new MySQL executor
func NewMySQLExecutor(cfg *config.Config) (*MySQLExecutor, error) {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
	)

	// Open connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return &MySQLExecutor{
		db:           db,
		queryTimeout: cfg.Executor.QueryTimeout,
		dbPrefix:     cfg.Executor.DBPrefix,
	}, nil
}

// Close closes the MySQL connection
func (e *MySQLExecutor) Close() error {
	return e.db.Close()
}

// GetDB returns the underlying database connection
func (e *MySQLExecutor) GetDB() *sql.DB {
	return e.db
}

// Execute executes SQL query in a sandboxed temporary database
func (e *MySQLExecutor) Execute(ctx context.Context, query string) (*domain.ExecuteResponse, error) {
	startTime := time.Now()

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, e.queryTimeout)
	defer cancel()

	// Create sandbox
	sandbox, err := NewSandbox(e, e.dbPrefix)
	if err != nil {
		return domain.NewErrorResponse(fmt.Sprintf("Failed to create sandbox: %v", err)), nil
	}

	// Ensure cleanup
	defer func() {
		if cleanupErr := sandbox.Cleanup(context.Background()); cleanupErr != nil {
			// Log cleanup error but don't fail the response
			fmt.Printf("WARNING: Failed to cleanup sandbox %s: %v\n", sandbox.dbName, cleanupErr)
		}
	}()

	// Execute query in sandbox
	output, err := sandbox.ExecuteQuery(execCtx, query)
	executionTime := time.Since(startTime).Milliseconds()

	if err != nil {
		// Check if it was a timeout
		if execCtx.Err() == context.DeadlineExceeded {
			return domain.NewErrorResponse(fmt.Sprintf("Query execution timeout exceeded (%v)", e.queryTimeout)), nil
		}
		return domain.NewErrorResponse(fmt.Sprintf("Query execution failed: %v", err)), nil
	}

	return domain.NewSuccessResponse(output, executionTime), nil
}
