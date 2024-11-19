package formats

import (
	"bytes"
	"fmt"
	"freb/config"
	"freb/models"
	"freb/utils"
	"github.com/go-shiori/go-epub"
	"os"
	"strconv"
)

// 章节
const (
	CHAPTER_FILE_PREFIX = `chapter_`
	HTML_P              = `<p>`
	HTMLP_END           = `</p>`
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

	// 添加 css
	e.InnerURL.css, err = e.AddCSS("assets/styles/main.css", "main.css")
	if err != nil {
		utils.Err(err)
		return
	}
	e.AddFont("assets/fonts/975MaruSC-Medium.ttf", "975MaruSC-Medium.ttf")
	_, err = e.AddCSS("assets/styles/fonts.css", "fonts.css")
	if err != nil {
		utils.Err(err)
		return
	}
	// 添加标题
	utils.Fmtf("初始化书籍标题:%s", e.Name)
	e.SetTitle(e.Book.Name)
	// 添加封面
	utils.Fmt("初始化书籍封面...")
	if e.Cover != "" {
		var image string
		image, err = e.AddImage(e.Cover, "cover.png")
		if err != nil {
			err = fmt.Errorf("添加封面失败 %w", err)
			return
		}
		err = e.SetCover(image, "")
		if err != nil {
			utils.Err(err)
			return
		}
	}

	// 添加作者
	utils.Fmtf("写入书籍作者:%s", e.Book.Author)
	e.SetAuthor(e.Book.Author)
	// 添加制作说明
	utils.Fmt("正在添加制作说明...")
	insPageCss, err := e.AddCSS("assets/styles/instruction.css", "instruction.css")
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
	// 内容简介
	utils.Fmt("正在添加内容简介...")
	logo, err := e.AddImage("assets/images/desc_logo.png", "desc_logo.png")
	if err != nil {
		return
	}
	_, err = e.AddSection(fmt.Sprintf(config.Cfg.Desc.Dom, logo, e.Book.Intro), config.Cfg.Desc.Title, "desc.xhtml", e.InnerURL.css)
	if err != nil {
		return
	}

	// 初始化章节静态资源
	utils.Fmt("初始化章节静态资源...")
	if e.Book.SubCover != "" {
		e.volImage, err = e.AddImage(e.Book.SubCover, "vol.png")
		if err != nil {
			utils.Err(err)
			return
		}
	}
	e.contentLogo, err = e.AddImage("assets/images/desc_logo.png", "content_logo.png")
	if err != nil {
		utils.Err(err)
		return
	}
	return
}

func (e *EpubFormat) GenContentPrefix(buf *bytes.Buffer, str string) {
	buf.WriteString(HTML_P + str + HTMLP_END)
}

func (e *EpubFormat) GenBookContent(index int, title, content bytes.Buffer) (err error) {
	utils.Fmtf("\r%s %d/%d", title.String(), index, len(e.ChapterUrls))
	_ = os.Stdout.Sync()

	if volNum, vol, isVol := utils.VolByDefaultReg(title.String()); isVol {
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Vol, e.volImage, volNum, vol),
			volNum+" "+vol, volNum+" "+vol+".xhtml", e.InnerURL.css)
		if err != nil {
			utils.Err(err)
			return
		}
		return
	}

	num, name, subNum := utils.ChapterTitleByDefaultReg(title.String())
	_, err = e.AddSection(fmt.Sprintf(config.Cfg.Chapter+content.String(), e.contentLogo, num, name, subNum),
		title.String(), CHAPTER_FILE_PREFIX+strconv.Itoa(index)+".xhtml", e.InnerURL.css)
	if err != nil {
		utils.Err(err)
		return
	}
	return
}

func (e *EpubFormat) Build() error {
	if e.Out != "" {
		return e.Epub.Write(e.Out)
	}
	return e.Epub.Write(e.Name + "-" + e.Book.Author + ".epub")
}
