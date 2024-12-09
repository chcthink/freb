package models

// Book 图书
// Name 书名 Author 作者
// Cover 封面 ContentImg 章封面 Vol 卷图
// Chapters 章节
type Book struct {
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
	IsDesc     bool
	IsOld      bool
	// Font     string
	Chapters []Chapter
}

type Chapter struct {
	Url      string
	Title    string
	Content  string
	Sections []Chapter
}
