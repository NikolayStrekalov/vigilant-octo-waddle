package server

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

var address string
var exitCodeTooManyArgs = 2

func Start() {
	r := chi.NewRouter()
	prepareRoutes(r)

	flag.StringVar(&address, "a", "localhost:8080", "Эндпоинт сервера HOST:PORT")
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.PrintDefaults()
		os.Exit(exitCodeTooManyArgs)
	}

	err := http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}
