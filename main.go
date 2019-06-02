package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	ListURL    = "https://www.dcard.tw/_api/forums/sex/posts?popular=false"
	ContentURL = "https://www.dcard.tw/_api/posts/"
	BoardURL   = "http://www.dcard.tw/f/sex/p/"
	TgbotToken = "755108266:AAFFw6H5k9LIMOcKlrp7Au622OL46JnGzec"
	ChatID     = -1001494629371
	limitNum   = 30
	SleepNum   = 30
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
	var LatestID int32 = 0

	bot, err := tgbotapi.NewBotAPI(TgbotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	for {
		log.Println(time.Now())
		log.Println("LatestID=", LatestID)

		var SendCount int = 0
		SexArticles, err := getLatestList(LatestID)
		if err != nil {
			log.Println(err)
			return
		}
		//fmt.Printf("%#v\n", SexArticles)

		for _, value := range SexArticles {
			var MsgContent string = ""

			fmt.Printf("文章連結:%s \n\n", BoardURL+strconv.FormatInt(int64(value.ID), 10))
			if value.ID > LatestID {
				LatestID = value.ID
			}
			fmt.Println("id=", value.ID)
			fmt.Printf("Title= %s \n", value.Title)

			MsgContent += fmt.Sprintf("%s \n", value.Title)
			MsgContent += fmt.Sprintf("文章連結:%s \n", BoardURL+strconv.FormatInt(int64(value.ID), 10))
			fmt.Println("-------Media-------")
			for _, MediaValue := range value.Media {
				fmt.Printf("PicURL= %s \n", MediaValue.PicURL)
				//MsgContent += fmt.Sprintf("PicURL= %s \n", MediaValue.PicURL)

				PhotoMsg := tgbotapi.NewPhotoShare(ChatID, MediaValue.PicURL)
				bot.Send(PhotoMsg)
			}

			fmt.Println("-------MediaMeta-------")
			for _, MediaMetaValue := range value.MediaMeta {
				fmt.Printf("id= %s \n", MediaMetaValue.(map[string]interface{})["id"].(string))
				fmt.Printf("url= %s \n", MediaMetaValue.(map[string]interface{})["url"].(string))
				fmt.Printf("type= %s \n\n", MediaMetaValue.(map[string]interface{})["type"].(string))
			}
			//fmt.Printf("MediaMeta=%v \n", value.MediaMeta)

			TextMsg := tgbotapi.NewMessage(ChatID, MsgContent)
			bot.Send(TextMsg)
			SendCount++
			fmt.Printf("%d new articles have been sent------------------------ \n", SendCount)
		}
		time.Sleep(SleepNum * time.Second)
	}
}

//get latest article list
func getLatestList(afterID int32) ([]ArticleInfo, error) {
	params := make(map[string]string)
	params["limit"] = strconv.FormatInt(int64(limitNum), 10)
	if afterID != 0 {
		params["after"] = strconv.FormatInt(int64(afterID), 10)
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
