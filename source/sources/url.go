package sources

import (
	"errors"
	"fmt"
	"freb/formatter"
	"freb/models"
	"freb/source"
	"freb/utils"
	"freb/utils/htmlx"
	"freb/utils/reg"
	"freb/utils/stdout"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
	"time"
)

type UrlSource struct {
}

func setChapterUrl(i int, title, url string, ef *formatter.EpubFormat) (index int) {
	index = i
	if index < len(ef.Sections) {
		if ef.Sections[index].IsVol {
			index++
		}
		title = reg.ChapterTitleWithoutNum(title)
		checkTitle := reg.ChapterTitleWithoutNum(ef.Sections[index].Title)
		if utils.SimilarStr(title, checkTitle) {
			ef.Sections[index].Url = url
			index++
		}
	}
	return
}

func getCatalog(ef *formatter.EpubFormat, doc *html.Node, catch *models.BookCatch) (err error) {
	var isCatalog bool
	var chapterIndex int
	if ef.BookConf.Catalog != "" {
		err = source.GetCatalogFromUrl(ef)
		if err != nil {
			return
		}
		isCatalog = true
	}

	tocUrl := htmlx.XPathFindStr(doc, catch.Toc)
	if !reg.CheckUrl(tocUrl) {
		tocUrl = strings.Join([]string{ef.BookConf.Url, tocUrl}, "")
	}

	doc, err = utils.GetDomByDefault(tocUrl)
	if err != nil {
		return err
	}
	var isReverse bool
	if sorting := htmlx.XPathFindStr(doc, catch.Sort); strings.Contains(sorting, "倒序") {
		isReverse = true
	}
	chapters := htmlquery.Find(doc, catch.Chapter.Element)
	total := len(chapters)
	if total == 0 {
		return errors.New("爬取错误: 章节数为 0")
	}
	total -= ef.BookConf.Jump
	if total <= 0 {
		return errors.New("跳过章节数[flag -j(jump)] 大于总章数")
	}
	htmlx.XpPathAscEach(chapters, func(i int, s *html.Node) {
		if i < ef.BookConf.Jump {
			return
		}
		if i == 0 && ef.Sections == nil {
			ef.Sections = make([]models.Section, total)
		}

		url := htmlx.XPathFindStr(s, catch.Chapter.Url)
		if !reg.CheckUrl(url) {
			tocUrl = strings.Join([]string{ef.BookConf.Url, url}, "")
		}
		title := htmlx.XPathFindStr(s, catch.Chapter.Title)

		if isCatalog {
			chapterIndex = setChapterUrl(chapterIndex, strings.TrimSpace(title), url, ef)
			if i == total-1 && ef.Sections[chapterIndex-1].Url == url {
				ef.Sections = ef.Sections[:chapterIndex]
			}
		} else {
			ef.Sections[i].Title = reg.PureTitle(title)
			ef.Sections[i].Url = url
		}

		// filter by config
		ef.Sections[i].Title = reg.RemoveTitleFromCfg(ef.Sections[i].Title)
	}, isReverse)
	if isCatalog {
		var cdbErrChapter [2]string
		var errChapter string
		for i := range ef.Sections {
			if ef.Sections[i].Url == "" && !ef.Sections[i].IsVol {
				errChapter = ef.Sections[i].Title
				if i > 0 {
					cdbErrChapter[0] = ef.Sections[i-1].Title
					cdbErrChapter[1] = ef.Sections[i-1].Url
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

func (u *UrlSource) GetBook(ef *formatter.EpubFormat, catch *models.BookCatch) (err error) {
	start := time.Now()
	doc, err := utils.GetDomByDefault(ef.BookConf.Url)
	if err != nil {
		return err
	}

	// get book basic info
	ef.Name = htmlx.XPathFindStr(doc, catch.Name.Selector)
	if ef.Author == "Unknown" {
		ef.Author = htmlx.XPathFindStr(doc, catch.Author.Selector)
	}
	ef.Intro = htmlx.XPathFindStr(doc, catch.Intro.Selector)
	ef.Intro, err = reg.Filters(catch.Intro.Filter, ef.Intro)
	if err != nil {
		return
	}

	// chapter
	stdout.Fmtln("正在获取目录信息...")
	err = getCatalog(ef, doc, catch)
	if err != nil {
		return err
	}
	stdout.Fmtfln("章节数: %d", len(ef.Sections))
	err = ef.InitBook()
	if err != nil {
		return err
	}
	// contents
	stdout.Fmtln("正在添加章节...")
	var volPath string
	for i, chapter := range ef.Sections {
		if chapter.Url == "" {
			continue
		}
		doc, err = utils.GetDomByDefault(chapter.Url)
		if err != nil {
			return
		}
		node := htmlquery.Find(doc, catch.Content.Selector)

		if htmlx.XPathFindStr(doc, catch.Title.Selector) == "" {
			return fmt.Errorf("当前章节爬取错误: %s %s", chapter.Title, chapter.Url)
		}

		var f func(int, *html.Node)
		f = func(index int, n *html.Node) {
			if n.Type == html.TextNode {
				raw := strings.TrimSpace(n.Data)
				if raw == "" || len([]rune(raw)) == 1 {
					return
				}
				// filter title in content
				if utils.SimilarStr(raw, ef.Sections[i].Title) && index <= 10 {
					return
				}
				if strings.Contains(raw, "本章完") {
					return
				}
				raw = reg.RemoveContentFromCfg(raw)
				if raw == "" {
					return
				}

				ef.Sections[i].Content += ef.GenLine(raw)
			}
			if n.FirstChild != nil {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(index, c)
				}
			}
		}
		for index, n := range node {
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

	totalTime := time.Since(start).Truncate(time.Second).String()
	stdout.Successfln("\n已生成书籍,使用时长: %s", totalTime)
	return
}
