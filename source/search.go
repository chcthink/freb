package source

import (
	"errors"
	"fmt"
	"freb/models"
	"freb/utils/htmlx"
	"github.com/antchfx/htmlquery"
	"net/url"
	"strings"
)

const (
	google         = `https://www.google.com/search?q=%s`
	site           = `site:%s`
	definiteSearch = `"%s" `
	bracket        = `(%s)`
	or             = ` OR `
)

const (
	noCatchFoundErr = "config.toml.book_catch数量为 0"
	unfoundedErr    = `未匹配%s`
)

// return https://www.google.com/search?q="abc"+(site:a.com+OR+site:b.com)
func genUrl(search string, catches map[string]*models.BookCatch) string {
	// "abc"+
	search = fmt.Sprintf(definiteSearch, search)

	var sites []string
	for _, catch := range catches {
		sites = append(sites, fmt.Sprintf(site, catch.SearchMatch))
	}
	// "abc"+(site:a.com+OR+site:b.com)
	search = strings.Join([]string{search, fmt.Sprintf(bracket, strings.Join(sites, or))}, "")
	search = url.QueryEscape(search)

	return fmt.Sprintf(google, search)
}

func Search(search string, catches map[string]*models.BookCatch) (bookUrl string, err error) {
	if len(catches) == 0 {
		return "", errors.New(noCatchFoundErr)
	}
	req := htmlx.GetWithUserAgent(genUrl(search, catches))
	doc, err := htmlx.TransDom2Doc(req)
	if err != nil {
		return
	}
	nodes := htmlquery.Find(doc, `//div[@class='dURPMd']//a[@jsname]/@href`)
	for _, node := range nodes {
		bookUrl = strings.TrimSpace(htmlquery.InnerText(node))
		for _, catch := range catches {
			if strings.Contains(bookUrl, catch.SearchMatch) {
				return
			}
		}
	}
	return "", fmt.Errorf(unfoundedErr, search)
}
