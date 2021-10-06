package config

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	RabbitMQ RabbitMQ
	Service  ServiceConfig
	Jeager   JeagerConfig
}

type RabbitMQ struct {
	Host           string
	Port           string
	User           string
	Password       string
	Exchange       string
	Queue          string
	RoutingKey     string
	ConsumerTag    string
	WorkerPoolSize int
}

type ServiceConfig struct {
	Port string
}

type JeagerConfig struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

func LoadConfig(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode config into struct, %v", err)
		return nil, err
	}
	return &c, nil
}

func GetConfig() (*Config, error) {
	configPath := GetConfigPath(os.Getenv("build_env"))
	v, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	config, err := ParseConfig(v)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetConfigPath(env string) string {
	if env == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
