package htmlx

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func XPathFindStr(node *html.Node, expr string) string {
	return htmlquery.InnerText(htmlquery.FindOne(node, expr))
}

func XpPathAscEach(nodes []*html.Node, f func(int, *html.Node), isReverse bool) {
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

func iterateChildNodes(node *html.Node) (dest []*html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		dest = append(dest, child)
	}
	return
}
