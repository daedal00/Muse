import React, { useState } from "react";
import { useQuery, useMutation } from "@apollo/client";
import {
  GET_SPOTIFY_AUTH_URL,
  GET_SPOTIFY_PLAYLISTS,
} from "../lib/graphql/queries";
import { IMPORT_SPOTIFY_PLAYLIST } from "../lib/graphql/mutations";

interface SpotifyPlaylist {
  id: string;
  name: string;
  description?: string;
  trackCount: number;
  isPublic: boolean;
  image?: string;
}

const SpotifyIntegration: React.FC = () => {
  const [selectedPlaylists, setSelectedPlaylists] = useState<Set<string>>(
    new Set()
  );
  const [importingPlaylists, setImportingPlaylists] = useState<Set<string>>(
    new Set()
  );

  const { data: authData, loading: authLoading } =
    useQuery(GET_SPOTIFY_AUTH_URL);
  const {
    data: playlistsData,
    loading: playlistsLoading,
    error: playlistsError,
    refetch,
  } = useQuery(GET_SPOTIFY_PLAYLISTS, {
    errorPolicy: "all",
  });

  const [importPlaylist] = useMutation(IMPORT_SPOTIFY_PLAYLIST, {
    onCompleted: (data) => {
      console.log(
        "Playlist imported successfully:",
        data.importSpotifyPlaylist
      );
    },
    onError: (error) => {
      console.error("Error importing playlist:", error);
    },
  });

  const handleSpotifyLogin = () => {
    if (authData?.spotifyAuthURL?.url) {
      window.location.href = authData.spotifyAuthURL.url;
    }
  };

  const handlePlaylistToggle = (playlistId: string) => {
    const newSelected = new Set(selectedPlaylists);
    if (newSelected.has(playlistId)) {
      newSelected.delete(playlistId);
    } else {
      newSelected.add(playlistId);
    }
    setSelectedPlaylists(newSelected);
  };

  const handleImportSelected = async () => {
    const playlistsToImport = Array.from(selectedPlaylists);
    const importing = new Set(playlistsToImport);
    setImportingPlaylists(importing);

    try {
      await Promise.all(
        playlistsToImport.map(async (playlistId) => {
          await importPlaylist({
            variables: { spotifyPlaylistId: playlistId },
          });
          setImportingPlaylists((prev) => {
            const newSet = new Set(prev);
            newSet.delete(playlistId);
            return newSet;
          });
        })
      );

      setSelectedPlaylists(new Set());
      alert("Selected playlists imported successfully!");
    } catch (error) {
      console.error("Error importing playlists:", error);
      alert("Error importing some playlists. Check console for details.");
    } finally {
      setImportingPlaylists(new Set());
    }
  };

  const handleImportSingle = async (playlistId: string) => {
    setImportingPlaylists((prev) => new Set(prev).add(playlistId));

    try {
      await importPlaylist({
        variables: { spotifyPlaylistId: playlistId },
      });
      alert("Playlist imported successfully!");
    } catch (error) {
      console.error("Error importing playlist:", error);
      alert("Error importing playlist. Check console for details.");
    } finally {
      setImportingPlaylists((prev) => {
        const newSet = new Set(prev);
        newSet.delete(playlistId);
        return newSet;
      });
    }
  };

  const playlists: SpotifyPlaylist[] = playlistsData?.spotifyPlaylists || [];
  const hasSpotifyConnection = !playlistsError || playlists.length > 0;

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <div className="flex items-center mb-6">
        <div className="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center mr-3">
          <svg
            className="w-5 h-5 text-white"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.48.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.42 1.56-.299.421-1.02.599-1.559.3z" />
          </svg>
        </div>
        <h2 className="text-2xl font-bold">Spotify Integration</h2>
      </div>

      {!hasSpotifyConnection ? (
        <div className="text-center py-8">
          <p className="text-gray-600 mb-4">
            Connect your Spotify account to import your playlists
          </p>
          <button
            onClick={handleSpotifyLogin}
            disabled={authLoading}
            className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center mx-auto"
          >
            {authLoading ? (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
            ) : (
              <svg
                className="w-5 h-5 mr-2"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.48.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.42 1.56-.299.421-1.02.599-1.559.3z" />
              </svg>
            )}
            Connect to Spotify
          </button>
        </div>
      ) : (
        <div>
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-semibold">Your Spotify Playlists</h3>
            <div className="space-x-2">
              <button
                onClick={() => refetch()}
                disabled={playlistsLoading}
                className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600 disabled:opacity-50"
              >
                {playlistsLoading ? "Refreshing..." : "Refresh"}
              </button>
              {selectedPlaylists.size > 0 && (
                <button
                  onClick={handleImportSelected}
                  disabled={importingPlaylists.size > 0}
                  className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  Import Selected ({selectedPlaylists.size})
                </button>
              )}
            </div>
          </div>

          {playlistsLoading && !playlists.length ? (
            <div className="flex justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
            </div>
          ) : playlists.length > 0 ? (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              {playlists.map((playlist) => (
                <div
                  key={playlist.id}
                  className="border rounded-lg p-4 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-start space-x-3">
                    <input
                      type="checkbox"
                      checked={selectedPlaylists.has(playlist.id)}
                      onChange={() => handlePlaylistToggle(playlist.id)}
                      className="mt-1"
                    />
                    <div className="flex-1">
                      {playlist.image && (
                        <img
                          src={playlist.image}
                          alt={playlist.name}
                          className="w-full h-32 object-cover rounded mb-2"
                        />
                      )}
                      <h4 className="font-semibold text-sm mb-1">
                        {playlist.name}
                      </h4>
                      {playlist.description && (
                        <p className="text-xs text-gray-600 mb-2 line-clamp-2">
                          {playlist.description}
                        </p>
                      )}
                      <div className="flex justify-between items-center text-xs text-gray-500 mb-2">
                        <span>{playlist.trackCount} tracks</span>
                        <span>{playlist.isPublic ? "Public" : "Private"}</span>
                      </div>
                      <button
                        onClick={() => handleImportSingle(playlist.id)}
                        disabled={importingPlaylists.has(playlist.id)}
                        className="w-full bg-green-600 text-white py-1 px-2 rounded text-xs hover:bg-green-700 disabled:opacity-50"
                      >
                        {importingPlaylists.has(playlist.id)
                          ? "Importing..."
                          : "Import"}
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <p>No playlists found in your Spotify account.</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default SpotifyIntegration;
