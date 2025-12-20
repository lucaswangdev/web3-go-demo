package listener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"web3-go-demo/db"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20 Transfer 事件 ABI
const erc20ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

// EventListener 事件监听器
type EventListener struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	database        *db.Database
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewEventListener 创建新的事件监听器
func NewEventListener(wsURL, contractAddr string, database *db.Database) (*EventListener, error) {
	// 连接到 WebSocket 节点
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// 解析合约 ABI
	contractABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventListener{
		client:          client,
		contractAddress: common.HexToAddress(contractAddr),
		contractABI:     contractABI,
		database:        database,
		ctx:             ctx,
		cancel:          cancel,
	}, nil
}

// Start 开始监听事件
func (l *EventListener) Start() error {
	log.Printf("开始监听合约事件: %s", l.contractAddress.Hex())

	// 创建查询过滤器
	query := ethereum.FilterQuery{
		Addresses: []common.Address{l.contractAddress},
	}

	// 创建日志通道
	logs := make(chan types.Log)

	// 订阅事件
	sub, err := l.client.SubscribeFilterLogs(l.ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	log.Println("成功订阅事件，等待 Transfer 事件...")

	// 监听事件
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("订阅错误: %v", err)
				// 实现重连逻辑
				log.Println("尝试重新连接...")
				time.Sleep(5 * time.Second)
				if err := l.reconnect(); err != nil {
					log.Printf("重连失败: %v", err)
				}
				return

			case vLog := <-logs:
				if err := l.processTransferEvent(vLog); err != nil {
					log.Printf("处理事件失败: %v", err)
				}

			case <-l.ctx.Done():
				log.Println("停止监听事件")
				return
			}
		}
	}()

	return nil
}

// processTransferEvent 处理 Transfer 事件
func (l *EventListener) processTransferEvent(vLog types.Log) error {
	// Transfer 事件结构
	type TransferEvent struct {
		From  common.Address
		To    common.Address
		Value *big.Int
	}

	var event TransferEvent

	// 解析事件数据
	// indexed 参数在 Topics 中，非 indexed 参数在 Data 中
	if len(vLog.Topics) < 3 {
		return fmt.Errorf("invalid Transfer event topics")
	}

	event.From = common.HexToAddress(vLog.Topics[1].Hex())
	event.To = common.HexToAddress(vLog.Topics[2].Hex())

	// 解析 value（非 indexed 参数）
	if err := l.contractABI.UnpackIntoInterface(&event, "Transfer", vLog.Data); err != nil {
		return fmt.Errorf("failed to unpack event data: %w", err)
	}

	log.Printf(
		"📤 Transfer 事件 | From: %s | To: %s | Amount: %s | Block: %d | TxHash: %s",
		event.From.Hex(),
		event.To.Hex(),
		event.Value.String(),
		vLog.BlockNumber,
		vLog.TxHash.Hex(),
	)

	// 获取区块时间戳
	block, err := l.client.BlockByNumber(l.ctx, big.NewInt(int64(vLog.BlockNumber)))
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	// 保存交易记录
	transaction := &db.Transaction{
		TxHash:      vLog.TxHash.Hex(),
		FromAddress: event.From.Hex(),
		ToAddress:   event.To.Hex(),
		Amount:      event.Value.String(),
		BlockNumber: vLog.BlockNumber,
		Timestamp:   time.Unix(int64(block.Time()), 0),
	}

	if err := l.database.SaveTransaction(transaction); err != nil {
		// 如果是重复键错误，可以忽略
		if !strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("failed to save transaction: %w", err)
		}
		log.Printf("交易 %s 已存在，跳过", vLog.TxHash.Hex())
	} else {
		log.Printf("✅ 交易已保存: %s", vLog.TxHash.Hex())
	}

	// 更新余额
	if err := l.database.ProcessTransfer(
		event.From.Hex(),
		event.To.Hex(),
		event.Value.String(),
	); err != nil {
		return fmt.Errorf("failed to process transfer: %w", err)
	}

	log.Println("✅ 余额已更新")

	return nil
}

// reconnect 重新连接
func (l *EventListener) reconnect() error {
	// 关闭旧连接
	l.client.Close()

	// 尝试重新连接
	// 注意：这里需要保存原始的 WebSocket URL
	// 实际使用时应该从配置中获取
	return fmt.Errorf("reconnect not fully implemented")
}

// Stop 停止监听
func (l *EventListener) Stop() {
	log.Println("停止事件监听器...")
	l.cancel()
	l.client.Close()
}
