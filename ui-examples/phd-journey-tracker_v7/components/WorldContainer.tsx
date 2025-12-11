import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronDown, Lock, Trophy, Star } from 'lucide-react';
import { WorldData, JourneyNodeData, Locale, getLocalizedText } from '../types';
import { JourneyNode } from './JourneyNode';
import { cn } from '../lib/utils';

interface WorldContainerProps {
  world: WorldData;
  index: number;
  locale: Locale;
  onNodeClick: (node: JourneyNodeData) => void;
}

export const WorldContainer: React.FC<WorldContainerProps> = ({ world, index, locale, onNodeClick }) => {
  const [isOpen, setIsOpen] = useState(world.status === 'active');
  const accentColor = world.color || '#3b82f6';
  const isLocked = world.status === 'locked';
  const isCompleted = world.status === 'completed';

  return (
    <div className={cn(
      "relative mb-8 rounded-3xl overflow-hidden border-2 transition-all duration-500",
      world.status === 'active' 
        ? "bg-white border-primary-200 shadow-xl scale-[1.01]" 
        : "bg-slate-50 border-slate-200 opacity-95"
    )}>
      
      {/* Level Watermark */}
      <div className="absolute right-0 top-0 text-[10rem] font-black text-slate-100 opacity-50 -z-0 pointer-events-none -mr-8 -mt-8 leading-none select-none">
        {world.order}
      </div>

      {/* World Header (Level Banner) */}
      <button 
        onClick={() => !isLocked && setIsOpen(!isOpen)}
        className={cn(
          "w-full flex items-center justify-between p-6 relative z-10",
          isLocked ? "cursor-not-allowed" : "hover:bg-slate-50/50"
        )}
      >
        <div className="flex items-center gap-5">
          {/* Level Icon / Badge */}
          <div className="relative">
             <div 
               className={cn(
                 "flex items-center justify-center w-14 h-14 rounded-2xl text-white font-bold text-xl shadow-lg transform rotate-3 transition-transform group-hover:rotate-6",
                 isLocked ? "bg-slate-300" : ""
               )}
               style={{ backgroundColor: isLocked ? undefined : accentColor }}
             >
               {isCompleted ? <Trophy size={24} /> : isLocked ? <Lock size={24} /> : world.order}
             </div>
             {isCompleted && (
               <div className="absolute -bottom-2 -right-2 bg-yellow-400 text-yellow-900 text-[0.6rem] font-bold px-2 py-0.5 rounded-full border border-white shadow-sm flex items-center gap-1">
                 <Star size={8} fill="currentColor" /> CLEARED
               </div>
             )}
          </div>

          <div className="text-left">
            <span className={cn(
                "text-xs font-bold uppercase tracking-widest mb-1 block",
                isLocked ? "text-slate-400" : "text-slate-500"
              )}>
                Level {world.order}
            </span>
            <h3 className={cn(
                "text-xl sm:text-2xl font-black tracking-tight", 
                isLocked ? "text-slate-400" : "text-slate-800"
              )}>
              {getLocalizedText(world.title, locale)}
            </h3>
            
            {/* XP / Progress Bar */}
            {!isLocked && (
              <div className="mt-3 flex items-center gap-3">
                 <div className="w-32 sm:w-48 h-2.5 bg-slate-200 rounded-full overflow-hidden shadow-inner">
                   <motion.div 
                     initial={{ width: 0 }}
                     animate={{ width: `${world.progress}%` }}
                     transition={{ duration: 1, delay: 0.2 }}
                     className="h-full rounded-full relative overflow-hidden"
                     style={{ backgroundColor: accentColor }}
                   >
                     {/* Gloss effect on bar */}
                     <div className="absolute top-0 left-0 w-full h-1/2 bg-white/30" />
                   </motion.div>
                 </div>
                 <span className="text-xs font-bold text-slate-400">
                   {Math.round(world.progress)}% XP
                 </span>
              </div>
            )}
          </div>
        </div>

        <div className={cn(
            "text-slate-300 transition-transform duration-300",
            isOpen ? "rotate-180" : "rotate-0"
          )}>
          {!isLocked && <ChevronDown size={28} strokeWidth={3} />}
        </div>
      </button>

      {/* Collapsible Content */}
      <AnimatePresence>
        {isOpen && !isLocked && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.4, ease: "easeInOut" }}
          >
            <div className="px-4 sm:px-8 pb-8 pt-2 relative z-10">
              <div className="relative pl-0 sm:pl-4">
                 {/* Visual connector from Header to first node */}
                 <div className="absolute left-[2.1rem] sm:left-[3.85rem] top-[-1.5rem] h-10 w-1.5 bg-slate-300/50 rounded-full -z-10" />

                 {world.nodes.map((node, i) => (
                   <JourneyNode 
                     key={node.id} 
                     node={node} 
                     index={i} 
                     isLast={i === world.nodes.length - 1} 
                     locale={locale}
                     onClick={onNodeClick}
                   />
                 ))}
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};