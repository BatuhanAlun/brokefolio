package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type FMPQuoteResponse []struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func StockPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	fmpAPIKey := os.Getenv("PRICE_API_KEY")
	if fmpAPIKey == "" {
		http.Error(w, "PRICE_API_KEY environment variable not set", http.StatusInternalServerError)
		return
	}

	fmpURL := fmt.Sprintf("https://financialmodelingprep.com/api/v3/quote/%s?apikey=%s", symbol, fmpAPIKey)

	resp, err := http.Get(fmpURL)
	if err != nil {
		log.Printf("Error fetching stock price from FMP for %s: %v", symbol, err)
		http.Error(w, "Failed to fetch stock price", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("FMP API returned non-OK status for %s: %d, body: %s", symbol, resp.StatusCode, string(bodyBytes))
		http.Error(w, fmt.Sprintf("External stock price API error: %d %s", resp.StatusCode, string(bodyBytes)), resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading FMP API response body: %v", err)
		http.Error(w, "Failed to read stock price response", http.StatusInternalServerError)
		return
	}

	var fmpResponse FMPQuoteResponse
	err = json.Unmarshal(body, &fmpResponse)
	if err != nil {
		log.Printf("Error unmarshalling FMP API response: %v, body: %s", err, string(body))
		http.Error(w, "Failed to parse stock price data", http.StatusInternalServerError)
		return
	}

	if len(fmpResponse) == 0 {
		http.Error(w, "Stock price not found for symbol", http.StatusNotFound)
		return
	}

	currentPrice := fmpResponse[0].Price
	if currentPrice == 0.0 {
		http.Error(w, "Invalid price returned for symbol", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{"price": currentPrice}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
