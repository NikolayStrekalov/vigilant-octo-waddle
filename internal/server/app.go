package server

import (
	"flag"
	"net/http"
	"os"
	"runtime"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/go-chi/chi/v5"
)

var address string
var exitCodeTooManyArgs = 2

func Start() {
	lgr := logger.InitLog()
	defer func() {
		if err := lgr.Sync(); err != nil {
			panic(err)
		}
	}()

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	prepareRoutes(r)

	flag.StringVar(&address, "a", "localhost:8080", "Эндпоинт сервера HOST:PORT")
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.PrintDefaults()
		runtime.Goexit()
		defer os.Exit(exitCodeTooManyArgs)
	}

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		address = envAddress
	}

	err := http.ListenAndServe(address, r)
	if err != nil {
		panic(err)
	}
}
