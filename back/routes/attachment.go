package routes

import (
	"note/db"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostAttachment(c *gin.Context) {
	userId := c.GetUint("userId")
	noteIdParam := c.Param("id")
	noteId, err := strconv.Atoi(noteIdParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid note ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid file"})
		return
	}

	name := file.Filename
	internalName := uuid.New().String()
	err = db.InsertAttachment(uint(noteId), name, internalName, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert attachment"})
		return
	}

	err = c.SaveUploadedFile(file, path.Join("uploads", internalName))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{"data": "Upload OK"})
}

func GetAttachmentFile(c *gin.Context) {
	userId := c.GetUint("userId")
	noteIdParam := c.Param("id")
	noteId, err := strconv.Atoi(noteIdParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid note ID"})
		return
	}

	attachmentIdParam := c.Param("attachmentId")
	attachmentId, err := strconv.Atoi(attachmentIdParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid attachment ID"})
		return
	}

	attachment, err := db.GetAttachment(uint(noteId), attachmentId, userId)
	if err != nil {
		c.JSON(404, gin.H{"error": "Attachment not found"})
		return
	}

	filePath := path.Join("uploads", attachment.FileUUID)
	c.FileAttachment(filePath, attachment.Name)
}

func DeleteAttachment(c *gin.Context) {
	userId := c.GetUint("userId")
	noteIdParam := c.Param("id")
	noteId, err := strconv.Atoi(noteIdParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid note id"})
		return
	}

	attachmentIdParam := c.Param("attachmentId")
	attachmentId, err := strconv.Atoi(attachmentIdParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid attachment id"})
		return
	}

	attachment, err := db.GetAttachment(uint(noteId), attachmentId, userId)
	if err != nil {
		c.JSON(404, gin.H{"error": "Attachment not found"})
		return
	}

	err = db.DeleteAttachment(uint(noteId), attachmentId, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete attachment"})
		return
	}

	err = os.Remove(path.Join("uploads", attachment.FileUUID))
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to delete attachment"})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
