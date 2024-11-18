package source

import (
	"bytes"
	"freb/formatter"
	"freb/formatter/formats"
	"freb/models"
	"freb/utils"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strconv"
	"strings"
)

type Source interface {
	GetBook(*models.Book) error
}

type UrlSource struct {
}

func (u *UrlSource) GetBook(book *models.Book) (err error) {
	doc, err := utils.GetDom(book.Url)
	if err != nil {
		return err
	}
	book.Name = doc.Find("div.booknav2 h1 a").Text()
	//cover, _ := doc.Find("div.bookimg2 img").Attr("src")
	// todo 下载文件
	//if !strings.Contains(cover, "nc.jpg") {
	//	s := strings.Split(book.Cover, "/")
	//	path := utils.DEFAULT_IMAGE_PATH + s[len(s)-1]
	//	err := utils.DownloadFile(utils.GetDomainFromUrl(book.Url)+book.Cover, path)
	//	if err != nil {
	//		return err
	//	}
	//	book.Cover = path
	//}
	book.Author = doc.Find("div.booknav2 p a[href*='author']").Text()
	book.Intro = doc.Find("div.content").Text()

	// chapter
	utils.Fmt("正在获取目录信息...")
	toc, _ := doc.Find("a.more-btn").Attr("href")
	doc, err = utils.GetDom(toc)
	if err != nil {
		return err
	}

	var total int
	doc.Find("div#catalog ul li").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			numStr, _ := s.Attr("data-num")
			num, _ := strconv.Atoi(numStr)
			total = num
			book.ChapterUrls = make([]string, num)
		}
		book.ChapterUrls[total-i-1], _ = s.Find("a").Attr("href")
	})
	utils.Fmtf("章节数: %d", len(book.ChapterUrls))
	// confirm format
	var novel formatter.Formatter
	switch book.Format {
	case "epub":
		novel = &formats.EpubFormat{Book: book}
		err = novel.InitBook()
		if err != nil {
			return err
		}
	}
	// contents
	utils.Fmt("正在添加章节...")
	for i, url := range book.ChapterUrls {
		doc, err = utils.GetDom(url)
		if err != nil {
			return
		}
		node := doc.Find("div.txtnav").Contents()
		var title bytes.Buffer
		title.WriteString(doc.Find("div.txtnav h1").Text())
		var buf bytes.Buffer
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.DataAtom == atom.Div || n.DataAtom == atom.H1 {
				return
			}
			if n.Type == html.TextNode {
				raw := strings.TrimSpace(n.Data)
				if raw == "" {
					return
				}
				if strings.Contains(n.Data, "本章完") || utils.IsTitle(n.Data) {
					return
				}
				if book.Format == "epub" {
					novel.GenContentPrefix(&buf, raw)
				} else {
					buf.WriteString(raw)
				}
			}
			if n.FirstChild != nil {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		}
		for _, n := range node.Nodes {
			f(n)
		}
		err = novel.GenBookContent(i+1, title, buf)
		if err != nil {
			return
		}
	}
	err = novel.Build()
	if err != nil {
		return
	}
	utils.Success("已生成书籍")
	return
}
