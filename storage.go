package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobankapi sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (store *PostgresStorage) Init() error {
	return store.createAccountTable()
}

func (store *PostgresStorage) createAccountTable() error {
	query := `create table if not exists account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	_, err := store.db.Exec(query)
	return err
}

func (store *PostgresStorage) CreateAccount(*Account) error {
	return nil
}

func (store *PostgresStorage) UpdateAccount(*Account) error {
	return nil
}

func (store *PostgresStorage) DeleteAccount(id int) error {
	return nil
}

func (store *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
