package routes

import (
	"note/noteConfig"

	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	c.JSON(200, gin.H{"version": noteConfig.GetVersion()})
}
