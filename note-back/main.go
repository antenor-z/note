package main

import (
	"note/auth"
	"note/db"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	auth.ConfigInit()

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/login", login)
	r.POST("/note", postNote)
	r.PUT("/note/:id", putNote)
	r.GET("/note", GetAllNotes)
	r.GET("/category", GetAllCategories)
	r.POST("/noteCat", GetNotesByCategory)
	r.DELETE("/note/:id", deleteNote)

	r.Run(":5000")

}

func GetAllCategories(c *gin.Context) {
	res, _ := db.GetAllCategories()
	c.JSON(200, gin.H{"data": res})
}

// curl -i -X POST -H "Content-Type: application/json" -d "{ \"title\": \"aaaa\", \"content\": \"aaad\", \"categories\": [\"a\", \"b\"] }" localhost:5000/note
func postNote(c *gin.Context) {
	var n Note
	err := c.ShouldBindJSON(&n)
	if err != nil {
		c.String(400, "error")
		return
	}

	db.InsertNote(n.Title, n.Content, n.Categories)

	c.JSON(200, gin.H{"data": n})
}

func GetAllNotes(c *gin.Context) {
	notes, err := db.GetAllNotes()
	if err != nil {
		panic("error on getNote")
	}
	c.JSON(200, gin.H{"data": notes})
}

// curl -i -X POST -H "Content-Type: application/json" -d "{ \"categories\": [\"c\"] }" localhost:5000/note
func GetNotesByCategory(c *gin.Context) {
	var nc NoteCategory
	err := c.ShouldBindJSON(&nc)
	if err != nil {
		c.String(400, "error")
		return
	}

	notes, err := db.GetNotesByCategory(nc.Categories)
	if err != nil {
		panic("error on getNote")
	}

	c.JSON(200, gin.H{"data": notes})
}

// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"title\": \"aaaa\", \"content\": \"aaad\", \"categories\": [\"f\", \"g\"]}" localhost:5000/note/:id
func putNote(c *gin.Context) {
	var n Note
	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	err2 := c.ShouldBindJSON(&n)
	if err2 != nil {
		c.String(400, "error")
		return
	}

	db.UpdateNote(noteId, n.Title, n.Content, n.Categories)
}

// curl -i -X DELETE -H "Content-Type: application/json" localhost:5000/note/:id
func deleteNote(c *gin.Context) {
	noteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}
	db.DeleteNote(noteId)
	c.JSON(200, gin.H{"status": "ok"})
}

// curl -i -X POST -H "Content-Type: application/json" -d "{ \"username\": \"a\", \"password\": \"123\"}" localhost:5000/login
func login(c *gin.Context) {
	var outside auth.Auth
	err := c.ShouldBindJSON(&outside)
	if err != nil {
		panic("No auth")
	}
	token, err := auth.Login(outside.Username, outside.Password)
	if err != nil {
		panic("No auth")
	}

	c.String(200, token)
}

type Note struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Categories []string `json:"categories" binding:"required"`
}
type NoteCategory struct {
	Categories []string `json:"categories" binding:"required"`
}
