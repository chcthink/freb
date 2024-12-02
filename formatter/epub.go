package formatter

import (
	"bytes"
	"fmt"
	"freb/config"
	"freb/models"
	"freb/utils"
	"github.com/go-shiori/go-epub"
	"path"
	"strconv"
	"strings"
)

// 章节
const (
	chapterFilePrefix = `chapter_`
	htmlP             = `<p>`
	htmlPEnd          = `</p>`
	percentSign       = "%"
)

type InnerURL struct {
	volImage    string
	css         string
	contentLogo string
}

type EpubFormat struct {
	*epub.Epub
	*models.Book
	*InnerURL
}

func (e *EpubFormat) InitBook() (err error) {
	e.Epub, err = epub.NewEpub(e.Book.Name)
	e.InnerURL = &InnerURL{}
	if err != nil {
		utils.Err(err)
		return
	}

	// 初始化书籍信息
	utils.Fmtf("初始化书籍信息:%s", e.Name)
	// 添加 css
	e.InnerURL.css, err = e.AddCSS(utils.LocalOrUrl("assets/styles/main.css"), "main.css")
	if err != nil {
		utils.Err(err)
		return
	}
	e.AddFont(utils.LocalOrUrl("assets/fonts/font.ttf"), "font.ttf")
	_, err = e.AddCSS(utils.LocalOrUrl("assets/styles/fonts.css"), "fonts.css")
	if err != nil {
		utils.Err(err)
		return
	}
	err = e.AddMetaINF(utils.LocalOrUrl("assets/META-INF/com.apple.ibooks.display-options.xml"))
	if err != nil {
		utils.Err(err)
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
		coverCss, err = e.AddCSS(utils.LocalOrUrl("assets/styles/cover.css"), "cover.css")
		if err != nil {
			utils.Err(err)
			return
		}
		err = e.SetCover(image, coverCss)
		if err != nil {
			utils.Err(err)
			return
		}
	}
	// 添加作者
	e.SetAuthor(e.Book.Author)
	// 添加制作说明
	if e.Book.Desc {
		utils.Fmt("正在添加制作说明...")
		var insPageCss string
		insPageCss, err = e.AddCSS(utils.LocalOrUrl("assets/styles/instruction.css"), "instruction.css")
		if err != nil {
			utils.Err(err)
			return
		}
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Instruction.Dom, e.Book.Name, e.Book.Author),
			config.Cfg.Instruction.Title, "instruction.xhtml", insPageCss)
		if err != nil {
			utils.Err(err)
			return
		}
	}
	// 内容简介
	if e.Book.Intro != "" {
		utils.Fmt("正在添加内容简介...")
		var logo string
		logo, err = e.AddImage(utils.LocalOrUrl("assets/images/desc_logo.png"), "desc_logo.png")
		if err != nil {
			return
		}
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Desc.Dom, logo, e.Book.Intro), config.Cfg.Desc.Title, "desc.xhtml", e.InnerURL.css)
		if err != nil {
			return
		}
	}
	if e.Book.Vol != "" {
		e.volImage, err = e.AddImage(e.Book.Vol, "vol.png")
		if err != nil {
			utils.Err(err)
			return
		}
	}
	if e.Book.SubCover != "" {
		e.contentLogo, err = e.AddImage(e.Book.SubCover, "sub_cover.png")
		if err != nil {
			utils.Err(err)
			return
		}
	} else {
		e.contentLogo, err = e.AddImage(utils.LocalOrUrl("assets/images/desc_logo.png"), "content_logo.png")
		if err != nil {
			utils.Err(err)
			return
		}
	}
	return
}

func (e *EpubFormat) GenLine(str string) string {
	str = utils.PureEscapeHtml(str)
	return htmlP + strings.ReplaceAll(str, percentSign, "%%") + htmlPEnd
}

func (e *EpubFormat) GenLine2Buffer(str string, buf *bytes.Buffer) {
	str = utils.PureEscapeHtml(str)
	buf.WriteString(htmlP + strings.ReplaceAll(str, percentSign, "%%") + htmlPEnd)
}

func (e *EpubFormat) GenBookContent(index int, vol string) (volPath string, err error) {
	title := e.Book.Chapters[index].Title
	fmt.Printf("\r[%d/%d]\033[K%s", index+1, len(e.Chapters), title)
	if index+1 == len(e.Chapters) {
		fmt.Println()
	}
	if volNum, vol, isVol := utils.VolByDefaultReg(title); isVol {
		volPath, err = e.AddSection(fmt.Sprintf(config.Cfg.Vol, e.volImage, volNum, vol),
			volNum+" "+vol, volNum+" "+vol+".xhtml", e.InnerURL.css)
		if err != nil {
			utils.Err(err)
			return
		}
		return
	}

	num, name, subNum := utils.ChapterTitleByDefaultReg(title)
	if vol == "" {
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Chapter+e.Book.Chapters[index].Content,
			e.contentLogo, num, name, subNum), num+" "+name,
			chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.InnerURL.css)
		if err != nil {
			utils.Err(err)
			return
		}
	} else {
		_, err = e.AddSubSection(vol, fmt.Sprintf(config.Cfg.Chapter+e.Book.Chapters[index].Content,
			e.contentLogo, num, name, subNum), num+" "+name,
			chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.InnerURL.css)
		if err != nil {
			utils.Err(err)
			return
		}
	}

	return
}

func (e *EpubFormat) Build() error {
	if e.Out != "" {
		return e.Epub.Write(e.Out)
	}
	return e.Epub.Write(e.Name + "-" + e.Book.Author + ".epub")
}
