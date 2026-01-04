import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_USERS } from "../../lib/graphql/queries";

export default function UsersPage() {
  const { data, loading, error, fetchMore } = useQuery(GET_USERS, {
    variables: { first: 20 },
  });

  const users = data?.users?.edges ?? [];
  const hasNextPage = data?.users?.pageInfo?.hasNextPage ?? false;
  const totalCount = data?.users?.totalCount ?? 0;

  const handleLoadMore = () => {
    if (hasNextPage) {
      fetchMore({
        variables: {
          first: 20,
          after: data?.users?.pageInfo?.endCursor,
        },
        updateQuery: (prev, { fetchMoreResult }) => {
          if (!fetchMoreResult?.users) return prev;
          return {
            users: {
              ...fetchMoreResult.users,
              edges: [...(prev.users?.edges ?? []), ...fetchMoreResult.users.edges],
            },
          };
        },
      });
    }
  };

  if (loading && users.length === 0) {
    return (
      <div className="flex justify-center py-12">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-emerald-500" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error</h3>
        <p className="text-rose-500">{error.message}</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <section className="card">
        <p className="text-xs uppercase tracking-[0.2em] text-slate-400">
          Community
        </p>
        <h1 className="text-3xl font-semibold mt-2">Discover Users</h1>
        <p className="muted mt-2">
          Explore music lovers and their ratings. {totalCount > 0 && `${totalCount} users on Muse.`}
        </p>
      </section>

      {users.length === 0 ? (
        <div className="card text-center py-12">
          <p className="text-slate-500 dark:text-slate-400">
            No users found. Be the first to sign up!
          </p>
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {users.map(({ node }: any) => (
            <Link
              key={node.id}
              href={`/users/${node.id}`}
              className="card group hover:border-emerald-500/50 transition-all duration-200"
            >
              <div className="flex items-center gap-4">
                {node.avatar ? (
                  <img
                    src={node.avatar}
                    alt={node.name}
                    className="h-14 w-14 rounded-full object-cover"
                  />
                ) : (
                  <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gradient-to-br from-emerald-400 to-emerald-600 text-xl font-bold text-white">
                    {node.name?.charAt(0)?.toUpperCase() || "?"}
                  </div>
                )}
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold truncate group-hover:text-emerald-500 transition-colors">
                    {node.name}
                  </h3>
                  {node.bio ? (
                    <p className="text-sm text-slate-500 dark:text-slate-400 line-clamp-2">
                      {node.bio}
                    </p>
                  ) : (
                    <p className="text-sm text-slate-400 dark:text-slate-500 italic">
                      No bio yet
                    </p>
                  )}
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      {hasNextPage && (
        <div className="flex justify-center">
          <button
            type="button"
            onClick={handleLoadMore}
            className="btn-secondary"
            disabled={loading}
          >
            {loading ? "Loading..." : "Load More"}
          </button>
        </div>
      )}
    </div>
  );
}
