package db

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("Init(): " + err.Error())
	}

	db.AutoMigrate(&Note{}, &Attachment{}, &ActiveSession{})
}

func CreateUser(username string, password string, totp string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	result := db.Create(&User{Username: username, Password: string(hashedPassword), Totp: totp})
	return result.Error
}

func EditUser(username string, password string, totp string) error {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
		user.Password = string(hashedPassword)
	}
	if totp != "" {
		user.Totp = totp
	}
	return result.Error
}

func DeleteUser(username string) error {
	result := db.Where("username = ?", username).Delete(&User{})
	return result.Error
}

func InsertNote(title string, content string, categoryNames []string, userId uint) error {
	var categories []*Category

	if len(categoryNames) == 1 && categoryNames[0] == "" {
		categoryNames = []string{"(uncategorized)"}
	}

	for _, name := range categoryNames {
		c := Category{Name: name, UserID: userId}
		db.FirstOrCreate(&c, Category{Name: name, UserID: userId})
		categories = append(categories, &c)
	}

	note := Note{Title: title, Content: content, Categories: categories, UserID: userId}
	result := db.Create(&note)
	return result.Error
}

func InsertAttachment(noteId uint, name string, fileUUID string, userId uint) error {
	attachment := Attachment{NoteID: noteId, Name: name, FileUUID: fileUUID, UserID: userId}
	result := db.Create(&attachment)
	return result.Error
}

func DeleteAttachment(noteId uint, attachmentId int, userId uint) error {
	result := db.Where("note_id = ? AND user_id = ?", noteId, userId).
		Delete(&Attachment{}, attachmentId)
	return result.Error
}

func GetAttachment(noteId uint, attachmentId int, userId uint) (Attachment, error) {
	var attachment Attachment
	err := db.Where("id = ? AND note_id = ? AND user_id = ?", attachmentId, noteId, userId).
		First(&attachment).Error
	return attachment, err
}

func GetAttachments(noteId uint, userId uint) ([]Attachment, error) {
	var attachments []Attachment
	err := db.Where("note_id = ? AND user_id = ?", noteId, userId).
		Find(&attachments).Error
	return attachments, err
}

func DeleteAllAttachments(noteId uint, userId uint) error {
	result := db.Where("note_id = ? AND user_id = ?", noteId, userId).
		Delete(&Attachment{})
	return result.Error
}

func GetAllNotes(userId uint) ([]Note, error) {
	var notes []Note
	err := db.Model(&Note{}).
		Order("notes.updated_at DESC").
		Preload("Categories").
		Preload("Attachments").
		Where("user_id = ?", userId).
		Find(&notes).Error
	return notes, err
}

func GetAllCategories(userId uint) ([]string, error) {
	var categories []Category
	res := db.Where("user_id = ?", userId).Find(&categories)
	var categoryList []string
	for _, category := range categories {
		categoryList = append(categoryList, category.Name)
	}
	return categoryList, res.Error
}

func DeleteNote(noteId int, userId uint) error {
	var note Note
	if err := db.Preload("Categories").
		Where("user_id = ?", userId).
		First(&note, noteId).Error; err != nil {
		return err
	}

	result := db.Exec("DELETE FROM notes WHERE id=? AND user_id = ?", noteId, userId)
	if result.Error != nil {
		return result.Error
	}

	db.Exec("DELETE FROM note_categories WHERE note_id=? AND user_id = ?", noteId, userId)

	for _, category := range note.Categories {
		var count int64
		db.Model(&Note{}).Joins("JOIN note_categories ON notes.id = note_categories.note_id").
			Where("note_categories.category_id = ?", category.ID).
			Where("user_id = ?", userId).
			Count(&count)
		if count == 0 {
			db.Delete(&Category{}, category.ID)
		}
	}

	return nil
}

func UpdateNote(noteId int, title string, content string, categoryNames []string, userId uint) error {
	var note Note
	if err := db.Preload("Categories").Where("user_id = ?", userId).First(&note, noteId).Error; err != nil {
		return err
	}

	// Store previous categories before updating
	prevCategories := note.Categories

	note.Title = title
	note.Content = content

	var categories []*Category
	for _, name := range categoryNames {
		var category Category
		db.FirstOrCreate(&category, Category{Name: name, UserID: userId})
		categories = append(categories, &category)
	}

	err := db.Model(&note).Association("Categories").Replace(categories)
	if err != nil {
		return err
	}

	for _, prevCat := range prevCategories {
		var count int64
		db.Model(&Note{}).Joins("JOIN note_categories ON notes.id = note_categories.note_id").
			Where("note_categories.category_id = ? AND user_id = ?", prevCat.ID, userId).Count(&count)

		if count == 0 {
			db.Delete(&Category{}, prevCat.ID)
		}
	}
	return db.Save(&note).Error
}

func GetNotesByCategory(categoryNames []string, userId uint) ([]Note, error) {
	var notes []Note
	err := db.Joins("JOIN note_categories ON note_categories.note_id = notes.id").
		Joins("JOIN categories ON categories.id = note_categories.category_id").
		Where("categories.name IN ?", categoryNames).
		Where("notes.user_id = ?", userId).
		Order("notes.updated_at DESC").
		Preload("Categories").
		Preload("Attachments").
		Find(&notes).Error
	return notes, err
}

func InsertSession(userId uint, token string) error {
	session := ActiveSession{Token: token, UserID: userId}
	result := db.Create(&session)
	return result.Error
}

func GetUserId(token string) (uint, error) {
	var activeSession ActiveSession
	result := db.Where("token = ?", token).First(&activeSession)
	if result.Error != nil {
		return 0, result.Error
	}
	return activeSession.UserID, nil
}

func DeleteSession(token string) error {
	result := db.Where("token = ?", token).Delete(&ActiveSession{})
	return result.Error
}

func CleanSessions() error {
	result := db.Where("created_at < ?", time.Now().AddDate(0, 0, -15)).Delete(&ActiveSession{})
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("created_at IS NULL").Delete(&ActiveSession{})
	return result.Error
}

func GetUser(username string) (User, error) {
	user := User{}
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
