import React, { useState } from "react";
import { useMutation } from "@apollo/client";
import {
  UPDATE_PROFILE,
  UPDATE_PROFILE_SETTINGS,
} from "../lib/graphql/mutations";
import ColorPicker, { GradientPicker, GRADIENT_PRESETS } from "./ColorPicker";

interface ProfileSettings {
  layout: "GRID" | "LIST" | "BENTO";
  primaryColor?: string | null;
  accentColor?: string | null;
  backgroundStyle?: "SOLID" | "GRADIENT" | "IMAGE" | null;
  backgroundValue?: string | null;
  pinnedAlbumIds: string[];
  pinnedTrackIds: string[];
  featuredArtistIds: string[];
  sectionsOrder: string[];
  showSpotifyStats: boolean;
  showListeningHistory: boolean;
  bioStyle?: "MINIMAL" | "DETAILED" | "QUOTE" | null;
  customTags: string[];
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

const DEFAULT_SECTIONS = [
  "featuredArtists",
  "topTracks",
  "pinnedAlbums",
  "recentReviews",
  "playlists",
];

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
  featuredArtists: "Featured Artists",
  pinnedAlbums: "Pinned Albums",
};

const BIO_STYLE_OPTIONS = [
  { value: "MINIMAL", label: "Minimal", description: "Clean and simple" },
  { value: "DETAILED", label: "Detailed", description: "Full description" },
  { value: "QUOTE", label: "Quote", description: "Styled as a quote" },
] as const;

export default function ProfileEditor({
  user,
  settings,
  onClose,
  onSave,
}: ProfileEditorProps) {
  const [name, setName] = useState(user.name);
  const [bio, setBio] = useState(user.bio || "");
  const [layout, setLayout] = useState<"GRID" | "LIST" | "BENTO">(
    settings?.layout || "GRID",
  );
  const [sectionsOrder, setSectionsOrder] = useState<string[]>(
    settings?.sectionsOrder || DEFAULT_SECTIONS,
  );
  const [showSpotifyStats, setShowSpotifyStats] = useState(
    settings?.showSpotifyStats ?? true,
  );
  const [showListeningHistory, setShowListeningHistory] = useState(
    settings?.showListeningHistory ?? true,
  );
  const [draggedSection, setDraggedSection] = useState<string | null>(null);

  // Theme customization state
  const [primaryColor, setPrimaryColor] = useState(
    settings?.primaryColor || "#1DB954",
  );
  const [accentColor, setAccentColor] = useState(
    settings?.accentColor || "#191414",
  );
  const [backgroundStyle, setBackgroundStyle] = useState<
    "SOLID" | "GRADIENT" | "IMAGE"
  >(settings?.backgroundStyle || "SOLID");
  const [backgroundValue, setBackgroundValue] = useState(
    settings?.backgroundValue || "#0f172a",
  );
  const [bioStyle, setBioStyle] = useState<"MINIMAL" | "DETAILED" | "QUOTE">(
    settings?.bioStyle || "MINIMAL",
  );
  const [customTags, setCustomTags] = useState<string[]>(
    settings?.customTags || [],
  );
  const [tagInput, setTagInput] = useState("");

  const [updateProfile, { loading: profileLoading }] =
    useMutation(UPDATE_PROFILE);
  const [updateProfileSettings, { loading: settingsLoading }] = useMutation(
    UPDATE_PROFILE_SETTINGS,
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

  const addTag = () => {
    const tag = tagInput.trim().toLowerCase();
    if (tag && !customTags.includes(tag) && customTags.length < 10) {
      setCustomTags([...customTags, tag]);
      setTagInput("");
    }
  };

  const removeTag = (tag: string) => {
    setCustomTags(customTags.filter((t) => t !== tag));
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
            primaryColor,
            accentColor,
            backgroundStyle,
            backgroundValue,
            sectionsOrder,
            showSpotifyStats,
            showListeningHistory,
            bioStyle,
            customTags,
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
      <div className="w-full max-w-3xl max-h-[90vh] overflow-y-auto rounded-2xl bg-white dark:bg-slate-900 shadow-xl">
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
              {/* Bio Style */}
              <div>
                <label className="block text-sm font-medium mb-2">
                  Bio Display Style
                </label>
                <div className="flex gap-3">
                  {BIO_STYLE_OPTIONS.map((option) => (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => setBioStyle(option.value)}
                      className={`flex-1 rounded-lg border-2 p-3 text-center transition-all ${
                        bioStyle === option.value
                          ? "border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20"
                          : "border-slate-200 dark:border-slate-700 hover:border-slate-300"
                      }`}
                    >
                      <p className="font-medium text-sm">{option.label}</p>
                      <p className="text-xs text-slate-500">
                        {option.description}
                      </p>
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </section>

          {/* Theme Customization */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Theme</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <ColorPicker
                label="Primary Color"
                value={primaryColor}
                onChange={setPrimaryColor}
              />
              <ColorPicker
                label="Accent Color"
                value={accentColor}
                onChange={setAccentColor}
              />
            </div>
          </section>

          {/* Background Customization */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Background</h3>
            <div className="space-y-4">
              <div className="flex gap-3">
                {(["SOLID", "GRADIENT"] as const).map((style) => (
                  <button
                    key={style}
                    type="button"
                    onClick={() => setBackgroundStyle(style)}
                    className={`flex-1 rounded-lg border-2 p-3 text-center transition-all ${
                      backgroundStyle === style
                        ? "border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20"
                        : "border-slate-200 dark:border-slate-700 hover:border-slate-300"
                    }`}
                  >
                    <span className="font-medium text-sm capitalize">
                      {style.toLowerCase()}
                    </span>
                  </button>
                ))}
              </div>
              {backgroundStyle === "SOLID" ? (
                <ColorPicker
                  label="Background Color"
                  value={backgroundValue}
                  onChange={setBackgroundValue}
                />
              ) : (
                <GradientPicker
                  label="Background Gradient"
                  value={backgroundValue}
                  onChange={setBackgroundValue}
                />
              )}
              {/* Preview */}
              <div
                className="h-24 rounded-xl shadow-inner border border-slate-200 dark:border-slate-700 flex items-center justify-center"
                style={{ background: backgroundValue }}
              >
                <span
                  className="font-semibold text-lg px-4 py-2 rounded-lg backdrop-blur-sm"
                  style={{
                    color: primaryColor,
                    backgroundColor: `${accentColor}80`,
                  }}
                >
                  Theme Preview
                </span>
              </div>
            </div>
          </section>

          {/* Custom Tags */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Music Vibes</h3>
            <p className="text-sm text-slate-500 dark:text-slate-400 mb-3">
              Add tags to describe your music taste (up to 10 tags).
            </p>
            <div className="flex gap-2 mb-3">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={(e) =>
                  e.key === "Enter" && (e.preventDefault(), addTag())
                }
                placeholder="indie, electronic, chill..."
                className="flex-1 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-800 px-4 py-2 focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
              <button
                type="button"
                onClick={addTag}
                disabled={!tagInput.trim() || customTags.length >= 10}
                className="btn-secondary disabled:opacity-50"
              >
                Add
              </button>
            </div>
            <div className="flex flex-wrap gap-2">
              {customTags.map((tag) => (
                <span
                  key={tag}
                  className="inline-flex items-center gap-1 rounded-full px-3 py-1 text-sm font-medium transition-all"
                  style={{
                    backgroundColor: `${primaryColor}20`,
                    color: primaryColor,
                    borderColor: primaryColor,
                    borderWidth: "1px",
                  }}
                >
                  #{tag}
                  <button
                    type="button"
                    onClick={() => removeTag(tag)}
                    className="ml-1 hover:opacity-70"
                  >
                    ×
                  </button>
                </span>
              ))}
              {customTags.length === 0 && (
                <span className="text-sm text-slate-400">No tags yet</span>
              )}
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

          {/* Privacy Settings */}
          <section>
            <h3 className="text-lg font-semibold mb-4">Privacy</h3>
            <div className="space-y-3">
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
              <label className="flex items-center justify-between rounded-lg border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 px-4 py-3">
                <div>
                  <p className="font-medium">Show Listening History</p>
                  <p className="text-sm text-slate-500 dark:text-slate-400">
                    Display your recent listening history on your profile
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={showListeningHistory}
                  onChange={(e) => setShowListeningHistory(e.target.checked)}
                  className="h-5 w-5 accent-emerald-500"
                />
              </label>
            </div>
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
