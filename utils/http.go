package utils

import (
	"bytes"
	"errors"
	"fmt"
	"freb/utils/reg"
	"freb/utils/stdout"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
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
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

const (
	githubRaw = "https://ghp.ci/https://raw.githubusercontent.com/chcthink/freb/refs/heads/main/"
)

func LocalOrDownload(path, tmpDir, from string) (source string, err error) {
	if filePath, isExist := IsFileInExecDir(path); !isExist {
		if from == "" {
			from = githubRaw
		}
		downloadUrl := from + path
		stdout.Fmtfln("正在从远程仓库下载文件: %s", downloadUrl)
		source, err = DownloadTmp(tmpDir, path, func() *http.Request {
			return GetWithUserAgent(downloadUrl)
		})
		return
	} else {
		return filePath, err
	}
}

func GetWithUserAgent(url string) (req *http.Request) {
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	return
}

func GetDomByDefault(url string) (doc *html.Node, err error) {
	req := GetWithUserAgent(url)
	return TransDom2Doc(req)
}

func TransDom2Doc(req *http.Request) (doc *html.Node, err error) {
	var body []byte
	body, err = TransDom2Bytes(req)
	unescapedBody := html.UnescapeString(string(body))
	doc, err = htmlquery.Parse(strings.NewReader(unescapedBody))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func TransDom2JSON(req *http.Request) (rest gjson.Result, err error) {
	var body []byte
	body, err = TransDom2Bytes(req)
	return gjson.ParseBytes(body), nil
}

const reqErr = "请求失败: %s"

func TransDom2Bytes(req *http.Request) (body []byte, err error) {
	url := req.URL.String()
	if !reg.CheckUrl(url) {
		return nil, errors.New(stdout.ErrUrl)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf(reqErr, url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("爬取错误:%s %s", req.URL.String(), resp.Status)
	}
	defer resp.Body.Close()

	var buf []byte
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	e, _, _ := charset.DetermineEncoding(buf, resp.Header.Get("Content-Type"))
	body, err = io.ReadAll(transform.NewReader(bytes.NewReader(buf), e.NewDecoder()))
	if err != nil {
		return nil, err
	}
	return body, nil
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

func DivideThousandURL(str, id string) string {
	bookId, _ := strconv.Atoi(id)
	mid := strconv.FormatFloat(math.Floor(float64(bookId)/1000.0), 'f', 0, 64)
	return fmt.Sprintf(str, mid, id, id)
}
