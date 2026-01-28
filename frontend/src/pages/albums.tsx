import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_ALBUMS } from "../lib/graphql/queries";

export default function AlbumsPage() {
  const { data, loading, error, fetchMore } = useQuery(GET_ALBUMS, {
    variables: { first: 10 },
  });

  const loadMore = () => {
    if (data?.albums.pageInfo.hasNextPage) {
      fetchMore({
        variables: {
          first: 10,
          after: data.albums.pageInfo.endCursor,
        },
      });
    }
  };

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-semibold mb-2">Albums</h1>
        <p className="text-slate-500 dark:text-slate-400">
          Browse and discover music
        </p>
      </div>

      {loading && !data && (
        <div className="flex justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
        </div>
      )}

      {error && (
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Error Loading Albums</h3>
          <p className="text-rose-500">{error.message}</p>
          <p className="text-sm text-rose-400 mt-2">
            {error.message.includes('relation "albums" does not exist') ||
            error.message.includes("SQLSTATE 42P01")
              ? "Database migrations are not applied yet. Run the backend migrations and retry."
              : "This might be expected if no albums are stored in the database yet."}
          </p>
        </div>
      )}

      {data?.albums && (
        <div className="space-y-6">
          <div className="card">
            <p className="text-slate-600 dark:text-slate-300">
              Total Albums: <strong>{data.albums.totalCount}</strong>
            </p>
          </div>

          {data.albums.edges.length > 0 ? (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.albums.edges.map(({ node: album }: { node: any }) => (
                <div key={album.id} className="card">
                  {album.coverImage && (
                    <img
                      src={album.coverImage}
                      alt={album.title}
                      className="w-full h-48 object-cover rounded mb-4"
                    />
                  )}
                  <h3 className="text-lg font-semibold mb-2">
                    <Link
                      href={`/albums/${album.id}`}
                      className="text-emerald-600 hover:text-emerald-500"
                    >
                      {album.title}
                    </Link>
                  </h3>
                  <p className="text-slate-600 dark:text-slate-300 mb-2">
                    by {album.artist.name}
                  </p>
                  {album.releaseDate && (
                    <p className="text-sm text-slate-500 dark:text-slate-400">
                      Released:{" "}
                      {new Date(album.releaseDate).toLocaleDateString()}
                    </p>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <div className="card text-center">
              <p className="text-amber-600 dark:text-amber-300">
                No albums found in the database.
              </p>
              <p className="text-sm text-amber-500 mt-2">
                Try using the search functionality to find albums from Spotify
                first.
              </p>
            </div>
          )}

          {data.albums.pageInfo.hasNextPage && (
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
