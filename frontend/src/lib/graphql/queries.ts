import { gql } from "@apollo/client";

export const GET_ME = gql`
  query GetMe {
    me {
      id
      name
      email
      bio
      avatar
    }
  }
`;

export const SEARCH_ALBUMS = gql`
  query SearchAlbums($input: AlbumSearchInput!) {
    searchAlbums(input: $input) {
      id
      title
      artist {
        id
        name
        externalSource
      }
      releaseDate
      coverImage
      externalSource
    }
  }
`;

export const SEARCH_ARTISTS = gql`
  query SearchArtists($input: ArtistSearchInput!) {
    searchArtists(input: $input) {
      id
      name
      externalSource
    }
  }
`;

export const GET_ALBUMS = gql`
  query GetAlbums($first: Int, $after: String) {
    albums(first: $first, after: $after) {
      totalCount
      edges {
        cursor
        node {
          id
          title
          artist {
            id
            name
          }
          releaseDate
          coverImage
        }
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
`;

export const GET_REVIEWS = gql`
  query GetReviews($first: Int, $after: String) {
    reviews(first: $first, after: $after) {
      totalCount
      edges {
        cursor
        node {
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
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
`;

export const GET_PLAYLISTS = gql`
  query GetPlaylists($first: Int, $after: String) {
    playlists(first: $first, after: $after) {
      totalCount
      edges {
        cursor
        node {
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
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
`;

// Add Spotify queries
export const GET_SPOTIFY_AUTH_URL = gql`
  query GetSpotifyAuthURL {
    spotifyAuthURL {
      url
    }
  }
`;

export const GET_SPOTIFY_PLAYLISTS = gql`
  query GetSpotifyPlaylists($input: SpotifyPlaylistsInput) {
    spotifyPlaylists(input: $input) {
      playlists {
        id
        name
        description
        image
        trackCount
        isPublic
      }
      totalCount
      hasNextPage
    }
  }
`;

export const GET_ALBUM_DETAILS = gql`
  query GetAlbumDetails($id: ID!) {
    albumDetails(id: $id) {
      id
      spotifyID
      title
      artist {
        id
        name
      }
      releaseDate
      coverImage
      tracks {
        id
        spotifyID
        title
        duration
        trackNumber
        featuredArtists {
          id
          name
        }
        averageRating
        totalReviews
      }
      averageRating
      totalReviews
      reviews(first: 10) {
        totalCount
        edges {
          node {
            id
            rating
            reviewText
            createdAt
            user {
              id
              name
            }
          }
        }
      }
    }
  }
`;

export const SEARCH_TRACKS = gql`
  query SearchTracks($input: TrackSearchInput!) {
    searchTracks(input: $input) {
      id
      title
      duration
      trackNumber
      album {
        id
        title
        artist {
          id
          name
        }
        releaseDate
        coverImage
        externalSource
      }
      artists {
        id
        name
        externalSource
      }
      externalSource
    }
  }
`;

export const GET_ARTIST = gql`
  query GetArtist($id: ID!) {
    artist(id: $id) {
      id
      spotifyID
      name
      albums(first: 20) {
        totalCount
        edges {
          node {
            id
            title
            releaseDate
            coverImage
          }
        }
      }
    }
  }
`;

export const GET_ARTIST_DETAILS = gql`
  query GetArtistDetails($id: ID!) {
    artistDetails(id: $id) {
      id
      spotifyID
      name
      image
      followers
      genres
      albums {
        id
        title
        releaseDate
        coverImage
      }
      topTracks {
        id
        title
        duration
        trackNumber
        album {
          id
          title
          coverImage
        }
      }
    }
  }
`;
