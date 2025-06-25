# Testing Your GraphQL Backend with Spotify Integration

## Setup

1. **Get Spotify Credentials**:

   - Go to https://developer.spotify.com/dashboard
   - Create a new application
   - Copy your Client ID and Client Secret

2. **Create `.env` file** in the `backend/` directory:

```env
SPOTIFY_CLIENT_ID=your_client_id_here
SPOTIFY_CLIENT_SECRET=your_client_secret_here
SPOTIFY_REDIRECT_URL=http://localhost:8080/callback
PORT=8080
```

3. **Start the server**:

```bash
cd backend
go run .
```

4. **Open GraphQL Playground**: http://localhost:8080

## Test Queries

### 1. Create a User

```graphql
mutation {
  createUser(
    name: "Test User"
    email: "test@example.com"
    password: "password123"
  ) {
    id
    name
    email
  }
}
```

### 2. Login

```graphql
mutation {
  login(email: "test@example.com", password: "password123")
}
```

### 3. Search Albums (Spotify Integration)

```graphql
query {
  searchAlbums(input: { query: "Abbey Road", limit: 5 }) {
    id
    title
    artist {
      name
    }
    releaseDate
    coverImage
    externalSource
  }
}
```

### 4. Search Artists (Spotify Integration)

```graphql
query {
  searchArtists(input: { query: "The Beatles", limit: 5 }) {
    id
    name
    externalSource
  }
}
```

### 5. Get Current User (with auth token)

```graphql
# Add to headers: { "Authorization": "Bearer YOUR_JWT_TOKEN" }
query {
  me {
    id
    name
    email
  }
}
```

### 6. Create a Review (with auth token)

```graphql
# First you'll need to create an album in your local store
# For now this will fail until you implement album creation from Spotify data
mutation {
  createReview(
    input: { albumId: "some-album-id", rating: 5, reviewText: "Amazing album!" }
  ) {
    id
    rating
    reviewText
    user {
      name
    }
  }
}
```

## What's Working

✅ **User Management**: Create users, login, JWT authentication
✅ **Spotify Search**: Search for albums and artists via Spotify API  
✅ **GraphQL Playground**: Full GraphQL introspection and testing
✅ **Authentication Middleware**: JWT token validation
✅ **Basic CRUD**: Create playlists, reviews (for local data)

## What Needs Implementation

🔧 **Album/Track Management**: Create local albums from Spotify data
🔧 **Playlist Track Management**: Add Spotify tracks to user playlists  
🔧 **User Profile Integration**: Link Spotify profiles to local users
🔧 **Advanced Features**: Recommendations, audio features, etc.
🔧 **Pagination**: Implement GraphQL cursor-based pagination
🔧 **Subscriptions**: Real-time updates for reviews, etc.

## Architecture Summary

Your GraphQL backend successfully:

1. **Accepts GraphQL queries** from your frontend
2. **Authenticates users** with JWT tokens
3. **Fetches data from Spotify** using the official API client
4. **Stores app-specific data** (users, reviews, playlists) in memory
5. **Combines external + internal data** in unified GraphQL responses
6. **Provides type-safe API** with full GraphQL introspection

This is a **solid foundation** for a music review/playlist application that leverages Spotify's data while adding your own unique features!
