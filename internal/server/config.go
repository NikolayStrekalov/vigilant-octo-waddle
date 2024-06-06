package server

type Config struct {
	Storage         StorageOperations
	Address         string
	FileStoragePath string
	StoreInterval   int
	RestoreStore    bool
}

var ServerConfig = Config{}

func (c *Config) IsSyncDump() bool {
	return c.StoreInterval == 0 && c.FileStoragePath != ""
}

func (c *Config) DumpEnabled() bool {
	return c.FileStoragePath != ""
}
