package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Port                  string `yaml:"port"`
	OpenGraphClientSecret string `yaml:"open_graph_client_secret"`
	OpenGraphClientID     string `yaml:"open_graph_client_id"`
	Database              string `yaml:"database"`
	RedisToken            string `yaml:"redis_token"`
	RedisUrl              string `yaml:"redis_url"`
}

func LoadConfig(path string) *Config {
	var cfg Config

	content, err := os.ReadFile(path)
	if err != nil {
		panic("Could not read config file info")
	}

	yaml.Unmarshal(content, &cfg)

	if port := os.Getenv("port"); port != "" {
		cfg.Port = port
	}
	if openGraphClientID := os.Getenv("open_graph_client_id"); openGraphClientID != "" {
		cfg.OpenGraphClientID = openGraphClientID
	}
	if openGraphClientSecret := os.Getenv("open_graph_client_secret"); openGraphClientSecret != "" {
		cfg.OpenGraphClientSecret = openGraphClientSecret
	}
	if database := os.Getenv("database"); database != "" {
		cfg.Database = database
	}
	if redisToken := os.Getenv("redis_token"); redisToken != "" {
		cfg.RedisToken = redisToken
	}
	if redisUrl := os.Getenv("redis_url"); redisUrl != "" {
		cfg.RedisUrl = redisUrl
	}

	return &cfg
}
