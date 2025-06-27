package models

type NewsAPIArticle struct {
	Source struct {
		Name string `json:"name"`
	} `json:"source"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

type NewsAPIResponse struct {
	Articles []NewsAPIArticle `json:"articles"`
}

type CryptoCompareNewsArticle struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Source    string `json:"source"`
	Published int64  `json:"published_on"`
}

type CryptoCompareNewsResponse struct {
	Data []CryptoCompareNewsArticle `json:"Data"`
}

type FrontendNewsArticle struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source"`
	Date   string `json:"date"`
}

type FrontendNewsResponse struct {
	Results []FrontendNewsArticle `json:"results"`
}
