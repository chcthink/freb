package reg

import (
	"fmt"
	"freb/models"
	"github.com/dlclark/regexp2"
	"regexp"
	"strings"
)

var (
	chapterNumRegs      []*Reg
	chapterPrologueRegs = &Reg{}
	chapterNumReg       = &Reg{}
	introReg            = &Reg{}
	chapterSubNumReg    = &Reg{}
	volReg              = &Reg{}
	authorReg           = &Reg{}
	endReg              = &Reg{}
	urlReg              *regexp.Regexp
	numReg              *regexp.Regexp
	c0ControlReg        *regexp.Regexp
)

var (
	rmTitleReg   []*Reg
	rmContentReg []*Reg
)

// init reg
func init() {
	urlReg = regexp.MustCompile("^https?://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
	numReg = regexp.MustCompile("[0-9]+")
	c0ControlReg = regexp.MustCompile(`[\x00-\x1F]`)
}

func InitCustomMatchReg(regs *models.Regs) {
	for _, pattern := range append(regs.ChapterTitle.Prologue, regs.ChapterTitle.Num...) {
		reg := new(Reg)
		reg.MustCompile(pattern)
		chapterNumRegs = append(chapterNumRegs, reg)
	}
	chapterPrologueRegs.MustCompile(strings.Join(regs.ChapterTitle.Prologue, "|"))
	introReg.MustCompile(regs.Intro)
	chapterNumReg.MustCompile(strings.Join(regs.ChapterTitle.Num, "|"))
	chapterSubNumReg.MustCompile(regs.ChapterTitle.SubNum)
	volReg.MustCompile(regs.Vol)
	authorReg.MustCompile(regs.Author)
	endReg.MustCompile(regs.End)

}

func InitCustomFilterReg(catch *models.BookCatch) {
	for _, reg := range catch.Title.Filter {
		r := new(Reg)
		r.MustCompile(reg)
		rmTitleReg = append(rmTitleReg, r)
	}
	for _, reg := range catch.Content.Filter {
		r := new(Reg)
		r.MustCompile(reg)
		rmContentReg = append(rmContentReg, r)
	}
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

func ChapterTitleWithoutNum(str string) (title string) {
	// str = html.UnescapeString(str)
	num := findNumInTitle(str)
	title = str
	if num != "" {
		title = strings.Split(title, num)[1]
	}
	return
}

func PureTitle(str string) (title string) {
	str = strings.TrimSpace(str)
	num := findNumInTitle(str)
	title = str
	if num != "" {
		subTitle := strings.TrimSpace(strings.Split(str, num)[1])
		if subTitle == "" {
			title = num
			return
		}
		if strings.HasSuffix(num, "：") {
			return strings.Join([]string{num, subTitle}, "")
		}
		return strings.Join([]string{num, subTitle}, " ")
	}
	return
}

func findNumInTitle(str string) (match string) {
	strs := strings.Split(str, " ")
	for _, s := range strs {
		tmp := strings.TrimSpace(s)
		if chapterPrologueRegs.MatchString(tmp) {
			return s
		}
		nums := chapterNumReg.FindAllString(tmp)
		for i := range nums {
			t := strings.Split(str, nums[i])
			if strings.TrimSpace(t[len(t)-1]) != "" {
				match = nums[i]
			}
		}
	}
	return
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
		index, length := authorReg.FindStringIndex(str)
		return index + length
	}
	return -1
}

func CheckNum(str string) bool {
	return numReg.MatchString(str)
}

func GetNum(str string) string {
	return numReg.FindString(str)
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

func RemoveTitleFromCfg(str string) (dst string) {
	dst = str
	for _, reg := range rmTitleReg {
		if reg.MatchString(dst) {
			dst = reg.ReplaceAllString(dst, "")
		}
	}
	return
}

func RemoveContentFromCfg(str string) (dst string) {
	dst = str
	for _, reg := range rmContentReg {
		if reg.MatchString(dst) {
			dst = reg.ReplaceAllString(dst, "")
		}
	}
	return
}

const (
	reg2MatchErr = "正则匹配异常: %s %s"
)

func MatchString(exp, str string) (dest string, err error) {
	match, err := regexp2.MustCompile(exp, regexp2.None).FindStringMatch(str)
	if err != nil {
		return "", fmt.Errorf(reg2MatchErr, exp, str)
	}
	dest = match.String()
	return
}

func Filters(exps []string, str string) (dest string, err error) {
	if len(exps) == 0 {
		return str, nil
	}
	var regs []*regexp2.Regexp
	for _, exp := range exps {
		var reg *regexp2.Regexp
		reg, err = regexp2.Compile(exp, regexp2.None)
		if err != nil {
			return "", fmt.Errorf(reg2MatchErr, exp, str)
		}
		regs = append(regs, reg)
	}
	dest = str
	for _, reg := range regs {
		if isExist, _ := reg.MatchString(dest); isExist {
			dest, _ = reg.Replace(dest, "", -1, -1)
		}
	}
	return
}

const (
	bookNameMarkPre = '《'
	bookNameMarkSuf = '》'
	hyphen          = '-'
)

func GetBookInfo(str string) (name, author string) {
	if strings.ContainsRune(str, bookNameMarkPre) && strings.ContainsRune(str, bookNameMarkSuf) {
		start := strings.IndexRune(str, bookNameMarkPre)
		end := strings.IndexRune(str, bookNameMarkSuf)
		if start >= end {
			return
		}
		name = strings.Replace(str[start:end], string(bookNameMarkPre), "", 1)
	}
	if index := AuthorIndex(str); index > 0 {
		author = str[index:]
	}
	if strings.ContainsRune(str, hyphen) {
		name = str[:strings.IndexRune(str, hyphen)]
		author = str[strings.IndexRune(str, hyphen):]
	}
	return
}
