package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type BuyDBHandler struct {
	DB *sql.DB
}

func NewBuyHandler(db *sql.DB) *BuyDBHandler {
	return &BuyDBHandler{DB: db}
}

func (h *BuyDBHandler) TradeHandler(w http.ResponseWriter, req *http.Request) {

	isAutheticated, ok := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)

	if !ok || !isAutheticated {
		http.Error(w, "User Auth Failed", http.StatusUnauthorized)
		log.Println("Auth Context Not Found")
		return
	}

	userID, ok := req.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not Found", http.StatusInternalServerError)
		log.Println("User ID not Found")
		return
	}

	var TradeReq models.TradeStruct

	err := json.NewDecoder(req.Body).Decode(&TradeReq)

	if err != nil {
		http.Error(w, "Something Went Wrong with Trade", http.StatusBadRequest)
		log.Printf("Server Could not Decode JSON %v", err)
		return
	}

	if TradeReq.Quantity <= 0 {
		http.Error(w, "Quantitiy need to be positive", http.StatusBadRequest)
		return
	}

	var totalTradeAmount float64 = float64(TradeReq.Quantity) * TradeReq.Price

	var portfolioID int
	err = h.DB.QueryRow("SELECT portfolio_id FROM Portfolios WHERE user_id = $1", userID).Scan(&portfolioID)

	if err == sql.ErrNoRows {
		log.Printf("Creating new portfolio for user_id: %s", userID)
		err = h.DB.QueryRow(
			`INSERT INTO Portfolios (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING portfolio_id`,
			userID, time.Now(), time.Now(),
		).Scan(&portfolioID)
		if err != nil {
			log.Printf("Failed to create new portfolio for user %s: %v", userID, err)
			http.Error(w, "Yeni portföy oluşturulurken hata oluştu.", http.StatusInternalServerError)
			return
		}
		log.Printf("New portfolio created with ID %d for user_id: %s", portfolioID, userID)
	} else if err != nil {
		log.Printf("Error querying portfolio for user %s: %v", userID, err)
		http.Error(w, "Portföy bilgisi alınırken hata oluştu.", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec(
		`INSERT INTO Transactions (portfolio_id, symbol, trade_type, quantity, price_per_unit, total_amount, transaction_date)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		portfolioID, TradeReq.Symbol, TradeReq.TradeType, TradeReq.Quantity, TradeReq.Price, totalTradeAmount, time.Now(),
	)

	if err != nil {
		log.Printf("Failed to insert transaction: %v", err)
		http.Error(w, "İşlem geçmişi kaydedilemedi.", http.StatusInternalServerError)
		return
	}

	var existingQuantity float64
	var existingAvgBuyPrice float64

	// Try to select existing asset
	err = h.DB.QueryRow(
		`SELECT quantity, average_buy_price FROM PortfolioAssets WHERE portfolio_id = $1 AND symbol = $2`,
		portfolioID, TradeReq.Symbol,
	).Scan(&existingQuantity, &existingAvgBuyPrice)

	if err == sql.ErrNoRows {
		// Asset does not exist for this portfolio, insert new
		_, err = h.DB.Exec(
			`INSERT INTO PortfolioAssets (portfolio_id, symbol, quantity, average_buy_price)
             VALUES ($1, $2, $3, $4)`,
			portfolioID, TradeReq.Symbol, TradeReq.Quantity, TradeReq.Price,
		)
		if err != nil {
			log.Printf("Failed to insert new portfolio asset: %v", err)
			http.Error(w, "Varlık portföye eklenirken hata oluştu.", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		// Other database error during select
		log.Printf("Error querying existing portfolio asset: %v", err)
		http.Error(w, "Varlık kontrol edilirken hata oluştu.", http.StatusInternalServerError)
		return
	} else {
		// Asset exists, update quantity and recalculate average_buy_price
		// This calculation assumes 'buy' operation. For 'sell', you'd subtract.
		newQuantity := existingQuantity + TradeReq.Quantity

		newAverageBuyPrice := ((existingAvgBuyPrice * existingQuantity) + (TradeReq.Price * TradeReq.Quantity)) / newQuantity

		_, err = h.DB.Exec(
			`UPDATE PortfolioAssets SET quantity = $1, average_buy_price = $2 WHERE portfolio_id = $3 AND symbol = $4`,
			newQuantity, newAverageBuyPrice, portfolioID, TradeReq.Symbol,
		)
		if err != nil {
			log.Printf("Failed to update portfolio asset: %v", err)
			http.Error(w, "Varlık portföyde güncellenirken hata oluştu.", http.StatusInternalServerError)
			return
		}
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Varlık portföyünüze başarıyla eklendi!"})

}
