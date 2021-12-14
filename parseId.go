package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"strconv"
	"strings"
)

func scrapImages(id string) []string {
	count := 0
	var images []string
	for {
		count++
		imageLink := "https://images.wbstatic.net/c516x688/new/" + id[0:4] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		resp, _ := http.Get(imageLink)
		if resp.StatusCode == 200 {
			images = append(images, imageLink)
		} else {
			return images
		}
	}
}

func scrapId(url string, category string) {
	c := colly.NewCollector()
	newElementsCount := 0
	c.OnHTML(".product-card__wrapper a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		id := strings.Split(link, "/")[2]
		if id != "basket" {
			imagesLinks := scrapImages(id)
			idInt, _ := strconv.Atoi(id)
			writeIdToPostgreSql(idInt, imagesLinks, category)
			newElementsCount++
		}
	})
	addrId := 0
	for {
		addrId++
		linkPage := url + "?sort=popular&page=" + strconv.Itoa(addrId)
		err := c.Visit(linkPage)
		if err != nil {
			fmt.Println(linkPage)
			return
		}
		//println(addrId, newElementsCount)
		if newElementsCount == 0 {
			break
		}
		newElementsCount = 0
	}

}

func scrapIds() {
	categories := readJson()
	for _, v := range categories.Categories {
		scrapId(v.PageUrl, v.Name)
	}
}
