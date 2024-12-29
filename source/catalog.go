package source

import (
	"errors"
	"freb/formatter"
	"freb/models"
	"freb/utils"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/transform"
	"strings"
)

type ScrapyConfig struct {
	Trans      transform.Transformer
	Catalog    string
	VolName    string
	VolFunc    func(*goquery.Selection) string
	Chapter    string
	NeedCookie bool
}

var (
	selectMap = map[string]ScrapyConfig{
		"qidian": {
			Trans:      transform.Nop,
			Catalog:    "#allCatalog",
			VolName:    ".volume-name",
			Chapter:    ".chapter-name",
			NeedCookie: true,
		},
		"fanqienovel": {
			Trans:   transform.Nop,
			Catalog: ".page-directory-content",
			VolName: ".volume",
			Chapter: ".chapter-item-title",
		},
	}
	// exclude vols and chapters
	passVols = []string{"第三方", "作品相关", "闲言碎语"}
	// exclude vols but include chapters
	excludeVols = []string{"正文", "VIP"}
)

func GetCatalogByUrl(ef *formatter.EpubFormat) error {
	req := utils.GetWithUserAgent(ef.Catalog.Url)
	if ef.Catalog.Cookie != "" {
		req.Header.Set("Cookie", ef.Catalog.Cookie)
	} else if CheckCookie(ef.Catalog.Url) {
		return errors.New("cookie is required")
	}
	for domain, config := range selectMap {
		if strings.Contains(ef.Catalog.Url, domain) {
			doc, err := utils.TransDom2Doc(ef.Catalog.Url, req, config.Trans)
			if err != nil {
				return err
			}

			doc.Find(config.Catalog).Children().Each(func(i int, s *goquery.Selection) {
				// filter vol
				vol := strings.TrimSpace(s.Find(config.VolName).Contents().First().Text())
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
	}
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
