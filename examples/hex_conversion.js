// JavaScript 16进制转换示例

console.log("=== JavaScript 16进制转换 ===");

// 方法1: parseInt() - 最常用
const hexString = "0x7a69";
const decimal1 = parseInt(hexString, 16);
console.log(`方法1 - parseInt(): ${hexString} = ${decimal1}`);

// 方法2: 直接使用 Number() - 自动识别0x前缀
const decimal2 = Number(hexString);
console.log(`方法2 - Number(): ${hexString} = ${decimal2}`);

// 方法3: 使用 BigInt (适用于大数)
const decimal3 = BigInt(hexString);
console.log(`方法3 - BigInt(): ${hexString} = ${decimal3}`);

// 方法4: 手动去除0x前缀
const hexWithoutPrefix = hexString.slice(2); // 去除 "0x"
const decimal4 = parseInt(hexWithoutPrefix, 16);
console.log(`方法4 - 去除前缀: ${hexWithoutPrefix} = ${decimal4}`);

// 反向转换: 10进制转16进制
const backToHex = decimal1.toString(16);
console.log(`反向转换: ${decimal1} = 0x${backToHex}`);

// 处理以太坊余额转换 (Wei to ETH)
console.log("\n=== 以太坊余额转换示例 ===");
const balanceHex = "0x21df03e39e1dd38c2b8"; // 从API返回的余额
const balanceWei = BigInt(balanceHex);
const balanceEth = Number(balanceWei) / Math.pow(10, 18);
console.log(`余额: ${balanceHex} Wei = ${balanceEth.toFixed(6)} ETH`);

// 实用函数
function hexToDecimal(hex) {
    return parseInt(hex, 16);
}

function weiToEth(weiHex) {
    const wei = BigInt(weiHex);
    const eth = Number(wei) / Math.pow(10, 18);
    return eth;
}

// 测试
console.log("\n=== 实用函数测试 ===");
console.log(`Chain ID: ${hexToDecimal("0x7a69")}`);
console.log(`余额: ${weiToEth("0x21df03e39e1dd38c2b8").toFixed(6)} ETH`);