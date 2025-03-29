package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Category struct {
	ID    uint    `gorm:"primaryKey" json:"id"`
	Name  string  `json:"name"`
	Notes []*Note `gorm:"many2many:note_categories;" json:"notes"`
}
type Note struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Categories  []*Category  `gorm:"many2many:note_categories;" json:"categories"`
	Attachments []Attachment `json:"attachments"`
}
type Attachment struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	NoteID   uint   `json:"noteId"`
	Name     string `json:"name"`
	FileUUID string `json:"path"`
}
type ActiveSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("Init(): " + err.Error())
	}

	db.AutoMigrate(&Note{}, &Attachment{}, &ActiveSession{})
}

func InsertNote(title string, content string, categoryNames []string) error {
	var categories []*Category

	if len(categoryNames) == 1 && categoryNames[0] == "" {
		categoryNames = []string{"(uncategorized)"}
	}

	for _, name := range categoryNames {
		c := Category{Name: name}
		db.FirstOrCreate(&c, Category{Name: name})
		categories = append(categories, &c)
	}

	note := Note{Title: title, Content: content, Categories: categories}
	result := db.Create(&note)
	return result.Error
}

func InsertAttachment(noteId uint, name string, fileUUID string) error {
	attachment := Attachment{NoteID: noteId, Name: name, FileUUID: fileUUID}
	result := db.Create(&attachment)
	return result.Error
}

func DeleteAttachment(noteId uint, attachmentId int) error {
	result := db.Where("note_id = ?", noteId).Delete(&Attachment{}, attachmentId)
	return result.Error
}

func GetAttachment(noteId uint, attachmentId int) (Attachment, error) {
	var attachment Attachment
	err := db.Where("id = ? AND note_id = ?", attachmentId, noteId).First(&attachment).Error
	return attachment, err
}

func GetAllNotes() ([]Note, error) {
	var notes []Note
	err := db.Model(&Note{}).
		Order("notes.updated_at DESC").
		Preload("Categories").
		Preload("Attachments").
		Find(&notes).Error
	return notes, err
}

func GetAllCategories() ([]string, error) {
	var categories []Category
	res := db.Find(&categories)
	var categoryList []string
	for _, category := range categories {
		categoryList = append(categoryList, category.Name)
	}
	return categoryList, res.Error
}

func DeleteNote(noteId int) error {
	var note Note
	if err := db.Preload("Categories").First(&note, noteId).Error; err != nil {
		return err
	}

	result := db.Exec("DELETE FROM notes WHERE id=?", noteId)
	if result.Error != nil {
		return result.Error
	}

	db.Exec("DELETE FROM note_categories WHERE note_id=?", noteId)

	for _, category := range note.Categories {
		var count int64
		db.Model(&Note{}).Joins("JOIN note_categories ON notes.id = note_categories.note_id").
			Where("note_categories.category_id = ?", category.ID).Count(&count)
		if count == 0 {
			db.Delete(&Category{}, category.ID)
		}
	}

	return nil
}

func UpdateNote(noteId int, title string, content string, categoryNames []string) error {
	var note Note
	if err := db.Preload("Categories").First(&note, noteId).Error; err != nil {
		return err
	}

	// Store previous categories before updating
	prevCategories := note.Categories

	note.Title = title
	note.Content = content

	var categories []*Category
	for _, name := range categoryNames {
		var category Category
		db.FirstOrCreate(&category, Category{Name: name})
		categories = append(categories, &category)
	}

	err := db.Model(&note).Association("Categories").Replace(categories)
	if err != nil {
		return err
	}

	for _, prevCat := range prevCategories {
		var count int64
		db.Model(&Note{}).Joins("JOIN note_categories ON notes.id = note_categories.note_id").
			Where("note_categories.category_id = ?", prevCat.ID).Count(&count)

		if count == 0 {
			db.Delete(&Category{}, prevCat.ID)
		}
	}
	return db.Save(&note).Error
}

func GetNotesByCategory(categoryNames []string) ([]Note, error) {
	var notes []Note
	err := db.Joins("JOIN note_categories ON note_categories.note_id = notes.id").
		Joins("JOIN categories ON categories.id = note_categories.category_id").
		Where("categories.name IN ?", categoryNames).
		Order("notes.updated_at DESC").
		Preload("Categories").
		Preload("Attachments").
		Find(&notes).Error
	return notes, err
}

func InsertSession(token string) error {
	session := ActiveSession{Token: token}
	result := db.Create(&session)
	return result.Error
}

func IsSessionValid(token string) bool {
	var activeSession ActiveSession
	result := db.Where("token = ?", token).Find(&activeSession)
	if result.Error != nil {
		return false
	}
	return activeSession.Token == token
}

func DeleteSession(token string) error {
	result := db.Where("token = ?", token).Delete(&ActiveSession{})
	return result.Error
}

func CleanSessions() error {
	result := db.Where("created_at < ?", time.Now().AddDate(0, 0, -7)).Delete(&ActiveSession{})
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("created_at IS NULL").Delete(&ActiveSession{})
	return result.Error
}
