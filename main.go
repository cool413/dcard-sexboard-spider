package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	ListURL    = "https://www.dcard.tw/_api/forums/sex/posts?popular=false"
	ContentURL = "https://www.dcard.tw/_api/posts/"
	BoardURL   = "http://www.dcard.tw/f/sex/p/"
)

type ArticleInfo struct {
	ID        int32   `json:"id"`
	Title     string  `json:"title"`
	Media     []Media `json:"media"`
	MediaMeta []interface{}
}

type Media struct {
	PicURL string `json:"url"`
}

func main() {
	var NowID int32 = 0

	for {
		log.Println(time.Now())
		log.Println("NowId=", NowID)

		SexArticles, err := getLatestList(NowID)
		if err != nil {
			log.Println(err)
			return
		}
		//fmt.Printf("%#v\n", SexArticles)

		for _, value := range SexArticles {
			fmt.Println("id=", value.ID)
			fmt.Printf("Title= %s \n", value.Title)
			fmt.Println("-------Media-------")
			for _, MediaValue := range value.Media {
				fmt.Printf("PicURL= %s \n", MediaValue.PicURL)
			}
			fmt.Println("-------MediaMeta-------")
			for _, MediaMetaValue := range value.MediaMeta {
				fmt.Printf("id= %s \n", MediaMetaValue.(map[string]interface{})["id"].(string))
				fmt.Printf("url= %s \n", MediaMetaValue.(map[string]interface{})["url"].(string))
				fmt.Printf("type= %s \n\n", MediaMetaValue.(map[string]interface{})["type"].(string))
			}
			//fmt.Printf("MediaMeta=%v \n", value.MediaMeta)

			if value.ID > NowID {
				NowID = value.ID
			}
		}

		time.Sleep(10 * time.Second)
	}

}

func getLatestList(LatestID int32) ([]ArticleInfo, error) {
	params := make(map[string]string)
	params["limit"] = "10"
	if LatestID != 0 {
		params["after"] = strconv.FormatInt(int64(LatestID), 10)
	}
	resp, err := Get(ListURL, params, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//fmt.Println(string(sitemap))

	var Articles []ArticleInfo
	if err := json.Unmarshal([]byte(sitemap), &Articles); err != nil {
		log.Println(string(sitemap))
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}

	return Articles, nil
}

//Get http get method
func Get(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
		//nil只能賦值給指標
		return nil, fmt.Errorf("new request is fail:%s", url)
	}

	q := req.URL.Query()
	if params != nil {
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	log.Printf("Go GET URL : %s \n", req.URL.String())

	return client.Do(req)
}
