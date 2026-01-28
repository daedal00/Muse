import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_ME } from "../../lib/graphql/queries";
import ThemeToggle from "./ThemeToggle";

const Header: React.FC = () => {
  const { data: user, loading } = useQuery(GET_ME);

  const handleLogout = () => {
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth-token");
      window.location.reload();
    }
  };

  return (
    <header className="sticky top-0 z-40 border-b border-slate-200/70 bg-white/80 backdrop-blur dark:border-slate-800/60 dark:bg-slate-950/80">
      <div className="app-container py-4">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <Link
            href="/"
            className="text-2xl font-semibold tracking-tight text-emerald-600"
          >
            Muse
          </Link>

          <nav className="flex flex-wrap items-center gap-4 text-sm font-semibold">
            <Link
              href="/"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Home
            </Link>
            <Link
              href="/profile"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Profile
            </Link>
            <Link
              href="/dashboard"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Dashboard
            </Link>
            <Link
              href="/search"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Search
            </Link>
            <Link
              href="/albums"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Albums
            </Link>
            <Link
              href="/reviews"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Reviews
            </Link>
            <Link
              href="/playlists"
              className="text-slate-600 transition hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
            >
              Playlists
            </Link>
          </nav>

          <div className="flex items-center gap-3">
            <ThemeToggle />
            {loading ? (
              <div className="h-8 w-20 animate-pulse rounded-full bg-slate-200 dark:bg-slate-800"></div>
            ) : user?.me ? (
              <div className="flex items-center gap-2 text-sm">
                <span className="text-slate-600 dark:text-slate-300">
                  Hi, {user.me.name}
                </span>
                <button
                  onClick={handleLogout}
                  className="text-rose-500 transition hover:text-rose-400"
                >
                  Logout
                </button>
              </div>
            ) : (
              <Link href="/auth" className="btn-primary">
                Login
              </Link>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
