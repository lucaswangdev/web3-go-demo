#!/usr/bin/env python3
"""
简单的Python脚本来测试Anvil RPC调用
使用方法: python3 test_rpc.py
"""

import requests
import json

def send_rpc_request(method, params=None):
    """发送JSON-RPC请求"""
    url = "http://127.0.0.1:8545"
    
    payload = {
        "jsonrpc": "2.0",
        "method": method,
        "params": params or [],
        "id": 1
    }
    
    try:
        response = requests.post(url, json=payload)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"❌ 请求失败: {e}")
        return None

def wei_to_eth(wei_hex):
    """将Wei转换为ETH"""
    wei = int(wei_hex, 16)
    eth = wei / (10 ** 18)
    return eth

def hex_to_dec(hex_str):
    """16进制转10进制"""
    return int(hex_str, 16)

def main():
    print("=== Anvil 本地节点 RPC 测试 ===\n")
    
    # 测试账户
    test_accounts = [
        "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
        "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
    ]
    
    # 1. 检查链ID
    print("1. 查询链ID:")
    result = send_rpc_request("eth_chainId")
    if result and 'result' in result:
        chain_id_hex = result['result']
        chain_id_dec = hex_to_dec(chain_id_hex)
        print(f"   Chain ID: {chain_id_hex} (十进制: {chain_id_dec})")
    else:
        print("   ❌ 获取链ID失败，请确保Anvil正在运行")
        return
    
    # 2. 查询账户余额
    print("\n2. 查询账户余额:")
    for i, account in enumerate(test_accounts):
        result = send_rpc_request("eth_getBalance", [account, "latest"])
        if result and 'result' in result:
            balance_hex = result['result']
            balance_eth = wei_to_eth(balance_hex)
            print(f"   账户 {i}: {account}")
            print(f"   余额: {balance_eth:.6f} ETH ({balance_hex} Wei)")
        else:
            print(f"   ❌ 获取账户 {i} 余额失败")
    
    # 3. 查询最新区块号
    print("\n3. 查询最新区块号:")
    result = send_rpc_request("eth_blockNumber")
    if result and 'result' in result:
        block_hex = result['result']
        block_dec = hex_to_dec(block_hex)
        print(f"   最新区块: {block_hex} (十进制: {block_dec})")
    
    # 4. 查询Gas价格
    print("\n4. 查询Gas价格:")
    result = send_rpc_request("eth_gasPrice")
    if result and 'result' in result:
        gas_price_hex = result['result']
        gas_price_dec = hex_to_dec(gas_price_hex)
        gas_price_gwei = gas_price_dec / (10 ** 9)  # 转换为Gwei
        print(f"   Gas价格: {gas_price_hex} ({gas_price_dec} Wei, {gas_price_gwei:.2f} Gwei)")
    
    # 5. 查询网络版本
    print("\n5. 查询网络版本:")
    result = send_rpc_request("net_version")
    if result and 'result' in result:
        net_version = result['result']
        print(f"   网络版本: {net_version}")
    
    print("\n✅ 测试完成！")
    print("\n💡 提示:")
    print("   - 如果看到错误，请确保运行: anvil")
    print("   - 默认端口是 8545")
    print("   - 可以使用 Ctrl+C 停止 Anvil")

if __name__ == "__main__":
    main()