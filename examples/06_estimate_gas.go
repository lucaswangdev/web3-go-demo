package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
示例 06: Gas 费估算

功能：
1. 估算简单 ETH 转账的 Gas
2. 估算合约调用的 Gas
3. 获取当前 Gas Price
4. 计算交易总成本
5. EIP-1559 Gas 费机制

使用场景：
- 在发送交易前估算成本
- 优化 Gas 费用
- 提供用户友好的费用显示

使用方式：
go run examples/06_estimate_gas.go
*/

func main() {
	fmt.Println("=== 示例 06: Gas 费估算 ===\n")

	// 连接到主网
	rpcURL := "https://mainnet.infura.io/v3/YOUR_API_KEY" // 替换为你的 API Key
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 检查网络
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("获取网络 ID 失败: %v", err)
	}
	fmt.Printf("连接到网络，Chain ID: %s\n\n", chainID.String())

	// ==== 示例 1: 简单 ETH 转账的 Gas 估算 ====
	fmt.Println("=== 示例 1: ETH 转账 Gas 估算 ===\n")

	fromAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	toAddress := common.HexToAddress("0x28C6c06298d514Db089934071355E5743bf21d60")
	value := big.NewInt(1000000000000000000) // 1 ETH

	// 构建交易消息
	msg := ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: value,
		Data:  nil,
	}

	// 估算 Gas Limit
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		log.Fatalf("估算 Gas 失败: %v", err)
	}

	// 获取建议的 Gas Price (Legacy)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("获取 Gas Price 失败: %v", err)
	}

	// 计算总成本
	totalCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)

	fmt.Printf("转账金额:    %s ETH\n", weiToEther(value))
	fmt.Printf("Gas Limit:   %d\n", gasLimit)
	fmt.Printf("Gas Price:   %s Gwei\n", weiToGwei(gasPrice))
	fmt.Printf("Gas 成本:    %s ETH\n", weiToEther(totalCost))
	fmt.Printf("总成本:      %s ETH\n", weiToEther(new(big.Int).Add(value, totalCost)))

	// 标准 ETH 转账固定为 21000 Gas
	fmt.Println("\n💡 提示：标准 ETH 转账的 Gas Limit 固定为 21000")

	// ==== 示例 2: 合约调用 Gas 估算 ====
	fmt.Println("\n=== 示例 2: ERC20 转账 Gas 估算 ===\n")

	// USDT 合约地址
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// ERC20 transfer 方法的数据
	// transfer(address to, uint256 amount)
	// 方法 ID: 0xa9059cbb
	// 参数 1: to address (32 bytes)
	// 参数 2: amount (32 bytes)
	transferAmount := big.NewInt(1000000) // 1 USDT (6 位小数)
	
	// 构建 calldata
	methodID := []byte{0xa9, 0x05, 0x9c, 0xbb} // transfer(address,uint256)
	toAddressBytes := common.LeftPadBytes(toAddress.Bytes(), 32)
	amountBytes := common.LeftPadBytes(transferAmount.Bytes(), 32)
	
	data := append(methodID, toAddressBytes...)
	data = append(data, amountBytes...)

	// 构建合约调用消息
	contractMsg := ethereum.CallMsg{
		From: fromAddress,
		To:   &usdtAddress,
		Data: data,
	}

	// 估算 Gas
	contractGasLimit, err := client.EstimateGas(ctx, contractMsg)
	if err != nil {
		log.Printf("估算 Gas 失败: %v\n", err)
		log.Println("提示：可能是账户余额不足或合约调用失败")
	} else {
		contractGasCost := new(big.Int).Mul(big.NewInt(int64(contractGasLimit)), gasPrice)
		
		fmt.Printf("合约:        USDT\n")
		fmt.Printf("方法:        transfer\n")
		fmt.Printf("转账金额:    %s USDT\n", formatUSDT(transferAmount))
		fmt.Printf("Gas Limit:   %d\n", contractGasLimit)
		fmt.Printf("Gas Price:   %s Gwei\n", weiToGwei(gasPrice))
		fmt.Printf("Gas 成本:    %s ETH\n", weiToEther(contractGasCost))
		
		fmt.Println("\n💡 提示：ERC20 转账的 Gas 通常在 45,000 - 65,000 之间")
	}

	// ==== 示例 3: EIP-1559 Gas 费 ====
	fmt.Println("\n=== 示例 3: EIP-1559 Gas 费（伦敦升级后）===\n")

	// 获取建议的 Gas Tip Cap (优先费)
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		log.Printf("获取 Gas Tip Cap 失败: %v", err)
	} else {
		fmt.Printf("Max Priority Fee (Tip): %s Gwei\n", weiToGwei(gasTipCap))
		
		// 获取最新区块的 Base Fee
		block, err := client.BlockByNumber(ctx, nil)
		if err != nil {
			log.Printf("获取区块失败: %v", err)
		} else {
			if block.BaseFee() != nil {
				baseFee := block.BaseFee()
				
				// Max Fee = Base Fee * 2 + Priority Fee (常用策略)
				maxFeePerGas := new(big.Int).Mul(baseFee, big.NewInt(2))
				maxFeePerGas.Add(maxFeePerGas, gasTipCap)
				
				fmt.Printf("Base Fee:               %s Gwei\n", weiToGwei(baseFee))
				fmt.Printf("Max Fee Per Gas:        %s Gwei\n", weiToGwei(maxFeePerGas))
				
				// 计算 EIP-1559 交易的成本
				eip1559Cost := new(big.Int).Mul(big.NewInt(21000), maxFeePerGas)
				fmt.Printf("预计最大成本 (21000 Gas): %s ETH\n", weiToEther(eip1559Cost))
				
				// 实际成本（通常更低）
				actualFee := new(big.Int).Add(baseFee, gasTipCap)
				actualCost := new(big.Int).Mul(big.NewInt(21000), actualFee)
				fmt.Printf("预计实际成本:            %s ETH\n", weiToEther(actualCost))
			}
		}
	}

	// ==== 示例 4: Gas 价格建议 ====
	fmt.Println("\n=== 示例 4: Gas 价格策略 ===\n")

	fmt.Printf("当前 Gas Price: %s Gwei\n\n", weiToGwei(gasPrice))

	// 快速、标准、慢速策略
	fastGasPrice := new(big.Int).Mul(gasPrice, big.NewInt(120))
	fastGasPrice.Div(fastGasPrice, big.NewInt(100)) // * 1.2

	normalGasPrice := gasPrice

	slowGasPrice := new(big.Int).Mul(gasPrice, big.NewInt(80))
	slowGasPrice.Div(slowGasPrice, big.NewInt(100)) // * 0.8

	fmt.Println("Gas 价格策略建议：")
	fmt.Printf("  🚀 快速 (1-2 分钟): %s Gwei\n", weiToGwei(fastGasPrice))
	fmt.Printf("  ⚡ 标准 (3-5 分钟): %s Gwei\n", weiToGwei(normalGasPrice))
	fmt.Printf("  🐌 慢速 (>10 分钟): %s Gwei\n", weiToGwei(slowGasPrice))

	fmt.Println("\n对应转账成本 (21000 Gas):")
	fmt.Printf("  快速: %s ETH\n", weiToEther(new(big.Int).Mul(big.NewInt(21000), fastGasPrice)))
	fmt.Printf("  标准: %s ETH\n", weiToEther(new(big.Int).Mul(big.NewInt(21000), normalGasPrice)))
	fmt.Printf("  慢速: %s ETH\n", weiToEther(new(big.Int).Mul(big.NewInt(21000), slowGasPrice)))

	// ==== 示例 5: 批量操作的 Gas 优化 ====
	fmt.Println("\n=== 示例 5: 批量操作 Gas 对比 ===\n")

	numTransfers := 10
	singleTxGas := uint64(21000)
	batchTxGas := uint64(21000 + (numTransfers-1)*15000) // 估算批量操作节省

	fmt.Printf("单笔转账 x %d 次:\n", numTransfers)
	fmt.Printf("  Gas: %d\n", singleTxGas*uint64(numTransfers))
	fmt.Printf("  成本: %s ETH\n", weiToEther(new(big.Int).Mul(
		big.NewInt(int64(singleTxGas*uint64(numTransfers))), 
		gasPrice,
	)))

	fmt.Printf("\n批量转账 (1 次交易):\n")
	fmt.Printf("  Gas: %d\n", batchTxGas)
	fmt.Printf("  成本: %s ETH\n", weiToEther(new(big.Int).Mul(
		big.NewInt(int64(batchTxGas)), 
		gasPrice,
	)))

	saved := float64(singleTxGas*uint64(numTransfers)-batchTxGas) / float64(singleTxGas*uint64(numTransfers)) * 100
	fmt.Printf("\n💰 节省: %.1f%%\n", saved)
}

// weiToEther 将 Wei 转换为 ETH
func weiToEther(wei *big.Int) string {
	fwei := new(big.Float).SetInt(wei)
	ether := new(big.Float).Quo(fwei, big.NewFloat(1e18))
	return ether.Text('f', 6)
}

// weiToGwei 将 Wei 转换为 Gwei
func weiToGwei(wei *big.Int) string {
	fwei := new(big.Float).SetInt(wei)
	gwei := new(big.Float).Quo(fwei, big.NewFloat(1e9))
	return gwei.Text('f', 2)
}

// formatUSDT 格式化 USDT 金额
func formatUSDT(value *big.Int) string {
	fValue := new(big.Float).SetInt(value)
	usdt := new(big.Float).Quo(fValue, big.NewFloat(1e6))
	return usdt.Text('f', 2)
}

/*
知识点：

1. Gas 相关概念
   ┌──────────────────────────────────────────┐
   │ Gas Limit:  交易可以使用的最大 Gas 数量    │
   │ Gas Used:   交易实际使用的 Gas 数量        │
   │ Gas Price:  每单位 Gas 的价格 (Gwei)      │
   │ Gas Fee:    Gas Used * Gas Price          │
   └──────────────────────────────────────────┘

2. Legacy vs EIP-1559
   Legacy:
   - Gas Price: 固定价格
   - 总费用 = Gas Limit * Gas Price
   
   EIP-1559:
   - Base Fee: 基础费用（协议决定，销毁）
   - Priority Fee: 小费（给矿工）
   - Max Fee = Base Fee + Priority Fee
   - 总费用 = Gas Used * (Base Fee + Priority Fee)

3. 常见操作的 Gas 消耗
   - ETH 转账:        21,000
   - ERC20 transfer:  45,000 - 65,000
   - Uniswap swap:    100,000 - 150,000
   - NFT mint:        50,000 - 100,000

4. Gas 优化策略
   - 选择低流量时段（周末、夜间）
   - 使用 L2 方案（Arbitrum, Optimism）
   - 批量操作
   - 优化合约代码

5. Gas Price 单位
   1 ETH = 10^9 Gwei = 10^18 Wei
   
实用建议：

1. 估算 Gas 时加 20% buffer
   estimatedGas := estimatedGas * 120 / 100

2. 监控 Gas Price
   - https://etherscan.io/gastracker
   - https://www.blocknative.com/gas-estimator

3. 使用 EIP-1559
   - 更可预测的费用
   - 避免过高支付
   - 通常比 Legacy 便宜

4. 大额交易
   - 多等几个区块确认
   - 使用较低的 Gas Price
   - 节省费用

进阶主题：
- 动态 Gas Price 调整
- MEV (Maximal Extractable Value)
- Flashbots 私有交易
- Gas Token 套利
*/
