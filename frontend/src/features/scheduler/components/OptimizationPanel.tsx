import React from 'react';
import { Zap, Play, RefreshCw, Gauge } from 'lucide-react';
import { cn } from '@/lib/utils';
import { SolverConfig } from '../types';

interface OptimizationPanelProps {
  config: SolverConfig;
  setConfig: (c: SolverConfig) => void;
  onRun: () => void;
  isRunning: boolean;
  score: number;
}

export const OptimizationPanel: React.FC<OptimizationPanelProps> = ({ 
  config, 
  setConfig, 
  onRun, 
  isRunning, 
  score 
}) => {
  return (
    <div className="absolute bottom-6 right-6 w-80 bg-slate-900 text-white rounded-2xl shadow-2xl border border-slate-700 overflow-hidden z-50 animate-in slide-in-from-right-10 fade-in duration-500">
       <div className="p-4 border-b border-slate-700 bg-slate-800/50 flex justify-between items-center">
          <div className="flex items-center gap-2 text-sm font-bold">
             <Zap size={16} className="text-yellow-400 fill-yellow-400" /> Auto-Scheduler
          </div>
          <div className="text-[10px] font-mono text-slate-400">SOLVER_V2.1</div>
       </div>
       
       <div className="p-5 space-y-5">
          <div className="space-y-3">
             <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-slate-400">
                <span>Optimization Goal</span>
             </div>
             
             <div className="space-y-4">
                <div>
                   <div className="flex justify-between text-[10px] font-bold mb-1.5">
                      <span className="text-indigo-400">Space Efficiency</span>
                      <span>{config.utilization_weight}%</span>
                   </div>
                   <input 
                     type="range" 
                     value={config.utilization_weight}
                     onChange={(e) => setConfig({...config, utilization_weight: parseInt(e.target.value), satisfaction_weight: 100 - parseInt(e.target.value)})}
                     className="w-full h-1.5 bg-slate-700 rounded-full appearance-none cursor-pointer accent-indigo-500"
                   />
                </div>
                
                <div>
                   <div className="flex justify-between text-[10px] font-bold mb-1.5">
                      <span className="text-emerald-400">Professor Preferences</span>
                      <span>{config.satisfaction_weight}%</span>
                   </div>
                   <input 
                     type="range" 
                     value={config.satisfaction_weight}
                     onChange={(e) => setConfig({...config, satisfaction_weight: parseInt(e.target.value), utilization_weight: 100 - parseInt(e.target.value)})}
                     className="w-full h-1.5 bg-slate-700 rounded-full appearance-none cursor-pointer accent-emerald-500"
                   />
                </div>
 
                <div className="flex items-center justify-between pt-2">
                   <span className="text-[10px] text-slate-300">Minimize Building Hops</span>
                   <input 
                     type="checkbox" 
                     checked={config.prioritize_buildings}
                     onChange={(e) => setConfig({...config, prioritize_buildings: e.target.checked})}
                     className="rounded border-slate-600 bg-slate-700 text-indigo-500 focus:ring-indigo-500/50"
                   />
                </div>
             </div>
          </div>
 
          <div className="p-3 bg-slate-800 rounded-xl border border-slate-700 flex items-center justify-between">
             <div className="flex items-center gap-3">
                <Gauge size={20} className="text-slate-400" />
                <div>
                   <div className="text-[10px] text-slate-400 uppercase font-bold">Fitness Score</div>
                   <div className="text-lg font-black">{score.toFixed(1)}</div>
                </div>
             </div>
             <div className={cn("h-2 w-2 rounded-full", score > 80 ? "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)]" : "bg-amber-500")} />
          </div>
 
          <button 
            onClick={onRun}
            disabled={isRunning}
            className="w-full py-3 bg-indigo-600 hover:bg-indigo-500 text-white rounded-xl font-bold text-sm shadow-lg shadow-indigo-900/20 transition-all flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
             {isRunning ? <RefreshCw size={16} className="animate-spin" /> : <Play size={16} fill="currentColor" />}
             {isRunning ? 'Solving...' : 'Run Optimization'}
          </button>
       </div>
    </div>
  );
};
