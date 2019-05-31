package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ListURL    = "https://www.dcard.tw/_api/forums/sex/posts?popular=false"
	ContentURL = "https://www.dcard.tw/_api/posts/"
	BoardURL   = "http://www.dcard.tw/f/sex/p/"
)

type Bird struct {
	Species     string
	Description string
}

type ArticleInfo struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Media     []Media `json:"media"`
	MediaMeta []interface{}
}

type Media struct {
	URL string
}

// type MediaMeta struct {
// 	ID            string `json:"id"`
// 	URL           string `json:"url"`
// 	type string `json:"type"`
// }

func main() {
	resp, err := http.Get(ListURL + "&limit=1")
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

	var SexArticles []ArticleInfo
	if err := json.Unmarshal([]byte(sitemap), &SexArticles); err != nil {
		fmt.Println("error:", err)
	}

	//fmt.Printf("%#v\n", SexArticles)
	for _, value := range SexArticles {
		fmt.Println("id=", value.ID)
		fmt.Printf("Title= %s \n", value.Title)
		fmt.Println("-------Media-------")
		for _, MediaValue := range value.Media {
			fmt.Printf("URL= %v \n", MediaValue)
		}
		fmt.Println("-------MediaMeta-------")
		for _, MediaMetaValue := range value.MediaMeta {
			fmt.Printf("id= %s \n", MediaMetaValue.(map[string]interface{})["id"].(string))
			fmt.Printf("url= %s \n", MediaMetaValue.(map[string]interface{})["url"].(string))
			fmt.Printf("type= %s \n\n", MediaMetaValue.(map[string]interface{})["type"].(string))
		}
		//fmt.Printf("MediaMeta=%v \n", value.MediaMeta)
	}

	//fmt.Println(SexArticles[0].MediaMeta[0].(map[string]interface{})["id"].(string))
	//fmt.Println(SexArticles.MediaMeta.(map[string]interface{})["normalizedUrl"].(string))

}
