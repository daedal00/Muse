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
      }
    }
  }
`;

export const GET_TRACK = gql`
  query GetTrack($id: ID!) {
    track(id: $id) {
      id
      title
      duration
      trackNumber
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
