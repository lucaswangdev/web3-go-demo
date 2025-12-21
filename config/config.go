package config

import (
	"fmt"
	"os"
)

// Config 存储应用程序配置
type Config struct {
	// Ethereum 配置
	EthereumRPC        string
	EthereumWSRPC      string
	ContractAddress    string
	StartBlock         uint64
	NetworkType        string // mainnet, testnet, local
	
	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	// API 配置
	APIPort    string
	
	// 日志配置
	LogLevel   string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	networkType := getEnv("NETWORK_TYPE", "mainnet")
	
	config := &Config{
		NetworkType:     networkType,
		ContractAddress: getEnv("CONTRACT_ADDRESS", "0xdAC17F958D2ee523a2206206994597C13D831ec7"), // USDT 默认
		StartBlock:      0,
		
		// 数据库配置
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "erc20_tracker"),
		
		// API 配置
		APIPort:    getEnv("API_PORT", "8080"),
		
		// 日志配置
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
	
	// 根据网络类型设置 RPC 端点
	switch networkType {
	case "local":
		config.EthereumRPC = getEnv("ETHEREUM_RPC", "http://127.0.0.1:8545")
		config.EthereumWSRPC = getEnv("ETHEREUM_WS_RPC", "ws://127.0.0.1:8545")
	case "testnet":
		config.EthereumRPC = getEnv("ETHEREUM_RPC", "https://sepolia.infura.io/v3/YOUR_API_KEY")
		config.EthereumWSRPC = getEnv("ETHEREUM_WS_RPC", "wss://sepolia.infura.io/ws/v3/YOUR_API_KEY")
	default: // mainnet
		config.EthereumRPC = getEnv("ETHEREUM_RPC", "https://mainnet.infura.io/v3/YOUR_API_KEY")
		config.EthereumWSRPC = getEnv("ETHEREUM_WS_RPC", "wss://mainnet.infura.io/ws/v3/YOUR_API_KEY")
	}
	
	return config
}

// GetDSN 返回 PostgreSQL 连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

// getEnv 从环境变量获取值，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
