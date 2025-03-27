package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"note/db"
	"note/noteConfig"

	"github.com/pquerna/otp/totp"
)

var authSecret noteConfig.Auth = noteConfig.GetAuthSecret()
var activeSessions = make(map[string]bool)

type AuthExternal struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Passcode string `json:"passcode" binding:"required"`
}

func tokenGenerator() string {
	b := make([]byte, 128)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func Login(outside AuthExternal) (string, error) {
	if outside.Username == noteConfig.GetAuthSecret().Username &&
		outside.Password == noteConfig.GetAuthSecret().Password &&
		totp.Validate(outside.Passcode, noteConfig.GetAuthSecret().Totp) {
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
