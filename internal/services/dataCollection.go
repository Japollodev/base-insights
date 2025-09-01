package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"crypto-insights/internal/config"
	"crypto-insights/internal/database"
)

type DataCollectionService struct {
	config             config.DataCollectionConfig
	blockchainService  *BlockchainService
	db                 *gorm.DB
	redis              *database.RedisClient
	logger             *logrus.Logger
	
	mu                 sync.RWMutex
	isRunning          bool
	lastProcessedBlock uint64
	stats              *ServiceStats
	
	// Cron jobs
	cron *cron.Cron
	
	// Channels for graceful shutdown
	stopChan chan struct{}
}

type ServiceStats struct {
	BlocksProcessed     uint64    `json:"blocksProcessed"`
	TransactionsProcessed uint64  `json:"transactionsProcessed"`
	Errors              uint64    `json:"errors"`
	LastUpdate          time.Time `json:"lastUpdate"`
}

func NewDataCollectionService(
	cfg config.DataCollectionConfig,
	blockchainService *BlockchainService,
	db *gorm.DB,
	redis *database.RedisClient,
	logger *logrus.Logger,
) *DataCollectionService {
	return &DataCollectionService{
		config:            cfg,
		blockchainService: blockchainService,
		db:                db,
		redis:             redis,
		logger:            logger,
		stats:             &ServiceStats{},
		cron:              cron.New(cron.WithSeconds()),
		stopChan:          make(chan struct{}),
	}
}

func (d *DataCollectionService) Start(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isRunning {
		d.logger.Warn("Data collection service is already running")
		return nil
	}

	// Initialize blockchain service
	if err := d.blockchainService.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize blockchain service: %w", err)
	}

	// Get last processed block from cache or database
	if err := d.loadLastProcessedBlock(ctx); err != nil {
		d.logger.Warnf("Failed to load last processed block: %v", err)
		d.lastProcessedBlock = 0
	}

	d.logger.Infof("Starting data collection service from block %d", d.lastProcessedBlock)

	// Start cron jobs
	d.startCronJobs()

	d.isRunning = true
	d.logger.Info("Data collection service started successfully")

	return nil
}

func (d *DataCollectionService) Stop(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isRunning {
		d.logger.Warn("Data collection service is not running")
		return nil
	}

	// Stop cron jobs
	d.cron.Stop()

	// Signal stop
	close(d.stopChan)

	d.isRunning = false
	d.logger.Info("Data collection service stopped")

	return nil
}

func (d *DataCollectionService) startCronJobs() {
	// Block scanning every 5 seconds
	d.cron.AddFunc("*/5 * * * * *", func() {
		if err := d.scanForNewBlocks(context.Background()); err != nil {
			d.logger.Errorf("Error in block scanning: %v", err)
			d.stats.Errors++
		}
	})

	// Price updates every minute
	d.cron.AddFunc("0 * * * * *", func() {
		if err := d.updatePrices(context.Background()); err != nil {
			d.logger.Errorf("Error updating prices: %v", err)
			d.stats.Errors++
		}
	})

	// Stats collection every minute
	d.cron.AddFunc("0 * * * * *", func() {
		if err := d.updateStats(context.Background()); err != nil {
			d.logger.Errorf("Error updating stats: %v", err)
		}
	})

	d.cron.Start()
}

func (d *DataCollectionService) loadLastProcessedBlock(ctx context.Context) error {
	// Try cache first
	cached, err := d.redis.Get("last_processed_block")
	if err == nil {
		if blockNum, err := strconv.ParseUint(cached, 10, 64); err == nil {
			d.lastProcessedBlock = blockNum
			return nil
		}
	}

	// Try database
	var lastBlock database.Block
	if err := d.db.Order("number DESC").First(&lastBlock).Error; err == nil {
		d.lastProcessedBlock = lastBlock.Number
		return nil
	}

	return fmt.Errorf("no last processed block found")
}

func (d *DataCollectionService) scanForNewBlocks(ctx context.Context) error {
	currentBlock, err := d.blockchainService.GetCurrentBlock(ctx)
	if err != nil {
		return err
	}

	if currentBlock.Number > d.lastProcessedBlock {
		d.logger.Infof("New blocks detected: %d to %d", d.lastProcessedBlock+1, currentBlock.Number)

		// Process blocks in batches
		batchSize := d.config.BatchSize
		for i := d.lastProcessedBlock + 1; i <= currentBlock.Number; i += uint64(batchSize) {
			endBlock := i + uint64(batchSize) - 1
			if endBlock > currentBlock.Number {
				endBlock = currentBlock.Number
			}

			if err := d.processBlockRange(ctx, i, endBlock); err != nil {
				return err
			}
		}

		d.lastProcessedBlock = currentBlock.Number
		if err := d.redis.SetEX("last_processed_block", d.lastProcessedBlock, time.Hour); err != nil {
			d.logger.Warnf("Failed to cache last processed block: %v", err)
		}

		d.stats.BlocksProcessed += (currentBlock.Number - d.lastProcessedBlock)
		d.stats.LastUpdate = time.Now()
	}

	return nil
}

func (d *DataCollectionService) processBlockRange(ctx context.Context, startBlock, endBlock uint64) error {
	for blockNumber := startBlock; blockNumber <= endBlock; blockNumber++ {
		if err := d.processBlock(ctx, blockNumber); err != nil {
			return err
		}

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (d *DataCollectionService) processBlock(ctx context.Context, blockNumber uint64) error {
	block, err := d.blockchainService.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return err
	}

	// Store block data
	if err := d.storeBlockData(block); err != nil {
		return err
	}

	// Process transactions
	if len(block.Transactions) > 0 {
		if err := d.processBlockTransactions(block); err != nil {
			return err
		}
	}

	d.logger.Debugf("Processed block %d with %d transactions", blockNumber, len(block.Transactions))
	return nil
}

func (d *DataCollectionService) storeBlockData(block *Block) error {
	// Check if block already exists
	var existingBlock database.Block
	if err := d.db.Where("number = ?", block.Number).First(&existingBlock).Error; err == nil {
		return nil // Block already exists
	}

	dbBlock := database.Block{
		Number:            block.Number,
		Hash:              block.Hash,
		ParentHash:        block.ParentHash,
		Timestamp:         block.Timestamp,
		GasUsed:           block.GasUsed,
		GasLimit:          block.GasLimit,
		Miner:             block.Miner,
		Difficulty:        block.Difficulty,
		TotalDifficulty:   block.TotalDifficulty,
		Size:              block.Size,
		ExtraData:         block.ExtraData,
		BaseFeePerGas:     block.BaseFeePerGas,
	}

	return d.db.Create(&dbBlock).Error
}

func (d *DataCollectionService) processBlockTransactions(block *Block) error {
	for _, tx := range block.Transactions {
		if err := d.processTransaction(&tx, block); err != nil {
			return err
		}
	}
	return nil
}

func (d *DataCollectionService) processTransaction(tx *Transaction, block *Block) error {
	// Check if transaction already exists
	var existingTx database.Transaction
	if err := d.db.Where("hash = ?", tx.Hash).First(&existingTx).Error; err == nil {
		return nil // Transaction already exists
	}

	dbTx := database.Transaction{
		Hash:             tx.Hash,
		BlockNumber:      tx.BlockNumber,
		BlockHash:        tx.BlockHash,
		From:             tx.From,
		To:               tx.To,
		Value:            tx.Value,
		Gas:              tx.Gas,
		GasPrice:         tx.GasPrice,
		Nonce:            tx.Nonce,
		Input:            tx.Input,
		TransactionIndex: tx.TransactionIndex,
	}

	if err := d.db.Create(&dbTx).Error; err != nil {
		return err
	}

	d.stats.TransactionsProcessed++
	return nil
}

func (d *DataCollectionService) updatePrices(ctx context.Context) error {
	gasData, err := d.blockchainService.GetGasPrice(ctx)
	if err != nil {
		return err
	}

	dbGasPrice := database.GasPrice{
		GasPrice:           gasData.GasPrice,
		GasPriceGwei:       gasData.GasPriceGwei,
		MaxFeePerGas:       gasData.MaxFeePerGas,
		MaxPriorityFeePerGas: gasData.MaxPriorityFeePerGas,
		BaseFeePerGas:      gasData.BaseFeePerGas,
		Timestamp:          time.Now(),
	}

	return d.db.Create(&dbGasPrice).Error
}

func (d *DataCollectionService) updateStats(ctx context.Context) error {
	stats := map[string]interface{}{
		"blocksProcessed":     d.stats.BlocksProcessed,
		"transactionsProcessed": d.stats.TransactionsProcessed,
		"errors":              d.stats.Errors,
		"lastUpdate":          d.stats.LastUpdate,
		"isRunning":           d.isRunning,
		"lastProcessedBlock":  d.lastProcessedBlock,
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return d.redis.SetEX("service_stats", string(statsJSON), 5*time.Minute)
}

func (d *DataCollectionService) GetStats() (*ServiceStats, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Try cache first
	cached, err := d.redis.Get("service_stats")
	if err == nil {
		var stats ServiceStats
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			return &stats, nil
		}
	}

	return d.stats, nil
}

func (d *DataCollectionService) GetServiceStatus() map[string]interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()

	stats, _ := d.GetStats()

	return map[string]interface{}{
		"isRunning":          d.isRunning,
		"lastProcessedBlock": d.lastProcessedBlock,
		"stats":              stats,
		"uptime":             time.Since(d.stats.LastUpdate).String(),
	}
}
