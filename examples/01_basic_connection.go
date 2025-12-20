package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 基础示例：连接节点并查询区块号
// 使用方式: go run examples/01_basic_connection.go
func main() {
	fmt.Println("=== 示例 1: 基础连接 ===\n")

	// 使用公共免费节点（无需 API Key）
	rpcURL := "https://eth.llamarpc.com"

	fmt.Println("正在连接以太坊主网...")
	fmt.Println("节点地址:", rpcURL)

	// 连接以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("❌ 连接失败:", err)
	}
	defer client.Close()

	fmt.Println("✅ 连接成功！\n")

	// 获取最新区块号
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("❌ 获取区块号失败:", err)
	}

	fmt.Println("📦 最新区块号:", blockNumber)
	fmt.Println("\n🎉 完成！")
}
