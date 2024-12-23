package utils

import (
	"fmt"
	"testing"
)

func TestPureText(t *testing.T) {
	fmt.Println(PureTitle("第5章 序"))
	// fmt.Println(utils.PureTitle("1. 第1章 序 拔剑（可看可不看）"))
	fmt.Println(PureTitle("2.第2章 是谁偷袭我？！"))
	fmt.Println(PureTitle("833.请假条"))
	fmt.Println(PureTitle("468.第462章 北海龟山"))
}
