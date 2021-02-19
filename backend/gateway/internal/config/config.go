package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Addr      string          `mapstructure:"addr"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Messenger MessengerConfig `mapstructure:"messenger"`
	Social    SocialConfig    `mapstructure:"social"`
	Jaeger    JaegerConfig    `mapstructure:"jaeger"`
	Consul    ConsulConfig    `mapstructure:"consul"`
	Logger    LoggerConfig    `mapstructure:"logger"`
}

// AuthConfig authorization service address.
//
// Addr - socket address (host + port).
type AuthConfig struct {
	Addr string `mapstructure:"addr"`
}

// MessengerConfig messenger service address.
//
// Addr - socket address (host + port).
type MessengerConfig struct {
	HTTPAddr string `mapstructure:"http_addr"`
	GRPCAddr string `mapstructure:"grpc_addr"`
}

// SocialConfig social-network service address.
//
// Addr - socket address (host + port).
type SocialConfig struct {
	Addr string `mapstructure:"addr"`
}

type JaegerConfig struct {
	AgentAddr   string `mapstructure:"agent_addr"`
	ServiceName string `mapstructure:"service_name"`
}

type ConsulConfig struct {
	Addr        string `mapstructure:"addr"`
	ServiceName string `mapstructure:"service_name"`
}

// LoggerConfig logger configuration.
//
// Level - level logging.
type LoggerConfig struct {
	Level string `mapstructure:"level"`
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
