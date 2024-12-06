package sources

import (
	"errors"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"os"
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
	var isReverse bool
	numStr, _ := doc.Find(tocSlt + ":last-child").Attr("data-num")
	if numStr == "1" {
		isReverse = true
	}
	total = doc.Find(tocSlt).Length()
	doc.Find(tocSlt).Each(func(i int, s *goquery.Selection) {
		index := i
		if isReverse {
			index = total - 1 - i
		}
		if i == 0 {
			book.Chapters = make([]models.Chapter, total)
		}
		book.Chapters[index].Url, _ = s.Find("a").Attr("href")
		book.Chapters[index].Url = utils.EmptyOrDomain(book.IsOld) + book.Chapters[index].Url
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
	var volPath string
	for i, chapter := range book.Chapters {
		if chapter.Url == "" {
			continue
		}
		doc, err = utils.GetDom(chapter.Url)
		if err != nil {
			return
		}

		node := doc.Find(contentSlt).Contents()
		book.Chapters[i].Title = strings.TrimSpace(doc.Find(titleSlt).Text())
		if book.Chapters[i].Title == "" {
			return errors.New("当前章节爬取错误")
		}
		book.Chapters[i].Title = utils.PureTitle(book.Chapters[i].Title)

		contentLen := len(node.Nodes)
		var f func(int, *html.Node)
		f = func(index int, n *html.Node) {
			if n.DataAtom == atom.Div || n.DataAtom == atom.H1 {
				return
			}
			if n.Type == html.TextNode {
				raw := strings.TrimSpace(n.Data)
				if raw == "" || len([]rune(raw)) == 1 {
					return
				}
				// filter title in content
				if strutil.Similarity(raw, book.Chapters[i].Title, metrics.NewJaro()) > 0.75 && index <= contentLen/3 {
					return
				}
				if strings.Contains(raw, "本章完") {
					return
				}
				book.Chapters[i].Content += ef.GenLine(raw)
			}
			if n.FirstChild != nil {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(index, c)
				}
			}
		}
		for index, n := range node.Nodes {
			f(index, n)
		}
		volPath, err = ef.GenBookContent(i, volPath)
		if err != nil {
			return
		}
		time.Sleep(time.Duration(config.Cfg.DelayTime) * time.Millisecond)
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
