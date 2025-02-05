package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

func tokenGenerator() string {
	b := make([]byte, 128)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

var authSecret Auth
var activeSessions = make(map[string]bool)

func Login(username string, password string) (string, error) {
	if username == authSecret.Username && password == authSecret.Password {
		token := tokenGenerator()
		activeSessions[token] = true
		return token, nil
	}
	return "", errors.New("invalid auth")
}

func Logout(token string) (string, error) {
	if activeSessions[token] {
		delete(activeSessions, token)
	}
	return "", errors.New("invalid auth")
}

func Validate(token string) bool {
	return activeSessions[token]
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

type Auth struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
