package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Sandbox represents an isolated MySQL database for query execution
type Sandbox struct {
	executor *MySQLExecutor
	dbName   string
}

// NewSandbox creates a new isolated sandbox database
func NewSandbox(executor *MySQLExecutor, dbPrefix string) (*Sandbox, error) {
	// Generate unique database name using UUID
	dbName := fmt.Sprintf("%s%s", dbPrefix, generateShortUUID())

	sandbox := &Sandbox{
		executor: executor,
		dbName:   dbName,
	}

	// Create the temporary database
	if err := sandbox.create(context.Background()); err != nil {
		return nil, err
	}

	return sandbox, nil
}

// create creates the temporary database
func (s *Sandbox) create(ctx context.Context) error {
	query := fmt.Sprintf("CREATE DATABASE `%s`", s.dbName)
	_, err := s.executor.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", s.dbName, err)
	}
	return nil
}

// Cleanup drops the temporary database
func (s *Sandbox) Cleanup(ctx context.Context) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", s.dbName)
	_, err := s.executor.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop database %s: %w", s.dbName, err)
	}
	return nil
}

// ExecuteQuery executes SQL query in the sandbox and returns formatted output
func (s *Sandbox) ExecuteQuery(ctx context.Context, query string) (string, error) {
	// First, switch to the sandbox database
	useQuery := fmt.Sprintf("USE `%s`", s.dbName)
	if _, err := s.executor.db.ExecContext(ctx, useQuery); err != nil {
		return "", fmt.Errorf("failed to switch to database %s: %w", s.dbName, err)
	}

	// Split query into individual statements
	statements := splitSQLStatements(query)
	if len(statements) == 0 {
		return "", fmt.Errorf("no valid SQL statements found")
	}

	var outputBuilder strings.Builder

	// Execute each statement
	for i, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		// Execute statement
		result, err := s.executeStatement(ctx, stmt)
		if err != nil {
			return "", fmt.Errorf("error in statement %d: %w", i+1, err)
		}

		// Append output
		if i > 0 {
			outputBuilder.WriteString("\n\n")
		}
		outputBuilder.WriteString(result)
	}

	return outputBuilder.String(), nil
}

// executeStatement executes a single SQL statement and formats the output
func (s *Sandbox) executeStatement(ctx context.Context, stmt string) (string, error) {
	// Determine if this is a SELECT query
	trimmedStmt := strings.TrimSpace(strings.ToUpper(stmt))
	isSelect := strings.HasPrefix(trimmedStmt, "SELECT") ||
		strings.HasPrefix(trimmedStmt, "SHOW") ||
		strings.HasPrefix(trimmedStmt, "DESCRIBE") ||
		strings.HasPrefix(trimmedStmt, "DESC") ||
		strings.HasPrefix(trimmedStmt, "EXPLAIN")

	if isSelect {
		return s.executeSelectStatement(ctx, stmt)
	}

	return s.executeNonSelectStatement(ctx, stmt)
}

// executeSelectStatement executes a SELECT-like statement and formats results as a table
func (s *Sandbox) executeSelectStatement(ctx context.Context, stmt string) (string, error) {
	rows, err := s.executor.db.QueryContext(ctx, stmt)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	return formatResultSet(rows)
}

// executeNonSelectStatement executes INSERT, UPDATE, DELETE, CREATE, etc.
func (s *Sandbox) executeNonSelectStatement(ctx context.Context, stmt string) (string, error) {
	result, err := s.executor.db.ExecContext(ctx, stmt)
	if err != nil {
		return "", err
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertID, _ := result.LastInsertId()

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Query OK, %d row(s) affected", rowsAffected))

	if lastInsertID > 0 {
		output.WriteString(fmt.Sprintf(" (last insert ID: %d)", lastInsertID))
	}

	return output.String(), nil
}

// splitSQLStatements splits SQL string into individual statements
func splitSQLStatements(query string) []string {
	// Simple split by semicolon
	// TODO: Handle strings with semicolons inside
	parts := strings.Split(query, ";")

	statements := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			statements = append(statements, trimmed)
		}
	}

	return statements
}

// generateShortUUID generates a short UUID for database names
func generateShortUUID() string {
	fullUUID := uuid.New().String()
	// Take first 8 characters (remove hyphens)
	return strings.ReplaceAll(fullUUID[:13], "-", "")
}
