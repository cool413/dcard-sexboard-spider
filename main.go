package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	ListURL    = "https://www.dcard.tw/_api/forums/sex/posts?popular=false"
	ContentURL = "https://www.dcard.tw/_api/posts/"
	LinkURL    = "http://www.dcard.tw/f/sex/p/"
	CommentURL = "http://dcard.tw/_api/posts/%d/comments"

	TgbotToken = "755108266:AAFFw6H5k9LIMOcKlrp7Au622OL46JnGzec"
	ChatID     = -1001494629371
	limitNum   = 10
	SleepNum   = 30
)

type ArticleInfo struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Media     []Media `json:"media"`
	MediaMeta []interface{}
}

type CommentInfo struct {
	Content   int64  `json:"content"`
	Title     string `json:"title"`
	MediaMeta []interface{}
}

type Media struct {
	PicURL string `json:"url"`
}

func main() {
	var LatestID int64 = 0

	bot, err := tgbotapi.NewBotAPI(TgbotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	for {
		fmt.Println(time.Now())
		fmt.Printf("------------------------ LatestID:%d ------------------------", LatestID)
		var SendCount int = 0

		SexArticles, err := getLatestList(LatestID)
		if err != nil {
			log.Println(err)
			return
		}
		//fmt.Printf("%#v\n", SexArticles)
		for _, value := range SexArticles {
			var MsgContent string = ""
			var ArticleID = value.ID

			if ArticleID > LatestID {
				LatestID = ArticleID
			}

			for _, MediaMetaValue := range value.MediaMeta {
				PicURL := MediaMetaValue.(map[string]interface{})["normalizedUrl"].(string)
				PicType := MediaMetaValue.(map[string]interface{})["type"].(string)

				if strings.Index(strings.ToLower(PicType), "thumbnail") == -1 {
					PhotoMsg := tgbotapi.NewPhotoShare(ChatID, PicURL)
					bot.Send(PhotoMsg)
				}

				// fmt.Printf("url= %s \n", MediaMetaValue.(map[string]interface{})["normalizedUrl"].(string))
				// fmt.Printf("type= %s \n\n", MediaMetaValue.(map[string]interface{})["type"].(string))
			}
			//fmt.Printf("MediaMeta=%v \n", value.MediaMeta)

			//get Comments
			// Comments, err := getComments(ArticleID)
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }
			// fmt.Printf("%#v\n", Comments)

			MsgContent += fmt.Sprintf("%s \n", value.Title)
			MsgContent += fmt.Sprintf("文章連結:%s \n", LinkURL+strconv.FormatInt(int64(ArticleID), 10))
			TextMsg := tgbotapi.NewMessage(ChatID, MsgContent)
			bot.Send(TextMsg)
			SendCount++

			fmt.Printf("%s \n", MsgContent)
			fmt.Printf("------------------------ %d new articles have been sent ------------------------ \n", SendCount)
		}
		time.Sleep(SleepNum * time.Second)
	}
}

//get latest article list
func getLatestList(afterID int64) ([]ArticleInfo, error) {
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

//get Comments
func getComments(articleID int64) ([]CommentInfo, error) {
	resp, err := Get(fmt.Sprintf(CommentURL, articleID), nil, nil)
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

	var Comments []CommentInfo
	if err := json.Unmarshal([]byte(sitemap), &Comments); err != nil {
		log.Println(string(sitemap))
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}

	return Comments, nil
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
