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

export const IMPORT_SPOTIFY_PLAYLIST = gql`
  mutation ImportSpotifyPlaylist($spotifyPlaylistId: String!) {
    importSpotifyPlaylist(spotifyPlaylistId: $spotifyPlaylistId) {
      id
      title
      description
      coverImage
      creator {
        id
        name
      }
      tracks {
        totalCount
      }
    }
  }
`;
