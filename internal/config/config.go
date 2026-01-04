package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost      string `mapstructure:"DB_HOST"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	DBPort      string `mapstructure:"DB_PORT"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	APIKey      string `mapstructure:"API_KEY"`
	DBSSLMode   string `mapstructure:"DB_SSL_MODE"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "ubah-ini")
	viper.SetDefault("DB_PASSWORD", "ubah-ini")
	viper.SetDefault("DB_NAME", "ubah-ini")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.BindEnv("SERVER_PORT", "PORT") // Bind PORT (Railway/Heroku standard) to SERVER_PORT
	viper.BindEnv("DATABASE_URL")        // Auto-bind DATABASE_URL
	viper.SetDefault("API_KEY", "ubah-ini")
	viper.SetDefault("DB_SSL_MODE", "disable")

	err = viper.ReadInConfig()
	// It's okay if config file doesn't exist, we fallback to env vars or defaults
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = nil
	}

	err = viper.Unmarshal(&config)

	// Helper function to read secret from file if the _FILE env var is set
	readSecret := func(envVar, fileEnvVar string, target *string) {
		if fileVal := viper.GetString(fileEnvVar); fileVal != "" {
			if content, err := os.ReadFile(fileVal); err == nil {
				*target = strings.TrimSpace(string(content))
			}
		}
	}

	readSecret("DB_PASSWORD", "DB_PASSWORD_FILE", &config.DBPassword)
	readSecret("API_KEY", "API_KEY_FILE", &config.APIKey)

	return
}
