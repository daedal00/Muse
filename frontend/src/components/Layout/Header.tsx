import React from "react";
import Link from "next/link";
import { useQuery } from "@apollo/client";
import { GET_ME } from "../../lib/graphql/queries";

const Header: React.FC = () => {
  const { data: user, loading } = useQuery(GET_ME);

  const handleLogout = () => {
    if (typeof window !== "undefined") {
      localStorage.removeItem("auth-token");
      window.location.reload();
    }
  };

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4 py-4">
        <div className="flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold text-blue-600">
            🎵 Muse
          </Link>

          <nav className="flex space-x-6">
            <Link href="/" className="text-gray-700 hover:text-blue-600">
              Home
            </Link>
            <Link href="/search" className="text-gray-700 hover:text-blue-600">
              Search
            </Link>
            <Link href="/albums" className="text-gray-700 hover:text-blue-600">
              Albums
            </Link>
            <Link href="/reviews" className="text-gray-700 hover:text-blue-600">
              Reviews
            </Link>
            <Link
              href="/playlists"
              className="text-gray-700 hover:text-blue-600"
            >
              Playlists
            </Link>
          </nav>

          <div className="flex items-center space-x-4">
            {loading ? (
              <div className="animate-pulse bg-gray-300 h-8 w-20 rounded"></div>
            ) : user?.me ? (
              <div className="flex items-center space-x-2">
                <span className="text-gray-700">Hello, {user.me.name}</span>
                <button
                  onClick={handleLogout}
                  className="text-red-600 hover:text-red-800"
                >
                  Logout
                </button>
              </div>
            ) : (
              <Link
                href="/auth"
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
              >
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
