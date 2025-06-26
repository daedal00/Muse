import { useState, useEffect } from "react";
import { useQuery, useMutation } from "@apollo/client";
import { GET_ME, GET_SPOTIFY_PLAYLISTS } from "../lib/graphql/queries";
import { IMPORT_SPOTIFY_PLAYLIST } from "../lib/graphql/mutations";

interface SpotifyPlaylist {
  id: string;
  name: string;
  description?: string;
  image?: string;
  trackCount: number;
  isPublic: boolean;
}

interface PlaylistCardProps {
  playlist: SpotifyPlaylist;
  onImport: (playlistId: string) => void;
  isImporting: boolean;
}

const PlaylistCard: React.FC<PlaylistCardProps> = ({
  playlist,
  onImport,
  isImporting,
}) => {
  return (
    <div className="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow p-4">
      <div className="flex gap-4">
        {playlist.image ? (
          <img
            src={playlist.image}
            alt={playlist.name}
            className="w-16 h-16 rounded-lg object-cover flex-shrink-0"
          />
        ) : (
          <div className="w-16 h-16 bg-gray-200 rounded-lg flex items-center justify-center flex-shrink-0">
            <svg
              className="w-8 h-8 text-gray-400"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M9.383 3.076A1 1 0 0110 4v12a1 1 0 01-1.707.707L4.586 13H2a1 1 0 01-1-1V8a1 1 0 011-1h2.586l3.707-3.707a1 1 0 011.09-.217zM15.657 6.343a1 1 0 011.414 0A9.972 9.972 0 0119 12a9.972 9.972 0 01-1.929 5.657 1 1 0 11-1.414-1.414A7.971 7.971 0 0017 12c0-1.594-.471-3.078-1.343-4.243a1 1 0 010-1.414z"
                clipRule="evenodd"
              />
            </svg>
          </div>
        )}

        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-gray-900 truncate">
            {playlist.name}
          </h3>
          {playlist.description && (
            <p className="text-sm text-gray-600 mt-1 line-clamp-2">
              {playlist.description}
            </p>
          )}
          <div className="flex items-center gap-4 mt-2 text-sm text-gray-500">
            <span>{playlist.trackCount} tracks</span>
            <span className="flex items-center gap-1">
              {playlist.isPublic ? (
                <>
                  <svg
                    className="w-4 h-4"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                    <path
                      fillRule="evenodd"
                      d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                  Public
                </>
              ) : (
                <>
                  <svg
                    className="w-4 h-4"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                  Private
                </>
              )}
            </span>
          </div>
        </div>

        <div className="flex-shrink-0">
          <button
            onClick={() => onImport(playlist.id)}
            disabled={isImporting}
            className="px-4 py-2 bg-green-500 hover:bg-green-600 disabled:bg-gray-400 text-white text-sm font-medium rounded-lg transition-colors"
          >
            {isImporting ? "Importing..." : "Import"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default function SpotifyPlaylistsPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const [importingPlaylist, setImportingPlaylist] = useState<string | null>(
    null
  );
  const playlistsPerPage = 10;

  const { data: userData, loading: userLoading } = useQuery(GET_ME);

  const { data, loading, error, refetch } = useQuery(GET_SPOTIFY_PLAYLISTS, {
    variables: {
      input: {
        limit: playlistsPerPage,
        offset: (currentPage - 1) * playlistsPerPage,
      },
    },
    skip: !userData?.me, // Only fetch if user is authenticated
  });

  const [importPlaylist] = useMutation(IMPORT_SPOTIFY_PLAYLIST, {
    onCompleted: (data) => {
      setImportingPlaylist(null);
      alert(
        `Successfully imported playlist: ${data.importSpotifyPlaylist.title}`
      );
    },
    onError: (error) => {
      setImportingPlaylist(null);
      console.error("Import failed:", error);
      alert(`Failed to import playlist: ${error.message}`);
    },
  });

  const handleImportPlaylist = async (spotifyPlaylistId: string) => {
    if (!userData?.me) {
      alert("You must be logged in to import playlists");
      return;
    }

    setImportingPlaylist(spotifyPlaylistId);
    try {
      await importPlaylist({
        variables: {
          spotifyPlaylistId,
        },
      });
    } catch (error) {
      console.error("Import error:", error);
    }
  };

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  const handleRefresh = () => {
    setCurrentPage(1);
    refetch({
      input: {
        limit: playlistsPerPage,
        offset: 0,
      },
    });
  };

  // Loading state
  if (userLoading) {
    return (
      <div className="flex justify-center items-center min-h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  // Not authenticated
  if (!userData?.me) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">
          Authentication Required
        </h2>
        <p className="text-gray-600 mb-6">
          You need to be logged in to view your Spotify playlists.
        </p>
        <a
          href="/auth"
          className="inline-block bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700"
        >
          Login / Register
        </a>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="text-center py-12">
        <h2 className="text-xl font-semibold text-gray-900 mb-2">
          Error Loading Playlists
        </h2>
        <p className="text-red-600 mb-4">{error.message}</p>
        <button
          onClick={handleRefresh}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Try Again
        </button>
      </div>
    );
  }

  const playlists = data?.spotifyPlaylists?.playlists || [];
  const totalCount = data?.spotifyPlaylists?.totalCount || 0;
  const hasNextPage = data?.spotifyPlaylists?.hasNextPage || false;
  const totalPages = Math.ceil(totalCount / playlistsPerPage);

  // Filter playlists based on search query
  const filteredPlaylists = playlists.filter((playlist: SpotifyPlaylist) =>
    playlist.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              Your Spotify Playlists
            </h1>
            <p className="text-gray-600">
              Import your Spotify playlists to Muse
            </p>
          </div>
          <button
            onClick={handleRefresh}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
          >
            <svg
              className={`w-4 h-4 ${loading ? "animate-spin" : ""}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
            Refresh
          </button>
        </div>

        {/* Search Bar */}
        <div className="relative mb-6">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg
              className="h-5 w-5 text-gray-400"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                clipRule="evenodd"
              />
            </svg>
          </div>
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search your playlists..."
            className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Playlist Count */}
        <div className="flex items-center justify-between text-sm text-gray-600 mb-4">
          <span>
            {filteredPlaylists.length} of {totalCount} playlists
            {searchQuery && ` (filtered)`}
          </span>
        </div>
      </div>

      {/* Loading State */}
      {loading && (
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500"></div>
          <span className="ml-2 text-gray-600">Loading playlists...</span>
        </div>
      )}

      {/* Content */}
      <div className="space-y-6">
        {!loading && filteredPlaylists.length === 0 && !searchQuery && (
          <div className="text-center py-12">
            <div className="mb-4">
              <svg
                className="w-16 h-16 text-gray-400 mx-auto"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No Spotify playlists found
            </h3>
            <p className="text-gray-600">
              Make sure you have connected your Spotify account and have created
              some playlists.
            </p>
          </div>
        )}

        {!loading && filteredPlaylists.length === 0 && searchQuery && (
          <div className="text-center py-12">
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              No playlists match your search
            </h3>
            <p className="text-gray-600">
              Try searching with different keywords.
            </p>
          </div>
        )}

        {filteredPlaylists.length > 0 && (
          <>
            {/* Playlists Grid */}
            <div className="space-y-4 mb-8">
              {filteredPlaylists.map((playlist: SpotifyPlaylist) => (
                <PlaylistCard
                  key={playlist.id}
                  playlist={playlist}
                  onImport={handleImportPlaylist}
                  isImporting={importingPlaylist === playlist.id}
                />
              ))}
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-between">
                <div className="text-sm text-gray-700">
                  Showing {(currentPage - 1) * playlistsPerPage + 1} to{" "}
                  {Math.min(currentPage * playlistsPerPage, totalCount)} of{" "}
                  {totalCount} playlists
                </div>

                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 1}
                    className="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Previous
                  </button>

                  <div className="flex items-center gap-1">
                    {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                      const pageNum = i + 1;
                      const isCurrentPage = pageNum === currentPage;

                      return (
                        <button
                          key={pageNum}
                          onClick={() => handlePageChange(pageNum)}
                          className={`px-3 py-2 text-sm font-medium rounded-md ${
                            isCurrentPage
                              ? "bg-blue-500 text-white"
                              : "text-gray-700 bg-white border border-gray-300 hover:bg-gray-50"
                          }`}
                        >
                          {pageNum}
                        </button>
                      );
                    })}

                    {totalPages > 5 && (
                      <>
                        <span className="px-2 text-gray-500">...</span>
                        <button
                          onClick={() => handlePageChange(totalPages)}
                          className={`px-3 py-2 text-sm font-medium rounded-md ${
                            totalPages === currentPage
                              ? "bg-blue-500 text-white"
                              : "text-gray-700 bg-white border border-gray-300 hover:bg-gray-50"
                          }`}
                        >
                          {totalPages}
                        </button>
                      </>
                    )}
                  </div>

                  <button
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={!hasNextPage}
                    className="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
