package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"freb/formatter"
	"freb/models"
	"freb/utils"
	"freb/utils/stdout"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/transform"
	"net/http"
	"strings"
)

type ScrapyConfig struct {
	Trans      transform.Transformer
	Catalog    string
	VolName    string
	VolFunc    func(*goquery.Selection) string
	Chapter    string
	NeedCookie bool
	IsJSON     bool
	Api        string
}

var (
	selectMap = map[string]ScrapyConfig{
		"qidian": {
			Trans:   transform.Nop,
			Catalog: ".volume-wrap",
			VolName: ".volume h3",
			Chapter: ".book_name a",
			Api:     "https://book.qidian.com/info/%s/",
			VolFunc: func(doc *goquery.Selection) (volName string) {
				volName = doc.Find(".volume h3").Contents().Not("a").Text()
				if strings.Contains(volName, "·") {
					volName = strings.Split(volName, "·")[0]
				}
				return
			},
		},
		"fanqienovel": {
			Trans:   transform.Nop,
			Catalog: ".page-directory-content",
			VolName: ".volume",
			Chapter: ".chapter-item-title",
		},
		"qimao": {
			Trans:  transform.Nop,
			IsJSON: true,
			Api:    "https://www.qimao.com/api/book/chapter-list?book_id=%s",
		},
	}
	// exclude vols and chapters
	passVols = []string{"第三方", "作品相关", "闲言碎语"}
	// exclude vols but include chapters
	excludeVols = []string{"正文", "VIP"}
)

func GetCatalogFromUrl(ef *formatter.EpubFormat) (err error) {
	var checkConfig ScrapyConfig
	for domain, config := range selectMap {
		if strings.Contains(ef.Catalog.Url, domain) {
			checkConfig = config
			break
		}
	}
	if checkConfig.Api != "" {
		ef.Catalog.Url = fmt.Sprintf(checkConfig.Api, utils.GetNum(ef.Catalog.Url))
	}
	req := utils.GetWithUserAgent(ef.Catalog.Url)
	if ef.Catalog.Cookie != "" {
		req.Header.Set("Cookie", ef.Catalog.Cookie)
	} else if CheckCookie(ef.Catalog.Url) {
		return errors.New("cookie is required")
	}
	if !checkConfig.IsJSON {
		err = GetCatalogByHTML(ef, checkConfig, req)
		return
	}
	err = GetCatalogByJSON(ef, req)
	return
}

func GetCatalogByHTML(ef *formatter.EpubFormat, config ScrapyConfig, req *http.Request) error {
	doc, err := utils.TransDom2Doc(req, config.Trans)
	if err != nil {
		return err
	}
	doc.Find(config.Catalog).Children().Each(func(i int, s *goquery.Selection) {
		// filter vol
		var vol string
		if config.VolFunc != nil {
			vol = config.VolFunc(s)
		} else {
			vol = strings.TrimSpace(s.Find(config.VolName).Contents().First().Text())
		}
		for _, pass := range passVols {
			if strings.Contains(vol, pass) {
				return
			}
		}
		var isExcludeVol bool
		for _, exclude := range excludeVols {
			if strings.Contains(vol, exclude) {
				isExcludeVol = true
				break
			}
		}
		if !isExcludeVol {
			ef.BookConf.Chapters = append(ef.BookConf.Chapters, models.Chapter{Title: vol, IsVol: true})
		}
		s.Find(config.Chapter).Each(func(j int, ss *goquery.Selection) {
			ef.BookConf.Chapters = append(ef.BookConf.Chapters, models.Chapter{Title: ss.Text()})
		})

	})
	return nil
}

func CheckCookie(url string) bool {
	for domain, config := range selectMap {
		if strings.Contains(url, domain) {
			if config.NeedCookie {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

type QiMaoData struct {
	Chapters []struct {
		Title string `json:"title"`
	}
}

type QiMaoJson struct {
	Data QiMaoData
}

func GetCatalogByJSON(ef *formatter.EpubFormat, req *http.Request) (err error) {
	if !utils.CheckUrl(ef.Catalog.Url) {
		return errors.New(stdout.ErrUrl)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var qmChapters QiMaoJson
	err = json.NewDecoder(resp.Body).Decode(&qmChapters)
	if err != nil {
		return
	}

	ef.Chapters = make([]models.Chapter, len(qmChapters.Data.Chapters))
	for i := range qmChapters.Data.Chapters {
		ef.Chapters[i] = models.Chapter{Title: qmChapters.Data.Chapters[i].Title}
	}
	return
}
