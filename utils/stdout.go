package utils

import (
	"errors"
	"github.com/fatih/color"
)

const ERR_URL = "url 输入错误"

var outB = color.New(color.Bold)

func Err(err error) {
	outB.Add(color.FgRed)
	_, _ = outB.Println(err.Error())
}

func UrlErr() {
	Err(errors.New(ERR_URL))
}

func Success(str string) {
	outB.Add(color.FgGreen)
	_, _ = outB.Println(str)
}

func Fmt(str string) {
	outB.Add(color.FgWhite)
	_, _ = outB.Println(str)
}

func Fmtf(str string, a ...interface{}) {
	outB.Add(color.FgWhite)
	_, _ = outB.Printf(str+"\n", a...)
}
