package utils

import (
	"errors"
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

func UrlErr() {
	Err(errors.New(ErrUrl))
}

func Success(str string) {
	outB.Add(color.FgGreen)
	_, _ = outB.Println(str)
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
