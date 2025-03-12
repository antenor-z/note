package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"note/db"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pquerna/otp/totp"
)

var authSecret AuthSecret
var activeSessions = make(map[string]bool)

func tokenGenerator() string {
	b := make([]byte, 128)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func Login(outside AuthExternal) (string, error) {
	if outside.Username == authSecret.Username &&
		outside.Password == authSecret.Password &&
		totp.Validate(outside.Passcode, authSecret.Totp) {
		token := tokenGenerator()
		activeSessions[token] = true
		db.InsertSession(token)
		db.CleanSessions()
		return token, nil
	}
	return "", errors.New("invalid auth")
}

func Logout(token string) (string, error) {
	if activeSessions[token] {
		delete(activeSessions, token)
		db.DeleteSession(token)
		return "logged out", nil
	}
	return "", errors.New("invalid auth")
}

func Validate(token string) bool {
	if activeSessions[token] {
		return true
	}
	if db.IsSessionValid(token) {
		activeSessions[token] = true
		return true
	}
	return false
}

func ConfigInit() {
	dat, err := os.ReadFile("auth.toml")
	if err != nil {
		panic("Error while opening config file")
	}
	_, err2 := toml.Decode(string(dat), &authSecret)
	if err2 != nil {
		panic("Error while reading toml")
	}
}

type AuthExternal struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Passcode string `json:"passcode" binding:"required"`
}

type AuthSecret struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Totp     string `json:"totp" binding:"required"`
}
