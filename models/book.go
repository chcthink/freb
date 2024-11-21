package models

import "bytes"

// Book 图书
// Name 书名 Cover 封面 Chapters 章节
// Author 作者 SubCover 章封面
type Book struct {
	Name     string
	Url      string
	Cover    string
	SubCover string
	Author   string
	Intro    string
	Format   string
	Out      string
	Font     string
	Chapters []Chapter
}

type Chapter struct {
	Url     string
	Title   *bytes.Buffer
	Content *bytes.Buffer
}
