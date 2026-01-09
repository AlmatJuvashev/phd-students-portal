import React, { useState } from 'react';
import { Menu, X, Layout, Settings, GraduationCap, Map, PenTool, Cpu, FileQuestion, BookOpen } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useNavigate } from 'react-router-dom';

// Note: In the root Main.tsx, we can't use useNavigate easily outside RouterProvider.
// But we can accept a generic 'navigate' function or use window.location if outside router context.
// Better: We wrap this component inside the Root Layout or make a wrapper in Main.tsx? 
// No, Main.tsx has RouterProvider. We can't put this *outside* RouterProvider if we want to use useNavigate inside.
// However, the router object itself can navigate. 

export const DevNavigation = ({ onNavigate }: { onNavigate: (path: string) => void }) => {
  const [isOpen, setIsOpen] = useState(false);

  // Close when clicking a link
  const nav = (path: string) => {
    onNavigate(path);
    setIsOpen(false);
  };

  return (
    <div className="fixed bottom-6 left-6 z-[9999] font-sans">
      <div className="relative">
        {isOpen && (
          <div className="absolute bottom-full left-0 mb-4 w-64 bg-slate-900 text-white rounded-2xl shadow-2xl overflow-hidden border border-slate-700 p-2 animate-in slide-in-from-bottom-5 fade-in origin-bottom-left">
            <div className="px-3 py-2 text-[10px] font-bold text-slate-500 uppercase tracking-wider flex justify-between items-center">
              <span>Quick Access</span>
              <span className="text-[9px] bg-amber-500/20 text-amber-500 px-1.5 py-0.5 rounded border border-amber-500/30">DEV MODE</span>
            </div>
            
            <div className="space-y-1">
              <button onClick={() => nav('/student/dashboard')} className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group">
                <GraduationCap size={16} className="text-emerald-400" />
                <div className="flex-1">
                  <div className="text-xs font-bold text-slate-200">Student Portal</div>
                </div>
              </button>

              <button onClick={() => nav('/teach/dashboard')} className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group">
                <PenTool size={16} className="text-amber-400" />
                <div className="flex-1">
                  <div className="text-xs font-bold text-slate-200">Instructor Portal</div>
                </div>
              </button>

              <button onClick={() => nav('/admin')} className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group">
                <Settings size={16} className="text-blue-400" />
                <div className="flex-1">
                  <div className="text-xs font-bold text-slate-200">Admin Dashboard</div>
                </div>
              </button>

              <div className="h-px bg-slate-800 my-1 mx-2" />
              <div className="px-3 py-1 text-[9px] font-bold text-slate-600 uppercase">New Features</div>

               <button onClick={() => nav('/admin/item-banks')} className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group">
                <FileQuestion size={16} className="text-rose-400" />
                <div className="flex-1">
                  <div className="text-xs font-bold text-slate-200">Item Banks</div>
                  <div className="text-[9px] text-slate-500">Assessment Polish</div>
                </div>
              </button>

              <button onClick={() => nav('/admin/programs')} className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group">
                <BookOpen size={16} className="text-purple-400" />
                <div className="flex-1">
                  <div className="text-xs font-bold text-slate-200">Curriculum Studio</div>
                  <div className="text-[9px] text-slate-500">Course & Program Builders</div>
                </div>
              </button>

              <button onClick={() => nav('/admin/scheduler')} className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group">
                <Cpu size={16} className="text-indigo-400" />
                <div className="flex-1">
                  <div className="text-xs font-bold text-slate-200">Scheduler AI</div>
                </div>
              </button>
              
               <div className="h-px bg-slate-800 my-1 mx-2" />
              
              <button 
                onClick={() => nav('/journey')}
                className="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-white/10 transition-colors text-left group"
              >
                <Map size={16} className="text-slate-500 group-hover:text-slate-300" />
                <div className="text-xs font-bold text-slate-500 group-hover:text-slate-300">Legacy Journey</div>
              </button>
            </div>
          </div>
        )}
        
        <button 
          onClick={() => setIsOpen(!isOpen)}
          className={cn(
            "h-12 w-12 rounded-full flex items-center justify-center shadow-2xl transition-all hover:scale-110 active:scale-95 border-4 z-[100]",
            isOpen ? "bg-slate-900 text-white border-slate-800" : "bg-white text-slate-900 border-slate-100 hover:border-indigo-100"
          )}
          title="Dev Navigation"
        >
          {isOpen ? <X size={20} /> : <Menu size={20} />}
        </button>
      </div>
    </div>
  );
};
