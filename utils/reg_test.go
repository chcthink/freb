package utils

import (
	"testing"
)

func TestPureText(t *testing.T) {
	// fmt.Println(PureTitle("833.请假条"))
	checkTitle := func(title, dstTitle, pureTitle string, t *testing.T) {
		if pureTitle != dstTitle {
			t.Errorf("\ntitle: %s \ndest titile: %s \nformat title: %s", title, dstTitle, pureTitle)
		}
	}
	t.Run("2.第2章 是谁偷袭我？！", func(t *testing.T) {
		title := "2.第2章 是谁偷袭我？！"
		dstTitle := "第2章 是谁偷袭我？！"
		checkTitle(title, dstTitle, PureTitle(title), t)
	})
	t.Run("1. 第1章 序 拔剑（可看可不看）", func(t *testing.T) {
		title := "1. 第1章 序 拔剑（可看可不看）"
		dstTitle := "序 拔剑（可看可不看）"
		checkTitle(title, dstTitle, PureTitle(title), t)
	})
	t.Run("第5章 序", func(t *testing.T) {
		title := "第5章 序"
		dstTitle := "序"
		checkTitle(title, dstTitle, PureTitle(title), t)
	})
	t.Run("468.第462章 北海龟山", func(t *testing.T) {
		title := "468.第462章 北海龟山"
		dstTitle := "第462章 北海龟山"
		checkTitle(title, dstTitle, PureTitle(title), t)

	})
	t.Run("第68章 援军 2", func(t *testing.T) {
		title := "第68章 援军 2"
		dstTitle := "第68章 援军 2"
		checkTitle(title, dstTitle, PureTitle(title), t)
	})
	t.Run("第255章 254：将他带往虚夜宫", func(t *testing.T) {
		title := "第255章 254：将他带往虚夜宫"
		dstTitle := "254：将他带往虚夜宫"
		checkTitle(title, dstTitle, PureTitle(title), t)
	})
	t.Run("羊没好", func(t *testing.T) {
		title := "羊没好"
		dstTitle := "羊没好"
		checkTitle(title, dstTitle, PureTitle(title), t)
	})
}
