#!/bin/bash
# docker-dev.sh
# Easy development startup script

set -e

echo "🐳 Starting Go gRPC Backend Development Environment..."

# Check if .env exists, if not copy from example
if [ ! -f .env ]; then
    echo "📋 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ Please review and update .env file if needed"
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "🔧 Building and starting services..."
docker-compose up --build

echo "🎉 Development environment is ready!"
echo ""
echo "📱 Services available at:"
echo "  - HTTP API: http://localhost:8080"
echo "  - gRPC API: localhost:50051"
echo "  - Swagger UI: http://localhost:8080/swagger-ui/"
echo "  - PostgreSQL: localhost:5432"
echo ""
echo "🛑 To stop services: docker-compose down"
echo "🗑️  To reset database: docker-compose down -v && docker-compose up --build"