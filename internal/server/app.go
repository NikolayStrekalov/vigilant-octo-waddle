package server

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/memstorage"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/pgstorage"

	chi "github.com/go-chi/chi/v5"
)

func Start() {
	var err error
	// Инициализируем логирование
	lgr := logger.InitLog()
	defer func() {
		_ = lgr.Sync()
	}()

	// Настраиваем приложение
	if err := fillConfig(); err != nil {
		flag.PrintDefaults()
		panic(err)
	}
	var storageClose func() error
	if ServerConfig.DatabaseDSN == "" {
		Storage, storageClose, err = memstorage.NewMemStorage(
			ServerConfig.FileStoragePath,
			ServerConfig.RestoreStore,
			ServerConfig.StoreInterval,
		)
	} else {
		// TODO: повторная инициализация при недоступности базы
		Storage, storageClose, err = pgstorage.NewPGStorage(
			ServerConfig.DatabaseDSN,
		)
	}
	if err != nil {
		logger.Info("error creating storage:", err)
	}
	defer func() {
		if storageClose != nil {
			_ = storageClose()
		}
	}()
	r := appRouter()

	// Дожидаемся выхода из этой функции
	var server *http.Server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		err := server.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	// Запускаем сервер
	server = &http.Server{Addr: ServerConfig.Address, Handler: r}
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func appRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	if ServerConfig.SignKey != "" {
		r.Use(shaMiddlewareBuilder(ServerConfig.SignKey))
	}
	r.Use(gzipMiddleware)
	prepareRoutes(r)
	return r
}
