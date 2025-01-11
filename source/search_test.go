package source

import (
	"fmt"
	"freb/config"
	"testing"
)

func TestGenUrl(t *testing.T) {
	config.InitConfig()
	url := genUrl("阵问长生", config.Cfg.BookCatch)
	fmt.Println(url)
}

func TestSearch(t *testing.T) {
	config.InitConfig()

	url, err := Search("阵问长生", config.Cfg.BookCatch)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(url)
}
