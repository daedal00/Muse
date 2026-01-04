import React from "react";
import SearchForm from "../components/SearchForm";

export default function SearchPage() {
  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-semibold mb-2">Search</h1>
        <p className="text-slate-500 dark:text-slate-400">
          Find albums, artists, and tracks to rate and review
        </p>
      </div>

      <SearchForm />
    </div>
  );
}
