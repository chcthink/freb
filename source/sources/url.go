package sources

import (
	"errors"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/utils"
	"freb/utils/stdout"
	"github.com/PuerkitoBio/goquery"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"os"
	"strconv"
	"strings"
	"time"
)

type UrlSource struct {
}

func (u *UrlSource) GetBook(ef formatter.EpubFormat) (err error) {
	start := time.Now()
	doc, err := utils.GetDom(utils.TocUrl(ef.Book.IsOld, ef.Book.Id))
	if err != nil {
		return err
	}
	ef.Book.Name = doc.Find("div.booknav2 h1 a").Text()

	if ef.Book.Author == "Unknown" {
		ef.Book.Author = doc.Find("div.booknav2 p a[href*='author']").Text()
	}

	var toc, tocSlt, titleSlt, contentSlt string
	if !ef.Book.IsOld {
		ef.Book.Intro = doc.Find("div.navtxt p:first-child").Text()

		toc, _ = doc.Find("a.btn:first-child").Attr("href")
		toc = utils.Domain() + toc

		titleSlt = "div.chaptertitle h1"
		contentSlt = "div.content"
		tocSlt = "div#chapters ul li"
	} else {
		ef.Book.Intro = doc.Find("div.content").Text()
		toc, _ = doc.Find("a.more-btn").Attr("href")

		titleSlt = "div.txtnav h1"
		contentSlt = "div.txtnav"
		tocSlt = "div#catalog ul li"
	}

	// chapter
	stdout.Fmt("正在获取目录信息...")
	doc, err = utils.GetDom(toc)
	if err != nil {
		return err
	}
	var isReverse bool
	firstChp, _ := doc.Find("div#catalog ul li" + ":nth-child(1)").Attr("data-num")
	secChp, _ := doc.Find("div#catalog ul li" + ":nth-child(2)").Attr("data-num")
	first, _ := strconv.Atoi(firstChp)
	sec, _ := strconv.Atoi(secChp)
	if first-sec > 0 {
		isReverse = true
	}
	var total int
	total = doc.Find(tocSlt).Length()
	if total == 0 {
		return errors.New("爬取错误: 章节数为 0")
	}
	total -= ef.Book.Jump
	if total <= 0 {
		return errors.New("跳过章节数[flag -j(jump)] 大于总章数")
	}
	doc.Find(tocSlt).Each(func(i int, s *goquery.Selection) {
		index := i
		if isReverse {
			index = total - 1 - i
			if i >= total {
				return
			}
		} else {
			if i < ef.Book.Jump {
				return
			}
		}
		if i == 0 {
			ef.Book.Chapters = make([]models.Chapter, total)
		}
		ef.Book.Chapters[index].Title = utils.PureTitle(s.Find("a").Text())

		ef.Book.Chapters[index].Url, _ = s.Find("a").Attr("href")
		ef.Book.Chapters[index].Url = utils.EmptyOrDomain(ef.Book.IsOld) + ef.Book.Chapters[index].Url
	})
	stdout.Fmtf("章节数: %d", len(ef.Book.Chapters))
	err = ef.InitBook()
	if err != nil {
		return err
	}
	// contents
	stdout.Fmt("正在添加章节...")
	var volPath string
	for i, chapter := range ef.Book.Chapters {
		if chapter.Url == "" {
			continue
		}
		doc, err = utils.GetDom(chapter.Url)
		if err != nil {
			return
		}

		node := doc.Find(contentSlt).Contents()
		if doc.Find(titleSlt).Text() == "" {
			return errors.New("当前章节爬取错误")
		}

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
				if strutil.Similarity(raw, ef.Book.Chapters[i].Title, metrics.NewJaro()) > 0.75 && index <= contentLen/3 {
					return
				}
				if strings.Contains(raw, "本章完") {
					return
				}
				ef.Book.Chapters[i].Content += ef.GenLine(raw)
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
		time.Sleep(time.Duration(ef.Delay) * time.Millisecond)
	}
	err = ef.Build()
	if err != nil {
		return
	}

	_ = os.RemoveAll(config.Cfg.TmpDir)
	totalTime := time.Since(start).Truncate(time.Second).String()
	stdout.Successf("\n已生成书籍,使用时长: %s", totalTime)
	return
}
