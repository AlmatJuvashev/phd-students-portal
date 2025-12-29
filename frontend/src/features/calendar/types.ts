
export type EventType = 'academic' | 'exam' | 'personal' | 'holiday' | 'meeting' | 'deadline';

export interface CalendarEvent {
  id: string;
  title: string;
  start: Date;
  end: Date;
  type: EventType;
  description?: string;
  location?: string;
  meeting_type?: 'online' | 'offline';
  meeting_url?: string;
  physical_address?: string;
  color?: string;
  attendees?: string[]; // user IDs
}
