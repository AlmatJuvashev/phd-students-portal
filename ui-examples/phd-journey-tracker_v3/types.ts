import { LucideIcon } from 'lucide-react';

export type Locale = 'ru' | 'kz' | 'en';

export type LocalizedText = {
  [key in Locale]?: string;
} | string;

export type NodeState = 'locked' | 'active' | 'waiting' | 'under_review' | 'needs_fixes' | 'done';
// Expanded node types based on JSON 'type' field and logical grouping
export type NodeType = 'task' | 'milestone' | 'gateway' | 'upload' | 'form' | 'confirmTask' | 'info';

export interface JourneyNodeData {
  id: string;
  module?: string; // e.g. "I", "II"
  title: LocalizedText;
  description?: LocalizedText;
  type: NodeType;
  // Computed state for UI
  state: NodeState; 
  // Raw data fields
  prerequisites?: string[];
  next?: string[];
  requirements?: any;
  screen?: any;
  meta?: any;
}

export interface WorldData {
  id: string;
  title: LocalizedText;
  description?: LocalizedText; // Mapped from first node or custom
  nodes: JourneyNodeData[];
  order: number;
  // Computed status for UI
  status: 'locked' | 'active' | 'completed';
  progress: number; // 0-100
  color?: string;
}

export interface Playbook {
  playbook_id: string;
  version: string;
  ui: {
    worlds_palette: string[];
    icons: Record<string, string>;
  };
  worlds: WorldData[];
}

// Helper to get text based on locale
export const getLocalizedText = (text: LocalizedText | undefined, locale: Locale): string => {
  if (!text) return '';
  if (typeof text === 'string') return text;
  return text[locale] || text['ru'] || Object.values(text)[0] || '';
};
