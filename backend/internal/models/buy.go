package models

// const tradeData = {
//     symbol: symbolForBackend,
//     quantity: quantity,
//     price: price,
//     type: tradeType
// };

type TradeStruct struct {
	Symbol    string  `json:"symbol"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	TradeType string  `json:"type"`
}
