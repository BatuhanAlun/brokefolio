package models

import "time"

type TradeStruct struct {
	Symbol    string  `json:"symbol"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	TradeType string  `json:"type"`
}

type PortfolioAsset struct {
	Symbol          string  `json:"symbol"`
	Quantity        float64 `json:"quantity"`
	AverageBuyPrice float64 `json:"averageBuyPrice"`
	CurrentPrice    float64 `json:"currentPrice"`
}

type PortfolioData struct {
	Holdings []PortfolioAsset `json:"holdings"`
	Error    string           `json:"error,omitempty"`
}

type Transaction struct {
	Timestamp time.Time `json:"timestamp"`
	Symbol    string    `json:"symbol"`
	Type      string    `json:"type"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
}

type TransactionsData struct {
	Transactions []Transaction `json:"transactions"`
	Error        string        `json:"error,omitempty"`
}

type PriceResponse struct {
	Price float64 `json:"price"`
	Error string  `json:"error"`
}
