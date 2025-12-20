package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== Web3 Go 进阶示例：带网络诊断和重试 ===\n")

	// ==================== 节点服务商配置 ====================
	nodes := []string{
		"https://eth.llamarpc.com",         // 公共免费节点（无需 API Key）
		"https://cloudflare-eth.com",       // Cloudflare 节点（无需 API Key）
		"https://mainnet.infura.io/v3/xxx", // 你的 Infura
		"https://rpc.ankr.com/eth",         // Ankr 公共节点
	}

	fmt.Println("📡 可用节点列表:")
	for i, node := range nodes {
		fmt.Printf("  %d. %s\n", i+1, node)
	}
	fmt.Println()

	// ==================== 测试网络连接 ====================
	fmt.Println("🔍 测试网络连接...")
	testHTTP()

	// ==================== 尝试连接节点 ====================
	var client *ethclient.Client
	var connectedNode string
	var err error

	for i, nodeURL := range nodes {
		fmt.Printf("\n尝试连接节点 %d/%d: %s\n", i+1, len(nodes), nodeURL)

		client, err = connectWithRetry(nodeURL, 2)
		if err == nil {
			connectedNode = nodeURL
			fmt.Println("✅ 连接成功！")
			break
		}

		fmt.Printf("❌ 连接失败: %v\n", err)
	}

	if client == nil {
		log.Fatal("\n❌ 所有节点都无法连接，请检查网络设置")
	}
	defer client.Close()

	fmt.Printf("\n✅ 当前使用节点: %s\n\n", connectedNode)

	// ==================== 获取链信息 ====================
	fmt.Println("📊 正在获取区块链信息...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. 获取最新区块号
	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal("❌ 获取区块号失败:", err)
	}
	fmt.Printf("📦 最新区块号: %d\n", blockNumber)

	// 2. 获取链 ID
	chainID, err := client.ChainID(ctx)
	if err == nil {
		fmt.Printf("🔗 链 ID: %d ", chainID.Int64())
		switch chainID.Int64() {
		case 1:
			fmt.Println("(以太坊主网)")
		case 5:
			fmt.Println("(Goerli 测试网)")
		case 11155111:
			fmt.Println("(Sepolia 测试网)")
		default:
			fmt.Println()
		}
	}

	// 3. 查询 Vitalik 的余额
	vitalikAddress := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	balance, err := client.BalanceAt(ctx, vitalikAddress, nil)
	if err == nil {
		ethBalance := new(big.Float).Quo(
			new(big.Float).SetInt(balance),
			new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
		)
		fmt.Printf("💰 Vitalik 地址余额: %.4f ETH\n", ethBalance)
	}

	// ==================== 成功信息 ====================
	fmt.Println("\n============================================================")
	fmt.Println("🎉 恭喜！你已经成功完成以下操作:")
	fmt.Println("  ✅ 连接到以太坊主网")
	fmt.Println("  ✅ 查询最新区块号")
	fmt.Println("  ✅ 查询账户余额")
	fmt.Println("============================================================")

	fmt.Println("\n📚 下一步学习:")
	fmt.Println("  1. 监听 USDT 合约的 Transfer 事件")
	fmt.Println("  2. 调用智能合约方法")
	fmt.Println("  3. 在测试网发送交易")
}

// connectWithRetry 带重试的连接函数
func connectWithRetry(nodeURL string, maxRetries int) (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = ethclient.DialContext(ctx, nodeURL)
		cancel()

		if err == nil {
			return client, nil
		}

		if i < maxRetries-1 {
			time.Sleep(1 * time.Second)
		}
	}

	return nil, err
}

// testHTTP 测试基础 HTTP 连接
func testHTTP() {
	testURL := "https://www.google.com"
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(testURL)
	if err != nil {
		fmt.Printf("⚠️  HTTP 测试失败: %v\n", err)
		fmt.Println("   可能原因：网络问题或需要代理")
	} else {
		resp.Body.Close()
		fmt.Printf("✅ HTTP 测试通过 (状态码: %d)\n", resp.StatusCode)
	}
}
