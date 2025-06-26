import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";

export default function SpotifyCallbackPage() {
  const router = useRouter();
  const [status, setStatus] = useState<"processing" | "success" | "error">(
    "processing"
  );
  const [message, setMessage] = useState("Processing Spotify connection...");

  useEffect(() => {
    const handleCallback = async () => {
      try {
        const { success, error } = router.query;

        if (error) {
          setStatus("error");
          let errorMessage = "Spotify authentication failed";

          switch (error) {
            case "access_denied":
              errorMessage =
                "You denied access to Spotify. Please try again and grant permission to continue.";
              break;
            case "invalid_client":
              errorMessage =
                "Application configuration error. Please contact support.";
              break;
            case "invalid_request":
              errorMessage = "Invalid request parameters. Please try again.";
              break;
            case "unauthorized_client":
              errorMessage =
                "Client authorization error. Please contact support.";
              break;
            case "unsupported_response_type":
              errorMessage =
                "Unsupported response type. Please contact support.";
              break;
            case "invalid_scope":
              errorMessage =
                "Invalid permissions requested. Please contact support.";
              break;
            case "missing_code":
              errorMessage =
                "Missing authorization code. Please try connecting again.";
              break;
            case "missing_state":
              errorMessage =
                "Missing security parameter. Please try connecting again.";
              break;
            case "state_expired":
              errorMessage =
                "Security token has expired. Please try connecting again.";
              break;
            case "state_invalid":
              errorMessage =
                "Invalid security token. Please try connecting again.";
              break;
            case "state_error":
              errorMessage =
                "Security validation error. Please try connecting again.";
              break;
            case "token_exchange_failed":
              errorMessage =
                "Failed to exchange authorization code. Please try again.";
              break;
            case "token_validation_failed":
              errorMessage = "Token validation failed. Please try again.";
              break;
            case "user_profile_failed":
              errorMessage =
                "Failed to retrieve Spotify profile. Please try again.";
              break;
            case "user_update_failed":
              errorMessage = "Failed to update user profile. Please try again.";
              break;
            case "auth_failed":
              errorMessage = "Authentication failed. Please try again.";
              break;
            case "service_unavailable":
              errorMessage =
                "Spotify service is temporarily unavailable. Please try again later.";
              break;
            case "unknown_error":
              errorMessage = "An unknown error occurred. Please try again.";
              break;
            default:
              errorMessage = `Spotify authentication failed: ${error}`;
          }

          setMessage(errorMessage);
          return;
        }

        if (success === "true") {
          setStatus("success");
          setMessage("Successfully connected to Spotify! Redirecting...");

          // Redirect to playlists page after a short delay
          setTimeout(() => {
            router.push("/playlists");
          }, 2000);
        } else {
          // Still processing or waiting for parameters
          if (router.isReady && !success && !error) {
            setStatus("error");
            setMessage("Unexpected response from Spotify authentication");
          }
        }
      } catch (err) {
        console.error("Error handling Spotify callback:", err);
        setStatus("error");
        setMessage("An unexpected error occurred");
      }
    };

    if (router.isReady) {
      handleCallback();
    }
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="bg-white p-8 rounded-lg shadow-md text-center max-w-md w-full">
        <div className="w-16 h-16 bg-green-500 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg
            className="w-8 h-8 text-white"
            fill="currentColor"
            viewBox="0 0 24 24"
          >
            <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.48.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.42 1.56-.299.421-1.02.599-1.559.3z" />
          </svg>
        </div>

        <h1 className="text-2xl font-bold mb-4">Spotify Integration</h1>

        {status === "processing" && (
          <>
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600 mx-auto mb-4"></div>
            <p className="text-gray-600">{message}</p>
          </>
        )}

        {status === "success" && (
          <>
            <div className="text-green-600 mb-4">
              <svg
                className="w-12 h-12 mx-auto"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <p className="text-green-600 font-semibold">{message}</p>
          </>
        )}

        {status === "error" && (
          <>
            <div className="text-red-600 mb-4">
              <svg
                className="w-12 h-12 mx-auto"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </div>
            <p className="text-red-600 font-semibold mb-4">{message}</p>
            <div className="space-y-2">
              <button
                onClick={() => router.push("/playlists")}
                className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 w-full"
              >
                Go to Playlists
              </button>
              <p className="text-sm text-gray-500">
                You can try connecting to Spotify again from the playlists page.
              </p>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
