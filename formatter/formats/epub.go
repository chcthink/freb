package formats

import (
	"bytes"
	"fmt"
	"freb/models"
	"freb/utils"
	"github.com/go-shiori/go-epub"
)

// instruction 制作说明
const (
	INSTRUCTION_TITLE = `制作说明`
	INSTRUCTION_HTML  = `
<h3 class="ver">制作说明</h3>
    <p class="ver-char"><span class="verchar_01">%s</span></p>
    <p class="ver-title_01">%s&#160;◎著</p>
    <br />
    <hr class="line" />
    <p class="ver-txt">
      制作：<svg
        xmlns="http://www.w3.org/2000/svg"
        height="14"
        width="14"
        viewBox="0 0 496 512"
      >
        <!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->
        <path
          d="M165.9 397.4c0 2-2.3 3.6-5.2 3.6-3.3 .3-5.6-1.3-5.6-3.6 0-2 2.3-3.6 5.2-3.6 3-.3 5.6 1.3 5.6 3.6zm-31.1-4.5c-.7 2 1.3 4.3 4.3 4.9 2.6 1 5.6 0 6.2-2s-1.3-4.3-4.3-5.2c-2.6-.7-5.5 .3-6.2 2.3zm44.2-1.7c-2.9 .7-4.9 2.6-4.6 4.9 .3 2 2.9 3.3 5.9 2.6 2.9-.7 4.9-2.6 4.6-4.6-.3-1.9-3-3.2-5.9-2.9zM244.8 8C106.1 8 0 113.3 0 252c0 110.9 69.8 205.8 169.5 239.2 12.8 2.3 17.3-5.6 17.3-12.1 0-6.2-.3-40.4-.3-61.4 0 0-70 15-84.7-29.8 0 0-11.4-29.1-27.8-36.6 0 0-22.9-15.7 1.6-15.4 0 0 24.9 2 38.6 25.8 21.9 38.6 58.6 27.5 72.9 20.9 2.3-16 8.8-27.1 16-33.7-55.9-6.2-112.3-14.3-112.3-110.5 0-27.5 7.6-41.3 23.6-58.9-2.6-6.5-11.1-33.3 2.6-67.9 20.9-6.5 69 27 69 27 20-5.6 41.5-8.5 62.8-8.5s42.8 2.9 62.8 8.5c0 0 48.1-33.6 69-27 13.7 34.7 5.2 61.4 2.6 67.9 16 17.7 25.8 31.5 25.8 58.9 0 96.5-58.9 104.2-114.8 110.5 9.2 7.9 17 22.9 17 46.4 0 33.7-.3 75.4-.3 83.6 0 6.5 4.6 14.4 17.3 12.1C428.2 457.8 496 362.9 496 252 496 113.3 383.5 8 244.8 8zM97.2 352.9c-1.3 1-1 3.3 .7 5.2 1.6 1.6 3.9 2.3 5.2 1 1.3-1 1-3.3-.7-5.2-1.6-1.6-3.9-2.3-5.2-1zm-10.8-8.1c-.7 1.3 .3 2.9 2.3 3.9 1.6 1 3.6 .7 4.3-.7 .7-1.3-.3-2.9-2.3-3.9-2-.6-3.6-.3-4.3 .7zm32.4 35.6c-1.6 1.3-1 4.3 1.3 6.2 2.3 2.3 5.2 2.6 6.5 1 1.3-1.3 .7-4.3-1.3-6.2-2.2-2.3-5.2-2.6-6.5-1zm-11.4-14.7c-1.6 1-1.6 3.6 0 5.9 1.6 2.3 4.3 3.3 5.6 2.3 1.6-1.3 1.6-3.9 0-6.2-1.4-2.3-4-3.3-5.6-2z"
        /></svg
      >&#160;chcthink/freb_bak&#160;
    </p>
    <p class="ver-txt">排版参考：阡陌居-笙歌夜夜</p>
    <p class="ver-txt">
      声明：本书仅作个人排版参考学习之用，请勿用于商业用途。如果喜欢本书，请购买正版。任何对本书的修改、加工、传播，请自负法律后果。
    </p>
    <hr class="line" />
    <p class="ver-note">
      注：为获得最佳阅读效果，请在多看设置中将排版设为“原版”（多看2.x版本）或“无”（多看3.x版本以上），背景为预设背景（不要自定义背景和字体颜色，以免整体配色出问题）；字体设置为“默认”（使用书中指定字体），字体大小为默认大小（一般手机上为+3，平板上为+2——即减小字体到最小值后，点击增大按钮的次数）。
    </p>`
)

// desc
// descCss 简介 css
const (
	DESC_TITLE = `内容简介`
	DESC_HTML  = `
<div class="pg">
  <img alt="logo" class="pg" src="%s" />
</div>
<h2 class="desc">内容简介</h2>
<p class="desc-2">%s</p>
`
)

// 卷
const (
	VOL_HTML = `
<div class="logo2">
	<img alt="logo" class="logo2" src="%s"/>
</div>
<div class="c1">
	<span>%s</span><br/>%s
</div>`
)

// 章节
const (
	CHAPTER_FILE_PREFIX = `chapter_`
	HTML_P              = `<p>`
	HTMLP_END           = `</p>`
	CHAPTER_HTML        = `
<div class="logo"><img alt="logo" class="logo" src="%s" /></div>
<h2><span class="num">%s</span><br /> %s<br /><span class="num-2">%s</span></h2>
`
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
	_, err = e.AddSection(fmt.Sprintf(INSTRUCTION_HTML, e.Book.Name, e.Book.Author), INSTRUCTION_TITLE,
		"instruction.xhtml", insPageCss)
	if err != nil {
		utils.Err(err)
		return
	}
	// 内容简介
	utils.Fmt("正在添加内容简介...")
	logo, err := e.AddImage("assets/images/desc_wawa.png", "desc_wawa.png")
	if err != nil {
		return
	}
	_, err = e.AddSection(fmt.Sprintf(DESC_HTML, logo, e.Book.Intro), DESC_TITLE, "desc.xhtml", e.InnerURL.css)
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
	e.contentLogo, err = e.AddImage("assets/images/desc_wawa.png", "content_logo.png")
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
	utils.Fmtf("%s %d/%d", title.String(), index, len(e.ChapterUrls))
	if volNum, vol, isVol := utils.VolByDefaultReg(title.String()); isVol {
		_, err = e.AddSection(fmt.Sprintf(VOL_HTML, e.volImage, volNum, vol),
			volNum+" "+vol, volNum+" "+vol+".xhtml", e.InnerURL.css)
		if err != nil {
			utils.Err(err)
			return
		}
		return
	}

	num, name, subNum := utils.ChapterTitleByDefaultReg(title.String())
	_, err = e.AddSection(fmt.Sprintf(CHAPTER_HTML+content.String(), e.contentLogo, num, name, subNum),
		title.String(), num+name+".xhtml", e.InnerURL.css)
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
