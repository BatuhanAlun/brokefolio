package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/models"
	"database/sql"
	"encoding/json"
	"fmt" // For string formatting
	"log"
	"net/http"
	"strings"
)

type PriceFetcher interface {
	GetPrice(symbol string) (float64, error)
}

type HTTPPriceFetcher struct{}

func (f *HTTPPriceFetcher) GetPrice(symbol string) (float64, error) {

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/crypto-price?symbol=%s", strings.ToUpper(symbol))) // Adjust port if different
	if err != nil {
		return 0, fmt.Errorf("failed to fetch price for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes := make([]byte, 1024)
		n, _ := resp.Body.Read(bodyBytes)
		return 0, fmt.Errorf("price API returned status %d for %s: %s", resp.StatusCode, symbol, string(bodyBytes[:n]))
	}

	var priceResp models.PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return 0, fmt.Errorf("failed to decode price response for %s: %w", symbol, err)
	}
	if priceResp.Price == 0 && priceResp.Error != "" {
		return 0, fmt.Errorf("price API error for %s: %s", symbol, priceResp.Error)
	}
	if priceResp.Price == 0 {
		return 0, fmt.Errorf("price not found for %s", symbol)
	}

	return priceResp.Price, nil
}

type PortfolioHandler struct {
	DB           *sql.DB
	PriceFetcher PriceFetcher
}

func NewPortfolioHandler(db *sql.DB, pf PriceFetcher) *PortfolioHandler {
	return &PortfolioHandler{DB: db, PriceFetcher: pf}
}

func (h *PortfolioHandler) PortfolioHandler(w http.ResponseWriter, req *http.Request) {
	isAuthenticated, ok := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)
	if !ok || !isAuthenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Auth Context Not Found")
		return
	}

	userID, ok := req.Context().Value(middleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not Found", http.StatusInternalServerError)
		log.Println("User ID not Found in context")
		return
	}

	var portfolioID int
	err := h.DB.QueryRow("SELECT portfolio_id FROM Portfolios WHERE user_id = $1", userID).Scan(&portfolioID)
	if err == sql.ErrNoRows {
		log.Printf("No portfolio found for user %s, returning empty holdings.", userID)
		json.NewEncoder(w).Encode(models.PortfolioData{Holdings: []models.PortfolioAsset{}})
		return
	} else if err != nil {
		log.Printf("Error querying portfolio for user %s: %v", userID, err)
		http.Error(w, "Portföy bilgisi alınırken hata oluştu.", http.StatusInternalServerError)
		return
	}

	rows, err := h.DB.Query(
		`SELECT symbol, quantity, average_buy_price FROM PortfolioAssets WHERE portfolio_id = $1`,
		portfolioID,
	)
	if err != nil {
		log.Printf("Error querying portfolio assets for portfolio ID %d: %v", portfolioID, err)
		http.Error(w, "Varlıklar çekilirken hata oluştu.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var holdings []models.PortfolioAsset
	for rows.Next() {
		var holding models.PortfolioAsset
		if err := rows.Scan(&holding.Symbol, &holding.Quantity, &holding.AverageBuyPrice); err != nil {
			log.Printf("Error scanning portfolio asset row: %v", err)
			// Continue processing other rows, or return error if critical
			continue
		}
		holdings = append(holdings, holding)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating portfolio asset rows: %v", err)
		http.Error(w, "Varlık verileri işlenirken hata oluştu.", http.StatusInternalServerError)
		return
	}

	// Fetch current prices for each holding
	for i := range holdings {
		price, err := h.PriceFetcher.GetPrice(holdings[i].Symbol)
		if err != nil {
			log.Printf("Warning: Could not fetch current price for %s: %v", holdings[i].Symbol, err)
			holdings[i].CurrentPrice = holdings[i].AverageBuyPrice // Fallback: use average buy price
		} else {
			holdings[i].CurrentPrice = price
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.PortfolioData{Holdings: holdings})
}

type TransactionsHandler struct {
	DB *sql.DB
}

func NewTransactionsHandler(db *sql.DB) *TransactionsHandler {
	return &TransactionsHandler{DB: db}
}

func (h *TransactionsHandler) TransactionsHandler(w http.ResponseWriter, req *http.Request) {
	isAuthenticated, ok := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)
	if !ok || !isAuthenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Auth Context Not Found")
		return
	}

	userID, ok := req.Context().Value(middleware.UserIDContextKey).(string)
	if !ok || userID == "" {
		http.Error(w, "User ID not Found", http.StatusInternalServerError)
		log.Println("User ID not Found in context")
		return
	}

	var portfolioID int
	err := h.DB.QueryRow("SELECT portfolio_id FROM Portfolios WHERE user_id = $1", userID).Scan(&portfolioID)
	if err == sql.ErrNoRows {
		// No portfolio yet, return empty transactions but no error
		log.Printf("No portfolio found for user %s, returning empty transactions.", userID)
		json.NewEncoder(w).Encode(models.TransactionsData{Transactions: []models.Transaction{}})
		return
	} else if err != nil {
		log.Printf("Error querying portfolio for user %s: %v", userID, err)
		http.Error(w, "Portföy bilgisi alınırken hata oluştu.", http.StatusInternalServerError)
		return
	}

	rows, err := h.DB.Query(
		`SELECT transaction_date, symbol, trade_type, quantity, price_per_unit
         FROM Transactions WHERE portfolio_id = $1 ORDER BY transaction_date DESC`,
		portfolioID,
	)
	if err != nil {
		log.Printf("Error querying transactions for portfolio ID %d: %v", portfolioID, err)
		http.Error(w, "İşlem geçmişi çekilirken hata oluştu.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.Timestamp, &transaction.Symbol, &transaction.Type, &transaction.Quantity, &transaction.Price); err != nil {
			log.Printf("Error scanning transaction row: %v", err)
			continue // Skip malformed row and continue
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating transaction rows: %v", err)
		http.Error(w, "İşlem geçmişi verileri işlenirken hata oluştu.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.TransactionsData{Transactions: transactions})
}
