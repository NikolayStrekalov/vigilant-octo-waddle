package server

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Address         string
	FileStoragePath string
	DatabaseDSN     string
	StoreInterval   int
	RestoreStore    bool
}

const defaultStoreInterval = 300 // seconds

var ServerConfig = Config{}

func (c *Config) IsSyncDump() bool {
	return c.StoreInterval == 0 && c.FileStoragePath != ""
}

func (c *Config) DumpEnabled() bool {
	return c.FileStoragePath != "" && ServerConfig.StoreInterval > 0
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
	flag.StringVar(
		&ServerConfig.DatabaseDSN,
		"d",
		"postgres://postgres:123@localhost:5432/postgres",
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
	return nil
}
