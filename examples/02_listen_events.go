package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20 Transfer 事件结构
type TransferEvent struct {
	From   common.Address
	To     common.Address
	Value  *big.Int
	TxHash string
	Block  uint64
}

// USDT 合约地址（以太坊主网）
const USDTAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

// ERC20 Transfer 事件 ABI
const TransferEventABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

func main() {
	fmt.Println("=== 示例 2: 监听 USDT Transfer 事件 ===\n")

	// ==================== 连接节点 ====================
	// 注意：监听事件建议使用 WebSocket 连接
	// 但公共 HTTP 节点也可以通过轮询方式使用
	rpcURL := "https://eth.llamarpc.com"

	fmt.Println("连接节点:", rpcURL)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Close()
	fmt.Println("✅ 连接成功\n")

	// ==================== 获取最新区块 ====================
	latestBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal("获取区块号失败:", err)
	}
	fmt.Printf("📦 当前区块: %d\n\n", latestBlock)

	// ==================== 查询历史 Transfer 事件 ====================
	// 从最近 100 个区块中查询 USDT Transfer 事件
	fromBlock := latestBlock - 100
	toBlock := latestBlock

	fmt.Printf("🔍 查询区块范围: %d - %d (最近 100 个区块)\n", fromBlock, toBlock)
	fmt.Println("📍 USDT 合约地址:", USDTAddress)
	fmt.Println("⏳ 正在查询事件...\n")

	// 构建过滤器
	contractAddress := common.HexToAddress(USDTAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{contractAddress},
	}

	// 执行查询
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal("查询事件失败:", err)
	}

	fmt.Printf("✅ 找到 %d 个 Transfer 事件\n", len(logs))
	fmt.Println(strings.Repeat("=", 80))

	// ==================== 解析事件 ====================
	// 解析 ABI
	contractAbi, err := abi.JSON(strings.NewReader(TransferEventABI))
	if err != nil {
		log.Fatal("解析 ABI 失败:", err)
	}

	// 显示前 10 个事件
	displayCount := 10
	if len(logs) < displayCount {
		displayCount = len(logs)
	}

	fmt.Printf("\n显示最近 %d 笔 USDT 转账:\n\n", displayCount)

	for i := 0; i < displayCount; i++ {
		vLog := logs[i]
		event, err := parseTransferEvent(vLog, &contractAbi)
		if err != nil {
			fmt.Printf("❌ 解析事件失败: %v\n", err)
			continue
		}

		// 将 USDT 从最小单位转换为实际金额（USDT 使用 6 位小数）
		valueFloat := new(big.Float).Quo(
			new(big.Float).SetInt(event.Value),
			big.NewFloat(1000000), // USDT 精度为 6
		)

		fmt.Printf("📝 转账 #%d\n", i+1)
		fmt.Printf("   区块: %d\n", event.Block)
		fmt.Printf("   交易: %s\n", event.TxHash)
		fmt.Printf("   从:   %s\n", event.From.Hex())
		fmt.Printf("   到:   %s\n", event.To.Hex())
		fmt.Printf("   金额: %.2f USDT\n", valueFloat)
		fmt.Println(strings.Repeat("-", 80))
	}

	// ==================== 统计信息 ====================
	fmt.Println("\n📊 统计信息:")
	fmt.Printf("  - 查询区块数: %d\n", toBlock-fromBlock+1)
	fmt.Printf("  - Transfer 事件数: %d\n", len(logs))
	fmt.Printf("  - 平均每区块: %.2f 笔\n", float64(len(logs))/float64(toBlock-fromBlock+1))

	// ==================== 下一步提示 ====================
	fmt.Println("\n💡 本示例演示:")
	fmt.Println("  ✅ 使用 FilterQuery 查询历史事件")
	fmt.Println("  ✅ 解析 Transfer 事件参数")
	fmt.Println("  ✅ 处理 ERC20 代币精度")

	fmt.Println("\n📚 进阶学习:")
	fmt.Println("  1. 使用 WebSocket 实时监听事件")
	fmt.Println("  2. 添加事件过滤器（只监听特定地址）")
	fmt.Println("  3. 将事件数据存储到数据库")
}

// parseTransferEvent 解析 Transfer 事件
func parseTransferEvent(vLog types.Log, contractAbi *abi.ABI) (*TransferEvent, error) {
	event := &TransferEvent{
		TxHash: vLog.TxHash.Hex(),
		Block:  vLog.BlockNumber,
	}

	// Transfer 事件的 Topics:
	// topics[0]: 事件签名 (keccak256("Transfer(address,address,uint256)"))
	// topics[1]: from address (indexed)
	// topics[2]: to address (indexed)
	// Data: value (not indexed)

	if len(vLog.Topics) < 3 {
		return nil, fmt.Errorf("无效的 Transfer 事件 topics")
	}

	// 解析 from (topics[1])
	event.From = common.HexToAddress(vLog.Topics[1].Hex())

	// 解析 to (topics[2])
	event.To = common.HexToAddress(vLog.Topics[2].Hex())

	// 解析 value (Data)
	var transferEvent struct {
		Value *big.Int
	}

	err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
	if err != nil {
		// 如果 UnpackIntoInterface 失败，尝试直接读取 Data
		if len(vLog.Data) > 0 {
			event.Value = new(big.Int).SetBytes(vLog.Data)
		} else {
			return nil, err
		}
	} else {
		event.Value = transferEvent.Value
	}

	return event, nil
}
