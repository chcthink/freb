package source

import (
	"fmt"
	"freb/formatter"
	"freb/models"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"testing"
)

func TestFanqie(t *testing.T) {
	var htmlFanQie = `<div class="page-body">
    <div class="page-abstract-header">
        <h2>作品简介</h2>
    </div>
    <div class="page-abstract-content">
        <p>（不后宫，不套路，不无敌，不系统，不无脑，不爽文，介意者慎入。）
            当我以为这只是寻常的一天时，却发现自己被捉到了终焉之地。
            当我以为只需要不断的参加死亡游戏就可以逃脱时，却发现众人开始觉醒超自然之力。
            当我以为这里是「造神之地」时，一切却又奔着湮灭走去。</p>
    </div>
    <div class="page-directory-header">
        <h3><span>目录</span><span class="directory-dot"></span>1360章</h3>
    </div>
    <div class="page-directory-content">
        <div class="">
            <div class="volume volume_first">第一卷：我听到了你们<span class="volume-dot"></span>共91章</div>
            <div class="chapter">
                <div class="chapter-item"><a href="/reader/7173216089122439711" class="chapter-item-title"
                        target="_blank">第1章 空屋</a></div>
                <div class="chapter-item"><a href="/reader/7173217408294453791" class="chapter-item-title"
                        target="_blank">第2章 说谎</a></div>
                <div class="chapter-item"><a href="/reader/7173615024101917184" class="chapter-item-title"
                        target="_blank">第3章 有技术的人</a></div>
            </div>
        </div>
    
        <div class="">
            <div class="volume">第十卷：这一切的终焉<span class="volume-dot"></span>共255章</div>
            <div class="chapter">
                <div class="chapter-item"><a href="/reader/7431902608589079102" class="chapter-item-title"
                        target="_blank">第1357章 「入梦」</a><span class="chapter-item-lock muyeicon-lock"></span></div>
                <div class="chapter-item"><a href="/reader/7431903018192224830" class="chapter-item-title"
                        target="_blank">第1358章 「十日终焉」</a><span class="chapter-item-lock muyeicon-lock"></span></div>
                <div class="chapter-item"><a href="/reader/7431941887952421438" class="chapter-item-title"
                        target="_blank">终：完结感言</a><span class="chapter-item-lock muyeicon-lock"></span></div>
            </div>
        </div>
    </div>
</div>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlFanQie))
	doc.Find(".page-directory-content").Children().Each(func(i int, s *goquery.Selection) {
		vol := strings.TrimSpace(s.Find(".volume").Text())
		if strings.Contains(vol, "·") {
			volSub := strings.Split(vol, "·")
			vol = volSub[len(volSub)-2]
		}
		fmt.Println(s.Find(".volume").Contents().First().Text())
		s.Find(".chapter-item-title").Each(func(j int, ss *goquery.Selection) {
			fmt.Println(ss.Text())
		})
	})
}

func TestQiDian(t *testing.T) {
	var ef formatter.EpubFormat
	ef.Book = &models.Book{
		Catalog: models.UrlWithCookie{
			Url:    "https://www.qidian.com/book/1035420986/",
			Cookie: "",
		},
	}
	err := GetCatalogByUrl(&ef)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ef.Chapters)
}
