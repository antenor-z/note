package routes

import (
	"note/db"
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
	res, _ := db.GetAllCategories()
	c.JSON(200, gin.H{"data": res})
}

func PostNote(c *gin.Context) {
	var n Note
	err := c.ShouldBindJSON(&n)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db.InsertNote(n.Title, n.Content, n.Categories)

	c.JSON(200, gin.H{"data": n})
}

func GetAllNotes(c *gin.Context) {
	notes, err := db.GetAllNotes()
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(200, gin.H{"data": notes})
}

func GetNotesByCategory(c *gin.Context) {
	var nc NoteCategory
	err := c.ShouldBindJSON(&nc)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	notes, err := db.GetNotesByCategory(nc.Categories)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	c.JSON(200, gin.H{"data": notes})
}

func PutNote(c *gin.Context) {
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

	db.UpdateNote(noteId, n.Title, n.Content, n.Categories)
	c.JSON(200, gin.H{"data": n})
}

func DeleteNote(c *gin.Context) {
	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	db.DeleteNote(noteId)
	c.JSON(200, gin.H{"status": "ok"})
}
