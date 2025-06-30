package main

import (
	"note/db"
	"note/middleware"
	"note/noteConfig"
	"note/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	noteConfig.ConfigInit()

	if noteConfig.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.MaxMultipartMemory = 256 << 20 // 256MB file max
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{noteConfig.GetDomain()},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/login", routes.Login)
	r.GET("/api/version", routes.Version)
	internal := r.Group("/")

	internal.Use(middleware.AuthMiddleware())
	internal.POST("/api/logout", routes.Logout)
	internal.POST("/api/note", routes.PostNote)
	internal.PUT("/api/note/:id", routes.PutNote)
	internal.GET("/api/note", routes.GetAllNotes)
	internal.GET("/api/isLogged", routes.IsLogged)
	internal.GET("/api/category", routes.GetAllCategories)
	internal.GET("/api/category/hidden", routes.GetAllCategoriesWithHidden)
	internal.POST("/api/note/category", routes.GetNotesByCategory)
	internal.DELETE("/api/note/:id", routes.DeleteNote)
	internal.POST("/api/note/:id/attachment", routes.PostAttachment)
	internal.DELETE("/api/note/:id/attachment/:attachmentId", routes.DeleteAttachment)
	internal.GET("/api/note/:id/attachment/:attachmentId/file", routes.GetAttachmentFile)
	r.Run(":5003")
}
