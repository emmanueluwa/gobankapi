package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

// starting up server
func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(server.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(server.handleGetAccountByID))

	log.Println("JSON API server running on port: ", server.listenAddress)
	http.ListenAndServe(server.listenAddress, router)
}

func (server *APIServer) handleAccount(writer http.ResponseWriter, request *http.Request) error {
	//handling what the request is
	//WOULD USE OF SWITCH BE BETTER?
	if request.Method == "GET" {
		return server.handleGetAccount(writer, request)
	}
	if request.Method == "POST" {
		return server.handleCreateAccount(writer, request)
	}
	if request.Method == "DELETE" {
		return server.handleDeleteAccount(writer, request)
	}

	return fmt.Errorf("Invalid method %s", request.Method)
}

func (server *APIServer) handleGetAccount(writer http.ResponseWriter, request *http.Request) error {
	accounts, err := server.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, accounts)
}

func (server *APIServer) handleGetAccountByID(writer http.ResponseWriter, request *http.Request) error {
	// the id comes as a string and needs to be converted
	idString := mux.Vars(request)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		return fmt.Errorf("invalid id %s", idString)
	}

	account, err := server.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, account)
}

func (server *APIServer) handleCreateAccount(writer http.ResponseWriter, request *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(request.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	//store the created account
	if err := server.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, account)
}

func (server *APIServer) handleDeleteAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (server *APIServer) handleTransfer(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

//helper functions...

// writing json for the API
func WriteJSON(writer http.ResponseWriter, status int, value any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(value)
}

// function sig we wish to use
type apiFunc func(http.ResponseWriter, *http.Request) error

// type for API errors
type ApiError struct {
	Error string `json:"error"`
}

// help decorate apifunc into the handlefunctions we wish to use
func makeHTTPHandlerFunc(function apiFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := function(writer, request); err != nil {
			//handle error
			WriteJSON(writer, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
