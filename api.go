package main

import "net/http"

type APIServer struct {
	listenAddress string
}

func newAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}

//starting up server
func (s *APIServer) Run() {

}

func (s *APIServer) handleAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleCreateAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(writer http.ResponseWriter, request *http.Request) error {
	return nil
}
