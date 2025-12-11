import React, { useMemo, useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Playbook, WorldData, JourneyNodeData, Locale, getLocalizedText } from '../types';
import { WorldContainer } from './WorldContainer';
import { Rocket, Sparkles, FileText, UploadCloud, Check, Map as MapIcon, Calendar, X, ListChecks, ArrowRight, MessageCircle } from 'lucide-react';
import { cn } from '../lib/utils';

// --- Sub-component: Modal for Checklist Logic ---
const NodeDetailsModal = ({ 
  node, 
  locale, 
  onClose 
}: { 
  node: JourneyNodeData; 
  locale: Locale; 
  onClose: () => void 
}) => {
  // State for checklist items
  const [checkedItems, setCheckedItems] = useState<Record<string, boolean>>({});

  // Filter only boolean fields for progress calculation
  const booleanFields = useMemo(() => {
    return node.requirements?.fields?.filter((f: any) => f.type === 'boolean') || [];
  }, [node]);

  const totalObjectives = booleanFields.length;
  const completedObjectives = booleanFields.filter((f: any) => checkedItems[f.key]).length;
  const progressPercent = totalObjectives > 0 ? (completedObjectives / totalObjectives) * 100 : 0;
  const isAllComplete = totalObjectives > 0 && progressPercent === 100;

  const toggleItem = (key: string) => {
    setCheckedItems(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
  };

  return (
    <div className="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4 bg-slate-950/60 backdrop-blur-sm" onClick={onClose}>
      <motion.div 
        initial={{ y: "100%" }}
        animate={{ y: 0 }}
        exit={{ y: "100%" }}
        className="bg-white w-full max-w-xl rounded-t-3xl sm:rounded-3xl shadow-2xl overflow-hidden border-t border-white/20 flex flex-col max-h-[90vh]"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Modal Header */}
        <div className="bg-slate-50 p-6 border-b border-slate-100 flex-shrink-0">
          <div className="flex items-start gap-4">
            <div className={cn(
              "p-4 shadow-lg rounded-2xl border border-slate-100 transition-colors duration-500",
              isAllComplete ? "bg-emerald-100 text-emerald-600" : "bg-white text-primary-600"
            )}>
              {isAllComplete ? <Check size={28} strokeWidth={3} /> : <FileText size={28} />}
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                 <span className="text-[10px] font-black uppercase tracking-widest bg-slate-200 text-slate-600 px-2 py-0.5 rounded">
                   {node.type}
                 </span>
                 <span className={cn(
                   "text-[10px] font-black uppercase tracking-widest px-2 py-0.5 rounded text-white transition-colors",
                   (node.state === 'done' || isAllComplete) ? 'bg-emerald-500' : 'bg-primary-500'
                 )}>
                   {isAllComplete ? 'COMPLETED' : node.state.replace('_', ' ')}
                 </span>
              </div>
              <h3 className="text-xl font-bold text-slate-900 leading-tight">
                {getLocalizedText(node.title, locale)}
              </h3>
            </div>
            <button 
              onClick={onClose}
              className="p-2 bg-slate-200 hover:bg-slate-300 rounded-full text-slate-500 transition-colors"
            >
              <X size={20} />
            </button>
          </div>

          {/* Quest Progress Bar */}
          {totalObjectives > 0 && (
            <div className="mt-6">
              <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-slate-400 mb-2">
                <span>Quest Progress</span>
                <span className={cn(isAllComplete && "text-emerald-500")}>
                  {completedObjectives} / {totalObjectives}
                </span>
              </div>
              <div className="h-2 w-full bg-slate-200 rounded-full overflow-hidden">
                <motion.div 
                  initial={{ width: 0 }}
                  animate={{ width: `${progressPercent}%` }}
                  className={cn(
                    "h-full rounded-full transition-colors duration-500",
                    isAllComplete ? "bg-emerald-500" : "bg-primary-500"
                  )}
                />
              </div>
            </div>
          )}
        </div>

        {/* Scrollable Content */}
        <div className="p-6 overflow-y-auto custom-scrollbar">
          <div className="prose prose-slate prose-sm max-w-none text-slate-600 bg-slate-50/50 p-4 rounded-2xl border border-slate-100 mb-6">
            <p>{getLocalizedText(node.description, locale) || "No detailed description available for this step."}</p>
          </div>

          {/* Checklist Form */}
          {node.requirements?.fields && (
            <div className="mb-8">
              <div className="flex items-center gap-2 mb-4">
                <ListChecks size={18} className="text-slate-400" />
                <h4 className="text-xs font-black text-slate-400 uppercase tracking-widest">
                  Objectives Checklist
                </h4>
              </div>
              
              <div className="space-y-3">
                {node.requirements.fields.map((field: any, i: number) => {
                  if (field.type === 'note') {
                    return (
                      <div key={i} className="pt-4 pb-1 border-b border-slate-100 mb-2">
                         <h5 className="text-sm font-bold text-slate-800 flex items-center gap-2">
                           <div className="w-1.5 h-1.5 rounded-full bg-slate-400" />
                           {getLocalizedText(field.label, locale)}
                         </h5>
                      </div>
                    );
                  }

                  if (field.type === 'boolean') {
                    const isChecked = checkedItems[field.key] || false;
                    return (
                      <motion.button
                        key={field.key}
                        onClick={() => toggleItem(field.key)}
                        whileTap={{ scale: 0.98 }}
                        className={cn(
                          "w-full flex items-center gap-4 p-3 rounded-xl border text-left transition-all duration-200 group",
                          isChecked 
                            ? "bg-emerald-50 border-emerald-200 shadow-sm" 
                            : "bg-white border-slate-200 hover:border-primary-300 hover:shadow-md"
                        )}
                      >
                        <div className={cn(
                          "w-6 h-6 rounded-md border-2 flex items-center justify-center transition-colors duration-200 flex-shrink-0",
                          isChecked
                            ? "bg-emerald-500 border-emerald-500 text-white"
                            : "bg-white border-slate-300 group-hover:border-primary-400"
                        )}>
                          {isChecked && <Check size={16} strokeWidth={4} />}
                        </div>
                        <span className={cn(
                          "text-sm font-medium transition-colors",
                          isChecked ? "text-emerald-800 line-through opacity-70" : "text-slate-700"
                        )}>
                          {getLocalizedText(field.label, locale)}
                        </span>
                      </motion.button>
                    );
                  }
                  return null;
                })}
              </div>
            </div>
          )}

          {/* Uploads Section */}
          {node.requirements?.uploads && (
             <div className="mb-8">
                <h4 className="text-xs font-black text-slate-400 uppercase tracking-widest mb-3">Required Uploads</h4>
                <ul className="space-y-3">
                   {node.requirements.uploads.map((u: any, i: number) => (
                     <li key={i} className="flex items-center gap-3 text-sm text-slate-700 p-3 bg-blue-50/50 border border-blue-100 rounded-xl border-dashed">
                       <div className="bg-blue-100 text-blue-600 p-1.5 rounded-lg">
                         <UploadCloud size={16} />
                       </div>
                       <span className="font-medium">{getLocalizedText(u.label, locale)}</span>
                     </li>
                   ))}
                </ul>
             </div>
          )}
        </div>
        
        {/* Footer Actions */}
        <div className="p-4 bg-white border-t border-slate-100 flex-shrink-0">
          <button 
            onClick={onClose}
            className={cn(
              "w-full py-4 font-bold text-lg rounded-2xl transition-all shadow-lg active:scale-[0.98] flex items-center justify-center gap-2",
              isAllComplete 
                ? "bg-emerald-500 hover:bg-emerald-600 text-white shadow-emerald-500/20" 
                : "bg-slate-900 hover:bg-slate-800 text-white shadow-slate-900/20"
            )}
          >
            {isAllComplete ? (
              <>Complete Mission <Check size={20} /></>
            ) : (
              "Close Details"
            )}
          </button>
        </div>
      </motion.div>
    </div>
  );
};

interface JourneyMapProps {
  playbook: Playbook;
  currentActiveNodeId?: string;
  locale: Locale;
  onOpenChat: () => void;
  onOpenCalendar: () => void;
}

export const JourneyMap: React.FC<JourneyMapProps> = ({ playbook, currentActiveNodeId, locale, onOpenChat, onOpenCalendar }) => {
  // --- State Processing Logic ---
  const processedWorlds = useMemo(() => {
    let foundActive = false;
    let globalTotalNodes = 0;
    let globalCompletedNodes = 0;

    const newWorlds: WorldData[] = playbook.worlds.map((world, wIndex) => {
      const worldColor = playbook.ui.worlds_palette[wIndex % playbook.ui.worlds_palette.length];
      let worldCompletedNodes = 0;

      const newNodes = world.nodes.map((node) => {
        let state: JourneyNodeData['state'] = 'locked';

        if (!currentActiveNodeId && wIndex === 0 && node === world.nodes[0]) {
           state = 'active'; 
           foundActive = true;
        } else if (foundActive) {
          state = 'locked';
        } else if (node.id === currentActiveNodeId) {
          state = 'active';
          foundActive = true;
        } else {
          state = 'done';
          worldCompletedNodes++;
          globalCompletedNodes++;
        }

        return { ...node, state };
      });

      let worldStatus: WorldData['status'] = 'locked';
      if (newNodes.some(n => n.state === 'active')) {
        worldStatus = 'active';
      } else if (newNodes.every(n => n.state === 'done')) {
        worldStatus = 'completed';
      } else if (newNodes.some(n => n.state === 'done')) {
        worldStatus = 'active'; 
      }
      
      globalTotalNodes += world.nodes.length;

      return {
        ...world,
        nodes: newNodes,
        status: worldStatus,
        color: worldColor,
        progress: (worldCompletedNodes / world.nodes.length) * 100
      };
    });

    return { worlds: newWorlds, globalProgress: (globalCompletedNodes / globalTotalNodes) * 100 };
  }, [playbook, currentActiveNodeId]);

  const [selectedNode, setSelectedNode] = useState<JourneyNodeData | null>(null);

  return (
    <div className="max-w-3xl mx-auto pb-24 relative">
      {/* Decorative Background Element */}
      <div className="absolute top-0 left-4 w-1 h-full bg-slate-200/50 -z-50 dashed-line-pattern" />

      {/* Sticky HUD Header */}
      <div className="sticky top-4 z-40 mb-10 mx-2 sm:mx-0">
        <div className="bg-slate-900/95 backdrop-blur-md rounded-2xl shadow-2xl border border-slate-700 p-4 text-white">
          <div className="flex justify-between items-center mb-3">
            <div className="flex items-center gap-3">
              <div className="bg-gradient-to-br from-primary-500 to-indigo-600 p-2.5 rounded-xl shadow-lg shadow-primary-500/30">
                <Rocket size={20} className="text-white" />
              </div>
              <div>
                <h1 className="text-lg font-bold leading-none tracking-tight">PhD Adventure</h1>
                <p className="text-xs text-slate-400 mt-1 font-medium">KazNMU Student Portal</p>
              </div>
            </div>
            
            <div className="text-right">
              <div className="text-2xl font-black tabular-nums leading-none">
                {Math.round(processedWorlds.globalProgress)}%
              </div>
              <div className="text-[10px] uppercase tracking-wider text-slate-400 font-bold">Total Progress</div>
            </div>
          </div>
          
          {/* Global Progress Bar (HUD Style) */}
          <div className="relative h-3 w-full bg-slate-800 rounded-full overflow-hidden shadow-inner border border-slate-700">
            <motion.div 
              initial={{ width: 0 }}
              animate={{ width: `${processedWorlds.globalProgress}%` }}
              transition={{ duration: 1.5, ease: "circOut" }}
              className="h-full bg-gradient-to-r from-emerald-400 via-primary-400 to-indigo-400 relative"
            >
               {/* Animated gloss */}
               <motion.div 
                 animate={{ x: ["-100%", "200%"] }}
                 transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
                 className="absolute top-0 left-0 w-1/3 h-full bg-white/30 skew-x-12 blur-sm"
               />
            </motion.div>
          </div>
          
          {processedWorlds.globalProgress === 100 && (
            <div className="mt-3 bg-emerald-500/20 border border-emerald-500/30 text-emerald-300 px-3 py-1.5 rounded-lg text-xs font-bold flex items-center justify-center gap-2 animate-pulse">
              <Sparkles size={14} /> QUEST COMPLETE: DEGREE AWARDED
            </div>
          )}
        </div>
      </div>

      {/* World List */}
      <div className="space-y-6 px-2 sm:px-0">
        {processedWorlds.worlds.map((world, index) => (
          <WorldContainer 
            key={world.id}
            world={world}
            index={index}
            locale={locale}
            onNodeClick={setSelectedNode}
          />
        ))}
      </div>

      {/* Detail Modal */}
      {selectedNode && (
        <NodeDetailsModal 
          node={selectedNode} 
          locale={locale} 
          onClose={() => setSelectedNode(null)} 
        />
      )}

      {/* Floating Action Buttons */}
      <div className="fixed bottom-6 right-6 z-50 flex flex-col gap-3">
        {/* Calendar FAB */}
        <motion.button
          onClick={onOpenCalendar}
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
          className="p-3 bg-white text-primary-600 border border-primary-100 rounded-full shadow-lg hover:shadow-xl transition-all flex items-center justify-center"
          title="Open Calendar"
        >
          <Calendar size={24} />
        </motion.button>

        {/* Chat FAB */}
        <motion.button
          onClick={onOpenChat}
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
          className="p-4 bg-primary-600 text-white rounded-full shadow-2xl shadow-primary-500/40 hover:bg-primary-500 transition-colors flex items-center justify-center relative"
          title="Open Chat"
        >
          <MessageCircle size={28} fill="currentColor" className="text-white" />
          <span className="absolute -top-1 -right-1 flex h-4 w-4">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-4 w-4 bg-red-500 text-[9px] font-bold items-center justify-center">2</span>
          </span>
        </motion.button>
      </div>
    </div>
  );
};