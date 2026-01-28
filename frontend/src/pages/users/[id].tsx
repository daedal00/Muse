import React from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_USER } from "../../lib/graphql/queries";
import StarRating from "../../components/StarRating";

export default function UserProfilePage() {
  const router = useRouter();
  const { id } = router.query;

  const { data, loading, error } = useQuery(GET_USER, {
    variables: { id },
    skip: !id,
  });

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-emerald-500" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error</h3>
        <p className="text-rose-500">{error.message}</p>
      </div>
    );
  }

  if (!data?.user) {
    return (
      <div className="card text-center py-12">
        <h2 className="text-xl font-semibold mb-2">User Not Found</h2>
        <p className="muted mb-4">This user doesn&apos;t exist or has been removed.</p>
        <Link href="/users" className="btn-primary">
          Browse Users
        </Link>
      </div>
    );
  }

  const user = data.user;
  const albumReviews = user.reviews?.edges ?? [];
  const trackReviews = user.trackReviews?.edges ?? [];
  const totalAlbumReviews = user.reviews?.totalCount ?? 0;
  const totalTrackReviews = user.trackReviews?.totalCount ?? 0;
  const settings = user.profileSettings;

  // Determine layout based on settings
  const layout = settings?.layout || "GRID";

  return (
    <div className="space-y-8">
      {/* Profile Header */}
      <section className="card">
        <div className="flex flex-wrap items-start gap-6">
          {user.avatar ? (
            <img
              src={user.avatar}
              alt={user.name}
              className="h-24 w-24 rounded-full object-cover"
            />
          ) : (
            <div className="flex h-24 w-24 items-center justify-center rounded-full bg-gradient-to-br from-emerald-400 to-emerald-600 text-3xl font-bold text-white">
              {user.name?.charAt(0)?.toUpperCase() || "?"}
            </div>
          )}
          <div className="flex-1">
            <p className="text-xs uppercase tracking-[0.2em] text-slate-400">
              Profile
            </p>
            <h1 className="text-3xl font-semibold mt-1">{user.name}</h1>
            {user.bio && (
              <p className="mt-3 text-slate-600 dark:text-slate-300">
                {user.bio}
              </p>
            )}
            <div className="mt-4 flex flex-wrap gap-4 text-sm text-slate-500 dark:text-slate-400">
              <span>{totalAlbumReviews} album reviews</span>
              <span>{totalTrackReviews} track ratings</span>
            </div>
          </div>
        </div>
      </section>

      {/* Album Reviews */}
      {albumReviews.length > 0 && (
        <section className="card">
          <h2 className="text-xl font-semibold">Album Reviews</h2>
          <p className="muted mt-1">
            Recent album reviews by {user.name}
          </p>
          <div
            className={`mt-6 ${
              layout === "GRID"
                ? "grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
                : layout === "LIST"
                ? "space-y-4"
                : "grid gap-4 sm:grid-cols-2"
            }`}
          >
            {albumReviews.map(({ node }: any) => (
              <Link
                key={node.id}
                href={`/albums/${node.album?.id}`}
                className={`group flex items-start gap-4 rounded-xl border border-slate-200 dark:border-slate-700 p-4 hover:border-emerald-500/50 transition-colors ${
                  layout === "BENTO" && albumReviews.indexOf({ node }) === 0
                    ? "sm:col-span-2"
                    : ""
                }`}
              >
                {node.album?.coverImage ? (
                  <img
                    src={node.album.coverImage}
                    alt={node.album.title}
                    className={`rounded-lg object-cover ${
                      layout === "LIST" ? "h-12 w-12" : "h-16 w-16"
                    }`}
                  />
                ) : (
                  <div
                    className={`rounded-lg bg-slate-200 dark:bg-slate-700 ${
                      layout === "LIST" ? "h-12 w-12" : "h-16 w-16"
                    }`}
                  />
                )}
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold truncate group-hover:text-emerald-500 transition-colors">
                    {node.album?.title}
                  </h3>
                  <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                    {node.album?.artist?.name}
                  </p>
                  <StarRating
                    value={node.rating}
                    label={`${node.rating}/5`}
                    className="mt-2"
                    size={14}
                  />
                  {node.reviewText && (
                    <p className="mt-2 text-sm text-slate-600 dark:text-slate-300 line-clamp-2">
                      {node.reviewText}
                    </p>
                  )}
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Track Reviews */}
      {trackReviews.length > 0 && (
        <section className="card">
          <h2 className="text-xl font-semibold">Track Ratings</h2>
          <p className="muted mt-1">
            Recent track ratings by {user.name}
          </p>
          <div
            className={`mt-6 ${
              layout === "GRID"
                ? "grid gap-4 sm:grid-cols-2 lg:grid-cols-3"
                : layout === "LIST"
                ? "space-y-4"
                : "grid gap-4 sm:grid-cols-2"
            }`}
          >
            {trackReviews.map(({ node }: any) => (
              <Link
                key={node.id}
                href={`/tracks/${node.track?.album?.id}`}
                className="group flex items-start gap-4 rounded-xl border border-slate-200 dark:border-slate-700 p-4 hover:border-emerald-500/50 transition-colors"
              >
                {node.track?.album?.coverImage ? (
                  <img
                    src={node.track.album.coverImage}
                    alt={node.track.album.title}
                    className={`rounded-lg object-cover ${
                      layout === "LIST" ? "h-12 w-12" : "h-16 w-16"
                    }`}
                  />
                ) : (
                  <div
                    className={`rounded-lg bg-slate-200 dark:bg-slate-700 ${
                      layout === "LIST" ? "h-12 w-12" : "h-16 w-16"
                    }`}
                  />
                )}
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold truncate group-hover:text-emerald-500 transition-colors">
                    {node.track?.title}
                  </h3>
                  <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                    {node.track?.album?.artist?.name} - {node.track?.album?.title}
                  </p>
                  <StarRating
                    value={node.rating}
                    label={`${node.rating}/5`}
                    className="mt-2"
                    size={14}
                  />
                  {node.reviewText && (
                    <p className="mt-2 text-sm text-slate-600 dark:text-slate-300 line-clamp-2">
                      {node.reviewText}
                    </p>
                  )}
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Empty State */}
      {albumReviews.length === 0 && trackReviews.length === 0 && (
        <div className="card text-center py-12">
          <p className="text-slate-500 dark:text-slate-400">
            {user.name} hasn&apos;t rated any music yet.
          </p>
        </div>
      )}
    </div>
  );
}
