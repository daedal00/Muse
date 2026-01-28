import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_ME } from "../lib/graphql/queries";

export default function DashboardPage() {
  const { data, loading, error } = useQuery(GET_ME);

  if (loading) {
    return (
      <div className="flex justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
      </div>
    );
  }

  const authError =
    error &&
    (error.message.toLowerCase().includes("unauthenticated") ||
      error.message.toLowerCase().includes("invalid token"));

  if (authError || (!loading && !data?.me)) {
    return (
      <div className="card text-center">
        <h2 className="text-xl font-semibold mb-2">You are not logged in</h2>
        <p className="muted mb-4">
          Log in to see your dashboard and saved activity.
        </p>
        <Link href="/auth" className="btn-primary">
          Go to Login
        </Link>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card">
        <h3 className="text-lg font-semibold mb-2">Error Loading Dashboard</h3>
        <p className="text-rose-500">{error.message}</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-semibold mb-2">Dashboard</h1>
        <p className="text-slate-600 dark:text-slate-300">
          Welcome back, {data.me.name}.
        </p>
      </div>

      <div className="card">
        <h2 className="text-xl font-semibold mb-4">Profile Summary</h2>
        <div className="space-y-1 text-slate-600 dark:text-slate-300">
          <p>
            <span className="font-medium">Name:</span> {data.me.name}
          </p>
          <p>
            <span className="font-medium">Email:</span> {data.me.email}
          </p>
          {data.me.bio && (
            <p>
              <span className="font-medium">Bio:</span> {data.me.bio}
            </p>
          )}
        </div>
      </div>

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Link
          href="/search"
          className="card hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Search Music</h3>
          <p className="text-slate-600 dark:text-slate-300 text-sm">
            Find albums, artists, and tracks from Spotify.
          </p>
        </Link>
        <Link
          href="/albums"
          className="card hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Albums</h3>
          <p className="text-slate-600 dark:text-slate-300 text-sm">
            Browse stored albums and view details.
          </p>
        </Link>
        <Link
          href="/reviews"
          className="card hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Reviews</h3>
          <p className="text-slate-600 dark:text-slate-300 text-sm">
            Read album reviews and ratings.
          </p>
        </Link>
        <Link
          href="/playlists"
          className="card hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Playlists</h3>
          <p className="text-slate-600 dark:text-slate-300 text-sm">
            Review playlist data and track lists.
          </p>
        </Link>
        <Link href="/profile" className="card hover:shadow-lg transition">
          <h3 className="text-lg font-semibold mb-2">Profile</h3>
          <p className="text-slate-600 dark:text-slate-300 text-sm">
            Manage Spotify imports and customize your home.
          </p>
        </Link>
      </div>
    </div>
  );
}
