package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	fmt.Println("=== 本地区块链测试示例 ===")
	
	// 连接到本地 Anvil 节点
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal("连接本地节点失败:", err)
	}
	defer client.Close()
	
	// 检查连接
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("获取链ID失败:", err)
	}
	fmt.Printf("连接成功! Chain ID: %d\n", chainID)
	
	// 使用 Anvil 默认的测试账户
	testAccounts := []struct {
		address    string
		privateKey string
	}{
		{
			address:    "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			privateKey: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		{
			address:    "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			privateKey: "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
	}
	
	// 检查账户余额
	fmt.Println("\n=== 账户余额 ===")
	for i, account := range testAccounts {
		address := common.HexToAddress(account.address)
		balance, err := client.BalanceAt(context.Background(), address, nil)
		if err != nil {
			log.Printf("获取账户 %d 余额失败: %v", i, err)
			continue
		}
		
		// 转换为 ETH (Wei to ETH)
		ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
		fmt.Printf("账户 %d (%s): %s ETH\n", i, account.address, ethBalance.String())
	}
	
	// 发送测试交易
	fmt.Println("\n=== 发送测试交易 ===")
	err = sendTestTransaction(client, testAccounts[0], testAccounts[1].address)
	if err != nil {
		log.Printf("发送交易失败: %v", err)
	}
	
	// 获取最新区块信息
	fmt.Println("\n=== 最新区块信息 ===")
	err = getLatestBlockInfo(client)
	if err != nil {
		log.Printf("获取区块信息失败: %v", err)
	}
	
	// 监听新区块 (演示几秒钟)
	fmt.Println("\n=== 监听新区块 (5秒) ===")
	err = listenToNewBlocks(client)
	if err != nil {
		log.Printf("监听区块失败: %v", err)
	}
}

// sendTestTransaction 发送测试交易
func sendTestTransaction(client *ethclient.Client, from struct {
	address    string
	privateKey string
}, toAddress string) error {
	
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(from.privateKey)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %v", err)
	}
	
	// 获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("获取公钥失败")
	}
	
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	
	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("获取nonce失败: %v", err)
	}
	
	// 设置交易参数
	value := big.NewInt(1000000000000000000) // 1 ETH in wei
	gasLimit := uint64(21000)                // 标准转账的gas限制
	
	// 获取建议的gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("获取gas价格失败: %v", err)
	}
	
	// 获取链ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("获取链ID失败: %v", err)
	}
	
	// 创建交易
	toAddr := common.HexToAddress(toAddress)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, nil)
	
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("签名交易失败: %v", err)
	}
	
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("发送交易失败: %v", err)
	}
	
	fmt.Printf("交易已发送! Hash: %s\n", signedTx.Hash().Hex())
	fmt.Printf("从 %s 向 %s 发送 1 ETH\n", fromAddress.Hex(), toAddress)
	
	return nil
}

// getLatestBlockInfo 获取最新区块信息
func getLatestBlockInfo(client *ethclient.Client) error {
	// 获取最新区块号
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("获取区块号失败: %v", err)
	}
	
	// 获取区块详情
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return fmt.Errorf("获取区块详情失败: %v", err)
	}
	
	fmt.Printf("最新区块号: %d\n", blockNumber)
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("区块时间: %d\n", block.Time())
	fmt.Printf("Gas使用量: %d\n", block.GasUsed())
	fmt.Printf("Gas限制: %d\n", block.GasLimit())
	
	return nil
}

// listenToNewBlocks 监听新区块 (HTTP轮询方式)
func listenToNewBlocks(client *ethclient.Client) error {
	fmt.Println("开始监听新区块 (HTTP轮询)...")
	
	// 获取初始区块号
	lastBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("获取初始区块号失败: %v", err)
	}
	
	// 轮询5秒
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			fmt.Println("监听结束")
			return nil
		case <-ticker.C:
			currentBlock, err := client.BlockNumber(context.Background())
			if err != nil {
				continue
			}
			
			if currentBlock > lastBlock {
				// 获取新区块详情
				block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(currentBlock)))
				if err != nil {
					continue
				}
				
				fmt.Printf("新区块: #%d, Hash: %s, 交易数: %d\n", 
					currentBlock, 
					block.Hash().Hex()[:10]+"...", 
					len(block.Transactions()))
				lastBlock = currentBlock
			}
		}
	}
}