import React from "react";
import Link from "next/link";
import { useMutation, useQuery } from "@apollo/client";
import {
  GET_FAVORITE_TRACKS,
  GET_ME,
  GET_PLAYLISTS,
  GET_SPOTIFY_PLAYLISTS,
  GET_SPOTIFY_STATUS,
  GET_SPOTIFY_TOP_TRACKS,
  GET_TOP_TRACKS,
} from "../lib/graphql/queries";
import { SPOTIFY_AUTH_URL } from "../lib/graphql/mutations";
import StarRating from "../components/StarRating";
import {
  defaultHomePreferences,
  loadHomePreferences,
} from "../lib/homePreferences";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

export default function HomePage() {
  const { data: meData } = useQuery(GET_ME);
  const isAuthed = Boolean(meData?.me);

  const [homePrefs, setHomePrefs] = React.useState(defaultHomePreferences);
  const [spotifyError, setSpotifyError] = React.useState<string | null>(null);

  React.useEffect(() => {
    setHomePrefs(loadHomePreferences());
  }, []);

  const { data: favoriteData } = useQuery(GET_FAVORITE_TRACKS, {
    variables: { limit: 6, minRating: 4 },
    skip: !isAuthed || !homePrefs.favorites,
  });

  const { data: topTracksData } = useQuery(GET_TOP_TRACKS, {
    variables: { limit: 6, sort: "RATING" },
    skip: !homePrefs.topTracks,
  });

  const { data: playlistsData } = useQuery(GET_PLAYLISTS, {
    variables: { first: 6 },
    skip: !homePrefs.playlists,
  });

  const { data: spotifyStatusData } = useQuery(GET_SPOTIFY_STATUS, {
    skip: !isAuthed || !homePrefs.spotify,
  });

  const spotifyConnected = Boolean(spotifyStatusData?.spotifyStatus?.connected);

  const { data: spotifyTopTracksData } = useQuery(GET_SPOTIFY_TOP_TRACKS, {
    variables: { limit: 6, timeRange: "SHORT_TERM" },
    skip: !spotifyConnected,
  });

  const { data: spotifyPlaylistsData } = useQuery(GET_SPOTIFY_PLAYLISTS, {
    variables: { limit: 6, offset: 0 },
    skip: !spotifyConnected,
  });

  const [spotifyAuthURL, { loading: spotifyAuthLoading }] = useMutation(
    SPOTIFY_AUTH_URL
  );

  const handleSpotifyConnect = async () => {
    setSpotifyError(null);
    if (!isAuthed || typeof window === "undefined") {
      return;
    }

    try {
      const redirectURI = `${window.location.origin}/profile`;
      const { data } = await spotifyAuthURL({
        variables: { redirectURI },
      });

      if (data?.spotifyAuthURL) {
        window.location.href = data.spotifyAuthURL;
      }
    } catch (err) {
      setSpotifyError(
        err instanceof Error ? err.message : "Failed to start Spotify auth"
      );
    }
  };

  const favorites = favoriteData?.favoriteTracks ?? [];
  const topTracks = topTracksData?.topTracks ?? [];
  const playlistEdges = playlistsData?.playlists?.edges ?? [];
  const playlists = playlistEdges.map((edge: any) => edge.node);
  const userPlaylists = isAuthed
    ? playlists.filter((playlist: any) => playlist.creator?.id === meData?.me?.id)
    : playlists;
  const spotifyTopTracks = spotifyTopTracksData?.spotifyTopTracks ?? [];
  const spotifyPlaylists = spotifyPlaylistsData?.spotifyPlaylists ?? [];

  const highlightTrack = topTracks[0]?.track;

  return (
    <div className="space-y-12">
      <section className="relative overflow-hidden rounded-3xl border border-slate-200/70 bg-slate-900 p-8 text-white shadow-sm md:p-12">
        <div className="absolute -right-24 top-0 h-64 w-64 rounded-full bg-emerald-500/30 blur-3xl" />
        <div className="absolute -left-32 bottom-0 h-72 w-72 rounded-full bg-amber-400/20 blur-3xl" />
        <div className="relative z-10 space-y-6">
          <span className="badge bg-white/15 text-white">Personal listening journal</span>
          <div>
            <h1 className="text-4xl font-semibold tracking-tight md:text-6xl">
              Muse
            </h1>
            <p className="mt-4 max-w-xl text-lg text-slate-200">
              Build a living archive of your favorite songs, review what sticks,
              and pull your Spotify stats into a home you actually want to look
              at.
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <Link href="/search" className="btn-primary">
              Search Spotify
            </Link>
            <Link href="/profile" className="btn-secondary">
              Open Profile
            </Link>
            {!isAuthed && (
              <Link href="/auth" className="btn-ghost">
                Login to personalize
              </Link>
            )}
          </div>
        </div>
        <div className="relative z-10 mt-10 grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-white/10 bg-white/10 p-4">
            <p className="text-xs uppercase tracking-[0.2em] text-slate-300">
              Top rated
            </p>
            {highlightTrack ? (
              <div className="mt-3">
                <p className="text-lg font-semibold">{highlightTrack.title}</p>
                <p className="text-sm text-slate-300">
                  {highlightTrack.album?.artist?.name}
                </p>
              </div>
            ) : (
              <p className="mt-3 text-sm text-slate-300">
                Rate some tracks to populate this section.
              </p>
            )}
          </div>
          <div className="rounded-2xl border border-white/10 bg-white/10 p-4">
            <p className="text-xs uppercase tracking-[0.2em] text-slate-300">
              Favorites
            </p>
            <p className="mt-3 text-sm text-slate-300">
              {favorites.length > 0
                ? `${favorites.length} tracks rated 4 stars and up.`
                : "Start rating tracks to build your list."}
            </p>
          </div>
          <div className="rounded-2xl border border-white/10 bg-white/10 p-4">
            <p className="text-xs uppercase tracking-[0.2em] text-slate-300">
              Spotify
            </p>
            <p className="mt-3 text-sm text-slate-300">
              {spotifyConnected
                ? "Connected and ready to import stats."
                : "Connect to pull in your Spotify insights."}
            </p>
          </div>
        </div>
      </section>

      {homePrefs.favorites && (
        <section className="space-y-4">
          <div className="flex flex-wrap items-end justify-between gap-4">
            <div>
              <h2 className="text-2xl font-semibold">Favorite tracks</h2>
              <p className="muted">Your 4-5 star picks, front and center.</p>
            </div>
            <Link
              href="/profile"
              className="text-sm font-semibold text-emerald-600 hover:text-emerald-500"
            >
              Manage favorites
            </Link>
          </div>
          {favorites.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {favorites.map((review: any) => (
                <div key={review.id} className="card-muted">
                  <div className="flex items-center gap-4">
                    {review.track?.album?.coverImage ? (
                      <img
                        src={review.track.album.coverImage}
                        alt={review.track.album.title}
                        className="h-16 w-16 rounded-xl object-cover"
                      />
                    ) : (
                      <div className="h-16 w-16 rounded-xl bg-slate-200/70 dark:bg-slate-800" />
                    )}
                    <div>
                      <p className="font-semibold">{review.track?.title}</p>
                      <p className="text-sm text-slate-600 dark:text-slate-300">
                        {review.track?.album?.artist?.name}
                      </p>
                      <StarRating
                        value={review.rating}
                        label={`${review.rating}/5`}
                        className="mt-2"
                        size={16}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="card-muted">
              <p className="muted">
                Rate some tracks to make your favorites show up here.
              </p>
            </div>
          )}
        </section>
      )}

      {homePrefs.topTracks && (
        <section className="space-y-4">
          <div>
            <h2 className="text-2xl font-semibold">Top tracks on Muse</h2>
            <p className="muted">
              Community-driven rankings based on track reviews.
            </p>
          </div>
          {topTracks.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {topTracks.map((insight: any) => (
                <div key={insight.track.id} className="card-muted">
                  <div className="flex items-center gap-4">
                    {insight.track?.album?.coverImage ? (
                      <img
                        src={insight.track.album.coverImage}
                        alt={insight.track.album.title}
                        className="h-16 w-16 rounded-xl object-cover"
                      />
                    ) : (
                      <div className="h-16 w-16 rounded-xl bg-slate-200/70 dark:bg-slate-800" />
                    )}
                    <div>
                      <p className="font-semibold">{insight.track.title}</p>
                      <p className="text-sm text-slate-600 dark:text-slate-300">
                        {insight.track.album?.artist?.name}
                      </p>
                      <StarRating
                        value={insight.averageRating}
                        label={`${insight.averageRating.toFixed(1)} (${insight.reviewCount})`}
                        className="mt-2"
                        size={16}
                      />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="card-muted">
              <p className="muted">
                No track reviews yet. Import tracks and start rating.
              </p>
            </div>
          )}
        </section>
      )}

      {homePrefs.playlists && (
        <section className="space-y-4">
          <div>
            <h2 className="text-2xl font-semibold">Playlists</h2>
            <p className="muted">
              {isAuthed
                ? "Your latest playlists from Muse."
                : "Community playlists in Muse."}
            </p>
          </div>
          {userPlaylists.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {userPlaylists.slice(0, 3).map((playlist: any) => (
                <div key={playlist.id} className="card-muted">
                  {playlist.coverImage ? (
                    <img
                      src={playlist.coverImage}
                      alt={playlist.title}
                      className="h-40 w-full rounded-xl object-cover"
                    />
                  ) : (
                    <div className="h-40 w-full rounded-xl bg-slate-200/70 dark:bg-slate-800" />
                  )}
                  <div className="mt-4">
                    <p className="font-semibold">{playlist.title}</p>
                    <p className="muted">
                      {playlist.creator?.name || "Muse user"}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="card-muted">
              <p className="muted">
                Create playlists in the Muse UI to surface them here.
              </p>
            </div>
          )}
        </section>
      )}

      {homePrefs.spotify && (
        <section className="space-y-4">
          <div>
            <h2 className="text-2xl font-semibold">Spotify highlights</h2>
            <p className="muted">
              Pull in your top tracks and playlists from Spotify.
            </p>
          </div>
          {!isAuthed ? (
            <div className="card-muted">
              <p className="muted">
                Log in to connect your Spotify account and import stats.
              </p>
              <Link href="/auth" className="btn-secondary mt-4">
                Login to connect
              </Link>
            </div>
          ) : spotifyConnected ? (
            <div className="grid gap-6 lg:grid-cols-2">
              <div className="card">
                <h3 className="text-lg font-semibold">Top tracks (recent)</h3>
                <div className="mt-4 space-y-4">
                  {spotifyTopTracks.length > 0 ? (
                    spotifyTopTracks.map((track: any) => (
                      <div
                        key={track.id}
                        className="flex items-center gap-4"
                      >
                        {track.album?.coverImage ? (
                          <img
                            src={track.album.coverImage}
                            alt={track.album.title}
                            className="h-12 w-12 rounded-lg object-cover"
                          />
                        ) : (
                          <div className="h-12 w-12 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                        )}
                        <div>
                          <p className="font-medium">{track.title}</p>
                          <p className="muted">
                            {track.artists
                              ?.map((artist: any) => artist.name)
                              .join(", ")}
                            {formatDuration(track.duration)
                              ? ` - ${formatDuration(track.duration)}`
                              : ""}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className="muted">No Spotify tracks loaded yet.</p>
                  )}
                </div>
              </div>
              <div className="card">
                <h3 className="text-lg font-semibold">Playlists</h3>
                <div className="mt-4 space-y-4">
                  {spotifyPlaylists.length > 0 ? (
                    spotifyPlaylists.slice(0, 4).map((playlist: any) => (
                      <div
                        key={playlist.id}
                        className="flex items-center gap-4"
                      >
                        {playlist.coverImage ? (
                          <img
                            src={playlist.coverImage}
                            alt={playlist.name}
                            className="h-12 w-12 rounded-lg object-cover"
                          />
                        ) : (
                          <div className="h-12 w-12 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                        )}
                        <div>
                          <p className="font-medium">{playlist.name}</p>
                          <p className="muted">
                            {playlist.trackCount} tracks -{" "}
                            {playlist.ownerName || "Spotify"}
                          </p>
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className="muted">No playlists found on Spotify.</p>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="card-muted">
              <p className="muted">
                Connect Spotify to import your top tracks, saved tracks, and
                playlists into Muse.
              </p>
              {spotifyError && (
                <p className="mt-2 text-sm text-rose-500">{spotifyError}</p>
              )}
              <button
                type="button"
                onClick={handleSpotifyConnect}
                className="btn-primary mt-4"
                disabled={spotifyAuthLoading}
              >
                {spotifyAuthLoading ? "Connecting..." : "Connect Spotify"}
              </button>
            </div>
          )}
        </section>
      )}
    </div>
  );
}
