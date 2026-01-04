import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { useQuery } from "@apollo/client";
import { GET_ALBUM } from "../../lib/graphql/queries";
import StarRating from "../../components/StarRating";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
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
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error Loading Album</h3>
        <p className="text-rose-500">{error.message}</p>
      </div>
    );
  }

  if (!data?.album) {
    return (
      <div className="card text-center">
        <p className="text-amber-600 dark:text-amber-300">Album not found.</p>
        <Link href="/albums" className="text-emerald-600 hover:text-emerald-500">
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
      <Link href="/albums" className="text-emerald-600 hover:text-emerald-500">
        {"<- Back to Albums"}
      </Link>

      <div className="card">
        <div className="flex flex-col md:flex-row md:items-start md:space-x-6">
          {album.coverImage && (
            <img
              src={album.coverImage}
              alt={album.title}
              className="w-48 h-48 object-cover rounded mb-4 md:mb-0"
            />
          )}
          <div className="space-y-2">
            <h1 className="text-3xl font-semibold">{album.title}</h1>
            <p className="text-slate-600 dark:text-slate-300">
              by {album.artist.name}
            </p>
            {album.releaseDate && (
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Released: {new Date(album.releaseDate).toLocaleDateString()}
              </p>
            )}
            <div className="text-sm text-slate-500 dark:text-slate-400">
              Tracks: {album.tracks.totalCount}
            </div>
            <div className="text-sm text-slate-500 dark:text-slate-400">
              Reviews: {album.reviews.totalCount}
            </div>
            <div className="pt-2">
              {averageRating !== null ? (
                <StarRating
                  value={averageRating}
                  label={`${averageRating.toFixed(1)}/5`}
                />
              ) : (
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  No ratings yet.
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <h2 className="text-2xl font-semibold mb-4">Tracks</h2>
        {album.tracks.edges.length > 0 ? (
          <ol className="space-y-2">
            {album.tracks.edges.map(({ node }: any, index: number) => (
              <li key={node.id} className="flex items-center justify-between">
                <div>
                  <span className="text-slate-500 dark:text-slate-400 mr-3">
                    {node.trackNumber || index + 1}.
                  </span>
                  <Link
                    href={`/tracks/${node.id}`}
                    className="text-emerald-600 hover:text-emerald-500"
                  >
                    {node.title}
                  </Link>
                </div>
                {formatDuration(node.duration) && (
                  <span className="text-sm text-slate-500 dark:text-slate-400">
                    {formatDuration(node.duration)}
                  </span>
                )}
              </li>
            ))}
          </ol>
        ) : (
          <p className="text-slate-500 dark:text-slate-400">
            No tracks stored for this album yet.
          </p>
        )}
      </div>

      <div className="card">
        <h2 className="text-2xl font-semibold mb-4">Reviews</h2>
        {album.reviews.edges.length > 0 ? (
          <div className="space-y-4">
            {album.reviews.edges.map(({ node }: any) => (
              <div
                key={node.id}
                className="border-b border-slate-200/70 pb-4 last:border-b-0 dark:border-slate-800/60"
              >
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-semibold">{node.user.name}</p>
                    <p className="text-sm text-slate-500 dark:text-slate-400">
                      {new Date(node.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="text-right">
                    <StarRating
                      value={node.rating}
                      label={`${node.rating}/5`}
                    />
                  </div>
                </div>
                {node.reviewText && (
                  <p className="text-slate-600 dark:text-slate-300 mt-2 italic">
                    "{node.reviewText}"
                  </p>
                )}
              </div>
            ))}
          </div>
        ) : (
          <p className="text-slate-500 dark:text-slate-400">
            No reviews yet for this album.
          </p>
        )}
      </div>
    </div>
  );
}
