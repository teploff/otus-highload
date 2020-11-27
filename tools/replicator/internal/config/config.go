package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Config struct {
	SlaveID    uint32          `mapstructure:"slave_id"`
	BinlogFile string          `mapstructure:"binlog_file"`
	BinlogPos  uint32          `mapstructure:"binlog_pos"`
	MySQL      MySQLConfig     `mapstructure:"mysql"`
	Tarantool  TarantoolConfig `mapstructure:"tarantool"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type TarantoolConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
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
