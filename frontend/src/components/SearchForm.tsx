import React, { useState } from "react";
import { useLazyQuery, useMutation } from "@apollo/client";
import { useRouter } from "next/router";
import {
  SEARCH_ALBUMS,
  SEARCH_ARTISTS,
  SEARCH_TRACKS,
} from "../lib/graphql/queries";
import { IMPORT_ALBUM, IMPORT_TRACK } from "../lib/graphql/mutations";

interface SearchFormProps {
  onResults?: (results: any, type: "albums" | "artists" | "tracks") => void;
}

const SearchForm: React.FC<SearchFormProps> = ({ onResults }) => {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [searchType, setSearchType] = useState<
    "albums" | "artists" | "tracks"
  >("albums");
  const [importingID, setImportingID] = useState<string | null>(null);
  const [importError, setImportError] = useState<string | null>(null);

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

  const [importAlbum] = useMutation(IMPORT_ALBUM);
  const [importTrack] = useMutation(IMPORT_TRACK);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    setImportError(null);

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

  const handleImport = async (type: "albums" | "tracks", spotifyID: string) => {
    setImportError(null);
    setImportingID(spotifyID);
    try {
      if (type === "albums") {
        const { data } = await importAlbum({
          variables: { spotifyAlbumID: spotifyID },
        });
        if (data?.importAlbum?.id) {
          await router.push(`/albums/${data.importAlbum.id}`);
        }
      } else {
        const { data } = await importTrack({
          variables: { spotifyTrackID: spotifyID },
        });
        if (data?.importTrack?.id) {
          await router.push(`/tracks/${data.importTrack.id}`);
        }
      }
    } catch (err) {
      setImportError(
        err instanceof Error ? err.message : "Failed to import item"
      );
    } finally {
      setImportingID(null);
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
      <h2 className="text-xl font-semibold mb-4">Search Music</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex space-x-4">
          <div className="flex-1">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search for albums, artists, or tracks..."
              className="input"
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
            {loading ? "Searching..." : "Search"}
          </button>
        </div>
      </form>

      {error && (
        <div className="mt-4 rounded-xl border border-rose-200 bg-rose-50 p-3 text-rose-700 dark:border-rose-500/30 dark:bg-rose-500/10 dark:text-rose-200">
          Error: {error.message}
        </div>
      )}

      {importError && (
        <div className="mt-4 rounded-xl border border-rose-200 bg-rose-50 p-3 text-rose-700 dark:border-rose-500/30 dark:bg-rose-500/10 dark:text-rose-200">
          Import Error: {importError}
        </div>
      )}

      {results && results.length > 0 && (
        <div className="mt-6">
          <h3 className="text-lg font-semibold mb-3">
            Search Results ({results.length})
          </h3>
          <div className="grid gap-4">
            {results.map((item: any) => (
              <div key={item.id} className="card-muted">
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
                      <p className="text-slate-600 dark:text-slate-300">
                        by {item.artist.map((a: any) => a.name).join(", ")}
                      </p>
                      {item.releaseDate && (
                        <p className="text-sm text-slate-500 dark:text-slate-400">
                          Released:{" "}
                          {new Date(item.releaseDate).toLocaleDateString()}
                        </p>
                      )}
                      <p className="text-xs text-emerald-600">
                        Source: {item.externalSource}
                      </p>
                    </div>
                    <div className="ml-auto">
                      <button
                        type="button"
                        onClick={() => handleImport("albums", item.id)}
                        disabled={importingID === item.id}
                        className="btn-secondary disabled:opacity-50"
                      >
                        {importingID === item.id
                          ? "Importing..."
                          : "Import & View"}
                      </button>
                    </div>
                  </div>
                ) : searchType === "tracks" ? (
                  <div className="flex items-start space-x-4">
                    {item.album?.coverImage && (
                      <img
                        src={item.album.coverImage}
                        alt={item.album.title}
                        className="w-16 h-16 object-cover rounded"
                      />
                    )}
                    <div className="flex-1">
                      <h4 className="font-semibold">{item.title}</h4>
                      <p className="text-slate-600 dark:text-slate-300">
                        by {item.artists.map((a: any) => a.name).join(", ")}
                      </p>
                      {item.album?.title && (
                        <p className="text-sm text-slate-500 dark:text-slate-400">
                          Album: {item.album.title}
                        </p>
                      )}
                      {formatDuration(item.duration) && (
                        <p className="text-sm text-slate-500 dark:text-slate-400">
                          Duration: {formatDuration(item.duration)}
                        </p>
                      )}
                      {item.trackNumber && (
                        <p className="text-sm text-slate-500 dark:text-slate-400">
                          Track {item.trackNumber}
                        </p>
                      )}
                      <p className="text-xs text-emerald-600">
                        Source: {item.externalSource}
                      </p>
                    </div>
                    <div>
                      <button
                        type="button"
                        onClick={() => handleImport("tracks", item.id)}
                        disabled={importingID === item.id}
                        className="btn-secondary disabled:opacity-50"
                      >
                        {importingID === item.id
                          ? "Importing..."
                          : "Import & View"}
                      </button>
                    </div>
                  </div>
                ) : (
                  <div>
                    <h4 className="font-semibold">{item.name}</h4>
                    <p className="text-xs text-emerald-600">
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
        <div className="mt-4 rounded-xl border border-amber-200 bg-amber-50 p-3 text-amber-700 dark:border-amber-400/30 dark:bg-amber-500/10 dark:text-amber-200">
          No results found for "{query}"
        </div>
      )}
    </div>
  );
};

export default SearchForm;
