import React, { useMemo, useState } from 'react';
import { motion } from 'framer-motion';
import { Playbook, WorldData, JourneyNodeData, Locale, getLocalizedText } from '../types';
import { WorldContainer } from './WorldContainer';
import { Rocket, Sparkles, FileText, UploadCloud, Check, Map as MapIcon, Calendar } from 'lucide-react';
import { cn } from '../lib/utils';

interface JourneyMapProps {
  playbook: Playbook;
  currentActiveNodeId?: string;
  locale: Locale;
}

export const JourneyMap: React.FC<JourneyMapProps> = ({ playbook, currentActiveNodeId, locale }) => {
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
        <div className="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-0 sm:p-4 bg-slate-950/60 backdrop-blur-sm" onClick={() => setSelectedNode(null)}>
          <motion.div 
            initial={{ y: "100%" }}
            animate={{ y: 0 }}
            exit={{ y: "100%" }}
            className="bg-white w-full max-w-lg rounded-t-3xl sm:rounded-3xl shadow-2xl overflow-hidden border-t border-white/20"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Modal Header */}
            <div className="bg-slate-50 p-6 border-b border-slate-100">
              <div className="flex items-start gap-4">
                <div className="p-4 bg-white shadow-lg rounded-2xl text-primary-600 border border-slate-100">
                  <FileText size={28} />
                </div>
                <div>
                  <div className="flex items-center gap-2 mb-1">
                     <span className="text-[10px] font-black uppercase tracking-widest bg-slate-200 text-slate-600 px-2 py-0.5 rounded">
                       {selectedNode.type}
                     </span>
                     <span className={cn(
                       "text-[10px] font-black uppercase tracking-widest px-2 py-0.5 rounded text-white",
                       selectedNode.state === 'done' ? 'bg-emerald-500' : 'bg-primary-500'
                     )}>
                       {selectedNode.state.replace('_', ' ')}
                     </span>
                  </div>
                  <h3 className="text-xl font-bold text-slate-900 leading-tight">
                    {getLocalizedText(selectedNode.title, locale)}
                  </h3>
                </div>
              </div>
            </div>

            {/* Modal Content */}
            <div className="p-6">
              <div className="prose prose-slate prose-sm max-w-none text-slate-600 bg-slate-50/50 p-4 rounded-2xl border border-slate-100 mb-6">
                <p>{getLocalizedText(selectedNode.description, locale) || "No detailed description available for this step."}</p>
              </div>

              {selectedNode.requirements && (
                 <div className="mb-8">
                    <h4 className="text-xs font-black text-slate-400 uppercase tracking-widest mb-3">Quest Requirements</h4>
                    <ul className="space-y-3">
                       {selectedNode.requirements.uploads?.map((u: any, i: number) => (
                         <li key={i} className="flex items-center gap-3 text-sm text-slate-700 p-3 bg-blue-50/50 border border-blue-100 rounded-xl">
                           <div className="bg-blue-100 text-blue-600 p-1.5 rounded-lg">
                             <UploadCloud size={16} />
                           </div>
                           <span className="font-medium">{getLocalizedText(u.label, locale)}</span>
                         </li>
                       ))}
                       {selectedNode.requirements.fields?.map((f: any, i: number) => (
                         <li key={i} className="flex items-center gap-3 text-sm text-slate-700 p-3 bg-slate-50 border border-slate-200 rounded-xl">
                           <div className="bg-slate-200 text-slate-500 p-1.5 rounded-lg">
                             <Check size={16} />
                           </div>
                           <span className="font-medium">{getLocalizedText(f.label, locale)}</span>
                         </li>
                       ))}
                    </ul>
                 </div>
              )}

              <button 
                onClick={() => setSelectedNode(null)}
                className="w-full py-4 bg-slate-900 hover:bg-slate-800 text-white font-bold text-lg rounded-2xl transition-all shadow-lg shadow-slate-900/20 active:scale-[0.98]"
              >
                Close Mission
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  );
};