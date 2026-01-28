import React, { useState } from "react";

interface ColorPickerProps {
  value: string;
  onChange: (color: string) => void;
  label?: string;
  presets?: string[];
}

const DEFAULT_PRESETS = [
  "#1DB954", // Spotify Green
  "#764ba2", // Purple
  "#667eea", // Indigo
  "#f093fb", // Pink
  "#f6d365", // Yellow
  "#48c6ef", // Cyan
  "#ee0979", // Hot Pink
  "#ff6a00", // Orange
  "#191414", // Spotify Black
  "#0f172a", // Slate 900
];

const GRADIENT_PRESETS = [
  "linear-gradient(135deg, #1DB954 0%, #191414 100%)",
  "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
  "linear-gradient(135deg, #f093fb 0%, #f5576c 100%)",
  "linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)",
  "linear-gradient(135deg, #f6d365 0%, #fda085 100%)",
  "linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%)",
  "linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)",
  "linear-gradient(135deg, #0f172a 0%, #1e3a5f 100%)",
];

export default function ColorPicker({
  value,
  onChange,
  label,
  presets = DEFAULT_PRESETS,
}: ColorPickerProps) {
  const [showCustom, setShowCustom] = useState(false);
  const [customColor, setCustomColor] = useState(value || "#1DB954");

  const handlePresetClick = (color: string) => {
    onChange(color);
    setShowCustom(false);
  };

  const handleCustomChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const color = e.target.value;
    setCustomColor(color);
    onChange(color);
  };

  const isGradient = value?.startsWith("linear-gradient");

  return (
    <div className="space-y-3">
      {label && (
        <label className="block text-sm font-medium text-slate-700 dark:text-slate-300">
          {label}
        </label>
      )}

      {/* Current Color Preview */}
      <div className="flex items-center gap-3">
        <div
          className="h-10 w-10 rounded-xl shadow-sm border border-slate-200 dark:border-slate-700"
          style={{ background: value || "#1DB954" }}
        />
        <span className="text-sm text-slate-600 dark:text-slate-400 font-mono">
          {value?.length > 30 ? "Gradient" : value}
        </span>
      </div>

      {/* Preset Colors */}
      <div className="flex flex-wrap gap-2">
        {presets.map((color) => (
          <button
            key={color}
            type="button"
            onClick={() => handlePresetClick(color)}
            className={`h-8 w-8 rounded-lg shadow-sm transition-transform hover:scale-110 ${
              value === color
                ? "ring-2 ring-emerald-500 ring-offset-2 dark:ring-offset-slate-900"
                : ""
            }`}
            style={{ background: color }}
            aria-label={`Select color ${color}`}
          />
        ))}
        <button
          type="button"
          onClick={() => setShowCustom(!showCustom)}
          className={`h-8 w-8 rounded-lg border-2 border-dashed border-slate-300 dark:border-slate-600 flex items-center justify-center text-slate-400 hover:border-emerald-500 hover:text-emerald-500 transition-colors ${
            showCustom ? "border-emerald-500 text-emerald-500" : ""
          }`}
          aria-label="Custom color"
        >
          <svg
            className="h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 6v6m0 0v6m0-6h6m-6 0H6"
            />
          </svg>
        </button>
      </div>

      {/* Custom Color Input */}
      {showCustom && (
        <div className="flex items-center gap-3 p-3 bg-slate-50 dark:bg-slate-800/50 rounded-xl">
          <input
            type="color"
            value={customColor}
            onChange={handleCustomChange}
            className="h-10 w-10 rounded-lg cursor-pointer border-0"
          />
          <input
            type="text"
            value={customColor}
            onChange={(e) => {
              setCustomColor(e.target.value);
              if (/^#[0-9A-Fa-f]{6}$/.test(e.target.value)) {
                onChange(e.target.value);
              }
            }}
            placeholder="#1DB954"
            className="flex-1 rounded-lg border border-slate-300 dark:border-slate-600 bg-white dark:bg-slate-900 px-3 py-2 text-sm font-mono focus:border-emerald-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>
      )}
    </div>
  );
}

export function GradientPicker({
  value,
  onChange,
  label,
}: {
  value: string;
  onChange: (gradient: string) => void;
  label?: string;
}) {
  return (
    <div className="space-y-3">
      {label && (
        <label className="block text-sm font-medium text-slate-700 dark:text-slate-300">
          {label}
        </label>
      )}

      {/* Current Gradient Preview */}
      <div
        className="h-16 w-full rounded-xl shadow-sm border border-slate-200 dark:border-slate-700"
        style={{ background: value || GRADIENT_PRESETS[0] }}
      />

      {/* Gradient Presets */}
      <div className="grid grid-cols-4 gap-2">
        {GRADIENT_PRESETS.map((gradient, index) => (
          <button
            key={index}
            type="button"
            onClick={() => onChange(gradient)}
            className={`h-10 rounded-lg shadow-sm transition-transform hover:scale-105 ${
              value === gradient
                ? "ring-2 ring-emerald-500 ring-offset-2 dark:ring-offset-slate-900"
                : ""
            }`}
            style={{ background: gradient }}
            aria-label={`Select gradient ${index + 1}`}
          />
        ))}
      </div>
    </div>
  );
}

export { GRADIENT_PRESETS, DEFAULT_PRESETS };
