package source

import (
	"fmt"
	"freb/formatter"
	"freb/models"
	"testing"
)

func TestGetCatalog(t *testing.T) {
	var ef formatter.EpubFormat

	t.Run("七猫", func(t *testing.T) {
		ef.BookConf = &models.BookConf{
			Catalog: models.UrlWithCookie{
				Url: "https://www.qimao.com/shuku/197091/",
			},
		}
		err := GetCatalogFromUrl(&ef)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(len(ef.Chapters))
	})

	t.Run("起点", func(t *testing.T) {
		ef.BookConf = &models.BookConf{
			Catalog: models.UrlWithCookie{
				Url: "https://book.qidian.com/info/1017281778",
			},
		}
		err := GetCatalogFromUrl(&ef)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(len(ef.Chapters))
	})

	t.Run("番茄", func(t *testing.T) {
		ef.BookConf = &models.BookConf{
			Catalog: models.UrlWithCookie{
				Url: "https://fanqienovel.com/page/7143038691944959011",
			},
		}
		err := GetCatalogFromUrl(&ef)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(len(ef.Chapters))
	})
}
