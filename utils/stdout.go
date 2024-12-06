package utils

import (
	"github.com/fatih/color"
)

const (
	ErrUrl  = "url 输入错误"
	newLine = "\n"
)

var outB = color.New(color.Bold)

func Err(err error) {
	outB.Add(color.FgRed)
	_, _ = outB.Println(err.Error())
}

func Errf(str string, a ...interface{}) {
	outB.Add(color.FgRed)
	_, _ = outB.Printf(str+newLine, a...)
}

func Successf(str string, a ...interface{}) {
	outB.Add(color.FgGreen)
	_, _ = outB.Printf(str+newLine, a...)
}

func Fmt(str string) {
	outB.Add(color.FgWhite)
	_, _ = outB.Println(str)
}

func Fmtf(str string, a ...interface{}) {
	outB.Add(color.FgWhite)
	_, _ = outB.Printf(str+newLine, a...)
}

func Warnf(str string, a ...interface{}) {
	outB.Add(color.FgYellow)
	_, _ = outB.Printf(str+newLine, a...)
}

func SysInfof(str string, a ...interface{}) {
	outB.Add(color.FgBlue)
	_, _ = outB.Printf(str+newLine, a...)
}
