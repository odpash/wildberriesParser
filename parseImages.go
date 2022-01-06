package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

func scrapImage(id string, category string) int {
	count := 0
	var images []string
	for {
		count++
		imageLink := ""
		if len(id) == 8 {
			imageLink = "https://images.wbstatic.net/c516x688/new/" + id[0:4] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		} else if len(id) == 7 {
			imageLink = "https://images.wbstatic.net/c516x688/new/" + id[0:3] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		}

		resp, e := http.Get(imageLink)
		if e != nil {
			strId, _ := strconv.Atoi(id)
			writeIdToPostgreSql(strId, images, category) // заменить запись
			return 1
		}
		if resp.StatusCode == 200 {
			images = append(images, imageLink)
		} else {
			strId, _ := strconv.Atoi(id)
			writeIdToPostgreSql(strId, images, category) // заменить запись
			return 1
		}
	}
}

func scrapImages() {
	var wg sync.WaitGroup
	for i, v := range getDbIds() {
		wg.Add(1)
		go func(id int, category string) {
			defer wg.Done()
			scrapImage(strconv.Itoa(id), category)
		}(v.id, v.category)
		if i%50 == 0 {
			fmt.Println(i)
			wg.Wait()
		}
	}
	wg.Wait()
}
