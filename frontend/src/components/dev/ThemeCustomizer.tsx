import React, { useState, useEffect, useCallback } from "react";
import { Palette, X, RotateCcw, Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface ThemeColors {
  primary: string;
  secondary: string;
}

// Convert hex to HSL
function hexToHsl(hex: string): string {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  if (!result) return "231 96% 72%";

  let r = parseInt(result[1], 16) / 255;
  let g = parseInt(result[2], 16) / 255;
  let b = parseInt(result[3], 16) / 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0,
    s = 0,
    l = (max + min) / 2;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    switch (max) {
      case r:
        h = ((g - b) / d + (g < b ? 6 : 0)) / 6;
        break;
      case g:
        h = ((b - r) / d + 2) / 6;
        break;
      case b:
        h = ((r - g) / d + 4) / 6;
        break;
    }
  }

  return `${Math.round(h * 360)} ${Math.round(s * 100)}% ${Math.round(l * 100)}%`;
}

// Convert HSL string to hex
function hslToHex(hsl: string): string {
  const parts = hsl.split(" ");
  if (parts.length < 3) return "#6c8cff";

  const h = parseInt(parts[0]) / 360;
  const s = parseInt(parts[1]) / 100;
  const l = parseInt(parts[2]) / 100;

  const hue2rgb = (p: number, q: number, t: number) => {
    if (t < 0) t += 1;
    if (t > 1) t -= 1;
    if (t < 1 / 6) return p + (q - p) * 6 * t;
    if (t < 1 / 2) return q;
    if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
    return p;
  };

  const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
  const p = 2 * l - q;

  const r = Math.round(hue2rgb(p, q, h + 1 / 3) * 255);
  const g = Math.round(hue2rgb(p, q, h) * 255);
  const b = Math.round(hue2rgb(p, q, h - 1 / 3) * 255);

  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
}

// Preset color schemes
const PRESETS: { name: string; primary: string; secondary: string }[] = [
  { name: "Blue", primary: "#6c8cff", secondary: "#a78bfa" },
  { name: "Emerald", primary: "#10b981", secondary: "#06b6d4" },
  { name: "Rose", primary: "#f43f5e", secondary: "#ec4899" },
  { name: "Amber", primary: "#f59e0b", secondary: "#ef4444" },
  { name: "Violet", primary: "#8b5cf6", secondary: "#d946ef" },
  { name: "Teal", primary: "#14b8a6", secondary: "#0ea5e9" },
];

const STORAGE_KEY = "phd-portal-theme";

export function ThemeCustomizer() {
  const [isOpen, setIsOpen] = useState(false);
  const [colors, setColors] = useState<ThemeColors>({
    primary: "#6c8cff",
    secondary: "#a78bfa",
  });
  const [copied, setCopied] = useState(false);

  // Load saved theme on mount
  useEffect(() => {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setColors(parsed);
        applyColors(parsed);
      } catch {
        // Ignore parse errors
      }
    }
  }, []);

  // Apply colors to CSS variables
  const applyColors = useCallback((c: ThemeColors) => {
    const root = document.documentElement;
    root.style.setProperty("--primary", hexToHsl(c.primary));
    root.style.setProperty("--sidebar-gradient-to", hexToHsl(c.secondary));
  }, []);

  // Handle color change
  const handleColorChange = (key: keyof ThemeColors, value: string) => {
    const newColors = { ...colors, [key]: value };
    setColors(newColors);
    applyColors(newColors);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newColors));
  };

  // Apply preset
  const applyPreset = (preset: (typeof PRESETS)[0]) => {
    const newColors = { primary: preset.primary, secondary: preset.secondary };
    setColors(newColors);
    applyColors(newColors);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newColors));
  };

  // Reset to default
  const resetToDefault = () => {
    const defaultColors = { primary: "#6c8cff", secondary: "#a78bfa" };
    setColors(defaultColors);
    applyColors(defaultColors);
    localStorage.removeItem(STORAGE_KEY);
  };

  // Copy CSS variables
  const copyCSS = () => {
    const css = `:root {
  --primary: ${hexToHsl(colors.primary)};
  --sidebar-gradient-to: ${hexToHsl(colors.secondary)};
}`;
    navigator.clipboard.writeText(css);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // Only show in development - hides in production builds
  if (import.meta.env.PROD) return null;

  return (
    <>
      {/* Toggle Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={cn(
          "fixed bottom-4 right-4 z-50 p-3 rounded-full shadow-lg transition-all",
          "bg-gradient-primary text-white hover:scale-110",
          isOpen && "rotate-180"
        )}
        title="Theme Customizer (Dev Only)"
      >
        <Palette className="h-5 w-5" />
      </button>

      {/* Panel */}
      {isOpen && (
        <div className="fixed bottom-20 right-4 z-50 w-72 glass rounded-xl border shadow-xl p-4 space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="font-semibold text-sm">Theme Customizer</h3>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => setIsOpen(false)}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>

          {/* Color Pickers */}
          <div className="space-y-3">
            <div>
              <label className="text-xs text-muted-foreground mb-1 block">
                Primary Color
              </label>
              <div className="flex gap-2 items-center">
                <input
                  type="color"
                  value={colors.primary}
                  onChange={(e) => handleColorChange("primary", e.target.value)}
                  className="w-10 h-10 rounded cursor-pointer border-0 p-0"
                />
                <input
                  type="text"
                  value={colors.primary}
                  onChange={(e) => handleColorChange("primary", e.target.value)}
                  className="flex-1 text-xs px-2 py-1.5 rounded border bg-background font-mono"
                />
              </div>
            </div>

            <div>
              <label className="text-xs text-muted-foreground mb-1 block">
                Secondary Color
              </label>
              <div className="flex gap-2 items-center">
                <input
                  type="color"
                  value={colors.secondary}
                  onChange={(e) =>
                    handleColorChange("secondary", e.target.value)
                  }
                  className="w-10 h-10 rounded cursor-pointer border-0 p-0"
                />
                <input
                  type="text"
                  value={colors.secondary}
                  onChange={(e) =>
                    handleColorChange("secondary", e.target.value)
                  }
                  className="flex-1 text-xs px-2 py-1.5 rounded border bg-background font-mono"
                />
              </div>
            </div>
          </div>

          {/* Presets */}
          <div>
            <label className="text-xs text-muted-foreground mb-2 block">
              Presets
            </label>
            <div className="grid grid-cols-6 gap-1.5">
              {PRESETS.map((preset) => (
                <button
                  key={preset.name}
                  onClick={() => applyPreset(preset)}
                  className="group relative w-8 h-8 rounded-lg overflow-hidden border hover:scale-110 transition-transform"
                  title={preset.name}
                  style={{
                    background: `linear-gradient(135deg, ${preset.primary}, ${preset.secondary})`,
                  }}
                />
              ))}
            </div>
          </div>

          {/* Preview */}
          <div
            className="h-12 rounded-lg flex items-center justify-center text-white text-sm font-medium"
            style={{
              background: `linear-gradient(135deg, ${colors.primary}, ${colors.secondary})`,
            }}
          >
            Live Preview
          </div>

          {/* Actions */}
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={resetToDefault}
              className="flex-1"
            >
              <RotateCcw className="h-3 w-3 mr-1" />
              Reset
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={copyCSS}
              className="flex-1"
            >
              {copied ? (
                <Check className="h-3 w-3 mr-1" />
              ) : (
                <Copy className="h-3 w-3 mr-1" />
              )}
              {copied ? "Copied!" : "Copy CSS"}
            </Button>
          </div>

          <p className="text-[10px] text-muted-foreground text-center">
            Dev mode only • Changes persist locally
          </p>
        </div>
      )}
    </>
  );
}

export default ThemeCustomizer;
