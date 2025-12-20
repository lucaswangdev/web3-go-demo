package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
示例 07: Nonce 管理

功能：
1. 理解 Nonce 的作用
2. 查询账户 Nonce
3. 并发发送交易的 Nonce 管理
4. 处理 Nonce 冲突和重试

使用场景：
- 高频交易
- 批量发送交易
- 交易加速和取消

使用方式：
go run examples/07_nonce_management.go
*/

// NonceManager Nonce 管理器
type NonceManager struct {
	client  *ethclient.Client
	address common.Address
	mu      sync.Mutex
	nonce   uint64
}

// NewNonceManager 创建 Nonce 管理器
func NewNonceManager(client *ethclient.Client, address common.Address) (*NonceManager, error) {
	ctx := context.Background()
	
	// 获取当前 nonce
	nonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, err
	}

	return &NonceManager{
		client:  client,
		address: address,
		nonce:   nonce,
	}, nil
}

// GetNonce 获取下一个可用的 nonce
func (nm *NonceManager) GetNonce() uint64 {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	currentNonce := nm.nonce
	nm.nonce++
	return currentNonce
}

// Reset 重置 nonce（从链上重新获取）
func (nm *NonceManager) Reset() error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	
	ctx := context.Background()
	nonce, err := nm.client.PendingNonceAt(ctx, nm.address)
	if err != nil {
		return err
	}
	
	nm.nonce = nonce
	return nil
}

// GetCurrentNonce 获取当前 nonce 值（不增加）
func (nm *NonceManager) GetCurrentNonce() uint64 {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	return nm.nonce
}

func main() {
	fmt.Println("=== 示例 07: Nonce 管理 ===\n")

	// 连接到测试网
	rpcURL := "https://sepolia.infura.io/v3/YOUR_API_KEY" // 替换为你的 API Key
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// ==== 示例 1: 理解 Nonce ====
	fmt.Println("=== 示例 1: 什么是 Nonce？ ===\n")
	
	fmt.Println("Nonce (Number used once) 是账户的交易序号：")
	fmt.Println("• 从 0 开始，每发送一笔交易 +1")
	fmt.Println("• 确保交易按顺序执行")
	fmt.Println("• 防止重放攻击")
	fmt.Println()

	// ==== 示例 2: 查询账户 Nonce ====
	fmt.Println("=== 示例 2: 查询账户 Nonce ===\n")

	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")

	// 方法 1: NonceAt - 已确认交易的 nonce
	confirmedNonce, err := client.NonceAt(ctx, address, nil)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("已确认的 Nonce: %d (下一笔已确认交易的序号)\n", confirmedNonce)
	}

	// 方法 2: PendingNonceAt - 包括待处理交易的 nonce
	pendingNonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("待处理的 Nonce: %d (下一笔交易应使用的序号)\n", pendingNonce)
	}

	if pendingNonce > confirmedNonce {
		fmt.Printf("\n💡 有 %d 笔交易在待处理队列中\n", pendingNonce-confirmedNonce)
	}

	// ==== 示例 3: Nonce 管理器 ====
	fmt.Println("\n=== 示例 3: 使用 Nonce 管理器 ===\n")

	manager, err := NewNonceManager(client, address)
	if err != nil {
		log.Printf("创建管理器失败: %v\n", err)
	} else {
		fmt.Printf("初始 Nonce: %d\n\n", manager.GetCurrentNonce())

		// 模拟获取多个 nonce
		fmt.Println("模拟批量获取 Nonce:")
		for i := 0; i < 5; i++ {
			nonce := manager.GetNonce()
			fmt.Printf("  交易 #%d: nonce = %d\n", i+1, nonce)
		}

		fmt.Printf("\n当前管理器 Nonce: %d\n", manager.GetCurrentNonce())
	}

	// ==== 示例 4: 并发安全的 Nonce 管理 ====
	fmt.Println("\n=== 示例 4: 并发安全测试 ===\n")

	manager2, err := NewNonceManager(client, address)
	if err != nil {
		log.Printf("创建管理器失败: %v\n", err)
	} else {
		var wg sync.WaitGroup
		nonceSet := make(map[uint64]bool)
		mu := sync.Mutex{}

		fmt.Println("并发获取 100 个 Nonce...")
		
		start := time.Now()
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				nonce := manager2.GetNonce()
				
				// 记录 nonce（检查是否有重复）
				mu.Lock()
				nonceSet[nonce] = true
				mu.Unlock()
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(start)

		fmt.Printf("✅ 完成！耗时: %v\n", elapsed)
		fmt.Printf("获取了 %d 个唯一的 Nonce\n", len(nonceSet))
		
		if len(nonceSet) == 100 {
			fmt.Println("✅ 测试通过：没有 Nonce 冲突")
		} else {
			fmt.Println("❌ 测试失败：存在 Nonce 冲突")
		}
	}

	// ==== 示例 5: Nonce 冲突处理 ====
	fmt.Println("\n=== 示例 5: Nonce 冲突场景 ===\n")

	fmt.Println("常见 Nonce 问题：")
	fmt.Println()
	
	fmt.Println("1. Nonce too low")
	fmt.Println("   原因：使用了已经使用过的 nonce")
	fmt.Println("   解决：重新获取最新的 pending nonce")
	fmt.Println()
	
	fmt.Println("2. Nonce too high")
	fmt.Println("   原因：跳过了某个 nonce")
	fmt.Println("   解决：确保 nonce 连续，或等待前面的交易完成")
	fmt.Println()
	
	fmt.Println("3. Replacement transaction underpriced")
	fmt.Println("   原因：用相同 nonce 发送了新交易，但 gas price 太低")
	fmt.Println("   解决：提高 gas price（至少提高 10%）")
	fmt.Println()

	// ==== 示例 6: 交易加速和取消 ====
	fmt.Println("=== 示例 6: 交易加速和取消 ===\n")

	fmt.Println("使用相同 Nonce 的技巧：")
	fmt.Println()
	
	fmt.Println("加速交易：")
	fmt.Println("  1. 使用原交易的相同 nonce")
	fmt.Println("  2. 提高 gas price (至少 +10%)")
	fmt.Println("  3. 保持其他参数不变")
	fmt.Println("  4. 重新签名并发送")
	fmt.Println()
	
	fmt.Println("取消交易：")
	fmt.Println("  1. 使用原交易的相同 nonce")
	fmt.Println("  2. 发送 0 ETH 给自己")
	fmt.Println("  3. 提高 gas price (至少 +10%)")
	fmt.Println("  4. 这会覆盖原交易")
	fmt.Println()

	// ==== 示例 7: 最佳实践 ====
	fmt.Println("=== 示例 7: Nonce 管理最佳实践 ===\n")

	fmt.Println("✅ 推荐做法：")
	fmt.Println("  1. 使用 PendingNonceAt 获取待处理的 nonce")
	fmt.Println("  2. 在内存中维护 nonce 计数器")
	fmt.Println("  3. 使用互斥锁保证并发安全")
	fmt.Println("  4. 定期与链上 nonce 同步")
	fmt.Println("  5. 记录每笔交易的 nonce")
	fmt.Println()
	
	fmt.Println("❌ 避免：")
	fmt.Println("  1. 每次都从链上查询 nonce（太慢）")
	fmt.Println("  2. 不处理并发场景")
	fmt.Println("  3. 不记录已发送的 nonce")
	fmt.Println("  4. 不处理 nonce 错误")
	fmt.Println()

	// ==== 示例 8: 生产级 Nonce 管理器 ====
	fmt.Println("=== 示例 8: 生产级功能建议 ===\n")

	fmt.Println("高级功能：")
	fmt.Println("  1. 持久化 nonce 状态（Redis/数据库）")
	fmt.Println("  2. 交易失败时的 nonce 回滚")
	fmt.Println("  3. 交易队列管理")
	fmt.Println("  4. 自动重试机制")
	fmt.Println("  5. 监控和告警")
	fmt.Println()

	fmt.Println("示例架构：")
	fmt.Println("  ┌─────────────────┐")
	fmt.Println("  │  应用层         │")
	fmt.Println("  └────────┬────────┘")
	fmt.Println("           │")
	fmt.Println("  ┌────────▼────────┐")
	fmt.Println("  │ Nonce Manager   │  ← 内存计数器 + 锁")
	fmt.Println("  └────────┬────────┘")
	fmt.Println("           │")
	fmt.Println("  ┌────────▼────────┐")
	fmt.Println("  │     Redis       │  ← 持久化状态")
	fmt.Println("  └────────┬────────┘")
	fmt.Println("           │")
	fmt.Println("  ┌────────▼────────┐")
	fmt.Println("  │   以太坊节点     │")
	fmt.Println("  └─────────────────┘")
	fmt.Println()
}

/*
知识点总结：

1. Nonce 规则
   ✓ 必须连续（0, 1, 2, 3...）
   ✓ 不能跳过
   ✓ 可以用相同 nonce 替换交易（提高 gas price）
   ✓ 区块链按 nonce 顺序执行交易

2. NonceAt vs PendingNonceAt
   ┌──────────────────┬─────────────────────┐
   │ NonceAt          │ 已确认交易的 nonce   │
   │ PendingNonceAt   │ 包括待处理交易       │
   └──────────────────┴─────────────────────┘
   
   建议：发送新交易时使用 PendingNonceAt

3. 并发场景
   问题：多个 goroutine 同时获取 nonce
   解决：使用 Mutex 或 Channel 保证原子性
   
   // ❌ 错误
   nonce, _ := client.PendingNonceAt(ctx, address)
   
   // ✅ 正确
   manager.GetNonce() // 线程安全

4. Nonce 错误处理
   错误类型                   解决方案
   ─────────────────────────────────────
   nonce too low              重新获取 nonce
   nonce too high             等待前面交易完成
   replacement underpriced    提高 gas price +10%

5. 交易替换
   用途：加速交易、取消交易
   方法：使用相同 nonce + 更高 gas price
   
   加速：发送相同交易 + 高 gas
   取消：发送 0 ETH 给自己 + 高 gas

实战经验：

1. 高频交易
   • 预先分配 nonce 范围
   • 使用内存队列
   • 批量发送

2. 交易监控
   • 记录每笔交易的 nonce
   • 监控未确认交易
   • 检测 nonce gap

3. 错误恢复
   • 交易失败时回滚 nonce
   • 定期与链上同步
   • 实现重试逻辑

4. 性能优化
   • 减少链上查询
   • 使用本地计数器
   • 批量处理交易

代码示例参考：
- Web3.js: web3.eth.getTransactionCount()
- Ethers.js: provider.getTransactionCount()
- go-ethereum: client.PendingNonceAt()
*/
