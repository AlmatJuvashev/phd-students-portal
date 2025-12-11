import React from 'react';
import { motion } from 'framer-motion';
import { 
  Check, 
  Lock, 
  FileText, 
  Flag, 
  Shield, 
  UploadCloud, 
  AlertCircle,
  ChevronRight,
  Info,
  GraduationCap,
  Star,
  MousePointerClick
} from 'lucide-react';
import { JourneyNodeData, NodeState, NodeType, Locale, getLocalizedText } from '../types';
import { JourneyPath } from './JourneyPath';
import { cn } from '../lib/utils';

interface JourneyNodeProps {
  node: JourneyNodeData;
  index: number;
  isLast: boolean;
  locale: Locale;
  onClick: (node: JourneyNodeData) => void;
}

const getNodeIcon = (type: NodeType) => {
  switch (type) {
    case 'milestone': return Flag;
    case 'gateway': return Shield;
    case 'upload': return UploadCloud;
    case 'form': return FileText;
    case 'confirmTask': return Check;
    case 'info': return Info;
    case 'task': return Star;
    default: return FileText;
  }
};

const getStatusStyles = (state: NodeState) => {
  switch (state) {
    case 'done':
      return {
        wrapper: 'border-emerald-500/50 bg-gradient-to-br from-emerald-50 to-white shadow-sm',
        iconBg: 'bg-gradient-to-br from-emerald-400 to-emerald-600 shadow-emerald-200',
        iconColor: 'text-white',
        text: 'text-slate-600',
        title: 'text-slate-700',
        badge: 'bg-emerald-100 text-emerald-700'
      };
    case 'active':
      return {
        wrapper: 'border-primary-500 bg-white ring-4 ring-primary-100 shadow-lg scale-[1.02]',
        iconBg: 'bg-gradient-to-br from-primary-500 to-blue-600 shadow-blue-200',
        iconColor: 'text-white',
        text: 'text-slate-600',
        title: 'text-primary-800',
        badge: 'bg-blue-100 text-blue-700'
      };
    case 'needs_fixes':
      return {
        wrapper: 'border-amber-500 bg-amber-50',
        iconBg: 'bg-amber-500',
        iconColor: 'text-white',
        text: 'text-amber-700',
        title: 'text-amber-900',
        badge: 'bg-amber-100 text-amber-800'
      };
    case 'waiting':
      return {
        wrapper: 'border-indigo-200 bg-slate-50',
        iconBg: 'bg-indigo-100',
        iconColor: 'text-indigo-400',
        text: 'text-slate-500',
        title: 'text-slate-600',
        badge: 'bg-indigo-50 text-indigo-500'
      };
    case 'locked':
    default:
      return {
        wrapper: 'border-slate-100 bg-slate-50/50 opacity-80',
        iconBg: 'bg-slate-100',
        iconColor: 'text-slate-300',
        text: 'text-slate-400',
        title: 'text-slate-500',
        badge: 'hidden'
      };
  }
};

export const JourneyNode: React.FC<JourneyNodeProps> = ({ node, isLast, locale, onClick }) => {
  const styles = getStatusStyles(node.state);
  const Icon = node.state === 'locked' ? Lock : getNodeIcon(node.type);
  const isInteractive = node.state !== 'locked';

  return (
    <div className="relative flex group">
      {/* Path Line */}
      <JourneyPath state={node.state} isLast={isLast} />

      {/* Node Content */}
      <div className="flex w-full mb-8 relative z-10 pl-2 sm:pl-4">
        
        {/* Player Avatar (Only for Active Node) */}
        {node.state === 'active' && (
          <motion.div 
            initial={{ y: -10, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            className="absolute left-[1.5rem] sm:left-[2.25rem] -top-5 z-30 pointer-events-none"
          >
            <motion.div
              animate={{ y: [0, -8, 0] }}
              transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
              className="bg-white p-1 rounded-full shadow-lg border-2 border-primary-500"
            >
              <div className="bg-primary-500 text-white p-1.5 rounded-full">
                <GraduationCap size={18} />
              </div>
            </motion.div>
          </motion.div>
        )}

        <motion.button
          layout
          onClick={() => isInteractive && onClick(node)}
          disabled={!isInteractive}
          whileHover={isInteractive ? { scale: 1.02, x: 5 } : {}}
          whileTap={isInteractive ? { scale: 0.98 } : {}}
          className={cn(
            "flex w-full items-center text-left rounded-2xl p-4 sm:p-5 transition-all duration-200 border-2 shadow-sm relative overflow-hidden group/btn",
            styles.wrapper,
            isInteractive ? "cursor-pointer hover:shadow-md" : "cursor-not-allowed"
          )}
        >
          {/* Active Pulse Background */}
          {node.state === 'active' && (
            <motion.div 
              className="absolute inset-0 bg-primary-50/30"
              animate={{ opacity: [0.3, 0.6, 0.3] }}
              transition={{ duration: 2, repeat: Infinity }}
            />
          )}

          {/* Icon Wrapper */}
          <div className="relative mr-5 flex-shrink-0 z-10">
            {node.state === 'active' && (
              <span className="absolute -inset-3 rounded-full border-2 border-primary-400 opacity-20 animate-ping" />
            )}
            
            <div className={cn(
              "w-12 h-12 sm:w-14 sm:h-14 rounded-2xl rotate-3 flex items-center justify-center shadow-md transition-transform duration-300 group-hover:rotate-6",
              styles.iconBg,
              styles.iconColor
            )}>
              <Icon size={24} strokeWidth={2.5} />
            </div>

            {/* Status Badge Overlays */}
            {node.state === 'done' && (
              <div className="absolute -bottom-1 -right-1 bg-white text-emerald-500 rounded-full p-1 border border-emerald-100 shadow-sm">
                <Check size={12} strokeWidth={4} />
              </div>
            )}
            {node.state === 'needs_fixes' && (
              <div className="absolute -top-1 -right-1 bg-amber-500 text-white rounded-full p-1 border-2 border-white animate-pulse">
                <AlertCircle size={12} />
              </div>
            )}
          </div>

          {/* Text Content */}
          <div className="flex-1 min-w-0 z-10">
             <div className="flex justify-between items-center">
               <div className="pr-2">
                  <div className="flex items-center gap-2 mb-1">
                    <span className={cn("text-xs font-bold uppercase tracking-wider px-2 py-0.5 rounded-md", styles.badge)}>
                      {node.type}
                    </span>
                    {node.state === 'active' && (
                      <span className="flex h-2 w-2 relative">
                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary-400 opacity-75"></span>
                        <span className="relative inline-flex rounded-full h-2 w-2 bg-primary-500"></span>
                      </span>
                    )}
                  </div>
                  <h4 className={cn("text-base sm:text-lg font-bold leading-tight", styles.title)}>
                    {getLocalizedText(node.title, locale)}
                  </h4>
                  {node.description && (
                    <p className={cn("text-xs sm:text-sm mt-1 line-clamp-1 opacity-90", styles.text)}>
                      {getLocalizedText(node.description, locale)}
                    </p>
                  )}
               </div>
               
               {isInteractive && (
                 <div className="hidden sm:flex flex-row items-center gap-2 transition-all">
                   <span className="text-xs font-bold text-primary-600 opacity-0 -translate-x-3 group-hover/btn:opacity-100 group-hover/btn:translate-x-0 transition-all duration-300 ease-out whitespace-nowrap">
                     View Details
                   </span>
                   <div className="flex items-center justify-center w-8 h-8 rounded-full bg-slate-100 text-slate-400 group-hover/btn:bg-primary-500 group-hover/btn:text-white transition-all duration-300 shadow-sm group-hover/btn:shadow-md group-hover/btn:scale-110">
                     <ChevronRight size={18} className="transition-transform duration-300 group-hover/btn:translate-x-0.5" />
                   </div>
                 </div>
               )}
             </div>
          </div>
          
          {/* Subtle "Click" icon watermark for interactivity hint */}
          {isInteractive && (
             <div className="absolute top-2 right-2 text-primary-300 opacity-0 group-hover/btn:opacity-20 transition-opacity duration-300">
               <MousePointerClick size={24} />
             </div>
          )}

        </motion.button>
      </div>
    </div>
  );
};