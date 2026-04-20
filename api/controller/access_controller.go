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
	var req model.CheckAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}

	jwtAddress := c.GetString("address")

	address := jwtAddress
	if req.Address != "" {
		if !service.IsValidAddress(req.Address) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid wallet address format"})
			return
		}
		if strings.ToLower(req.Address) != strings.ToLower(jwtAddress) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "not authorized to query other addresses"})
			return
		}
		address = req.Address
	}

	if !service.IsValidAddress(req.ContractAddress) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid contract address format"})
		return
	}

	hasAccess, err := ctrl.ethService.CheckERC20Balance(req.ChainID, address, req.ContractAddress)
	if err != nil {
		log.Printf("链上查询失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "on-chain query failed"})
		return
	}

	c.JSON(http.StatusOK, model.CheckAccessResponse{HasAccess: hasAccess})
}

// HandleCheckNFT 处理 ERC-721 NFT 查询请求
func (ctrl *AccessController) HandleCheckNFT(c *gin.Context) {
	var req model.CheckNFTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}

	jwtAddress := c.GetString("address")

	address := jwtAddress
	if req.Address != "" {
		if !service.IsValidAddress(req.Address) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid wallet address format"})
			return
		}
		if strings.ToLower(req.Address) != strings.ToLower(jwtAddress) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "not authorized to query other addresses"})
			return
		}
		address = req.Address
	}

	if !service.IsValidAddress(req.ContractAddress) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid contract address format"})
		return
	}

	hasNFT, tokenIDs, err := ctrl.ethService.CheckERC721Ownership(req.ChainID, address, req.ContractAddress)
	if err != nil {
		log.Printf("ERC-721 query failed: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "on-chain query failed"})
		return
	}

	ids := make([]string, 0, len(tokenIDs))
	for _, id := range tokenIDs {
		ids = append(ids, id.String())
	}

	c.JSON(http.StatusOK, model.CheckNFTResponse{HasNFT: hasNFT, TokenIDs: ids})
}

// HandleCheckNFT1155 处理 ERC-1155 NFT 查询请求
func (ctrl *AccessController) HandleCheckNFT1155(c *gin.Context) {
	var req model.CheckNFT1155Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return
	}

	jwtAddress := c.GetString("address")

	address := jwtAddress
	if req.Address != "" {
		if !service.IsValidAddress(req.Address) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid wallet address format"})
			return
		}
		if strings.ToLower(req.Address) != strings.ToLower(jwtAddress) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "not authorized to query other addresses"})
			return
		}
		address = req.Address
	}

	if !service.IsValidAddress(req.ContractAddress) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid contract address format"})
		return
	}

	hasNFT, balance, err := ctrl.ethService.CheckERC1155Balance(req.ChainID, address, req.ContractAddress, req.TokenID)
	if err != nil {
		log.Printf("ERC-1155 query failed: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "on-chain query failed"})
		return
	}

	balanceStr := "0"
	if balance != nil {
		balanceStr = balance.String()
	}

	c.JSON(http.StatusOK, model.CheckNFT1155Response{HasNFT: hasNFT, Balance: balanceStr})
}
