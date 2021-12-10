package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)


func main() {
	//scrapId("https://www.wildberries.ru/catalog/zhenshchinam/bolshie-razmery/bele")
	scrapCategories()
}


func scrapCategoriesCycle(c []byte) {
	jsonparser.ArrayEach(c, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		cNew, _, _, childError := jsonparser.Get(value, "childs")
		if childError != nil {
			name, _, _, _ := jsonparser.Get(value, "name")
			pageUrl, _, _, _ := jsonparser.Get(value, "pageUrl")
			if strings.Contains(string(pageUrl), "catalog") {
				fmt.Println("Ребенка нет:", string(name), string(pageUrl))
			}
		} else {
			scrapCategoriesCycle(cNew)
		}
})
}

func scrapCategories() {
	url := "https://www.wildberries.ru/gettopmenuinner?lang=ru"
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	c, _, _, _ := jsonparser.Get(body, "value", "menu")
	scrapCategoriesCycle(c)

}


func scrapId(url string) {
	c := colly.NewCollector()
	newElementsCount := 0
	c.OnHTML(".product-card__wrapper a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		id := strings.Split(link, "/")[2]
		if id != "basket" {
			println(id)
			newElementsCount++
		}
	})

	addrId := 0
	for {
		addrId++
		linkPage := url + "?sort=popular&page=" + strconv.Itoa(addrId)
		c.Visit(linkPage)
		println(addrId, newElementsCount)
		if newElementsCount == 0 {
			break
		}
		newElementsCount = 0
	}

}