package main

import "math/rand"

type Account struct {
	ID         int
	FirstName  string
	LastName   string
	CardNumber int64
	Balance    int64
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:         rand.Intn(10000),
		FirstName:  firstName,
		LastName:   lastName,
		CardNumber: int64(rand.Intn(100000)),
	}
}
