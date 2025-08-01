package routes

import (
	"note/fileserver"
	"path"

	"github.com/gin-gonic/gin"
)

func Ls(c *gin.Context) {
	unsafePath := c.GetString("unsafePath")
	userId := c.GetUint("userId")
	fileList, err := fileserver.Ls(unsafePath, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(200, gin.H{"data": fileList})
}

func Mkdir(c *gin.Context) {
	unsafePath := c.GetString("unsafePath")
	userId := c.GetUint("userId")
	err := fileserver.Mkdir(unsafePath, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "ok"})
}

func Rm(c *gin.Context) {
	unsafePath := c.GetString("unsafePath")
	userId := c.GetUint("userId")
	err := fileserver.Rm(unsafePath, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "ok"})
}

func ReadFile(c *gin.Context) {
	unsafePath := c.GetString("unsafePath")
	userId := c.GetUint("userId")
	fullPath, fileName, err := fileserver.GetFullPathAndName(unsafePath, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.FileAttachment(fullPath, fileName)
}

func WriteFile(c *gin.Context) {
	unsafePath := c.GetString("unsafePath")
	userId := c.GetUint("userId")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid file"})
		return
	}

	fullPath, _, err := fileserver.GetFullPathAndName(unsafePath, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	c.SaveUploadedFile(file, path.Join(fullPath, path.Clean(file.Filename)))

	c.JSON(200, gin.H{"data": "Upload OK"})
}
