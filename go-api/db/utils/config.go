package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER                 string        `mapstructure:"DB_DRIVER"`
	PGUSER                    string        `mapstructure:"PGUSER"`
	PGHOST                    string        `mapstructure:"PGHOST"`
	PGPORT                    string        `mapstructure:"PGPORT"`
	POSTGRES_DB               string        `mapstructure:"POSTGRES_DB"`
	POSTGRES_PASSWORD         string        `mapstructure:"POSTGRES_PASSWORD"`
	DB_SOURCE                 string        `mapstructure:"DB_SOURCE"`
	DB_SOURCE_DEVELOPMENT     string        `mapstructure:"DB_SOURCE_DEVELOPMENT"`
	MIGRATION_SOURCE          string        `mapstructure:"MIGRATION_SOURCE"`
	SERVER_ADDRESS            string        `mapstructure:"SERVER_ADDRESS"`
	PUBLIC_URL                string        `mapstructure:"PUBLIC_URL"`
	TOKEN_DURATION            time.Duration `mapstructure:"TOKEN_DURATION"`
	TOKEN_SYMMETRY_KEY        string        `mapstructure:"TOKEN_SYMMETRY_KEY"`
	PYTHON_APP_TOKEN_DURATION time.Duration `mapstructure:"PYTHON_APP_TOKEN_DURATION"`
	EMAIL_SENDER_NAME         string        `mapstructure:"EMAIL_SENDER_NAME"`
	EMAIL_SENDER_ADDRESS      string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EMAIL_SENDER_PASSWORD     string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	MPESA_SHORT_CODE          string        `mapstructure:"MPESA_SHORT_CODE"`
	CONSUMER_KEY              string        `mapstructure:"CONSUMER_KEY"`
	CONSUMER_SECRET           string        `mapstructure:"CONSUMER_SECRET"`
	PASSKEY                   string        `mapstructure:"PASSKEY"`
	REDIS_ADDRESS             string        `mapstructure:"REDIS_ADDRESS"`
	REDIS_ADDRESS_DEVELOPMENT string        `mapstructure:"REDIS_ADDRESS_DEVELOPMENT"`
	REDIS_URI                 string        `mapstructure:"REDIS_URI"`
	REDIS_PASSWORD            string        `mapstructure:"REDIS_PASSWORD"`
	POSTGRES_SNAPSHOTS        string        `mapstructure:"POSTGRES_SNAPSHOTS"`
	HOST_SNAPSHOTS            string        `mapstructure:"HOST_SNAPSHOTS"`
}

func ReadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// viper.SetConfigFile("app.env")
	// viper.SetConfigType()

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.Unmarshal(&config)

	return config, nil
}
