#!/bin/bash
# Quick Start Script for Docker Deployment

echo "🐳 E-Commerce Backend - Docker Quick Start"
echo "=========================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed!"
    echo "Please install Docker Desktop from: https://www.docker.com/products/docker-desktop"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed!"
    echo "Please install Docker Compose"
    exit 1
fi

echo "✅ Docker is installed"
echo "✅ Docker Compose is installed"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.docker .env
    echo "✅ .env file created"
    echo ""
    echo "⚠️  IMPORTANT: Edit .env file and add your S3 credentials!"
    echo "   You can skip S3 config for testing, but image upload won't work."
    echo ""
    read -p "Press Enter to continue or Ctrl+C to exit and edit .env..."
fi

echo ""
echo "🏗️  Building and starting services..."
echo "This may take a few minutes on first run..."
echo ""

# Build and start services
docker-compose up --build -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
if [ "$(docker-compose ps | grep -c Up)" -ge 3 ]; then
    echo ""
    echo "✅ All services are running!"
    echo ""
    echo "📍 Your application is available at:"
    echo "   - API:              http://localhost:8080"
    echo "   - Swagger UI:       http://localhost:8080/swagger/index.html"
    echo "   - Health Check:     http://localhost:8080/health"
    echo "   - phpMyAdmin:       http://localhost:8081"
    echo "   - Redis Commander:  http://localhost:8082"
    echo ""
    echo "📊 View logs:"
    echo "   docker-compose logs -f"
    echo ""
    echo "🛑 Stop services:"
    echo "   docker-compose down"
    echo ""
    echo "🎉 Happy testing!"
else
    echo ""
    echo "❌ Some services failed to start"
    echo "Check logs with: docker-compose logs"
fi
