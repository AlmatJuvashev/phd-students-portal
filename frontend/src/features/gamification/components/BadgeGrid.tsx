import React from 'react';
import { Badge } from '../api';
import { Hexagon } from 'lucide-react';

interface BadgeGridProps {
  badges: Badge[];
  earnedBadges: Set<string>; // Set of Badge IDs
}

export const BadgeGrid: React.FC<BadgeGridProps> = ({ badges, earnedBadges }) => {
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
      {badges.map(badge => {
        const isEarned = earnedBadges.has(badge.id);
        return (
            <div 
                key={badge.id}
                className={`relative group p-4 rounded-2xl border transition-all ${
                    isEarned 
                        ? 'bg-white border-indigo-100 shadow-sm hover:shadow-md hover:border-indigo-200' 
                        : 'bg-slate-50 border-slate-100 opacity-60 grayscale'
                }`}
            >
                <div className="flex flex-col items-center text-center space-y-2">
                    <div className={`w-12 h-12 rounded-full flex items-center justify-center mb-1 ${
                        isEarned ? 'bg-indigo-50 text-indigo-600' : 'bg-slate-200 text-slate-400'
                    }`}>
                        {badge.icon_url ? (
                             <img src={badge.icon_url} alt={badge.name} className="w-8 h-8 object-contain" />
                        ) : (
                             <Hexagon size={24} />
                        )}
                    </div>
                    <h4 className="font-bold text-sm text-slate-900 leading-tight">{badge.name}</h4>
                    <span className="text-[10px] uppercase font-bold tracking-wider text-slate-400">
                        {badge.xp_reward} XP
                    </span>
                </div>
                
                {/* Tooltip-ish description on hover */}
                <div className="absolute inset-0 bg-white/95 backdrop-blur-sm rounded-2xl p-4 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center text-center border border-indigo-100 shadow-lg pointer-events-none z-10">
                    <p className="text-xs text-slate-600 font-medium">{badge.description}</p>
                </div>
            </div>
        );
      })}
    </div>
  );
};
