import React from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_PLAYLIST } from "../../lib/graphql/queries";
import StarRating from "../../components/StarRating";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

export default function PlaylistDetailPage() {
  const router = useRouter();
  const { id } = router.query;

  const { data, loading, error, fetchMore } = useQuery(GET_PLAYLIST, {
    variables: { id, tracksFirst: 20 },
    skip: !id,
  });

  const loadMoreTracks = () => {
    if (data?.playlist?.tracks?.pageInfo?.hasNextPage) {
      fetchMore({
        variables: {
          id,
          tracksFirst: 20,
          tracksAfter: data.playlist.tracks.pageInfo.endCursor,
        },
        updateQuery: (prevResult, { fetchMoreResult }) => {
          if (!fetchMoreResult?.playlist?.tracks) {
            return prevResult;
          }
          return {
            ...prevResult,
            playlist: {
              ...prevResult.playlist,
              tracks: {
                ...fetchMoreResult.playlist.tracks,
                edges: [
                  ...(prevResult.playlist?.tracks?.edges ?? []),
                  ...fetchMoreResult.playlist.tracks.edges,
                ],
              },
            },
          };
        },
      });
    }
  };

  if (loading && !data) {
    return (
      <div className="flex justify-center py-12">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-emerald-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error Loading Playlist</h3>
        <p className="text-rose-500">{error.message}</p>
      </div>
    );
  }

  if (!data?.playlist) {
    return (
      <div className="card text-center">
        <p className="text-slate-500">Playlist not found</p>
        <Link href="/playlists" className="btn-primary mt-4 inline-block">
          Back to Playlists
        </Link>
      </div>
    );
  }

  const playlist = data.playlist;
  const tracks = playlist.tracks?.edges ?? [];
  const totalTracks = playlist.tracks?.totalCount ?? 0;

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col md:flex-row gap-6">
        {playlist.coverImage ? (
          <img
            src={playlist.coverImage}
            alt={playlist.title}
            className="w-48 h-48 object-cover rounded-2xl shadow-lg flex-shrink-0"
          />
        ) : (
          <div className="w-48 h-48 rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center flex-shrink-0">
            <span className="text-6xl">🎵</span>
          </div>
        )}

        <div className="flex-1">
          <p className="text-xs uppercase tracking-[0.2em] text-slate-400 mb-2">
            Playlist
          </p>
          <h1 className="text-3xl font-bold mb-2">{playlist.title}</h1>
          {playlist.description && (
            <p className="text-slate-600 dark:text-slate-300 mb-4">
              {playlist.description}
            </p>
          )}
          <div className="flex items-center gap-4 text-sm text-slate-500 dark:text-slate-400">
            <span>
              Created by{" "}
              <Link
                href={`/users/${playlist.creator.id}`}
                className="font-medium text-emerald-600 hover:text-emerald-500"
              >
                {playlist.creator.name}
              </Link>
            </span>
            <span>·</span>
            <span>{totalTracks} tracks</span>
            <span>·</span>
            <span>{new Date(playlist.createdAt).toLocaleDateString()}</span>
          </div>
        </div>
      </div>

      {/* Tracks */}
      <div className="card">
        <h2 className="text-xl font-semibold mb-4">Tracks</h2>

        {tracks.length > 0 ? (
          <div className="space-y-1">
            {tracks.map(({ node: track }: { node: any }, index: number) => (
              <Link
                key={track.id}
                href={`/tracks/${track.id}`}
                className="flex items-center gap-4 p-3 -mx-3 rounded-xl hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors group"
              >
                <span className="w-6 text-sm text-slate-400 text-right">
                  {index + 1}
                </span>

                {track.album?.coverImage ? (
                  <img
                    src={track.album.coverImage}
                    alt={track.album.title}
                    className="w-10 h-10 rounded-lg object-cover flex-shrink-0"
                  />
                ) : (
                  <div className="w-10 h-10 rounded-lg bg-slate-200 dark:bg-slate-700 flex-shrink-0" />
                )}

                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{track.title}</p>
                  <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                    {track.album?.artist?.name}
                    {track.album?.title && ` · ${track.album.title}`}
                  </p>
                </div>

                {track.averageRating && (
                  <StarRating
                    value={track.averageRating}
                    size={14}
                    className="flex-shrink-0"
                  />
                )}

                {formatDuration(track.duration) && (
                  <span className="text-sm text-slate-400 dark:text-slate-500 w-12 text-right">
                    {formatDuration(track.duration)}
                  </span>
                )}

                <span className="text-sm text-emerald-600 dark:text-emerald-400 opacity-0 group-hover:opacity-100 transition-opacity">
                  →
                </span>
              </Link>
            ))}
          </div>
        ) : (
          <p className="text-slate-500 dark:text-slate-400">
            This playlist is empty.
          </p>
        )}

        {data?.playlist?.tracks?.pageInfo?.hasNextPage && (
          <button
            onClick={loadMoreTracks}
            disabled={loading}
            className="btn-secondary mt-6 w-full"
          >
            {loading ? "Loading..." : "Load More Tracks"}
          </button>
        )}
      </div>
    </div>
  );
}
