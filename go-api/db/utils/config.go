package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER                 string        `mapstructure:"DB_DRIVER"`
	DB_SOURCE                 string        `mapstructure:"DB_SOURCE"`
	SERVER_ADDRESS            string        `mapstructure:"SERVER_ADDRESS"`
	TOKEN_DURATION            time.Duration `mapstructure:"TOKEN_DURATION"`
	TOKEN_SYMMETRY_KEY        string        `mapstructure:"TOKEN_SYMMETRY_KEY"`
	PYTHON_APP_TOKEN_DURATION time.Duration `mapstructure:"PYTHON_APP_TOKEN_DURATION"`
	EMAIL_SENDER_NAME         string        `mapstructure:"EMAIL_SENDER_NAME"`
	EMAIL_SENDER_ADDRESS      string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EMAIL_SENDER_PASSWORD     string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func ReadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.Unmarshal(&config)

	return
}
