package conf

import (
	"path/filepath"

	"github.com/xqa/chathub/cmd/flags"
	"github.com/xqa/chathub/pkg/utils/random"
)

type Database struct {
	DBFile      string `json:"db_file" env:"DB_FILE"`
	TablePrefix string `json:"table_prefix" env:"DB_TABLE_PREFIX"`
}

type LogConfig struct {
	Enable     bool   `json:"enable" env:"LOG_ENABLE"`
	Name       string `json:"name" env:"LOG_NAME"`
	MaxSize    int    `json:"max_size" env:"MAX_SIZE"`
	MaxBackups int    `json:"max_backups" env:"MAX_BACKUPS"`
	MaxAge     int    `json:"max_age" env:"MAX_AGE"`
	Compress   bool   `json:"compress" env:"COMPRESS"`
}

type Config struct {
	Force   bool   `json:"force" env:"FORCE"`
	Address string `json:"address" env:"ADDR"`
	Port    int    `json:"port" env:"PORT"`

	JwtSecret      string    `json:"jwt_secret" env:"JWT_SECRET"`
	TokenExpiresIn int       `json:"token_expires_in" env:"TOKEN_EXPIRES_IN"`
	Database       Database  `json:"database"`
	Log            LogConfig `json:"log"`
	MaxConnections int       `json:"max_connections" env:"MAX_CONNECTIONS"`
}

func DefaultConfig() *Config {
	logPath := filepath.Join(flags.DataDir, "log/log.log")
	dbPath := filepath.Join(flags.DataDir, "data.db")
	return &Config{
		Address:        "0.0.0.0",
		Port:           5244,
		JwtSecret:      random.String(16),
		TokenExpiresIn: 48,
		Database: Database{
			TablePrefix: "x_",
			DBFile:      dbPath,
		},
		Log: LogConfig{
			Enable:     true,
			Name:       logPath,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     28,
		},
		MaxConnections: 0,
	}
}
