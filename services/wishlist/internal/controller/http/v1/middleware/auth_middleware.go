package middleware

import (
	"net/http"

	"github.com/soumirel/wishlister/wishlist/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	userIdHeader = "X-User-Id"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader(userIdHeader)
		if userId == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		authCtx := auth.NewCtx(c.Request.Context(), auth.Auth{
			UserID: userId,
		})
		c.Request = c.Request.WithContext(authCtx)
		c.Next()
	}
}
