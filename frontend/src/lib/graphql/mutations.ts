import { gql } from "@apollo/client";

export const CREATE_USER = gql`
  mutation CreateUser($name: String!, $email: String!, $password: String!) {
    createUser(name: $name, email: $email, password: $password) {
      id
      name
      email
    }
  }
`;

export const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password)
  }
`;

export const CREATE_REVIEW = gql`
  mutation CreateReview($input: CreateReviewInput!) {
    createReview(input: $input) {
      id
      rating
      reviewText
      createdAt
      user {
        id
        name
      }
      album {
        id
        title
        artist {
          name
        }
      }
    }
  }
`;

export const CREATE_PLAYLIST = gql`
  mutation CreatePlaylist($input: CreatePlaylistInput!) {
    createPlaylist(input: $input) {
      id
      title
      description
      coverImage
      createdAt
      creator {
        id
        name
      }
    }
  }
`;

export const ADD_TRACK_TO_PLAYLIST = gql`
  mutation AddTrackToPlaylist($playlistId: ID!, $trackId: ID!) {
    addTrackToPlaylist(playlistId: $playlistId, trackId: $trackId) {
      id
      title
      tracks {
        totalCount
        edges {
          node {
            id
            title
            artist {
              name
            }
          }
        }
      }
    }
  }
`;

export const IMPORT_ALBUM = gql`
  mutation ImportAlbum($spotifyAlbumID: ID!) {
    importAlbum(spotifyAlbumID: $spotifyAlbumID) {
      id
      title
    }
  }
`;

export const IMPORT_TRACK = gql`
  mutation ImportTrack($spotifyTrackID: ID!) {
    importTrack(spotifyTrackID: $spotifyTrackID) {
      id
      title
      album {
        id
      }
    }
  }
`;

export const CREATE_TRACK_REVIEW = gql`
  mutation CreateTrackReview($input: CreateTrackReviewInput!) {
    createTrackReview(input: $input) {
      id
      rating
      reviewText
      createdAt
      user {
        id
        name
      }
      track {
        id
        title
      }
    }
  }
`;

export const SPOTIFY_AUTH_URL = gql`
  mutation SpotifyAuthURL($redirectURI: String) {
    spotifyAuthURL(redirectURI: $redirectURI)
  }
`;

export const DISCONNECT_SPOTIFY = gql`
  mutation DisconnectSpotify {
    disconnectSpotify
  }
`;

export const IMPORT_SPOTIFY_TOP_TRACKS = gql`
  mutation ImportSpotifyTopTracks($limit: Int, $timeRange: SpotifyTimeRange) {
    importSpotifyTopTracks(limit: $limit, timeRange: $timeRange) {
      importedTracks
      importedAlbums
      importedPlaylists
    }
  }
`;

export const IMPORT_SPOTIFY_SAVED_TRACKS = gql`
  mutation ImportSpotifySavedTracks($limit: Int, $offset: Int) {
    importSpotifySavedTracks(limit: $limit, offset: $offset) {
      importedTracks
      importedAlbums
      importedPlaylists
    }
  }
`;

export const IMPORT_SPOTIFY_PLAYLIST = gql`
  mutation ImportSpotifyPlaylist($spotifyPlaylistID: ID!) {
    importSpotifyPlaylist(spotifyPlaylistID: $spotifyPlaylistID) {
      importedTracks
      importedAlbums
      importedPlaylists
    }
  }
`;

export const UPDATE_PROFILE = gql`
  mutation UpdateProfile($input: UpdateProfileInput!) {
    updateProfile(input: $input) {
      id
      name
      bio
      avatar
    }
  }
`;

export const UPDATE_PROFILE_SETTINGS = gql`
  mutation UpdateProfileSettings($input: UpdateProfileSettingsInput!) {
    updateProfileSettings(input: $input) {
      layout
      primaryColor
      accentColor
      backgroundStyle
      backgroundValue
      pinnedAlbumIds
      pinnedTrackIds
      featuredArtistIds
      sectionsOrder
      showSpotifyStats
      showListeningHistory
      bioStyle
      customTags
    }
  }
`;
export const CREATE_COMMENT = gql`
  mutation CreateComment($input: CreateCommentInput!) {
    createComment(input: $input) {
      id
      content
      createdAt
      user {
        id
        name
        avatar
      }
    }
  }
`;
