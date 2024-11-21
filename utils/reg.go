package utils

import (
	"regexp"
	"strconv"
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

func GetDomainFromUrl(url string) string {
	reg := regexp.MustCompile(urlDomain)
	return reg.FindString(url)
}

const (
	ChapterNumReg    = "第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集卷部]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^\\d+、|^引子$|^楔子$|^章节目录|^章节|^序章"
	ChapterSubNumReg = "[(（][0-9一二三四五六七八九十零〇百千两 ][)）]"
	volReg           = "^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]"
)

func ChapterTitleByDefaultReg(str string) (num, title, subNum string) {
	numReg := regexp.MustCompile(ChapterNumReg)
	SubNumReg := regexp.MustCompile(ChapterSubNumReg)
	num = numReg.FindString(str)
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

func VolByDefaultReg(str string) (num, title string, isVol bool) {
	reg := regexp.MustCompile(volReg)
	num = reg.FindString(str)
	if num != "" {
		isVol = true
		title = strings.TrimSpace(strings.Split(str, num)[1])
	}
	return
}

func IsTitle(str string) bool {
	reg := regexp.MustCompile(ChapterNumReg)
	return reg.MatchString(str)
}

const (
	bookIdReg = "[0-9]+.htm"
)

func BookId(url string) int {
	reg := regexp.MustCompilePOSIX(bookIdReg)
	suffix := reg.FindString(url)
	if strings.Contains(suffix, ".htm") {
		ret, _ := strconv.Atoi(strings.Split(suffix, ".")[0])
		return ret
	}
	return 0
}
