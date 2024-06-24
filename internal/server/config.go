package server

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
)

type Config struct {
	Address         string `json:"server"`
	FileStoragePath string `json:"dumpPath"`
	DatabaseDSN     string `json:"dsn"`
	StoreInterval   int    `json:"interval"`
	RestoreStore    bool   `json:"restore"`
}

const defaultStoreInterval = 300 // seconds

var ServerConfig = Config{}

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
	flag.StringVar(
		&ServerConfig.DatabaseDSN,
		"d",
		"",
		"Адрес базы данных PostgreSQL",
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
	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		ServerConfig.DatabaseDSN = envDatabase
	}

	ServerConfig.log()
	return nil
}

func (s *Config) log() {
	lg, err := json.Marshal(ServerConfig)
	if err != nil {
		logger.Info("error serializing config:", err)
	}
	logger.Info("config:", string(lg))
}
