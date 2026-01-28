import React from "react";
import Link from "next/link";
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
      <div>
        <h1 className="text-3xl font-semibold mb-2">Playlists</h1>
        <p className="text-slate-500 dark:text-slate-400">
          Curated collections from the community
        </p>
      </div>

      {loading && !data && (
        <div className="flex justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
        </div>
      )}

      {error && (
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">
            Error Loading Playlists
          </h3>
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
              {data.playlists.edges.map(({ node: playlist }: { node: any }) => (
                <Link
                  key={playlist.id}
                  href={`/playlists/${playlist.id}`}
                  className="card hover:ring-2 hover:ring-emerald-500/50 transition-all group"
                >
                  {playlist.coverImage ? (
                    <img
                      src={playlist.coverImage}
                      alt={playlist.title}
                      className="w-full h-48 object-cover rounded-xl mb-4 group-hover:opacity-90 transition-opacity"
                    />
                  ) : (
                    <div className="w-full h-48 rounded-xl bg-gradient-to-br from-emerald-400 to-emerald-600 mb-4 flex items-center justify-center">
                      <span className="text-5xl">🎵</span>
                    </div>
                  )}
                  <h3 className="text-lg font-semibold mb-1 group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">
                    {playlist.title}
                  </h3>
                  {playlist.description && (
                    <p className="text-sm text-slate-600 dark:text-slate-300 mb-3 line-clamp-2">
                      {playlist.description}
                    </p>
                  )}
                  <div className="text-sm text-slate-500 dark:text-slate-400">
                    <p>
                      by{" "}
                      <span className="font-medium">
                        {playlist.creator.name}
                      </span>
                    </p>
                  </div>
                </Link>
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
    </div>
  );
}
