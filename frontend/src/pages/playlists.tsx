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
        <h1 className="text-3xl font-semibold mb-4">Playlists</h1>
        <p className="text-lg text-slate-600 dark:text-slate-300">
          Browse user playlists and test the playlist system
        </p>
      </div>

      {loading && !data && (
        <div className="flex justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
        </div>
      )}

      {error && (
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Error Loading Playlists</h3>
          <p className="text-rose-500">{error.message}</p>
          <p className="text-sm text-rose-400 mt-2">
            {error.message.includes('relation "playlists" does not exist') ||
            error.message.includes("SQLSTATE 42P01")
              ? "Database migrations are not applied yet. Run the backend migrations and retry."
              : "This might be expected if no playlists exist in the database yet."}
          </p>
        </div>
      )}

      {data?.playlists && (
        <div className="space-y-6">
          <div className="card">
            <p className="text-slate-600 dark:text-slate-300">
              Total Playlists: <strong>{data.playlists.totalCount}</strong>
            </p>
          </div>

          {data.playlists.edges.length > 0 ? (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.playlists.edges.map(({ node: playlist }) => (
                <div key={playlist.id} className="card">
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
                    <p className="text-slate-600 dark:text-slate-300 mb-3">
                      {playlist.description}
                    </p>
                  )}
                  <div className="text-sm text-slate-500 dark:text-slate-400 border-t border-slate-200/70 dark:border-slate-800/60 pt-3">
                    <p>
                      Created by <strong>{playlist.creator.name}</strong>
                    </p>
                    <p>{new Date(playlist.createdAt).toLocaleDateString()}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="card text-center">
              <p className="text-amber-600 dark:text-amber-300">
                No playlists found in the database.
              </p>
              <p className="text-sm text-amber-500 mt-2">
                Playlists will appear here once users start creating them.
              </p>
            </div>
          )}

          {data.playlists.pageInfo.hasNextPage && (
            <div className="text-center">
              <button
                onClick={loadMore}
                disabled={loading}
                className="btn-primary disabled:opacity-50"
              >
                {loading ? "Loading..." : "Load More"}
              </button>
            </div>
          )}
        </div>
      )}

      <div className="card-muted">
        <h3 className="text-lg font-semibold mb-3">Testing Notes</h3>
        <ul className="list-disc list-inside space-y-2 text-slate-600 dark:text-slate-300">
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
