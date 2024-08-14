package models

type TransferRequest struct {
	AccountNumber uint64  `json:"account_number"`
	Amount        float64 `json:"amount"`
}
