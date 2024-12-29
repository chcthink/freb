package models

// BookConf 图书配置
// Name 书名 Author 作者
// Cover 封面 ContentImg 章封面 Vol 卷图
// Chapters 章节
type BookConf struct {
	Name       string
	Id         string
	Cover      string
	ContentImg string
	IntroImg   string
	Vol        string
	Author     string
	Intro      string
	Format     string
	Out        string
	Lang       string
	Path       string
	Jump       int
	Delay      int
	IsDesc     bool
	IsOld      bool
	// Font     string
	Chapters []Chapter
	Catalog  UrlWithCookie
}

type Chapter struct {
	Url     string
	Title   string
	Content string
	IsVol   bool
}

type UrlWithCookie struct {
	Url    string
	Cookie string
}
