package utils

import (
	"regexp"
	"strings"
)

var (
	chapterNumRegs    []*regexp.Regexp
	chapterPreNumRegs []*regexp.Regexp
	introReg          *regexp.Regexp
	chapterSubNumReg  *regexp.Regexp
	volReg            *regexp.Regexp
	authorReg         *regexp.Regexp
	endReg            *regexp.Regexp
	urlReg            *regexp.Regexp
	numReg            *regexp.Regexp
	c0ControlReg      *regexp.Regexp
)

// init reg
func init() {
	var chapterPreNum = []string{
		`^引子$`,
		`^楔子$`,
		`^序章?`,
	}
	var chapterNum = []string{
		`^章节(目录)?`,
		`第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷部]`,
		`^\d+\.?`,
		`^[Ss]ection.{1,20}$`,
		`^[Cc]hapter.{1,20}$`,
		`^[Pp]age.{1,20}$`,
	}
	for _, pattern := range append(chapterPreNum) {
		re := regexp.MustCompile(pattern)
		chapterPreNumRegs = append(chapterPreNumRegs, re)
	}
	for _, pattern := range append(chapterPreNum, chapterNum...) {
		re := regexp.MustCompile(pattern)
		chapterNumRegs = append(chapterNumRegs, re)
	}
	introReg = regexp.MustCompile("(文章|内容)简介([:：])?")
	chapterSubNumReg = regexp.MustCompile("[(（][0-9一二三四五六七八九十零〇百千两上中下 ][)）]")
	volReg = regexp.MustCompile("^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]")
	authorReg = regexp.MustCompile("作者([:：])?")
	endReg = regexp.MustCompile("大结局|最终话")
	urlReg = regexp.MustCompile("^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
	numReg = regexp.MustCompile("[0-9]+")
	c0ControlReg = regexp.MustCompile(`[\x00-\x1F]`)
}

func CheckUrl(url string) bool {
	if len(url) == 0 {
		return false
	}
	return urlReg.MatchString(url)
}

func ChapterTitleByDefaultReg(str string) (num, title, subNum string) {
	num = findNumInTitle(str)
	subNum = chapterSubNumReg.FindString(str)
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
	str = strings.TrimSpace(str)
	nStr := strings.Split(str, " ")
	// find latest num
	num := findNumInTitle(str)
	if len(nStr) > 1 {
		for i := len(nStr) - 1; i >= 0; i-- {
			num = findNumInTitle(nStr[i])
			if num != "" {
				break
			}
		}
	}

	if num != "" {
		subTitle := strings.TrimSpace(strings.Split(str, num)[1])
		if subTitle == "" {
			title = num
			return
		}
		return strings.Join([]string{num, subTitle}, " ")
	}
	return
}

func checkPreNumInTitle(str string) bool {
	for _, reg := range chapterPreNumRegs {
		if reg.MatchString(str) {
			return reg.MatchString(str)
		}
	}
	return false
}

func findNumInTitle(str string) string {
	for _, reg := range chapterNumRegs {
		if reg.MatchString(str) {
			return reg.FindString(str)
		}
	}
	return ""
}

func VolByDefaultReg(str string) (num, title string, isVol bool) {
	num = volReg.FindString(str)
	if num != "" {
		isVol = true
		title = strings.TrimSpace(strings.Split(str, num)[1])
	}
	return
}

func CheckVol(str string) bool {
	return volReg.MatchString(str)
}

func CheckTitle(str string) bool {
	for _, reg := range chapterNumRegs {
		if reg.MatchString(str) {
			return true
		}
	}
	return false
}

func CheckEnd(str string) bool {
	return endReg.MatchString(str)
}

func GetAuthor(str string) (isAuthor bool, author string) {
	isAuthor = authorReg.MatchString(str)
	if isAuthor {
		author = strings.TrimSpace(authorReg.ReplaceAllString(str, ""))
	}
	return
}

func AuthorIndex(str string) int {
	if authorReg.MatchString(str) {
		return authorReg.FindStringIndex(str)[1]
	}
	return -1
}

func CheckNum(str string) bool {
	return numReg.MatchString(str)
}

func GetIntro(str string) (isIntro bool, intro string) {
	isIntro = introReg.MatchString(str)
	if isIntro {
		intro = strings.TrimSpace(introReg.ReplaceAllString(str, ""))
	}
	return
}

func ReplaceC0Control(str string) string {
	return c0ControlReg.ReplaceAllString(str, "")
}
