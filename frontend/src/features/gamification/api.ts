import { api } from '@/api/client';

export type UserXP = {
  user_id: string;
  tenant_id: string;
  total_xp: number;
  level: number;
  current_streak: number;
  longest_streak: number;
  last_activity_date: string;
};

export type Badge = {
  id: string;
  code: string;
  name: string;
  description: string;
  icon_url: string;
  category: string;
  xp_reward: number;
  rarity: 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary';
  is_active: boolean;
};

export type UserBadge = {
  id: string;
  user_id: string;
  badge_id: string;
  earned_at: string;
  progress: number;
  badge_name: string;
  badge_icon: string;
  badge_desc: string;
};

export type LeaderboardEntry = {
  user_id: string;
  total_xp: number;
  level: number;
  first_name: string;
  last_name: string;
  avatar_url: string;
};

export const getMyStats = () => api.get<UserXP>('/gamification/stats');
export const getLeaderboard = (limit = 10) => api.get<LeaderboardEntry[]>(`/gamification/leaderboard?limit=${limit}`);
export const getMyBadges = () => api.get<UserBadge[]>('/gamification/badges/mine');
export const listAllBadges = () => api.get<Badge[]>('/gamification/badges');
