package utils

import (
	"fmt"
	"freb/models"
	"freb/utils/reg"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"os"
	"path/filepath"
	"strings"
)

var (
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

func IsFileInWorkDir(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsFileInExecDir(path string) (filePath string, isExist bool) {
	if strings.Contains(models.Version, "dev") {
		filePath, _ = findProjectRoot()
		filePath = filepath.Join(filePath, path)
	} else {
		execPath, _ := os.Executable()
		execDir := filepath.Dir(execPath)
		filePath = filepath.Join(execDir, path)
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return
	}
	return filePath, true
}

const (
	bookNameMarkPre = '《'
	bookNameMarkSuf = '》'
	hyphen          = '-'
)

func GetBookInfo(str string) (name, author string) {
	if strings.ContainsRune(str, bookNameMarkPre) && strings.ContainsRune(str, bookNameMarkSuf) {
		start := strings.IndexRune(str, bookNameMarkPre)
		end := strings.IndexRune(str, bookNameMarkSuf)
		if start >= end {
			return
		}
		name = strings.Replace(str[start:end], string(bookNameMarkPre), "", 1)
	}
	if index := reg.AuthorIndex(str); index > 0 {
		author = str[index:]
	}
	if strings.ContainsRune(str, hyphen) {
		name = str[:strings.IndexRune(str, hyphen)]
		author = str[strings.IndexRune(str, hyphen):]
	}
	return
}

func SimilarStr(str1, str2 string) bool {
	metricsList := []strutil.StringMetric{
		metrics.NewJaro(),
		metrics.NewSmithWatermanGotoh(),
		metrics.NewJaroWinkler(),
	}

	// 遍历所有的度量方法，检查是否有一个相似度大于阈值
	for _, metric := range metricsList {
		if strutil.Similarity(str1, str2, metric) > 0.7 {
			return true
		}
	}
	return false
}

func findProjectRoot() (string, error) {
	work, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(work, "go.mod")); err == nil {
			return work, nil
		}
		parent := filepath.Dir(work)
		if parent == work {
			return "", fmt.Errorf("未找到 go.mod 文件")
		}
		work = parent
	}
}
