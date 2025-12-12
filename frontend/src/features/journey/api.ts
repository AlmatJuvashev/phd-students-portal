import { api } from "@/api/client";

export interface ScoreboardEntry {
  user_id: string;
  name: string;
  avatar: string;
  score: number;
  rank: number;
}

export interface ScoreboardResponse {
  top_5: ScoreboardEntry[];
  average_score: number;
  me: ScoreboardEntry | null;
  total_users: number;
}

export const journeyApi = {
  getScoreboard: () => api<ScoreboardResponse>("/journey/scoreboard"),
};
