# Crypto Insights - Go Service Conversion

This project has been converted from Node.js to Go, providing a high-performance blockchain analytics platform for Coinbase's Base L2 network.

## Architecture

The Go service follows functional programming principles and is structured as follows:

### Backend (Go)
- **Main Service**: `cmd/main.go` - Application entry point
- **Configuration**: `internal/config/` - Configuration management using Viper
- **Database**: `internal/database/` - PostgreSQL with GORM and Redis client
- **Services**: `internal/services/` - Blockchain and data collection services
- **HTTP Server**: `internal/server/` - RESTful API using Gin framework

### Frontend (React + Material-UI)
- **Dashboard**: Real-time blockchain statistics and charts
- **Blockchain**: Live blockchain data and gas price monitoring
- **Analytics**: Historical data analysis and visualization
- **Data Collection**: Service control and monitoring

## Features

- **Real-time Block Processing**: Continuously monitors Base L2 for new blocks
- **Transaction Analysis**: Processes and stores transaction data
- **Gas Price Monitoring**: Tracks gas price trends
- **Performance Metrics**: Service health and processing statistics
- **RESTful API**: Comprehensive API for frontend integration
- **Database Persistence**: PostgreSQL for structured data, Redis for caching

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Node.js 18+ (for frontend)

## Quick Start

### 1. Backend Setup

```bash
# Install Go dependencies
go mod tidy

# Set environment variables
export DATABASE_HOST=localhost
export DATABASE_USER=postgres
export DATABASE_PASSWORD=password
export DATABASE_NAME=crypto_insights
export REDIS_HOST=localhost
export REDIS_PORT=6379

# Run the service
go run cmd/main.go
```

### 2. Using Docker

```bash
# Start all services
docker-compose -f docker-compose.go.yml up -d

# View logs
docker-compose -f docker-compose.go.yml logs -f app
```

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm start
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Blockchain
- `GET /api/v1/blockchain/status` - Blockchain service status
- `GET /api/v1/blockchain/current-block` - Current block information
- `GET /api/v1/blockchain/block/:number` - Block by number
- `GET /api/v1/blockchain/gas-price` - Current gas prices
- `GET /api/v1/blockchain/network-info` - Network information

### Data Collection
- `GET /api/v1/data/status` - Data collection service status
- `GET /api/v1/data/stats` - Processing statistics
- `POST /api/v1/data/start` - Start data collection
- `POST /api/v1/data/stop` - Stop data collection

### Analytics
- `GET /api/v1/analytics/blocks` - Block data
- `GET /api/v1/analytics/transactions` - Transaction data
- `GET /api/v1/analytics/gas-prices` - Gas price history

## Configuration

The service uses a YAML configuration file (`config/config.yaml`) with the following sections:

- **Server**: HTTP server configuration
- **Database**: PostgreSQL connection settings
- **Redis**: Redis connection settings
- **Blockchain**: Base L2 RPC configuration
- **Data Collection**: Processing intervals and batch sizes

## Key Differences from Node.js Version

### Performance Improvements
- **Concurrent Processing**: Go's goroutines for parallel block processing
- **Memory Efficiency**: Lower memory footprint and better garbage collection
- **Faster Startup**: Reduced initialization time
- **Better Resource Management**: More efficient database connections

### Architecture Changes
- **Functional Approach**: Service methods are pure functions where possible
- **Strong Typing**: Compile-time type safety
- **Error Handling**: Explicit error handling with Go's error interface
- **Context Support**: Built-in context for request cancellation and timeouts

### Database Layer
- **GORM**: Object-relational mapping with automatic migrations
- **Connection Pooling**: Efficient database connection management
- **Transaction Support**: ACID-compliant database operations

## Monitoring and Logging

- **Structured Logging**: JSON-formatted logs using Logrus
- **Metrics**: Real-time processing statistics
- **Health Checks**: Service health monitoring
- **Performance Tracking**: Processing rates and error tracking

## Development

### Adding New Features
1. Define interfaces in the appropriate service package
2. Implement the functionality following Go best practices
3. Add corresponding API endpoints
4. Update the frontend to consume new endpoints

### Testing
```bash
# Run Go tests
go test ./...

# Run frontend tests
cd frontend && npm test
```

### Building
```bash
# Build Go binary
go build -o bin/crypto-insights cmd/main.go

# Build frontend
cd frontend && npm run build
```

## Deployment

### Production Considerations
- **Environment Variables**: Use proper environment variable management
- **Database**: Configure production PostgreSQL with proper security
- **Redis**: Set up Redis with authentication and persistence
- **Monitoring**: Implement application monitoring and alerting
- **SSL/TLS**: Configure HTTPS for production deployment

### Scaling
- **Horizontal Scaling**: Run multiple instances behind a load balancer
- **Database**: Consider read replicas for analytics queries
- **Caching**: Implement Redis clustering for high availability

## Troubleshooting

### Common Issues
1. **Database Connection**: Verify PostgreSQL is running and accessible
2. **Redis Connection**: Check Redis server status and configuration
3. **Blockchain RPC**: Ensure Base L2 RPC endpoint is accessible
4. **Port Conflicts**: Verify ports 8080, 5432, and 6379 are available

### Logs
- **Application Logs**: Check Go service logs for errors
- **Database Logs**: Review PostgreSQL logs for connection issues
- **Frontend Logs**: Check browser console for API errors

## Contributing

1. Follow Go coding standards and conventions
2. Add tests for new functionality
3. Update documentation for API changes
4. Use meaningful commit messages

## License

MIT License - see LICENSE file for details
