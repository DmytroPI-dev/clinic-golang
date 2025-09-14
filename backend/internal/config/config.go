package config

import (
	"github.com/spf13/viper"
)

// Config for all application configuration
// Values are to bee read from env or config via Viper

type Config struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	DB_DSN        string `mapstructure:"DB_DSN"`
	AdminRole     string `mapstructure:"ADMIN_ROLE"`
	SessionSecret string `mapstructure:"SESSION_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
