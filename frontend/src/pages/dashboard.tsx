import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_ME } from "../lib/graphql/queries";

export default function DashboardPage() {
  const { data, loading, error } = useQuery(GET_ME);

  if (loading) {
    return (
      <div className="flex justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const authError =
    error &&
    (error.message.toLowerCase().includes("unauthenticated") ||
      error.message.toLowerCase().includes("invalid token"));

  if (authError || (!loading && !data?.me)) {
    return (
      <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg text-center">
        <h2 className="text-xl font-semibold text-yellow-800 mb-2">
          You are not logged in
        </h2>
        <p className="text-yellow-700 mb-4">
          Log in to see your dashboard and saved activity.
        </p>
        <Link
          href="/auth"
          className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Go to Login
        </Link>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 p-6 rounded-lg">
        <h3 className="text-lg font-semibold text-red-800 mb-2">
          Error Loading Dashboard
        </h3>
        <p className="text-red-600">{error.message}</p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Dashboard</h1>
        <p className="text-gray-600">Welcome back, {data.me.name}.</p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-semibold mb-4">Profile Summary</h2>
        <div className="space-y-1 text-gray-700">
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
          className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Search Music</h3>
          <p className="text-gray-600 text-sm">
            Find albums, artists, and tracks from Spotify.
          </p>
        </Link>
        <Link
          href="/albums"
          className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Albums</h3>
          <p className="text-gray-600 text-sm">
            Browse stored albums and view details.
          </p>
        </Link>
        <Link
          href="/reviews"
          className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Reviews</h3>
          <p className="text-gray-600 text-sm">
            Read album reviews and ratings.
          </p>
        </Link>
        <Link
          href="/playlists"
          className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition"
        >
          <h3 className="text-lg font-semibold mb-2">Playlists</h3>
          <p className="text-gray-600 text-sm">
            Review playlist data and track lists.
          </p>
        </Link>
      </div>
    </div>
  );
}
