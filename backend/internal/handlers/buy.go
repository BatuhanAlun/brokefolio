package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	var tradeReq models.TradeStruct // Renamed to tradeReq for consistency
	err := json.NewDecoder(req.Body).Decode(&tradeReq)

	// Normalize trade type to uppercase
	tradeReq.TradeType = strings.ToUpper(tradeReq.TradeType)

	if err != nil {
		http.Error(w, "Something Went Wrong with Trade", http.StatusBadRequest)
		log.Printf("Server Could not Decode JSON: %v", err)
		return
	}

	if tradeReq.Quantity <= 0 {
		http.Error(w, "Quantity needs to be positive", http.StatusBadRequest)
		return
	}

	var totalTradeAmount float64 = tradeReq.Quantity * tradeReq.Price

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

	// Insert transaction record regardless of BUY/SELL
	_, err = h.DB.Exec(
		`INSERT INTO Transactions (portfolio_id, symbol, trade_type, quantity, price_per_unit, total_amount, transaction_date)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		portfolioID, tradeReq.Symbol, tradeReq.TradeType, tradeReq.Quantity, tradeReq.Price, totalTradeAmount, time.Now(),
	)
	if err != nil {
		log.Printf("Failed to insert transaction: %v", err)
		http.Error(w, "İşlem geçmişi kaydedilemedi.", http.StatusInternalServerError)
		return
	}

	// --- Start of Corrected PortfolioAsset Update Logic ---
	var existingQuantity float64
	var existingAvgBuyPrice float64

	// Try to select existing asset
	err = h.DB.QueryRow(
		`SELECT quantity, average_buy_price FROM PortfolioAssets WHERE portfolio_id = $1 AND symbol = $2`,
		portfolioID, tradeReq.Symbol,
	).Scan(&existingQuantity, &existingAvgBuyPrice)

	if err == sql.ErrNoRows {
		// Asset does not exist for this portfolio
		if tradeReq.TradeType == "SELL" {
			// Cannot sell an asset that doesn't exist in holdings
			http.Error(w, "Bu varlığa sahip değilsiniz, satış yapılamaz.", http.StatusBadRequest)
			return
		}
		// If it's a BUY and asset doesn't exist, insert new
		_, err = h.DB.Exec(
			`INSERT INTO PortfolioAssets (portfolio_id, symbol, quantity, average_buy_price)
             VALUES ($1, $2, $3, $4)`,
			portfolioID, tradeReq.Symbol, tradeReq.Quantity, tradeReq.Price,
		)
		if err != nil {
			log.Printf("Failed to insert new portfolio asset for BUY: %v", err)
			http.Error(w, "Varlık portföye eklenirken hata oluştu.", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		// Other database error during select
		log.Printf("Error querying existing portfolio asset: %v", err)
		http.Error(w, "Varlık kontrol edilirken hata oluştu.", http.StatusInternalServerError)
		return
	} else {
		// Asset exists, now update based on trade type
		var newQuantity float64
		var newAverageBuyPrice float64

		switch tradeReq.TradeType {
		case "BUY":
			newQuantity = existingQuantity + tradeReq.Quantity
			// Recalculate average buy price for BUYs
			newAverageBuyPrice = ((existingAvgBuyPrice * existingQuantity) + (tradeReq.Price * tradeReq.Quantity)) / newQuantity

		case "SELL":
			if tradeReq.Quantity > existingQuantity {
				http.Error(w, fmt.Sprintf("Yetersiz miktar! Elinizde %.2f adet %s var, %.2f adet satmaya çalışıyorsunuz.", existingQuantity, tradeReq.Symbol, tradeReq.Quantity), http.StatusBadRequest)
				return
			}
			newQuantity = existingQuantity - tradeReq.Quantity

			totalHoldingPrice := newQuantity * tradeReq.Price
			// For sells, average_buy_price does not change unless it's a partial sell that leaves some.
			// If the quantity becomes 0, the asset should ideally be removed.
			newAverageBuyPrice = totalHoldingPrice / newQuantity // Average buy price remains the same for sells

		default:
			http.Error(w, "Geçersiz işlem tipi.", http.StatusBadRequest)
			return
		}

		if newQuantity <= 0.000001 { // Use a small epsilon for float comparison to account for tiny remainders
			// If quantity becomes zero or negative (due to float precision), remove the asset
			_, err = h.DB.Exec(
				`DELETE FROM PortfolioAssets WHERE portfolio_id = $1 AND symbol = $2`,
				portfolioID, tradeReq.Symbol,
			)
			if err != nil {
				log.Printf("Failed to delete portfolio asset (quantity became zero): %v", err)
				http.Error(w, "Portföy varlığı silinirken hata oluştu.", http.StatusInternalServerError)
				return
			}
			log.Printf("Asset %s for portfolio %d removed as quantity became zero or less.", tradeReq.Symbol, portfolioID)
		} else {
			// Update existing asset with new quantity and (potentially) new average buy price
			_, err = h.DB.Exec(
				`UPDATE PortfolioAssets SET quantity = $1, average_buy_price = $2 WHERE portfolio_id = $3 AND symbol = $4`,
				newQuantity, newAverageBuyPrice, portfolioID, tradeReq.Symbol,
			)
			if err != nil {
				log.Printf("Failed to update portfolio asset: %v", err)
				http.Error(w, "Varlık portföyde güncellenirken hata oluştu.", http.StatusInternalServerError)
				return
			}
		}
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "İşlem Başarılı!"})

}
