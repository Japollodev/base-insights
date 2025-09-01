package services

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"

	"crypto-insights/internal/config"
)

type BlockchainService struct {
	client     *ethclient.Client
	config     config.BlockchainConfig
	logger     *logrus.Logger
	isReady    bool
	lastBlock  uint64
}

type Block struct {
	Number            uint64    `json:"number"`
	Hash              string    `json:"hash"`
	ParentHash        string    `json:"parentHash"`
	Timestamp         time.Time `json:"timestamp"`
	Transactions      []Transaction `json:"transactions"`
	GasUsed           uint64    `json:"gasUsed"`
	GasLimit          uint64    `json:"gasLimit"`
	Miner             string    `json:"miner"`
	Difficulty        string    `json:"difficulty"`
	TotalDifficulty   string    `json:"totalDifficulty"`
	Size              uint64    `json:"size"`
	ExtraData         string    `json:"extraData"`
	BaseFeePerGas     string    `json:"baseFeePerGas"`
}

type Transaction struct {
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              uint64 `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Nonce            uint64 `json:"nonce"`
	Input            string `json:"input"`
	BlockNumber      uint64 `json:"blockNumber"`
	BlockHash        string `json:"blockHash"`
	TransactionIndex uint64 `json:"transactionIndex"`
}

type GasData struct {
	GasPrice           string `json:"gasPrice"`
	GasPriceGwei       string `json:"gasPriceGwei"`
	MaxFeePerGas       string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	BaseFeePerGas      string `json:"baseFeePerGas"`
}

func NewBlockchainService(cfg config.BlockchainConfig, logger *logrus.Logger) *BlockchainService {
	return &BlockchainService{
		config: cfg,
		logger: logger,
	}
}

func (b *BlockchainService) Initialize(ctx context.Context) error {
	client, err := ethclient.DialContext(ctx, b.config.RPCURL)
	if err != nil {
		return err
	}

	b.client = client

	// Test connection
	blockNumber, err := b.client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	b.lastBlock = blockNumber
	b.isReady = true

	b.logger.Infof("Blockchain service initialized. Connected to Base L2 at block %d", blockNumber)
	return nil
}

func (b *BlockchainService) GetCurrentBlock(ctx context.Context) (*Block, error) {
	if !b.isReady {
		return nil, ErrServiceNotReady
	}

	blockNumber, err := b.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	block, err := b.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}

	return b.convertBlock(block), nil
}

func (b *BlockchainService) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error) {
	if !b.isReady {
		return nil, ErrServiceNotReady
	}

	block, err := b.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}

	return b.convertBlock(block), nil
}

func (b *BlockchainService) GetGasPrice(ctx context.Context) (*GasData, error) {
	if !b.isReady {
		return nil, ErrServiceNotReady
	}

	gasPrice, err := b.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to Gwei
	gwei := new(big.Int).Div(gasPrice, big.NewInt(1e9))

	return &GasData{
		GasPrice:     gasPrice.String(),
		GasPriceGwei: gwei.String(),
		// Note: EIP-1559 fields would need additional RPC calls
		MaxFeePerGas:       "",
		MaxPriorityFeePerGas: "",
		BaseFeePerGas:      "",
	}, nil
}

func (b *BlockchainService) GetNetworkInfo(ctx context.Context) (map[string]interface{}, error) {
	if !b.isReady {
		return nil, ErrServiceNotReady
	}

	chainID, err := b.client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	blockNumber, err := b.client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"chainId":      chainID.String(),
		"isListening":  true, // Assume true if we can connect
		"currentBlock": blockNumber,
		"networkName":  "Base L2",
		"rpcUrl":       b.config.RPCURL,
		"explorerUrl":  b.config.ExplorerURL,
	}, nil
}

func (b *BlockchainService) IsReady() bool {
	return b.isReady
}

func (b *BlockchainService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"isInitialized": b.isReady,
		"isConnected":   b.isReady,
		"lastBlockNumber": b.lastBlock,
		"rpcUrl":        b.config.RPCURL,
		"chainId":       b.config.ChainID,
	}
}

func (b *BlockchainService) convertBlock(block *types.Block) *Block {
	transactions := make([]Transaction, 0)
	for i, tx := range block.Transactions() {
		transactions = append(transactions, Transaction{
			Hash:             tx.Hash().Hex(),
			From:             "", // Would need to get from receipt
			To:               tx.To().Hex(),
			Value:            tx.Value().String(),
			Gas:              tx.Gas(),
			GasPrice:         tx.GasPrice().String(),
			Nonce:            tx.Nonce(),
			Input:            common.Bytes2Hex(tx.Data()),
			BlockNumber:      block.NumberU64(),
			BlockHash:        block.Hash().Hex(),
			TransactionIndex: uint64(i),
		})
	}

	return &Block{
		Number:          block.NumberU64(),
		Hash:            block.Hash().Hex(),
		ParentHash:      block.ParentHash().Hex(),
		Timestamp:       time.Unix(int64(block.Time()), 0),
		Transactions:    transactions,
		GasUsed:         block.GasUsed(),
		GasLimit:        block.GasLimit(),
		Miner:           block.Coinbase().Hex(),
		Difficulty:      block.Difficulty().String(),
		TotalDifficulty: block.Difficulty().String(), // Would need to calculate
		Size:            uint64(block.Size()),
		ExtraData:       common.Bytes2Hex(block.Extra()),
		BaseFeePerGas:   block.BaseFee().String(),
	}
}

// Errors
var ErrServiceNotReady = errors.New("blockchain service not ready")
