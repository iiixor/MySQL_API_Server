package domain

// ExecuteResponse represents the response from SQL query execution
type ExecuteResponse struct {
	// Success indicates if the query executed without errors
	Success bool `json:"success"`

	// Output contains the raw MySQL text output (tables, messages, etc.)
	Output string `json:"output"`

	// ExecutionTimeMs is the time taken to execute the query in milliseconds
	ExecutionTimeMs int64 `json:"execution_time_ms"`

	// Error contains the error message if execution failed
	Error string `json:"error"`
}

// NewSuccessResponse creates a successful response
func NewSuccessResponse(output string, executionTimeMs int64) *ExecuteResponse {
	return &ExecuteResponse{
		Success:         true,
		Output:          output,
		ExecutionTimeMs: executionTimeMs,
		Error:           "",
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(errorMsg string) *ExecuteResponse {
	return &ExecuteResponse{
		Success:         false,
		Output:          "",
		ExecutionTimeMs: 0,
		Error:           errorMsg,
	}
}
