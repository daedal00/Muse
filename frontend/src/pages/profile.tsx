import React, { useState } from "react";
import Link from "next/link";
import { useMutation, useQuery } from "@apollo/client";
import {
  GET_FAVORITE_TRACKS,
  GET_PROFILE,
  GET_SPOTIFY_PLAYLISTS,
  GET_SPOTIFY_SAVED_TRACKS,
  GET_SPOTIFY_STATUS,
  GET_SPOTIFY_TOP_TRACKS,
  GET_TOP_TRACKS,
} from "../lib/graphql/queries";
import {
  DISCONNECT_SPOTIFY,
  IMPORT_SPOTIFY_PLAYLIST,
  IMPORT_SPOTIFY_SAVED_TRACKS,
  IMPORT_SPOTIFY_TOP_TRACKS,
  SPOTIFY_AUTH_URL,
} from "../lib/graphql/mutations";
import StarRating from "../components/StarRating";
import ProfileEditor from "../components/ProfileEditor";
import {
  defaultHomePreferences,
  loadHomePreferences,
  saveHomePreferences,
} from "../lib/homePreferences";

const formatDuration = (durationSeconds?: number) => {
  if (!durationSeconds && durationSeconds !== 0) return null;
  const minutes = Math.floor(durationSeconds / 60);
  const seconds = durationSeconds % 60;
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

export default function ProfilePage() {
  const {
    data: profileData,
    loading: profileLoading,
    error: profileError,
    fetchMore,
  } = useQuery(GET_PROFILE, {
    variables: { trackReviewsFirst: 6 },
  });

  const isAuthed = Boolean(profileData?.me);

  const { data: favoriteData } = useQuery(GET_FAVORITE_TRACKS, {
    variables: { limit: 6, minRating: 4 },
    skip: !isAuthed,
  });

  const { data: topTracksData } = useQuery(GET_TOP_TRACKS, {
    variables: { limit: 6, sort: "RATING" },
  });

  const { data: spotifyStatusData } = useQuery(GET_SPOTIFY_STATUS, {
    skip: !isAuthed,
  });

  const spotifyConnected = Boolean(spotifyStatusData?.spotifyStatus?.connected);

  const { data: spotifyTopTracksData } = useQuery(GET_SPOTIFY_TOP_TRACKS, {
    variables: { limit: 10, timeRange: "MEDIUM_TERM" },
    skip: !spotifyConnected,
  });

  const { data: spotifySavedTracksData } = useQuery(GET_SPOTIFY_SAVED_TRACKS, {
    variables: { limit: 10, offset: 0 },
    skip: !spotifyConnected,
  });

  const { data: spotifyPlaylistsData } = useQuery(GET_SPOTIFY_PLAYLISTS, {
    variables: { limit: 10, offset: 0 },
    skip: !spotifyConnected,
  });

  const [spotifyAuthURL, { loading: spotifyAuthLoading }] =
    useMutation(SPOTIFY_AUTH_URL);
  const [disconnectSpotify, { loading: disconnectLoading }] =
    useMutation(DISCONNECT_SPOTIFY);
  const [importSpotifyTopTracks, { loading: importTopLoading }] = useMutation(
    IMPORT_SPOTIFY_TOP_TRACKS,
  );
  const [importSpotifySavedTracks, { loading: importSavedLoading }] =
    useMutation(IMPORT_SPOTIFY_SAVED_TRACKS);
  const [importSpotifyPlaylist] = useMutation(IMPORT_SPOTIFY_PLAYLIST);

  const [homePrefs, setHomePrefs] = React.useState(defaultHomePreferences);
  const [spotifyError, setSpotifyError] = React.useState<string | null>(null);
  const [importMessage, setImportMessage] = React.useState<string | null>(null);
  const [importingPlaylist, setImportingPlaylist] = React.useState<
    string | null
  >(null);
  const [showProfileEditor, setShowProfileEditor] = useState(false);

  React.useEffect(() => {
    setHomePrefs(loadHomePreferences());
  }, []);

  const handleHomePrefToggle = (key: keyof typeof homePrefs) => {
    const next = { ...homePrefs, [key]: !homePrefs[key] };
    setHomePrefs(next);
    saveHomePreferences(next);
  };

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
        err instanceof Error ? err.message : "Failed to start Spotify auth",
      );
    }
  };

  const handleSpotifyDisconnect = async () => {
    setSpotifyError(null);
    try {
      await disconnectSpotify();
      window.location.reload();
    } catch (err) {
      setSpotifyError(
        err instanceof Error ? err.message : "Failed to disconnect Spotify",
      );
    }
  };

  const handleImportTopTracks = async () => {
    setImportMessage(null);
    try {
      const { data } = await importSpotifyTopTracks({
        variables: { limit: 20, timeRange: "MEDIUM_TERM" },
      });
      if (data?.importSpotifyTopTracks) {
        setImportMessage(
          `Imported ${data.importSpotifyTopTracks.importedTracks} tracks.`,
        );
      }
    } catch (err) {
      setImportMessage(
        err instanceof Error ? err.message : "Failed to import top tracks",
      );
    }
  };

  const handleImportSavedTracks = async () => {
    setImportMessage(null);
    try {
      const { data } = await importSpotifySavedTracks({
        variables: { limit: 20, offset: 0 },
      });
      if (data?.importSpotifySavedTracks) {
        setImportMessage(
          `Imported ${data.importSpotifySavedTracks.importedTracks} tracks.`,
        );
      }
    } catch (err) {
      setImportMessage(
        err instanceof Error ? err.message : "Failed to import saved tracks",
      );
    }
  };

  const handleImportPlaylist = async (playlistID: string) => {
    setImportMessage(null);
    setImportingPlaylist(playlistID);
    try {
      const { data } = await importSpotifyPlaylist({
        variables: { spotifyPlaylistID: playlistID },
      });
      if (data?.importSpotifyPlaylist) {
        setImportMessage(
          `Imported ${data.importSpotifyPlaylist.importedPlaylists} playlist.`,
        );
      }
    } catch (err) {
      setImportMessage(
        err instanceof Error ? err.message : "Failed to import playlist",
      );
    } finally {
      setImportingPlaylist(null);
    }
  };

  if (profileLoading && !profileData) {
    return (
      <div className="flex justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-emerald-500"></div>
      </div>
    );
  }

  const authError =
    profileError &&
    (profileError.message.toLowerCase().includes("unauthenticated") ||
      profileError.message.toLowerCase().includes("invalid token"));

  if (authError || (!profileLoading && !profileData?.me)) {
    return (
      <div className="card text-center">
        <h2 className="text-xl font-semibold mb-2">You are not logged in</h2>
        <p className="muted mb-4">
          Log in to see your profile, ratings, and Spotify stats.
        </p>
        <Link href="/auth" className="btn-primary">
          Go to Login
        </Link>
      </div>
    );
  }

  if (profileError) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error Loading Profile</h3>
        <p className="text-rose-500">{profileError.message}</p>
      </div>
    );
  }

  const trackReviews = profileData?.me?.trackReviews?.edges ?? [];
  const favorites = favoriteData?.favoriteTracks ?? [];
  const topTracks = topTracksData?.topTracks ?? [];
  const spotifyTopTracks = spotifyTopTracksData?.spotifyTopTracks ?? [];
  const spotifySavedTracks = spotifySavedTracksData?.spotifySavedTracks ?? [];
  const spotifyPlaylists = spotifyPlaylistsData?.spotifyPlaylists ?? [];

  // Extract theme settings with defaults
  const settings = profileData?.me?.profileSettings;
  const primaryColor = settings?.primaryColor || "#1DB954";
  const accentColor = settings?.accentColor || "#191414";
  const backgroundStyle = settings?.backgroundStyle || "SOLID";
  const backgroundValue = settings?.backgroundValue || "";
  const bioStyle = settings?.bioStyle || "MINIMAL";
  const customTags = settings?.customTags || [];

  // Generate background CSS
  const getProfileBackground = () => {
    if (backgroundStyle === "GRADIENT") {
      return {
        background:
          backgroundValue ||
          `linear-gradient(135deg, ${primaryColor} 0%, ${accentColor} 100%)`,
      };
    }
    if (backgroundValue) {
      return { background: backgroundValue };
    }
    return {
      background: `linear-gradient(135deg, ${primaryColor}20 0%, ${accentColor}20 100%)`,
    };
  };

  // Determine if using dark background for text contrast
  const isGradientBg = backgroundStyle === "GRADIENT";

  // Bio rendering based on style
  const renderBio = () => {
    const bio = profileData?.me?.bio;
    if (!bio) return null;

    if (bioStyle === "QUOTE") {
      return (
        <blockquote
          className="mt-4 pl-4 border-l-4 italic"
          style={{
            borderColor: primaryColor,
            color: isGradientBg ? "white" : undefined,
          }}
        >
          &ldquo;{bio}&rdquo;
        </blockquote>
      );
    }

    if (bioStyle === "DETAILED") {
      return (
        <div
          className="mt-4 p-4 rounded-xl backdrop-blur-sm"
          style={{
            backgroundColor: isGradientBg
              ? "rgba(255,255,255,0.1)"
              : `${primaryColor}10`,
            borderColor: `${primaryColor}30`,
            borderWidth: "1px",
          }}
        >
          <p
            className={
              isGradientBg
                ? "text-white/90"
                : "text-slate-600 dark:text-slate-300"
            }
          >
            {bio}
          </p>
        </div>
      );
    }

    // MINIMAL
    return (
      <p
        className={`mt-3 text-sm ${isGradientBg ? "text-white/80" : "text-slate-600 dark:text-slate-300"}`}
      >
        {bio}
      </p>
    );
  };

  const loadMoreTrackReviews = () => {
    if (profileData?.me?.trackReviews?.pageInfo?.hasNextPage) {
      fetchMore({
        variables: {
          trackReviewsFirst: 6,
          trackReviewsAfter: profileData.me.trackReviews.pageInfo.endCursor,
        },
        updateQuery: (prevResult, { fetchMoreResult }) => {
          if (!fetchMoreResult?.me?.trackReviews) {
            return prevResult;
          }
          return {
            ...prevResult,
            me: {
              ...prevResult.me,
              trackReviews: {
                ...fetchMoreResult.me.trackReviews,
                edges: [
                  ...(prevResult.me?.trackReviews?.edges ?? []),
                  ...fetchMoreResult.me.trackReviews.edges,
                ],
              },
            },
          };
        },
      });
    }
  };

  return (
    <div className="space-y-10">
      {/* Profile Header with Custom Theme */}
      <section
        className={`rounded-2xl p-8 shadow-lg relative overflow-hidden ${
          isGradientBg
            ? "profile-bg-animated"
            : "border border-slate-200/70 dark:border-slate-800/60"
        }`}
        style={getProfileBackground()}
      >
        {/* Gradient overlay for better text visibility */}
        {isGradientBg && (
          <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent" />
        )}

        <div className="relative flex flex-wrap items-start justify-between gap-6">
          <div className="flex items-start gap-5">
            {profileData?.me?.avatar ? (
              <img
                src={profileData.me.avatar}
                alt={profileData.me.name}
                className="h-24 w-24 rounded-full object-cover shadow-xl"
                style={{
                  border: `4px solid ${isGradientBg ? "rgba(255,255,255,0.3)" : primaryColor}`,
                  boxShadow: `0 10px 30px ${primaryColor}40`,
                }}
              />
            ) : (
              <div
                className="flex h-24 w-24 items-center justify-center rounded-full text-3xl font-bold text-white shadow-xl"
                style={{
                  background: `linear-gradient(135deg, ${primaryColor} 0%, ${accentColor} 100%)`,
                  border: "4px solid rgba(255,255,255,0.3)",
                  boxShadow: `0 10px 30px ${primaryColor}40`,
                }}
              >
                {profileData?.me?.name?.charAt(0)?.toUpperCase() || "?"}
              </div>
            )}
            <div>
              <p
                className={`text-xs uppercase tracking-[0.2em] ${
                  isGradientBg ? "text-white/60" : "text-slate-400"
                }`}
              >
                Profile
              </p>
              <h1
                className={`text-3xl font-semibold mt-2 ${isGradientBg ? "text-white" : ""}`}
                style={!isGradientBg ? { color: primaryColor } : {}}
              >
                {profileData?.me?.name}
              </h1>
              <p className={`mt-1 ${isGradientBg ? "text-white/70" : "muted"}`}>
                {profileData?.me?.email}
              </p>
              {renderBio()}

              {/* Custom Tags */}
              {customTags.length > 0 && (
                <div className="flex flex-wrap gap-2 mt-4">
                  {customTags.map((tag: string) => (
                    <span
                      key={tag}
                      className="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium"
                      style={{
                        backgroundColor: isGradientBg
                          ? "rgba(255,255,255,0.2)"
                          : `${primaryColor}20`,
                        color: isGradientBg ? "white" : primaryColor,
                        borderColor: isGradientBg
                          ? "rgba(255,255,255,0.3)"
                          : primaryColor,
                        borderWidth: "1px",
                      }}
                    >
                      #{tag}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
          <div className="flex flex-col items-end gap-3">
            <button
              type="button"
              onClick={() => setShowProfileEditor(true)}
              className={`px-5 py-2 rounded-full font-semibold text-sm transition-all ${
                isGradientBg
                  ? "bg-white/20 text-white backdrop-blur-sm hover:bg-white/30 border border-white/30"
                  : "btn-secondary"
              }`}
              style={
                !isGradientBg
                  ? { borderColor: primaryColor, color: primaryColor }
                  : {}
              }
            >
              Edit Profile
            </button>
            <div
              className={`text-sm ${isGradientBg ? "text-white/70" : "text-slate-500 dark:text-slate-400"}`}
            >
              <p>
                Track ratings: {profileData?.me?.trackReviews?.totalCount ?? 0}
              </p>
            </div>
          </div>
        </div>
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <div className="card">
          <h2 className="text-xl font-semibold">Home customization</h2>
          <p className="muted mt-2">
            Pick what shows up on your personalized home.
          </p>
          <div className="mt-4 space-y-3">
            {(
              [
                { key: "favorites", label: "Favorites" },
                { key: "topTracks", label: "Top tracks" },
                { key: "playlists", label: "Playlists" },
                { key: "spotify", label: "Spotify highlights" },
              ] as const
            ).map(({ key, label }) => (
              <label key={key} className="flex items-center gap-3">
                <input
                  type="checkbox"
                  checked={homePrefs[key]}
                  onChange={() => handleHomePrefToggle(key)}
                  className="h-4 w-4 accent-emerald-500"
                />
                <span className="text-sm font-medium">{label}</span>
              </label>
            ))}
          </div>
        </div>

        <div className="card">
          <h2 className="text-xl font-semibold">Spotify connection</h2>
          <p className="muted mt-2">
            {spotifyConnected
              ? `Connected as ${spotifyStatusData?.spotifyStatus?.displayName || "Spotify user"}.`
              : "Connect Spotify to import stats and playlists."}
          </p>
          {spotifyError && (
            <p className="mt-2 text-sm text-rose-500">{spotifyError}</p>
          )}
          <div className="mt-4 flex flex-wrap gap-3">
            {spotifyConnected ? (
              <button
                type="button"
                onClick={handleSpotifyDisconnect}
                className="btn-secondary"
                disabled={disconnectLoading}
              >
                {disconnectLoading ? "Disconnecting..." : "Disconnect"}
              </button>
            ) : (
              <button
                type="button"
                onClick={handleSpotifyConnect}
                className="btn-primary"
                disabled={spotifyAuthLoading}
              >
                {spotifyAuthLoading ? "Connecting..." : "Connect Spotify"}
              </button>
            )}
          </div>
        </div>
      </section>

      <section className="card">
        <div className="flex items-end justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold flex items-center gap-2">
              <span
                className="w-1 h-5 rounded-full"
                style={{ backgroundColor: primaryColor }}
              />
              Recent track ratings
            </h2>
            <p className="muted">Your latest track reviews.</p>
          </div>
        </div>
        <div className="mt-4 space-y-4">
          {trackReviews.length > 0 ? (
            trackReviews.map(({ node }: any) => (
              <div key={node.id} className="flex items-center gap-4">
                {node.track?.album?.coverImage ? (
                  <img
                    src={node.track.album.coverImage}
                    alt={node.track.album.title}
                    className="h-14 w-14 rounded-lg object-cover"
                  />
                ) : (
                  <div className="h-14 w-14 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                )}
                <div>
                  <p className="font-medium">{node.track?.title}</p>
                  <p className="muted">{node.track?.album?.artist?.name}</p>
                  <StarRating
                    value={node.rating}
                    label={`${node.rating}/5`}
                    className="mt-1"
                    size={16}
                  />
                </div>
              </div>
            ))
          ) : (
            <p className="muted">No track ratings yet.</p>
          )}
        </div>
        {profileData?.me?.trackReviews?.pageInfo?.hasNextPage && (
          <button
            type="button"
            onClick={loadMoreTrackReviews}
            className="btn-secondary mt-6"
          >
            Load more
          </button>
        )}
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <div className="card">
          <h2 className="text-xl font-semibold">Favorite tracks</h2>
          <p className="muted mt-2">Tracks you rated 4-5 stars.</p>
          <div className="mt-4 space-y-4">
            {favorites.length > 0 ? (
              favorites.map((review: any) => (
                <div key={review.id} className="flex items-center gap-4">
                  {review.track?.album?.coverImage ? (
                    <img
                      src={review.track.album.coverImage}
                      alt={review.track.album.title}
                      className="h-12 w-12 rounded-lg object-cover"
                    />
                  ) : (
                    <div className="h-12 w-12 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                  )}
                  <div>
                    <p className="font-medium">{review.track?.title}</p>
                    <p className="muted">{review.track?.album?.artist?.name}</p>
                    <StarRating
                      value={review.rating}
                      label={`${review.rating}/5`}
                      className="mt-1"
                      size={16}
                    />
                  </div>
                </div>
              ))
            ) : (
              <p className="muted">No favorites yet.</p>
            )}
          </div>
        </div>

        <div className="card">
          <h2 className="text-xl font-semibold">Top tracks on Muse</h2>
          <p className="muted mt-2">
            Based on average ratings from the community.
          </p>
          <div className="mt-4 space-y-4">
            {topTracks.length > 0 ? (
              topTracks.map((insight: any) => (
                <div key={insight.track.id} className="flex items-center gap-4">
                  {insight.track?.album?.coverImage ? (
                    <img
                      src={insight.track.album.coverImage}
                      alt={insight.track.album.title}
                      className="h-12 w-12 rounded-lg object-cover"
                    />
                  ) : (
                    <div className="h-12 w-12 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                  )}
                  <div>
                    <p className="font-medium">{insight.track.title}</p>
                    <p className="muted">{insight.track.album?.artist?.name}</p>
                    <StarRating
                      value={insight.averageRating}
                      label={`${insight.averageRating.toFixed(1)} (${insight.reviewCount})`}
                      className="mt-1"
                      size={16}
                    />
                  </div>
                </div>
              ))
            ) : (
              <p className="muted">No top tracks yet.</p>
            )}
          </div>
        </div>
      </section>

      <section className="card">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold">Spotify imports</h2>
            <p className="muted">
              Import Spotify stats and bring tracks into Muse.
            </p>
          </div>
          {spotifyConnected && (
            <div className="flex flex-wrap gap-3">
              <button
                type="button"
                onClick={handleImportTopTracks}
                className="btn-secondary"
                disabled={importTopLoading}
              >
                {importTopLoading ? "Importing..." : "Import top tracks"}
              </button>
              <button
                type="button"
                onClick={handleImportSavedTracks}
                className="btn-secondary"
                disabled={importSavedLoading}
              >
                {importSavedLoading ? "Importing..." : "Import saved tracks"}
              </button>
            </div>
          )}
        </div>
        {importMessage && (
          <p className="mt-3 text-sm text-emerald-600">{importMessage}</p>
        )}

        {!spotifyConnected ? (
          <div className="mt-4 card-muted">
            <p className="muted">Connect Spotify to import data.</p>
          </div>
        ) : (
          <div className="mt-6 grid gap-6 lg:grid-cols-3">
            <div>
              <h3 className="text-lg font-semibold">Top tracks</h3>
              <div className="mt-3 space-y-3">
                {spotifyTopTracks.length > 0 ? (
                  spotifyTopTracks.map((track: any) => (
                    <div key={track.id} className="flex items-center gap-3">
                      {track.album?.coverImage ? (
                        <img
                          src={track.album.coverImage}
                          alt={track.album.title}
                          className="h-10 w-10 rounded-lg object-cover"
                        />
                      ) : (
                        <div className="h-10 w-10 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                      )}
                      <div>
                        <p className="text-sm font-medium">{track.title}</p>
                        <p className="muted">
                          {track.artists
                            ?.map((artist: any) => artist.name)
                            .join(", ")}
                        </p>
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="muted">No Spotify top tracks yet.</p>
                )}
              </div>
            </div>
            <div>
              <h3 className="text-lg font-semibold">Saved tracks</h3>
              <div className="mt-3 space-y-3">
                {spotifySavedTracks.length > 0 ? (
                  spotifySavedTracks.map((track: any) => (
                    <div key={track.id} className="flex items-center gap-3">
                      {track.album?.coverImage ? (
                        <img
                          src={track.album.coverImage}
                          alt={track.album.title}
                          className="h-10 w-10 rounded-lg object-cover"
                        />
                      ) : (
                        <div className="h-10 w-10 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                      )}
                      <div>
                        <p className="text-sm font-medium">{track.title}</p>
                        <p className="muted">
                          {formatDuration(track.duration)
                            ? `Duration ${formatDuration(track.duration)}`
                            : ""}
                        </p>
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="muted">No Spotify saved tracks yet.</p>
                )}
              </div>
            </div>
            <div>
              <h3 className="text-lg font-semibold">Playlists</h3>
              <div className="mt-3 space-y-3">
                {spotifyPlaylists.length > 0 ? (
                  spotifyPlaylists.map((playlist: any) => (
                    <div
                      key={playlist.id}
                      className="flex items-center justify-between gap-3"
                    >
                      <div className="flex items-center gap-3">
                        {playlist.coverImage ? (
                          <img
                            src={playlist.coverImage}
                            alt={playlist.name}
                            className="h-10 w-10 rounded-lg object-cover"
                          />
                        ) : (
                          <div className="h-10 w-10 rounded-lg bg-slate-200/70 dark:bg-slate-800" />
                        )}
                        <div>
                          <p className="text-sm font-medium">{playlist.name}</p>
                          <p className="muted">{playlist.trackCount} tracks</p>
                        </div>
                      </div>
                      <button
                        type="button"
                        onClick={() => handleImportPlaylist(playlist.id)}
                        className="btn-ghost"
                        disabled={importingPlaylist === playlist.id}
                      >
                        {importingPlaylist === playlist.id
                          ? "Importing..."
                          : "Import"}
                      </button>
                    </div>
                  ))
                ) : (
                  <p className="muted">No Spotify playlists yet.</p>
                )}
              </div>
            </div>
          </div>
        )}
      </section>

      {/* Profile Editor Modal */}
      {showProfileEditor && profileData?.me && (
        <ProfileEditor
          user={{
            id: profileData.me.id,
            name: profileData.me.name,
            bio: profileData.me.bio,
            avatar: profileData.me.avatar,
          }}
          settings={profileData.me.profileSettings}
          onClose={() => setShowProfileEditor(false)}
          onSave={() => {
            setShowProfileEditor(false);
            // Refetch profile data after saving
            window.location.reload();
          }}
        />
      )}
    </div>
  );
}
