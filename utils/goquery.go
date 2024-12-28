package utils

import (
	"github.com/PuerkitoBio/goquery"
)

func AscEach(s *goquery.Selection, f func(int, *goquery.Selection), isReverse bool) {
	if isReverse {
		nodes := s.Nodes
		total := len(nodes)
		for i := total - 1; i >= 0; i-- {
			node := nodes[i]
			doc := goquery.NewDocumentFromNode(node)
			f(total-1-i, doc.Selection)
		}
	} else {
		s.Each(f)
	}
}
