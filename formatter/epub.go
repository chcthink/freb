package formatter

import (
	"bytes"
	"fmt"
	"freb/config"
	"freb/models"
	"freb/utils"
	"freb/utils/stdout"
	"github.com/go-shiori/go-epub"
	"path"
	"strconv"
	"strings"
)

// 章节
const (
	chapterFilePrefix = `chapter_`
	volFilePrefix     = `vol_`
	htmlP             = `<p>`
	htmlPEnd          = `</p>`
	percentSign       = "%"
)

type AssetsPath struct {
	CommonCss      string
	CoverCss       string
	FontCss        string
	InstructionCss string
	Font           string
	MetaInf        string
}

type Inner struct {
	volImage    string
	contentLogo string
	css         string
	volIndex    int
}

type EpubFormat struct {
	*epub.Epub
	*models.Book
	*Inner
	*AssetsPath
}

func (e *EpubFormat) InitBook() (err error) {
	e.Epub, err = epub.NewEpub(e.Book.Name)
	e.Inner = &Inner{}
	if err != nil {
		stdout.Err(err)
		return
	}

	// 初始化书籍信息
	stdout.Fmtf("初始化书籍信息:%s", e.Name)
	// 添加 css
	e.Inner.css, err = e.AddCSS(e.AssetsPath.CommonCss, "main.css")
	if err != nil {
		stdout.Err(err)
		return
	}
	_, _ = e.AddFont(e.AssetsPath.Font, "font.ttf")
	_, err = e.AddCSS(e.AssetsPath.FontCss, "fonts.css")
	if err != nil {
		stdout.Err(err)
		return
	}
	err = e.AddMetaINF(e.MetaInf)
	if err != nil {
		stdout.Err(err)
		return
	}
	e.SetLang(e.Book.Lang)
	// 添加标题
	e.SetTitle(e.Book.Name)
	// 添加封面
	if e.Cover != "" {
		var image, coverCss string

		image, err = e.AddImage(e.Cover, path.Base(e.Cover))
		if err != nil {
			err = fmt.Errorf("添加封面失败 %w", err)
			return
		}
		coverCss, err = e.AddCSS(e.CoverCss, "cover.css")
		if err != nil {
			stdout.Err(err)
			return
		}
		err = e.SetCover(image, coverCss)
		if err != nil {
			stdout.Err(err)
			return
		}
	}
	// 添加作者
	e.SetAuthor(e.Book.Author)
	// 添加制作说明
	if e.Book.IsDesc {
		stdout.Fmt("正在添加制作说明...")
		var insPageCss string
		insPageCss, err = e.AddCSS(e.InstructionCss, "instruction.css")
		if err != nil {
			stdout.Err(err)
			return
		}
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Instruction.Dom, e.Book.Name, e.Book.Author),
			config.Cfg.Instruction.Title, "instruction.xhtml", insPageCss)
		if err != nil {
			stdout.Err(err)
			return
		}
	}
	// 内容简介
	if e.Book.Intro != "" {
		stdout.Fmt("正在添加内容简介...")
		var logo string
		logo, err = e.AddImage(e.IntroImg, path.Base(e.IntroImg))
		if err != nil {
			return
		}
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Desc.Dom, logo, e.Book.Intro), config.Cfg.Desc.Title, "desc.xhtml", e.Inner.css)
		if err != nil {
			return
		}
	}
	if e.Book.Vol != "" {
		e.volImage, err = e.AddImage(e.Vol, path.Base(e.Vol))
		if err != nil {
			stdout.Err(err)
			return
		}
	}
	if e.Book.ContentImg != "" {
		e.contentLogo, err = e.AddImage(e.Book.ContentImg, path.Base(e.ContentImg))
		if err != nil {
			stdout.Err(err)
			return
		}
	}
	return
}

func cleanHTML(str string) string {
	str = utils.PureEscapeHtml(str)
	return utils.ReplaceC0Control(str)
}

func genLine(str string) string {
	return htmlP + strings.ReplaceAll(str, percentSign, "%%") + htmlPEnd
}

func (e *EpubFormat) GenLine(str string) string {
	str = cleanHTML(str)
	return genLine(str)
}

func (e *EpubFormat) GenLine2Buffer(str string, buf *bytes.Buffer) {
	str = cleanHTML(str)
	buf.WriteString(genLine(str))
}

func (e *EpubFormat) GenBookContent(index int, vol string) (volPath string, err error) {
	title := e.Book.Chapters[index].Title
	fmt.Printf("\r[%d/%d]\033[K%s", index+1, len(e.Chapters), title)
	if volNum, volTitle, isVol := utils.VolByDefaultReg(title); isVol {
		e.volIndex += 1
		volPath, err = e.AddSection(fmt.Sprintf(config.Cfg.Vol, e.volImage, volNum, volTitle),
			volNum+" "+volTitle, volFilePrefix+strconv.Itoa(e.volIndex)+".xhtml", e.Inner.css)
		if err != nil {
			stdout.Err(err)
			return
		}
		return
	}

	num, name, subNum := utils.ChapterTitleByDefaultReg(title)
	if vol == "" {
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Chapter+e.Book.Chapters[index].Content,
			e.contentLogo, num, name, subNum), strings.Join([]string{num, name, subNum}, " "),
			chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.Inner.css)
		if err != nil {
			stdout.Err(err)
			return
		}
	} else {
		_, err = e.AddSubSection(vol, fmt.Sprintf(config.Cfg.Chapter+e.Book.Chapters[index].Content,
			e.contentLogo, num, name, subNum), strings.Join([]string{num, name, subNum}, " "),
			chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.Inner.css)
		if err != nil {
			stdout.Err(err)
			return
		}
		volPath = vol
	}

	return
}

func (e *EpubFormat) Build() error {
	if e.Out != "" {
		return e.Epub.Write(e.Out)
	}
	return e.Epub.Write(e.Name + "-" + e.Book.Author + ".epub")
}
