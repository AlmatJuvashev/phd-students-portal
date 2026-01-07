import React from 'react';
import { Progress } from '@/components/ui/progress';
import { Trophy, Star } from 'lucide-react';
import { UserXP } from '../api';

interface XPProgressBarProps {
  stats: UserXP;
}

// Simple level config for frontend calculation (keep in sync with backend or fetch)
const LEVELS = [
  { level: 1, xp: 0 },
  { level: 2, xp: 100 },
  { level: 3, xp: 300 },
  { level: 4, xp: 600 },
  { level: 5, xp: 1000 },
  { level: 6, xp: 1500 },
  { level: 7, xp: 2100 },
  { level: 8, xp: 2800 },
  { level: 9, xp: 3600 },
  { level: 10, xp: 4500 },
];

export const XPProgressBar: React.FC<XPProgressBarProps> = ({ stats }) => {
  const currentLevel = LEVELS.find(l => l.level === stats.level) || LEVELS[0];
  const nextLevel = LEVELS.find(l => l.level === stats.level + 1);
  
  let progress = 100;
  let nextXP = stats.total_xp; // Default if max level
  
  if (nextLevel) {
    const range = nextLevel.xp - currentLevel.xp;
    const current = stats.total_xp - currentLevel.xp;
    progress = Math.min(100, Math.max(0, (current / range) * 100));
    nextXP = nextLevel.xp;
  }

  return (
    <div className="bg-white p-4 rounded-xl border border-slate-100 shadow-sm">
      <div className="flex justify-between items-center mb-2">
        <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600">
                <Trophy size={16} />
            </div>
            <div>
                <div className="text-xs font-bold text-slate-400 uppercase">Level {stats.level}</div>
                <div className="font-bold text-slate-700">Learner</div>
            </div>
        </div>
        <div className="text-right">
             <div className="text-xs font-bold text-slate-400 uppercase">Total XP</div>
             <div className="font-black text-indigo-600 text-lg flex items-center gap-1 justify-end">
                {stats.total_xp} <Star size={12} className="fill-indigo-600" />
             </div>
        </div>
      </div>
      
      <div className="space-y-1">
        <Progress value={progress} className="h-2 bg-slate-100" />
        {nextLevel && (
            <div className="flex justify-between text-[10px] font-bold text-slate-400 uppercase tracking-wider">
                <span>{stats.total_xp} XP</span>
                <span>{nextLevel.xp} XP to Level {nextLevel.level}</span>
            </div>
        )}
      </div>
    </div>
  );
};
