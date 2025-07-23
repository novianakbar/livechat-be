#!/bin/bash

# LiveChat Backend Setup Script
# This script helps you set up the LiveChat backend environment

set -e

echo "🚀 LiveChat Backend Setup Script"
echo "=================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
echo "✅ Go version: $GO_VERSION"

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Please install PostgreSQL 12+ first."
    echo "   Ubuntu/Debian: sudo apt-get install postgresql postgresql-contrib"
    echo "   macOS: brew install postgresql"
    exit 1
fi

# Check if Redis is installed
if ! command -v redis-server &> /dev/null; then
    echo "❌ Redis is not installed. Please install Redis 6+ first."
    echo "   Ubuntu/Debian: sudo apt-get install redis-server"
    echo "   macOS: brew install redis"
    exit 1
fi

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "📦 Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

echo "✅ All dependencies are installed"

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy
go mod download

# Copy environment file
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env file created. Please edit it with your database credentials."
else
    echo "✅ .env file already exists"
fi

# Database setup
echo ""
echo "🗄️  Database Setup"
echo "=================="

# Source environment variables
if [ -f ".env" ]; then
    export $(cat .env | xargs)
fi

# Create database
echo "Creating database '$DB_NAME'..."
if createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" 2>/dev/null; then
    echo "✅ Database '$DB_NAME' created successfully"
else
    echo "⚠️  Database '$DB_NAME' might already exist"
fi

# Run migrations
echo "🔄 Running database migrations..."
if migrate -path migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up; then
    echo "✅ Database migrations completed successfully"
else
    echo "❌ Failed to run migrations. Please check your database connection."
    exit 1
fi

# Build application
echo ""
echo "🔨 Building application..."
if go build -o dist/livechat-be cmd/main.go; then
    echo "✅ Application built successfully"
else
    echo "❌ Failed to build application"
    exit 1
fi

# Setup complete
echo ""
echo "🎉 Setup Complete!"
echo "=================="
echo ""
echo "Next steps:"
echo "1. Edit .env file with your configuration"
echo "2. Start Redis server: redis-server"
echo "3. Start PostgreSQL server"
echo "4. Run the application: make dev"
echo ""
echo "Default credentials:"
echo "  Admin: admin@livechat.com / password"
echo "  Agent: agent1@livechat.com / password"
echo ""
echo "API will be available at: http://localhost:8080"
echo "WebSocket endpoint: ws://localhost:8080/ws/chat"
echo ""
echo "Import Postman collection from: docs/LiveChat_API.postman_collection.json"
echo ""
echo "Happy coding! 🚀"
