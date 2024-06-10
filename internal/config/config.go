// Package config provides functionality to load configuration settings
// from a file and environment variables using the viper library.
package config

import (
	"github.com/spf13/viper"
	"go-proxy-server/pkg/logger"
)

// Config struct holds the configuration settings.
type Config struct {
	Port      string
	LogLevel  string
	Servers   []string
	JwtSecret string
	RedisAddr string
	RedisPass string
	RedisDB   int
}

// Loader defines an interface for loading configuration settings.
type Loader interface {
	LoadConfig() *Config
}

// ViperLoader implements Loader interface using viper.
type ViperLoader struct{}

// NewViperLoader returns a new instance of ViperLoader.
func NewViperLoader() *ViperLoader {
	return &ViperLoader{}
}

// LoadConfig loads configuration settings from a config file and environment variables.
func (v *ViperLoader) LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.NewLogger("error").Fatalf("Error reading config file: %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.NewLogger("error").Fatalf("Unable to decode into config struct: %s", err)
	}

	return &config
}
