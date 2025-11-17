package main

import (
	"net/http"
)

func main() {
	const listenAddr = ":8080"
	const filePathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	srv := http.Server{
		Handler: mux,
		Addr:    listenAddr,
	}

	srv.ListenAndServe()
}
