import React, { useState, useEffect } from "react";
import { useMutation } from "@apollo/client";
import { UPDATE_PROFILE, UPDATE_PROFILE_SETTINGS } from "../lib/graphql/mutations";

interface ProfileSettings {
  layout: "GRID" | "LIST" | "BENTO";
  pinnedAlbumIds: string[];
  pinnedTrackIds: string[];
  sectionsOrder: string[];
  showSpotifyStats: boolean;
}

interface ProfileEditorProps {
  user: {
    id: string;
    name: string;
    bio?: string | null;
    avatar?: string | null;
  };
  settings?: ProfileSettings | null;
  onClose: () => void;
  onSave: () => void;
}

const DEFAULT_SECTIONS = ["favorites", "recentReviews", "topTracks", "playlists"];

const LAYOUT_OPTIONS = [
  { value: "GRID", label: "Grid", description: "Classic grid layout" },
  { value: "LIST", label: "List", description: "Compact list view" },
  { value: "BENTO", label: "Bento", description: "Modern asymmetric layout" },
] as const;

const SECTION_LABELS: Record<string, string> = {
  favorites: "Favorite Tracks",
  recentReviews: "Recent Reviews",
  topTracks: "Top Tracks",
  playlists: "Playlists",
  spotifyStats: "Spotify Stats",
};

export default function ProfileEditor({
  user,
  settings,
  onClose,
  onSave,
}: ProfileEditorProps) {
  const [name, setName] = useState(user.name);
  const [bio, setBio] = useState(user.bio || "");
  const [layout, setLayout] = useState<"GRID" | "LIST" | "BENTO">(
    settings?.layout || "GRID"
  );
  const [sectionsOrder, setSectionsOrder] = useState<string[]>(
    settings?.sectionsOrder || DEFAULT_SECTIONS
  );
  const [showSpotifyStats, setShowSpotifyStats] = useState(
    settings?.showSpotifyStats ?? true
  );
  const [draggedSection, setDraggedSection] = useState<string | null>(null);

  const [updateProfile, { loading: profileLoading }] = useMutation(UPDATE_PROFILE);
  const [updateProfileSettings, { loading: settingsLoading }] = useMutation(
    UPDATE_PROFILE_SETTINGS
  );

  const loading = profileLoading || settingsLoading;

  const handleDragStart = (section: string) => {
    setDraggedSection(section);
  };

  const handleDragOver = (e: React.DragEvent, targetSection: string) => {
    e.preventDefault();
    if (!draggedSection || draggedSection === targetSection) return;

    const newOrder = [...sectionsOrder];
    const draggedIdx = newOrder.indexOf(draggedSection);
    const targetIdx = newOrder.indexOf(targetSection);

    newOrder.splice(draggedIdx, 1);
    newOrder.splice(targetIdx, 0, draggedSection);
    setSectionsOrder(newOrder);
  };

  const handleDragEnd = () => {
    setDraggedSection(null);
  };

  const moveSection = (section: string, direction: "up" | "down") => {
    const idx = sectionsOrder.indexOf(section);
    if (
      (direction === "up" && idx === 0) ||
      (direction === "down" && idx === sectionsOrder.length - 1)
    ) {
      return;
    }

    const newOrder = [...sectionsOrder];
    const newIdx = direction === "up" ? idx - 1 : idx + 1;
    [newOrder[idx], newOrder[newIdx]] = [newOrder[newIdx], newOrder[idx]];
    setSectionsOrder(newOrder);
  };

  const handleSave = async () => {
    try {
      await updateProfile({
        variables: {
          input: {
            name,
            bio: bio || null,
          },
        },
      });

      await updateProfileSettings({
        variables: {
          input: {
            layout,
            sectionsOrder,
            showSpotifyStats,
          },
        },
      });

      onSave();
    } catch (err) {
      console.error("Failed to save profile:", err);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="w-full max-w-2xl max-h-[90vh] overflow-y-auto rounded-2xl bg-white dark:bg-slate-900 shadow-xl">
        <div className="sticky top-0 z-10 flex items-center justify-between border-b border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 px-6 py-4">
          <h2 className="text-xl font-semibold">Edit Profile</h2>
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg p-2 hover:bg-slate-100 dark:hover:bg-slate-800"
          >
            <svg
              className="h-5 w-5"
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
          </button>
        </div>

        <div className="space-y-8 p-6">
          {/* Basic Info */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Basic Info</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Name</label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-4 py-2 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">Bio</label>
                <textarea
                  value={bio}
                  onChange={(e) => setBio(e.target.value)}
                  rows={3}
                  placeholder="Tell others about your music taste..."
                  className="w-full rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-4 py-2 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500 resize-none"
                />
              </div>
            </div>
          </section>

          {/* Layout Selection */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Profile Layout</h3>
            <div className="grid grid-cols-3 gap-4">
              {LAYOUT_OPTIONS.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => setLayout(option.value)}
                  className={`rounded-xl border-2 p-4 text-left transition-all ${
                    layout === option.value
                      ? "border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20"
                      : "border-slate-200 dark:border-slate-700 hover:border-slate-300 dark:hover:border-slate-600"
                  }`}
                >
                  <div className="mb-2">
                    {option.value === "GRID" && (
                      <div className="grid grid-cols-2 gap-1">
                        <div className="h-3 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-3 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-3 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-3 rounded bg-slate-300 dark:bg-slate-600" />
                      </div>
                    )}
                    {option.value === "LIST" && (
                      <div className="space-y-1">
                        <div className="h-2 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-2 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-2 rounded bg-slate-300 dark:bg-slate-600" />
                      </div>
                    )}
                    {option.value === "BENTO" && (
                      <div className="grid grid-cols-2 gap-1">
                        <div className="col-span-2 h-3 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-4 rounded bg-slate-300 dark:bg-slate-600" />
                        <div className="h-4 rounded bg-slate-300 dark:bg-slate-600" />
                      </div>
                    )}
                  </div>
                  <p className="font-medium text-sm">{option.label}</p>
                  <p className="text-xs text-slate-500 dark:text-slate-400">
                    {option.description}
                  </p>
                </button>
              ))}
            </div>
          </section>

          {/* Section Order */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Section Order</h3>
            <p className="text-sm text-slate-500 dark:text-slate-400 mb-4">
              Drag sections to reorder them on your profile.
            </p>
            <div className="space-y-2">
              {sectionsOrder.map((section, idx) => (
                <div
                  key={section}
                  draggable
                  onDragStart={() => handleDragStart(section)}
                  onDragOver={(e) => handleDragOver(e, section)}
                  onDragEnd={handleDragEnd}
                  className={`flex items-center justify-between rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-3 cursor-move transition-all ${
                    draggedSection === section
                      ? "opacity-50 scale-[0.98]"
                      : "hover:border-slate-300 dark:hover:border-slate-600"
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <svg
                      className="h-5 w-5 text-slate-400"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M4 8h16M4 16h16"
                      />
                    </svg>
                    <span className="font-medium">
                      {SECTION_LABELS[section] || section}
                    </span>
                  </div>
                  <div className="flex items-center gap-1">
                    <button
                      type="button"
                      onClick={() => moveSection(section, "up")}
                      disabled={idx === 0}
                      className="rounded p-1 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-30"
                    >
                      <svg
                        className="h-4 w-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 15l7-7 7 7"
                        />
                      </svg>
                    </button>
                    <button
                      type="button"
                      onClick={() => moveSection(section, "down")}
                      disabled={idx === sectionsOrder.length - 1}
                      className="rounded p-1 hover:bg-slate-100 dark:hover:bg-slate-700 disabled:opacity-30"
                    >
                      <svg
                        className="h-4 w-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 9l-7 7-7-7"
                        />
                      </svg>
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </section>

          {/* Spotify Stats Toggle */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Privacy</h3>
            <label className="flex items-center justify-between rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-3">
              <div>
                <p className="font-medium">Show Spotify Stats</p>
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  Display your Spotify listening data on your public profile
                </p>
              </div>
              <input
                type="checkbox"
                checked={showSpotifyStats}
                onChange={(e) => setShowSpotifyStats(e.target.checked)}
                className="h-5 w-5 accent-emerald-500"
              />
            </label>
          </section>
        </div>

        <div className="sticky bottom-0 flex items-center justify-end gap-3 border-t border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 px-6 py-4">
          <button
            type="button"
            onClick={onClose}
            className="btn-secondary"
            disabled={loading}
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={handleSave}
            className="btn-primary"
            disabled={loading}
          >
            {loading ? "Saving..." : "Save Changes"}
          </button>
        </div>
      </div>
    </div>
  );
}
