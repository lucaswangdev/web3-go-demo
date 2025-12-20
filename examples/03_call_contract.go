package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ==================== USDT 合约地址 ====================
const USDTAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"

// ==================== 一些知名地址 ====================
var (
	VitalikAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045") // Vitalik
	BinanceAddress = common.HexToAddress("0xF977814e90dA44bFA03b6295A0616a897441aceC") // Binance 热钱包
)

func main() {
	fmt.Println("=== 示例 3: 调用智能合约方法 ===\n")

	// ==================== 连接节点 ====================
	rpcURL := "https://eth.llamarpc.com"
	fmt.Println("连接节点:", rpcURL)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Close()
	fmt.Println("✅ 连接成功\n")

	// ==================== 调用合约（方式1：手动构造） ====================
	fmt.Println("📋 方式 1: 使用 CallContract 手动调用")
	fmt.Println(strings.Repeat("=", 80))

	contractAddress := common.HexToAddress(USDTAddress)

	// 查询 USDT 总供应量
	totalSupply, err := getTotalSupply(client, contractAddress)
	if err != nil {
		log.Printf("⚠️  查询总供应量失败: %v\n", err)
	} else {
		// USDT 精度为 6
		supply := new(big.Float).Quo(
			new(big.Float).SetInt(totalSupply),
			big.NewFloat(1000000),
		)
		fmt.Printf("\n💰 USDT 总供应量: %.2f USDT\n", supply)
		fmt.Printf("   (约 %s 美元)\n", formatBigNumber(totalSupply))
	}

	// 查询 Vitalik 的 USDT 余额
	vitalikBalance, err := balanceOf(client, contractAddress, VitalikAddress)
	if err != nil {
		log.Printf("⚠️  查询余额失败: %v\n", err)
	} else {
		balance := new(big.Float).Quo(
			new(big.Float).SetInt(vitalikBalance),
			big.NewFloat(1000000),
		)
		fmt.Printf("\n👤 Vitalik 的 USDT 余额: %.2f USDT\n", balance)
	}

	// 查询 Binance 热钱包的 USDT 余额
	binanceBalance, err := balanceOf(client, contractAddress, BinanceAddress)
	if err != nil {
		log.Printf("⚠️  查询余额失败: %v\n", err)
	} else {
		balance := new(big.Float).Quo(
			new(big.Float).SetInt(binanceBalance),
			big.NewFloat(1000000),
		)
		fmt.Printf("\n🏦 Binance 热钱包 USDT 余额: %.2f USDT\n", balance)
		fmt.Printf("   (约 %s 美元)\n", formatBigNumber(binanceBalance))
	}

	// ==================== 查询代币信息 ====================
	fmt.Println("\n\n📋 代币信息")
	fmt.Println(strings.Repeat("=", 80))

	// Name
	name, err := getName(client, contractAddress)
	if err == nil {
		fmt.Printf("📝 名称: %s\n", name)
	}

	// Symbol
	symbol, err := getSymbol(client, contractAddress)
	if err == nil {
		fmt.Printf("🔤 符号: %s\n", symbol)
	}

	// Decimals
	decimals, err := getDecimals(client, contractAddress)
	if err == nil {
		fmt.Printf("🔢 精度: %d\n", decimals)
	}

	// ==================== 批量查询多个地址 ====================
	fmt.Println("\n\n📋 批量查询地址余额")
	fmt.Println(strings.Repeat("=", 80))

	addresses := map[string]common.Address{
		"Vitalik": VitalikAddress,
		"Binance": BinanceAddress,
		"USDT合约":  contractAddress,
	}

	for name, addr := range addresses {
		// 查询 ETH 余额
		ethBalance, err := client.BalanceAt(context.Background(), addr, nil)
		if err != nil {
			continue
		}

		// 查询 USDT 余额
		usdtBalance, err := balanceOf(client, contractAddress, addr)
		if err != nil {
			continue
		}

		ethFloat := new(big.Float).Quo(
			new(big.Float).SetInt(ethBalance),
			big.NewFloat(1e18),
		)

		usdtFloat := new(big.Float).Quo(
			new(big.Float).SetInt(usdtBalance),
			big.NewFloat(1e6),
		)

		fmt.Printf("\n💼 %s\n", name)
		fmt.Printf("   地址: %s\n", addr.Hex())
		fmt.Printf("   ETH:  %.4f\n", ethFloat)
		fmt.Printf("   USDT: %.2f\n", usdtFloat)
	}

	// ==================== 总结 ====================
	fmt.Println("\n\n" + strings.Repeat("=", 80))
	fmt.Println("🎉 完成！你已经学会:")
	fmt.Println("  ✅ 调用智能合约的 view 方法（不消耗 Gas）")
	fmt.Println("  ✅ 查询 ERC20 代币余额和信息")
	fmt.Println("  ✅ 批量查询多个地址")
	fmt.Println("  ✅ 处理不同精度的代币")

	fmt.Println("\n📚 下一步学习:")
	fmt.Println("  1. 在测试网发送交易（需要签名）")
	fmt.Println("  2. 调用合约的 write 方法（transfer, approve）")
	fmt.Println("  3. 部署新的智能合约")
}

// ==================== 辅助函数 ====================

// getTotalSupply 查询代币总供应量
// 调用合约的 totalSupply() 方法
func getTotalSupply(client *ethclient.Client, contractAddr common.Address) (*big.Int, error) {
	// totalSupply() 的方法签名
	// keccak256("totalSupply()") 的前 4 字节 = 0x18160ddd
	data := common.FromHex("0x18160ddd")

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(result), nil
}

// balanceOf 查询地址的代币余额
// 调用合约的 balanceOf(address) 方法
func balanceOf(client *ethclient.Client, contractAddr, account common.Address) (*big.Int, error) {
	// balanceOf(address) 的方法签名
	// keccak256("balanceOf(address)") 的前 4 字节 = 0x70a08231
	methodID := common.FromHex("0x70a08231")

	// 参数：地址（补齐到 32 字节）
	paddedAddress := common.LeftPadBytes(account.Bytes(), 32)

	// 组合 data
	data := append(methodID, paddedAddress...)

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(result), nil
}

// getName 查询代币名称
func getName(client *ethclient.Client, contractAddr common.Address) (string, error) {
	// name() 的方法签名 = 0x06fdde03
	data := common.FromHex("0x06fdde03")

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	// 解析字符串返回值（简化版本）
	if len(result) > 64 {
		offset := new(big.Int).SetBytes(result[0:32]).Uint64()
		length := new(big.Int).SetBytes(result[offset : offset+32]).Uint64()
		stringData := result[offset+32 : offset+32+length]
		return string(stringData), nil
	}

	return "Unknown", nil
}

// getSymbol 查询代币符号
func getSymbol(client *ethclient.Client, contractAddr common.Address) (string, error) {
	// symbol() 的方法签名 = 0x95d89b41
	data := common.FromHex("0x95d89b41")

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	// 解析字符串返回值
	if len(result) > 64 {
		offset := new(big.Int).SetBytes(result[0:32]).Uint64()
		length := new(big.Int).SetBytes(result[offset : offset+32]).Uint64()
		stringData := result[offset+32 : offset+32+length]
		return string(stringData), nil
	}

	return "Unknown", nil
}

// getDecimals 查询代币精度
func getDecimals(client *ethclient.Client, contractAddr common.Address) (uint8, error) {
	// decimals() 的方法签名 = 0x313ce567
	data := common.FromHex("0x313ce567")

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, err
	}

	if len(result) > 0 {
		decimals := new(big.Int).SetBytes(result)
		return uint8(decimals.Uint64()), nil
	}

	return 0, fmt.Errorf("无法解析 decimals")
}

// formatBigNumber 格式化大数字（添加千分位逗号）
func formatBigNumber(n *big.Int) string {
	// 将 big.Int 除以 1,000,000 得到实际美元数
	dollars := new(big.Int).Div(n, big.NewInt(1000000))
	str := dollars.String()

	// 添加千分位逗号
	var result []rune
	for i, r := range reverseString(str) {
		if i > 0 && i%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, r)
	}

	return reverseString(string(result))
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
