package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
)

type Categories struct {
	Categories []Category
}

type Category struct {
	Name    string
	PageUrl string
}

func main() {
	//scrapId("https://www.wildberries.ru/catalog/zhenshchinam/bolshie-razmery/bele")
	//scrapCategories()
	//scrapItem("7851246")
	//scrapId("https://www.wildberries.ru/catalog/budushchie-mamy/aksessuary")
	//fmt.Println(scrapImages("9510116"))
}

func scrapItem(id string) {
	url := "https://wbxcatalog-ru.wildberries.ru/nm-2-card/catalog?locale=ru&nm=" + id
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	c, _, _, _ := jsonparser.Get(body, "data", "products")
	_, err := jsonparser.ArrayEach(c, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, _, _, _ := jsonparser.Get(value, "id")
		price, _, _, _ := jsonparser.Get(value, "priceU")
		salePrice, _, _, _ := jsonparser.Get(value, "salePriceU")
		fmt.Println(string(id), string(price), string(salePrice))
	})
	if err != nil {
		return
	}
}
