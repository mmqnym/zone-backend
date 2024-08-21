package models

type TransactionRecord struct {
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	TxId      int64  `json:"txId"`
	Type      string `json:"type"`
}
