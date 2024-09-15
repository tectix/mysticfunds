package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string   `mapstructure:"SERVICE_NAME"`
	GRPCPort    int      `mapstructure:"GRPC_PORT"`
	LogLevel    string   `mapstructure:"LOG_LEVEL"`
	JWTSecret   string   `mapstructure:"JWT_SECRET"`
	DB          DBConfig `mapstructure:",squash"`
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	if config.DB.Name == "" {
		config.DB.Name = os.Getenv("SERVICE_NAME")
	}

	return &config, nil
}
