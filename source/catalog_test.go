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
				Url:    "https://www.qidian.com/book/1035420986/",
				Cookie: "e1=%7B%22l6%22%3A%221%22%2C%22l7%22%3A%22%22%2C%22l1%22%3A%22%22%2C%22l3%22%3A%22%22%2C%22pid%22%3A%22qd_P_xiangqing%22%2C%22eid%22%3A%22%22%7D; e2=%7B%22l6%22%3A%221%22%2C%22l7%22%3A%22%22%2C%22l1%22%3A%22%22%2C%22l3%22%3A%22%22%2C%22pid%22%3A%22qd_P_xiangqing%22%2C%22eid%22%3A%22%22%7D; newstatisticUUID=1728217849_1979792855; _csrfToken=OKiYHYB6eLNnXewMkNxbA67cFsgwhIkuCjp6qmEo; fu=1949850521; supportWebp=true; traffic_search_engine=; se_ref=; e1=%7B%22l6%22%3A%22%22%2C%22l7%22%3A%22%22%2C%22l1%22%3A3%2C%22l3%22%3A%22%22%2C%22pid%22%3A%22qd_p_qidian%22%2C%22eid%22%3A%22qd_A17%22%7D; e2=%7B%22l6%22%3A%22%22%2C%22l7%22%3A%22%22%2C%22l1%22%3A3%2C%22l3%22%3A%22%22%2C%22pid%22%3A%22qd_p_qidian%22%2C%22eid%22%3A%22qd_A16%22%7D; supportwebp=true; traffic_utm_referer=; x-waf-captcha-referer=; w_tsfp=ltvuV0MF2utBvS0Q6q/tk0quET0gfDo4h0wpEaR0f5thQLErU5mG1IJ9v8vwNXzd58xnvd7DsZoyJTLYCJI3dwMXR8WYJYlH3V6WwYd33olBCUExEJKKD1RKd+527DQXdXhCNxS00jA8eIUd379yilkMsyN1zap3TO14fstJ019E6KDQmI5uDW3HlFWQRzaLbjcMcuqPr6g18L5a5TrctlquKFx9UbhF1EHBhysXXXpw4EO7Ju0LNxz7d8enSqA=",
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
