# Muse - Music Rating & Review Platform

## Project Overview

Muse is a full-stack music rating and review platform that integrates with Spotify to allow users to discover, rate, and discuss albums and tracks. Users can connect their Spotify account, search for music, leave ratings/reviews, and engage in discussions through comments.

---

## Architecture

### Backend (`/backend`)

- **Language:** Go 1.21+
- **GraphQL:** gqlgen (schema-first)
- **Database:** PostgreSQL (Neon cloud)
- **Cache/Sessions:** Redis (optional, graceful fallback)
- **Authentication:** JWT-based with Spotify OAuth2

### Frontend (`/frontend`)

- **Framework:** Next.js 14 (Pages Router)
- **Language:** TypeScript
- **Styling:** Vanilla CSS with dark mode support
- **GraphQL Client:** Apollo Client
- **State:** React hooks + Apollo cache

---

## Key Directories

```
/backend
├── server.go              # Main entry point, HTTP routes, middleware
├── graph/
│   ├── schema.graphqls    # GraphQL schema definition
│   ├── schema.resolvers.go# Resolver implementations
│   ├── resolver.go        # Dependency injection, initialization
│   └── model/             # Generated GraphQL models
├── internal/
│   ├── config/            # Environment configuration
│   ├── database/          # PostgreSQL and Redis connections
│   ├── models/            # Database models (User, Album, Track, Review, Comment, etc.)
│   ├── repository/        # Repository interfaces and implementations
│   │   ├── interfaces.go  # Repository contracts
│   │   ├── postgres/      # PostgreSQL implementations
│   │   └── redis/         # Redis implementations (sessions, cache)
│   └── spotify/           # Spotify API client and services
├── migrations/            # SQL migrations (001-004)
├── certs/                 # Self-signed SSL certificates for local HTTPS
└── auth/                  # Password hashing, JWT utilities

/frontend
├── src/
│   ├── pages/             # Next.js pages
│   │   ├── albums/        # Album detail pages
│   │   │   ├── [id].tsx   # Album by internal ID
│   │   │   └── spotify/[id].tsx # Spotify album preview/import
│   │   ├── tracks/        # Track detail pages
│   │   ├── search.tsx     # Search page
│   │   ├── profile.tsx    # User profile
│   │   └── auth.tsx       # Login/Register
│   ├── components/        # Reusable UI components
│   │   ├── SearchForm.tsx # Unified search with results
│   │   ├── StarRating.tsx # 5-star rating component
│   │   └── CommentSection.tsx # Comments for albums/tracks
│   └── lib/graphql/       # GraphQL queries and mutations
│       ├── queries.ts
│       └── mutations.ts
└── public/                # Static assets
```

---

## Environment Configuration

### Backend (`/backend/.env`)

```env
SPOTIFY_CLIENT_ID=<spotify_client_id>
SPOTIFY_CLIENT_SECRET=<spotify_client_secret>
SPOTIFY_REDIRECT_URL=https://127.0.0.1:8080/spotify/callback
FRONTEND_URL=http://localhost:3000
PORT=8080
JWT_SECRET=<secret>
DATABASE_URL=<neon_postgres_url>
REDIS_URL=redis://localhost:6379  # Optional
```

### Frontend (`/frontend/.env.local`)

```env
NEXT_PUBLIC_GRAPHQL_URL=https://localhost:8080/query
```

---

## Running the Application

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL (Neon cloud or local)
- Redis (optional, for caching/sessions)

### Backend

```bash
cd backend
go mod download
go run server.go
# Server runs at https://localhost:8080 (with TLS certs)
```

### Frontend

```bash
cd frontend
npm install
npm run dev
# App runs at http://localhost:3000
```

### Database Migrations

```bash
cd backend
go run ./cmd/migrate/main.go up
```

---

## Key Features

### Implemented

- ✅ User authentication (email/password + JWT)
- ✅ Spotify OAuth2 integration
- ✅ Album/Track search via Spotify API
- ✅ Import albums/tracks from Spotify to local DB
- ✅ Star ratings (1-5) for albums and tracks
- ✅ Text reviews for albums and tracks
- ✅ Comments on albums and tracks
- ✅ User profiles with review history
- ✅ Real-time subscriptions (GraphQL)
- ✅ Dark mode UI

### Database Schema (Key Tables)

- `users` - User accounts
- `artists` - Music artists
- `albums` - Albums (with Spotify ID)
- `tracks` - Tracks (with Spotify ID)
- `reviews` - Album reviews (rating + text)
- `track_reviews` - Track reviews
- `comments` - Discussion comments on albums/tracks
- `playlists` - User playlists
- `spotify_tokens` - OAuth tokens per user

---

## GraphQL Schema Highlights

### Queries

- `album(id: ID!)` - Get album with tracks, reviews, comments
- `track(id: ID!)` - Get track with reviews, comments
- `searchAlbums/searchArtists/searchTracks` - Search Spotify
- `me` - Current authenticated user

### Mutations

- `createUser`, `login` - Authentication
- `createReview`, `createTrackReview` - Submit ratings
- `createComment` - Add comments
- `importAlbum`, `importTrack` - Import from Spotify
- `spotifyAuthURL` - Get OAuth authorization URL

---

## Known Issues / TODO

### Current Blockers

- **Spotify OAuth:** "INVALID_CLIENT: Insecure redirect URI" error
  - Dashboard has `https://127.0.0.1:8080/spotify/callback` whitelisted
  - Backend sends matching URI but Spotify still rejects
  - May need to verify exact match or wait for dashboard propagation

### Future Improvements

- Add playlist management UI
- Implement user following system
- Add listening statistics from Spotify
- Improve search result display for artists
- Add pagination to search results
- Mobile-responsive refinements

---

## Development Notes

### GraphQL Code Generation

```bash
cd backend/graph
go run github.com/99designs/gqlgen generate
```

### HTTPS for Local Development

Self-signed certificates are in `/backend/certs/`. The server auto-detects and uses them.
When accessing `https://localhost:8080`, accept the self-signed certificate warning.

### Logging

Backend logs all requests with timing:

```
[REQUEST] POST /query from 127.0.0.1
[AUTH] ✅ User authenticated: <user-id>
[MUTATION] CreateReview completed - Duration: 15ms
[RESPONSE] POST /query - Status: 200 - Duration: 150ms
```

---

## File Naming Conventions

- Go: `snake_case.go`
- TypeScript/React: `PascalCase.tsx` for components, `camelCase.ts` for utilities
- GraphQL: `schema.graphqls`
- Migrations: `NNN_description.up.sql` / `.down.sql`
