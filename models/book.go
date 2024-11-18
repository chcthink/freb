package models

// Book 图书
// Name 书名 Cover 封面 ChapterUrls 章节
// Author 作者 SubCover 章封面
type Book struct {
	Name        string
	Url         string
	Cover       string
	SubCover    string
	Author      string
	Intro       string
	Format      string
	Out         string
	ChapterUrls []string
}
