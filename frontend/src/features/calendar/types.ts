
export type EventType = 'academic' | 'exam' | 'personal' | 'holiday' | 'meeting' | 'deadline';

export interface CalendarEvent {
  id: string;
  title: string;
  start: Date;
  end: Date;
  type: EventType;
  description?: string;
  location?: string;
}
