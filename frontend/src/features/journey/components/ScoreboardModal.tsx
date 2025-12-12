import React, { useEffect, useState, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslation } from 'react-i18next';
import { X, Trophy, Medal, TrendingUp, Crown, Star, Flame, Zap, Award } from 'lucide-react';
import { cn } from '@/lib/utils';
import { journeyApi, ScoreboardResponse, ScoreboardEntry } from '../api';

interface ScoreboardModalProps {
  currentXP: number;
  onClose: () => void;
}

// Internal interface for UI usage
interface DisplayEntry extends ScoreboardEntry {
  isCurrentUser: boolean;
  isAverage: boolean;
  badges?: ('streak' | 'speed' | 'veteran')[];
}

export const ScoreboardModal: React.FC<ScoreboardModalProps> = ({ currentXP, onClose }) => {
  const { t } = useTranslation('common');
  const [data, setData] = useState<ScoreboardResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    journeyApi.getScoreboard()
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const { top3, restList, userRank, totalUsers } = useMemo(() => {
    if (!data) return { top3: [], restList: [], userRank: 0, totalUsers: 0 };

    // 1. Convert Top 5 and Me to DisplayEntry list
    // We only have Top 5 + Me from API.
    // Ideally user list would be longer for "The Chase" section.
    // For now we work with what we have.
    // Ensure uniqueness by ID
    const uniqueMap = new Map<string, DisplayEntry>();

    const processEntry = (e: ScoreboardEntry) => {
       const isMe = e.user_id === data.me?.user_id;
       // Mock badges for flavor
       const badges: ('streak' | 'speed' | 'veteran')[] = [];
       if (e.rank === 1) badges.push('veteran');
       if (e.rank === 2) badges.push('speed');
       if (e.rank === 3) badges.push('streak');
       if (isMe) badges.push('streak'); // Give user a streak badge for fun

       uniqueMap.set(e.user_id, {
           ...e,
           isCurrentUser: isMe,
           isAverage: false,
           badges
       });
    };

    data.top_5.forEach(processEntry);
    if (data.me) processEntry(data.me);

    let allStudents = Array.from(uniqueMap.values());
    allStudents.sort((a, b) => b.score - a.score);

    // 2. Insert Average
    const averageEntry: any = {
         user_id: 'avg',
         name: t('scoreboard.cohort_average'),
         score: data.average_score,
         isAverage: true,
         rank: -1,
         avatar: '',
         isCurrentUser: false
    };
    // Add average to the list for specific display logic
    // Usually average is shown in the list.
    allStudents.push(averageEntry);
    allStudents.sort((a, b) => b.score - a.score);

    // 3. Split
    // The "Podium" needs the actual Top 3 students (not average)
    // Filter out average for podium calculation to be safe
    const studentsOnly = allStudents.filter(s => !s.isAverage);
    const top3 = studentsOnly.slice(0, 3);
    
    // The rest list should contain everyone else + average, sorted.
    // Provide visually distinct list.
    // We exclude the top 3 (already on podium) from the list view?
    // v9 example: top3 are on podium, restStudents in list.
    const podiumIds = new Set(top3.map(s => s.user_id));
    const restList = allStudents.filter(s => !podiumIds.has(s.user_id));

    return { 
        top3, 
        restList, 
        userRank: data.me?.rank || 0,
        totalUsers: data.total_users 
    };
  }, [data, t]);

  const renderBadge = (type: string) => {
    switch(type) {
      case 'streak': return <div className="text-orange-500" title={t('scoreboard.badges.streak')}><Flame size={14} fill="currentColor" /></div>;
      case 'speed': return <div className="text-yellow-500" title={t('scoreboard.badges.speed')}><Zap size={14} fill="currentColor" /></div>;
      case 'veteran': return <div className="text-purple-500" title={t('scoreboard.badges.veteran')}><Award size={14} fill="currentColor" /></div>;
      default: return null;
    }
  };

  const PodiumStep = ({ student, rank }: { student: DisplayEntry, rank: number }) => {
    const isFirst = rank === 1;
    const isMe = student.isCurrentUser;
    
    // Initials generator
    const initials = student.name 
       ? student.name.split(' ').map((n: string) => n[0]).join('').substring(0,2).toUpperCase()
       : '??';
    
    let colors = {
      bg: "bg-slate-100",
      text: "text-slate-600",
      border: "border-slate-200",
      gradient: "from-slate-300 to-slate-400",
      shadow: "shadow-slate-500/20"
    };

    if (rank === 1) {
      colors = {
        bg: "bg-gradient-to-b from-yellow-100 to-white",
        text: "text-amber-700",
        border: "border-yellow-300",
        gradient: "from-amber-300 to-yellow-500",
        shadow: "shadow-amber-500/40"
      };
    } else if (rank === 2) {
      colors = {
        bg: "bg-gradient-to-b from-slate-100 to-white",
        text: "text-slate-700",
        border: "border-slate-300",
        gradient: "from-slate-300 to-slate-500",
        shadow: "shadow-slate-500/30"
      };
    } else if (rank === 3) {
      colors = {
        bg: "bg-gradient-to-b from-orange-50 to-white",
        text: "text-orange-800",
        border: "border-orange-300",
        gradient: "from-orange-400 to-amber-700",
        shadow: "shadow-orange-500/30"
      };
    }

    return (
      <div className={cn("flex flex-col items-center", isFirst ? "-mt-6 z-10" : "z-0")}>
        {/* Crown for #1 */}
        {isFirst && (
          <motion.div 
            initial={{ y: 10, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            className="mb-2"
          >
            <Crown size={32} className="text-amber-400 fill-amber-400 drop-shadow-md" />
          </motion.div>
        )}
        
        {/* Avatar */}
        <div className="relative mb-3">
           <div className={cn(
             "w-16 h-16 rounded-full flex items-center justify-center text-xl font-bold text-white shadow-xl bg-gradient-to-br ring-4 ring-white",
             colors.gradient,
             isMe && "ring-primary-400 scale-105"
           )}>
             {student.avatar || initials}
           </div>
           <div className={cn(
             "absolute -bottom-2 left-1/2 -translate-x-1/2 w-6 h-6 rounded-full bg-white flex items-center justify-center text-xs font-bold border-2 shadow-sm",
             colors.text,
             colors.border
           )}>
             {rank}
           </div>
        </div>

        {/* Name & Score */}
        <div className="text-center">
          <div className={cn("font-bold text-sm flex items-center justify-center gap-1", isMe ? "text-primary-700" : "text-slate-800")}>
             {student.name}
             {isMe && <span className="text-[9px] bg-primary-600 text-white px-1 rounded-sm">{t('scoreboard.you')}</span>}
          </div>
          <div className={cn("text-xs font-black", colors.text)}>
            {student.score.toLocaleString()} XP
          </div>
        </div>
        
        {/* Podium Base */}
        <motion.div 
          initial={{ height: 0 }}
          animate={{ height: isFirst ? 100 : rank === 2 ? 70 : 50 }}
          className={cn(
            "w-20 sm:w-24 mt-3 rounded-t-lg border-x border-t relative overflow-hidden",
            colors.bg,
            colors.border
          )}
        >
          <div className="absolute inset-0 bg-white/30 opacity-50" />
        </motion.div>
      </div>
    );
  };
  
  // Fallback XP if loading
  const displayXP = data?.me?.score ?? currentXP;

  return (
    <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" onClick={onClose}>
      <motion.div 
        initial={{ opacity: 0, scale: 0.9, y: 20 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        exit={{ opacity: 0, scale: 0.9 }}
        className="bg-white w-full max-w-lg rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[85vh] sm:max-h-[90vh]"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="relative bg-slate-900 text-white p-6 overflow-hidden flex-shrink-0">
          <div className="absolute top-0 right-0 w-64 h-64 bg-primary-500/20 rounded-full blur-3xl -mr-20 -mt-20" />
          <div className="absolute bottom-0 left-0 w-32 h-32 bg-amber-500/20 rounded-full blur-3xl -ml-10 -mb-10" />
          
          <div className="relative z-10 flex justify-between items-start mb-4">
            <div>
              <div className="flex items-center gap-2 mb-1 text-amber-400 font-bold uppercase tracking-wider text-xs">
                <Trophy size={14} /> {t('scoreboard.year')}
              </div>
              <h2 className="text-2xl font-black tracking-tight">{t('scoreboard.title')}</h2>
            </div>
            <button 
              onClick={onClose}
              className="p-2 bg-white/10 hover:bg-white/20 rounded-full transition-colors"
            >
              <X size={20} />
            </button>
          </div>

          {/* User Stats Mini-Summary */}
          <div className="relative z-10 flex items-center gap-4 bg-white/10 rounded-xl p-3 border border-white/5 backdrop-blur-md">
            <div className="flex-1">
               <div className="text-[10px] text-slate-400 uppercase font-bold mb-0.5">{t('scoreboard.your_position')}</div>
               <div className="flex items-baseline gap-1">
                 <span className="text-xl font-bold text-white">#{userRank > 0 ? userRank : '-'}</span>
                 <span className="text-xs text-slate-400">/ {totalUsers > 0 ? totalUsers : '?'}</span>
               </div>
            </div>
            <div className="w-px h-8 bg-white/10" />
            <div className="flex-1">
               <div className="text-[10px] text-slate-400 uppercase font-bold mb-0.5">{t('scoreboard.your_xp')}</div>
               <div className="text-xl font-bold text-white">{displayXP.toLocaleString()}</div>
            </div>
             <div className="w-px h-8 bg-white/10" />
             <div className="flex-1">
                <div className="text-[10px] text-slate-400 uppercase font-bold mb-0.5">{t('scoreboard.top_percent')}</div>
                <div className="text-xl font-bold text-white">
                    {userRank > 0 && totalUsers ? 
                        `Top ${Math.round((userRank / totalUsers) * 100)}%` 
                        : "-"
                    }
                </div>
             </div>
          </div>
        </div>

        {/* Content */}
        {loading ? (
             <div className="flex-1 flex items-center justify-center min-h-[300px] text-slate-400 animate-pulse">
                {t('scoreboard.loading')}
             </div>
        ) : (
        <div className="flex-1 overflow-y-auto custom-scrollbar bg-slate-50">
           
           {/* PODIUM Section (Only if we have top students) */}
           {top3.length > 0 && (
               <div className="bg-white pt-8 pb-4 px-4 shadow-sm border-b border-slate-200">
                 <div className="flex justify-center items-end gap-2 sm:gap-4">
                   {/* 2nd Place */}
                   {top3[1] ? <PodiumStep student={top3[1]} rank={2} /> : <div className="w-20 sm:w-24" />}
                   
                   {/* 1st Place */}
                   {top3[0] ? <PodiumStep student={top3[0]} rank={1} /> : <div className="w-20 sm:w-24" />}
                   
                   {/* 3rd Place */}
                   {top3[2] ? <PodiumStep student={top3[2]} rank={3} /> : <div className="w-20 sm:w-24" />}
                 </div>
               </div>
           )}

           {/* LIST Section */}
           <div className="p-4 space-y-2">
             <h3 className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3 px-2">{t('scoreboard.the_chase')}</h3>
             
             {restList.length === 0 && (
                 <div className="text-center text-sm text-slate-400 py-4">{t('scoreboard.no_contenders')}</div>
             )}

             {restList.map((entry) => {
               // Rank needs to be calculated or taken from entry.rank
               // For average, we don't have a rank.
               // For students, we have entry.rank.
               const rank = entry.rank > 0 ? entry.rank : '-';

               if (entry.isAverage) {
                 return (
                   <div key="avg" className="py-2 flex items-center gap-4 px-2 opacity-60">
                      <div className="w-6 text-center text-xs font-bold text-slate-400">-</div>
                      <div className="flex-1 flex items-center gap-3">
                        <div className="flex-1 border-b border-dashed border-slate-300 h-px" />
                        <div className="flex items-center gap-1.5 text-[10px] font-bold text-slate-500 uppercase px-2 py-0.5 bg-slate-200 rounded-full whitespace-nowrap">
                           <TrendingUp size={12} /> Avg: {entry.score}
                        </div>
                        <div className="flex-1 border-b border-dashed border-slate-300 h-px" />
                      </div>
                   </div>
                 );
               }

               const isMe = entry.isCurrentUser;
               
               // Initials
                const initials = entry.name 
                   ? entry.name.split(' ').map((n: string) => n[0]).join('').substring(0,2).toUpperCase()
                   : '??';

               return (
                 <motion.div
                   key={entry.user_id}
                   initial={{ opacity: 0, y: 10 }}
                   animate={{ opacity: 1, y: 0 }}
                   className={cn(
                     "flex items-center gap-3 p-3 rounded-xl transition-all border group",
                     isMe 
                       ? "bg-white border-primary-200 shadow-md ring-1 ring-primary-100 z-10" 
                       : "bg-white border-transparent hover:border-slate-200 hover:shadow-sm"
                   )}
                 >
                    <div className="w-6 text-center font-mono text-sm font-bold text-slate-400">
                      {rank}
                    </div>

                    <div className={cn(
                      "w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold",
                      isMe ? "bg-primary-100 text-primary-700" : "bg-slate-100 text-slate-500"
                    )}>
                      {entry.avatar || initials}
                    </div>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                         <span className={cn("font-medium text-sm truncate", isMe ? "text-slate-900" : "text-slate-600")}>
                           {entry.name}
                         </span>
                         <div className="flex items-center gap-1">
                            {entry.badges?.map((b, idx) => (
                              <div key={idx} className="opacity-0 group-hover:opacity-100 transition-opacity">
                                {renderBadge(b)}
                              </div>
                            ))}
                         </div>
                      </div>
                    </div>

                    <div className="text-right font-bold text-sm text-slate-600 tabular-nums">
                       {entry.score}
                    </div>
                 </motion.div>
               );
             })}
           </div>
        </div>
        )}
      </motion.div>
    </div>
  );
};
