package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"note/db"
	"note/noteConfig"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"
)

var authSecret noteConfig.Auth = noteConfig.GetAuthSecret()

type safeSession struct {
	active map[string]bool
	mutex  sync.RWMutex
}

func (s *safeSession) add(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.active[token] = true
}

func (s *safeSession) exists(token string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.active[token]
}

func (s *safeSession) remove(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.active, token)
}

var sessionsCache safeSession = safeSession{active: make(map[string]bool)}

func cachePurge() {
	timer := time.NewTicker(time.Hour * 1)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			go func() {
				sessionsCache.mutex.Lock()
				defer sessionsCache.mutex.Unlock()
				clear(sessionsCache.active)
			}()
		}
	}
}
func init() {
	go cachePurge()
}

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
		sessionsCache.add(token)
		db.InsertSession(token)
		db.CleanSessions()
		return token, nil
	}
	return "", errors.New("invalid auth")
}

func Logout(token string) (string, error) {
	if sessionsCache.exists(token) {
		sessionsCache.remove(token)
		db.DeleteSession(token)
		return "logged out", nil
	}
	return "", errors.New("invalid auth")
}

func Validate(token string) bool {
	if sessionsCache.exists(token) {
		return true
	}
	if db.IsSessionValid(token) {
		sessionsCache.add(token)
		return true
	}
	return false
}
