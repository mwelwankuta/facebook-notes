package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	yaml "gopkg.in/yaml.v3"
)

const (
	RedirectURI = "http://localhost:8080/callback"
	FbAuthURL   = "https://www.facebook.com/v16.0/dialog/oauth"
	FbTokenURL  = "https://graph.facebook.com/v16.0/oauth/access_token"
	FbGraphAPI  = "https://graph.facebook.com/me"
)

type Config struct {
	Port                  string `yaml:"port"`
	OpenGraphClientSecret string `yaml:"open_graph_client_secret"`
	OpenGraphClientID     string `yaml:"open_graph_client_id"`
	Database              string `yaml:"database"`
	RedisToken            string `yaml:"redis_token"`
	RedisUrl              string `yaml:"redis_url"`
	JwtSecret             string `yaml:"jwt_secret"`
}

// LoadConfig loads the configuration from a file or environment variables if the file is not found
func LoadConfig(path string) (*Config, error) {
	var cfg Config

	content, err := os.ReadFile(path)
	if err != nil {
		panic("Could not read config file info")
	}

	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return &cfg, err
	}

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
	if jwtSecret := os.Getenv("jwt_secret"); jwtSecret != "" {
		cfg.JwtSecret = jwtSecret
	}

	return &cfg, nil
}
