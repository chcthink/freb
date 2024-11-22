package source

import (
	"bytes"
	"freb/config"
	"freb/formatter"
	"freb/formatter/formats"
	"freb/models"
	"freb/utils"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"os"
	"strconv"
	"strings"
	"time"
)

type UrlSource struct {
}

func (u *UrlSource) GetBook(book *models.Book) (err error) {
	start := time.Now()
	doc, err := utils.GetDom(book.Url)
	if err != nil {
		return err
	}
	book.Name = doc.Find("div.booknav2 h1 a").Text()

	if book.Cover == "" || !utils.IsFileExist(book.Cover) {
		book.Cover, err = utils.DownloadCover(book.Url)
		if err != nil {
			return err
		}
	}

	if book.Author == "Unknown" {
		book.Author = doc.Find("div.booknav2 p a[href*='author']").Text()
	}
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
			book.Chapters = make([]models.Chapter, num)
		}
		book.Chapters[i].Title = &bytes.Buffer{}
		book.Chapters[i].Content = &bytes.Buffer{}
		book.Chapters[total-i-1].Url, _ = s.Find("a").Attr("href")
	})
	utils.Fmtf("章节数: %d", len(book.Chapters))
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
	for i, chapter := range book.Chapters {
		doc, err = utils.GetDom(chapter.Url)
		if err != nil {
			return
		}

		node := doc.Find("div.txtnav").Contents()
		book.Chapters[i].Title.WriteString(doc.Find("div.txtnav h1").Text())
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
					novel.GenContentPrefix(i, raw)
				} else {
					book.Chapters[i].Content.WriteString(raw)
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

		err = novel.GenBookContent(i)
		if err != nil {
			return
		}
	}
	err = novel.Build()
	if err != nil {
		return
	}

	_ = os.RemoveAll(config.Cfg.TmpDir)
	totalTime := time.Since(start).Truncate(time.Second).String()
	utils.Successf("\n已生成书籍,使用时长: %s", totalTime)
	return
}
