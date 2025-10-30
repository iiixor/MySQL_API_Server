package domain

// ExecuteRequest represents a SQL query execution request from the client
type ExecuteRequest struct {
	// Query contains the SQL code to execute
	// Can contain multiple statements separated by semicolons
	Query string `json:"query" binding:"required"`
}

// Validate performs basic validation on the request
func (r *ExecuteRequest) Validate() error {
	if r.Query == "" {
		return ErrEmptyQuery
	}

	if len(r.Query) > MaxQueryLength {
		return ErrQueryTooLong
	}

	return nil
}
