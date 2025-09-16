// config contains  DatabaseConfig model
// it needed for connection to data base
// Why uses it os.Getenv?
// No external dependencies - Standard library only
// Simplicity - Easy to understand and maintain
// Lightweight - No additional package overhead
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// env
const (
	db_host         = "DB_HOST"
	db_port         = "DB_PORT"
	db_user         = "DB_USER"
	db_password     = "DB_PASSWORD"
	db_name         = "DB_NAME"
	db_sslmode      = "DB_SSL_MODE"
	db_maxconns     = "DB_MAX_CONNS"
	db_minconns     = "DB_MIN_CONNS"
	db_life         = "DB_CONN_MAX_LIFETIME"
	db_idle_time    = "DB_CONN_MAX_IDLE_TIME"
	db_timeout      = "DB_TIMEOUT"
	db_health_check = "DB_HEALTH_CHECK_PERIOD"
)

// default values
const (
	df_host                = "localhost"
	df_port                = 5432
	df_user                = "postgres"
	df_password            = ""
	df_name                = "mydb"
	df_sslmode             = "disable"
	df_maxconns            = 20
	df_minconns            = 1
	df_lifetime            = time.Hour
	df_lifeidletime        = 30 * time.Minute
	df_timeout             = 5 * time.Second
	df_health_check_period = time.Minute
)

// DatabaseConfig
type DatabaseConfig struct {
	Host              string
	Port              int
	User              string
	Password          string
	Name              string
	SSlMode           string
	MaxConns          int
	MinConns          int
	ConnMaxLifeTime   time.Duration
	ConnMaxIdleTime   time.Duration
	Timeout           time.Duration
	HealthCheckPeriod time.Duration
}

// LoadConfig returns data base config
func LoadConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:              getEnv(db_host, df_host),
		Port:              getEnvAsInt(db_port, df_port),
		User:              getEnv(db_user, df_user),
		Password:          getEnv(db_password, df_password),
		Name:              getEnv(db_name, df_name),
		SSlMode:           getEnv(db_sslmode, df_sslmode),
		MaxConns:          getEnvAsInt(db_maxconns, df_maxconns),
		MinConns:          getEnvAsInt(db_minconns, df_minconns),
		ConnMaxLifeTime:   getEnvAsDuration(db_life, df_lifetime),
		ConnMaxIdleTime:   getEnvAsDuration(db_idle_time, df_lifeidletime),
		Timeout:           getEnvAsDuration(db_timeout, df_timeout),
		HealthCheckPeriod: getEnvAsDuration(db_health_check, df_health_check_period),
	}
}

// ConnectionString reutrn a string for connect to data base
func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSlMode,
	)
}

// there are helpers...
func getEnv(key, defaultvalue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultvalue
}

func getEnvAsInt(key string, defaultvalue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultvalue
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return defaultvalue
	}
	return v
}

func getEnvAsDuration(key string, defaultvalue time.Duration) time.Duration {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultvalue
	}
	t, err := time.ParseDuration(value)
	if err != nil {
		return defaultvalue
	}
	return t
}
