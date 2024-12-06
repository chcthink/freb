package utils

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// req
const (
	userAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	oldDomain    = "https://69shuba.cx"
	domain       = "https://www.69yuedu.net"
	oldSearchUrl = "https://69shuba.cx/modules/article/search.php"
	searchUrl    = "https://www.69yuedu.net/modules/article/search.php"
	coverUrl     = "https://www.69yuedu.net/files/article/image/%s/cover.jpg"
	oldToc       = "https://69shuba.cx/book/%s.htm"
	tocPage      = "https://www.69yuedu.net/article/%s.html"
)

func NewGet(url string) (req *http.Request) {
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	return
}

// GetDom 获取 HTML DOM
func GetDom(url string) (doc *goquery.Document, err error) {
	if !CheckUrl(url) {
		return nil, errors.New(ErrUrl)
	}
	req := NewGet(url)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err = goquery.NewDocumentFromReader(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func DownloadTmp(dir, filename string, handler func() *http.Request) (path string, err error) {
	if handler != nil {
		req := handler()
		if req == nil {
			return
		}
		filename = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		paths := strings.Split(req.URL.Path, "/")
		name := paths[len(paths)-1]
		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			err = fmt.Errorf("下载错误:%w", err)
			return
		}
		defer resp.Body.Close()
		// 检查 HTTP 响应状态码
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("无法请求地址: %s", resp.Status)
			return
		}
		// 创建文件
		path = filepath.Join(dir, filename+filepath.Ext(name))
		var out *os.File
		out, err = os.Create(path)
		if err != nil {
			return
		}
		defer out.Close()

		// 将响应体写入文件
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return
		}
	}
	return
}

func Domain() string {
	return domain
}

func SearchUrl(isOld bool) string {
	if isOld {
		return oldSearchUrl
	}
	return searchUrl
}

func TocUrl(isOld bool, id string) string {
	if isOld {
		return fmt.Sprintf(oldToc, id)
	}
	return fmt.Sprintf(tocPage, id)
}

func CoverUrl(isOld bool, id string) string {
	if isOld {
		bookId, _ := strconv.Atoi(id)
		mid := strconv.FormatFloat(math.Floor(float64(bookId)/1000.0), 'f', 0, 64)
		return strings.Join([]string{domain, "fengmian", mid, id, id + "s.jpg"}, "/")
	}
	return fmt.Sprintf(coverUrl, id)
}

const (
	githubRaw = "https://ghp.ci/https://raw.githubusercontent.com/chcthink/freb/refs/heads/main/"
)

func LocalOrUrl(path string) string {
	if !IsFileExist(path) {
		return githubRaw + path
	}
	return path
}

func LocalOrDownload(path, tmpDir string) (source string, err error) {
	if !IsFileExist(path) {
		source, err = DownloadTmp(tmpDir, path, func() *http.Request {
			return NewGet(githubRaw + path)
		})
		return
	}

	return path, err
}

func EmptyOrDomain(isOld bool) string {
	if isOld {
		return ""
	}
	return domain
}
