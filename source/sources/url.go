package sources

import (
	"errors"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/utils"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type UrlSource struct {
}

func (u *UrlSource) GetBook(book *models.Book) (err error) {
	start := time.Now()
	doc, err := utils.GetDom(utils.TocUrl(book.IsOld, book.Id))
	if err != nil {
		return err
	}
	book.Name = doc.Find("div.booknav2 h1 a").Text()

	if book.Author == "Unknown" {
		book.Author = doc.Find("div.booknav2 p a[href*='author']").Text()
	}

	var toc, tocSlt, titleSlt, contentSlt string
	if !book.IsOld {
		book.Intro = doc.Find("div.navtxt p:first-child").Text()

		toc, _ = doc.Find("a.btn:first-child").Attr("href")
		toc = utils.Domain() + toc

		titleSlt = "div.chaptertitle h1"
		contentSlt = "div.content"
		tocSlt = "div#chapters ul li"
	} else {
		book.Intro = doc.Find("div.content").Text()
		toc, _ = doc.Find("a.more-btn").Attr("href")

		titleSlt = "div.txtnav h1"
		contentSlt = "div.txtnav"
		tocSlt = "div#catalog ul li"
	}

	// chapter
	utils.Fmt("正在获取目录信息...")
	doc, err = utils.GetDom(toc)
	if err != nil {
		return err
	}

	var total int
	numStr, _ := doc.Find(tocSlt + ":last-child").Attr("data-num")
	if numStr == "1" {
		numStr, _ = doc.Find(tocSlt + ":first-child").Attr("data-num")
	}
	total, _ = strconv.Atoi(numStr)
	doc.Find(tocSlt).Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			book.Chapters = make([]models.Chapter, total)
		}
		book.Chapters[i].Url, _ = s.Find("a").Attr("href")
		book.Chapters[i].Url = utils.EmptyOrDomain(book.IsOld) + book.Chapters[i].Url
	})
	if len(book.Chapters) == 0 {
		return errors.New("爬取错误: 章节数为 0")
	}
	utils.Fmtf("章节数: %d", len(book.Chapters))
	// confirm format
	var ef formatter.EpubFormat
	ef.Book = book
	err = ef.InitBook()
	if err != nil {
		return err
	}
	// contents
	utils.Fmt("正在添加章节...")
	// return
	var volPath string
	for i, chapter := range book.Chapters {
		doc, err = utils.GetDom(chapter.Url)
		if err != nil {
			return
		}

		node := doc.Find(contentSlt).Contents()
		book.Chapters[i].Title = doc.Find(titleSlt).Text()
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
				if utils.EqAllWithoutSpace(book.Chapters[i].Title, n.Data) {
					return
				}
				if strings.Contains(n.Data, "本章完") {
					return
				}
				book.Chapters[i].Content = ef.GenLine(raw)
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
		volPath, err = ef.GenBookContent(i, volPath)
		if err != nil {
			return
		}
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	err = ef.Build()
	if err != nil {
		return
	}

	_ = os.RemoveAll(config.Cfg.TmpDir)
	totalTime := time.Since(start).Truncate(time.Second).String()
	utils.Successf("\n已生成书籍,使用时长: %s", totalTime)
	return
}
