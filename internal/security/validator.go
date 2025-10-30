package security

import (
	"regexp"
	"strings"

	"mysql-tui-editor/server/internal/domain"
)

// Validator validates SQL queries for security threats
type Validator struct {
	dangerousPatterns []*regexp.Regexp
}

// NewValidator creates a new SQL validator
func NewValidator() *Validator {
	return &Validator{
		dangerousPatterns: compileDangerousPatterns(),
	}
}

// Validate checks if the query contains dangerous commands
func (v *Validator) Validate(query string) error {
	// Normalize query: convert to lowercase and remove extra whitespace
	normalizedQuery := normalizeQuery(query)

	// Check for DROP DATABASE specifically (most critical)
	if containsDropDatabase(normalizedQuery) {
		return domain.ErrDropDatabase
	}

	// Check for other dangerous patterns
	for _, pattern := range v.dangerousPatterns {
		if pattern.MatchString(normalizedQuery) {
			return domain.ErrDangerousCommand
		}
	}

	return nil
}

// compileDangerousPatterns compiles all dangerous SQL patterns
func compileDangerousPatterns() []*regexp.Regexp {
	dangerousCommands := []string{
		// System commands
		`\bSHUTDOWN\b`,
		`\bRESTART\b`,

		// File operations
		`\bLOAD_FILE\s*\(`,
		`\bINTO\s+OUTFILE\b`,
		`\bINTO\s+DUMPFILE\b`,
		`\bLOAD\s+DATA\s+INFILE\b`,

		// User management
		`\bCREATE\s+USER\b`,
		`\bDROP\s+USER\b`,
		`\bALTER\s+USER\b`,
		`\bRENAME\s+USER\b`,
		`\bGRANT\b`,
		`\bREVOKE\b`,
		`\bSET\s+PASSWORD\b`,

		// Dangerous system variables
		`\bSET\s+GLOBAL\b`,
		`\bSET\s+@@GLOBAL\b`,

		// Process/thread manipulation
		`\bKILL\b`,

		// Plugin operations
		`\bINSTALL\s+PLUGIN\b`,
		`\bUNINSTALL\s+PLUGIN\b`,
	}

	patterns := make([]*regexp.Regexp, 0, len(dangerousCommands))
	for _, cmd := range dangerousCommands {
		// Case-insensitive matching
		pattern := regexp.MustCompile(`(?i)` + cmd)
		patterns = append(patterns, pattern)
	}

	return patterns
}

// containsDropDatabase checks specifically for DROP DATABASE command
func containsDropDatabase(query string) bool {
	// Multiple patterns to catch various forms
	dropDBPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bDROP\s+DATABASE\b`),
		regexp.MustCompile(`(?i)\bDROP\s+SCHEMA\b`),
	}

	for _, pattern := range dropDBPatterns {
		if pattern.MatchString(query) {
			return true
		}
	}

	return false
}

// normalizeQuery normalizes SQL query for validation
func normalizeQuery(query string) string {
	// Convert to lowercase for case-insensitive matching
	normalized := strings.ToLower(query)

	// Replace multiple whitespace with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized = spaceRegex.ReplaceAllString(normalized, " ")

	// Trim leading/trailing whitespace
	normalized = strings.TrimSpace(normalized)

	return normalized
}

// IsSafeQuery performs a comprehensive safety check
func (v *Validator) IsSafeQuery(query string) bool {
	return v.Validate(query) == nil
}
