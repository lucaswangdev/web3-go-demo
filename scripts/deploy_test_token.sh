#!/bin/bash

# 部署测试 ERC20 代币到本地 Anvil 网络

echo "=== 部署测试 ERC20 代币 ==="

# 检查 Anvil 是否运行
if ! curl -s -X POST -H "Content-Type: application/json" -d '{"jsonrpc": "2.0", "method": "eth_chainId", "params": [], "id": 1}' http://127.0.0.1:8545 > /dev/null; then
    echo "❌ Anvil 未运行，请先启动:"
    echo "   anvil"
    exit 1
fi

echo "✅ Anvil 节点运行正常"

# 使用 Foundry 部署合约
echo "📦 部署 TestToken 合约..."

# 创建 foundry.toml 配置文件
cat > foundry.toml << EOF
[profile.default]
src = "contracts"
out = "out"
libs = ["lib"]
solc_version = "0.8.19"

[rpc_endpoints]
local = "http://127.0.0.1:8545"
EOF

# 编译合约
echo "🔨 编译合约..."
forge build

if [ $? -ne 0 ]; then
    echo "❌ 合约编译失败"
    exit 1
fi

# 部署合约
echo "🚀 部署合约到本地网络..."
DEPLOY_RESULT=$(forge create --rpc-url http://127.0.0.1:8545 \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --broadcast \
  contracts/TestToken.sol:TestToken)

if [ $? -ne 0 ]; then
    echo "❌ 合约部署失败"
    exit 1
fi

# 提取合约地址
CONTRACT_ADDRESS=$(echo "$DEPLOY_RESULT" | grep "Deployed to:" | awk '{print $3}')

if [ -z "$CONTRACT_ADDRESS" ]; then
    echo "❌ 无法获取合约地址"
    exit 1
fi

echo "✅ 合约部署成功!"
echo "📍 合约地址: $CONTRACT_ADDRESS"

# 更新 .env 文件
echo "📝 更新 .env 文件..."
sed -i.bak "s/CONTRACT_ADDRESS=.*/CONTRACT_ADDRESS=$CONTRACT_ADDRESS/" .env

echo ""
echo "🎉 部署完成！"
echo ""
echo "📋 合约信息:"
echo "   地址: $CONTRACT_ADDRESS"
echo "   名称: Test Token (TEST)"
echo "   总供应量: 1,000,000 TEST"
echo "   持有者: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
echo ""
echo "🔄 现在可以重新运行你的项目:"
echo "   go run cmd/tracker/main.go"
echo ""
echo "🧪 测试转账:"
echo "   cast send $CONTRACT_ADDRESS \"transfer(address,uint256)\" 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 1000000000000000000 \\"
echo "     --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \\"
echo "     --rpc-url http://127.0.0.1:8545"