package controller

import (
	"net/http"

	"chain-access/api/model"
	"chain-access/api/service"

	"github.com/gin-gonic/gin"
)

// ChainController 链信息控制器
type ChainController struct {
	ethService service.EthereumService
}

// NewChainController 创建链信息控制器
func NewChainController(ethService service.EthereumService) *ChainController {
	return &ChainController{ethService: ethService}
}

// HandleGetChains 返回支持的链列表
func (ctrl *ChainController) HandleGetChains(c *gin.Context) {
	chains := ctrl.ethService.GetChains()
	var infos []model.ChainInfo
	for _, ch := range chains {
		infos = append(infos, model.ChainInfo{
			ID:      ch.ID,
			Name:    ch.Name,
			ChainID: ch.ChainID,
		})
	}
	c.JSON(http.StatusOK, model.ChainsResponse{Chains: infos})
}
