import React from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { useMutation, useQuery } from "@apollo/client";
import { GET_ME, GET_TRACK } from "../../lib/graphql/queries";
import { CREATE_TRACK_REVIEW } from "../../lib/graphql/mutations";
import StarRating from "../../components/StarRating";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

export default function TrackDetailPage() {
  const router = useRouter();
  const trackID = typeof router.query.id === "string" ? router.query.id : "";

  const { data, loading, error } = useQuery(GET_TRACK, {
    variables: { id: trackID, reviewsFirst: 20 },
    skip: !trackID,
  });

  const { data: meData } = useQuery(GET_ME);
  const isAuthed = Boolean(meData?.me);

  const [selectedRating, setSelectedRating] = React.useState<number>(0);
  const [saveMessage, setSaveMessage] = React.useState<string | null>(null);

  const [createTrackReview, { loading: saving, error: saveError }] =
    useMutation(CREATE_TRACK_REVIEW, {
      refetchQueries: trackID
        ? [
            {
              query: GET_TRACK,
              variables: { id: trackID, reviewsFirst: 20 },
            },
          ]
        : [],
      awaitRefetchQueries: true,
    });

  React.useEffect(() => {
    if (!data?.track?.reviews?.edges || !meData?.me?.id) {
      return;
    }
    const existing = data.track.reviews.edges.find(
      (edge: any) => edge.node.user.id === meData.me.id
    );
    if (existing?.node?.rating) {
      setSelectedRating(existing.node.rating);
    } else {
      setSelectedRating(0);
    }
  }, [data?.track?.reviews?.edges, meData?.me?.id]);

  const handleSaveRating = async () => {
    if (!trackID) return;
    if (selectedRating < 1) {
      setSaveMessage("Select a rating before saving.");
      return;
    }

    setSaveMessage(null);
    try {
      await createTrackReview({
        variables: {
          input: {
            trackId: trackID,
            rating: selectedRating,
          },
        },
      });
      setSaveMessage("Rating saved.");
    } catch (err) {
      setSaveMessage(
        err instanceof Error ? err.message : "Failed to save rating"
      );
    }
  };

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
        <h3 className="text-lg font-semibold mb-2">Error Loading Track</h3>
        <p className="text-rose-500">{error.message}</p>
      </div>
    );
  }

  if (!data?.track) {
    return (
      <div className="card text-center">
        <p className="text-amber-600 dark:text-amber-300">Track not found.</p>
        <Link href="/search" className="text-emerald-600 hover:text-emerald-500">
          Back to Search
        </Link>
      </div>
    );
  }

  const track = data.track;
  const trackAverage =
    typeof track.averageRating === "number" ? track.averageRating : null;
  const albumAverage =
    typeof track.album?.averageRating === "number"
      ? track.album.averageRating
      : null;

  return (
    <div className="space-y-8">
      <Link href="/search" className="text-emerald-600 hover:text-emerald-500">
        {"<- Back to Search"}
      </Link>

      <div className="card">
        <div className="flex flex-col md:flex-row md:items-start md:space-x-6">
          {track.album?.coverImage && (
            <img
              src={track.album.coverImage}
              alt={track.album.title}
              className="w-48 h-48 object-cover rounded mb-4 md:mb-0"
            />
          )}
          <div className="space-y-2">
            <h1 className="text-3xl font-semibold">{track.title}</h1>
            <p className="text-slate-600 dark:text-slate-300">
              Album:{" "}
              {track.album?.id ? (
                <Link
                  href={`/albums/${track.album.id}`}
                  className="text-emerald-600 hover:text-emerald-500"
                >
                  {track.album.title}
                </Link>
              ) : (
                track.album?.title
              )}
            </p>
            {track.album?.artist?.name && (
              <p className="text-slate-600 dark:text-slate-300">
                Artist: {track.album.artist.name}
              </p>
            )}
            {track.trackNumber && (
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Track {track.trackNumber}
              </p>
            )}
            {formatDuration(track.duration) && (
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Duration: {formatDuration(track.duration)}
              </p>
            )}
            {track.album?.releaseDate && (
              <p className="text-sm text-slate-500 dark:text-slate-400">
                Released: {new Date(track.album.releaseDate).toLocaleDateString()}
              </p>
            )}
            <p className="text-sm text-slate-500 dark:text-slate-400">
              Reviews: {track.reviews.totalCount}
            </p>
            <div className="pt-2 space-y-2">
              {trackAverage !== null ? (
                <StarRating
                  value={trackAverage}
                  label={`${trackAverage.toFixed(1)}/5 track rating`}
                />
              ) : (
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  No track ratings yet.
                </p>
              )}
              {albumAverage !== null ? (
                <StarRating
                  value={albumAverage}
                  label={`${albumAverage.toFixed(1)}/5 album rating`}
                />
              ) : (
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  No album ratings yet.
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <h2 className="text-2xl font-semibold mb-4">Rate this track</h2>
        {isAuthed ? (
          <div className="space-y-4">
            <StarRating value={selectedRating} onChange={setSelectedRating} />
            <button
              type="button"
              onClick={handleSaveRating}
              className="btn-primary"
              disabled={saving}
            >
              {saving ? "Saving..." : "Save rating"}
            </button>
            {saveMessage && (
              <p className="text-sm text-emerald-600">{saveMessage}</p>
            )}
            {saveError && (
              <p className="text-sm text-rose-500">{saveError.message}</p>
            )}
          </div>
        ) : (
          <div>
            <p className="muted">
              Login to add a rating to this track.
            </p>
            <Link href="/auth" className="btn-secondary mt-4">
              Login to rate
            </Link>
          </div>
        )}
      </div>

      <div className="card">
        <h2 className="text-2xl font-semibold mb-4">Track reviews</h2>
        {track.reviews.edges.length > 0 ? (
          <div className="space-y-4">
            {track.reviews.edges.map(({ node }: any) => (
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
                  <StarRating value={node.rating} label={`${node.rating}/5`} />
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
            No track reviews yet.
          </p>
        )}
      </div>
    </div>
  );
}
