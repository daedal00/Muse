import React, { useState } from "react";
import { useLazyQuery } from "@apollo/client";
import { useRouter } from "next/router";
import {
  SEARCH_ALBUMS,
  SEARCH_ARTISTS,
  SEARCH_TRACKS,
} from "../lib/graphql/queries";

interface SearchFormProps {
  onResults?: (results: any, type: "albums" | "artists" | "tracks") => void;
}

const SearchForm: React.FC<SearchFormProps> = ({ onResults }) => {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [searchType, setSearchType] = useState<"albums" | "artists" | "tracks">(
    "albums",
  );
  const [navigatingID, setNavigatingID] = useState<string | null>(null);

  const [
    searchAlbums,
    { loading: albumsLoading, data: albumsData, error: albumsError },
  ] = useLazyQuery(SEARCH_ALBUMS);
  const [
    searchArtists,
    { loading: artistsLoading, data: artistsData, error: artistsError },
  ] = useLazyQuery(SEARCH_ARTISTS);
  const [
    searchTracks,
    { loading: tracksLoading, data: tracksData, error: tracksError },
  ] = useLazyQuery(SEARCH_TRACKS);

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
    } else if (searchType === "artists") {
      searchArtists({ variables: { input: searchInput } });
    } else {
      searchTracks({ variables: { input: searchInput } });
    }
  };

  // Navigate to preview page (no import until user takes action)
  const handleView = (type: "albums" | "tracks", spotifyID: string) => {
    setNavigatingID(spotifyID);
    if (type === "albums") {
      router.push(`/albums/spotify/${spotifyID}`);
    } else {
      router.push(`/tracks/spotify/${spotifyID}`);
    }
  };

  const formatDuration = (durationSeconds?: number) => {
    if (!durationSeconds && durationSeconds !== 0) return null;
    const minutes = Math.floor(durationSeconds / 60);
    const seconds = durationSeconds % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  const loading = albumsLoading || artistsLoading || tracksLoading;
  const error = albumsError || artistsError || tracksError;
  const results =
    searchType === "albums"
      ? albumsData?.searchAlbums
      : searchType === "artists"
        ? artistsData?.searchArtists
        : tracksData?.searchTracks;

  React.useEffect(() => {
    if (results && onResults) {
      onResults(results, searchType);
    }
  }, [results, searchType, onResults]);

  return (
    <div className="w-full">
      <form onSubmit={handleSubmit} className="mb-8">
        <div className="relative group">
          <div className="flex items-center w-full p-2 rounded-2xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 shadow-sm transition-all focus-within:ring-2 focus-within:ring-emerald-500/20 focus-within:border-emerald-500/50">
            {/* Type Selector */}
            <div className="relative flex-shrink-0">
              <select
                value={searchType}
                onChange={(e) =>
                  setSearchType(
                    e.target.value as "albums" | "artists" | "tracks",
                  )
                }
                className="appearance-none bg-slate-50 dark:bg-slate-700/50 pl-4 pr-10 py-3 rounded-xl text-slate-700 dark:text-slate-200 font-medium text-sm focus:outline-none cursor-pointer border-r border-transparent hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
                title="Select search type"
              >
                <option value="albums">Albums</option>
                <option value="artists">Artists</option>
                <option value="tracks">Tracks</option>
              </select>
              <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="m6 9 6 6 6-6" />
                </svg>
              </div>
            </div>

            <div className="w-px h-8 bg-slate-200 dark:bg-slate-700 mx-2"></div>

            {/* Input */}
            <div className="flex-1 relative">
              <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder={`Search for ${searchType}...`}
                className="w-full bg-transparent px-2 py-3 text-lg outline-none placeholder:text-slate-400 text-slate-900 dark:text-slate-100"
              />
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading || !query.trim()}
              className="ml-2 flex-shrink-0 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl px-6 py-3 font-medium transition-all shadow-md shadow-emerald-900/10 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              {loading ? (
                <>
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  <span>Searching</span>
                </>
              ) : (
                <>
                  <span>Search</span>
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <circle cx="11" cy="11" r="8" />
                    <path d="m21 21-4.3-4.3" />
                  </svg>
                </>
              )}
            </button>
          </div>
        </div>
      </form>

      {error && (
        <div className="mt-4 rounded-xl border border-rose-200 bg-rose-50 p-3 text-rose-700 dark:border-rose-500/30 dark:bg-rose-500/10 dark:text-rose-200">
          Error: {error.message}
        </div>
      )}

      {results && results.length > 0 && (
        <div className="mt-6">
          <p className="text-sm text-slate-500 dark:text-slate-400 mb-3">
            {results.length} results
          </p>
          <div className="space-y-2">
            {results.map((item: any) => (
              <div
                key={item.id}
                className="flex items-center gap-3 p-3 rounded-xl bg-slate-50 dark:bg-slate-800/50 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer group"
                onClick={() =>
                  searchType === "albums"
                    ? handleView("albums", item.id)
                    : searchType === "tracks"
                      ? handleView("tracks", item.id)
                      : null
                }
              >
                {searchType === "albums" ? (
                  <>
                    {item.coverImage ? (
                      <img
                        src={item.coverImage}
                        alt={item.title}
                        className="w-12 h-12 object-cover rounded-lg flex-shrink-0"
                      />
                    ) : (
                      <div className="w-12 h-12 rounded-lg bg-slate-200 dark:bg-slate-700 flex-shrink-0" />
                    )}
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-base truncate">
                        {item.title}
                      </p>
                      <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                        {item.artist.map((a: any) => a.name).join(", ")}
                        {item.releaseDate &&
                          ` · ${new Date(item.releaseDate).getFullYear()}`}
                      </p>
                    </div>
                    <span className="text-sm text-emerald-600 dark:text-emerald-400 opacity-0 group-hover:opacity-100 transition-opacity">
                      {navigatingID === item.id ? "Opening..." : "View →"}
                    </span>
                  </>
                ) : searchType === "tracks" ? (
                  <>
                    {item.album?.coverImage ? (
                      <img
                        src={item.album.coverImage}
                        alt={item.album.title}
                        className="w-12 h-12 object-cover rounded-lg flex-shrink-0"
                      />
                    ) : (
                      <div className="w-12 h-12 rounded-lg bg-slate-200 dark:bg-slate-700 flex-shrink-0" />
                    )}
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-base truncate">
                        {item.title}
                      </p>
                      <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                        {item.artists.map((a: any) => a.name).join(", ")}
                        {item.album?.title && ` · ${item.album.title}`}
                      </p>
                    </div>
                    {formatDuration(item.duration) && (
                      <span className="text-sm text-slate-400 dark:text-slate-500">
                        {formatDuration(item.duration)}
                      </span>
                    )}
                    <span className="text-sm text-emerald-600 dark:text-emerald-400 opacity-0 group-hover:opacity-100 transition-opacity">
                      {navigatingID === item.id ? "Opening..." : "View →"}
                    </span>
                  </>
                ) : (
                  <>
                    <div className="w-12 h-12 rounded-lg bg-slate-200 dark:bg-slate-700 flex-shrink-0 flex items-center justify-center">
                      <span className="text-lg">🎤</span>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-base truncate">
                        {item.name}
                      </p>
                      <p className="text-sm text-slate-500 dark:text-slate-400">
                        Artist
                      </p>
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
      <div className="grid grid-cols-1 gap-4">
        {results?.map((item: any) => (
          <div
            key={item.id}
            onClick={() => {
              if (searchType === "artists") return; // Artist view not implemented yet?
              handleView(searchType as any, item.id);
            }}
            className={`flex items-center p-4 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-800/50 hover:border-emerald-500/50 hover:shadow-lg hover:shadow-emerald-500/10 transition-all cursor-pointer group ${searchType === "artists" ? "cursor-default" : ""}`}
          >
            {/* Image */}
            <div className="w-16 h-16 flex-shrink-0 bg-slate-200 dark:bg-slate-700 rounded-lg overflow-hidden mr-4">
              {item.coverImage || item.images?.[0]?.url ? (
                <img
                  src={item.coverImage || item.images?.[0]?.url}
                  alt={item.title || item.name}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-slate-400">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <circle cx="12" cy="12" r="10" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                </div>
              )}
            </div>

            {/* Info */}
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-lg text-slate-900 dark:text-slate-100 truncate pr-4 group-hover:text-emerald-500 transition-colors">
                {item.title || item.name}
              </h3>
              <div className="flex items-center text-sm text-slate-500 dark:text-slate-400 mt-1">
                {item.artist && (
                  <>
                    <span className="truncate max-w-[200px]">
                      {Array.isArray(item.artist)
                        ? item.artist.map((a: any) => a.name).join(", ")
                        : item.artist.name}
                    </span>
                    {item.releaseDate && <span className="mx-2">•</span>}
                  </>
                )}
                {item.artists && (
                  <>
                    <span className="truncate max-w-[200px]">
                      {item.artists.map((a: any) => a.name).join(", ")}
                    </span>
                    {item.album && <span className="mx-2">•</span>}
                    {item.album && (
                      <span className="truncate max-w-[150px]">
                        {item.album.title}
                      </span>
                    )}
                  </>
                )}
                {item.releaseDate && (
                  <span>{item.releaseDate.split("-")[0]}</span>
                )}
              </div>
            </div>

            {/* Action */}
            <div className="ml-4">
              {navigatingID === item.id ? (
                <div className="w-8 h-8 rounded-full border-2 border-emerald-500/30 border-t-emerald-500 animate-spin" />
              ) : (
                <div className="w-10 h-10 rounded-full bg-slate-100 dark:bg-slate-700 md:group-hover:bg-emerald-500 md:group-hover:text-white flex items-center justify-center transition-colors text-slate-400">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="m9 18 6-6-6-6" />
                  </svg>
                </div>
              )}
            </div>
          </div>
        ))}
        {results && results.length === 0 && query && (
          <div className="text-center py-12">
            <p className="text-slate-500 text-lg">
              No results found for "{query}"
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default SearchForm;
