package source

import (
	"fmt"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/utils/htmlx"
	"freb/utils/reg"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

func GetCatalogFromUrl(ef *formatter.EpubFormat) (err error) {
	var checkConfig *models.InfoSelector
	for domain, conf := range config.Cfg.InfoSelector {
		if strings.Contains(ef.Catalog, domain) {
			checkConfig = conf
			break
		}
	}
	if checkConfig == nil {
		return fmt.Errorf("未匹配目录来源: %s", ef.Catalog)
	}
	if checkConfig.Api != "" {
		ef.Catalog = fmt.Sprintf(checkConfig.Api, reg.GetNum(ef.Catalog))
	}
	req := htmlx.GetWithUserAgent(ef.Catalog)
	if !checkConfig.IsJSON {
		err = GetCatalogByHTML(ef, checkConfig, req)
		return
	}
	err = GetCatalogByJSON(ef, checkConfig, req)
	return
}

func GetCatalogByHTML(ef *formatter.EpubFormat, conf *models.InfoSelector, req *http.Request) error {
	doc, err := htmlx.TransDom2Doc(req)
	if err != nil {
		return err
	}
	nodes := htmlquery.Find(doc, conf.Catalog)
	var isContinue bool
VOL:
	for i, node := range nodes {
		var vol string
		vol, err = htmlx.XPathFindStr(node, conf.VolName)
		if err != nil {
			return err
		}
		// avoid filter exclude vol name when normal vol name  follow exclude vol name
		if i == 0 {
			var secVol string
			secVol, err = htmlx.XPathFindStr(nodes[len(nodes)-1], conf.VolName)
			if err != nil {
				return err
			}
			if secVol != vol {
				isContinue = true
			}
		}
		for _, pass := range conf.PassVols {
			if strings.Contains(vol, pass) {
				continue VOL
			}
		}
		var isExcludeVol bool
		if isContinue {
			isExcludeVol = true
		} else {
			for _, exclude := range conf.ExcludeVols {
				if strings.Contains(vol, exclude) {
					isExcludeVol = true
					break
				}
			}
		}

		if !isExcludeVol {
			ef.Sections = append(ef.Sections, models.Section{Title: vol, IsVol: true})
			conf.ExcludeVols = append(conf.ExcludeVols, vol)
		}
		for _, subNode := range htmlquery.Find(node, conf.Chapter) {
			ef.Sections = append(ef.Sections, models.Section{Title: htmlquery.InnerText(subNode)})
		}
	}
	return nil
}

func GetCatalogByJSON(ef *formatter.EpubFormat, conf *models.InfoSelector, req *http.Request) (err error) {
	var rest gjson.Result
	rest, err = htmlx.TransDom2JSON(req)
	if err != nil {
		return
	}
	if conf.Catalog != "" {
		nodes := rest.Get(conf.Catalog).Array()
	VOL:
		for _, node := range nodes {
			vol := strings.TrimSpace(node.Get(conf.VolName).String())
			for _, pass := range conf.PassVols {
				if strings.Contains(vol, pass) {
					continue VOL
				}
			}
			var isExcludeVol bool
			for _, exclude := range conf.ExcludeVols {
				if strings.Contains(vol, exclude) {
					isExcludeVol = true
					break
				}
			}
			if !isExcludeVol {
				ef.Sections = append(ef.Sections, models.Section{Title: vol, IsVol: true})
				conf.ExcludeVols = append(conf.ExcludeVols, vol)
			}
			for _, subNode := range node.Get(conf.Chapter).Array() {
				ef.Sections = append(ef.Sections, models.Section{Title: subNode.String()})
			}
		}
		return nil

	}
	chpts := rest.Get(conf.Chapter).Array()
	ef.Sections = make([]models.Section, len(chpts))
	for i, chpt := range chpts {
		ef.Sections[i] = models.Section{Title: chpt.String()}
	}
	return
}
