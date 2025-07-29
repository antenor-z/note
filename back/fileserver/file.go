package fileserver

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func sanitizePath(unsafePath string) string {
	// No hidden files. Path traversal not allowed
	unsafePath = strings.ReplaceAll(unsafePath, ".", "")
	partsUnsafe := strings.Split(unsafePath, "/")
	parts := []string{}
	for _, part := range partsUnsafe {
		part = strings.Trim(part, "/")
		part = strings.Trim(part, " ")
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, "/")
}

func createUserRootIfNotExists(userID uint) {
	os.Mkdir("uploads/fileserver/"+fmt.Sprint(userID), os.ModePerm)
}

func getPath(path string, userID uint) string {
	createUserRootIfNotExists(userID)
	return sanitizePath(fmt.Sprintf("uploads/fileserver/%d/%s", userID, path))
}

type File struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"isDirectory"`
	ModTime time.Time `json:"modifiedOn"`
	Size    int64     `json:"size"`
}

func Ls(path string, userID uint) ([]File, error) {
	safePath := getPath(path, userID)
	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, err
	}
	var files []File
	for _, entry := range entries {
		info, _ := entry.Info()
		files = append(files, File{
			Name:    entry.Name(),
			Path:    fmt.Sprintf("%s/%s", safePath, entry.Name()),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})

	}
	return files, nil
}

func Mkdir(path string, userID uint) error {
	return os.Mkdir(getPath(path, userID), os.ModePerm)
}

func WriteFile(path string, userID uint, fileContent []byte) error {
	return os.WriteFile(getPath(path, userID), fileContent, os.ModePerm)
}

func GetFullPathAndName(path string, userID uint) (string, string, error) {
	safePath := getPath(path, userID)
	_, err := os.Stat(safePath)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(safePath, "/")
	return safePath, parts[len(parts)-1], nil
}

func Rm(path string, userID uint) error {
	return os.Remove(getPath(path, userID))
}
