package sources

import (
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
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
