package stdout

import (
	"github.com/fatih/color"
)

const (
	ErrUrl  = "url 输入错误"
	newLine = "\n"
)

var outB = color.New(color.Bold)

func Errln(err error) {
	outB.Add(color.FgRed)
	_, _ = outB.Println(err.Error())
}

func Errfln(str string, a ...interface{}) {
	outB.Add(color.FgRed)
	_, _ = outB.Printf(str+newLine, a...)
}

func Successfln(str string, a ...interface{}) {
	outB.Add(color.FgGreen)
	_, _ = outB.Printf(str+newLine, a...)
}

func Fmtln(str string) {
	outB.Add(color.FgWhite)
	_, _ = outB.Println(str)
}

func Fmt(str string) {
	outB.Add(color.FgWhite)
	_, _ = outB.Print(str)
}

func Fmtfln(str string, a ...interface{}) {
	outB.Add(color.FgWhite)
	_, _ = outB.Printf(str+newLine, a...)
}

func Warnf(str string, a ...interface{}) {
	outB.Add(color.FgYellow)
	_, _ = outB.Printf(str, a...)
}

func SysInfofln(str string, a ...interface{}) {
	outB.Add(color.FgBlue)
	_, _ = outB.Printf(str+newLine, a...)
}

func Contentln(str string) {
	outB.Add(color.FgGreen)
	_, _ = outB.Println(str)
}
