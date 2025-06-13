package db

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Totp     string `json:"totp"`
}
type Category struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	Name   string  `json:"name"`
	Notes  []*Note `gorm:"many2many:note_categories;" json:"notes"`
	UserID uint
	User   User
}
type Note struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Categories  []*Category  `gorm:"many2many:note_categories;" json:"categories"`
	Attachments []Attachment `json:"attachments"`
	UserID      uint
	User        User
}
type Attachment struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	NoteID   uint   `json:"noteId"`
	Name     string `json:"name"`
	FileUUID string `json:"path"`
	UserID   uint
	User     User
}
type ActiveSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    uint
	User      User
}
