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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
示例 05: 解析事件日志

功能：
1. 查询历史事件日志
2. 解析 indexed 和非 indexed 参数
3. 过滤特定地址的事件
4. 批量处理历史数据

使用场景：
- 同步历史数据
- 分析链上活动
- 审计和合规

使用方式：
go run examples/05_parse_logs.go
*/

// ERC20 Transfer 事件 ABI
const transferEventABI = `[{
	"anonymous": false,
	"inputs": [
		{"indexed": true, "name": "from", "type": "address"},
		{"indexed": true, "name": "to", "type": "address"},
		{"indexed": false, "name": "value", "type": "uint256"}
	],
	"name": "Transfer",
	"type": "event"
}]`

// Transfer 事件结构
type TransferEvent struct {
	From        common.Address
	To          common.Address
	Value       *big.Int
	TxHash      string
	BlockNumber uint64
}

func main() {
	fmt.Println("=== 示例 05: 解析事件日志 ===\n")

	// 连接到主网
	rpcURL := "https://mainnet.infura.io/v3/YOUR_API_KEY" // 替换为你的 API Key
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// USDT 合约地址
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// 解析合约 ABI
	contractABI, err := abi.JSON(strings.NewReader(transferEventABI))
	if err != nil {
		log.Fatalf("解析 ABI 失败: %v", err)
	}

	// ==== 示例 1: 查询指定区块范围的事件 ====
	fmt.Println("=== 示例 1: 查询最近 10 个区块的 Transfer 事件 ===\n")
	
	// 获取最新区块号
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("获取最新区块失败: %v", err)
	}

	fromBlock := latestBlock - 10
	toBlock := latestBlock

	fmt.Printf("查询区块范围: %d - %d\n", fromBlock, toBlock)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{usdtAddress},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Fatalf("查询日志失败: %v", err)
	}

	fmt.Printf("找到 %d 个事件\n\n", len(logs))

	// 解析前 5 个事件（避免输出太多）
	maxDisplay := 5
	if len(logs) > maxDisplay {
		fmt.Printf("显示前 %d 个事件:\n\n", maxDisplay)
	}

	for i, vLog := range logs {
		if i >= maxDisplay {
			break
		}

		event, err := parseTransferEvent(vLog, contractABI)
		if err != nil {
			log.Printf("解析事件失败: %v", err)
			continue
		}

		displayEvent(event, i+1)
	}

	// ==== 示例 2: 查询特定地址的转入事件 ====
	fmt.Println("\n=== 示例 2: 查询 Binance 热钱包的转入事件 ===\n")

	binanceAddress := common.HexToAddress("0x28C6c06298d514Db089934071355E5743bf21d60")
	
	// Transfer 事件的 topic0 是事件签名的 hash
	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferTopic := crypto.Keccak256Hash(transferEventSignature)

	// 构建查询：to = binanceAddress
	// topic1 = from (indexed)
	// topic2 = to (indexed)
	queryToAddress := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{usdtAddress},
		Topics: [][]common.Hash{
			{transferTopic},                       // topic0: 事件签名
			nil,                                    // topic1: from（任意地址）
			{common.BytesToHash(binanceAddress.Bytes())}, // topic2: to（目标地址）
		},
	}

	logsToAddress, err := client.FilterLogs(ctx, queryToAddress)
	if err != nil {
		log.Fatalf("查询日志失败: %v", err)
	}

	fmt.Printf("找到 %d 个转入事件\n\n", len(logsToAddress))

	for i, vLog := range logsToAddress {
		if i >= 3 { // 只显示前 3 个
			break
		}

		event, err := parseTransferEvent(vLog, contractABI)
		if err != nil {
			continue
		}

		displayEvent(event, i+1)
	}

	// ==== 示例 3: 查询特定地址的转出事件 ====
	fmt.Println("\n=== 示例 3: 查询 Binance 热钱包的转出事件 ===\n")

	queryFromAddress := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{usdtAddress},
		Topics: [][]common.Hash{
			{transferTopic},                       // topic0: 事件签名
			{common.BytesToHash(binanceAddress.Bytes())}, // topic1: from（目标地址）
			nil,                                    // topic2: to（任意地址）
		},
	}

	logsFromAddress, err := client.FilterLogs(ctx, queryFromAddress)
	if err != nil {
		log.Fatalf("查询日志失败: %v", err)
	}

	fmt.Printf("找到 %d 个转出事件\n\n", len(logsFromAddress))

	for i, vLog := range logsFromAddress {
		if i >= 3 {
			break
		}

		event, err := parseTransferEvent(vLog, contractABI)
		if err != nil {
			continue
		}

		displayEvent(event, i+1)
	}

	// ==== 示例 4: 统计分析 ====
	fmt.Println("\n=== 示例 4: 统计分析 ===\n")

	totalTransfers := len(logs)
	totalValue := big.NewInt(0)

	for _, vLog := range logs {
		event, err := parseTransferEvent(vLog, contractABI)
		if err != nil {
			continue
		}
		totalValue.Add(totalValue, event.Value)
	}

	fmt.Printf("总交易数: %d\n", totalTransfers)
	fmt.Printf("总交易量: %s USDT\n", formatUSDT(totalValue))
	
	if totalTransfers > 0 {
		avgValue := new(big.Int).Div(totalValue, big.NewInt(int64(totalTransfers)))
		fmt.Printf("平均金额: %s USDT\n", formatUSDT(avgValue))
	}
}

// parseTransferEvent 解析 Transfer 事件
func parseTransferEvent(vLog types.Log, contractABI abi.ABI) (*TransferEvent, error) {
	event := &TransferEvent{
		TxHash:      vLog.TxHash.Hex(),
		BlockNumber: vLog.BlockNumber,
	}

	// 解析 indexed 参数（在 Topics 中）
	if len(vLog.Topics) < 3 {
		return nil, fmt.Errorf("topics 数量不足")
	}

	// topics[0] 是事件签名
	// topics[1] 是 from (indexed)
	// topics[2] 是 to (indexed)
	event.From = common.HexToAddress(vLog.Topics[1].Hex())
	event.To = common.HexToAddress(vLog.Topics[2].Hex())

	// 解析非 indexed 参数（在 Data 中）
	// 方法 1: 使用 ABI 自动解析
	err := contractABI.UnpackIntoInterface(&event, "Transfer", vLog.Data)
	if err != nil {
		return nil, fmt.Errorf("解析 data 失败: %v", err)
	}

	// 方法 2: 手动解析 (备选)
	// value := new(big.Int).SetBytes(vLog.Data)
	// event.Value = value

	return event, nil
}

// displayEvent 显示事件信息
func displayEvent(event *TransferEvent, index int) {
	fmt.Printf("事件 #%d:\n", index)
	fmt.Printf("  From:        %s\n", event.From.Hex())
	fmt.Printf("  To:          %s\n", event.To.Hex())
	fmt.Printf("  Value:       %s USDT\n", formatUSDT(event.Value))
	fmt.Printf("  TxHash:      %s\n", event.TxHash)
	fmt.Printf("  BlockNumber: %d\n", event.BlockNumber)
	fmt.Println()
}

// formatUSDT 格式化 USDT 金额（6 位小数）
func formatUSDT(value *big.Int) string {
	fValue := new(big.Float).SetInt(value)
	usdt := new(big.Float).Quo(fValue, big.NewFloat(1e6))
	return usdt.Text('f', 2)
}

/*
知识点：

1. Event Topics 结构
   ┌────────────────────────────────────┐
   │ topics[0]: 事件签名 (Keccak256)     │
   │ topics[1]: 第 1 个 indexed 参数     │
   │ topics[2]: 第 2 个 indexed 参数     │
   │ topics[3]: 第 3 个 indexed 参数     │
   │ data:      非 indexed 参数           │
   └────────────────────────────────────┘

2. indexed vs 非 indexed
   - indexed: 可以作为过滤条件，存储在 topics 中
   - 非 indexed: 不能过滤，存储在 data 中
   - 最多 3 个 indexed 参数（topic0 是事件签名）

3. 事件签名计算
   - Transfer(address,address,uint256)
   - Keccak256 hash
   - 结果：0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef

4. 批量查询建议
   - 一次查询不超过 1000 个区块
   - 使用分页处理大量数据
   - 考虑使用专业索引服务（The Graph）

5. 性能优化
   - 使用 topics 过滤减少数据量
   - 缓存已解析的事件
   - 并发处理多个查询

使用场景：

1. 历史数据同步
   - 服务启动时同步历史数据
   - 填补丢失的事件

2. 链上分析
   - 统计交易量
   - 分析用户行为
   - 监控大额转账

3. 审计和合规
   - 追踪特定地址活动
   - 生成交易报告
   - 资金流向分析

进阶主题：
- 多个合约的批量查询
- 使用 The Graph 进行复杂查询
- 事件数据的存储和索引
- 实时监听 + 历史同步
*/
