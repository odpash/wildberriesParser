package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const categoryFilename = "category.json"

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
