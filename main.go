package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/varun-muthanna/URL-Shortner/handler"
)

func main() {

	h := handler.NewURLShort()

	m := mux.NewRouter()

	getrouter := m.Methods("GET").Subrouter()
	postrouter := m.Methods("POST").Subrouter()

	getrouter.HandleFunc("/", h.ServeForm)
	getrouter.HandleFunc("/short/{shortcode:[a-zA-Z0-9]+}", h.HandleRedirect)

	postrouter.HandleFunc("/shorten", h.HandleShorten)

	fmt.Println("URL Shortener is running on :3031")
	err := http.ListenAndServe(":3031", m)

	if err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}

}
