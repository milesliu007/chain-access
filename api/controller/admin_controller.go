package controller

import (
	"net/http"
	"strconv"

	"chain-access/api/model"
	"chain-access/api/service"

	"github.com/gin-gonic/gin"
)

// AdminController handles admin API requests
type AdminController struct {
	adminService service.AdminService
	authService  service.AuthService
}

// NewAdminController creates admin controller
func NewAdminController(adminSvc service.AdminService, authSvc service.AuthService) *AdminController {
	return &AdminController{adminService: adminSvc, authService: authSvc}
}

// HandleAdminLogin handles admin login via ERC721
// POST /admin/login
func (ctrl *AdminController) HandleAdminLogin(c *gin.Context) {
	var req model.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := ctrl.adminService.VerifyAdminAccess(req.Address, req.Signature, req.ChainID, req.NFTContract)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.AdminLoginResponse{Token: token})
}

// HandleListBalances returns paginated user balances
// GET /admin/balances?page=1&size=20&address=0x...
func (ctrl *AdminController) HandleListBalances(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	address := c.Query("address")

	resp, err := ctrl.adminService.ListBalances(page, size, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "failed to fetch balances"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
