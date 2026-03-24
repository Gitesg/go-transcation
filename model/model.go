package model

type Card struct {
	CardHolder string
	PinHash    string
	Balance    float64
	Status     string
}

type TransactionRequest struct {
	CardNumber string  `json:"cardNumber"`
	Pin        string  `json:"pin"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
}

type TransactionHistory struct {
	TransactionId string  `json:"transactionId"`
	CardNumber    string  `json:"cardNumber"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"`
}

type APIResponse struct {
	Status   string      `json:"status"`
	RespCode string      `json:"respCode"`
	Message  string      `json:"message,omitempty"`
	Balance  float64     `json:"balance,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
