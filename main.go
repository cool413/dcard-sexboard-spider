package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"mvdan.cc/xurls"
)

const (
	ListURL    = "https://www.dcard.tw/_api/forums/sex/posts?popular=false"
	ContentURL = "https://www.dcard.tw/_api/posts/%d"
	LinkURL    = "http://www.dcard.tw/f/sex/p/"
	CommentURL = "http://dcard.tw/_api/posts/%d/comments"

	LimitNum = 20
	SleepNum = 60
)

func main() {
	var latestID int64 = 0
	tgBotToken := os.Getenv("TgBotToken")
	tgChannelID, err := strconv.ParseInt(os.Getenv("TgChannelId"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Println(err)
		return
	}
	bot.Debug = false

	for {
		fmt.Println(time.Now())
		fmt.Printf("------------------------ LatestId: %d ------------------------\n", latestID)
		sendCount := 0

		sexArticles, err := getLatestList(latestID)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Total article: %d \n", len(sexArticles))

		for _, value := range sexArticles {
			msgContent := ""
			articleID := value.ID

			articleContent, err := getContent(articleID)
			content := articleContent.Content
			if err != nil {
				log.Println(err)
				return
			}

			comments, err := getComments(articleID)
			if err != nil {
				log.Println(err)
				return
			}

			if articleID > latestID {
				latestID = articleID
			}

			for _, mediaMetaValue := range value.MediaMeta {
				picURL := mediaMetaValue.(map[string]interface{})["normalizedUrl"].(string)
				picType := mediaMetaValue.(map[string]interface{})["type"].(string)

				if strings.Index(strings.ToLower(picType), "thumbnail") == -1 {
					photoMsg := tgbotapi.NewPhotoShare(tgChannelID, picURL)

					bot.Send(photoMsg)
				}
			}

			msgContent += fmt.Sprintf("標題: %s \n", value.Title)
			msgContent += getContentURL(content)
			msgContent += fmt.Sprintf("文章連結: %s \n", LinkURL+strconv.FormatInt(int64(articleID), 10))
			textMsg := tgbotapi.NewMessage(tgChannelID, msgContent)
			myMsg, _ := bot.Send(textMsg)

			for _, value := range comments {
				content := value.Content

				for _, mediaMetaValue := range value.MediaMeta {
					picURL := mediaMetaValue.(map[string]interface{})["normalizedUrl"].(string)
					picType := mediaMetaValue.(map[string]interface{})["type"].(string)

					if strings.Index(strings.ToLower(picType), "thumbnail") == -1 {
						photoMsg := tgbotapi.NewPhotoShare(tgChannelID, picURL)
						photoMsg.Caption = content
						photoMsg.ReplyToMessageID = myMsg.MessageID

						bot.Send(photoMsg)
					}
				}
			}

			sendCount++
			fmt.Printf("%s \n", msgContent)
			fmt.Printf("------------------------ Sent %d articles ------------------------ \n", sendCount)
		}
		time.Sleep(SleepNum * time.Second)
	}
}

//Get the latest article list
func getLatestList(afterID int64) ([]ArticleList, error) {
	params := make(map[string]string)
	params["limit"] = strconv.FormatInt(int64(LimitNum), 10)
	if afterID != 0 {
		params["after"] = strconv.FormatInt(int64(afterID), 10)
	}

	resp, err := Get(ListURL, params, nil)
	if err != nil {
		log.Println(err)
		return []ArticleList{}, err
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []ArticleList{}, err
	}
	//fmt.Println(string(sitemap))

	var articleList []ArticleList
	if err := json.Unmarshal([]byte(sitemap), &articleList); err != nil {
		log.Println(string(sitemap))
		return []ArticleList{}, fmt.Errorf("getLatestList json.Unmarshal error:%s", err)
	}

	return articleList, nil
}

//Get article content
func getContent(articleID int64) (ArticleContent, error) {
	resp, err := Get(fmt.Sprintf(ContentURL, articleID), nil, nil)
	if err != nil {
		log.Println(err)
		return ArticleContent{}, err
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ArticleContent{}, err
	}
	//fmt.Println(string(sitemap))

	var articleContent ArticleContent
	if err := json.Unmarshal([]byte(sitemap), &articleContent); err != nil {
		log.Println(string(sitemap))
		return ArticleContent{}, fmt.Errorf("getContent json.Unmarshal error:%s", err)
	}

	return articleContent, nil
}

//Get article comments
func getComments(articleID int64) ([]ArticleComment, error) {
	resp, err := Get(fmt.Sprintf(CommentURL, articleID), nil, nil)
	if err != nil {
		log.Println(err)
		return []ArticleComment{}, err
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []ArticleComment{}, err
	}
	//fmt.Println(string(sitemap))

	var articleComment []ArticleComment
	if err := json.Unmarshal([]byte(sitemap), &articleComment); err != nil {
		log.Println(string(sitemap))
		return []ArticleComment{}, fmt.Errorf("getComments json.Unmarshal error:%s", err)
	}

	return articleComment, nil
}

//Get the URL in the article
func getContentURL(content string) string {
	var URLstring string

	URLary := xurls.Strict().FindAllString(content, -1)
	for _, URL := range URLary {
		if strings.Index(URL, "i.imgur.com") == -1 {
			URLstring += fmt.Sprintf("內文連結: %s \n", URL)
		}
	}

	return URLstring
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
