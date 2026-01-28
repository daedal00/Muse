import React, { useState } from "react";
import { useRouter } from "next/router";
import Link from "next/link";
import { useQuery, useMutation } from "@apollo/client";
import { GET_SPOTIFY_ALBUM_PREVIEW } from "../../../lib/graphql/queries";
import { IMPORT_ALBUM } from "../../../lib/graphql/mutations";

export default function SpotifyAlbumPreviewPage() {
  const router = useRouter();
  const { id: spotifyID } = router.query;
  const [importing, setImporting] = useState(false);

  const { data, loading, error } = useQuery(GET_SPOTIFY_ALBUM_PREVIEW, {
    variables: { spotifyID },
    skip: !spotifyID,
  });

  const [importAlbum] = useMutation(IMPORT_ALBUM);

  const handleImportAndRate = async () => {
    if (!spotifyID) return;
    setImporting(true);
    try {
      const { data: importData } = await importAlbum({
        variables: { spotifyAlbumID: spotifyID },
      });
      if (importData?.importAlbum?.id) {
        router.push(`/albums/${importData.importAlbum.id}`);
      }
    } catch (err) {
      console.error("Failed to import album:", err);
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
        <h3 className="text-lg font-semibold mb-2">Error Loading Album</h3>
        <p className="text-rose-500">{error.message}</p>
        <Link href="/search" className="btn-primary mt-4 inline-block">
          Back to Search
        </Link>
      </div>
    );
  }

  const album = data?.spotifyAlbum;

  if (!album) {
    return (
      <div className="card text-center">
        <p className="text-slate-500">Album not found</p>
        <Link href="/search" className="btn-primary mt-4 inline-block">
          Back to Search
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-8">
      {/* Album Header */}
      <div className="flex flex-col sm:flex-row gap-6 items-start">
        {album.coverImage ? (
          <img
            src={album.coverImage}
            alt={album.title}
            className="w-48 h-48 object-cover rounded-2xl shadow-lg flex-shrink-0"
          />
        ) : (
          <div className="w-48 h-48 rounded-2xl bg-gradient-to-br from-emerald-400 to-emerald-600 flex items-center justify-center flex-shrink-0">
            <span className="text-6xl">💿</span>
          </div>
        )}

        <div className="flex-1">
          <p className="text-xs uppercase tracking-[0.2em] text-slate-400 mb-2">
            Album Preview
          </p>
          <h1 className="text-3xl font-bold mb-2">{album.title}</h1>
          <p className="text-lg text-slate-600 dark:text-slate-300 mb-2">
            {album.artist?.map((a: any) => a.name).join(", ")}
          </p>
          {album.releaseDate && (
            <p className="text-sm text-slate-500 dark:text-slate-400 mb-4">
              {new Date(album.releaseDate).getFullYear()}
            </p>
          )}

          <div className="flex gap-3">
            <button
              onClick={handleImportAndRate}
              disabled={importing}
              className="btn-primary"
            >
              {importing ? "Opening..." : "Rate This Album"}
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
          This is a preview from Spotify. Click "Rate This Album" to add it to your
          collection and leave a review.
        </p>
      </div>
    </div>
  );
}
