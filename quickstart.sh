#!/bin/bash

# Base L2 Insights Platform - Quick Start Script
# This script will help you get the platform running quickly

set -e

echo "🚀 Base L2 Insights Platform - Quick Start"
echo "=========================================="
echo ""

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ first."
    echo "   Visit: https://nodejs.org/"
    exit 1
fi

# Check Node.js version
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
    echo "❌ Node.js version $(node -v) is too old. Please upgrade to Node.js 18+"
    exit 1
fi

echo "✅ Node.js $(node -v) detected"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. You can still run the platform locally."
    echo "   To use Docker, install Docker Desktop from: https://www.docker.com/"
    DOCKER_AVAILABLE=false
else
    echo "✅ Docker $(docker --version) detected"
    DOCKER_AVAILABLE=true
fi

# Check if Docker Compose is available
if [ "$DOCKER_AVAILABLE" = true ] && ! command -v docker-compose &> /dev/null; then
    echo "⚠️  Docker Compose is not available. You can still run the platform locally."
    DOCKER_COMPOSE_AVAILABLE=false
elif [ "$DOCKER_AVAILABLE" = true ]; then
    echo "✅ Docker Compose $(docker-compose --version) detected"
    DOCKER_COMPOSE_AVAILABLE=true
fi

echo ""
echo "📋 Setup Steps:"
echo "================"

# Step 1: Environment setup
echo ""
echo "1️⃣  Setting up environment variables..."
if [ ! -f .env ]; then
    if [ -f env.example ]; then
        cp env.example .env
        echo "   ✅ Created .env file from template"
        echo "   💡 Please edit .env file with your configuration"
    else
        echo "   ❌ env.example not found"
        exit 1
    fi
else
    echo "   ✅ .env file already exists"
fi

# Step 2: Install dependencies
echo ""
echo "2️⃣  Installing dependencies..."
if [ ! -d "node_modules" ]; then
    npm install
    echo "   ✅ Dependencies installed"
else
    echo "   ✅ Dependencies already installed"
fi

# Step 3: Test setup
echo ""
echo "3️⃣  Testing setup..."
node scripts/test-setup.js

# Step 4: Choose startup method
echo ""
echo "🚀 Choose your startup method:"
echo "==============================="

if [ "$DOCKER_COMPOSE_AVAILABLE" = true ]; then
    echo ""
    echo "🐳 Option 1: Docker Compose (Recommended)"
    echo "   - Runs PostgreSQL, Redis, and the app in containers"
    echo "   - Includes pgAdmin and Redis Commander for management"
    echo "   - Command: docker-compose up -d"
    echo ""
    echo "💻 Option 2: Local Development"
    echo "   - Requires PostgreSQL and Redis installed locally"
    echo "   - Command: npm run dev"
    echo ""
    
    read -p "Choose option (1 or 2): " choice
    
    case $choice in
        1)
            echo ""
            echo "🐳 Starting with Docker Compose..."
            echo "   Starting services in background..."
            docker-compose up -d
            
            echo ""
            echo "⏳ Waiting for services to be ready..."
            sleep 10
            
            echo ""
            echo "🔍 Checking service status..."
            docker-compose ps
            
            echo ""
            echo "🌐 Your platform should be available at:"
            echo "   - Main App: http://localhost:3000"
            echo "   - API: http://localhost:3000/api"
            echo "   - pgAdmin: http://localhost:5050"
            echo "   - Redis Commander: http://localhost:8081"
            echo ""
            echo "📊 To view logs: docker-compose logs -f"
            echo "🛑 To stop: docker-compose down"
            ;;
        2)
            echo ""
            echo "💻 Starting local development..."
            echo "   Please ensure PostgreSQL and Redis are running locally"
            echo "   Then run: npm run dev"
            ;;
        *)
            echo "❌ Invalid choice. Please run the script again."
            exit 1
            ;;
    esac
else
    echo ""
    echo "💻 Local Development Only"
    echo "   Please ensure PostgreSQL and Redis are installed and running"
    echo "   Then run: npm run dev"
fi

echo ""
echo "🎉 Setup complete! Your Base L2 Insights platform is ready."
echo ""
echo "📚 Next steps:"
echo "   - Visit the API documentation at http://localhost:3000/api"
echo "   - Check the health endpoint at http://localhost:3000/health"
echo "   - Explore blockchain data at http://localhost:3000/api/blockchain"
echo ""
echo "🔧 For development:"
echo "   - Edit .env file to customize settings"
echo "   - Modify src/ directory for custom logic"
echo "   - Add new API endpoints in src/api/routes/"
echo ""
echo "📖 For more information, see README.md"
echo ""
