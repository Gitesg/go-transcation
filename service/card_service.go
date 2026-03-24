package service

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"../model"
	"../repository"
)

type Service struct {
	repo *repository.Repo
}

func NewService(r *repository.Repo) *Service {
	return &Service{repo: r}
}

func (s *Service) Transaction(req model.TransactionRequest) (string, error) {

	s.repo.Mu.Lock()
	defer s.repo.Mu.Unlock()

	card, ok := s.repo.Cards[req.CardNumber]
	if !ok {
		return "", fmt.Errorf("invalid card number")
	}

	if card.Status != "active" {
		return "", fmt.Errorf("card is not active")
	}

	if !checkPin(card.PinHash, req.Pin) {
		return "", fmt.Errorf("invalid pin")
	}

	if card.Balance < req.Amount {

		s.repo.History[req.CardNumber] = append(
			s.repo.History[req.CardNumber],
			model.TransactionHistory{
				TransactionId: generateTransactionId(),
				CardNumber:    req.CardNumber,
				Type:          req.Type,
				Amount:        req.Amount,
				Status:        "FAILED_INSUFFICIENT_BALANCE",
				Timestamp:     time.Now().Format(time.RFC3339),
			},
		)

		return "", fmt.Errorf("insufficient balance")
	}

	card.Balance -= req.Amount
	s.repo.Cards[req.CardNumber] = card

	s.repo.History[req.CardNumber] = append(
		s.repo.History[req.CardNumber],
		model.TransactionHistory{
			TransactionId: generateTransactionId(),
			CardNumber:    req.CardNumber,
			Type:          req.Type,
			Amount:        req.Amount,
			Status:        "SUCCESS",
			Timestamp:     time.Now().Format(time.RFC3339),
		},
	)

	return "SUCCESS", nil
}

func (s *Service) GetBalance(cardNumber, pin string) (float64, error) {

	s.repo.Mu.RLock()
	defer s.repo.Mu.RUnlock()

	card, ok := s.repo.Cards[cardNumber]
	if !ok {
		return 0, fmt.Errorf("invalid card number")
	}

	if card.Status != "active" {
		return 0, fmt.Errorf("card is not active")
	}

	if !checkPin(card.PinHash, pin) {
		return 0, fmt.Errorf("invalid pin")
	}

	return card.Balance, nil
}

func (s *Service) GetHistory(cardNumber, pin string) ([]model.TransactionHistory, error) {

	s.repo.Mu.RLock()
	defer s.repo.Mu.RUnlock()

	card, ok := s.repo.Cards[cardNumber]
	if !ok {
		return nil, fmt.Errorf("invalid card number")
	}

	if card.Status != "active" {
		return nil, fmt.Errorf("card is not active")
	}

	if !checkPin(card.PinHash, pin) {
		return nil, fmt.Errorf("invalid pin")
	}

	return s.repo.History[cardNumber], nil
}

func hashPin(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPin(hash, pin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin))
	return err == nil
}

func generateTransactionId() string {
	return fmt.Sprintf("txn-%d", time.Now().UnixNano())
}
