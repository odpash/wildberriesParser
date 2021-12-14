package main

type Categories struct {
	Categories []Category
}

type Category struct {
	Name    string
	PageUrl string
}

func main() {
	// scrapCategories() // used to get categories name and ids
	// scrapIds()  // used to get ids and images
	scrapItems() // used to get item info (such as price, sale price, color, size, count)
}
