
import React from 'react';
import { Loader2, Cloud, CheckCircle2, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';

export type AutosaveStatus = 'saved' | 'saving' | 'unsaved' | 'error';

interface AutosaveIndicatorProps {
  status: AutosaveStatus;
  lastSaved?: Date | null;
  className?: string;
}

export const AutosaveIndicator: React.FC<AutosaveIndicatorProps> = ({ status, lastSaved, className }) => {
  return (
    <div className={cn("flex items-center gap-2 text-xs font-medium transition-colors select-none", className)}>
      {status === 'saving' && (
        <>
          <Loader2 size={14} className="animate-spin text-slate-400" />
          <span className="text-slate-400">Saving...</span>
        </>
      )}
      
      {status === 'saved' && (
        <div className="flex items-center gap-1.5 text-slate-400 group relative cursor-help">
          <Cloud size={14} className="text-emerald-500/80" />
          <span>Saved</span>
          {lastSaved && (
            <span className="opacity-0 group-hover:opacity-100 transition-opacity absolute top-full mt-1 left-1/2 -translate-x-1/2 bg-slate-800 text-white px-2 py-1 rounded text-[10px] whitespace-nowrap z-50">
              Last saved at {lastSaved.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
            </span>
          )}
        </div>
      )}

      {status === 'unsaved' && (
        <>
          <div className="w-2 h-2 rounded-full bg-amber-400 ring-2 ring-amber-100" />
          <span className="text-amber-600">Unsaved changes</span>
        </>
      )}

      {status === 'error' && (
        <>
          <AlertCircle size={14} className="text-red-500" />
          <span className="text-red-600">Failed to save</span>
        </>
      )}
    </div>
  );
};
