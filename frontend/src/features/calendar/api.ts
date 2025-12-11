import { CalendarEvent } from './types';

const MOCK_EVENTS: CalendarEvent[] = [
  {
    id: '1',
    title: 'Research Methodology Lecture',
    start: new Date(new Date().setHours(10, 0, 0, 0)),
    end: new Date(new Date().setHours(11, 30, 0, 0)),
    type: 'academic',
    location: 'Hall A, Building 4',
    description: 'Attendance is mandatory for all 1st year PhD students.'
  },
  {
    id: '2',
    title: 'Consultation with Prof. Ivanov',
    start: new Date(new Date().setDate(new Date().getDate() + 1)), // Tomorrow
    end: new Date(new Date().setDate(new Date().getDate() + 1)),
    type: 'academic',
    location: 'Office 302',
    description: 'Discuss chapter 2 corrections.'
  },
  {
    id: '3',
    title: 'Preliminary Exam',
    start: new Date(new Date().setDate(new Date().getDate() + 5)),
    end: new Date(new Date().setDate(new Date().getDate() + 5)),
    type: 'exam',
    location: 'Exam Center',
    description: 'Bring ID and registration slip.'
  },
  {
    id: '4',
    title: 'Doctoral Day Off',
    start: new Date(new Date().setDate(new Date().getDate() - 2)),
    end: new Date(new Date().setDate(new Date().getDate() - 2)),
    type: 'holiday',
    description: 'University closed.'
  },
  {
    id: '5',
    title: 'Gym Session',
    start: new Date(new Date().setHours(18, 0, 0, 0)),
    end: new Date(new Date().setHours(19, 30, 0, 0)),
    type: 'personal',
    location: 'Campus Gym'
  }
];

// Helper to simulate API delay
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export const fetchEvents = async (start: Date, end: Date): Promise<CalendarEvent[]> => {
  await delay(400);
  return MOCK_EVENTS;
};

export const createEvent = async (event: Omit<CalendarEvent, 'id'>): Promise<CalendarEvent> => {
  await delay(400);
  const newEvent = { ...event, id: Math.random().toString(36).substr(2, 9) };
  MOCK_EVENTS.push(newEvent);
  return newEvent;
};

export const updateEvent = async (updatedEvent: CalendarEvent): Promise<CalendarEvent> => {
  await delay(400);
  const index = MOCK_EVENTS.findIndex(e => e.id === updatedEvent.id);
  if (index !== -1) {
    MOCK_EVENTS[index] = updatedEvent;
  }
  return updatedEvent;
};

export const deleteEvent = async (id: string): Promise<boolean> => {
  await delay(300);
  const index = MOCK_EVENTS.findIndex(e => e.id === id);
  if (index !== -1) {
    MOCK_EVENTS.splice(index, 1);
  }
  return true;
};
