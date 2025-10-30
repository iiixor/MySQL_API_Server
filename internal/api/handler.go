package api

import (
	"net/http"
	"time"

	"mysql-tui-editor/server/internal/domain"
	"mysql-tui-editor/server/internal/executor"
	"mysql-tui-editor/server/internal/security"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	executor  *executor.MySQLExecutor
	validator *security.Validator
}

// NewHandler creates a new HTTP handler
func NewHandler(executor *executor.MySQLExecutor, validator *security.Validator) *Handler {
	return &Handler{
		executor:  executor,
		validator: validator,
	}
}

// ExecuteQuery handles POST /api/v1/execute
func (h *Handler) ExecuteQuery(c *gin.Context) {
	var req domain.ExecuteRequest

	// Bind JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse("Invalid request format: "+err.Error()))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(err.Error()))
		return
	}

	// Validate SQL security
	if err := h.validator.Validate(req.Query); err != nil {
		c.JSON(http.StatusForbidden, domain.NewErrorResponse("Security validation failed: "+err.Error()))
		return
	}

	// Execute query
	startTime := time.Now()
	response, err := h.executor.Execute(c.Request.Context(), req.Query)
	executionTime := time.Since(startTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse("Internal server error: "+err.Error()))
		return
	}

	// Log execution
	logQueryExecution(req.Query, response.Success, executionTime)

	// Return response
	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, response) // Still 200 OK, but success=false
	}
}

// HealthCheck handles GET /api/v1/health
func (h *Handler) HealthCheck(c *gin.Context) {
	// Check MySQL connection
	if err := h.executor.GetDB().Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"message": "MySQL connection failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Server is running",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// logQueryExecution logs query execution details
func logQueryExecution(query string, success bool, duration time.Duration) {
	// Simple logging - in production, use structured logger
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}

	// Truncate query for logging
	truncatedQuery := query
	if len(query) > 100 {
		truncatedQuery = query[:100] + "..."
	}

	// Replace newlines for cleaner logs
	truncatedQuery = truncateForLog(truncatedQuery)

	// Log
	println(time.Now().Format("2006-01-02 15:04:05"), status, duration.String(), truncatedQuery)
}

// truncateForLog truncates and cleans string for logging
func truncateForLog(s string) string {
	// Remove newlines and extra spaces
	cleaned := ""
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' {
			cleaned += " "
		} else {
			cleaned += string(r)
		}
	}
	return cleaned
}
