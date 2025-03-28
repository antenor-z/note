package routes

import (
	"net/http"
	"note/auth"
	"note/noteConfig"

	"github.com/gin-gonic/gin"
)

func IsLogged(c *gin.Context) {
	c.JSON(200, gin.H{"data": "ok"})
}

func Login(c *gin.Context) {
	var outside auth.AuthExternal
	err := c.ShouldBindJSON(&outside)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	token, err := auth.Login(outside)
	if err != nil {
		c.JSON(401, gin.H{"error": "Wrong credential"})
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("auth_token", token, 604800, "/", noteConfig.GetDomain(), true, true)

	c.JSON(200, gin.H{"message": "Login successful"})
}
func Logout(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err == nil {
		auth.Logout(token)
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("auth_token", "", -1, "/", noteConfig.GetDomain(), true, true)
	c.JSON(200, gin.H{"data": "Logged out ok"})
}
