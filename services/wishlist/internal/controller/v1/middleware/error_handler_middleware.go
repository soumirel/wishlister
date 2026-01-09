package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		var statusCode int
		switch err.Type {
		case gin.ErrorTypeBind:
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		// TODO: Errors code defining
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})

	}
}
