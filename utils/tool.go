package utils

import (
	"fmt"
	"freb/models"
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

func IsDev() bool {
	return strings.Contains(models.Version, "dev")
}

func IsFileInExecDir(path string) (filePath string, isExist bool) {
	if IsDev() {
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
