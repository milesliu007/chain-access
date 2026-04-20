package router

import (
	"net/http"

	"chain-access/api/controller"
	"chain-access/api/middleware"
	"chain-access/api/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 创建并配置路由
func SetupRouter(
	allowedOrigins []string,
	authService service.AuthService,
	ethService service.EthereumService,
) *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 控制器
	authCtrl := controller.NewAuthController(authService)
	accessCtrl := controller.NewAccessController(ethService)
	chainCtrl := controller.NewChainController(ethService)

	// 认证路由（无需 JWT）
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/challenge", authCtrl.HandleChallenge)
		authGroup.POST("/verify", authCtrl.HandleVerify)
	}

	// 公开 API（无需 JWT）
	r.GET("/chains", chainCtrl.HandleGetChains)

	// 权限查询路由（需要 JWT）
	protectedGroup := r.Group("/")
	protectedGroup.Use(middleware.JWTMiddleware(authService))
	{
		protectedGroup.POST("/check-access", accessCtrl.HandleCheckAccess)
		protectedGroup.POST("/check-nft", accessCtrl.HandleCheckNFT)
		protectedGroup.POST("/check-nft1155", accessCtrl.HandleCheckNFT1155)
	}

	// 前端静态文件
	r.StaticFile("/", "./frontend/dist/index.html")
	r.Static("/assets", "./frontend/dist/assets")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}
