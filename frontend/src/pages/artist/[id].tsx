import { useRouter } from "next/router";
import Layout from "../../components/Layout/Layout";
import { useQuery } from "@apollo/client";
import { GET_ARTIST_DETAILS } from "../../lib/graphql/queries";

export default function ArtistDetailsPage() {
  const router = useRouter();
  const { id } = router.query;

  const { data, loading, error, refetch } = useQuery(GET_ARTIST_DETAILS, {
    variables: { id: id as string },
    skip: !id,
  });

  const formatDuration = (seconds?: number) => {
    if (!seconds) return "";
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
  };

  const handleAlbumClick = (albumId: string) => {
    router.push(`/album/${albumId}`);
  };

  if (loading) {
    return (
      <Layout>
        <div className="flex justify-center items-center min-h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">
            Error loading artist
          </h2>
          <p className="text-gray-600 mb-4">{error.message}</p>
          <button
            onClick={() => refetch()}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Try Again
          </button>
        </div>
      </Layout>
    );
  }

  if (!data?.artistDetails) {
    return (
      <Layout>
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold text-gray-900 mb-2">
            Artist not found
          </h2>
          <button
            onClick={() => router.back()}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Go Back
          </button>
        </div>
      </Layout>
    );
  }

  const artist = data.artistDetails;

  return (
    <Layout>
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Artist Header */}
        <div className="bg-white rounded-lg shadow-md p-8 mb-8">
          <div className="flex items-start gap-6">
            {artist.image && (
              <img
                src={artist.image}
                alt={artist.name}
                className="w-40 h-40 object-cover rounded-lg shadow-md"
              />
            )}
            <div className="flex-1">
              <h1 className="text-4xl font-bold text-gray-900 mb-2">
                {artist.name}
              </h1>
              {artist.followers && (
                <p className="text-gray-600 mb-2">
                  {artist.followers.toLocaleString()} followers
                </p>
              )}
              {artist.genres && artist.genres.length > 0 && (
                <div className="flex flex-wrap gap-2 mb-4">
                  {artist.genres.map((genre, index) => (
                    <span
                      key={index}
                      className="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full"
                    >
                      {genre}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Top Tracks */}
        {artist.topTracks && artist.topTracks.length > 0 && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-8">
            <h2 className="text-2xl font-semibold text-gray-900 mb-6">
              Top Tracks
            </h2>
            <div className="space-y-3">
              {artist.topTracks.map((track, index) => (
                <div
                  key={track.id}
                  className="flex items-center justify-between p-3 hover:bg-gray-50 rounded-lg transition-colors"
                >
                  <div className="flex items-center gap-4">
                    <span className="text-gray-400 text-sm font-mono w-6">
                      {index + 1}
                    </span>
                    {track.album?.coverImage && (
                      <img
                        src={track.album.coverImage}
                        alt={track.album.title}
                        className="w-12 h-12 object-cover rounded"
                      />
                    )}
                    <div>
                      <h3 className="font-medium text-gray-900">
                        {track.title}
                      </h3>
                      {track.album && (
                        <p className="text-sm text-gray-600">
                          {track.album.title}
                        </p>
                      )}
                    </div>
                  </div>
                  {track.duration && (
                    <span className="text-sm text-gray-500 font-mono">
                      {formatDuration(track.duration)}
                    </span>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Albums */}
        {artist.albums && artist.albums.length > 0 && (
          <div className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-semibold text-gray-900 mb-6">
              Albums ({artist.albums.length})
            </h2>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
              {artist.albums.map((album) => (
                <div
                  key={album.id}
                  onClick={() => handleAlbumClick(album.id)}
                  className="cursor-pointer group"
                >
                  <div className="bg-gray-100 rounded-lg overflow-hidden mb-3 group-hover:shadow-lg transition-shadow">
                    {album.coverImage ? (
                      <img
                        src={album.coverImage}
                        alt={album.title}
                        className="w-full aspect-square object-cover group-hover:scale-105 transition-transform duration-200"
                      />
                    ) : (
                      <div className="w-full aspect-square bg-gray-200 flex items-center justify-center">
                        <span className="text-gray-400 text-sm">No Cover</span>
                      </div>
                    )}
                  </div>
                  <h3 className="font-medium text-gray-900 group-hover:text-blue-600 transition-colors line-clamp-2">
                    {album.title}
                  </h3>
                  {album.releaseDate && (
                    <p className="text-sm text-gray-500 mt-1">
                      {new Date(album.releaseDate).getFullYear()}
                    </p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Empty State */}
        {(!artist.albums || artist.albums.length === 0) &&
          (!artist.topTracks || artist.topTracks.length === 0) && (
            <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
              <p className="text-yellow-800">
                No albums or tracks found for this artist.
              </p>
            </div>
          )}
      </div>
    </Layout>
  );
}
