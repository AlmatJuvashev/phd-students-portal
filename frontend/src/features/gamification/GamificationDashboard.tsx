import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getMyStats, getMyBadges, listAllBadges } from './api';
import { XPProgressBar } from './components/XPProgressBar';
import { BadgeGrid } from './components/BadgeGrid';
import { LeaderboardWidget } from './components/LeaderboardWidget';
import { Loader2 } from 'lucide-react';

export const GamificationDashboard: React.FC = () => {
    
  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['gamification', 'stats'],
    queryFn: getMyStats,
  });

  const { data: myBadges, isLoading: myBadgesLoading } = useQuery({
    queryKey: ['gamification', 'badges', 'mine'],
    queryFn: getMyBadges,
  });

  const { data: allBadges, isLoading: allBadgesLoading } = useQuery({
    queryKey: ['gamification', 'badges', 'all'],
    queryFn: listAllBadges,
  });

  if (statsLoading || myBadgesLoading || allBadgesLoading) {
      return <div className="p-12 flex justify-center"><Loader2 className="animate-spin text-slate-400" /></div>;
  }

  const earnedBadgeIds = new Set(myBadges?.map(b => b.badge_id) || []);

  return (
    <div className="max-w-6xl mx-auto space-y-8 animate-in fade-in duration-500">
        <div>
          <h1 className="text-3xl font-black text-slate-900 tracking-tight">My Achievements</h1>
          <p className="text-slate-500 font-medium mt-1">Track your progress and badges.</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-8">
                {/* Stats & Progress */}
                {stats && <XPProgressBar stats={stats} />}
                
                {/* Badges Section */}
                <div>
                   <h3 className="text-lg font-bold text-slate-900 mb-4 flex items-center gap-2">
                      Badges Collection <span className="text-xs font-normal text-slate-400 bg-slate-100 px-2 py-1 rounded-full">{myBadges?.length || 0} / {allBadges?.length || 0}</span>
                   </h3>
                   {allBadges && <BadgeGrid badges={allBadges} earnedBadges={earnedBadgeIds} />}
                </div>
            </div>

            <div className="space-y-6">
                <LeaderboardWidget />
            </div>
        </div>
    </div>
  );
};
