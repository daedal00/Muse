import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import Layout from "../../components/Layout/Layout";
import { useQuery, useMutation } from "@apollo/client";
import { GET_ALBUM_DETAILS } from "../../lib/graphql/queries";
import { CREATE_REVIEW } from "../../lib/graphql/mutations";

interface TrackDetailsProps {
  track: {
    id: string;
    title: string;
    duration?: number;
    trackNumber?: number;
    featuredArtists: Array<{
      id: string;
      name: string;
    }>;
    averageRating?: number;
    totalReviews: number;
  };
  onRate: (trackId: string, rating: number) => void;
}

const TrackDetails: React.FC<TrackDetailsProps> = ({ track, onRate }) => {
  const [userRating, setUserRating] = useState<number>(0);
  const [showRatingInput, setShowRatingInput] = useState(false);

  const formatDuration = (seconds?: number) => {
    if (!seconds) return "";
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
  };

  const handleRatingSubmit = () => {
    if (userRating > 0) {
      onRate(track.id, userRating);
      setShowRatingInput(false);
      setUserRating(0);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-4 mb-3 hover:shadow-md transition-shadow">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-3">
            {track.trackNumber && (
              <span className="text-gray-400 text-sm font-mono w-6">
                {track.trackNumber}
              </span>
            )}
            <div>
              <h3 className="font-medium text-gray-900">{track.title}</h3>
              {track.featuredArtists.length > 0 && (
                <p className="text-sm text-gray-600">
                  feat. {track.featuredArtists.map((a) => a.name).join(", ")}
                </p>
              )}
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4">
          {/* Rating display */}
          <div className="text-center">
            {track.averageRating && (
              <div className="flex items-center gap-1">
                <span className="text-yellow-500">★</span>
                <span className="text-sm font-medium">
                  {track.averageRating.toFixed(1)}
                </span>
                <span className="text-xs text-gray-500">
                  ({track.totalReviews})
                </span>
              </div>
            )}
            {!track.averageRating && track.totalReviews === 0 && (
              <span className="text-xs text-gray-400">No ratings</span>
            )}
          </div>

          {/* Rating input */}
          {showRatingInput ? (
            <div className="flex items-center gap-2">
              <div className="flex gap-1">
                {[1, 2, 3, 4, 5].map((rating) => (
                  <button
                    key={rating}
                    onClick={() => setUserRating(rating)}
                    className={`text-lg ${
                      userRating >= rating
                        ? "text-yellow-500"
                        : "text-gray-300 hover:text-yellow-400"
                    }`}
                  >
                    ★
                  </button>
                ))}
              </div>
              <button
                onClick={handleRatingSubmit}
                disabled={userRating === 0}
                className="px-2 py-1 text-xs bg-blue-500 text-white rounded disabled:bg-gray-300"
              >
                Rate
              </button>
              <button
                onClick={() => {
                  setShowRatingInput(false);
                  setUserRating(0);
                }}
                className="px-2 py-1 text-xs bg-gray-500 text-white rounded"
              >
                Cancel
              </button>
            </div>
          ) : (
            <button
              onClick={() => setShowRatingInput(true)}
              className="px-3 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded"
            >
              Rate
            </button>
          )}

          {track.duration && (
            <span className="text-sm text-gray-500 font-mono">
              {formatDuration(track.duration)}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

export default function AlbumDetailsPage() {
  const router = useRouter();
  const { id } = router.query;

  const { data, loading, error, refetch } = useQuery(GET_ALBUM_DETAILS, {
    variables: { id: id as string },
    skip: !id,
  });

  const [createReview] = useMutation(CREATE_REVIEW, {
    onCompleted: () => {
      // Refetch album details to show updated ratings
      refetch();
    },
    onError: (error) => {
      console.error("Failed to create review:", error);
      alert("Failed to submit rating. Please try again.");
    },
  });

  const handleTrackRating = async (trackId: string, rating: number) => {
    try {
      await createReview({
        variables: {
          input: {
            trackId,
            rating,
          },
        },
      });
    } catch (error) {
      console.error("Error rating track:", error);
    }
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
            Error loading album
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

  if (!data?.albumDetails) {
    return (
      <Layout>
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold text-gray-900">
            Album not found
          </h2>
        </div>
      </Layout>
    );
  }

  const album = data.albumDetails;

  return (
    <Layout>
      <div className="max-w-6xl mx-auto px-4 py-8">
        {/* Album Header */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-8">
          <div className="flex gap-6">
            {album.coverImage && (
              <img
                src={album.coverImage}
                alt={album.title}
                className="w-48 h-48 rounded-lg shadow-md object-cover"
              />
            )}
            <div className="flex-1">
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                {album.title}
              </h1>
              <p className="text-xl text-gray-700 mb-4">
                by {album.artist.name}
              </p>
              {album.releaseDate && (
                <p className="text-gray-600 mb-4">
                  Released: {new Date(album.releaseDate).toLocaleDateString()}
                </p>
              )}

              {/* Album Rating */}
              <div className="flex items-center gap-4 mb-4">
                {album.averageRating && (
                  <div className="flex items-center gap-2">
                    <span className="text-yellow-500 text-xl">★</span>
                    <span className="text-lg font-semibold">
                      {album.averageRating.toFixed(1)}
                    </span>
                    <span className="text-gray-600">
                      ({album.totalReviews} reviews)
                    </span>
                  </div>
                )}
                {!album.averageRating && (
                  <span className="text-gray-500">No ratings yet</span>
                )}
              </div>

              <div className="text-gray-600">{album.tracks.length} tracks</div>
            </div>
          </div>
        </div>

        {/* Track List */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Tracks</h2>
          <div className="space-y-2">
            {album.tracks.map((track: any) => (
              <TrackDetails
                key={track.id}
                track={track}
                onRate={handleTrackRating}
              />
            ))}
          </div>
        </div>

        {/* Album Reviews Section */}
        {album.reviews && album.reviews.totalCount > 0 && (
          <div className="bg-white rounded-lg shadow-sm p-6 mt-8">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">
              Album Reviews
            </h2>
            <div className="space-y-4">
              {album.reviews.edges.map(({ node: review }: any) => (
                <div key={review.id} className="border-b border-gray-200 pb-4">
                  <div className="flex items-center gap-3 mb-2">
                    <div className="flex items-center gap-1">
                      {[...Array(5)].map((_, i) => (
                        <span
                          key={i}
                          className={`text-sm ${
                            i < review.rating
                              ? "text-yellow-500"
                              : "text-gray-300"
                          }`}
                        >
                          ★
                        </span>
                      ))}
                    </div>
                    <span className="font-medium text-gray-900">
                      {review.user.name}
                    </span>
                    <span className="text-sm text-gray-500">
                      {new Date(review.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                  {review.reviewText && (
                    <p className="text-gray-700">{review.reviewText}</p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}
