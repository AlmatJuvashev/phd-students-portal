
import React from 'react';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { cva, type VariantProps } from 'class-variance-authority'; // Standardizing on CVA if available, but for now mimicking the provided code accurately

// --- Presence / Collaboration ---
export const AvatarGroup = ({ users }: { users: { initials: string, color: string }[] }) => (
  <div className="flex -space-x-2 overflow-hidden">
    {users.map((user, i) => (
      <div 
        key={i} 
        className={cn(
          "inline-block h-8 w-8 rounded-full ring-2 ring-white flex items-center justify-center text-[10px] font-bold text-white shadow-sm transition-transform hover:scale-110 hover:z-10 cursor-help",
          user.color
        )}
        title="Active Editor"
      >
        {user.initials}
      </div>
    ))}
    <div className="h-8 w-8 rounded-full bg-slate-100 border-2 border-dashed border-slate-300 flex items-center justify-center text-slate-400 cursor-pointer hover:bg-slate-50 transition-colors">
       <span className="text-xs font-bold">+</span>
    </div>
  </div>
);

export const Badge = ({ children, variant = 'default', className }: { children?: React.ReactNode, variant?: 'default' | 'outline' | 'secondary' | 'success' | 'warning' | 'destructive' | 'purple', className?: string }) => {
  const variants: Record<string, string> = {
    default: "bg-slate-900 text-white hover:bg-slate-800",
    outline: "text-slate-950 border border-slate-200 hover:bg-slate-100",
    secondary: "bg-slate-100 text-slate-900 hover:bg-slate-200",
    success: "bg-emerald-100 text-emerald-700 border border-emerald-200",
    warning: "bg-amber-100 text-amber-700 border border-amber-200",
    destructive: "bg-red-100 text-red-700 border border-red-200",
    purple: "bg-purple-100 text-purple-700 border border-purple-200"
  };
  return (
    <div className={cn("inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-slate-950 focus:ring-offset-2", variants[variant] || variants.default, className)}>
      {children}
    </div>
  );
};

export const Card = ({ children, className, ...props }: any) => (
  <div className={cn("rounded-xl border border-slate-200 bg-white text-slate-950 shadow-sm", className)} {...props}>
    {children}
  </div>
);

export const Button = ({ children, variant = 'primary', size = 'md', className, icon: Icon, isLoading, ...props }: any) => {
  const variants: Record<string, string> = {
    primary: "bg-indigo-600 text-white hover:bg-indigo-700 shadow-md shadow-indigo-200 active:scale-95",
    secondary: "bg-white text-slate-900 border border-slate-200 hover:bg-slate-50 shadow-sm",
    ghost: "hover:bg-slate-100 text-slate-700",
    danger: "bg-red-600 text-white hover:bg-red-700 shadow-sm",
    outline: "border border-slate-200 bg-transparent hover:bg-slate-100 text-slate-900",
    glass: "bg-white/20 backdrop-blur-md border border-white/30 text-white hover:bg-white/30"
  };
  const sizes: Record<string, string> = {
    sm: "h-8 px-3 text-xs",
    md: "h-10 px-4 py-2",
    lg: "h-11 px-8",
    icon: "h-10 w-10 p-2 flex items-center justify-center"
  };
  
  return (
    <button 
      className={cn(
        "inline-flex items-center justify-center whitespace-nowrap rounded-lg text-sm font-bold ring-offset-white transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 gap-2",
        variants[variant] || variants.primary,
        sizes[size] || sizes.md,
        className
      )}
      {...props}
    >
      {isLoading && <Loader2 className="h-4 w-4 animate-spin" />}
      {!isLoading && Icon && <Icon className="h-4 w-4" />}
      {children}
    </button>
  );
};

export const IconButton = ({ icon: Icon, onClick, className, variant = 'ghost', ...props }: any) => {
   const baseClass = "p-2 rounded-lg transition-all flex items-center justify-center active:scale-90";
   const variants: Record<string, string> = {
     ghost: "text-slate-400 hover:text-slate-700 hover:bg-slate-100",
     danger: "text-red-400 hover:text-red-700 hover:bg-red-50",
     primary: "text-indigo-600 bg-indigo-50 hover:bg-indigo-100",
     dark: "bg-slate-900 text-slate-400 hover:text-white"
   };
   return (
     <button onClick={onClick} className={cn(baseClass, variants[variant], className)} {...props}>
       <Icon size={16} />
     </button>
   );
};

export const Input = ({ className, ...props }: any) => (
  <input
    className={cn(
      "flex h-10 w-full rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm ring-offset-white file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-slate-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-600 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 transition-shadow",
      className
    )}
    {...props}
  />
);

export const Switch = ({ checked, onCheckedChange }: { checked: boolean, onCheckedChange: (v: boolean) => void }) => (
  <button 
    onClick={() => onCheckedChange(!checked)}
    className={cn("w-10 h-5 rounded-full relative transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500", checked ? "bg-indigo-600" : "bg-slate-300")}
  >
    <div className={cn("absolute top-1 w-3 h-3 rounded-full bg-white transition-all shadow-sm", checked ? "left-6" : "left-1")} />
  </button>
);

export const Tabs = ({ tabs, activeTab, onChange }: { tabs: string[], activeTab: string, onChange: (t: string) => void }) => (
  <div className="inline-flex h-10 items-center justify-center rounded-lg bg-slate-100 p-1 text-slate-500 w-full">
    {tabs.map(tab => (
      <button
        key={tab}
        onClick={() => onChange(tab)}
        className={cn(
          "inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1.5 text-xs font-bold uppercase tracking-wide ring-offset-white transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-950 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 flex-1",
          activeTab === tab ? "bg-white text-slate-950 shadow-sm" : "hover:bg-slate-200/50 hover:text-slate-900"
        )}
      >
        {tab}
      </button>
    ))}
  </div>
);

export const Tooltip = ({ text, children }: { text: string, children?: React.ReactNode }) => (
  <div className="group relative inline-block">
    {children}
    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-2 py-1 bg-slate-900 text-white text-[10px] font-bold rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none z-50">
      {text}
      <div className="absolute top-full left-1/2 -translate-x-1/2 border-4 border-transparent border-t-slate-900" />
    </div>
  </div>
);
