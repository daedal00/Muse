import React, { useState } from "react";
import { useLazyQuery } from "@apollo/client";
import { SEARCH_ALBUMS, SEARCH_ARTISTS } from "../lib/graphql/queries";

interface SearchFormProps {
  onResults?: (results: any, type: "albums" | "artists") => void;
}

const SearchForm: React.FC<SearchFormProps> = ({ onResults }) => {
  const [query, setQuery] = useState("");
  const [searchType, setSearchType] = useState<"albums" | "artists">("albums");

  const [
    searchAlbums,
    { loading: albumsLoading, data: albumsData, error: albumsError },
  ] = useLazyQuery(SEARCH_ALBUMS);
  const [
    searchArtists,
    { loading: artistsLoading, data: artistsData, error: artistsError },
  ] = useLazyQuery(SEARCH_ARTISTS);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    const searchInput = {
      query: query.trim(),
      limit: 10,
      offset: 0,
      source: "SPOTIFY" as const,
    };

    if (searchType === "albums") {
      searchAlbums({ variables: { input: searchInput } });
    } else {
      searchArtists({ variables: { input: searchInput } });
    }
  };

  const loading = albumsLoading || artistsLoading;
  const error = albumsError || artistsError;
  const results =
    searchType === "albums"
      ? albumsData?.searchAlbums
      : artistsData?.searchArtists;

  React.useEffect(() => {
    if (results && onResults) {
      onResults(results, searchType);
    }
  }, [results, searchType, onResults]);

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-4">Search Music</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex space-x-4">
          <div className="flex-1">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search for albums or artists..."
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <select
            value={searchType}
            onChange={(e) =>
              setSearchType(e.target.value as "albums" | "artists")
            }
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="albums">Albums</option>
            <option value="artists">Artists</option>
          </select>

          <button
            type="submit"
            disabled={loading || !query.trim()}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Searching..." : "Search"}
          </button>
        </div>
      </form>

      {error && (
        <div className="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          Error: {error.message}
        </div>
      )}

      {results && results.length > 0 && (
        <div className="mt-6">
          <h3 className="text-lg font-semibold mb-3">
            Search Results ({results.length})
          </h3>
          <div className="grid gap-4">
            {results.map((item: any) => (
              <div
                key={item.id}
                className="border border-gray-200 rounded-lg p-4"
              >
                {searchType === "albums" ? (
                  <div className="flex items-center space-x-4">
                    {item.coverImage && (
                      <img
                        src={item.coverImage}
                        alt={item.title}
                        className="w-16 h-16 object-cover rounded"
                      />
                    )}
                    <div>
                      <h4 className="font-semibold">{item.title}</h4>
                      <p className="text-gray-600">
                        by {item.artist.map((a: any) => a.name).join(", ")}
                      </p>
                      {item.releaseDate && (
                        <p className="text-sm text-gray-500">
                          Released:{" "}
                          {new Date(item.releaseDate).toLocaleDateString()}
                        </p>
                      )}
                      <p className="text-xs text-blue-600">
                        Source: {item.externalSource}
                      </p>
                    </div>
                  </div>
                ) : (
                  <div>
                    <h4 className="font-semibold">{item.name}</h4>
                    <p className="text-xs text-blue-600">
                      Source: {item.externalSource}
                    </p>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {results && results.length === 0 && (
        <div className="mt-4 p-3 bg-yellow-100 border border-yellow-400 text-yellow-700 rounded">
          No results found for "{query}"
        </div>
      )}
    </div>
  );
};

export default SearchForm;
