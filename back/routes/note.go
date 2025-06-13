package routes

import (
	"note/db"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Note struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content"`
	Categories []string `json:"categories" binding:"required"`
}
type NoteCategory struct {
	Categories []string `json:"categories" binding:"required"`
}

func GetAllCategories(c *gin.Context) {
	userId := c.GetUint("userId")
	res, _ := db.GetAllCategories(userId)
	c.JSON(200, gin.H{"data": res})
}

func PostNote(c *gin.Context) {
	userId := c.GetUint("userId")
	var n Note
	err := c.ShouldBindJSON(&n)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db.InsertNote(n.Title, n.Content, n.Categories, userId)

	c.JSON(200, gin.H{"data": n})
}

func GetAllNotes(c *gin.Context) {
	userId := c.GetUint("userId")
	notes, err := db.GetAllNotes(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(200, gin.H{"data": notes})
}

func GetNotesByCategory(c *gin.Context) {
	userId := c.GetUint("userId")
	var nc NoteCategory
	err := c.ShouldBindJSON(&nc)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	notes, err := db.GetNotesByCategory(nc.Categories, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(200, gin.H{"data": notes})
}

func PutNote(c *gin.Context) {
	userId := c.GetUint("userId")
	var n Note
	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err2 := c.ShouldBindJSON(&n)
	if err2 != nil {
		c.String(400, "error")
		return
	}

	db.UpdateNote(noteId, n.Title, n.Content, n.Categories, userId)
	c.JSON(200, gin.H{"data": n})
}

func DeleteNote(c *gin.Context) {
	userId := c.GetUint("userId")
	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	attachments, err := db.GetAttachments(uint(noteId), userId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Couldn't get attachments"})
		return
	}
	for _, attachment := range attachments {
		os.Remove(path.Join("uploads", attachment.FileUUID))
	}
	db.DeleteAllAttachments(uint(noteId), userId)
	db.DeleteNote(noteId, userId)
	c.JSON(200, gin.H{"status": "ok"})
}
