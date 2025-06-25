package spotify

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client wraps the Spotify client with additional functionality
type Client struct {
	spotify     *spotify.Client
	auth        *spotifyauth.Authenticator
	clientID    string
	clientSecret string
}

// Config holds the configuration for the Spotify client
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// NewClient creates a new Spotify client instance
func NewClient(config Config) *Client {
	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(config.RedirectURL),
		spotifyauth.WithScopes(config.Scopes...),
	)

	return &Client{
		auth:         auth,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
	}
}

// NewClientFromEnv creates a new Spotify client from environment variables
func NewClientFromEnv() *Client {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURL := os.Getenv("SPOTIFY_REDIRECT_URL")
	
	if clientID == "" || clientSecret == "" {
		log.Fatal("SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET must be set")
	}
	
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/callback"
	}

	return NewClient(Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserTopRead,
		},
	})
}

// GetClientCredentialsClient returns a client using client credentials flow (for public data)
func (c *Client) GetClientCredentialsClient(ctx context.Context) (*spotify.Client, error) {
	config := &clientcredentials.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	
	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get token: %w", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	return spotify.New(httpClient), nil
}

// GetAuthorizedClient returns a client from an authorization code (for user-specific data)
func (c *Client) GetAuthorizedClient(ctx context.Context, token *oauth2.Token) *spotify.Client {
	httpClient := c.auth.Client(ctx, token)
	return spotify.New(httpClient)
}

// GetAuthURL generates the authorization URL for the OAuth flow
func (c *Client) GetAuthURL(state string) string {
	return c.auth.AuthURL(state)
}

// GetAuthURLWithPKCE generates the authorization URL with PKCE parameters
func (c *Client) GetAuthURLWithPKCE(state, codeChallenge string) string {
	return c.auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("client_id", c.clientID),
	)
}

// ExchangeCode exchanges an authorization code for a token
func (c *Client) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.auth.Exchange(ctx, code)
}

// ExchangeCodeWithPKCE exchanges an authorization code for a token using PKCE
func (c *Client) ExchangeCodeWithPKCE(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	return c.auth.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
}

// SearchService provides search functionality
type SearchService struct {
	client *spotify.Client
}

// Search performs a search across multiple types
func (s *SearchService) Search(ctx context.Context, query string, searchType spotify.SearchType) (*spotify.SearchResult, error) {
	return s.client.Search(ctx, query, searchType)
}

// SearchTracks searches for tracks only
func (s *SearchService) SearchTracks(ctx context.Context, query string, options ...spotify.RequestOption) (*spotify.SearchResult, error) {
	return s.client.Search(ctx, query, spotify.SearchTypeTrack, options...)
}

// SearchAlbums searches for albums only
func (s *SearchService) SearchAlbums(ctx context.Context, query string, options ...spotify.RequestOption) (*spotify.SearchResult, error) {
	return s.client.Search(ctx, query, spotify.SearchTypeAlbum, options...)
}

// SearchArtists searches for artists only
func (s *SearchService) SearchArtists(ctx context.Context, query string, options ...spotify.RequestOption) (*spotify.SearchResult, error) {
	return s.client.Search(ctx, query, spotify.SearchTypeArtist, options...)
}

// SearchPlaylists searches for playlists only
func (s *SearchService) SearchPlaylists(ctx context.Context, query string, options ...spotify.RequestOption) (*spotify.SearchResult, error) {
	return s.client.Search(ctx, query, spotify.SearchTypePlaylist, options...)
}

// UserService provides user-related functionality
type UserService struct {
	client *spotify.Client
}

// GetCurrentUser gets the current user's profile
func (u *UserService) GetCurrentUser(ctx context.Context) (*spotify.PrivateUser, error) {
	return u.client.CurrentUser(ctx)
}

// GetUserProfile gets a user's public profile
func (u *UserService) GetUserProfile(ctx context.Context, userID spotify.ID) (*spotify.User, error) {
	return u.client.GetUsersPublicProfile(ctx, userID)
}

// GetUserPlaylists gets a user's playlists
func (u *UserService) GetUserPlaylists(ctx context.Context, userID string, options ...spotify.RequestOption) (*spotify.SimplePlaylistPage, error) {
	return u.client.GetPlaylistsForUser(ctx, userID, options...)
}

// GetCurrentUserPlaylists gets the current user's playlists
func (u *UserService) GetCurrentUserPlaylists(ctx context.Context, options ...spotify.RequestOption) (*spotify.SimplePlaylistPage, error) {
	return u.client.CurrentUsersPlaylists(ctx, options...)
}

// PlaylistService provides playlist functionality
type PlaylistService struct {
	client *spotify.Client
}

// GetPlaylist gets a playlist by ID
func (p *PlaylistService) GetPlaylist(ctx context.Context, playlistID spotify.ID, options ...spotify.RequestOption) (*spotify.FullPlaylist, error) {
	return p.client.GetPlaylist(ctx, playlistID, options...)
}

// GetPlaylistItems gets items from a playlist with pagination support
func (p *PlaylistService) GetPlaylistItems(ctx context.Context, playlistID spotify.ID, options ...spotify.RequestOption) (*spotify.PlaylistItemPage, error) {
	return p.client.GetPlaylistItems(ctx, playlistID, options...)
}

// GetFeaturedPlaylists gets featured playlists
func (p *PlaylistService) GetFeaturedPlaylists(ctx context.Context, options ...spotify.RequestOption) (string, *spotify.SimplePlaylistPage, error) {
	return p.client.FeaturedPlaylists(ctx, options...)
}

// TrackService provides track functionality
type TrackService struct {
	client *spotify.Client
}

// GetTrack gets a track by ID
func (t *TrackService) GetTrack(ctx context.Context, trackID spotify.ID, options ...spotify.RequestOption) (*spotify.FullTrack, error) {
	return t.client.GetTrack(ctx, trackID, options...)
}

// GetAudioFeatures gets audio features for tracks
func (t *TrackService) GetAudioFeatures(ctx context.Context, trackIDs ...spotify.ID) ([]*spotify.AudioFeatures, error) {
	return t.client.GetAudioFeatures(ctx, trackIDs...)
}

// GetRecommendations gets track recommendations
func (t *TrackService) GetRecommendations(ctx context.Context, seeds spotify.Seeds, trackAttributes *spotify.TrackAttributes, options ...spotify.RequestOption) (*spotify.Recommendations, error) {
	return t.client.GetRecommendations(ctx, seeds, trackAttributes, options...)
}

// AlbumService provides album functionality
type AlbumService struct {
	client *spotify.Client
}

// GetAlbum gets an album by ID
func (a *AlbumService) GetAlbum(ctx context.Context, albumID spotify.ID, options ...spotify.RequestOption) (*spotify.FullAlbum, error) {
	return a.client.GetAlbum(ctx, albumID, options...)
}

// GetAlbumTracks gets tracks from an album
func (a *AlbumService) GetAlbumTracks(ctx context.Context, albumID spotify.ID, options ...spotify.RequestOption) (*spotify.SimpleTrackPage, error) {
	return a.client.GetAlbumTracks(ctx, albumID, options...)
}

// ArtistService provides artist functionality
type ArtistService struct {
	client *spotify.Client
}

// GetArtist gets an artist by ID
func (a *ArtistService) GetArtist(ctx context.Context, artistID spotify.ID) (*spotify.FullArtist, error) {
	return a.client.GetArtist(ctx, artistID)
}

// GetArtistTopTracks gets an artist's top tracks
func (a *ArtistService) GetArtistTopTracks(ctx context.Context, artistID spotify.ID, country string) ([]spotify.FullTrack, error) {
	return a.client.GetArtistsTopTracks(ctx, artistID, country)
}

// GetArtistAlbums gets an artist's albums
func (a *ArtistService) GetArtistAlbums(ctx context.Context, artistID spotify.ID, albumTypes []spotify.AlbumType, options ...spotify.RequestOption) (*spotify.SimpleAlbumPage, error) {
	return a.client.GetArtistAlbums(ctx, artistID, albumTypes, options...)
}

// Services provides access to all Spotify services
type Services struct {
	Search   *SearchService
	User     *UserService
	Playlist *PlaylistService
	Track    *TrackService
	Album    *AlbumService
	Artist   *ArtistService
	client   *spotify.Client
}

// NewServices creates service instances from a Spotify client
func NewServices(client *spotify.Client) *Services {
	return &Services{
		Search:   &SearchService{client: client},
		User:     &UserService{client: client},
		Playlist: &PlaylistService{client: client},
		Track:    &TrackService{client: client},
		Album:    &AlbumService{client: client},
		Artist:   &ArtistService{client: client},
		client:   client,
	}
}

// GetSpotifyClient returns the underlying Spotify client for direct access
func (s *Services) GetSpotifyClient() *spotify.Client {
	return s.client
}

// NextPlaylistPage fetches the next page for playlist items
func (s *Services) NextPlaylistPage(ctx context.Context, page *spotify.PlaylistItemPage) error {
	return s.client.NextPage(ctx, page)
}

// NextSimplePlaylistPage fetches the next page for simple playlists
func (s *Services) NextSimplePlaylistPage(ctx context.Context, page *spotify.SimplePlaylistPage) error {
	return s.client.NextPage(ctx, page)
}

// NextTrackPage fetches the next page for tracks
func (s *Services) NextTrackPage(ctx context.Context, page *spotify.SimpleTrackPage) error {
	return s.client.NextPage(ctx, page)
}

// NextAlbumPage fetches the next page for albums
func (s *Services) NextAlbumPage(ctx context.Context, page *spotify.SimpleAlbumPage) error {
	return s.client.NextPage(ctx, page)
}

// PaginationHelper provides utilities for handling pagination
type PaginationHelper struct {
	client *spotify.Client
}

// GetAllPlaylistItems gets all items from a playlist across all pages
func (p *PaginationHelper) GetAllPlaylistItems(ctx context.Context, playlistID spotify.ID) ([]spotify.PlaylistItem, error) {
	var allItems []spotify.PlaylistItem
	
	items, err := p.client.GetPlaylistItems(ctx, playlistID)
	if err != nil {
		return nil, err
	}
	
	allItems = append(allItems, items.Items...)
	
	for {
		err = p.client.NextPage(ctx, items)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, items.Items...)
	}
	
	return allItems, nil
}

// NewPaginationHelper creates a new pagination helper
func NewPaginationHelper(client *spotify.Client) *PaginationHelper {
	return &PaginationHelper{client: client}
}
