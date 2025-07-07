package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Postgres struct {
		DB_HOST     string
		DB_PORT     string
		DB_USER     string
		DB_PASSWORD string
		DB_NAME     string
	}
}

func InitConfig() (*Config, error) {

	viper.SetConfigFile("/root/config/config.yaml")

	err := viper.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("problems with read in config err: %v", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("problems with unmarshall config err: %v", err)
	}

	return &cfg, nil
}
