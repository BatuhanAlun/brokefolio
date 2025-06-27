package handlers

import (
	"brokefolio/backend/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func fetchStockNews(apiKey string) ([]models.FrontendNewsArticle, error) {

	newsAPIURL := fmt.Sprintf("https://newsapi.org/v2/everything?q=stock+market+OR+finance+OR+investment&language=en&sortBy=publishedAt&pageSize=5&apiKey=%s", apiKey)

	resp, err := http.Get(newsAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching stock news: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("external NewsAPI returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading NewsAPI response body: %w", err)
	}

	var newsAPIResponse models.NewsAPIResponse
	err = json.Unmarshal(body, &newsAPIResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling NewsAPI response: %w, body: %s", err, string(body))
	}

	var stockNews []models.FrontendNewsArticle
	for _, article := range newsAPIResponse.Articles {
		t, err := time.Parse(time.RFC3339, article.PublishedAt)
		formattedDate := "Unknown Date"
		if err == nil {
			formattedDate = t.Format("02 Jan 2006 15:04")
		}
		stockNews = append(stockNews, models.FrontendNewsArticle{
			Title:  article.Title,
			URL:    article.URL,
			Source: article.Source.Name,
			Date:   formattedDate,
		})
	}
	return stockNews, nil
}

func fetchCryptoNewsExternal(apiKey string) ([]models.FrontendNewsArticle, error) {

	cryptoAPIURL := fmt.Sprintf("https://newsapi.org/v2/everything?q=cryptocurrency+OR+bitcoin+OR+ethereum&language=en&sortBy=publishedAt&pageSize=5&apiKey=%s", apiKey)

	resp, err := http.Get(cryptoAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching crypto news: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("external crypto NewsAPI returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading crypto NewsAPI response body: %w", err)
	}

	var newsAPIResponse models.NewsAPIResponse
	err = json.Unmarshal(body, &newsAPIResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling crypto NewsAPI response: %w, body: %s", err, string(body))
	}

	var cryptoNews []models.FrontendNewsArticle
	for _, article := range newsAPIResponse.Articles {
		t, err := time.Parse(time.RFC3339, article.PublishedAt)
		formattedDate := "Unknown Date"
		if err == nil {
			formattedDate = t.Format("02 Jan 2006 15:04")
		}
		cryptoNews = append(cryptoNews, models.FrontendNewsArticle{
			Title:  article.Title,
			URL:    article.URL,
			Source: article.Source.Name,
			Date:   formattedDate,
		})
	}
	return cryptoNews, nil
}

func CombinedNewsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	newsAPIKey := os.Getenv("NEWS_API_KEY")
	if newsAPIKey == "" {
		http.Error(w, "NEWS_API_KEY environment variable not set", http.StatusInternalServerError)
		return
	}

	var (
		stockNews  []models.FrontendNewsArticle
		cryptoNews []models.FrontendNewsArticle
		wg         sync.WaitGroup
		errStock   error
		errCrypto  error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		stockNews, errStock = fetchStockNews(newsAPIKey)
		if errStock != nil {
			log.Printf("Error fetching stock news: %v", errStock)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cryptoNews, errCrypto = fetchCryptoNewsExternal(newsAPIKey)

		if errCrypto != nil {
			log.Printf("Error fetching crypto news: %v", errCrypto)
		}
	}()

	wg.Wait()

	var allNews []models.FrontendNewsArticle
	allNews = append(allNews, stockNews...)
	allNews = append(allNews, cryptoNews...)

	if len(allNews) > 10 {
		allNews = allNews[:10]
	}

	response := models.FrontendNewsResponse{Results: allNews}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding combined news response: %v", err)
		http.Error(w, "Failed to encode news response", http.StatusInternalServerError)
	}
}
