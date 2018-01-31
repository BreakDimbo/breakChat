package main

import (
	"bChat/guard/action"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/debug", action.Debug).Methods("GET")
	router.HandleFunc("/insert", action.EntryPlug).Methods("POST")
	return router
}
