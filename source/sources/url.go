package sources

import (
	"errors"
	"fmt"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/source"
	"freb/utils"
	"freb/utils/stdout"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"os"
	"strings"
	"time"
)

type UrlSource struct {
}

func setChapterUrl(i int, title, url string, ef *formatter.EpubFormat) (index int) {
	index = i
	if index < len(ef.Book.Chapters) {
		if ef.Book.Chapters[index].IsVol {
			index++
		}
		title = utils.ChapterTitleWithoutNum(title)
		checkTitle := utils.ChapterTitleWithoutNum(ef.Book.Chapters[index].Title)
		if utils.SimilarStr(title, checkTitle) {
			ef.Book.Chapters[index].Url = url
			index++
		}
	}
	return
}

func getCatalog(ef *formatter.EpubFormat, doc *goquery.Document) (err error) {
	var toc, tocSlt string
	if !ef.Book.IsOld {
		toc, _ = doc.Find("a.btn:first-child").Attr("href")
		toc = utils.Domain() + toc
		tocSlt = "div#chapters ul li"
	} else {
		toc, _ = doc.Find("a.more-btn").Attr("href")
		tocSlt = "div#catalog ul li"
	}

	var isCatalog bool
	var chapterIndex int
	if ef.Book.Catalog.Url != "" {
		err = source.GetCatalogByUrl(ef)
		if err != nil {
			return
		}
		isCatalog = true
	}

	doc, err = utils.GetDomByDefault(toc)
	if err != nil {
		return err
	}
	var isReverse bool
	if sorting := doc.Find(".sorting a[style]").Text(); strings.Contains(sorting, "倒序") {
		isReverse = true
	}
	total := doc.Find(tocSlt).Length()
	if total == 0 {
		return errors.New("爬取错误: 章节数为 0")
	}
	total -= ef.Book.Jump
	if total <= 0 {
		return errors.New("跳过章节数[flag -j(jump)] 大于总章数")
	}
	utils.AscEach(doc.Find(tocSlt), func(i int, s *goquery.Selection) {
		if i < ef.Book.Jump {
			return
		}
		if i == 0 && ef.Book.Chapters == nil {
			ef.Book.Chapters = make([]models.Chapter, total)
		}

		url, _ := s.Find("a").Attr("href")
		url = utils.EmptyOrDomain(ef.Book.IsOld) + url
		if isCatalog {
			chapterIndex = setChapterUrl(chapterIndex, strings.TrimSpace(s.Find("a").Text()), url, ef)
			if i == total-1 && ef.Chapters[chapterIndex-1].Url == url {
				ef.Chapters = ef.Chapters[:chapterIndex]
			}
		} else {
			ef.Book.Chapters[i].Title = utils.PureTitle(s.Find("a").Text())
			ef.Book.Chapters[i].Url = url
		}
	}, isReverse)
	if isCatalog {
		var cdbErrChapter [2]string
		var errChapter string
		for i := range ef.Chapters {
			if ef.Chapters[i].Url == "" && !ef.Chapters[i].IsVol {
				errChapter = ef.Book.Chapters[i].Title
				if i > 0 {
					cdbErrChapter[0] = ef.Book.Chapters[i-1].Title
					cdbErrChapter[1] = ef.Book.Chapters[i-1].Url
				}
				break
			}
		}
		if cdbErrChapter[0] != "" {
			stdout.Warnf("可能匹配URL错误章节:\n%s\t%s\n", cdbErrChapter[0], cdbErrChapter[1])
		}
		if errChapter != "" {
			stdout.Warnf("空 URL 匹配起始章节: %s\n", errChapter)
			var isContinue string
			stdout.Fmt("以上章节合并失败,请确认是否继续生成 EPub (y/n):")
			_, err = fmt.Scanln(&isContinue)
			if err != nil {
				return err
			}
			if isContinue != "y" {
				return errors.New("终止生成 EPub")
			}
		}
	}
	return
}

func (u *UrlSource) GetBook(ef *formatter.EpubFormat) (err error) {
	start := time.Now()
	doc, err := utils.GetDomByDefault(utils.TocUrl(ef.Book.IsOld, ef.Book.Id))
	if err != nil {
		return err
	}
	ef.Book.Name = doc.Find("div.booknav2 h1 a").Text()

	if ef.Book.Author == "Unknown" {
		ef.Book.Author = doc.Find("div.booknav2 p a[href*='author']").Text()
	}
	var titleSlt, contentSlt string
	if !ef.Book.IsOld {
		ef.Book.Intro = doc.Find("div.navtxt p:first-child").Text()

		titleSlt = "div.chaptertitle h1"
		contentSlt = "div.content"
	} else {
		ef.Book.Intro = doc.Find("div.content").Text()

		titleSlt = "div.txtnav h1"
		contentSlt = "div.txtnav"
	}

	// chapter
	stdout.Fmtln("正在获取目录信息...")
	err = getCatalog(ef, doc)
	if err != nil {
		return err
	}
	stdout.Fmtfln("章节数: %d", len(ef.Book.Chapters))
	err = ef.InitBook()
	if err != nil {
		return err
	}
	// contents
	stdout.Fmtln("正在添加章节...")
	var volPath string
	for i, chapter := range ef.Book.Chapters {
		if chapter.Url == "" {
			continue
		}
		doc, err = utils.GetDomByDefault(chapter.Url)
		if err != nil {
			return
		}

		node := doc.Find(contentSlt).Contents().Not("h1,div")
		if doc.Find(titleSlt).Text() == "" {
			return errors.New("当前章节爬取错误")
		}

		var f func(int, *html.Node)
		f = func(index int, n *html.Node) {
			if n.Type == html.TextNode {
				raw := strings.TrimSpace(n.Data)
				if raw == "" || len([]rune(raw)) == 1 {
					return
				}
				// filter title in content
				if utils.SimilarStr(raw, ef.Book.Chapters[i].Title) && index <= 10 {
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
	stdout.Successfln("\n已生成书籍,使用时长: %s", totalTime)
	return
}
