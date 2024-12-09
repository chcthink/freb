package utils

import (
	"regexp"
	"strings"
)

// reg
const (
	checkURL  = "^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]"
	urlDomain = "(https?://)[^/]+"
)

const (
	DefaultImagePath = "assets/images/"
)

func CheckUrl(url string) bool {
	if len(url) == 0 {
		return false
	}
	re, _ := regexp.Compile(checkURL)
	return re.MatchString(url)
}

const (
	IntroReg         = "(文章|内容)简介([:：])?"
	ChapterNumReg    = "第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷部]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^引子$|^楔子$|^章节目录|^章节|^序章"
	ChapterSubNumReg = "[(（][0-9一二三四五六七八九十零〇百千两 ][)）]"
	volReg           = "^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]"
	authorReg        = "作者([:：])?"
	endReg           = "大结局|最终话"
)

func ChapterTitleByDefaultReg(str string) (num, title, subNum string) {
	numTitleReg := regexp.MustCompile(ChapterNumReg)
	SubNumReg := regexp.MustCompile(ChapterSubNumReg)
	num = numTitleReg.FindString(str)
	subNum = SubNumReg.FindString(str)
	title = str
	if num != "" {
		title = strings.Split(title, num)[1]
	}
	if subNum != "" {
		title = strings.Split(title, subNum)[0]
	}
	if title != "" {
		title = strings.TrimSpace(title)
	}
	return
}

func PureTitle(str string) (title string) {
	numTitleReg := regexp.MustCompile(ChapterNumReg)
	num := numTitleReg.FindString(str)
	title = str
	if num != "" {
		return num + " " + strings.TrimSpace(strings.Split(title, num)[1])
	}
	return str
}

func VolByDefaultReg(str string) (num, title string, isVol bool) {
	reg := regexp.MustCompile(volReg)
	num = reg.FindString(str)
	if num != "" {
		isVol = true
		title = strings.TrimSpace(strings.Split(str, num)[1])
	}
	return
}

func CheckVol(str string) bool {
	reg := regexp.MustCompile(volReg)
	return reg.MatchString(str)
}

func CheckTitle(str string) bool {
	reg := regexp.MustCompile(ChapterNumReg)
	return reg.MatchString(str)
}

func CheckEnd(str string) bool {
	reg := regexp.MustCompile(endReg)
	return reg.MatchString(str)
}

func GetAuthor(str string) (isAuthor bool, author string) {
	reg := regexp.MustCompile(authorReg)
	isAuthor = reg.MatchString(str)
	if isAuthor {
		author = strings.TrimSpace(reg.ReplaceAllString(str, ""))
	}
	return
}

func AuthorIndex(str string) int {
	reg := regexp.MustCompile(authorReg)
	if reg.MatchString(str) {
		return reg.FindStringIndex(str)[1]
	}
	return -1
}

const (
	numReg = "[0-9]+"
)

func CheckDomain(str string) bool {
	reg := regexp.MustCompilePOSIX(numReg)
	return reg.MatchString(str)
}

func GetIntro(str string) (isIntro bool, intro string) {
	reg := regexp.MustCompile(IntroReg)
	isIntro = reg.MatchString(str)
	if isIntro {
		intro = strings.TrimSpace(reg.ReplaceAllString(str, ""))
	}
	return
}

const (
	c0ControlReg = `[\x00-\x1F]`
)

func ReplaceC0Control(str string) string {
	reg := regexp.MustCompile(c0ControlReg)
	return reg.ReplaceAllString(str, "")
}
