# Stocet Core Service

The core backend service for the Stocet platform, providing RBAC, authentication, and all API endpoints for the frontend.

## Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control (Admin, User, Researcher)
- **News Endpoints**: Read news articles from MongoDB (acting as the sole gateway for aggregated news)
- **Market Data**: Access to market prices, tickers, price history, and technical indicators
- **Brokers Management**: CRUD operations for broker information (Admin only for create/update/delete)
- **Education**: Courses, lessons, quizzes, and glossary management
- **Chat Service**: REST API and WebSocket support for real-time messaging
- **Alerts & Watchlist**: User alerts for price/volume monitoring, and customizable stock watchlists
- **Analytics**: Dashboard metrics, market overview, and platform health
- **Admin Tools**: User management, platform analytics, and content moderation

## Architecture

The service follows clean architecture principles:

```
Routes → Controllers → Services → Repositories
```

### Database Strategy

- **PostgreSQL**: Core service data (users, brokers, education, chat, alerts)
- **MongoDB**: Read-only access to news and market data (shared with news-aggregation service)

## Tech Stack

- **Go 1.21+**: Programming language
- **Gin**: HTTP web framework
- **GORM**: ORM for PostgreSQL
- **MongoDB Driver**: For reading news/market data
- **JWT**: Authentication tokens
- **WebSocket**: Real-time chat support
- **Docker**: Containerization

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- MongoDB (for reading news/market data - can be external)
- Docker and Docker Compose (optional)

### Local Development

1. **Clone the repository** (if not already done)

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL** (if not using Docker):
   ```bash
   # Make sure PostgreSQL is running and accessible
   ```

5. **Start MongoDB** (or ensure news-aggregation service MongoDB is accessible):
   ```bash
   # MongoDB should be running on the configured URL
   ```

6. **Run migrations** (auto-migrated on startup):
   ```bash
   # Migrations run automatically when the service starts
   ```

7. **Start the service**:
   ```bash
   go run cmd/server/main.go
   ```

   The service will be available at `http://localhost:8080`

### Docker Deployment

1. **Build and start all services**:
   ```bash
   docker-compose up -d --build
   ```

2. **View logs**:
   ```bash
   docker-compose logs -f core-service
   ```

3. **Stop services**:
   ```bash
   docker-compose down
   ```

## Configuration

### Environment Variables

See `.env.example` for all available configuration options:

- `SERVER_PORT`: Server port (default: 8080)
- `SERVER_HOST`: Server host (default: 0.0.0.0)
- `POSTGRES_*`: PostgreSQL connection settings
- `MONGODB_URL`: MongoDB connection string
- `MONGODB_DB_NAME`: MongoDB database name
- `JWT_SECRET`: Secret key for JWT tokens (required in production)
- `ENVIRONMENT`: Environment (development/production)

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/signup` - User registration

### News
- `GET /api/v1/news` - List news articles (with filters)
- `DELETE /api/v1/news/:id` - Delete news (Admin only)

### Market Data
- `GET /api/v1/market-prices` - Get all market prices
- `GET /api/v1/market-prices/tickers` - Get ticker list
- `GET /api/v1/market-prices/:ticker` - Get specific ticker price
- `GET /api/v1/market-prices/:ticker/history` - Get price history
- `GET /api/v1/market-prices/:ticker/indicators` - Get technical indicators

### Brokers
- `GET /api/v1/brokers` - List brokers
- `GET /api/v1/brokers/:id` - Get broker details
- `POST /api/v1/brokers` - Create broker (Admin only)
- `PUT /api/v1/brokers/:id` - Update broker (Admin only)
- `DELETE /api/v1/brokers/:id` - Delete broker (Admin only)

### Education
- `GET /api/v1/education/courses` - List courses
- `GET /api/v1/education/courses/:id` - Get course details
- `GET /api/v1/education/lessons/:id` - Get lesson details
- `POST /api/v1/education/quizzes/:lesson_id/submit` - Submit quiz
- `GET /api/v1/education/glossary` - Get glossary terms
- Admin CRUD endpoints for courses, lessons, and glossary

### Chat
- `GET /api/v1/chat/channels` - List channels
- `GET /api/v1/chat/channels/:id/messages` - Get messages
- `POST /api/v1/chat/channels/:id/messages` - Send message
- `POST /api/v1/chat/channels` - Create channel (Admin only)

### Alerts
- `GET /api/v1/alerts` - Get user's alerts
- `POST /api/v1/alerts` - Create alert
- `PUT /api/v1/alerts/:id` - Update alert
- `DELETE /api/v1/alerts/:id` - Delete alert

### Analytics
- `GET /api/v1/analytics/dashboard` - Dashboard metrics
- `GET /api/v1/analytics/market-overview` - Market overview
- `GET /api/v1/analytics/user-activity/:user_id` - User activity (Admin only)
- `GET /api/v1/analytics/platform-health` - Platform health (Admin only)

### Dashboard
- `GET /api/v1/dashboard` - Get dashboard data (requires authentication)

### Admin
- `GET /api/v1/admin/users` - List users (Admin only)
- `PUT /api/v1/admin/users/:id` - Update user (Admin only)
- `DELETE /api/v1/admin/users/:id` - Delete user (Admin only)

## Authentication

Most endpoints require authentication via JWT token. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Roles

- **admin**: Full access to all endpoints, can manage users, brokers, education content, and delete news
- **user**: Standard user access to dashboard, alerts, chat, and read-only access to most content
- **researcher**: (Future role) Enhanced access for research features

## Subscription Model

The service includes infrastructure for subscription-based features:
- `free`: Basic features
- `paid`: Most features (excluding API access)
- `premium`: All features including API access

Currently, all users default to `free` subscription. The subscription model is ready for future implementation.

## Database Migrations

Migrations run automatically on service startup using GORM's AutoMigrate feature. Tables are created/updated as needed.

## Health Check

- `GET /health` - Returns service health status

## Development

### Project Structure

```
Stocet-Core-Service/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration and constants
│   ├── routes/         # HTTP route definitions
│   ├── controllers/    # Request handlers
│   ├── services/       # Business logic
│   ├── repositories/  # Data access layer
│   ├── models/         # Domain models
│   ├── middleware/    # HTTP middleware
│   ├── utils/          # Utility functions
│   └── database/       # Database setup
├── docker/             # Docker configuration
└── README.md
```

## License

MIT
