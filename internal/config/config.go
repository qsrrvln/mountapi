package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	APIKey     string `mapstructure:"API_KEY"`
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
	viper.SetDefault("API_KEY", "ubah-ini")

	err = viper.ReadInConfig()
	// It's okay if config file doesn't exist, we fallback to env vars or defaults
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = nil
	}

	err = viper.Unmarshal(&config)
	return
}
