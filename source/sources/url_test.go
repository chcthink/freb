package sources

import (
	"fmt"
	"freb/config"
	"freb/formatter"
	"freb/models"
	"freb/utils"
	"freb/utils/reg"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
	"testing"
	"time"
)

func TestSimilarity(t *testing.T) {
	// raw := "第87章 您已被……移出群聊"
	// title := "87.第87章您已被移出群聊"
	// raw := "第101章 抄家"
	// title := "第九十六章 抄家"
	// raw := "第5章 序"
	// title := "序"
	raw := "第二百章十九  难题"
	title := "第233章 难题"
	algorithms := map[string]strutil.StringMetric{
		"newHamming":         metrics.NewHamming(),
		"Levenshtein":        metrics.NewLevenshtein(),
		"Jaccard":            metrics.NewJaccard(),
		"Jaro":               metrics.NewJaro(),
		"JaroWinkler":        metrics.NewJaroWinkler(),
		"OverlapCoefficient": metrics.NewOverlapCoefficient(),
		"SorensenDice":       metrics.NewSorensenDice(),
		"SmithWatermanGotoh": metrics.NewSmithWatermanGotoh(),
	}
	timeMin := 10 * time.Second
	algo := ""
	similar := 0.0
	similarAlgo := ""
	similarTime := 10 * time.Second
	mostSimilar := 0.0
	for k, v := range algorithms {
		tb := time.Now()
		t.Run(k, func(t *testing.T) {
			similar = strutil.Similarity(raw, title, v)
		})
		since := time.Since(tb)
		if since < timeMin {
			algo = k
			timeMin = since
		}
		if similar > mostSimilar {
			mostSimilar = similar
			similarAlgo = k
			similarTime = since
		}
		fmt.Printf("algo: %s\nsimilar: %.2f\ntime: %s\n", k, similar, since)
	}
	fmt.Printf("\nmin algo: %s,time: %s\n", algo, timeMin)
	fmt.Printf("most similar algo: %s,time: %s, similar: %f\n", similarAlgo, similarTime, mostSimilar)

}

var htmlstr = `<div class="txtnav">
            <h1 class="hide720">title</h1>
            <div class="txtinfo hide720"><span>2025-01-06</span></div>
            <div id="txtright">
                <script>loadAdv(2, 0);</script>
            </div>
            第1章 title
			<p>冲云破雾</p>
			<p>为小失大</p><div class="contentadv"><script>loadAdv(7,3);</script></div>远年近岁<p></p>
			<p>汉票签处</p>
			<p>丁宁周至</p>

            <div class="bottom-ad">
                <script>loadAdv(3, 0);</script>
            </div>
        </div>`

func TestContent(t *testing.T) {

	config.InitConfig()
	var bookCatch *models.BookCatch
	for domain, catch := range config.Cfg.BookCatch {
		if strings.Contains("https://69shuba.cx/book/85745.htm", domain) {
			bookCatch = catch
		}
	}
	doc, _ := htmlquery.Parse(strings.NewReader(htmlstr))

	title := "第1章 title"
	content := ""
	var ef formatter.EpubFormat
	node := htmlquery.Find(doc, bookCatch.Content.Selector)
	var f func(int, *html.Node)
	f = func(index int, n *html.Node) {
		if n.Type == html.TextNode {
			raw := strings.TrimSpace(n.Data)

			if raw == "" || len([]rune(raw)) == 1 {
				return
			}
			// filter title in content
			if utils.SimilarStr(raw, title) && index <= 10 {
				return
			}
			if strings.Contains(raw, "本章完") {
				return
			}
			raw = reg.RemoveContentFromCfg(raw)
			if raw == "" {
				return
			}

			content += ef.GenLine(raw)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(index, c)
			}
		}
	}
	for index, n := range node {
		f(index, n)
	}
	fmt.Println(content)

}
