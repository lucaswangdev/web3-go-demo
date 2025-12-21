package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// JSON-RPC 请求结构
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSON-RPC 响应结构
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	fmt.Println("=== 直接调用 Anvil JSON-RPC API ===")
	
	rpcURL := "http://127.0.0.1:8545"
	
	// 测试账户地址
	testAccounts := []string{
		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
	}
	
	// 1. 检查链ID
	fmt.Println("\n1. 查询链ID:")
	chainID, err := getChainID(rpcURL)
	if err != nil {
		log.Printf("获取链ID失败: %v", err)
		return
	}
	fmt.Printf("Chain ID: %s (十进制: %d)\n", chainID, hexToDecimal(chainID))
	
	// 2. 查询账户余额
	fmt.Println("\n2. 查询账户余额:")
	for i, address := range testAccounts {
		balance, err := getBalance(rpcURL, address)
		if err != nil {
			log.Printf("获取账户 %d 余额失败: %v", i, err)
			continue
		}
		
		// 转换为ETH
		ethBalance := weiToEth(balance)
		fmt.Printf("账户 %d (%s): %s ETH\n", i, address, ethBalance)
	}
	
	// 3. 获取最新区块号
	fmt.Println("\n3. 查询最新区块号:")
	blockNumber, err := getBlockNumber(rpcURL)
	if err != nil {
		log.Printf("获取区块号失败: %v", err)
		return
	}
	fmt.Printf("最新区块号: %s (十进制: %d)\n", blockNumber, hexToDecimal(blockNumber))
	
	// 4. 获取网络版本
	fmt.Println("\n4. 查询网络版本:")
	netVersion, err := getNetVersion(rpcURL)
	if err != nil {
		log.Printf("获取网络版本失败: %v", err)
		return
	}
	fmt.Printf("网络版本: %s\n", netVersion)
	
	// 5. 获取Gas价格
	fmt.Println("\n5. 查询Gas价格:")
	gasPrice, err := getGasPrice(rpcURL)
	if err != nil {
		log.Printf("获取Gas价格失败: %v", err)
		return
	}
	fmt.Printf("Gas价格: %s Wei (十进制: %d)\n", gasPrice, hexToDecimal(gasPrice))
}

// 发送JSON-RPC请求
func sendJSONRPCRequest(url string, method string, params []interface{}) (*JSONRPCResponse, error) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	
	var response JSONRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	if response.Error != nil {
		return nil, fmt.Errorf("RPC错误: %s", response.Error.Message)
	}
	
	return &response, nil
}

// 获取链ID
func getChainID(rpcURL string) (string, error) {
	resp, err := sendJSONRPCRequest(rpcURL, "eth_chainId", []interface{}{})
	if err != nil {
		return "", err
	}
	
	return resp.Result.(string), nil
}

// 获取账户余额
func getBalance(rpcURL, address string) (string, error) {
	params := []interface{}{address, "latest"}
	resp, err := sendJSONRPCRequest(rpcURL, "eth_getBalance", params)
	if err != nil {
		return "", err
	}
	
	return resp.Result.(string), nil
}

// 获取最新区块号
func getBlockNumber(rpcURL string) (string, error) {
	resp, err := sendJSONRPCRequest(rpcURL, "eth_blockNumber", []interface{}{})
	if err != nil {
		return "", err
	}
	
	return resp.Result.(string), nil
}

// 获取网络版本
func getNetVersion(rpcURL string) (string, error) {
	resp, err := sendJSONRPCRequest(rpcURL, "net_version", []interface{}{})
	if err != nil {
		return "", err
	}
	
	return resp.Result.(string), nil
}

// 获取Gas价格
func getGasPrice(rpcURL string) (string, error) {
	resp, err := sendJSONRPCRequest(rpcURL, "eth_gasPrice", []interface{}{})
	if err != nil {
		return "", err
	}
	
	return resp.Result.(string), nil
}

// 16进制转十进制
func hexToDecimal(hexStr string) int64 {
	// 移除0x前缀
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	
	decimal, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0
	}
	
	return decimal
}

// Wei转ETH
func weiToEth(weiHex string) string {
	// 移除0x前缀
	if len(weiHex) > 2 && weiHex[:2] == "0x" {
		weiHex = weiHex[2:]
	}
	
	// 转换为big.Int
	wei := new(big.Int)
	wei.SetString(weiHex, 16)
	
	// 1 ETH = 10^18 Wei
	ethDivisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	
	// 转换为ETH
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(ethDivisor))
	
	return eth.Text('f', 6) // 保留6位小数
}