package utils

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"net/http"
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
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)
	filePath = filepath.Join(execDir, path)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return
	}
	return filePath, true
}

func PureEscapeHtml(str string) string {
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	return strings.ReplaceAll(str, "&", "&amp;")
}

const (
	defaultImgDir = "assets/images/"
)

func SetImage(from, dir, filename string, handler func() *http.Request) (path string, err error) {
	if IsImgFile(from) {
		if IsFileInWorkDir(from) {
			path = from
			return
		}
	}
	if CheckUrl(from) {
		path, err = DownloadTmp(dir, filename, func() *http.Request {
			return NewGetWithUserAgent(from)
		})
		if path != "" {
			return
		}
	}
	if handler != nil {
		path, err = DownloadTmp(dir, filename, handler)
		return
	}
	if filePath, isExist := IsFileInExecDir(defaultImgDir + filename); isExist {
		path = filePath
		return
	}

	path, err = DownloadTmp(dir, filename, func() *http.Request {
		return NewGetWithUserAgent(githubRaw + defaultImgDir + filename)
	})
	return
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
	if index := AuthorIndex(str); index > 0 {
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
