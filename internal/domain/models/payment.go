package models

import "time"

type Payment struct {
	ID            string    `json:"id"`
	RequestID     string    `json:"requestId"`
	HelperID      string    `json:"helperId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	EscrowExpires time.Time `json:"escrowExpires"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Payout struct {
	ID        string    `json:"id"`
	HelperID  string    `json:"helperId"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	Processed time.Time `json:"processed"`
}
