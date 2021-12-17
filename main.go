package main

import (
	"database/sql"
	"github.com/lib/pq"
)

type Categories struct {
	Categories []Category
}

type Category struct {
	Name    string
	PageUrl string
}

func ab() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	id := 5
	var images []string
	category := "test0"
	db.Exec("insert into items (id, imagelinks, count, category) values ($1, $2, 0, $3)",
		id, pq.Array(images), category)
	db.Exec("update items set imagelinks = $2, category = $3 where id = $1",
		id, pq.Array(images), category)
} //  ON conflict ($1) do
func main() {
	//scrapCategories() // used to get categories name and ids
	//scrapIds()  // used to get ids and images
	scrapItems() // used to get item info (such as price, sale price, color, size, count)
	//ab()
}
