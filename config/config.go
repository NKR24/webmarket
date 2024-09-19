package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	KafkaBrokers string
	ElasticHost  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		DBHost:       viper.GetString("DB_HOST"),
		DBPort:       viper.GetString("DB_PORT"),
		DBUser:       viper.GetString("DB_USER"),
		DBPassword:   viper.GetString("DB_PASSWORD"),
		DBName:       viper.GetString("DB_NAME"),
		KafkaBrokers: viper.GetString("KAFKA_BROKERS"),
		ElasticHost:  viper.GetString("ELASTIC_HOST"),
	}

	return config, nil
}
