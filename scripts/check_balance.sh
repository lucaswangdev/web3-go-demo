#!/bin/bash

# 本地区块链余额查询脚本
# 使用方法: ./scripts/check_balance.sh [账户地址]

RPC_URL="http://127.0.0.1:8545"
DEFAULT_ACCOUNT="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

# 如果提供了参数，使用参数作为账户地址，否则使用默认地址
ACCOUNT=${1:-$DEFAULT_ACCOUNT}

echo "=== Anvil 本地节点状态查询 ==="
echo "RPC URL: $RPC_URL"
echo "查询账户: $ACCOUNT"
echo ""

# 1. 检查节点是否运行
echo "1. 检查节点连接..."
CHAIN_ID=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "eth_chainId", "params": [], "id": 1}' \
  $RPC_URL | jq -r '.result')

if [ "$CHAIN_ID" = "null" ] || [ -z "$CHAIN_ID" ]; then
    echo "❌ 无法连接到本地节点，请确保 Anvil 正在运行:"
    echo "   anvil"
    exit 1
fi

# 转换链ID为十进制
CHAIN_ID_DEC=$((16#${CHAIN_ID#0x}))
echo "✅ 节点连接成功"
echo "   Chain ID: $CHAIN_ID (十进制: $CHAIN_ID_DEC)"

# 2. 查询账户余额
echo ""
echo "2. 查询账户余额..."
BALANCE_HEX=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\": \"2.0\", \"method\": \"eth_getBalance\", \"params\": [\"$ACCOUNT\", \"latest\"], \"id\": 1}" \
  $RPC_URL | jq -r '.result')

if [ "$BALANCE_HEX" = "null" ] || [ -z "$BALANCE_HEX" ]; then
    echo "❌ 获取余额失败"
    exit 1
fi

# 使用 Python 转换 Wei 到 ETH (如果有Python的话)
if command -v python3 &> /dev/null; then
    BALANCE_ETH=$(python3 -c "print(int('$BALANCE_HEX', 16) / 10**18)")
    echo "✅ 账户余额: $BALANCE_ETH ETH"
else
    echo "✅ 账户余额: $BALANCE_HEX Wei (16进制)"
fi

# 3. 查询最新区块
echo ""
echo "3. 查询最新区块..."
BLOCK_NUMBER_HEX=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}' \
  $RPC_URL | jq -r '.result')

BLOCK_NUMBER_DEC=$((16#${BLOCK_NUMBER_HEX#0x}))
echo "✅ 最新区块号: $BLOCK_NUMBER_DEC"

# 4. 查询Gas价格
echo ""
echo "4. 查询Gas价格..."
GAS_PRICE_HEX=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "method": "eth_gasPrice", "params": [], "id": 1}' \
  $RPC_URL | jq -r '.result')

GAS_PRICE_DEC=$((16#${GAS_PRICE_HEX#0x}))
echo "✅ Gas价格: $GAS_PRICE_DEC Wei"

echo ""
echo "=== 常用测试账户 ==="
echo "账户0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
echo "账户1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
echo ""
echo "使用方法:"
echo "  查看账户0余额: $0"
echo "  查看账户1余额: $0 0x70997970C51812dc3A010C7d01b50e0d17dc79C8"