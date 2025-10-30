package executor

import (
	"database/sql"
	"fmt"
	"strings"
)

// formatResultSet formats SQL query results as a text table (MySQL CLI style)
func formatResultSet(rows *sql.Rows) (string, error) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("failed to get columns: %w", err)
	}

	if len(columns) == 0 {
		return "Empty result set", nil
	}

	// Collect all rows
	var results [][]string
	columnCount := len(columns)

	for rows.Next() {
		// Create a slice of sql.RawBytes to hold each column
		values := make([]sql.RawBytes, columnCount)
		valuePtrs := make([]interface{}, columnCount)
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert values to strings
		row := make([]string, columnCount)
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = string(val)
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating rows: %w", err)
	}

	// If no results
	if len(results) == 0 {
		return "Empty set", nil
	}

	// Calculate column widths
	colWidths := make([]int, columnCount)
	for i, col := range columns {
		colWidths[i] = len(col)
	}

	for _, row := range results {
		for i, val := range row {
			if len(val) > colWidths[i] {
				colWidths[i] = len(val)
			}
		}
	}

	// Build the table
	var output strings.Builder

	// Top border
	output.WriteString(buildBorder(colWidths))
	output.WriteString("\n")

	// Header row
	output.WriteString(buildRow(columns, colWidths))
	output.WriteString("\n")

	// Middle border
	output.WriteString(buildBorder(colWidths))
	output.WriteString("\n")

	// Data rows
	for _, row := range results {
		output.WriteString(buildRow(row, colWidths))
		output.WriteString("\n")
	}

	// Bottom border
	output.WriteString(buildBorder(colWidths))
	output.WriteString("\n")

	// Row count
	rowCount := len(results)
	if rowCount == 1 {
		output.WriteString("1 row in set")
	} else {
		output.WriteString(fmt.Sprintf("%d rows in set", rowCount))
	}

	return output.String(), nil
}

// buildBorder creates a border line for the table
func buildBorder(colWidths []int) string {
	var border strings.Builder
	border.WriteString("+")
	for _, width := range colWidths {
		border.WriteString(strings.Repeat("-", width+2))
		border.WriteString("+")
	}
	return border.String()
}

// buildRow creates a data row for the table
func buildRow(values []string, colWidths []int) string {
	var row strings.Builder
	row.WriteString("|")
	for i, val := range values {
		row.WriteString(" ")
		row.WriteString(val)
		row.WriteString(strings.Repeat(" ", colWidths[i]-len(val)))
		row.WriteString(" |")
	}
	return row.String()
}
