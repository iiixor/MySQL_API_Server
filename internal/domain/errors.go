package domain

import "errors"

// Common domain errors
var (
	// Request validation errors
	ErrEmptyQuery   = errors.New("query cannot be empty")
	ErrQueryTooLong = errors.New("query exceeds maximum allowed length")

	// Security errors
	ErrDangerousCommand = errors.New("query contains dangerous commands that are not allowed")
	ErrDropDatabase     = errors.New("DROP DATABASE command is not allowed")

	// Execution errors
	ErrExecutionTimeout = errors.New("query execution timeout exceeded")
	ErrDatabaseCreation = errors.New("failed to create temporary database")
	ErrDatabaseCleanup  = errors.New("failed to cleanup temporary database")

	// Connection errors
	ErrDatabaseConnection = errors.New("failed to connect to MySQL server")
)

// Constants
const (
	// MaxQueryLength is the maximum allowed length for a SQL query
	MaxQueryLength = 1024 * 1024 // 1MB
)
