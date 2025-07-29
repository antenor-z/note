package routes

import (
	"note/fileserver"

	"github.com/gin-gonic/gin"
)

func Ls(c *gin.Context) {
	path, ok := c.Params.Get("path")
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	userId := c.GetUint("userId")
	fileList, err := fileserver.Ls(path, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(200, gin.H{"data": fileList})
}

func Mkdir(c *gin.Context) {
	path, ok := c.Params.Get("path")
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userId := c.GetUint("userId")
	err := fileserver.Mkdir(path, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "ok"})
}

func Rm(c *gin.Context) {
	path, ok := c.Params.Get("path")
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userId := c.GetUint("userId")
	err := fileserver.Rm(path, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "ok"})
}

func ReadFile(c *gin.Context) {
	path, ok := c.Params.Get("path")
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	userId := c.GetUint("userId")
	filePath, fileName, err := fileserver.GetFullPathAndName(path, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.FileAttachment(filePath, fileName)
}
