package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Addr    string        `mapstructure:"addr"`
	Storage StorageConfig `mapstructure:"storage"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	Logger  LoggerConfig  `mapstructure:"logger"`
}

type StorageConfig struct {
	DSN             string        `mapstructure:"dsn"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	AttemptCount    int           `mapstructure:"attempt_count"`
}

type JWTConfig struct {
	Secret                 string
	AccessTokenTimeExpire  time.Duration `mapstructure:"access_token_time_expire"`
	RefreshTokenTimeExpire time.Duration `mapstructure:"refresh_token_time_expire"`
}

// LoggerConfig logger configuration.
//
// Level - logging level.
type LoggerConfig struct {
	Level string `yaml:"level"`
}

// Load create configuration from file & environments.
func Load(path string) (*Config, error) {
	dir, file := filepath.Split(path)
	viper.SetConfigName(strings.TrimSuffix(file, filepath.Ext(file)))
	viper.AddConfigPath(dir)

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("fail to decode into struct, %w", err)
	}

	return &cfg, nil
}
