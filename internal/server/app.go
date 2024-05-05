package server

import (
	"net/http"
)

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
