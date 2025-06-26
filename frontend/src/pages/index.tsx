import React from "react";
import { useQuery } from "@apollo/client";
import { GET_ME } from "../lib/graphql/queries";

export default function HomePage() {
  const { data: user, loading, error } = useQuery(GET_ME);

  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Welcome to Muse 🎵
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Your music discovery and review platform
        </p>
      </div>

      {/* User Status */}
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h2 className="text-2xl font-semibold mb-4">Authentication Status</h2>
        {loading && (
          <div className="flex items-center space-x-2">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            <span>Checking authentication...</span>
          </div>
        )}
        {error && (
          <div className="text-red-600">
            <p>Authentication Error: {error.message}</p>
            <p className="text-sm mt-1">
              This might be expected if you're not logged in.
            </p>
          </div>
        )}
        {!loading && !error && user?.me && (
          <div className="text-green-600">
            <p>
              ✅ Authenticated as: <strong>{user.me.name}</strong>
            </p>
            <p className="text-sm text-gray-600">Email: {user.me.email}</p>
            {user.me.bio && (
              <p className="text-sm text-gray-600">Bio: {user.me.bio}</p>
            )}
          </div>
        )}
        {!loading && !user?.me && (
          <div className="text-yellow-600">
            <p>⚠️ Not authenticated</p>
            <p className="text-sm">Visit the Login page to authenticate.</p>
          </div>
        )}
      </div>

      {/* API Testing Features */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-3">🔍 Search</h3>
          <p className="text-gray-600 mb-4">
            Test album and artist search functionality using Spotify API
            integration.
          </p>
          <a
            href="/search"
            className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Try Search
          </a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-3">💿 Albums</h3>
          <p className="text-gray-600 mb-4">
            Browse stored albums and test pagination with GraphQL connections.
          </p>
          <a
            href="/albums"
            className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            View Albums
          </a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-3">⭐ Reviews</h3>
          <p className="text-gray-600 mb-4">
            View and create album reviews to test the review system.
          </p>
          <a
            href="/reviews"
            className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Browse Reviews
          </a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-3">🎵 Playlists</h3>
          <p className="text-gray-600 mb-4">
            Create and manage playlists to test the playlist functionality.
          </p>
          <a
            href="/playlists"
            className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            View Playlists
          </a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-3">🔐 Authentication</h3>
          <p className="text-gray-600 mb-4">
            Test user registration and login functionality.
          </p>
          <a
            href="/auth"
            className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          >
            Login/Register
          </a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h3 className="text-lg font-semibold mb-3">🛠️ GraphQL Playground</h3>
          <p className="text-gray-600 mb-4">
            Direct access to the GraphQL playground for advanced testing.
          </p>
          <a
            href="http://localhost:8080"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-block bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700"
          >
            Open Playground
          </a>
        </div>
      </div>

      {/* Backend Connection Test */}
      <div className="bg-yellow-50 border border-yellow-200 p-6 rounded-lg">
        <h3 className="text-lg font-semibold mb-3">🔗 Backend Connection</h3>
        <p className="text-gray-700 mb-3">
          This frontend is configured to connect to your backend at{" "}
          <code className="bg-gray-100 px-2 py-1 rounded text-sm">
            {process.env.NEXT_PUBLIC_GRAPHQL_URL ||
              "http://localhost:8080/query"}
          </code>
        </p>
        <p className="text-sm text-gray-600">
          Make sure your backend is running on port 8080 for the frontend to
          work properly.
        </p>
      </div>
    </div>
  );
}
