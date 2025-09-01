package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"crypto-insights/internal/services"
)

type Server struct {
	server                *http.Server
	blockchainService     *services.BlockchainService
	dataCollectionService *services.DataCollectionService
	logger                *logrus.Logger
}

func New(
	port string,
	blockchainService *services.BlockchainService,
	dataCollectionService *services.DataCollectionService,
	logger *logrus.Logger,
) *Server {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	server := &Server{
		blockchainService:     blockchainService,
		dataCollectionService: dataCollectionService,
		logger:                logger,
	}

	// Setup routes
	server.setupRoutes(router)

	httpServer := &http.Server{
		Addr:    port,
		Handler: router,
	}

	server.server = httpServer
	return server
}

func (s *Server) setupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", s.healthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Blockchain routes
		blockchain := api.Group("/blockchain")
		{
			blockchain.GET("/status", s.getBlockchainStatus)
			blockchain.GET("/current-block", s.getCurrentBlock)
			blockchain.GET("/block/:number", s.getBlockByNumber)
			blockchain.GET("/gas-price", s.getGasPrice)
			blockchain.GET("/network-info", s.getNetworkInfo)
		}

		// Data collection routes
		data := api.Group("/data")
		{
			data.GET("/status", s.getDataCollectionStatus)
			data.GET("/stats", s.getDataCollectionStats)
			data.POST("/start", s.startDataCollection)
			data.POST("/stop", s.stopDataCollection)
		}

		// Analytics routes
		analytics := api.Group("/analytics")
		{
			analytics.GET("/blocks", s.getBlocks)
			analytics.GET("/transactions", s.getTransactions)
			analytics.GET("/gas-prices", s.getGasPrices)
		}
	}
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "crypto-insights",
	})
}

func (s *Server) getBlockchainStatus(c *gin.Context) {
	status := s.blockchainService.GetStatus()
	c.JSON(http.StatusOK, status)
}

func (s *Server) getCurrentBlock(c *gin.Context) {
	block, err := s.blockchainService.GetCurrentBlock(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, block)
}

func (s *Server) getBlockByNumber(c *gin.Context) {
	blockNumber := c.Param("number")
	// Parse block number and get block
	c.JSON(http.StatusOK, gin.H{"message": "Get block by number", "number": blockNumber})
}

func (s *Server) getGasPrice(c *gin.Context) {
	gasData, err := s.blockchainService.GetGasPrice(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gasData)
}

func (s *Server) getNetworkInfo(c *gin.Context) {
	info, err := s.blockchainService.GetNetworkInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (s *Server) getDataCollectionStatus(c *gin.Context) {
	status := s.dataCollectionService.GetServiceStatus()
	c.JSON(http.StatusOK, status)
}

func (s *Server) getDataCollectionStats(c *gin.Context) {
	stats, err := s.dataCollectionService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (s *Server) startDataCollection(c *gin.Context) {
	if err := s.dataCollectionService.Start(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data collection started"})
}

func (s *Server) stopDataCollection(c *gin.Context) {
	if err := s.dataCollectionService.Stop(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data collection stopped"})
}

func (s *Server) getBlocks(c *gin.Context) {
	// Get blocks from database
	c.JSON(http.StatusOK, gin.H{"message": "Get blocks"})
}

func (s *Server) getTransactions(c *gin.Context) {
	// Get transactions from database
	c.JSON(http.StatusOK, gin.H{"message": "Get transactions"})
}

func (s *Server) getGasPrices(c *gin.Context) {
	// Get gas prices from database
	c.JSON(http.StatusOK, gin.H{"message": "Get gas prices"})
}

func (s *Server) Start() error {
	s.logger.Infof("Starting server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}
