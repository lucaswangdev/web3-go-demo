package db

import (
	"time"
)

// Balance 地址余额表
type Balance struct {
	Address   string    `gorm:"primaryKey;type:varchar(42)" json:"address"`
	Balance   string    `gorm:"type:varchar(78);not null" json:"balance"` // 使用字符串存储大数
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Transaction 交易记录表
type Transaction struct {
	TxHash      string    `gorm:"primaryKey;type:varchar(66)" json:"tx_hash"`
	FromAddress string    `gorm:"type:varchar(42);index" json:"from_address"`
	ToAddress   string    `gorm:"type:varchar(42);index" json:"to_address"`
	Amount      string    `gorm:"type:varchar(78);not null" json:"amount"` // 使用字符串存储大数
	BlockNumber uint64    `gorm:"index;not null" json:"block_number"`
	Timestamp   time.Time `gorm:"index;not null" json:"timestamp"`
}

// TableName 指定 Balance 表名
func (Balance) TableName() string {
	return "balances"
}

// TableName 指定 Transaction 表名
func (Transaction) TableName() string {
	return "transactions"
}
