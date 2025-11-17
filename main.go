package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	srv.ListenAndServe()
}
