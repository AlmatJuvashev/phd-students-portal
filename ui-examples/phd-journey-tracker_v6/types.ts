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

// --- Chat Types ---

export type UserRole = 'student' | 'advisor' | 'admin' | 'secretary';

export interface ChatUser {
  id: string;
  name: string;
  avatar?: string; // Initials or URL
  role: UserRole;
  isOnline?: boolean;
}

export interface Attachment {
  id: string;
  type: 'image' | 'file';
  name: string;
  url: string;
  size?: string;
}

export interface Message {
  id: string;
  senderId: string;
  content: string;
  timestamp: string; // ISO string
  status: 'sending' | 'sent' | 'read';
  attachments?: Attachment[];
}

export interface ChatRoom {
  id: string;
  name: string;
  type: 'private' | 'group' | 'channel';
  participants: ChatUser[];
  lastMessage?: Message;
  unreadCount: number;
  isArchived?: boolean;
  messages: Message[]; // In a real app, this would be fetched separately
}

// --- Calendar Types ---

export type EventType = 'academic' | 'exam' | 'personal' | 'holiday';

export interface CalendarEvent {
  id: string;
  title: string;
  start: Date;
  end: Date;
  type: EventType;
  description?: string;
  location?: string;
}