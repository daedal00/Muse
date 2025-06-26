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
