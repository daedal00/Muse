import React from "react";
import { useQuery } from "@apollo/client";
import { GET_PLAYLISTS } from "../lib/graphql/queries";

export default function PlaylistsPage() {
  const { data, loading, error, fetchMore } = useQuery(GET_PLAYLISTS, {
    variables: { first: 10 },
  });

  const loadMore = () => {
    if (data?.playlists.pageInfo.hasNextPage) {
      fetchMore({
        variables: {
          first: 10,
          after: data.playlists.pageInfo.endCursor,
        },
      });
    }
  };

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Playlists</h1>
        <p className="text-lg text-gray-600">
          Browse user playlists and test the playlist system
        </p>
      </div>

      {loading && !data && (
        <div className="flex justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      )}

      {error && (
        <div className="bg-red-50 border border-red-200 p-6 rounded-lg">
          <h3 className="text-lg font-semibold text-red-800 mb-2">
            Error Loading Playlists
          </h3>
          <p className="text-red-600">{error.message}</p>
          <p className="text-sm text-red-500 mt-2">
            This might be expected if no playlists exist in the database yet.
          </p>
        </div>
      )}

      {data?.playlists && (
        <div className="space-y-6">
          <div className="bg-white p-4 rounded-lg shadow-md">
            <p className="text-gray-600">
              Total Playlists: <strong>{data.playlists.totalCount}</strong>
            </p>
          </div>

          {data.playlists.edges.length > 0 ? (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.playlists.edges.map(({ node: playlist }) => (
                <div
                  key={playlist.id}
                  className="bg-white p-6 rounded-lg shadow-md"
                >
                  {playlist.coverImage && (
                    <img
                      src={playlist.coverImage}
                      alt={playlist.title}
                      className="w-full h-48 object-cover rounded mb-4"
                    />
                  )}
                  <h3 className="text-lg font-semibold mb-2">
                    {playlist.title}
                  </h3>
                  {playlist.description && (
                    <p className="text-gray-600 mb-3">{playlist.description}</p>
                  )}
                  <div className="text-sm text-gray-500 border-t pt-3">
                    <p>
                      Created by <strong>{playlist.creator.name}</strong>
                    </p>
                    <p>{new Date(playlist.createdAt).toLocaleDateString()}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
              <p className="text-yellow-800">
                No playlists found in the database.
              </p>
              <p className="text-sm text-yellow-600 mt-2">
                Playlists will appear here once users start creating them.
              </p>
            </div>
          )}

          {data.playlists.pageInfo.hasNextPage && (
            <div className="text-center">
              <button
                onClick={loadMore}
                disabled={loading}
                className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {loading ? "Loading..." : "Load More"}
              </button>
            </div>
          )}
        </div>
      )}

      <div className="bg-blue-50 border border-blue-200 p-6 rounded-lg">
        <h3 className="text-lg font-semibold mb-3">🧪 Testing Notes</h3>
        <ul className="list-disc list-inside space-y-2 text-gray-700">
          <li>This page tests the GraphQL playlists query with pagination</li>
          <li>Playlists are created by authenticated users</li>
          <li>Each playlist can have a title, description, and cover image</li>
          <li>
            To test: create a user, login, then use the createPlaylist mutation
          </li>
        </ul>
      </div>
    </div>
  );
}
