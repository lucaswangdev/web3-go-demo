package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"web3-go-demo/api"
	"web3-go-demo/config"
	"web3-go-demo/db"
	"web3-go-demo/listener"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 启动 ERC20 代币余额追踪服务...")

	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  无法加载 .env 文件: %v (将使用环境变量或默认值)", err)
	} else {
		log.Println("✅ .env 文件加载成功")
	}

	// 加载配置
	cfg := config.LoadConfig()
	log.Printf("配置加载完成: Contract=%s, Network=%s", cfg.ContractAddress, cfg.NetworkType)

	// 连接数据库
	database, err := db.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer database.Close()

	// 创建事件监听器
	eventListener, err := listener.NewEventListener(
		cfg.EthereumWSRPC,
		cfg.ContractAddress,
		database,
	)
	if err != nil {
		log.Fatalf("创建事件监听器失败: %v", err)
	}
	defer eventListener.Stop()

	// 启动事件监听
	if err := eventListener.Start(); err != nil {
		log.Fatalf("启动事件监听失败: %v", err)
	}

	// 创建并启动 API 服务器
	server := api.NewServer(database)

	// 在 goroutine 中启动 API 服务器
	go func() {
		if err := server.Start(cfg.APIPort); err != nil {
			log.Fatalf("API 服务器启动失败: %v", err)
		}
	}()

	log.Println("✅ 服务启动成功")
	log.Printf("📡 监听合约: %s", cfg.ContractAddress)
	log.Printf("🌐 API 服务: http://localhost:%s", cfg.APIPort)

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 正在关闭服务...")
}
