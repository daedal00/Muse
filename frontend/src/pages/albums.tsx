import React from "react";
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
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Albums</h1>
        <p className="text-lg text-gray-600">
          Browse stored albums and test pagination
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
            Error Loading Albums
          </h3>
          <p className="text-red-600">{error.message}</p>
          <p className="text-sm text-red-500 mt-2">
            This might be expected if no albums are stored in the database yet.
          </p>
        </div>
      )}

      {data?.albums && (
        <div className="space-y-6">
          <div className="bg-white p-4 rounded-lg shadow-md">
            <p className="text-gray-600">
              Total Albums: <strong>{data.albums.totalCount}</strong>
            </p>
          </div>

          {data.albums.edges.length > 0 ? (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {data.albums.edges.map(({ node: album }) => (
                <div
                  key={album.id}
                  className="bg-white p-6 rounded-lg shadow-md"
                >
                  {album.coverImage && (
                    <img
                      src={album.coverImage}
                      alt={album.title}
                      className="w-full h-48 object-cover rounded mb-4"
                    />
                  )}
                  <h3 className="text-lg font-semibold mb-2">{album.title}</h3>
                  <p className="text-gray-600 mb-2">by {album.artist.name}</p>
                  {album.releaseDate && (
                    <p className="text-sm text-gray-500">
                      Released:{" "}
                      {new Date(album.releaseDate).toLocaleDateString()}
                    </p>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
              <p className="text-yellow-800">
                No albums found in the database.
              </p>
              <p className="text-sm text-yellow-600 mt-2">
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
          <li>This page tests the GraphQL albums query with pagination</li>
          <li>Albums are stored in your database after being searched/added</li>
          <li>Pagination uses cursor-based pagination (GraphQL connections)</li>
          <li>If no albums appear, the database might be empty initially</li>
        </ul>
      </div>
    </div>
  );
}
