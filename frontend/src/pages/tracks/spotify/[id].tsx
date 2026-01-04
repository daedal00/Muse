import React, { useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import { useQuery, useMutation } from "@apollo/client";
import { GET_SPOTIFY_TRACK_PREVIEW } from "../../../lib/graphql/queries";
import { IMPORT_TRACK } from "../../../lib/graphql/mutations";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

export default function SpotifyTrackPreviewPage() {
  const router = useRouter();
  const { id: spotifyID } = router.query;
  const [importing, setImporting] = useState(false);

  const { data, loading, error } = useQuery(GET_SPOTIFY_TRACK_PREVIEW, {
    variables: { spotifyID },
    skip: !spotifyID,
  });

  const [importTrack] = useMutation(IMPORT_TRACK);

  const handleImportAndRate = async () => {
    if (!spotifyID) return;
    setImporting(true);
    try {
      const { data: importData } = await importTrack({
        variables: { spotifyTrackID: spotifyID },
      });
      if (importData?.importTrack?.id) {
        router.push(`/tracks/${importData.importTrack.id}`);
      }
    } catch (err) {
      console.error("Failed to import track:", err);
      setImporting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-emerald-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error Loading Track</h3>
        <p className="text-rose-500">{error.message}</p>
        <Link href="/search" className="btn-primary mt-4 inline-block">
          Back to Search
        </Link>
      </div>
    );
  }

  const track = data?.spotifyTrack;

  if (!track) {
    return (
      <div className="card text-center">
        <p className="text-slate-500">Track not found</p>
        <Link href="/search" className="btn-primary mt-4 inline-block">
          Back to Search
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-8">
      {/* Track Header */}
      <div className="flex flex-col sm:flex-row gap-6 items-start">
        {track.album?.coverImage ? (
          <img
            src={track.album.coverImage}
            alt={track.album.title}
            className="w-48 h-48 object-cover rounded-2xl shadow-lg flex-shrink-0"
          />
        ) : (
          <div className="w-48 h-48 rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center flex-shrink-0">
            <span className="text-6xl">🎵</span>
          </div>
        )}

        <div className="flex-1">
          <p className="text-xs uppercase tracking-[0.2em] text-slate-400 mb-2">
            Track Preview
          </p>
          <h1 className="text-3xl font-bold mb-2">{track.title}</h1>
          <p className="text-lg text-slate-600 dark:text-slate-300 mb-1">
            {track.artists?.map((a: any) => a.name).join(", ")}
          </p>
          {track.album?.title && (
            <p className="text-sm text-slate-500 dark:text-slate-400 mb-1">
              from {track.album.title}
            </p>
          )}
          <div className="flex items-center gap-3 text-sm text-slate-500 dark:text-slate-400 mb-4">
            {formatDuration(track.duration) && (
              <span>{formatDuration(track.duration)}</span>
            )}
            {track.trackNumber && (
              <>
                <span>·</span>
                <span>Track {track.trackNumber}</span>
              </>
            )}
          </div>

          <div className="flex gap-3">
            <button
              onClick={handleImportAndRate}
              disabled={importing}
              className="btn-primary"
            >
              {importing ? "Opening..." : "Rate This Track"}
            </button>
            <Link href="/search" className="btn-secondary">
              Back to Search
            </Link>
          </div>
        </div>
      </div>

      {/* Preview Notice */}
      <div className="card-muted">
        <p className="text-sm text-slate-600 dark:text-slate-300">
          This is a preview from Spotify. Click "Rate This Track" to add it to your
          collection and leave a review.
        </p>
      </div>
    </div>
  );
}
