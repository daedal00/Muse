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
  const [searchType, setSearchType] = useState<
    "albums" | "artists" | "tracks"
  >("albums");
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
    <div className="card">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex gap-3">
          <div className="flex-1">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search albums, artists, or tracks..."
              className="input text-lg"
            />
          </div>

          <select
            value={searchType}
            onChange={(e) =>
              setSearchType(
                e.target.value as "albums" | "artists" | "tracks"
              )
            }
            className="select"
          >
            <option value="albums">Albums</option>
            <option value="artists">Artists</option>
            <option value="tracks">Tracks</option>
          </select>

          <button
            type="submit"
            disabled={loading || !query.trim()}
            className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "..." : "Search"}
          </button>
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
                      <p className="font-medium text-base truncate">{item.title}</p>
                      <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                        {item.artist.map((a: any) => a.name).join(", ")}
                        {item.releaseDate && ` · ${new Date(item.releaseDate).getFullYear()}`}
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
                      <p className="font-medium text-base truncate">{item.title}</p>
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
                      <p className="font-medium text-base truncate">{item.name}</p>
                      <p className="text-sm text-slate-500 dark:text-slate-400">Artist</p>
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {results && results.length === 0 && (
        <div className="mt-4 rounded-xl border border-amber-200 bg-amber-50 p-3 text-amber-700 dark:border-amber-400/30 dark:bg-amber-500/10 dark:text-amber-200">
          No results found for "{query}"
        </div>
      )}
    </div>
  );
};

export default SearchForm;
