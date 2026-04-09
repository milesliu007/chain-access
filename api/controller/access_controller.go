package controller

import (
	"log"
	"net/http"
	"strings"

	"chain-access/api/model"
	"chain-access/api/service"

	"github.com/gin-gonic/gin"
)

// AccessController 权限查询控制器
type AccessController struct {
	ethService service.EthereumService
}

// NewAccessController 创建权限查询控制器
func NewAccessController(ethService service.EthereumService) *AccessController {
	return &AccessController{ethService: ethService}
}

// HandleCheckAccess 处理权限查询请求
func (ctrl *AccessController) HandleCheckAccess(c *gin.Context) {
	// 从 JWT context 获取认证地址，防止越权查询
	jwtAddress := c.GetString("address")
	queryAddress := c.Query("address")
	contractAddress := c.Query("contract_address")
	chainID := c.DefaultQuery("chain_id", "ethereum")

	if contractAddress == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "missing required parameter: contract_address"})
		return
	}

	// 如果提供了 query address，必须与 JWT 地址一致
	address := jwtAddress
	if queryAddress != "" {
		if !service.IsValidAddress(queryAddress) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid wallet address format"})
			return
		}
		if strings.ToLower(queryAddress) != strings.ToLower(jwtAddress) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "not authorized to query other addresses"})
			return
		}
		address = queryAddress
	}

	if !service.IsValidAddress(contractAddress) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid contract address format"})
		return
	}

	hasAccess, err := ctrl.ethService.CheckERC20Balance(chainID, address, contractAddress)
	if err != nil {
		log.Printf("链上查询失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "on-chain query failed"})
		return
	}

	c.JSON(http.StatusOK, model.CheckAccessResponse{HasAccess: hasAccess})
}
