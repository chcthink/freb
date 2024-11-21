package utils

import (
	"errors"
	"fmt"
	"freb/config"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

// req
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Safari/605.1.15"

func NewReq(url string) (req *http.Request) {
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	// can't confirm this cookie is work in connecting check
	// req.AddCookie(&http.Cookie{Name: "shuba", Value: "11072-13848-21591-4277"})
	return
}

// GetDom 获取 HTML DOM
func GetDom(url string) (doc *goquery.Document, err error) {
	if !CheckUrl(url) {
		return nil, errors.New(ErrUrl)
	}
	req := NewReq(url)
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

func DownloadCover(url string) (path string, err error) {

	// 创建 HTTP 请求
	domain := GetDomainFromUrl(url)
	bookId := BookId(url)
	bookIdStr := strconv.Itoa(bookId)
	mid := strconv.FormatFloat(math.Floor(float64(bookId)/1000.0), 'f', 0, 64)
	url = domain + "/fengmian/" + mid + "/" + bookIdStr + "/" + bookIdStr + "s.jpg"
	req := NewReq(url)
	req.Header.Set("Referer", "https://69shuba.cx/modules/article/search.php")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("下载封面错误:%w", err)
		return
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println(url)
		err = fmt.Errorf("无法请求封面地址: %s", resp.Status)
		return
	}
	// 创建文件
	path = config.Cfg.TmpDir + "/cover.jpg"
	out, err := os.Create(path)
	if err != nil {
		return
	}
	defer out.Close()

	// 将响应体写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return
}
