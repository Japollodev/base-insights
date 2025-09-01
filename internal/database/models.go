package database

import (
	"time"

	"gorm.io/gorm"
)

type Block struct {
	ID              uint      `gorm:"primaryKey"`
	Number          uint64    `gorm:"uniqueIndex;not null"`
	Hash            string    `gorm:"uniqueIndex;not null"`
	ParentHash      string    `gorm:"not null"`
	Timestamp       time.Time `gorm:"not null"`
	GasUsed         uint64    `gorm:"not null"`
	GasLimit        uint64    `gorm:"not null"`
	Miner           string    `gorm:"not null"`
	Difficulty      string    `gorm:"not null"`
	TotalDifficulty string    `gorm:"not null"`
	Size            uint64    `gorm:"not null"`
	ExtraData       string    `gorm:"not null"`
	BaseFeePerGas   string    `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	Hash            string    `gorm:"uniqueIndex;not null"`
	BlockNumber     uint64    `gorm:"not null"`
	BlockHash       string    `gorm:"not null"`
	From            string    `gorm:"not null"`
	To              string    `gorm:"not null"`
	Value           string    `gorm:"not null"`
	Gas             uint64    `gorm:"not null"`
	GasPrice        string    `gorm:"not null"`
	Nonce           uint64    `gorm:"not null"`
	Input           string    `gorm:"not null"`
	TransactionIndex uint64   `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type GasPrice struct {
	ID                 uint      `gorm:"primaryKey"`
	GasPrice           string    `gorm:"not null"`
	GasPriceGwei       string    `gorm:"not null"`
	MaxFeePerGas       string    `gorm:"not null"`
	MaxPriorityFeePerGas string  `gorm:"not null"`
	BaseFeePerGas      string    `gorm:"not null"`
	Timestamp          time.Time `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// BeforeCreate hook to set timestamps
func (b *Block) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	return nil
}

func (g *GasPrice) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now
	return nil
}
