package main

import (
	"github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func scrapItem(id string, category string) int {
	url := "https://wbxcatalog-ru.wildberries.ru/nm-2-card/catalog?locale=ru&nm=" + id
	res, err := http.Get(url)
	if err != nil {
		time.Sleep(time.Second * 3)
		return scrapItem(id, category)
	}
	body, e := ioutil.ReadAll(res.Body)
	if e != nil {
		time.Sleep(time.Second * 3)
		return scrapItem(id, category)
	}
	c, _, _, _ := jsonparser.Get(body, "data", "products")
	_, err = jsonparser.ArrayEach(c, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		price, _, _, error1 := jsonparser.Get(value, "priceU")
		salePrice, _, _, error2 := jsonparser.Get(value, "salePriceU")

		colorsObj, _, _, _ := jsonparser.Get(value, "colors")
		sizeObj, _, _, _ := jsonparser.Get(value, "sizes")
		var colors, sizes []string
		count := 0
		_, err1 := jsonparser.ArrayEach(colorsObj, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			color, _, _, e1 := jsonparser.Get(value, "name")
			if e1 == nil {
				colors = append(colors, string(color))
			}
		})
		if err1 != nil {
			return
		}
		_, err2 := jsonparser.ArrayEach(sizeObj, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			size, _, _, e1 := jsonparser.Get(value, "name")
			if e1 == nil {
				sizes = append(sizes, string(size))
			}
			stockObj, _, _, _ := jsonparser.Get(value, "stocks")
			_, err3 := jsonparser.ArrayEach(stockObj, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				cc, _, _, e1 := jsonparser.Get(value, "qty")
				if e1 == nil {
					ccInt, _ := strconv.Atoi(string(cc))
					count += ccInt
				}
			})
			if err3 != nil {
				return
			}
		})
		if err2 != nil {
			return
		}
		if error1 != nil && error2 != nil {
			return
		}
		var priceF, salePriceF float64
		if error1 != nil {
			salePriceS := string(salePrice)
			salePriceS = salePriceS[:len(salePriceS)-2] + "." + salePriceS[len(salePriceS)-2:]
			salePriceF, _ = strconv.ParseFloat(salePriceS, 8)
			priceF = salePriceF
		} else if error2 != nil {
			priceS := string(price)
			priceS = priceS[:len(priceS)-2] + "." + priceS[len(priceS)-2:]
			priceF, _ = strconv.ParseFloat(priceS, 8)
			salePriceF = priceF
		} else {
			priceS := string(price)
			salePriceS := string(salePrice)
			priceS = priceS[:len(priceS)-2] + "." + priceS[len(priceS)-2:]
			salePriceS = salePriceS[:len(salePriceS)-2] + "." + salePriceS[len(salePriceS)-2:]
			priceF, _ = strconv.ParseFloat(priceS, 8)
			salePriceF, _ = strconv.ParseFloat(salePriceS, 8)
		}
		idInt, _ := strconv.Atoi(id)
		updateItemInfoPostgreSql(idInt, float32(priceF), float32(salePriceF), colors, sizes, count, category)
	})
	if err != nil {
		time.Sleep(time.Second * 3)
		return scrapItem(id, category)
	}
	return 1
}

func scrapItems() {
	var wg sync.WaitGroup
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f20597c3014e4699969af0244a66a6f8@o1108001.ingest.sentry.io/6135375",
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("[4/4] Скрипт парсера товаров запущен!")
	count := 0
	data := GetDbIds()
	for i, v := range data {
		count += 1
		wg.Add(1)
		go func(id int, category string) {
			defer wg.Done()
			scrapItem(strconv.Itoa(id), category)
		}(v.id, v.category)
		if i%50 == 0 {
			wg.Wait()
			if i%50000 == 0 {
				sentry.CaptureMessage("[4/4] Обработано " + strconv.Itoa(count) + " из " + strconv.Itoa(len(data)))
			}

		}
	}
	wg.Wait()
	sentry.CaptureMessage("[4/4] Скрипт парсера товаров завершен!")
}

func mainItems() {
	time.Sleep(500 * time.Second)
	for {
		scrapItems()
	}
}
