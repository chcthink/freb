package htmlx

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

const exprErr = "无法解析 xpath或请求异常\nxpath: %s \nhtml: %s"

func XPathFindStr(node *html.Node, expr string) (string, error) {
	dest := htmlquery.FindOne(node, expr)
	if dest == nil {
		return "", fmt.Errorf(exprErr, expr, fmt.Sprintf("%q", node.Data))
	}
	return strings.TrimSpace(htmlquery.InnerText(dest)), nil
}

func XPathAscEach(nodes []*html.Node, f func(int, *html.Node), isReverse bool) {
	if isReverse {
		total := len(nodes)
		for i := total - 1; i >= 0; i-- {
			f(total-1-i, nodes[i])
		}
	} else {
		for i := range nodes {
			f(i, nodes[i])
		}
	}
}
