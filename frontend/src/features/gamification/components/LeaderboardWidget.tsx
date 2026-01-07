import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Trophy, Medal, Crown } from 'lucide-react';
import { getLeaderboard } from '../api';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Skeleton } from '@/components/ui/skeleton';

export const LeaderboardWidget: React.FC = () => {
  const { data: leaders, isLoading } = useQuery({
    queryKey: ['gamification', 'leaderboard'],
    queryFn: () => getLeaderboard(5),
  });

  if (isLoading) {
    return (
        <div className="space-y-3">
             {[1,2,3].map(i => <Skeleton key={i} className="h-12 w-full rounded-xl" />)}
        </div>
    );
  }

  const getRankIcon = (index: number) => {
    switch (index) {
        case 0: return <Crown className="w-4 h-4 text-amber-400 fill-amber-400" />;
        case 1: return <Medal className="w-4 h-4 text-slate-300 fill-slate-300" />;
        case 2: return <Medal className="w-4 h-4 text-amber-700 fill-amber-700" />;
        default: return <span className="text-sm font-bold text-slate-400 w-4 text-center">{index + 1}</span>;
    }
  };

  return (
    <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden">
      <div className="p-4 border-b border-slate-100 bg-slate-50/50">
        <h3 className="font-bold text-slate-900 flex items-center gap-2">
            <Trophy className="w-4 h-4 text-indigo-500" /> Leaderboard
        </h3>
      </div>
      <div className="divide-y divide-slate-50">
        {leaders?.map((entry, index) => (
            <div key={entry.user_id} className="p-3 flex items-center gap-3 hover:bg-slate-50 transition-colors">
                <div className="flex-shrink-0 w-6 flex justify-center">
                    {getRankIcon(index)}
                </div>
                <Avatar className="h-8 w-8 border border-slate-200">
                    <AvatarImage src={entry.avatar_url} />
                    <AvatarFallback className="text-[10px] font-bold bg-indigo-50 text-indigo-600">
                        {entry.first_name?.[0]}{entry.last_name?.[0]}
                    </AvatarFallback>
                </Avatar>
                <div className="flex-1 min-w-0">
                    <div className="text-sm font-bold text-slate-700 truncate">
                        {entry.first_name} {entry.last_name}
                    </div>
                    <div className="text-[10px] text-slate-400 font-medium uppercase tracking-wider">
                        Level {entry.level}
                    </div>
                </div>
                <div className="text-sm font-bold text-indigo-600 font-mono">
                    {entry.total_xp}
                </div>
            </div>
        ))}
         {leaders?.length === 0 && (
            <div className="p-8 text-center text-slate-400 text-sm">No active users yet.</div>
        )}
      </div>
    </div>
  );
};
