package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func main () {
	rpcURL := "http://127.0.0.1:8545"

	fmt.Println("xxx")
		// 1. 检查链ID
	fmt.Println("\n1. 查询链ID:")
	chainID, err := getChainID(rpcURL)
	if err != nil {
		log.Printf("获取链ID失败: %v", err)
		return
	}
	fmt.Printf("Chain ID: %s (十进制: %d)\n", chainID, hexToDecimal(chainID))
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