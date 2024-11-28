package models

import "bytes"

// Book 图书
// Name 书名 Author 作者
// Cover 封面 SubCover 章封面 Vol 卷图
// Chapters 章节
type Book struct {
	Name     string
	Id       string
	Cover    string
	SubCover string
	Vol      string
	Author   string
	Intro    string
	Format   string
	Out      string
	Lang     string
	Desc     bool
	IsOld    bool
	// Font     string
	Chapters []Chapter
}

type Chapter struct {
	Url     string
	Title   *bytes.Buffer
	Content *bytes.Buffer
}
