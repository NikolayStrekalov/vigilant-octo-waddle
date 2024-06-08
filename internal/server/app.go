package server

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/memstorage"

	chi "github.com/go-chi/chi/v5"
)

func Start() {
	// Инициализируем логирование
	lgr := logger.InitLog()
	defer func() {
		_ = lgr.Sync()
	}()

	// Настраиваем приложение
	r := appRouter()
	if err := fillConfig(); err != nil {
		flag.PrintDefaults()
		panic(err)
	}
	Storage = memstorage.NewMemStorage(ServerConfig.IsSyncDump(), ServerConfig.FileStoragePath)

	// Восстанавливаем данные
	if ServerConfig.RestoreStore {
		Storage.Restore()
	}
	// Периодически сохраняем данные
	if ServerConfig.DumpEnabled() {
		go periodicDump()
	}

	// При выходе сохраняем данные
	var server *http.Server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		Storage.Dump()
		err := server.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	// Запускаем сервер
	server = &http.Server{Addr: ServerConfig.Address, Handler: r}
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func appRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(gzipMiddleware)
	prepareRoutes(r)
	return r
}

func periodicDump() {
	for {
		time.Sleep(time.Duration(ServerConfig.StoreInterval) * time.Second)
		Storage.Dump()
	}
}
