package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"note/db"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type safeSession struct {
	active map[string]uint
	mutex  sync.RWMutex
}

func (s *safeSession) add(token string, userID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.active[token] = userID
}

func (s *safeSession) getUserId(token string) (uint, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.active[token] != 0 {
		return s.active[token], nil
	}
	return 0, errors.New("no user found")
}

func (s *safeSession) remove(token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.active, token)
}

var sessionsCache safeSession = safeSession{active: make(map[string]uint)}

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
	user, err := db.GetUser(outside.Username)
	if err != nil {
		return "", errors.New("invalid auth")
	}
	pwHashMismatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(outside.Password))
	if outside.Username == user.Username &&
		pwHashMismatch == nil &&
		totp.Validate(outside.Passcode, user.Totp) {
		token := tokenGenerator()
		sessionsCache.add(token, user.ID)
		db.InsertSession(user.ID, token)
		db.CleanSessions()
		return token, nil
	}
	return "", errors.New("invalid auth")
}

func Logout(token string) (string, error) {
	_, err := sessionsCache.getUserId(token)
	if err != nil {
		return "", errors.New("invalid auth")
	}
	sessionsCache.remove(token)
	db.DeleteSession(token)
	return "logged out", nil
}

func GetLoggedUserId(token string) (uint, error) {
	userId, err := sessionsCache.getUserId(token)
	if err == nil {
		return userId, nil
	}
	userIdFromDb, err := db.GetUserId(token)
	if err != nil {
		return 0, errors.New("invalid auth")
	}
	sessionsCache.add(token, userIdFromDb)
	return userIdFromDb, nil
}
