package source

import (
	"fmt"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestGetCatalog(t *testing.T) {
	var ef formatter.EpubFormat
	err := config.InitConfig()
	if err != nil {
		t.Error(err)
	}

	t.Run("七猫", func(t *testing.T) {
		ef.BookConf = &models.BookConf{
			Catalog: "https://www.qimao.com/shuku/197091/",
		}
		err = GetCatalogFromUrl(&ef)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(len(ef.Sections))
		spew.Dump(ef.Sections[:2])
	})

	t.Run("起点", func(t *testing.T) {
		ef.BookConf = &models.BookConf{
			Catalog: "https://book.qidian.com/info/1017281778",
		}
		err := GetCatalogFromUrl(&ef)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(len(ef.Sections))
		spew.Dump(ef.Sections[:20])
	})

	t.Run("番茄", func(t *testing.T) {
		ef.BookConf = &models.BookConf{
			Catalog: "https://fanqienovel.com/page/7143038691944959011",
		}
		err := GetCatalogFromUrl(&ef)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(len(ef.Sections))
		spew.Dump(ef.Sections[:20])
	})
}
