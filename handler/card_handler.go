package handler

import (
	"bank-transaction/model"
	"bank-transaction/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}
func (h *Handler) Transaction(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := h.service.Transaction(req)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "FAILED",
			"respCode": "99",
			"message":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   result,
		"respCode": "00",
		"message":  "Transaction successful",
	})
}
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := r.URL.Query().Get("cardNumber")
	pin := r.URL.Query().Get("pin")

	balance, err := h.service.GetBalance(cardNumber, pin)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "FAILED",
			"respCode": "99",
			"message":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "SUCCESS",
		"respCode": "00",
		"balance":  balance,
	})
}func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := r.URL.Query().Get("cardNumber")
	pin := r.URL.Query().Get("pin")

	history, err := h.service.GetHistory(cardNumber, pin)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "FAILED",
			"respCode": "99",
			"message":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "SUCCESS",
		"respCode": "00",
		"data":     history,
	})
}