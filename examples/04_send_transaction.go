package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
示例 04: 发送交易

功能：
1. 发送 ETH 到指定地址
2. 签名交易
3. 广播交易
4. 等待交易确认

⚠️ 警告：
- 这个示例涉及私钥操作，仅用于测试网！
- 永远不要在代码中硬编码私钥
- 生产环境应使用 KMS 或硬件钱包

使用方式：
go run examples/04_send_transaction.go

前提条件：
- 有测试网账户并充值了测试币
- 设置环境变量 PRIVATE_KEY（不带 0x 前缀）
*/

func main() {
	fmt.Println("=== 示例 04: 发送交易 ===\n")

	// 连接到测试网（Sepolia）
	// 注意：使用测试网，不要在主网测试！
	rpcURL := "https://sepolia.infura.io/v3/YOUR_API_KEY" // 替换为你的 API Key
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 检查连接
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("获取网络 ID 失败: %v", err)
	}
	fmt.Printf("连接到网络，Chain ID: %s\n\n", chainID.String())

	// 1. 加载私钥
	// ⚠️ 实际使用时应从环境变量或安全存储读取
	privateKeyHex := "your_private_key_here" // 不带 0x 前缀
	
	// 示例：如何生成私钥（仅用于测试）
	// privateKey, err := crypto.GenerateKey()
	
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("私钥加载失败: %v\n提示：请替换为你的测试网私钥", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("获取公钥失败")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("发送地址: %s\n", fromAddress.Hex())

	// 2. 查询账户余额
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatalf("查询余额失败: %v", err)
	}
	fmt.Printf("当前余额: %s ETH\n\n", weiToEther(balance))

	// 检查余额是否足够
	if balance.Cmp(big.NewInt(0)) == 0 {
		log.Fatal("❌ 余额不足，请先从水龙头获取测试币")
	}

	// 3. 准备交易
	toAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e") // 示例地址
	value := big.NewInt(10000000000000000) // 0.01 ETH (以 Wei 为单位)
	
	// 获取 nonce（交易序号）
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}
	fmt.Printf("Nonce: %d\n", nonce)

	// 获取建议的 Gas Price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("获取 Gas Price 失败: %v", err)
	}
	fmt.Printf("Gas Price: %s Gwei\n", weiToGwei(gasPrice))

	// Gas Limit（对于简单的 ETH 转账，21000 是标准值）
	gasLimit := uint64(21000)

	// 4. 创建交易
	tx := types.NewTransaction(
		nonce,
		toAddress,
		value,
		gasLimit,
		gasPrice,
		nil, // data（对于 ETH 转账为空）
	)

	// 5. 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("签名交易失败: %v", err)
	}

	// 6. 发送交易
	fmt.Println("\n=== 发送交易 ===")
	fmt.Printf("接收地址: %s\n", toAddress.Hex())
	fmt.Printf("转账金额: %s ETH\n", weiToEther(value))
	fmt.Printf("Gas Limit: %d\n", gasLimit)
	fmt.Printf("Gas Price: %s Gwei\n", weiToGwei(gasPrice))
	fmt.Printf("预计费用: %s ETH\n", weiToEther(new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)))
	fmt.Println()

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("发送交易失败: %v", err)
	}

	txHash := signedTx.Hash()
	fmt.Printf("✅ 交易已发送！\n")
	fmt.Printf("交易哈希: %s\n", txHash.Hex())
	fmt.Printf("区块浏览器: https://sepolia.etherscan.io/tx/%s\n\n", txHash.Hex())

	// 7. 等待交易确认（可选）
	fmt.Println("⏳ 等待交易确认...")
	receipt, err := waitForReceipt(client, txHash)
	if err != nil {
		log.Fatalf("获取交易收据失败: %v", err)
	}

	// 8. 显示交易结果
	fmt.Println("\n=== 交易确认 ===")
	fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Printf("状态: %d (1=成功, 0=失败)\n", receipt.Status)
	
	if receipt.Status == 1 {
		actualCost := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), gasPrice)
		fmt.Printf("实际费用: %s ETH\n", weiToEther(actualCost))
		fmt.Println("\n✅ 交易成功！")
	} else {
		fmt.Println("\n❌ 交易失败")
	}
}

// waitForReceipt 等待交易被打包
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	
	// 轮询交易收据
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		
		// 如果是 "not found" 错误，继续等待
		// 否则返回错误
		if err.Error() != "not found" {
			return nil, err
		}
		
		// 每 3 秒查询一次
		fmt.Print(".")
		// time.Sleep(3 * time.Second)
	}
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

/*
运行步骤：

1. 获取测试网 API Key：
   - https://infura.io/ 注册并获取 Sepolia API Key

2. 获取测试币：
   - https://sepoliafaucet.com/
   - 或者 https://faucet.quicknode.com/ethereum/sepolia

3. 准备私钥：
   - 使用 MetaMask 或其他钱包导出私钥
   - 设置环境变量：export PRIVATE_KEY=your_key
   - 或者在代码中直接替换（仅测试！）

4. 运行：
   go run examples/04_send_transaction.go

5. 在 Etherscan 查看交易：
   https://sepolia.etherscan.io/tx/YOUR_TX_HASH

安全提示：
1. ❌ 永远不要在生产代码中硬编码私钥
2. ❌ 不要将包含私钥的代码提交到 Git
3. ✅ 使用环境变量存储私钥
4. ✅ 生产环境使用 KMS（AWS KMS, Google Cloud KMS）
5. ✅ 考虑使用硬件钱包（Ledger, Trezor）

进阶主题：
- EIP-1559 交易（动态 Gas Fees）
- 合约交互交易（带 data 字段）
- 批量发送交易
- 私钥加密存储
- 多签钱包
*/
