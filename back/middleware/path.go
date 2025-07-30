package middleware

import "github.com/gin-gonic/gin"

func GetPathMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Query("path")
		if path == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing path"})
			return
		}

		c.Set("unsafePath", path)
		c.Next()
	}
}
