package security

import (
	"testing"

	"mysql-tui-editor/server/internal/domain"
)

func TestValidator_DropDatabase(t *testing.T) {
	validator := NewValidator()

	dangerousQueries := []string{
		"DROP DATABASE test;",
		"drop database test;",
		"DROP   DATABASE   test;",
		"DROP DATABASE IF EXISTS test;",
		"DROP SCHEMA test;",
		"drop schema test;",
		"SELECT * FROM users; DROP DATABASE test;",
	}

	for _, query := range dangerousQueries {
		err := validator.Validate(query)
		if err != domain.ErrDropDatabase {
			t.Errorf("Expected ErrDropDatabase for query: %s, got: %v", query, err)
		}
	}
}

func TestValidator_DangerousCommands(t *testing.T) {
	validator := NewValidator()

	dangerousQueries := []string{
		"SHUTDOWN;",
		"LOAD_FILE('/etc/passwd');",
		"SELECT * INTO OUTFILE '/tmp/data.txt' FROM users;",
		"CREATE USER 'hacker'@'localhost';",
		"GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';",
		"SET GLOBAL max_connections = 10000;",
		"KILL 123;",
		"INSTALL PLUGIN malicious SONAME 'plugin.so';",
	}

	for _, query := range dangerousQueries {
		err := validator.Validate(query)
		if err == nil {
			t.Errorf("Expected error for dangerous query: %s", query)
		}
	}
}

func TestValidator_SafeQueries(t *testing.T) {
	validator := NewValidator()

	safeQueries := []string{
		"SELECT * FROM users;",
		"CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100));",
		"INSERT INTO users VALUES (1, 'John');",
		"UPDATE users SET name = 'Jane' WHERE id = 1;",
		"DELETE FROM users WHERE id = 1;",
		"DROP TABLE users;",
		"ALTER TABLE users ADD COLUMN email VARCHAR(255);",
		"CREATE INDEX idx_name ON users(name);",
		"SHOW TABLES;",
		"DESCRIBE users;",
	}

	for _, query := range safeQueries {
		err := validator.Validate(query)
		if err != nil {
			t.Errorf("Expected no error for safe query: %s, got: %v", query, err)
		}
	}
}

func TestValidator_CaseInsensitive(t *testing.T) {
	validator := NewValidator()

	queries := []string{
		"dRoP dAtAbAsE test;",
		"DrOp DaTaBaSe test;",
		"LOAD_file('/etc/passwd');",
		"load_FILE('/etc/passwd');",
	}

	for _, query := range queries {
		err := validator.Validate(query)
		if err == nil {
			t.Errorf("Expected error for query with mixed case: %s", query)
		}
	}
}

func TestValidator_MultipleStatements(t *testing.T) {
	validator := NewValidator()

	query := `
		CREATE TABLE users (id INT PRIMARY KEY);
		INSERT INTO users VALUES (1);
		DROP DATABASE malicious;
		SELECT * FROM users;
	`

	err := validator.Validate(query)
	if err != domain.ErrDropDatabase {
		t.Errorf("Expected ErrDropDatabase for multi-statement query with DROP DATABASE")
	}
}

func TestValidator_IsSafeQuery(t *testing.T) {
	validator := NewValidator()

	if !validator.IsSafeQuery("SELECT * FROM users;") {
		t.Error("Expected safe query to return true")
	}

	if validator.IsSafeQuery("DROP DATABASE test;") {
		t.Error("Expected dangerous query to return false")
	}
}
