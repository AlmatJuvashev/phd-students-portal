import { CalendarEvent } from '../../types';
import { addDays, addWeeks, addMonths, isBefore, isSameDay } from 'date-fns';

const MOCK_EVENTS: CalendarEvent[] = [
  {
    id: '1',
    title: 'Research Methodology Lecture',
    start: new Date(new Date().setHours(10, 0, 0, 0)),
    end: new Date(new Date().setHours(11, 30, 0, 0)),
    type: 'academic',
    location: 'Hall A, Building 4',
    description: 'Attendance is mandatory for all 1st year PhD students.',
    recurrence: 'weekly',
    recurrenceEnd: new Date(new Date().setMonth(new Date().getMonth() + 2))
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
    description: 'University closed.',
    recurrence: 'weekly'
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

const expandRecurringEvents = (events: CalendarEvent[], startRange: Date, endRange: Date): CalendarEvent[] => {
  const expanded: CalendarEvent[] = [];
  
  events.forEach(event => {
    // If not recurring, check if it roughly falls in range or just add it (simple mock logic)
    if (!event.recurrence || event.recurrence === 'none') {
      expanded.push(event);
      return;
    }

    // Recurrence Expansion Logic
    let currentStart = new Date(event.start);
    let currentEnd = new Date(event.end);
    const stopDate = event.recurrenceEnd ? new Date(event.recurrenceEnd) : endRange;
    // Cap strictly at end of fetched range + buffer or the recurrence end, whichever is sooner
    const effectiveStopDate = isBefore(stopDate, endRange) ? stopDate : endRange;

    // We start generating from the event start. 
    // In a real optimized backend, we'd skip to `startRange`.
    let safetyCounter = 0;
    while (isBefore(currentStart, addDays(effectiveStopDate, 1)) && safetyCounter < 500) {
       // Only add if it's within the requested window (or slightly before to handle overlaps)
       if (!isBefore(currentStart, startRange) || !isBefore(addDays(currentEnd, 1), startRange)) {
           expanded.push({
             ...event,
             id: `${event.id}_${currentStart.getTime()}`, // Composite ID for instance
             start: new Date(currentStart),
             end: new Date(currentEnd),
             recurrence: 'none' // Instance is not itself the rule container for UI purposes
           });
       }

       // Advance date
       if (event.recurrence === 'daily') {
         currentStart = addDays(currentStart, 1);
         currentEnd = addDays(currentEnd, 1);
       } else if (event.recurrence === 'weekly') {
         currentStart = addWeeks(currentStart, 1);
         currentEnd = addWeeks(currentEnd, 1);
       } else if (event.recurrence === 'monthly') {
         currentStart = addMonths(currentStart, 1);
         currentEnd = addMonths(currentEnd, 1);
       }
       safetyCounter++;
    }
  });

  return expanded;
};

export const fetchEvents = async (start: Date, end: Date): Promise<CalendarEvent[]> => {
  await delay(300);
  // Expand events based on the requested window. 
  // If no window provided (defaults), use a generous window around now.
  const windowStart = start || addMonths(new Date(), -6);
  const windowEnd = end || addMonths(new Date(), 6);
  
  return expandRecurringEvents(MOCK_EVENTS, windowStart, windowEnd);
};

export const createEvent = async (event: Omit<CalendarEvent, 'id'>): Promise<CalendarEvent> => {
  await delay(300);
  const newEvent = { ...event, id: Math.random().toString(36).substr(2, 9) };
  MOCK_EVENTS.push(newEvent);
  return newEvent;
};

export const updateEvent = async (updatedEvent: CalendarEvent): Promise<CalendarEvent> => {
  await delay(300);
  
  // Logic to handle "Edit Series"
  // If ID has an underscore, it's an instance. We extract the base ID.
  const baseId = updatedEvent.id.includes('_') ? updatedEvent.id.split('_')[0] : updatedEvent.id;
  
  const index = MOCK_EVENTS.findIndex(e => e.id === baseId);
  if (index !== -1) {
    // We update the master record. 
    // Note: This resets the start/end to the edited instance's times if we aren't careful.
    // In a simple "Edit Series" from an instance, we usually want to keep the original start time 
    // but update title/desc/recurrence.
    // HOWEVER, for this mock, we will assume the user is editing the 'definition' of the event.
    // If they changed the time on an instance, and saved as series, the series shifts.
    
    // We need to preserve the ID as the base ID
    MOCK_EVENTS[index] = { ...updatedEvent, id: baseId };
  }
  return updatedEvent;
};

export const deleteEvent = async (id: string): Promise<boolean> => {
  await delay(300);
  const baseId = id.includes('_') ? id.split('_')[0] : id;
  
  const index = MOCK_EVENTS.findIndex(e => e.id === baseId);
  if (index !== -1) {
    MOCK_EVENTS.splice(index, 1);
  }
  return true;
};