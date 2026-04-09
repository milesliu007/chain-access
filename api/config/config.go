package config

import (
	"fmt"
	"os"
	"strings"
)

// Config 应用配置
type Config struct {
	Chains         []ChainConfig
	JWTSecret      []byte
	Port           string
	AllowedOrigins []string
	HTTPProxy      string
}

// LoadConfig 从环境变量加载配置，缺少必要配置时拒绝启动
func LoadConfig() (*Config, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET is too short, at least 32 bytes required")
	}

	infuraKey := os.Getenv("INFURA_API_KEY")
	if infuraKey == "" {
		return nil, fmt.Errorf("INFURA_API_KEY environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := []string{"http://localhost:" + port}
	if allowedOrigins != "" {
		for _, o := range strings.Split(allowedOrigins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins = append(origins, o)
			}
		}
	}

	return &Config{
		Chains:         DefaultChains(infuraKey),
		JWTSecret:      []byte(jwtSecret),
		Port:           port,
		AllowedOrigins: origins,
		HTTPProxy:      os.Getenv("HTTP_PROXY"),
	}, nil
}
