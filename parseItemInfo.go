package main

import (
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strconv"
)

func scrapItem(id string, category string) {
	url := "https://wbxcatalog-ru.wildberries.ru/nm-2-card/catalog?locale=ru&nm=" + id
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	c, _, _, _ := jsonparser.Get(body, "data", "products")
	_, err := jsonparser.ArrayEach(c, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
		return
	}
}

func scrapItems() {
	for _, v := range getDbIds() {
		scrapItem(strconv.Itoa(v.id), v.category)
		// time.Sleep(time.Second)
	}
}
