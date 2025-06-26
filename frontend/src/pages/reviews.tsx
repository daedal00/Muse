import React from "react";
import { useQuery } from "@apollo/client";
import { GET_REVIEWS } from "../lib/graphql/queries";
import Layout from "../components/Layout/Layout";

export default function ReviewsPage() {
  const { data, loading, error, fetchMore } = useQuery(GET_REVIEWS, {
    variables: { first: 10 },
  });

  const loadMore = () => {
    if (data?.reviews.pageInfo.hasNextPage) {
      fetchMore({
        variables: {
          first: 10,
          after: data.reviews.pageInfo.endCursor,
        },
      });
    }
  };

  const formatRating = (rating: number) => {
    return "⭐".repeat(rating) + "☆".repeat(5 - rating);
  };

  return (
    <Layout>
      <div className="space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">Reviews</h1>
          <p className="text-lg text-gray-600">
            Browse album reviews and test the review system
          </p>
        </div>

        {loading && !data && (
          <div className="flex justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 p-6 rounded-lg">
            <h3 className="text-lg font-semibold text-red-800 mb-2">
              Error Loading Reviews
            </h3>
            <p className="text-red-600">{error.message}</p>
            <p className="text-sm text-red-500 mt-2">
              This might be expected if no reviews exist in the database yet.
            </p>
          </div>
        )}

        {data?.reviews && (
          <div className="space-y-6">
            <div className="bg-white p-4 rounded-lg shadow-md">
              <p className="text-gray-600">
                Total Reviews: <strong>{data.reviews.totalCount}</strong>
              </p>
            </div>

            {data.reviews.edges.length > 0 ? (
              <div className="space-y-4">
                {data.reviews.edges.map(({ node: review }: { node: any }) => (
                  <div
                    key={review.id}
                    className="bg-white p-6 rounded-lg shadow-md"
                  >
                    <div className="flex justify-between items-start mb-4">
                      <div>
                        <h3 className="text-lg font-semibold">
                          {review.album.title}
                        </h3>
                        <p className="text-gray-600">
                          by {review.album.artist.name}
                        </p>
                      </div>
                      <div className="text-right">
                        <div className="text-xl mb-1">
                          {formatRating(review.rating)}
                        </div>
                        <p className="text-sm text-gray-500">
                          {review.rating}/5
                        </p>
                      </div>
                    </div>

                    {review.reviewText && (
                      <p className="text-gray-700 mb-4 italic">
                        "{review.reviewText}"
                      </p>
                    )}

                    <div className="text-sm text-gray-500 border-t pt-3">
                      <p>
                        Reviewed by <strong>{review.user.name}</strong>
                      </p>
                      <p>{new Date(review.createdAt).toLocaleDateString()}</p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
                <p className="text-yellow-800">
                  No reviews found in the database.
                </p>
                <p className="text-sm text-yellow-600 mt-2">
                  Reviews will appear here once users start reviewing albums.
                </p>
              </div>
            )}

            {data.reviews.pageInfo.hasNextPage && (
              <div className="text-center">
                <button
                  onClick={loadMore}
                  disabled={loading}
                  className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
                >
                  {loading ? "Loading..." : "Load More"}
                </button>
              </div>
            )}
          </div>
        )}

        <div className="bg-blue-50 border border-blue-200 p-6 rounded-lg">
          <h3 className="text-lg font-semibold mb-3">🧪 Testing Notes</h3>
          <ul className="list-disc list-inside space-y-2 text-gray-700">
            <li>This page tests the GraphQL reviews query with pagination</li>
            <li>Reviews are created by authenticated users</li>
            <li>Each review includes a 1-5 star rating and optional text</li>
            <li>
              To test: create a user, login, then use the createReview mutation
            </li>
          </ul>
        </div>
      </div>
    </Layout>
  );
}
