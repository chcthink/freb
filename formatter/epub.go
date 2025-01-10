package formatter

import (
	"bytes"
	"fmt"
	"freb/config"
	"freb/models"
	"freb/utils/reg"
	"freb/utils/stdout"
	"github.com/go-shiori/go-epub"
	"html"
	"path/filepath"
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

type EpubFormat struct {
	Name   string
	Author string
	Lang   string
	Intro  string
	Out    string
	*epub.Epub
	*models.BookConf
	*models.Images
	*models.Assets
	VolIndex int
	Sections []models.Section
}

func (e *EpubFormat) Init() {
	e.BookConf = &models.BookConf{}
	e.Images = &models.Images{
		Cover:       "cover.jpg",
		ContentLogo: "content_logo.jpg",
		IntroImg:    "intro_logo.jpg",
		VolImg:      "vol.jpg",
	}
	e.Assets = &models.Assets{}
}

func (e *EpubFormat) InitEPub() (err error) {
	e.Epub, err = epub.NewEpub(e.Name)
	if err != nil {
		stdout.Errln(err)
		return
	}

	// 初始化书籍信息
	stdout.Fmtfln("初始化书籍信息:%s", e.Name)
	// 添加 MainCss
	e.Assets.MainCss, err = e.AddCSS(e.Assets.MainCss, "main.css")
	if err != nil {
		stdout.Errln(err)
		return
	}
	_, _ = e.AddFont(e.Assets.Font, "font.ttf")
	_, err = e.AddCSS(e.Assets.FontCss, "fonts.css")
	if err != nil {
		stdout.Errln(err)
		return
	}
	err = e.AddMetaINF(e.MetaInf)
	if err != nil {
		stdout.Errln(err)
		return
	}
	e.SetLang(e.Lang)
	// 添加标题
	e.SetTitle(e.Name)
	// 添加封面
	if e.Cover != "" {
		var image, coverCss string

		image, err = e.AddImage(e.Cover, filepath.Base(e.Cover))
		if err != nil {
			err = fmt.Errorf("添加封面失败 %w", err)
			return
		}
		coverCss, err = e.AddCSS(e.CoverCss, "cover.css")
		if err != nil {
			stdout.Errln(err)
			return
		}
		err = e.SetCover(image, coverCss)
		if err != nil {
			stdout.Errln(err)
			return
		}
	}
	// 添加作者
	e.SetAuthor(e.Author)
	// 添加制作说明
	if e.BookConf.IsDesc {
		stdout.Fmtln("正在添加制作说明...")
		var insPageCss string
		insPageCss, err = e.AddCSS(e.InstructionCss, "instruction.css")
		if err != nil {
			stdout.Errln(err)
			return
		}
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Instruction.Dom, e.Name, e.Author),
			config.Cfg.Instruction.Title, "instruction.xhtml", insPageCss)
		if err != nil {
			stdout.Errln(err)
			return
		}
	}
	// 内容简介
	if e.Intro != "" {
		stdout.Fmtln("正在添加内容简介...")
		var logo string
		logo, err = e.AddImage(e.IntroImg, filepath.Base(e.IntroImg))
		if err != nil {
			return
		}
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Desc.Dom, logo, e.Intro), config.Cfg.Desc.Title, "desc.xhtml", e.Assets.MainCss)
		if err != nil {
			return
		}
	}
	if e.VolImg != "" {
		e.ColImg, err = e.AddImage(e.VolImg, filepath.Base(e.VolImg))
		if err != nil {
			stdout.Errln(err)
			return
		}
	}
	if e.ContentLogo != "" {
		e.ContentLogo, err = e.AddImage(e.ContentLogo, filepath.Base(e.ContentLogo))
		if err != nil {
			stdout.Errln(err)
			return
		}
	}
	return
}

func cleanHTML(str string) string {
	str = html.UnescapeString(str)
	return reg.ReplaceC0Control(str)
}

func genLine(str string) string {
	return htmlP + strings.ReplaceAll(str, percentSign, "%%") + htmlPEnd
}

func (e *EpubFormat) GenLine(str string) string {
	// str = cleanHTML(str)
	return genLine(str)
}

func (e *EpubFormat) GenLine2Buffer(str string, buf *bytes.Buffer) {
	str = cleanHTML(str)
	buf.WriteString(genLine(str))
}

func (e *EpubFormat) GenBookContent(index int, vol string) (volPath string, err error) {
	title := e.Sections[index].Title
	fmt.Printf("\r[%d/%d]\033[K%s", index+1, len(e.Sections), title)
	if volNum, volTitle, isVol := reg.VolByDefaultReg(title); isVol {
		e.VolIndex += 1
		volPath, err = e.AddSection(fmt.Sprintf(config.Cfg.Style.Vol, e.ColImg, volNum, volTitle),
			volNum+" "+volTitle, volFilePrefix+strconv.Itoa(e.VolIndex)+".xhtml", e.Assets.MainCss)
		if err != nil {
			stdout.Errln(err)
			return
		}
		return
	}

	num, name, subNum := reg.ChapterTitleByDefaultReg(title)
	if vol == "" {
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Chapter+e.Sections[index].Content,
			e.ContentLogo, num, name, subNum), strings.Join([]string{num, name, subNum}, " "),
			chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.Assets.MainCss)
		if err != nil {
			stdout.Errln(err)
			return
		}
	} else {
		_, err = e.AddSubSection(vol, fmt.Sprintf(config.Cfg.Chapter+e.Sections[index].Content,
			e.ContentLogo, num, name, subNum), strings.Join([]string{num, name, subNum}, " "),
			chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.Assets.MainCss)
		if err != nil {
			stdout.Errln(err)
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
	return e.Epub.Write(e.Name + "-" + e.Author + ".epub")
}
