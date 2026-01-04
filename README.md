# Muse - Music Discovery & Review Platform

A music discovery and review platform inspired by Letterboxd, but for music. Rate albums, create playlists, discover new music, and share your musical journey with others.

## 🎵 What is Muse?

Muse allows users to:

- **Rate and review albums** (1-5 stars) with detailed reviews
- **Create and share playlists** with other users
- **Discover new music** through personalized recommendations
- **Track your musical journey** with a personal "My Muse" profile
- **Search and explore** artists, albums, and tracks via Spotify integration
- **Convert playlists** between different streaming services (planned)

## 🏗️ Current Architecture

### Backend (Implemented)

- **GraphQL API** built with Go and `gqlgen`
- **PostgreSQL** database for all persistent data
- **Redis** for session management and caching (optional)
- **Spotify API integration** for music data and search (optional for non-search use cases)
- **JWT Authentication** for secure user sessions
- **Docker support** for easy deployment

### Frontend (Implemented for API Testing)

- **Next.js (TypeScript)** API testing UI (basic pages)
- **React Native** for mobile apps (planned)

## 📁 Project Structure

```
Muse/
├── backend/                 # Go GraphQL backend (implemented)
│   ├── graph/              # GraphQL schema and resolvers
│   ├── internal/           # Internal packages
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connections (PostgreSQL, Redis)
│   │   ├── models/         # Database models
│   │   ├── repository/     # Data access layer
│   │   └── spotify/        # Spotify API integration
│   ├── migrations/         # Database migrations
│   └── server.go          # Main application entry
├── frontend/               # Next.js frontend (API testing UI)
└── mobile/                # React Native app (planned)
```

## 🚀 Current Implementation Status

### ✅ Completed Backend Features

- **Core Data Models**: Users, Artists, Albums, Tracks, Reviews, Playlists
- **GraphQL API**: Complete schema with all core types and operations
- **Database Layer**: Full CRUD operations for all entities
- **Authentication**: JWT-based user authentication
- **Spotify Integration**: Search artists, albums, and tracks
- **Session Management**: Redis repositories implemented (JWT flow does not persist sessions yet)
- **Comprehensive Testing**: All repositories tested and working
- **CI/CD Pipeline**: GitHub Actions with testing and deployment

### 🔧 Backend - Nearly Complete

- **GraphQL Resolvers**: Core queries/mutations implemented
- **Pagination**: Cursor-based pagination implemented
- **Error Handling**: Comprehensive error handling throughout

### ⚠️ Known Gaps

- **Product UI**: The current frontend is a testing harness, not a full UX

### 📋 Planned Features

- **Frontend Development**: Product UI (current frontend is a testing harness)
- **Mobile App**: React Native for iOS/Android
- **Playlist Conversion**: Convert playlists between streaming services
- **Advanced Recommendations**: ML-based music recommendations
- **Social Features**: Follow users, share activity feeds
- **Real-time Updates**: WebSocket support for live updates

## 🛠️ Tech Stack

| Component          | Technology                     | Status         |
| ------------------ | ------------------------------ | -------------- |
| **Backend API**    | Go + GraphQL (gqlgen)          | ✅ Implemented |
| **Database**       | PostgreSQL                     | ✅ Implemented |
| **Caching**        | Redis                          | ✅ Implemented |
| **Authentication** | JWT                            | ✅ Implemented |
| **External APIs**  | Spotify Web API                | ✅ Implemented |
| **Testing**        | Go testing + integration tests | ✅ Implemented |
| **CI/CD**          | GitHub Actions                 | ✅ Implemented |
| **Frontend**       | Next.js (TypeScript)           | ✅ Implemented |
| **Mobile**         | React Native                   | 🔄 Planned     |
| **Deployment**     | Docker + Kubernetes            | 🔄 Planned     |

## 🎯 Core Features

### User Experience

1. **Account Management**: Register, login, profile management
2. **Music Discovery**: Search via Spotify, browse curated lists
3. **Review System**: Rate albums 1-5 stars with optional text reviews
4. **Playlist Management**: Create, organize, and share playlists
5. **Personal Profile**: "My Muse" page showing ratings, favorites, activity

### Technical Features

1. **GraphQL API**: Flexible, efficient data fetching
2. **Real-time Search**: Spotify API integration for music data
3. **Scalable Architecture**: Repository pattern, clean separation of concerns
4. **Comprehensive Testing**: Unit and integration test coverage
5. **Production Ready**: Docker deployment, CI/CD pipeline

## 📊 Data Model

### Core Entities

- **Users**: Account information, preferences, authentication
- **Artists**: Music artists with Spotify integration
- **Albums**: Album metadata, cover art, release information
- **Tracks**: Individual songs with duration, track numbers
- **Reviews**: User ratings and reviews for albums
- **Playlists**: User-created track collections

### Relationships

- Users create Reviews for Albums
- Users create Playlists containing Tracks
- Albums belong to Artists and contain Tracks
- All entities support pagination via GraphQL connections

## 🔧 Development

### Backend Setup

```bash
cd backend
go mod tidy
go run cmd/migrate/main.go up
go run server.go
```

See [`backend/README.md`](backend/README.md) for detailed setup instructions.

### Testing

```bash
cd backend
go test ./...
```

### GraphQL Playground

Once running, visit `http://localhost:8080` for the GraphQL playground.

## 🚦 Next Steps

### Immediate (Backend Completion)

1. **Product UI**: Replace the testing harness with full UX flows
2. **Set up production database** (NeonDB or similar)
3. **Deploy backend** to production environment

### Frontend Development

1. **React Web App**: Main user interface
2. **Authentication Flow**: Login/register pages
3. **Core Pages**: Search, albums, reviews, playlists
4. **Responsive Design**: Mobile-friendly web interface

### Advanced Features

1. **Mobile Apps**: React Native for iOS/Android
2. **Playlist Conversion**: Cross-platform playlist import/export
3. **Recommendation Engine**: ML-based music suggestions
4. **Social Features**: User following, activity feeds

## 📄 Documentation

- [`backend/README.md`](backend/README.md) - Backend setup and API documentation
- [`backend/IMPLEMENTATION_TODO.md`](backend/IMPLEMENTATION_TODO.md) - Remaining tasks
- [`backend/TESTING_AND_CICD.md`](backend/TESTING_AND_CICD.md) - Testing and deployment
- [`.github/workflows/ci.yml`](.github/workflows/ci.yml) - CI/CD pipeline

## 🤝 Contributing

1. Check the implementation TODO for available tasks
2. All backend repositories are implemented and tested
3. Frontend development is ready to begin
4. Follow existing patterns for consistency

## 📝 License

[MIT License](LICENSE)

---

**Muse** - Discover, review, and share the music you love. 🎵
