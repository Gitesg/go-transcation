package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

var (
	cards   = make(map[string]Card)
	history = make(map[string][]TransactionHistory)
	mu      sync.RWMutex
)

func sendJSON(w http.ResponseWriter, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func hash(pin string) string {
	h := sha256.New()
	h.Write([]byte(pin))
	return hex.EncodeToString(h.Sum(nil))
}

func generateTransactionId() string {
	return fmt.Sprintf("txn-%d", time.Now().UnixNano())
}

func transaction(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		sendJSON(w, APIResponse{"FAILED", "01", "Method not allowed", 0, nil})
		return
	}

	var body TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendJSON(w, APIResponse{"FAILED", "02", "Invalid request body", 0, nil})
		return
	}

	if body.Amount <= 0 {
		sendJSON(w, APIResponse{"FAILED", "03", "Invalid amount", 0, nil})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	card, ok := cards[body.CardNumber]
	if !ok {
		sendJSON(w, APIResponse{"FAILED", "05", "Invalid card", 0, nil})
		return
	}

	if strings.ToUpper(card.Status) != "ACTIVE" {
		sendJSON(w, APIResponse{"FAILED", "04", "Card blocked", 0, nil})
		return
	}

	if card.PinHash != hash(body.Pin) {
		sendJSON(w, APIResponse{"FAILED", "06", "Invalid PIN", 0, nil})
		return
	}

	var status string

	switch body.Type {
	case "withdraw":
		if card.Balance < body.Amount {
			status = "FAILED"
			history[body.CardNumber] = append(history[body.CardNumber], TransactionHistory{
				TransactionId: generateTransactionId(),
				CardNumber:    body.CardNumber,
				Type:          body.Type,
				Amount:        body.Amount,
				Status:        status,
				Timestamp:     time.Now().Format(time.RFC3339),
			})
			sendJSON(w, APIResponse{"FAILED", "99", "Insufficient balance", card.Balance, nil})
			return
		}
		card.Balance -= body.Amount

	case "topup":
		card.Balance += body.Amount

	default:
		sendJSON(w, APIResponse{"FAILED", "07", "Invalid transaction type", 0, nil})
		return
	}

	cards[body.CardNumber] = card

	status = "SUCCESS"
	history[body.CardNumber] = append(history[body.CardNumber], TransactionHistory{
		TransactionId: generateTransactionId(),
		CardNumber:    body.CardNumber,
		Type:          body.Type,
		Amount:        body.Amount,
		Status:        status,
		Timestamp:     time.Now().Format(time.RFC3339),
	})

	sendJSON(w, APIResponse{"SUCCESS", "00", "", card.Balance, nil})
}

func balance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		sendJSON(w, APIResponse{"FAILED", "01", "Method not allowed", 0, nil})
		return
	}

	cardNumber := strings.TrimPrefix(r.URL.Path, "/api/card/balance/")
	pin := r.URL.Query().Get("pin")

	mu.RLock()
	defer mu.RUnlock()

	card, ok := cards[cardNumber]
	if !ok {
		sendJSON(w, APIResponse{"FAILED", "05", "Invalid card", 0, nil})
		return
	}

	if card.PinHash != hash(pin) {
		sendJSON(w, APIResponse{"FAILED", "06", "Invalid PIN", 0, nil})
		return
	}

	sendJSON(w, APIResponse{"SUCCESS", "00", "", card.Balance, nil})
}

func transactionHistory(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		sendJSON(w, APIResponse{"FAILED", "01", "Method not allowed", 0, nil})
		return
	}

	cardNumber := strings.TrimPrefix(r.URL.Path, "/api/card/transactions/")
	pin := r.URL.Query().Get("pin")

	mu.RLock()
	defer mu.RUnlock()

	card, ok := cards[cardNumber]
	if !ok {
		sendJSON(w, APIResponse{"FAILED", "05", "Invalid card", 0, nil})
		return
	}

	if card.PinHash != hash(pin) {
		sendJSON(w, APIResponse{"FAILED", "06", "Invalid PIN", 0, nil})
		return
	}

	sendJSON(w, APIResponse{
		Status:   "SUCCESS",
		RespCode: "00",
		Data:     history[cardNumber],
	})
}

func main() {

	cards["4123456789012345"] = Card{
		CardHolder: "John Doe",
		PinHash:    hash("1234"),
		Balance:    1000,
		Status:     "ACTIVE",
	}

	http.HandleFunc("/api/transaction", transaction)
	http.HandleFunc("/api/card/balance/", balance)
	http.HandleFunc("/api/card/transactions/", transactionHistory)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
