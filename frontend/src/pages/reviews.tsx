import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_REVIEWS } from "../lib/graphql/queries";
import StarRating from "../components/StarRating";

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

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-semibold mb-2">Reviews</h1>
        <p className="text-slate-500 dark:text-slate-400">
          See what people are saying about their favorite music
        </p>
      </div>

      {loading && !data && (
        <div className="flex justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
        </div>
      )}

      {error && (
        <div className="card">
          <h3 className="text-lg font-semibold mb-2">Error Loading Reviews</h3>
          <p className="text-rose-500">{error.message}</p>
          <p className="text-sm text-rose-400 mt-2">
            {error.message.includes('relation "reviews" does not exist') ||
            error.message.includes("SQLSTATE 42P01")
              ? "Database migrations are not applied yet. Run the backend migrations and retry."
              : "This might be expected if no reviews exist in the database yet."}
          </p>
        </div>
      )}

      {data?.reviews && (
        <div className="space-y-6">
          <div className="card">
            <p className="text-slate-600 dark:text-slate-300">
              Total Reviews: <strong>{data.reviews.totalCount}</strong>
            </p>
          </div>

          {data.reviews.edges.length > 0 ? (
            <div className="space-y-4">
              {data.reviews.edges.map(({ node: review }) => (
                <div key={review.id} className="card">
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <h3 className="text-lg font-semibold">
                        <Link
                          href={`/albums/${review.album.id}`}
                          className="text-emerald-600 hover:text-emerald-500"
                        >
                          {review.album.title}
                        </Link>
                      </h3>
                      <p className="text-slate-600 dark:text-slate-300">
                        by {review.album.artist.name}
                      </p>
                    </div>
                    <div className="text-right">
                      <StarRating
                        value={review.rating}
                        label={`${review.rating}/5`}
                      />
                    </div>
                  </div>

                  {review.reviewText && (
                    <p className="text-slate-600 dark:text-slate-300 mb-4 italic">
                      "{review.reviewText}"
                    </p>
                  )}

                  <div className="text-sm text-slate-500 dark:text-slate-400 border-t border-slate-200/70 dark:border-slate-800/60 pt-3">
                    <p>
                      Reviewed by <strong>{review.user.name}</strong>
                    </p>
                    <p>{new Date(review.createdAt).toLocaleDateString()}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="card text-center">
              <p className="text-amber-600 dark:text-amber-300">
                No reviews found in the database.
              </p>
              <p className="text-sm text-amber-500 mt-2">
                Reviews will appear here once users start reviewing albums.
              </p>
            </div>
          )}

          {data.reviews.pageInfo.hasNextPage && (
            <div className="text-center">
              <button
                onClick={loadMore}
                disabled={loading}
                className="btn-primary disabled:opacity-50"
              >
                {loading ? "Loading..." : "Load More"}
              </button>
            </div>
          )}
        </div>
      )}

    </div>
  );
}
