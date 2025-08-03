# Pantry Butler

A smart pantry management system built with Go, GraphQL, and MongoDB.

## 🚀 Features

- **Pantry Management**: Track ingredients and quantities
- **Recipe Management**: Store and manage recipes
- **User Management**: User authentication and authorization
- **GraphQL API**: Modern API with GraphQL
- **MongoDB**: NoSQL database for flexible data storage
- **Docker Support**: Containerized deployment
- **CI/CD Pipeline**: Automated testing and deployment
- **Telegram Notifications**: Real-time CI/CD notifications

## 📱 Telegram Bot Integration

This project includes Telegram bot integration for real-time notifications about:
- Code pushes and deployments
- Pull request events
- CI/CD pipeline status
- Release notifications
- Workflow failures

### Quick Setup

1. **Create a Telegram bot** via [@BotFather](https://t.me/botfather)
2. **Get your chat ID** using [@userinfobot](https://t.me/userinfobot)
3. **Add GitHub secrets**:
   - `TELEGRAM_BOT_TOKEN`: Your bot token
   - `TELEGRAM_CHAT_ID`: Your chat ID
4. **Test the integration**:
   ```bash
   ./scripts/test_telegram_bot.sh YOUR_BOT_TOKEN YOUR_CHAT_ID
   ```

📖 **Detailed Setup**: See [Telegram Bot Setup Guide](docs/TELEGRAM_BOT_SETUP.md)  
📋 **Quick Reference**: See [Telegram Quick Reference](docs/TELEGRAM_QUICK_REFERENCE.md)

## 🛠️ Development

### Prerequisites

- Go 1.21+
- MongoDB 7.0+
- Docker (optional)

### Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/pantry_butler.git
   cd pantry_butler
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up MongoDB**:
   ```bash
   docker-compose up -d mongodb
   ```

4. **Run migrations**:
   ```bash
   ./migrate.sh
   ```

5. **Start the server**:
   ```bash
   go run cmd/pantry_butler/main.go
   ```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test ./internal/usecase/test/
```

## 📚 Documentation

- [API Request Flow](docs/API_REQUEST_FLOW.md)
- [Gin Integration](docs/GIN_INTEGRATION.md)
- [Request Extraction Examples](docs/REQUEST_EXTRACTION_EXAMPLES.md)
- [Testing Guide](TESTING.md)
- [Telegram Bot Setup](docs/TELEGRAM_BOT_SETUP.md)
- [Telegram Quick Reference](docs/TELEGRAM_QUICK_REFERENCE.md)

## 🚀 Deployment

### Docker

```bash
# Build and run with Docker Compose
docker-compose up --build

# Or build manually
docker build -t pantry_butler .
docker run -p 8080:8080 pantry_butler
```

### Environment Variables

```bash
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=pantry_butler
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.