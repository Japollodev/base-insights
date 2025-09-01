package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server                ServerConfig     `mapstructure:"server"`
	Database              DatabaseConfig   `mapstructure:"database"`
	Redis                 RedisConfig      `mapstructure:"redis"`
	Blockchain            BlockchainConfig `mapstructure:"blockchain"`
	DataCollection        DataCollectionConfig `mapstructure:"data_collection"`
	LogLevel              string           `mapstructure:"log_level"`
	GracefulShutdownTimeout time.Duration  `mapstructure:"graceful_shutdown_timeout"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type BlockchainConfig struct {
	RPCURL       string `mapstructure:"rpc_url"`
	ChainID      int64  `mapstructure:"chain_id"`
	ExplorerURL  string `mapstructure:"explorer_url"`
}

type DataCollectionConfig struct {
	BlockScanInterval      time.Duration `mapstructure:"block_scan_interval"`
	TransactionScanInterval time.Duration `mapstructure:"transaction_scan_interval"`
	PriceUpdateInterval    time.Duration `mapstructure:"price_update_interval"`
	BatchSize              int           `mapstructure:"batch_size"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// Set defaults
	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("server.read_timeout", 30*time.Second)
	viper.SetDefault("server.write_timeout", 30*time.Second)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("blockchain.rpc_url", "https://mainnet.base.org")
	viper.SetDefault("blockchain.chain_id", 8453)
	viper.SetDefault("data_collection.block_scan_interval", 5*time.Second)
	viper.SetDefault("data_collection.transaction_scan_interval", 3*time.Second)
	viper.SetDefault("data_collection.price_update_interval", 1*time.Minute)
	viper.SetDefault("data_collection.batch_size", 10)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("graceful_shutdown_timeout", 30*time.Second)

	// Environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional, continue with defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
