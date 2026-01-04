import React from "react";
import SearchForm from "../components/SearchForm";

export default function SearchPage() {
  return (
    <div className="space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-semibold mb-4">Search Music</h1>
        <p className="text-lg text-slate-600 dark:text-slate-300">
          Test the Spotify API integration by searching for albums, artists, and
          tracks
        </p>
      </div>

      <SearchForm />

      <div className="card-muted">
        <h3 className="text-lg font-semibold mb-3">Testing Notes</h3>
        <ul className="list-disc list-inside space-y-2 text-slate-600 dark:text-slate-300">
          <li>This search uses your backend's Spotify API integration</li>
          <li>Results are fetched in real-time from Spotify</li>
          <li>
            Try searching for popular artists like "Taylor Swift" or "The
            Beatles"
          </li>
          <li>
            Switch between Albums, Artists, and Tracks to test different search
            types
          </li>
          <li>Check the Network tab in DevTools to see GraphQL queries</li>
        </ul>
      </div>
    </div>
  );
}
