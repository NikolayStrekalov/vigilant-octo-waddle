package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	chi "github.com/go-chi/chi/v5"
)

const defaultStoreInterval = 300 // seconds

func Start() {
	// Инициализируем логирование
	lgr := logger.InitLog()
	defer func() {
		if err := lgr.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			lgr.Infoln("lgr err", err)
			panic(err)
		}
	}()

	// Настраиваем приложение
	r := appRouter()
	if err := fillConfig(); err != nil {
		flag.PrintDefaults()
		panic(err)
	}

	// Восстанавливаем данные
	if ServerConfig.RestoreStore {
		storage.Restore(ServerConfig.FileStoragePath)
	}
	// Периодически сохраняем данные
	if ServerConfig.DumpEnabled() && ServerConfig.StoreInterval > 0 {
		go periodicDump()
	}

	// При выходе сохраняем данные
	var server *http.Server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if ServerConfig.FileStoragePath != "" {
			storage.Dump(ServerConfig.FileStoragePath)
		}
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

func fillConfig() error {
	flag.StringVar(&ServerConfig.Address, "a", "localhost:8080", "Эндпоинт сервера HOST:PORT")
	flag.IntVar(
		&ServerConfig.StoreInterval,
		"i",
		defaultStoreInterval,
		"Интервал времени сохранения текущих значений (0 - синхронная запись)",
	)
	flag.StringVar(
		&ServerConfig.FileStoragePath,
		"f",
		"/tmp/metrics-db.json",
		"Полное имя файла, куда сохраняются текущие значения",
	)
	flag.BoolVar(
		&ServerConfig.RestoreStore,
		"r",
		true,
		"Загружать или нет ранее сохранённые значения из указанного файла",
	)
	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("too many args")
	}

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		ServerConfig.Address = envAddress
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return fmt.Errorf("can't parse STORE_INTERVAL: %w", err)
		}
		ServerConfig.StoreInterval = value
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		ServerConfig.FileStoragePath = envFileStoragePath
	}
	if envRestoreStore := os.Getenv("RESTORE"); envRestoreStore != "" {
		value, err := strconv.ParseBool(envRestoreStore)
		if err != nil {
			return fmt.Errorf("can't parse RESTORE: %w", err)
		}
		ServerConfig.RestoreStore = value
	}
	return nil
}

func periodicDump() {
	for {
		time.Sleep(time.Duration(ServerConfig.StoreInterval) * time.Second)
		storage.Dump(ServerConfig.FileStoragePath)
	}
}
