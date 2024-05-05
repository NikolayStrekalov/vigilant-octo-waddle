package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Start() {
	r := chi.NewRouter()
	prepareRoutes(r)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
