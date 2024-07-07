package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DBUrl string `mapstructure:"DB_URL"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

func init() {
	_, err := LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
}
