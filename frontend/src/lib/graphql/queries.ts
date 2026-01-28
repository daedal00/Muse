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
        releaseDate
        coverImage
        artist {
          id
          name
          externalSource
        }
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

export const GET_ALBUM = gql`
  query GetAlbum(
    $id: ID!
    $tracksFirst: Int
    $tracksAfter: String
    $reviewsFirst: Int
    $reviewsAfter: String
  ) {
    album(id: $id) {
      id
      title
      releaseDate
      coverImage
      averageRating
      artist {
        id
        name
      }
      tracks(first: $tracksFirst, after: $tracksAfter) {
        totalCount
        edges {
          cursor
          node {
            id
            title
            duration
            trackNumber
          }
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
      reviews(first: $reviewsFirst, after: $reviewsAfter) {
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
          }
        }
        pageInfo {
          hasNextPage
          endCursor
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
      comments(first: 20) {
        totalCount
        edges {
          node {
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
      }
    }
  }
`;

export const GET_TRACK = gql`
  query GetTrack($id: ID!, $reviewsFirst: Int, $reviewsAfter: String) {
    track(id: $id) {
      id
      title
      duration
      trackNumber
      averageRating
      album {
        id
        title
        releaseDate
        coverImage
        averageRating
        artist {
          id
          name
        }
      }
      reviews(first: $reviewsFirst, after: $reviewsAfter) {
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
          }
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
      comments(first: 20) {
        totalCount
        edges {
          node {
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

export const GET_PLAYLIST = gql`
  query GetPlaylist($id: ID!, $tracksFirst: Int, $tracksAfter: String) {
    playlist(id: $id) {
      id
      title
      description
      coverImage
      createdAt
      creator {
        id
        name
      }
      tracks(first: $tracksFirst, after: $tracksAfter) {
        totalCount
        edges {
          cursor
          node {
            id
            title
            duration
            trackNumber
            averageRating
            album {
              id
              title
              coverImage
              artist {
                id
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
  }
`;

export const GET_TOP_TRACKS = gql`
  query GetTopTracks($limit: Int, $sort: TopTrackSort) {
    topTracks(limit: $limit, sort: $sort) {
      averageRating
      reviewCount
      track {
        id
        title
        album {
          id
          title
          coverImage
          artist {
            name
          }
        }
      }
    }
  }
`;

export const GET_FAVORITE_TRACKS = gql`
  query GetFavoriteTracks($limit: Int, $minRating: Int) {
    favoriteTracks(limit: $limit, minRating: $minRating) {
      id
      rating
      reviewText
      createdAt
      track {
        id
        title
        album {
          id
          title
          coverImage
          artist {
            name
          }
        }
      }
    }
  }
`;

export const GET_PROFILE = gql`
  query GetProfile($trackReviewsFirst: Int, $trackReviewsAfter: String) {
    me {
      id
      name
      email
      bio
      avatar
      profileSettings {
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
      trackReviews(first: $trackReviewsFirst, after: $trackReviewsAfter) {
        totalCount
        edges {
          cursor
          node {
            id
            rating
            reviewText
            createdAt
            track {
              id
              title
              album {
                id
                title
                coverImage
                artist {
                  name
                }
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
  }
`;

export const GET_SPOTIFY_STATUS = gql`
  query GetSpotifyStatus {
    spotifyStatus {
      connected
      displayName
      scope
      expiresAt
    }
  }
`;

export const GET_SPOTIFY_TOP_TRACKS = gql`
  query GetSpotifyTopTracks($limit: Int, $timeRange: SpotifyTimeRange) {
    spotifyTopTracks(limit: $limit, timeRange: $timeRange) {
      id
      title
      duration
      trackNumber
      album {
        id
        title
        coverImage
        artist {
          name
        }
      }
      artists {
        id
        name
      }
      externalSource
    }
  }
`;

export const GET_SPOTIFY_SAVED_TRACKS = gql`
  query GetSpotifySavedTracks($limit: Int, $offset: Int) {
    spotifySavedTracks(limit: $limit, offset: $offset) {
      id
      title
      duration
      trackNumber
      album {
        id
        title
        coverImage
        artist {
          name
        }
      }
      artists {
        id
        name
      }
      externalSource
    }
  }
`;

export const GET_SPOTIFY_PLAYLISTS = gql`
  query GetSpotifyPlaylists($limit: Int, $offset: Int) {
    spotifyPlaylists(limit: $limit, offset: $offset) {
      id
      name
      description
      coverImage
      ownerName
      trackCount
      externalSource
    }
  }
`;

// User queries for social features
export const GET_USERS = gql`
  query GetUsers($first: Int, $after: String) {
    users(first: $first, after: $after) {
      totalCount
      edges {
        cursor
        node {
          id
          name
          bio
          avatar
        }
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
`;

export const GET_USER = gql`
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      name
      bio
      avatar
      profileSettings {
        layout
        primaryColor
        accentColor
        backgroundStyle
        backgroundValue
        pinnedAlbumIds
        pinnedAlbums {
          id
          title
          coverImage
          artist {
            id
            name
          }
        }
        pinnedTrackIds
        pinnedTracks {
          id
          title
          album {
            id
            title
            coverImage
            artist {
              id
              name
            }
          }
        }
        featuredArtistIds
        featuredArtists {
          id
          name
        }
        sectionsOrder
        showSpotifyStats
        showListeningHistory
        bioStyle
        customTags
      }
      reviews(first: 10) {
        totalCount
        edges {
          node {
            id
            rating
            reviewText
            createdAt
            album {
              id
              title
              coverImage
              artist {
                name
              }
            }
          }
        }
      }
      trackReviews(first: 10) {
        totalCount
        edges {
          node {
            id
            rating
            reviewText
            createdAt
            track {
              id
              title
              album {
                id
                title
                coverImage
                artist {
                  name
                }
              }
            }
          }
        }
      }
    }
  }
`;

// Preview queries - view Spotify data without storing to DB
export const GET_SPOTIFY_ALBUM_PREVIEW = gql`
  query GetSpotifyAlbumPreview($spotifyID: String!) {
    spotifyAlbum(spotifyID: $spotifyID) {
      id
      title
      releaseDate
      coverImage
      externalSource
      artist {
        id
        name
        externalSource
      }
    }
  }
`;

export const GET_SPOTIFY_TRACK_PREVIEW = gql`
  query GetSpotifyTrackPreview($spotifyID: String!) {
    spotifyTrack(spotifyID: $spotifyID) {
      id
      title
      duration
      trackNumber
      externalSource
      album {
        id
        title
        coverImage
        releaseDate
        externalSource
        artist {
          id
          name
          externalSource
        }
      }
      artists {
        id
        name
        externalSource
      }
    }
  }
`;
