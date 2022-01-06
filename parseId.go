package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gocolly/colly"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func scrapId(url string, category string, pageNum int, readOnly bool) int {
	c := colly.NewCollector()
	pagesCountInt := 0
	c.OnHTML(".goods-count span", func(e *colly.HTMLElement) {
		itemsCount := ""
		for i := 0; i < len(e.Text); i++ {
			if strings.ContainsAny(string(e.Text[i]), "0123456789") {
				itemsCount += string(e.Text[i])
			}
		}
		pagesCountInt, _ = strconv.Atoi(itemsCount)
		pagesCountInt = pagesCountInt/100 + 1
	})

	c.OnHTML(".product-card__wrapper a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		id := strings.Split(link, "/")[2]
		if id != "basket" && !readOnly {
			//imagesLinks := scrapImages(id)
			var imagesLinks []string
			idInt, _ := strconv.Atoi(id)
			writeIdToPostgreSql(idInt, imagesLinks, category)
		}
	})

	for {
		linkPage := url + "?sort=popular&page=" + strconv.Itoa(pageNum)
		err := c.Visit(linkPage)
		if err != nil {
			divizion := rand.Intn(1000)
			fmt.Println("Request error. Sleep ", divizion, " microsecs and continue")
			time.Sleep(time.Millisecond * time.Duration(divizion))
			scrapId(url, category, pageNum, false)
		}
		return pagesCountInt

		//println(addrId, newElementsCount)

	}

}

type arr struct {
	pagesCount int
}

func scrapIds() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f20597c3014e4699969af0244a66a6f8@o1108001.ingest.sentry.io/6135375",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("[2/4] Скрипт парсера ID запущен!")
	var wg sync.WaitGroup
	categories := readJson()
	summaryCount := 0
	start_first := time.Now()
	var sp []arr
	for _, v := range categories.Categories {
		pagesCount := scrapId(v.PageUrl, v.Name, 1, true)
		summaryCount += pagesCount
		sp = append(sp, arr{pagesCount: pagesCount})
	}
	sentry.CaptureMessage("[2/4] Все страницы были получены за: " + time.Since(start_first).String() + " Количество:" + strconv.Itoa(summaryCount))
	nowCount := 0
	for x, v := range categories.Categories {
		start := time.Now()
		pagesCount := sp[x].pagesCount
		for i := 1; i <= pagesCount; i++ {
			wg.Add(1)
			go func(v Category, i int, readOnly bool) {
				defer wg.Done()
				scrapId(v.PageUrl, v.Name, i, false)
			}(v, i, false)
			if i%50 == 0 {
				wg.Wait()
			}
		}
		if pagesCount != 0 {
			wg.Wait()
		}
		nowCount += pagesCount
		sentry.CaptureMessage("[2/4] Обработка группы " + strconv.Itoa(x+1) + "/" + strconv.Itoa(len(sp)) +
			" завершена за " + time.Since(start).String() + "!\nПолучено данных: " + strconv.Itoa(pagesCount) +
			".\nСуммарно обработано " + strconv.Itoa(nowCount) + "/" + strconv.Itoa(summaryCount) + ".\n" +
			"Осталось обработать " + strconv.Itoa(summaryCount-nowCount) + " страниц.\nВремя в работе: " + time.Since(start_first).String())
	}
	sentry.CaptureMessage("[2/4] Парсер ID завершил работу за " + time.Since(start_first).String())
}

func main() {
	for {
		scrapIds() // How to start? | Easy! | go run parseId.go db.go Interfaces.go
	}
}
