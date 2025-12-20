package db

import (
	"fmt"
	"log"
	"math/big"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库连接封装
type Database struct {
	DB *gorm.DB
}

// NewDatabase 创建新的数据库连接
func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&Balance{}, &Transaction{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("数据库连接成功，表迁移完成")

	return &Database{DB: db}, nil
}

// UpdateBalance 更新或创建地址余额
func (d *Database) UpdateBalance(address, balance string) error {
	result := d.DB.Model(&Balance{}).
		Where("address = ?", address).
		Updates(Balance{
			Address: address,
			Balance: balance,
		})

	if result.RowsAffected == 0 {
		// 如果没有更新任何行，则创建新记录
		return d.DB.Create(&Balance{
			Address: address,
			Balance: balance,
		}).Error
	}

	return result.Error
}

// GetBalance 查询地址余额
func (d *Database) GetBalance(address string) (*Balance, error) {
	var balance Balance
	err := d.DB.Where("address = ?", address).First(&balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// SaveTransaction 保存交易记录
func (d *Database) SaveTransaction(tx *Transaction) error {
	return d.DB.Create(tx).Error
}

// GetTransactionsByAddress 查询地址的交易记录
func (d *Database) GetTransactionsByAddress(address string, limit, offset int) ([]Transaction, error) {
	var transactions []Transaction
	err := d.DB.Where("from_address = ? OR to_address = ?", address, address).
		Order("block_number DESC, timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

// GetTransactionByHash 根据交易哈希查询交易
func (d *Database) GetTransactionByHash(txHash string) (*Transaction, error) {
	var transaction Transaction
	err := d.DB.Where("tx_hash = ?", txHash).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ProcessTransfer 处理转账事件（更新余额）
func (d *Database) ProcessTransfer(from, to, amount string) error {
	// 开启事务
	return d.DB.Transaction(func(tx *gorm.DB) error {
		// 如果 from 不是零地址，则减少 from 的余额
		if from != "0x0000000000000000000000000000000000000000" {
			var fromBalance Balance
			if err := tx.Where("address = ?", from).First(&fromBalance).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// 如果不存在，可能是第一次查询，跳过（或者设置为0）
					log.Printf("警告: 地址 %s 余额不存在，跳过减少操作", from)
				} else {
					return err
				}
			} else {
				// 计算新余额
				oldBalance := new(big.Int)
				oldBalance.SetString(fromBalance.Balance, 10)
				
				transferAmount := new(big.Int)
				transferAmount.SetString(amount, 10)
				
				newBalance := new(big.Int).Sub(oldBalance, transferAmount)
				
				// 更新余额
				if err := tx.Model(&Balance{}).
					Where("address = ?", from).
					Update("balance", newBalance.String()).Error; err != nil {
					return err
				}
			}
		}
		
		// 如果 to 不是零地址，则增加 to 的余额
		if to != "0x0000000000000000000000000000000000000000" {
			var toBalance Balance
			if err := tx.Where("address = ?", to).First(&toBalance).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// 如果不存在，则创建新记录
					if err := tx.Create(&Balance{
						Address: to,
						Balance: amount,
					}).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				// 计算新余额
				oldBalance := new(big.Int)
				oldBalance.SetString(toBalance.Balance, 10)
				
				transferAmount := new(big.Int)
				transferAmount.SetString(amount, 10)
				
				newBalance := new(big.Int).Add(oldBalance, transferAmount)
				
				// 更新余额
				if err := tx.Model(&Balance{}).
					Where("address = ?", to).
					Update("balance", newBalance.String()).Error; err != nil {
					return err
				}
			}
		}
		
		return nil
	})
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
