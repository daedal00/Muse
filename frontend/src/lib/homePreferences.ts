export type HomeSectionKey = "favorites" | "topTracks" | "playlists" | "spotify";

export type HomePreferences = Record<HomeSectionKey, boolean>;

export const defaultHomePreferences: HomePreferences = {
  favorites: true,
  topTracks: true,
  playlists: true,
  spotify: true,
};

export const loadHomePreferences = (): HomePreferences => {
  if (typeof window === "undefined") {
    return defaultHomePreferences;
  }

  try {
    const raw = window.localStorage.getItem("muse-home-preferences");
    if (!raw) {
      return defaultHomePreferences;
    }
    const parsed = JSON.parse(raw) as Partial<HomePreferences>;
    return { ...defaultHomePreferences, ...parsed };
  } catch {
    return defaultHomePreferences;
  }
};

export const saveHomePreferences = (prefs: HomePreferences) => {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.setItem("muse-home-preferences", JSON.stringify(prefs));
};
