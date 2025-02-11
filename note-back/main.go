package main

import (
	"net/http"
	"note/auth"
	"note/db"
	"note/noteConfig"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	auth.ConfigInit()
	noteConfig.ConfigInit()

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{noteConfig.GetDomain()},
		AllowMethods:     []string{"PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", login)
	r.POST("/logout", logout)
	internal := r.Group("/")

	internal.Use(authMiddleware())
	internal.POST("/api/note", postNote)
	internal.PUT("/api/note/:id", putNote)
	internal.GET("/api/note", GetAllNotes)
	internal.GET("/api/isLogged", isLogged)
	internal.GET("/api/category", GetAllCategories)
	internal.POST("/api/noteCat", GetNotesByCategory)
	internal.DELETE("/api/note/:id", deleteNote)
	r.Run(":5003")

}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil || auth.Validate(token) == false {
			c.JSON(401, "Unauthorized")
			c.Abort()
		}
	}

}

func isLogged(c *gin.Context) {
	c.JSON(200, gin.H{"data": "ok"})
	return
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
	c.JSON(200, gin.H{"data": n})
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
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	token, err := auth.Login(outside.Username, outside.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Wrong credential"})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("auth_token", token, 3600, "/", noteConfig.GetDomain(), true, true)

	c.JSON(200, gin.H{"message": "Login successful"})
}
func logout(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err == nil {
		auth.Logout(token)
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("auth_token", "", -1, "/", noteConfig.GetDomain(), true, true)
	c.JSON(200, gin.H{"data": "Logged out ok"})
}

type Note struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Categories []string `json:"categories" binding:"required"`
}
type NoteCategory struct {
	Categories []string `json:"categories" binding:"required"`
}
