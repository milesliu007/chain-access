package main

import (
	"log"

	"chain-access/api/config"
	"chain-access/api/repository"
	"chain-access/api/router"
	"chain-access/api/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化数据层
	challengeRepo := repository.NewMemoryChallengeRepository()

	// 初始化业务逻辑层
	authService := service.NewAuthService(cfg.JWTSecret, challengeRepo)

	ethService, err := service.NewEthereumService(cfg.Chains, cfg.HTTPProxy)
	if err != nil {
		log.Fatalf("Ethereum 服务初始化失败: %v", err)
	}
	defer ethService.Close()

	// 创建路由并启动
	r := router.SetupRouter(cfg.AllowedOrigins, authService, ethService)

	log.Printf("服务启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
