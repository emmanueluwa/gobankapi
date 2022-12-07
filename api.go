package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// writing json for the API
func WriteJSON(writer http.ResponseWriter, status int, value any) error {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(writer).Encode(value)
}

// function sig we wish to use
type apiFunc func(http.ResponseWriter, *http.Request) error

// type for API errors
type ApiError struct {
	Error string
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

type APIServer struct {
	listenAddress string
}

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}

// starting up server
func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(server.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(server.handleGetAccount))

	log.Println("JSON API seerver running on port: ", server.listenAddress)
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
	id := mux.Vars(request)["id"]

	fmt.Println(id)

	// account := NewAccount("Wise", "O")
	return WriteJSON(writer, http.StatusOK, &Account{})
}

func (server *APIServer) handleCreateAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (server *APIServer) handleDeleteAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (server *APIServer) handleTransfer(writer http.ResponseWriter, request *http.Request) error {
	return nil
}
