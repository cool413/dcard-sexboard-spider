package main

type ArticleList struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Media []struct {
		URL string `json:"url"`
	} `json:"media"`
	MediaMeta []interface{} `json:"mediaMeta"`
}

type ArticleContent struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}

type ArticleComment struct {
	ID        string        `json:"id"`
	Content   string        `json:"content"`
	MediaMeta []interface{} `json:"mediaMeta"`
}
