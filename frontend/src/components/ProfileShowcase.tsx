import React from "react";
import Link from "next/link";

interface Album {
  id: string;
  title: string;
  coverImage?: string | null;
  artist: { id: string; name: string };
}

interface Track {
  id: string;
  title: string;
  album?: {
    id: string;
    title: string;
    coverImage?: string | null;
    artist: { id: string; name: string };
  } | null;
}

interface Artist {
  id: string;
  name: string;
}

interface ProfileShowcaseProps {
  pinnedAlbums?: Album[];
  pinnedTracks?: Track[];
  featuredArtists?: Artist[];
  customTags?: string[];
  primaryColor?: string;
  accentColor?: string;
  layout?: "GRID" | "LIST" | "BENTO";
}

export default function ProfileShowcase({
  pinnedAlbums = [],
  pinnedTracks = [],
  featuredArtists = [],
  customTags = [],
  primaryColor = "#1DB954",
  accentColor = "#191414",
  layout = "GRID",
}: ProfileShowcaseProps) {
  const hasContent =
    pinnedAlbums.length > 0 ||
    pinnedTracks.length > 0 ||
    featuredArtists.length > 0;

  if (!hasContent && customTags.length === 0) {
    return null;
  }

  return (
    <div className="space-y-8">
      {/* Custom Tags / Vibes */}
      {customTags.length > 0 && (
        <div className="flex flex-wrap gap-2 justify-center">
          {customTags.map((tag) => (
            <span
              key={tag}
              className="inline-flex items-center rounded-full px-4 py-1.5 text-sm font-medium transition-all hover:scale-105"
              style={{
                backgroundColor: `${primaryColor}20`,
                color: primaryColor,
                borderColor: primaryColor,
                borderWidth: "1px",
              }}
            >
              #{tag}
            </span>
          ))}
        </div>
      )}

      {/* Featured Artists */}
      {featuredArtists.length > 0 && (
        <section>
          <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span
              className="w-1 h-5 rounded-full"
              style={{ backgroundColor: primaryColor }}
            />
            Featured Artists
          </h3>
          <div
            className={`${layout === "BENTO" ? "bento-grid" : "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6"} gap-4`}
          >
            {featuredArtists.map((artist, i) => (
              <Link
                key={artist.id}
                href={`/artists/${artist.id}`}
                className={`group text-center p-4 rounded-2xl transition-all hover:scale-105 ${
                  layout === "BENTO" && i === 0 ? "bento-featured" : ""
                }`}
                style={{
                  backgroundColor: `${accentColor}30`,
                }}
              >
                <div
                  className="w-16 h-16 mx-auto mb-3 rounded-full flex items-center justify-center text-2xl font-bold"
                  style={{
                    background: `linear-gradient(135deg, ${primaryColor} 0%, ${accentColor} 100%)`,
                    color: "#fff",
                  }}
                >
                  {artist.name.charAt(0).toUpperCase()}
                </div>
                <p className="font-medium truncate">{artist.name}</p>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Pinned Tracks */}
      {pinnedTracks.length > 0 && (
        <section>
          <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span
              className="w-1 h-5 rounded-full"
              style={{ backgroundColor: primaryColor }}
            />
            Top Picks
          </h3>
          <div
            className={`${layout === "BENTO" ? "bento-grid-tracks" : layout === "LIST" ? "space-y-2" : "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4"} gap-4`}
          >
            {pinnedTracks.map((track, i) => (
              <Link
                key={track.id}
                href={`/tracks/${track.id}`}
                className={`group flex ${layout === "LIST" ? "flex-row items-center gap-4" : "flex-col"} rounded-xl overflow-hidden transition-all hover:scale-[1.02] ${
                  layout === "BENTO" && i === 0 ? "bento-featured" : ""
                }`}
                style={{
                  backgroundColor: `${accentColor}20`,
                }}
              >
                <div
                  className={`${layout === "LIST" ? "w-14 h-14" : "aspect-square w-full"} relative overflow-hidden`}
                >
                  {track.album?.coverImage ? (
                    <img
                      src={track.album.coverImage}
                      alt={track.title}
                      className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                    />
                  ) : (
                    <div
                      className="w-full h-full flex items-center justify-center"
                      style={{
                        background: `linear-gradient(135deg, ${primaryColor} 0%, ${accentColor} 100%)`,
                      }}
                    >
                      <svg
                        className="w-8 h-8 text-white/70"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z" />
                      </svg>
                    </div>
                  )}
                </div>
                <div className={`${layout === "LIST" ? "flex-1" : "p-3"}`}>
                  <p className="font-medium truncate">{track.title}</p>
                  {track.album?.artist && (
                    <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                      {track.album.artist.name}
                    </p>
                  )}
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Pinned Albums */}
      {pinnedAlbums.length > 0 && (
        <section>
          <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span
              className="w-1 h-5 rounded-full"
              style={{ backgroundColor: primaryColor }}
            />
            Favorite Albums
          </h3>
          <div
            className={`${layout === "BENTO" ? "bento-grid-albums" : layout === "LIST" ? "space-y-2" : "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4"} gap-4`}
          >
            {pinnedAlbums.map((album, i) => (
              <Link
                key={album.id}
                href={`/albums/${album.id}`}
                className={`group rounded-xl overflow-hidden transition-all hover:scale-[1.02] ${
                  layout === "BENTO" && i === 0 ? "bento-featured" : ""
                }`}
                style={{
                  backgroundColor: `${accentColor}20`,
                }}
              >
                <div className="aspect-square relative overflow-hidden">
                  {album.coverImage ? (
                    <img
                      src={album.coverImage}
                      alt={album.title}
                      className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                    />
                  ) : (
                    <div
                      className="w-full h-full flex items-center justify-center"
                      style={{
                        background: `linear-gradient(135deg, ${primaryColor} 0%, ${accentColor} 100%)`,
                      }}
                    >
                      <svg
                        className="w-12 h-12 text-white/70"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 14.5c-2.49 0-4.5-2.01-4.5-4.5S9.51 7.5 12 7.5s4.5 2.01 4.5 4.5-2.01 4.5-4.5 4.5zm0-5.5c-.55 0-1 .45-1 1s.45 1 1 1 1-.45 1-1-.45-1-1-1z" />
                      </svg>
                    </div>
                  )}
                </div>
                <div className="p-3">
                  <p className="font-medium truncate">{album.title}</p>
                  <p className="text-sm text-slate-500 dark:text-slate-400 truncate">
                    {album.artist.name}
                  </p>
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
