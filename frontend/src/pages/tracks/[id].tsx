import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { useQuery } from "@apollo/client";
import { GET_TRACK } from "../../lib/graphql/queries";

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

export default function TrackDetailPage() {
  const router = useRouter();
  const trackID = typeof router.query.id === "string" ? router.query.id : "";

  const { data, loading, error } = useQuery(GET_TRACK, {
    variables: { id: trackID },
    skip: !trackID,
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
          Error Loading Track
        </h3>
        <p className="text-red-600">{error.message}</p>
      </div>
    );
  }

  if (!data?.track) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
        <p className="text-yellow-800">Track not found.</p>
        <Link href="/search" className="text-blue-600 hover:text-blue-800">
          Back to Search
        </Link>
      </div>
    );
  }

  const track = data.track;
  const averageRating =
    typeof track.album?.averageRating === "number"
      ? track.album.averageRating
      : null;

  return (
    <div className="space-y-8">
      <Link href="/search" className="text-blue-600 hover:text-blue-800">
        {"<- Back to Search"}
      </Link>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <div className="flex flex-col md:flex-row md:items-start md:space-x-6">
          {track.album?.coverImage && (
            <img
              src={track.album.coverImage}
              alt={track.album.title}
              className="w-48 h-48 object-cover rounded mb-4 md:mb-0"
            />
          )}
          <div className="space-y-2">
            <h1 className="text-3xl font-bold text-gray-900">{track.title}</h1>
            <p className="text-gray-600">
              Album:{" "}
              {track.album?.id ? (
                <Link
                  href={`/albums/${track.album.id}`}
                  className="text-blue-600 hover:text-blue-800"
                >
                  {track.album.title}
                </Link>
              ) : (
                track.album?.title
              )}
            </p>
            {track.album?.artist?.name && (
              <p className="text-gray-600">Artist: {track.album.artist.name}</p>
            )}
            {track.trackNumber && (
              <p className="text-sm text-gray-500">
                Track {track.trackNumber}
              </p>
            )}
            {formatDuration(track.duration) && (
              <p className="text-sm text-gray-500">
                Duration: {formatDuration(track.duration)}
              </p>
            )}
            {track.album?.releaseDate && (
              <p className="text-sm text-gray-500">
                Released:{" "}
                {new Date(track.album.releaseDate).toLocaleDateString()}
              </p>
            )}
            <div className="pt-2">
              {averageRating !== null ? (
                <div className="text-lg">
                  {formatRating(averageRating)}{" "}
                  <span className="text-sm text-gray-600">
                    ({averageRating.toFixed(1)}/5 album rating)
                  </span>
                </div>
              ) : (
                <p className="text-sm text-gray-500">
                  No album ratings yet.
                </p>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
