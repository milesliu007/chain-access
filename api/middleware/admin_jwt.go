package middleware

import (
	"net/http"
	"strings"

	"chain-access/api/model"
	"chain-access/api/service"

	"github.com/gin-gonic/gin"
)

// AdminJWTMiddleware validates JWT for admin routes (same JWT, just checks presence)
func AdminJWTMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid Authorization format"})
			c.Abort()
			return
		}

		address, err := authService.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "JWT validation failed"})
			c.Abort()
			return
		}

		c.Set("address", address)
		c.Next()
	}
}
