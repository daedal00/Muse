package model

import (
	"testing"
)

func TestUser_Fields(t *testing.T) {
	user := User{
		ID:    "user-123",
		Name:  "John Doe",
		Email: "john@example.com",
		Bio:   stringPtr("A music lover"),
		Avatar: stringPtr("https://example.com/avatar.jpg"),
	}

	if user.ID != "user-123" {
		t.Errorf("Expected ID 'user-123', got %s", user.ID)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got %s", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got %s", user.Email)
	}

	if user.Bio == nil || *user.Bio != "A music lover" {
		t.Errorf("Expected Bio 'A music lover', got %v", user.Bio)
	}
}

func TestAlbum_Fields(t *testing.T) {
	album := Album{
		ID:          "album-123",
		SpotifyID:   stringPtr("spotify-123"),
		Title:       "Test Album",
		ReleaseDate: stringPtr("2023-01-01"),
		CoverImage:  stringPtr("https://example.com/cover.jpg"),
	}

	if album.ID != "album-123" {
		t.Errorf("Expected ID 'album-123', got %s", album.ID)
	}

	if album.Title != "Test Album" {
		t.Errorf("Expected Title 'Test Album', got %s", album.Title)
	}

	if album.SpotifyID == nil || *album.SpotifyID != "spotify-123" {
		t.Errorf("Expected SpotifyID 'spotify-123', got %v", album.SpotifyID)
	}
}

func TestTrack_Fields(t *testing.T) {
	track := Track{
		ID:          "track-123",
		SpotifyID:   stringPtr("spotify-track-123"),
		Title:       "Test Track",
		Duration:    int32Ptr(180),
		TrackNumber: int32Ptr(1),
	}

	if track.ID != "track-123" {
		t.Errorf("Expected ID 'track-123', got %s", track.ID)
	}

	if track.Title != "Test Track" {
		t.Errorf("Expected Title 'Test Track', got %s", track.Title)
	}

	if track.Duration == nil || *track.Duration != 180 {
		t.Errorf("Expected Duration 180, got %v", track.Duration)
	}

	if track.TrackNumber == nil || *track.TrackNumber != 1 {
		t.Errorf("Expected TrackNumber 1, got %v", track.TrackNumber)
	}
}

func TestExternalSource_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		source ExternalSource
		valid  bool
	}{
		{
			name:   "valid spotify source",
			source: ExternalSourceSpotify,
			valid:  true,
		},
		{
			name:   "valid musicbrainz source",
			source: ExternalSourceMusicbrainz,
			valid:  true,
		},
		{
			name:   "invalid source",
			source: ExternalSource("INVALID"),
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.source.IsValid(); got != tt.valid {
				t.Errorf("ExternalSource.IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestExternalSource_String(t *testing.T) {
	tests := []struct {
		name     string
		source   ExternalSource
		expected string
	}{
		{
			name:     "spotify source",
			source:   ExternalSourceSpotify,
			expected: "SPOTIFY",
		},
		{
			name:     "musicbrainz source",
			source:   ExternalSourceMusicbrainz,
			expected: "MUSICBRAINZ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.source.String(); got != tt.expected {
				t.Errorf("ExternalSource.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func BenchmarkExternalSource_IsValid(b *testing.B) {
	source := ExternalSourceSpotify

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = source.IsValid()
	}
}

func BenchmarkExternalSource_String(b *testing.B) {
	source := ExternalSourceSpotify

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = source.String()
	}
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
} 