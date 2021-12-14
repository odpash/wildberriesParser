package main

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"time"
)

type Item struct {
	imageLinks, colors, sizes []string
	prices, salePrices        []float32
	infoDate                  []string
	count                     int
	category                  string
}

type Id struct {
	id       int
	category string
}

const categoryFilename = "category.json"
const connStr = "user=postgres password=12345 dbname=wildberries sslmode=disable"

func writeIdToPostgreSql(id int, images []string, category string) {
	status, _ := getAllByIdPostgreSql(id)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if status {
		db.Exec("update items set imagelinks = $1 where id = $2",
			pq.Array(images), id)
	} else {
		db.Exec("insert into items (id, imagelinks, count, category) values ($1, $2, 0, $3)",
			id, pq.Array(images), category)
	}
}

func updateItemInfoPostgreSql(id int, priceF float32, salePriceF float32, colors []string, sizes []string, count int, category string) {
	status, item := getAllByIdPostgreSql(id)
	if status == false {
		return
	}
	dt := time.Now()
	item.infoDate = append(item.infoDate, dt.Format("01-02-2006"))
	item.count = count
	item.prices = append(item.prices, priceF)
	item.salePrices = append(item.salePrices, salePriceF)
	item.colors = colors
	item.category = category
	item.sizes = sizes
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec("update items set infodate = $1, count = $2, prices = $3, saleprices = $4, colors = $5, sizes = $6, category = $7 where id = $8",
		pq.Array(item.infoDate), item.count, pq.Array(item.prices), pq.Array(item.salePrices), pq.Array(item.colors), pq.Array(item.sizes), item.category, id)
}

func getAllByIdPostgreSql(id int) (bool, Item) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, e := db.Query("Select imagelinks, colors, sizes, saleprices, prices, infodate, count, category from items where id = $1", id)
	item := Item{}
	if e != nil {
		return false, item
	}
	for res.Next() {
		res.Scan(pq.Array(&item.imageLinks), pq.Array(&item.colors), pq.Array(&item.sizes), pq.Array(&item.salePrices), pq.Array(&item.prices), pq.Array(&item.infoDate), &item.count, &item.category)
		return true, item
	}
	return false, item
}

func writeJson(info Categories) {
	rawDataOut, err := json.MarshalIndent(&info, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}
	err = ioutil.WriteFile(categoryFilename, rawDataOut, 0)
	if err != nil {
		log.Fatal("Cannot write updated settings file:", err)
	}
}

func readJson() Categories {
	var newCategories Categories
	rawDataIn, err := ioutil.ReadFile(categoryFilename)
	if err != nil {
		log.Fatal("Cannot load settings:", err)
	}
	err = json.Unmarshal(rawDataIn, &newCategories)
	if err != nil {
		log.Fatal("Invalid settings format:", err)
	}
	return newCategories
}

func getDbIds() []Id {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, _ := db.Query("Select id, category from items")
	var ids []Id
	for res.Next() {
		id := Id{}
		res.Scan(&id.id, &id.category)
		ids = append(ids, id)
	}
	return ids
}
