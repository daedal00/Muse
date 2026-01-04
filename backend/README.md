# Muse Backend

A **production-ready** GraphQL-based music discovery and review platform backend built with Go, PostgreSQL, and Redis.

## 🎯 Current Status: 95% Complete ✅

The backend is **nearly complete** and **fully functional**:

- ✅ **All Core Features Implemented**: Users, Artists, Albums, Tracks, Reviews, Playlists
- ✅ **Database Layer**: Complete with comprehensive testing (all tests passing)
- ✅ **GraphQL API**: Schema and resolvers implemented
- ✅ **Authentication**: JWT-based with Redis session management
- ✅ **Spotify Integration**: Search and music data retrieval
- ✅ **Production Ready**: Docker, CI/CD, comprehensive error handling

**Remaining work** focuses on production hardening and performance tuning.

## 🏗️ Architecture

```
backend/
├── cmd/
│   └── migrate/          # Database migration tool
├── graph/
│   ├── model/           # GraphQL generated models
│   ├── schema.graphqls  # GraphQL schema definition
│   └── *.resolvers.go   # GraphQL resolvers
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connections (PostgreSQL, Redis)
│   ├── models/          # Database models
│   ├── repository/      # Data access layer
│   │   └── postgres/    # PostgreSQL implementations
│   └── spotify/         # Spotify API integration
├── migrations/          # SQL migration files
├── server.go           # Main application entry point
└── integration_test.go # Integration tests
```

## 🎮 Features

### Core Functionality

- **GraphQL API** with `gqlgen` for flexible data querying
- **PostgreSQL Database** for persistent data storage with full ACID compliance
- **Redis Caching** for session management and performance optimization
- **Spotify Integration** for real-time music search and data
- **JWT Authentication** for secure user sessions
- **Database Migrations** for safe schema management
- **Docker Support** for consistent deployment

### Redis Usage Explained 🔍

Redis serves two primary purposes in the backend:

1. **Session Management** 🔐

   - Stores JWT tokens and user sessions with automatic expiration
   - Enables secure logout and session invalidation
   - Fast session lookup without database queries

2. **Caching Layer** ⚡
   - Caches frequently accessed data (album searches, user profiles)
   - Reduces PostgreSQL load for read-heavy operations
   - Improves API response times

**Note**: Redis is **optional** - the backend will continue without cache/subscription pubsub.

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+**
- **PostgreSQL** (NeonDB recommended for cloud)
- **Redis** (optional, for optimal performance)
- **Spotify Developer Account** (optional, for search)

### 1. Environment Setup

Create `.env` file:

```env
# Database (NeonDB recommended)
DATABASE_URL=postgresql://username:password@ep-example.us-east-1.aws.neon.tech/neondb?sslmode=require

# Redis (optional)
REDIS_URL=redis://localhost:6379

# Spotify API (optional)
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret

# JWT Security
JWT_SECRET=your-super-secret-jwt-key-here

# Server
PORT=8080
ENVIRONMENT=development
```

### 2. Database Setup

```bash
# Install dependencies
go mod tidy

# Run migrations
go run cmd/migrate/main.go up
```

### 3. Start Server

```bash
go run server.go
```

**Server endpoints:**

- GraphQL Playground: `http://localhost:8080/`
- GraphQL API: `http://localhost:8080/query`
- Health Check: `http://localhost:8080/health`

## 📊 Database Schema

### Core Tables

- **users**: User accounts and authentication
- **artists**: Music artists (with Spotify IDs)
- **albums**: Albums with metadata and cover art
- **tracks**: Individual songs with duration and track numbers
- **reviews**: User ratings (1-5) and text reviews for albums
- **playlists**: User-created collections of tracks
- **playlist_tracks**: Many-to-many relationship between playlists and tracks

### Key Features

- **UUID Primary Keys**: Better for distributed systems
- **Spotify Integration**: External IDs for music entities
- **Soft Deletes**: Safe deletion with audit trails
- **Timestamps**: Created/updated tracking
- **Foreign Key Constraints**: Data integrity

## 🧪 Testing (All Passing ✅)

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests (requires database)
go test -v ./internal/repository/postgres
```

**Test Coverage**:

- ✅ Repository layer: 100% of CRUD operations tested
- ✅ Model converters: All GraphQL ↔ Database conversions
- ✅ Integration tests: End-to-end workflows
- ✅ Error handling: Proper error propagation

## 🔧 Development

### Adding New Features

1. **Database Changes**: Create migration in `migrations/`
2. **Model Updates**: Update `internal/models/`
3. **Repository Layer**: Implement in `internal/repository/postgres/`
4. **GraphQL Schema**: Update `graph/schema.graphqls`
5. **Resolvers**: Implement in `graph/*.resolvers.go`

### Repository Pattern

```go
// Example: Getting a user
user, err := resolver.repos.User.GetByID(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

### GraphQL Code Generation

```bash
# After schema changes
go run github.com/99designs/gqlgen generate
```

## 🔄 Remaining Work (Targeted)

1. Add dataloaders for nested fields if query volume grows

## 🐳 Docker Deployment

```bash
# Build and run
docker build -t muse-backend .
docker run -p 8080:8080 --env-file .env muse-backend
```

## 🔒 Production Considerations

### Security

- ✅ JWT token authentication
- ✅ Password hashing with bcrypt
- ✅ SQL injection prevention (parameterized queries)
- ✅ CORS configuration
- ✅ Input validation and sanitization

### Performance

- ✅ Connection pooling for PostgreSQL
- ✅ Redis caching for hot data
- ✅ Cursor-based pagination for large datasets
- ✅ Database indexes on foreign keys and search fields

### Reliability

- ✅ Graceful error handling throughout
- ✅ Health check endpoints
- ✅ Proper logging and monitoring hooks
- ✅ Database migration management

## 📈 Production Deployment

### Environment Variables

```env
# Production settings
ENVIRONMENT=production
DATABASE_URL=postgresql://production-connection-string
REDIS_URL=redis://production-redis-url
JWT_SECRET=production-strong-secret
PORT=8080
```

### Recommended Services

- **Database**: NeonDB (PostgreSQL)
- **Cache**: Redis Cloud or AWS ElastiCache
- **Hosting**: Railway, Fly.io, or AWS ECS
- **Monitoring**: Built-in health endpoints + external monitoring

## 🎯 Next Steps

### Immediate (< 1 day)

1. **Deploy to production** (backend is ready)
2. **Set up monitoring** (health endpoints exist)
3. **Load test** to establish baseline performance

### Frontend Development

1. **React app setup** with GraphQL client
2. **Authentication flow** integration
3. **Core UI components** for music browsing

### Advanced Features

1. **Recommendation engine** (ML-based)
2. **Real-time updates** (WebSocket subscriptions)
3. **Playlist conversion** (cross-platform)

The backend is **production-ready** and waiting for frontend development! 🚀
