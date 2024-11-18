package utils

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// reg
const (
	CHECK_URL  = "^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]"
	URL_DOMAIN = "(https?://)[^/]+"
)

const (
	DEFAULT_IMAGE_PATH = "assets/images/"
)

func CheckUrl(url string) bool {
	if len(url) == 0 {
		return false
	}
	re, _ := regexp.Compile(CHECK_URL)
	return re.MatchString(url)
}

func GetDomainFromUrl(url string) string {
	reg := regexp.MustCompile(URL_DOMAIN)
	return reg.FindString(url)
}

func DownloadFile(url string, filepath string) error {
	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将响应体写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

const (
	ChapterNumReg    = "第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷部]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^\\d+、|^引子$|^楔子$|^章节目录|^章节|^序章"
	ChapterSubNumReg = "[(（][0-9一二三四五六七八九十零〇百千两 ][)）]"
	volReg           = "^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]"
)

func ChapterTitleByDefaultReg(str string) (num, title, subNum string) {
	numReg := regexp.MustCompile(ChapterNumReg)
	SubNumReg := regexp.MustCompile(ChapterSubNumReg)
	num = numReg.FindString(str)
	subNum = SubNumReg.FindString(str)
	title = str
	if num != "" {
		title = strings.Split(title, num)[1]
	}
	if subNum != "" {
		title = strings.Split(title, subNum)[0]
	}
	if title != "" {
		title = strings.TrimSpace(title)
	}
	return
}

func VolByDefaultReg(str string) (num, title string, isVol bool) {
	reg := regexp.MustCompile(volReg)
	num = reg.FindString(str)
	if num != "" {
		isVol = true
		title = strings.TrimSpace(strings.Split(str, num)[1])
	}
	return
}

func IsTitle(str string) bool {
	reg := regexp.MustCompile(ChapterNumReg)
	return reg.MatchString(str)
}

// req
const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Safari/605.1.15"

// GetDom 获取 HTML DOM
func GetDom(url string) (doc *goquery.Document, err error) {
	if !CheckUrl(url) {
		return nil, errors.New(ERR_URL)
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	// can't confirm this cookie is work in connecting check
	req.AddCookie(&http.Cookie{Name: "shuba", Value: "11072-13848-21591-4277"})
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
