package dto

import "time"

type Note struct {
	Title      string     `json:"title" binding:"required"`
	Priority   uint       `json:"priority"`
	IsHidden   bool       `json:"isHidden"`
	Deadline   *time.Time `json:"deadline"`
	Content    string     `json:"content"`
	Categories []string   `json:"categories" binding:"required"`
}
type NoteCategory struct {
	Categories []string `json:"categories" binding:"required"`
	ShowHidden bool     `json:"showHidden"`
	Priority   *uint    `json:"priority"`
}
