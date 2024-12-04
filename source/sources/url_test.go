package sources

import (
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"testing"
	"time"
)

func TestSimilarity(t *testing.T) {
	raw := "第87章 您已被……移出群聊"
	title := "87.第87章您已被移出群聊"
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
		fmt.Println(k, similar, since)
	}
	fmt.Printf("\nmin algo: %s,time: %s\n", algo, timeMin)
}
