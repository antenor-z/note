package fileserver

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const UPLOAD_FOLDER = "uploads/fileserver/"

func createUserRootIfNotExists(userID uint) {
	os.Mkdir(path.Join(UPLOAD_FOLDER, fmt.Sprint(userID)), os.ModePerm)
}

func getSafePath(unsafePath string, userID uint) string {
	createUserRootIfNotExists(userID)
	cleanPath := filepath.Clean(unsafePath)
	userRoot := filepath.Join(UPLOAD_FOLDER, fmt.Sprint(userID))
	fullPath := filepath.Join(userRoot, cleanPath)

	return fullPath
}

func formatSize(size int64) string {
	if size > 1024*1024*1024 {
		return fmt.Sprintf("%d GiB", size/1024/1024/1024)
	} else if size > 1024*1024 {
		return fmt.Sprintf("%d MiB", size/1024/1024)
	} else if size > 1024 {
		return fmt.Sprintf("%d KiB", size/1024)
	} else {
		return fmt.Sprintf("%d B", size)
	}
}

type File struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"isDirectory"`
	ModTime time.Time `json:"modifiedOn"`
	Size    string    `json:"size"`
}

func Ls(unsafePath string, userID uint) ([]File, error) {
	safePath := getSafePath(unsafePath, userID)
	entries, err := os.ReadDir(safePath)
	if err != nil {
		return nil, err
	}
	var files []File = make([]File, 0)
	var directories []File = make([]File, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if entry.IsDir() {
			directories = append(directories, File{
				Name:    entry.Name(),
				Path:    path.Join(unsafePath, entry.Name()),
				IsDir:   entry.IsDir(),
				ModTime: info.ModTime(),
				Size:    "(DIR)",
			})
		} else {
			files = append(files, File{
				Name:    entry.Name(),
				Path:    path.Join(unsafePath, entry.Name()),
				IsDir:   entry.IsDir(),
				ModTime: info.ModTime(),
				Size:    formatSize(info.Size()),
			})
		}
	}

	return append(directories, files...), nil
}

func Mkdir(unsafePath string, userID uint) error {
	return os.Mkdir(getSafePath(unsafePath, userID), os.ModePerm)
}

func GetFullPathAndName(unsafePath string, userID uint) (string, string, error) {
	// Returns safe path, file name and error
	safePath := getSafePath(unsafePath, userID)
	_, err := os.Stat(safePath)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(safePath, "/")
	return safePath, parts[len(parts)-1], nil
}

func Rm(unsafePath string, userID uint) error {
	// Removes only empty folders or files
	return os.Remove(getSafePath(unsafePath, userID))
}
