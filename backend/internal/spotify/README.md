# Spotify Client

A comprehensive Spotify client wrapper built on top of the official `github.com/zmb3/spotify/v2` library. This client provides organized services for different Spotify API operations and supports both client credentials and authorization code flows.

## Features

- **Multiple Authentication Methods**:
  - Client Credentials Flow (for public data)
  - Authorization Code Flow (for user-specific data)
  - PKCE support for enhanced security
- **Organized Service Architecture**:
  - SearchService: Track, album, artist, and playlist search
  - UserService: User profiles and playlists
  - PlaylistService: Playlist operations and featured playlists
  - TrackService: Track details, audio features, and recommendations
  - AlbumService: Album details and tracks
  - ArtistService: Artist information and related content
- **Built-in Pagination Support**: Helper methods for handling paginated responses
- **Type-safe API**: Leverages the official Spotify library's type safety

## Environment Variables

Set these environment variables in your `.env` file:

```env
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
SPOTIFY_REDIRECT_URL=http://localhost:8080/callback  # Optional, defaults to this
```

## Quick Start

### Client Credentials Flow (Public Data)

```go
package main

import (
    "context"
    "log"

    "github.com/daedal00/muse/backend/internal/spotify"
    spotifyapi "github.com/zmb3/spotify/v2"
)

func main() {
    ctx := context.Background()

    // Initialize client from environment
    client := spotify.NewClientFromEnv()

    // Get client credentials client (for public data)
    spotifyClient, err := client.GetClientCredentialsClient(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Create services
    services := spotify.NewServices(spotifyClient)

    // Search for tracks
    results, err := services.Search.SearchTracks(ctx, "holiday", spotifyapi.Limit(5))
    if err != nil {
        log.Fatal(err)
    }

    // Use the results...
}
```

### Authorization Code Flow (User Data)

```go
func handleAuthFlow() {
    ctx := context.Background()
    client := spotify.NewClientFromEnv()

    // Step 1: Generate auth URL
    state := "random-state-string"
    authURL := client.GetAuthURL(state)

    // Redirect user to authURL...

    // Step 2: Handle callback and exchange code
    token, err := client.ExchangeCode(ctx, authCode)
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Create authenticated client
    spotifyClient := client.GetAuthorizedClient(ctx, token)
    services := spotify.NewServices(spotifyClient)

    // Get current user
    user, err := services.User.GetCurrentUser(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Logged in as: %s", user.DisplayName)
}
```

## Service Examples

### Search Operations

```go
services := spotify.NewServices(spotifyClient)

// Search tracks
tracks, err := services.Search.SearchTracks(ctx, "query", spotifyapi.Limit(10))

// Search albums
albums, err := services.Search.SearchAlbums(ctx, "query", spotifyapi.Limit(10))

// Search artists
artists, err := services.Search.SearchArtists(ctx, "query", spotifyapi.Limit(10))

// Search playlists
playlists, err := services.Search.SearchPlaylists(ctx, "query", spotifyapi.Limit(10))

// Multi-type search
results, err := services.Search.Search(ctx, "query",
    spotifyapi.SearchTypeTrack|spotifyapi.SearchTypeArtist)
```

### User Operations

```go
// Get current user (requires authorization)
user, err := services.User.GetCurrentUser(ctx)

// Get user profile (public data)
profile, err := services.User.GetUserProfile(ctx, spotifyapi.ID("username"))

// Get user's playlists
playlists, err := services.User.GetCurrentUserPlaylists(ctx, spotifyapi.Limit(20))

// Get another user's public playlists
userPlaylists, err := services.User.GetUserPlaylists(ctx, "username", spotifyapi.Limit(20))
```

### Playlist Operations

```go
// Get playlist details
playlist, err := services.Playlist.GetPlaylist(ctx, spotifyapi.ID("playlist_id"))

// Get playlist items with pagination
items, err := services.Playlist.GetPlaylistItems(ctx, spotifyapi.ID("playlist_id"))

// Get featured playlists
message, playlists, err := services.Playlist.GetFeaturedPlaylists(ctx, spotifyapi.Limit(10))
```

### Track Operations

```go
// Get track details
track, err := services.Track.GetTrack(ctx, spotifyapi.ID("track_id"))

// Get audio features
features, err := services.Track.GetAudioFeatures(ctx, trackID1, trackID2)

// Get recommendations
recommendations, err := services.Track.GetRecommendations(ctx,
    spotifyapi.Seeds{
        Artists: []spotifyapi.ID{"artist_id"},
        Genres:  []string{"rock", "pop"},
        Tracks:  []spotifyapi.ID{"track_id"},
    },
    spotifyapi.NewTrackAttributes().
        MaxValence(0.8).
        TargetEnergy(0.6),
    spotifyapi.Country("US"),
    spotifyapi.Limit(20),
)
```

### Album Operations

```go
// Get album details
album, err := services.Album.GetAlbum(ctx, spotifyapi.ID("album_id"))

// Get album tracks
tracks, err := services.Album.GetAlbumTracks(ctx, spotifyapi.ID("album_id"),
    spotifyapi.Market("US"))
```

### Artist Operations

```go
// Get artist details
artist, err := services.Artist.GetArtist(ctx, spotifyapi.ID("artist_id"))

// Get artist's top tracks
topTracks, err := services.Artist.GetArtistTopTracks(ctx, spotifyapi.ID("artist_id"), "US")

// Get artist's albums
albums, err := services.Artist.GetArtistAlbums(ctx, spotifyapi.ID("artist_id"),
    []spotifyapi.AlbumType{spotifyapi.AlbumTypeAlbum, spotifyapi.AlbumTypeSingle})
```

## Pagination

### Automatic Pagination Helper

```go
paginationHelper := spotify.NewPaginationHelper(spotifyClient)

// Get all playlist items across all pages
allItems, err := paginationHelper.GetAllPlaylistItems(ctx, spotifyapi.ID("playlist_id"))
```

### Manual Pagination

```go
// Get first page
items, err := services.Playlist.GetPlaylistItems(ctx, spotifyapi.ID("playlist_id"))

// Iterate through pages
for {
    // Process current page items...
    for _, item := range items.Items {
        // Handle item...
    }

    // Get next page
    err = services.NextPlaylistPage(ctx, items)
    if err == spotifyapi.ErrNoMorePages {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Options

### Custom Configuration

```go
client := spotify.NewClient(spotify.Config{
    ClientID:     "your_client_id",
    ClientSecret: "your_client_secret",
    RedirectURL:  "http://localhost:3000/callback",
    Scopes: []string{
        spotifyauth.ScopeUserReadPrivate,
        spotifyauth.ScopePlaylistReadPrivate,
        // Add more scopes as needed...
    },
})
```

### Available Scopes

Common scopes you might need:

- `spotifyauth.ScopeUserReadPrivate`: Read user profile
- `spotifyauth.ScopeUserReadEmail`: Read user email
- `spotifyauth.ScopePlaylistReadPrivate`: Read private playlists
- `spotifyauth.ScopePlaylistReadCollaborative`: Read collaborative playlists
- `spotifyauth.ScopeUserLibraryRead`: Read user's saved tracks
- `spotifyauth.ScopeUserTopRead`: Read user's top tracks and artists

## Error Handling

```go
results, err := services.Search.SearchTracks(ctx, "query")
if err != nil {
    // Handle specific errors
    if spotifyErr, ok := err.(*spotifyapi.Error); ok {
        switch spotifyErr.Status {
        case 401:
            // Handle authentication error
        case 403:
            // Handle forbidden error
        case 429:
            // Handle rate limiting
        default:
            // Handle other API errors
        }
    }
    // Handle other errors
}
```

## Integration with GraphQL

This client is designed to work well with GraphQL resolvers. You can integrate it into your GraphQL schema by:

1. Creating a Spotify client instance in your resolver constructor
2. Using the services in your field resolvers
3. Leveraging the pagination helpers for GraphQL pagination

Example GraphQL integration:

```go
type Resolver struct {
    spotifyServices *spotify.Services
}

func NewResolver(clientID, clientSecret string) *Resolver {
    client := spotify.NewClient(spotify.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        // ... other config
    })

    spotifyClient, _ := client.GetClientCredentialsClient(context.Background())

    return &Resolver{
        spotifyServices: spotify.NewServices(spotifyClient),
    }
}

func (r *queryResolver) SearchTracks(ctx context.Context, query string) ([]*Track, error) {
    results, err := r.spotifyServices.Search.SearchTracks(ctx, query, spotify.Limit(20))
    if err != nil {
        return nil, err
    }

    // Convert Spotify tracks to your GraphQL types...
    return convertTracks(results.Tracks.Tracks), nil
}
```

## Examples

See `example.go` for comprehensive usage examples including:

- Basic client credentials usage
- Authorization flow implementation
- Pagination handling
- Audio features analysis
- User profile operations

Run examples with your Spotify credentials set in environment variables.
