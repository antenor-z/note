package middleware

import (
	"note/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil || !auth.Validate(token) {
			c.JSON(401, "Unauthorized")
			c.Abort()
		}
	}
}
