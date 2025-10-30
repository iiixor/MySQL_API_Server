package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	Executor ExecutorConfig `mapstructure:"executor"`
	Security SecurityConfig `mapstructure:"security"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// MySQLConfig holds MySQL connection configuration
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// ExecutorConfig holds query execution configuration
type ExecutorConfig struct {
	QueryTimeout time.Duration `mapstructure:"query_timeout"`
	DBPrefix     string        `mapstructure:"db_prefix"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimitPerSecond int `mapstructure:"rate_limit_per_second"`
	RateLimitBurst     int `mapstructure:"rate_limit_burst"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Allow environment variable overrides
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "35s")
	viper.SetDefault("server.write_timeout", "35s")
	viper.SetDefault("server.shutdown_timeout", "5s")

	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "rootpassword")
	viper.SetDefault("mysql.max_open_conns", 25)
	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.conn_max_lifetime", "5m")

	viper.SetDefault("executor.query_timeout", "30s")
	viper.SetDefault("executor.db_prefix", "student_db_")

	viper.SetDefault("security.rate_limit_per_second", 10)
	viper.SetDefault("security.rate_limit_burst", 20)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}
