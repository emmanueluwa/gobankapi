package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
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

	router.HandleFunc("/account/{id}", WithJWTAuth(makeHTTPHandlerFunc(server.handleGetAccountByID), server.store))

	//using post request instead of get, to reduce exposure of account number in browser history/webserver logs
	router.HandleFunc("/transfer", makeHTTPHandlerFunc(server.handleTransfer))

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
	if request.Method == "GET" {
		id, err := getID(request)
		if err != nil {
			//return clean json error
			return err
		}

		account, err := server.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(writer, http.StatusOK, account)
	}

	if request.Method == "DELETE" {
		return server.handleDeleteAccount(writer, request)
	}

	return fmt.Errorf("method %s not allowed", request.Method)
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

	//generate jwt token
	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT token: ", tokenString)

	return WriteJSON(writer, http.StatusOK, account)
}

func (server *APIServer) handleDeleteAccount(writer http.ResponseWriter, request *http.Request) error {
	id, err := getID(request)
	if err != nil {
		//return clean json error
		return err
	}

	if err := server.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, map[string]int{"deleted": id})
}

func (server *APIServer) handleTransfer(writer http.ResponseWriter, request *http.Request) error {
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(request.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer request.Body.Close()

	return WriteJSON(writer, http.StatusOK, transferRequest)
}

/***


helper functions...


**/

// writing json for the API
func WriteJSON(writer http.ResponseWriter, status int, value any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)

	return json.NewEncoder(writer).Encode(value)
}

func createJWT(account *Account) (string, error) {
	// Create the Claims
	claims := &jwt.MapClaims{
		"ExpiresAt":     jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"accountNumber": account.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//convert string to byte
	return token.SignedString([]byte(secret))
}

func accessDenied(writer http.ResponseWriter) {
	WriteJSON(writer, http.StatusForbidden, ApiError{Error: "access denied"})
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store Storage) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Calling JWT auth middleware")

		tokenString := request.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil {
			accessDenied(writer)
			return
		}
		if !token.Valid {
			accessDenied(writer)
			return
		}

		userID, err := getID(request)
		if err != nil {
			accessDenied(writer)
			return
		}
		account, err := store.GetAccountByID(userID)
		//never return errors that give hints on what to do(hackers), use logs to assist debugging
		if err != nil {
			accessDenied(writer)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		// panic(reflect.TypeOf(claims["accountNumber"])) --->> prints float64
		// use own claims instead of default to avoid having to hack desired result
		if account.Number != int64(claims["accountNumber"].(float64)) {
			accessDenied(writer)
			return
		}

		handlerFunc(writer, request)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

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

func getID(request *http.Request) (int, error) {
	// the id comes as a string and needs to be converted
	idString := mux.Vars(request)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		return id, fmt.Errorf("invalid id %s", idString)
	}
	//ensuring useful error is given
	return id, nil
}
