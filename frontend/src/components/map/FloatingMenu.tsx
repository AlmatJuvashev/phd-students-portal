import { useState, useEffect, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { 
  MessageCircle, 
  Calendar, 
  Settings, 
  Palette, 
  Plus, 
  X,
  RotateCcw, 
  Copy, 
  Check,
  Map as MapIcon
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface MenuItem {
  id: string;
  icon: React.ReactNode;
  label: string;
  onClick: () => void;
  color: string;
}

interface ThemeColors {
  primary: string;
  secondary: string;
}

// Helper functions for theme
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

const PRESETS: { name: string; primary: string; secondary: string }[] = [
  { name: "Blue", primary: "#6c8cff", secondary: "#a78bfa" },
  { name: "Emerald", primary: "#10b981", secondary: "#06b6d4" },
  { name: "Rose", primary: "#f43f5e", secondary: "#ec4899" },
  { name: "Amber", primary: "#f59e0b", secondary: "#ef4444" },
  { name: "Violet", primary: "#8b5cf6", secondary: "#d946ef" },
  { name: "Teal", primary: "#14b8a6", secondary: "#0ea5e9" },
];

const STORAGE_KEY = "phd-portal-theme";

export function FloatingMenu() {
  const [isOpen, setIsOpen] = useState(false);
  const [showThemePanel, setShowThemePanel] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const isChat = location.pathname === "/chat";

  // Theme State
  const [colors, setColors] = useState<ThemeColors>({
    primary: "#6c8cff",
    secondary: "#a78bfa",
  });
  const [copied, setCopied] = useState(false);

  // Load saved theme
  useEffect(() => {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setColors(parsed);
        applyColors(parsed);
      } catch {}
    }
  }, []);

  const applyColors = useCallback((c: ThemeColors) => {
    const root = document.documentElement;
    root.style.setProperty("--primary", hexToHsl(c.primary));
    root.style.setProperty("--sidebar-gradient-to", hexToHsl(c.secondary));
  }, []);

  const handleColorChange = (key: keyof ThemeColors, value: string) => {
    const newColors = { ...colors, [key]: value };
    setColors(newColors);
    applyColors(newColors);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newColors));
  };

  const applyPreset = (preset: (typeof PRESETS)[0]) => {
    const newColors = { primary: preset.primary, secondary: preset.secondary };
    setColors(newColors);
    applyColors(newColors);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(newColors));
  };

  const resetToDefault = () => {
    const defaultColors = { primary: "#6c8cff", secondary: "#a78bfa" };
    setColors(defaultColors);
    applyColors(defaultColors);
    localStorage.removeItem(STORAGE_KEY);
  };

  const menuItems: MenuItem[] = [
    {
      id: isChat ? "journey" : "chat",
      icon: isChat ? <MapIcon size={20} /> : <MessageCircle size={20} />,
      label: isChat ? "Journey Map" : "Chat",
      onClick: () => navigate(isChat ? "/journey" : "/chat"),
      color: "bg-blue-500",
    },
    {
      id: "calendar",
      icon: <Calendar size={20} />,
      label: "Calendar",
      onClick: () => navigate("/calendar"),
      color: "bg-green-500",
    },
    {
      id: "settings",
      icon: <Settings size={20} />,
      label: "Settings",
      onClick: () => navigate("/settings"),
      color: "bg-slate-500",
    },
    {
      id: "theme",
      icon: <Palette size={20} />,
      label: "Theme",
      onClick: () => setShowThemePanel(true),
      color: "bg-purple-500",
    },
  ];

  // Semicircle configuration
  const radius = 80;
  const startAngle = 180;
  const endAngle = 270;
  const totalAngle = endAngle - startAngle;

  return (
    <>
      {/* Theme Panel Modal/Overlay */}
      <AnimatePresence>
        {showThemePanel && (
           <div className="fixed inset-0 z-[70] flex items-center justify-center p-4 bg-black/20 backdrop-blur-sm" onClick={() => setShowThemePanel(false)}>
              <motion.div 
                initial={{ scale: 0.9, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                exit={{ scale: 0.9, opacity: 0 }}
                className="bg-white dark:bg-slate-900 rounded-2xl shadow-2xl p-6 w-full max-w-sm border border-slate-200 dark:border-slate-800"
                onClick={e => e.stopPropagation()}
              >
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="font-bold text-lg">Theme Settings</h3>
                    <Button variant="ghost" size="icon" onClick={() => setShowThemePanel(false)}>
                        <X size={20} />
                    </Button>
                  </div>

                  <div className="space-y-4">
                     {/* Presets */}
                     <div>
                        <label className="text-xs font-bold uppercase text-slate-400 mb-2 block">Presets</label>
                        <div className="grid grid-cols-6 gap-2">
                            {PRESETS.map((preset) => (
                                <button
                                    key={preset.name}
                                    onClick={() => applyPreset(preset)}
                                    className="w-8 h-8 rounded-full border-2 border-transparent hover:scale-110 transition-transform shadow-sm"
                                    style={{ background: `linear-gradient(135deg, ${preset.primary}, ${preset.secondary})` }}
                                    title={preset.name}
                                />
                            ))}
                        </div>
                     </div>

                     {/* Manual Pickers */}
                     <div className="space-y-3 pt-2 border-t border-slate-100">
                        <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">Primary</span>
                            <input 
                                type="color" 
                                value={colors.primary} 
                                onChange={(e) => handleColorChange("primary", e.target.value)}
                                className="w-8 h-8 rounded cursor-pointer border-0 p-0"
                            />
                        </div>
                        <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">Secondary</span>
                             <input 
                                type="color" 
                                value={colors.secondary} 
                                onChange={(e) => handleColorChange("secondary", e.target.value)}
                                className="w-8 h-8 rounded cursor-pointer border-0 p-0"
                            />
                        </div>
                     </div>

                     <div className="pt-4 flex gap-2">
                        <Button variant="outline" size="sm" className="flex-1" onClick={resetToDefault}>
                            <RotateCcw className="mr-2 h-4 w-4" /> Reset
                        </Button>
                     </div>
                  </div>
              </motion.div>
           </div>
        )}
      </AnimatePresence>

      <div className={cn(
        "fixed right-6 z-[60] flex items-center justify-center pointer-events-none transition-all duration-300",
        isChat ? "bottom-24" : "bottom-6"
      )}>
        
        {/* Menu Items Container - pointer-events-auto for children */}
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none"> 
            <AnimatePresence>
            {isOpen && (
                <>
                {menuItems.map((item, index) => {
                    const angle = startAngle + (index / (menuItems.length - 1)) * totalAngle;
                    const radian = (angle * Math.PI) / 180;
                    const x = Math.cos(radian) * radius;
                    const y = Math.sin(radian) * radius;

                    return (
                    <motion.button
                        key={item.id}
                        initial={{ x: 0, y: 0, scale: 0, opacity: 0 }}
                        animate={{ x, y, scale: 1, opacity: 1 }}
                        exit={{ x: 0, y: 0, scale: 0, opacity: 0 }}
                        transition={{ 
                        type: "spring", 
                        stiffness: 300, 
                        damping: 20, 
                        delay: index * 0.05 
                        }}
                        onClick={() => {
                        item.onClick();
                        setIsOpen(false);
                        }}
                        className={cn(
                        "absolute w-12 h-12 rounded-full flex items-center justify-center text-white shadow-xl pointer-events-auto",
                        item.color
                        )}
                        title={item.label}
                    >
                        {item.icon}
                    </motion.button>
                    );
                })}
                </>
            )}
            </AnimatePresence>
        </div>

        {/* Main Toggle Button */}
        <motion.button
            onClick={() => setIsOpen(!isOpen)}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            className={cn(
            "relative w-16 h-16 rounded-full flex items-center justify-center text-white shadow-2xl transition-colors z-[65] pointer-events-auto",
            isOpen ? "bg-slate-800" : "bg-primary hover:bg-primary/90"
            )}
        >
            <motion.div
                animate={{ rotate: isOpen ? 45 : 0 }}
                transition={{ duration: 0.2 }}
            >
                <Plus size={32} />
            </motion.div>
        </motion.button>
      </div>
    </>
  );
}
