package utils

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	bookSupports  = []string{".txt", ".epub"}
	imageSupports = []string{".jpg", ".jpeg", ".png", ".svg", ".webp"}
)

func CheckFileType(filename string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, fileType := range exts {
		if ext == fileType {
			return true
		}
	}
	return false
}

func IsImgFile(filename string) bool {
	return CheckFileType(filename, imageSupports)
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ReplaceTitle(str, title string) (ret string) {
	ret = strings.ReplaceAll(str, title, "")
	return strings.ReplaceAll(ret, strings.ReplaceAll(title, " ", ""), "")
}
