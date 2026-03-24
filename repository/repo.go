package repository

import (
	"sync"

	"../model"
)

type Repo struct {
	Cards   map[string]model.Card
	History map[string][]model.TransactionHistory
	Mu      sync.RWMutex
}

func NewRepo() *Repo {
	return &Repo{
		Cards:   make(map[string]model.Card),
		History: make(map[string][]model.TransactionHistory),
	}
}
