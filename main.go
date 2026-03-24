package main

import (
	"bank-transaction/handler"
	"bank-transaction/repository"
	"bank-transaction/service"
	"net/http"
)

func main() {
	repo := repository.NewRepo()
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	http.HandleFunc("/api/transaction", handler.Transaction)
	http.HandleFunc("/api/card/balance", handler.GetBalance)
	http.HandleFunc("/api/card/transactions", handler.GetHistory)

	http.ListenAndServe(":8080", nil)
}
