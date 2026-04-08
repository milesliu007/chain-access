package controller

import (
	"net/http"

	"chain-access/api/model"
	"chain-access/api/service"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// HandleChallenge 处理 challenge 请求
func (ctrl *AuthController) HandleChallenge(c *gin.Context) {
	var req model.ChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request parameters: " + err.Error()})
		return
	}

	if !service.IsValidAddress(req.Address) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Ethereum address format"})
		return
	}

	challenge := ctrl.authService.GenerateChallenge(req.Address)
	c.JSON(http.StatusOK, model.ChallengeResponse{Challenge: challenge})
}

// HandleVerify 处理签名验证请求
func (ctrl *AuthController) HandleVerify(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request parameters: " + err.Error()})
		return
	}

	if !service.IsValidAddress(req.Address) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Ethereum address format"})
		return
	}

	token, err := ctrl.authService.VerifySignature(req.Address, req.Signature)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.VerifyResponse{Token: token})
}
