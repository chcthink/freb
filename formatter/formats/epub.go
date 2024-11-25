package formats

import (
	"fmt"
	"freb/config"
	"freb/models"
	"freb/utils"
	"github.com/go-shiori/go-epub"
	"path"
	"strconv"
)

// 章节
const (
	chapterFilePrefix = `chapter_`
	htmlP             = `<p>`
	htmlPEnd          = `</p>`
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
	e.InnerURL.css, err = e.AddCSS("assets/styles/main.css", "main.css")
	if err != nil {
		utils.Err(err)
		return
	}
	e.AddFont("assets/fonts/font.ttf", "font.ttf")
	_, err = e.AddCSS("assets/styles/fonts.css", "fonts.css")
	if err != nil {
		utils.Err(err)
		return
	}
	err = e.AddMetaINF("assets/META-INF/com.apple.ibooks.display-options.xml")
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
		coverCss, err = e.AddCSS("assets/styles/cover.css", "cover.css")
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
		insPageCss, err = e.AddCSS("assets/styles/instruction.css", "instruction.css")
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
		logo, err = e.AddImage("assets/images/desc_logo.png", "desc_logo.png")
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
		e.contentLogo, err = e.AddImage("assets/images/desc_logo.png", "content_logo.png")
		if err != nil {
			utils.Err(err)
			return
		}
	}
	return
}

func (e *EpubFormat) GenContentPrefix(i int, str string) {
	e.Book.Chapters[i].Content.WriteString(htmlP + str + htmlPEnd)
}

func (e *EpubFormat) GenBookContent(index int) (err error) {
	title := e.Book.Chapters[index].Title.String()
	fmt.Printf("\r[%d/%d]\033[K%s", index+1, len(e.Chapters), title)
	if index+1 == len(e.Chapters) {
		fmt.Println()
	}
	if volNum, vol, isVol := utils.VolByDefaultReg(title); isVol {
		_, err = e.AddSection(fmt.Sprintf(config.Cfg.Vol, e.volImage, volNum, vol),
			volNum+" "+vol, volNum+" "+vol+".xhtml", e.InnerURL.css)
		if err != nil {
			utils.Err(err)
			return
		}
		return
	}

	num, name, subNum := utils.ChapterTitleByDefaultReg(title)
	_, err = e.AddSection(fmt.Sprintf(config.Cfg.Chapter+e.Book.Chapters[index].Content.String(),
		e.contentLogo, num, name, subNum), num+" "+name,
		chapterFilePrefix+strconv.Itoa(index+1)+".xhtml", e.InnerURL.css)
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
