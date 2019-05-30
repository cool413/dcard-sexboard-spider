package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ArticleInfo struct {
	Id    int
	Title string
	Media []Media
}

type Media struct {
	Url string
}

func main() {
	resp, err := http.Get("https://www.dcard.tw/_api/forums/sex/posts?limit=10")
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var articlesInfo []ArticleInfo
	json.Unmarshal([]byte(sitemap), &articlesInfo)
	fmt.Printf("articlesInfo :%+v", articlesInfo)
	fmt.Printf("\n ArticleInfo[0]: %s", articlesInfo[0].Media[0].Url)
}
