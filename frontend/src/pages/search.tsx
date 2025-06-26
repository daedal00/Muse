import { useState, useCallback } from "react";
import { useRouter } from "next/router";
import { useLazyQuery, useMutation } from "@apollo/client";
import {
  SEARCH_ALBUMS,
  SEARCH_ARTISTS,
  SEARCH_TRACKS,
} from "../lib/graphql/queries";
import { CREATE_REVIEW } from "../lib/graphql/mutations";
import { debounce } from "lodash";
import Layout from "../components/Layout/Layout";

interface AlbumSearchResult {
  id: string;
  title: string;
  artist: Array<{
    id: string;
    name: string;
    externalSource: string;
  }>;
  releaseDate?: string;
  coverImage?: string;
  externalSource: string;
}

interface ArtistSearchResult {
  id: string;
  name: string;
  externalSource: string;
}

interface TrackSearchResult {
  id: string;
  title: string;
  duration?: number;
  trackNumber?: number;
  album?: {
    id: string;
    title: string;
    artist: Array<{
      id: string;
      name: string;
    }>;
    releaseDate?: string;
    coverImage?: string;
    externalSource: string;
  };
  artists: Array<{
    id: string;
    name: string;
    externalSource: string;
  }>;
  externalSource: string;
}

export default function SearchPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");
  const [searchType, setSearchType] = useState<"albums" | "artists" | "tracks">(
    "albums"
  );
  const [albumResults, setAlbumResults] = useState<AlbumSearchResult[]>([]);
  const [artistResults, setArtistResults] = useState<ArtistSearchResult[]>([]);
  const [trackResults, setTrackResults] = useState<TrackSearchResult[]>([]);

  const [searchAlbums, { loading: albumsLoading }] = useLazyQuery(
    SEARCH_ALBUMS,
    {
      onCompleted: (data) => {
        setAlbumResults(data.searchAlbums || []);
      },
      onError: (error) => {
        console.error("Album search error:", error);
        setAlbumResults([]);
      },
    }
  );

  const [searchArtists, { loading: artistsLoading }] = useLazyQuery(
    SEARCH_ARTISTS,
    {
      onCompleted: (data) => {
        setArtistResults(data.searchArtists || []);
      },
      onError: (error) => {
        console.error("Artist search error:", error);
        setArtistResults([]);
      },
    }
  );

  const [searchTracks, { loading: tracksLoading }] = useLazyQuery(
    SEARCH_TRACKS,
    {
      onCompleted: (data) => {
        setTrackResults(data.searchTracks || []);
      },
      onError: (error) => {
        console.error("Track search error:", error);
        setTrackResults([]);
      },
    }
  );

  const [createReview] = useMutation(CREATE_REVIEW, {
    onCompleted: () => {
      alert("Rating submitted successfully!");
    },
    onError: (error) => {
      console.error("Failed to create review:", error);
      alert("Failed to submit rating. Please try again.");
    },
  });

  const loading = albumsLoading || artistsLoading || tracksLoading;

  // Debounced search function
  const debouncedSearch = useCallback(
    debounce((query: string, type: "albums" | "artists" | "tracks") => {
      if (query.trim().length >= 2) {
        const searchInput = {
          query: query.trim(),
          limit: 20,
          offset: 0,
          source: "SPOTIFY" as const,
        };

        if (type === "albums") {
          searchAlbums({ variables: { input: searchInput } });
        } else if (type === "artists") {
          searchArtists({ variables: { input: searchInput } });
        } else if (type === "tracks") {
          searchTracks({ variables: { input: searchInput } });
        }
      } else {
        setAlbumResults([]);
        setArtistResults([]);
        setTrackResults([]);
      }
    }, 500),
    [searchAlbums, searchArtists, searchTracks]
  );

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const query = e.target.value;
    setSearchQuery(query);
    debouncedSearch(query, searchType);
  };

  const formatDuration = (seconds?: number) => {
    if (!seconds) return "";
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
  };

  const handleAlbumClick = (albumId: string) => {
    router.push(`/album/${albumId}`);
  };

  const handleArtistClick = (artistId: string) => {
    router.push(`/artist/${artistId}`);
  };

  // Get current results based on search type
  const getCurrentResults = () => {
    switch (searchType) {
      case "albums":
        return albumResults;
      case "artists":
        return artistResults;
      case "tracks":
        return trackResults;
      default:
        return [];
    }
  };

  const currentResults = getCurrentResults();

  return (
    <Layout>
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">
            Search Music
          </h1>

          {/* Search Type Tabs */}
          <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg mb-6">
            {["albums", "artists", "tracks"].map((type) => (
              <button
                key={type}
                onClick={() => {
                  setSearchType(type as "albums" | "artists" | "tracks");
                  if (searchQuery.trim().length >= 2) {
                    debouncedSearch(
                      searchQuery,
                      type as "albums" | "artists" | "tracks"
                    );
                  }
                }}
                className={`flex-1 py-2 px-4 rounded-md text-sm font-medium capitalize transition-colors ${
                  searchType === type
                    ? "bg-white text-blue-600 shadow"
                    : "text-gray-600 hover:text-gray-900"
                }`}
              >
                {type}
              </button>
            ))}
          </div>

          {/* Search Input */}
          <div className="relative">
            <input
              type="text"
              placeholder={`Search for ${searchType}...`}
              value={searchQuery}
              onChange={handleSearchChange}
              className="w-full px-4 py-3 pl-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg
                className="h-5 w-5 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
          </div>

          {loading && (
            <div className="flex justify-center items-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
            </div>
          )}
        </div>

        {/* Results */}
        {currentResults.length > 0 && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">
              {searchType.charAt(0).toUpperCase() + searchType.slice(1)} (
              {currentResults.length})
            </h2>

            {searchType === "albums" && (
              <div className="grid gap-4">
                {(currentResults as AlbumSearchResult[]).map((album) => (
                  <div
                    key={album.id}
                    onClick={() => handleAlbumClick(album.id)}
                    className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
                  >
                    <div className="flex items-center space-x-4">
                      {album.coverImage && (
                        <img
                          src={album.coverImage}
                          alt={album.title}
                          className="w-16 h-16 object-cover rounded"
                        />
                      )}
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900 hover:text-blue-600">
                          {album.title}
                        </h3>
                        <p className="text-gray-600">
                          by{" "}
                          {album.artist
                            .map((a) => (
                              <span
                                key={a.id}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  handleArtistClick(a.id);
                                }}
                                className="hover:text-blue-600 cursor-pointer"
                              >
                                {a.name}
                              </span>
                            ))
                            .reduce(
                              (prev, curr, index) =>
                                index === 0 ? [curr] : [...prev, ", ", curr],
                              [] as React.ReactNode[]
                            )}
                        </p>
                        {album.releaseDate && (
                          <p className="text-sm text-gray-500">
                            Released:{" "}
                            {new Date(album.releaseDate).toLocaleDateString()}
                          </p>
                        )}
                        <p className="text-xs text-blue-600 mt-1">
                          Source: Spotify
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {searchType === "artists" && (
              <div className="grid gap-4">
                {(currentResults as ArtistSearchResult[]).map((artist) => (
                  <div
                    key={artist.id}
                    onClick={() => handleArtistClick(artist.id)}
                    className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
                  >
                    <div className="flex items-center space-x-4">
                      <div className="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center">
                        <svg
                          className="w-8 h-8 text-gray-400"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path
                            fillRule="evenodd"
                            d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </div>
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900 hover:text-blue-600">
                          {artist.name}
                        </h3>
                        <p className="text-xs text-blue-600 mt-1">
                          Source: Spotify
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            {searchType === "tracks" && (
              <div className="grid gap-4">
                {(currentResults as TrackSearchResult[]).map((track) => (
                  <div
                    key={track.id}
                    className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
                  >
                    <div className="flex items-center space-x-4">
                      {track.album?.coverImage && (
                        <img
                          src={track.album.coverImage}
                          alt={track.album.title}
                          className="w-16 h-16 object-cover rounded"
                        />
                      )}
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900">
                          {track.title}
                        </h3>
                        <p className="text-gray-600">
                          by{" "}
                          {track.artists
                            .map((artist) => (
                              <span
                                key={artist.id}
                                onClick={() => handleArtistClick(artist.id)}
                                className="hover:text-blue-600 cursor-pointer"
                              >
                                {artist.name}
                              </span>
                            ))
                            .reduce(
                              (prev, curr, index) =>
                                index === 0 ? [curr] : [...prev, ", ", curr],
                              [] as React.ReactNode[]
                            )}
                        </p>
                        {track.album && (
                          <p
                            className="text-sm text-gray-500 hover:text-blue-600 cursor-pointer"
                            onClick={() => handleAlbumClick(track.album!.id)}
                          >
                            from "{track.album.title}"
                          </p>
                        )}
                        <div className="flex items-center gap-4 mt-2">
                          {track.duration && (
                            <span className="text-xs text-gray-500">
                              {formatDuration(track.duration)}
                            </span>
                          )}
                          <span className="text-xs text-blue-600">
                            Source: Spotify
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* No Results */}
        {!loading && searchQuery.length >= 2 && currentResults.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-600">
              No {searchType} found for "{searchQuery}"
            </p>
          </div>
        )}

        {/* Instructions */}
        {searchQuery.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-600">
              Enter at least 2 characters to search for {searchType}
            </p>
          </div>
        )}
      </div>
    </Layout>
  );
}
