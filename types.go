package main

import (
	"math/rand"
	"time"
)

type CreateAccount struct {
	Firstname string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	CardNumber int64     `json:"cardNumber"`
	Balance    int64     `json:"balance"`
	CreatedAt  time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:         rand.Intn(10000),
		FirstName:  firstName,
		LastName:   lastName,
		CardNumber: int64(rand.Intn(100000)),
		CreatedAt:  time.Now().UTC(),
	}
}
