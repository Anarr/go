package config

import "github.com/spf13/viper"

type (
	rabbitMqConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	redisConfig struct {
		Port int
	}

	AppConfig struct {
		RabbitMQ rabbitMqConfig `json:"rabbitmq"`
		Redis    redisConfig    `json:"redis"`
	}
)

// InitAppConfig initialize application configuration
func InitAppConfig() (*AppConfig, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	ac := &AppConfig{}
	err = viper.Unmarshal(ac)
	if err != nil {
		return nil, err
	}
	return ac, nil
}
