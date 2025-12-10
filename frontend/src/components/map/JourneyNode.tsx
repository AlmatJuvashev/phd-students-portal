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
  MousePointerClick,
  Users,
  GitMerge,
  Clock,
  ExternalLink,
  LucideIcon
} from 'lucide-react';
import { NodeVM, NodeState, t } from "@/lib/playbook";
import { JourneyPath } from './JourneyPath';
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";

interface JourneyNodeProps {
  node: NodeVM;
  index: number;
  isLast: boolean;
  onClick: (node: NodeVM) => void;
}

const getNodeIcon = (node: NodeVM): LucideIcon => {
    // Map specific IDs if needed, similar to original NodeToken
     const idMap: Record<string, LucideIcon> = {
        S1_text_ready: FileText,
        S1_antiplag: Shield,
        E3_hearing_nk: Users,
    };
    if (idMap[node.id]) return idMap[node.id];

  switch (node.type) {
    case 'decision': return GitMerge;
    case 'meeting': return Users;
    case 'waiting': return Clock;
    case 'external': return ExternalLink;
    case 'boss': return Star; // Was Trophy
    case 'gateway': return Shield;
    case 'upload': return UploadCloud;
    case 'form': return FileText;
    case 'confirmTask': return Check;
    case 'info': return Info;
    default: return FileText;
  }
};

const getStatusStyles = (state: NodeState) => {
  switch (state) {
    case 'done': // emerald
      return {
        wrapper: 'border-emerald-500/50 bg-gradient-to-br from-emerald-50 to-white shadow-sm',
        iconBg: 'bg-gradient-to-br from-emerald-400 to-emerald-600 shadow-emerald-200',
        iconColor: 'text-white',
        text: 'text-slate-600 dark:text-slate-300',
        title: 'text-slate-700 dark:text-slate-200',
        badge: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900 dark:text-emerald-300'
      };
    case 'submitted': // amber/orange (waiting for review)
        return {
            wrapper: 'border-orange-200 bg-orange-50 dark:bg-orange-950/10',
            iconBg: 'bg-orange-100 dark:bg-orange-900',
            iconColor: 'text-orange-600 dark:text-orange-400',
            text: 'text-slate-500',
            title: 'text-slate-700 dark:text-slate-200',
            badge: 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300'
        };
    case 'active': // primary blue
      return {
        wrapper: 'border-primary bg-white dark:bg-slate-900 ring-4 ring-primary/20 shadow-lg scale-[1.02]',
        iconBg: 'bg-gradient-to-br from-primary to-blue-600 shadow-blue-200',
        iconColor: 'text-primary-foreground',
        text: 'text-slate-600 dark:text-slate-300',
        title: 'text-primary dark:text-primary-foreground',
        badge: 'bg-primary/20 text-primary'
      };
    case 'needs_fixes': // red
      return {
        wrapper: 'border-destructive/50 bg-destructive/5',
        iconBg: 'bg-destructive',
        iconColor: 'text-destructive-foreground',
        text: 'text-destructive/80',
        title: 'text-destructive',
        badge: 'bg-destructive/20 text-destructive'
      };
    case 'waiting': // indigo / slate
      return {
        wrapper: 'border-indigo-200 bg-slate-50 dark:bg-slate-900/50',
        iconBg: 'bg-indigo-100 dark:bg-indigo-900',
        iconColor: 'text-indigo-400 dark:text-indigo-300',
        text: 'text-slate-500 dark:text-slate-400',
        title: 'text-slate-600 dark:text-slate-300',
        badge: 'bg-indigo-50 text-indigo-500 dark:bg-indigo-900/50 dark:text-indigo-300'
      };
    case 'locked':
    default:
      return {
        wrapper: 'border-slate-100 bg-slate-50/50 opacity-80 dark:bg-slate-900/20 dark:border-slate-800',
        iconBg: 'bg-slate-100 dark:bg-slate-800',
        iconColor: 'text-slate-300 dark:text-slate-600',
        text: 'text-slate-400 dark:text-slate-600',
        title: 'text-slate-500 dark:text-slate-500',
        badge: 'hidden'
      };
  }
};

export const JourneyNode: React.FC<JourneyNodeProps> = ({ node, isLast, index, onClick }) => {
  const styles = getStatusStyles(node.state);
  const Icon = node.state === 'locked' ? Lock : getNodeIcon(node);
  const isInteractive = node.state !== 'locked';
  
  // Title mapping using our t() helper
  const title = t(node.title, node.id);

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
              className="bg-background p-1 rounded-full shadow-lg border-2 border-primary"
            >
              <div className="bg-primary text-primary-foreground p-1.5 rounded-full">
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
              className="absolute inset-0 bg-primary/5"
              animate={{ opacity: [0.3, 0.6, 0.3] }}
              transition={{ duration: 2, repeat: Infinity }}
            />
          )}

          {/* Icon Wrapper */}
          <div className="relative mr-5 flex-shrink-0 z-10">
            {node.state === 'active' && (
              <span className="absolute -inset-3 rounded-full border-2 border-primary/40 opacity-20 animate-ping" />
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
              <div className="absolute -bottom-1 -right-1 bg-background text-emerald-500 rounded-full p-1 border border-emerald-100 shadow-sm">
                <Check size={12} strokeWidth={4} />
              </div>
            )}
            {node.state === 'needs_fixes' && (
              <div className="absolute -top-1 -right-1 bg-destructive text-destructive-foreground rounded-full p-1 border-2 border-white animate-pulse">
                <AlertCircle size={12} />
              </div>
            )}
          </div>

          {/* Text Content */}
          <div className="flex-1 min-w-0 z-10">
             <div className="flex justify-between items-center text-foreground">
               <div className="pr-2">
                  <div className="flex items-center gap-2 mb-1">
                    <span className={cn("text-xs font-bold uppercase tracking-wider px-2 py-0.5 rounded-md", styles.badge)}>
                      {node.type}
                    </span>
                    {node.state === 'active' && (
                      <span className="flex h-2 w-2 relative">
                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                        <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                      </span>
                    )}
                    {['submitted', 'needs_fixes'].includes(node.state || '') && (
                         <span className={cn("text-xs font-bold uppercase tracking-wider px-2 py-0.5 rounded-md ml-2", 
                             node.state === 'submitted' ? "bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300" : 
                             "bg-destructive/20 text-destructive"
                         )}>
                           {node.state === 'submitted' ? t({ru: "На проверке", kz: "Тексерілуде", en: "Under Review"}, "Pending") : 
                            t({ru: "Требует исправления", kz: "Түзету қажет", en: "Needs Fixes"}, "Fixes")}
                         </span>
                    )}
                  </div>
                  <h4 className={cn("text-base sm:text-lg font-bold leading-tight", styles.title)}>
                    {title}
                  </h4>
                  {/* We don't really have description in NodeVM yet, but if we did: */}
                  {/* {node.description && (
                    <p className={cn("text-xs sm:text-sm mt-1 line-clamp-1 opacity-90", styles.text)}>
                      {t(node.description)}
                    </p>
                  )} */}
               </div>
               
               {isInteractive && (
                 <div className="hidden sm:flex flex-row items-center gap-2 transition-all">
                   <span className="text-xs font-bold text-primary opacity-0 -translate-x-3 group-hover/btn:opacity-100 group-hover/btn:translate-x-0 transition-all duration-300 ease-out whitespace-nowrap">
                     Details
                   </span>
                   <div className="flex items-center justify-center w-8 h-8 rounded-full bg-slate-100 text-slate-400 group-hover/btn:bg-primary group-hover/btn:text-primary-foreground transition-all duration-300 shadow-sm group-hover/btn:shadow-md group-hover/btn:scale-110 dark:bg-slate-800">
                     <ChevronRight size={18} className="transition-transform duration-300 group-hover/btn:translate-x-0.5" />
                   </div>
                 </div>
               )}
             </div>
          </div>
          
          {/* Subtle "Click" icon watermark for interactivity hint */}
          {isInteractive && (
             <div className="absolute top-2 right-2 text-primary/30 opacity-0 group-hover/btn:opacity-20 transition-opacity duration-300">
               <MousePointerClick size={24} />
             </div>
          )}

        </motion.button>
      </div>
    </div>
  );
};
