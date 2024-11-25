package utils

import (
	"errors"
	"fmt"
	"freb/config"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// req
const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Safari/605.1.15"
	domain    = "https://69shuba.cx"
	tocPage   = "https://69shuba.cx/book/%s.htm"
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

func DownloadTmp(filename string, handler func() *http.Request) (path string, err error) {
	if handler != nil {
		req := handler()
		paths := strings.Split(req.URL.Path, "/")
		name := paths[len(paths)-1]
		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			err = fmt.Errorf("下载封面错误:%w", err)
			return
		}
		defer resp.Body.Close()

		// 检查 HTTP 响应状态码
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("无法请求地址: %s", resp.Status)
			return
		}
		// 创建文件
		path = filepath.Join(config.Cfg.TmpDir, filename+filepath.Ext(name))
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

func TocUrl() string {
	return tocPage
}
