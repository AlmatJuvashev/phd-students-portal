import React, { useState, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslation } from 'react-i18next';
import { ChevronDown, Lock, Trophy, Star, Zap } from 'lucide-react';
import { NodeVM, t } from "@/lib/playbook";
import { JourneyNode } from './JourneyNode';
import { cn } from "@/lib/utils";

interface WorldData {
    id: string;
    title: string;
    nodes: NodeVM[];
    order?: number;
}

interface WorldContainerProps {
  world: WorldData;
  index: number;
  onNodeClick: (node: NodeVM) => void;
  isLocked?: boolean;
}

export const WorldContainer: React.FC<WorldContainerProps> = ({ world, index, onNodeClick, isLocked = false }) => {
    const { t: T } = useTranslation('common');
    // Calculate world stats
    const stats = useMemo(() => {
        const total = world.nodes.length;
        const done = world.nodes.filter(n => n.state === 'done').length;
        const active = world.nodes.some(n => n.state === 'active');
        const completed = done === total;
        const progress = total > 0 ? (done / total) * 100 : 0;
        return { total, done, active, completed, progress }; 
    }, [world.nodes]);

  const [isOpen, setIsOpen] = useState(stats.active || !isLocked);
  const [showXpTooltip, setShowXpTooltip] = useState(false);
  
  // Use a reliable color palette or passed in props. 
  const accentColor = '#3b82f6'; // Default primary blue

    // XP Calculation Stats
  // Logic mirrors JourneyMap: World 3 gives 0 XP, others give 100 per node
  // world.order might be missing, use index + 1
  const order = world.order || (index + 1);
  const isXpEligible = order !== 3;
  const currentXP = isXpEligible ? stats.done * 100 : 0;
  const potentialXP = isXpEligible ? stats.total * 100 : 0;

  return (
    <div className={cn(
      "relative mb-8 rounded-3xl overflow-hidden border-2 transition-all duration-500",
      stats.active 
        ? "bg-card border-primary/20 shadow-xl scale-[1.01]" 
        : "bg-muted/50 border-border opacity-95"
    )}>
      
      {/* Level Watermark */}
      <div className="absolute right-0 top-0 text-[10rem] font-black text-foreground opacity-[0.03] -z-0 pointer-events-none -mr-8 -mt-8 leading-none select-none">
        {index + 1}
      </div>

      {/* World Header (Level Banner) */}
      <button 
        onClick={() => !isLocked && setIsOpen(!isOpen)}
        className={cn(
          "w-full flex items-center justify-between p-6 relative z-10",
          isLocked ? "cursor-not-allowed" : "hover:bg-muted/20"
        )}
      >
        <div className="flex items-center gap-5">
          {/* Level Icon / Badge */}
          <div className="relative">
             <div 
               className={cn(
                 "flex items-center justify-center w-14 h-14 rounded-2xl text-white font-bold text-xl shadow-lg transform rotate-3 transition-transform group-hover:rotate-6",
                 isLocked ? "bg-slate-300 dark:bg-slate-700" : "bg-primary"
               )}
               // style={{ backgroundColor: isLocked ? undefined : accentColor }} 
             >
               {stats.completed ? <Trophy size={24} /> : isLocked ? <Lock size={24} /> : (index + 1)}
             </div>
             {stats.completed && (
               <div className="absolute -bottom-2 -right-2 bg-yellow-400 text-yellow-900 text-[0.6rem] font-bold px-2 py-0.5 rounded-full border border-white shadow-sm flex items-center gap-1">
                 <Star size={8} fill="currentColor" /> {T('map.done_suffix')}
               </div>
             )}
          </div>

          <div className="text-left">
            <span className={cn(
                "text-xs font-bold uppercase tracking-widest mb-1 block",
                isLocked ? "text-muted-foreground" : "text-muted-foreground"
              )}>
                 {T('world.level')} {index + 1}
            </span>
            <h3 className={cn(
                "text-xl sm:text-2xl font-black tracking-tight", 
                isLocked ? "text-muted-foreground" : "text-foreground"
              )}>
              {world.title}
            </h3>
            
            {/* XP / Progress Bar */}
            {!isLocked && (
              <div 
                className="mt-3 flex items-center gap-3 relative group/progress"
                 onMouseEnter={() => setShowXpTooltip(true)}
                 onMouseLeave={() => setShowXpTooltip(false)}
                 onClick={(e) => e.stopPropagation()} // Prevent collapse on bar click
              >
                 <div className="w-32 sm:w-48 h-2.5 bg-secondary rounded-full overflow-hidden shadow-inner cursor-help transition-transform group-hover/progress:scale-y-125 group-hover/progress:shadow-md">
                   <motion.div 
                     initial={{ width: 0 }}
                     animate={{ width: `${stats.progress}%` }}
                     transition={{ duration: 1, delay: 0.2 }}
                     className="h-full rounded-full relative overflow-hidden bg-primary"
                   >
                     {/* Gloss effect on bar */}
                     <div className="absolute top-0 left-0 w-full h-1/2 bg-white/30" />

                     {/* Shimmer effect when active */}
                     {stats.active && (
                       <motion.div 
                         className="absolute inset-0 bg-gradient-to-r from-transparent via-white/40 to-transparent w-1/2 -skew-x-12"
                         animate={{ x: ['-200%', '300%'] }}
                         transition={{ duration: 2, repeat: Infinity, repeatDelay: 1 }}
                       />
                     )}
                   </motion.div>
                 </div>
                 <span className="text-xs font-bold text-muted-foreground group-hover/progress:text-foreground transition-colors">
                   {Math.round(stats.progress)}% XP
                 </span>

                 {/* Tooltip */}
                 <AnimatePresence>
                   {showXpTooltip && (
                     <motion.div
                       initial={{ opacity: 0, y: 10, scale: 0.9 }}
                       animate={{ opacity: 1, y: 0, scale: 1 }}
                       exit={{ opacity: 0, y: 10, scale: 0.9 }}
                       className="absolute left-0 bottom-full mb-3 bg-slate-900/95 backdrop-blur-md text-white text-xs rounded-xl py-3 px-4 shadow-xl z-50 min-w-[160px] border border-slate-700 pointer-events-none"
                     >
                       <div className="flex flex-col gap-2">
                          <div className="flex items-center justify-between border-b border-slate-700/50 pb-2">
                             <span className="font-bold text-slate-200">{T('world.stats_title')}</span>
                             <div className="flex items-center gap-1 text-[10px] bg-slate-800 px-1.5 py-0.5 rounded text-slate-400">
                               <Zap size={10} className="text-yellow-400 fill-current" />
                               {isXpEligible ? T('world.ranked') : T('world.unranked')}
                             </div>
                          </div>
                          
                          <div className="space-y-1">
                             <div className="flex justify-between items-center text-[11px]">
                               <span className="text-slate-400 font-medium">{T('world.progress')}</span>
                               <span className={cn("font-bold", stats.progress === 100 ? "text-emerald-400" : "text-white")}>
                                 {Math.round(stats.progress)}%
                               </span>
                             </div>
                             
                             <div className="flex justify-between items-center text-[11px]">
                               <span className="text-slate-400 font-medium">{T('world.xp_earned')}</span>
                               <span className="text-amber-400 font-mono font-bold">
                                 +{currentXP.toLocaleString()} <span className="text-slate-600 font-normal">/ {potentialXP}</span>
                               </span>
                             </div>
                             
                             <div className="flex justify-between items-center text-[11px]">
                               <span className="text-slate-400 font-medium">{T('world.nodes')}</span>
                               <span className="text-slate-300 font-mono">
                                 {stats.done} / {stats.total}
                               </span>
                             </div>
                          </div>
                       </div>
                       
                       {/* Tooltip Arrow */}
                       <div className="absolute -bottom-1.5 left-6 w-3 h-3 bg-slate-900 border-r border-b border-slate-700 rotate-45 transform skew-x-12" />
                     </motion.div>
                   )}
                 </AnimatePresence>
              </div>
            )}
          </div>
        </div>

        <div className={cn(
            "text-muted-foreground transition-transform duration-300",
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
                 <div className="absolute left-[2.1rem] sm:left-[3.85rem] top-[-1.5rem] h-10 w-1.5 bg-slate-300/50 dark:bg-slate-700/50 rounded-full -z-10" />

                 {world.nodes.map((node, i) => (
                   <JourneyNode 
                     key={node.id} 
                     node={node} 
                     index={i} 
                     isLast={i === world.nodes.length - 1} 
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
