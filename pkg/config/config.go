package config

import (
	"fmt"
	"os"
	"strconv"

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

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use environment variables only
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Override with environment variables if they exist
	if grpcPortStr := os.Getenv("GRPC_PORT"); grpcPortStr != "" {
		if grpcPort, err := strconv.Atoi(grpcPortStr); err == nil {
			config.GRPCPort = grpcPort
		}
	}
	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		config.ServiceName = serviceName
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.DB.Host = dbHost
	}
	if dbPortStr := os.Getenv("DB_PORT"); dbPortStr != "" {
		if dbPort, err := strconv.Atoi(dbPortStr); err == nil {
			config.DB.Port = dbPort
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.DB.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" || os.Getenv("DB_PASSWORD") == "" {
		config.DB.Password = os.Getenv("DB_PASSWORD")
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.DB.Name = dbName
	}

	// Set defaults
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	if config.DB.Host == "" {
		config.DB.Host = "localhost"
	}

	if config.DB.Port == 0 {
		config.DB.Port = 5432
	}

	if config.DB.User == "" {
		config.DB.User = "mysticfunds"
	}

	if config.DB.Password == "" {
		config.DB.Password = "mysticfunds"
	}

	// No default GRPC port - should come from config file

	if config.DB.Name == "" {
		config.DB.Name = os.Getenv("SERVICE_NAME")
	}

	return &config, nil
}

func (c *Config) GetString(key, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}
