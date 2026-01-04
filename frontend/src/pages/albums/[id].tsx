import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { useQuery } from "@apollo/client";
import { GET_ALBUM } from "../../lib/graphql/queries";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

const formatRating = (rating: number) => {
  const rounded = Math.round(rating);
  return "*".repeat(rounded) + "-".repeat(5 - rounded);
};

export default function AlbumDetailPage() {
  const router = useRouter();
  const albumID = typeof router.query.id === "string" ? router.query.id : "";

  const { data, loading, error } = useQuery(GET_ALBUM, {
    variables: { id: albumID, tracksFirst: 50, reviewsFirst: 20 },
    skip: !albumID,
  });

  if (loading) {
    return (
      <div className="flex justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 p-6 rounded-lg">
        <h3 className="text-lg font-semibold text-red-800 mb-2">
          Error Loading Album
        </h3>
        <p className="text-red-600">{error.message}</p>
      </div>
    );
  }

  if (!data?.album) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
        <p className="text-yellow-800">Album not found.</p>
        <Link href="/albums" className="text-blue-600 hover:text-blue-800">
          Back to Albums
        </Link>
      </div>
    );
  }

  const album = data.album;
  const averageRating =
    typeof album.averageRating === "number" ? album.averageRating : null;

  return (
    <div className="space-y-8">
      <Link href="/albums" className="text-blue-600 hover:text-blue-800">
        {"<- Back to Albums"}
      </Link>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <div className="flex flex-col md:flex-row md:items-start md:space-x-6">
          {album.coverImage && (
            <img
              src={album.coverImage}
              alt={album.title}
              className="w-48 h-48 object-cover rounded mb-4 md:mb-0"
            />
          )}
          <div className="space-y-2">
            <h1 className="text-3xl font-bold text-gray-900">{album.title}</h1>
            <p className="text-gray-600">by {album.artist.name}</p>
            {album.releaseDate && (
              <p className="text-sm text-gray-500">
                Released: {new Date(album.releaseDate).toLocaleDateString()}
              </p>
            )}
            <div className="text-sm text-gray-500">
              Tracks: {album.tracks.totalCount}
            </div>
            <div className="text-sm text-gray-500">
              Reviews: {album.reviews.totalCount}
            </div>
            <div className="pt-2">
              {averageRating !== null ? (
                <div className="text-lg">
                  {formatRating(averageRating)}{" "}
                  <span className="text-sm text-gray-600">
                    ({averageRating.toFixed(1)}/5)
                  </span>
                </div>
              ) : (
                <p className="text-sm text-gray-500">No ratings yet.</p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-2xl font-semibold mb-4">Tracks</h2>
        {album.tracks.edges.length > 0 ? (
          <ol className="space-y-2">
            {album.tracks.edges.map(({ node }: any, index: number) => (
              <li key={node.id} className="flex items-center justify-between">
                <div>
                  <span className="text-gray-500 mr-3">
                    {node.trackNumber || index + 1}.
                  </span>
                  <Link
                    href={`/tracks/${node.id}`}
                    className="text-blue-600 hover:text-blue-800"
                  >
                    {node.title}
                  </Link>
                </div>
                {formatDuration(node.duration) && (
                  <span className="text-sm text-gray-500">
                    {formatDuration(node.duration)}
                  </span>
                )}
              </li>
            ))}
          </ol>
        ) : (
          <p className="text-gray-500">No tracks stored for this album yet.</p>
        )}
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-2xl font-semibold mb-4">Reviews</h2>
        {album.reviews.edges.length > 0 ? (
          <div className="space-y-4">
            {album.reviews.edges.map(({ node }: any) => (
              <div key={node.id} className="border-b pb-4 last:border-b-0">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-semibold">{node.user.name}</p>
                    <p className="text-sm text-gray-500">
                      {new Date(node.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="text-right">
                    <div>{formatRating(node.rating)}</div>
                    <div className="text-sm text-gray-500">
                      {node.rating}/5
                    </div>
                  </div>
                </div>
                {node.reviewText && (
                  <p className="text-gray-700 mt-2 italic">
                    "{node.reviewText}"
                  </p>
                )}
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500">No reviews yet for this album.</p>
        )}
      </div>
    </div>
  );
}
