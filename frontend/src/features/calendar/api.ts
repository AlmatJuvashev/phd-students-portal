import { api } from "@/api/client";
import { CalendarEvent } from './types';

// Helper to map backend event to frontend CalendarEvent
const mapBackendEvent = (e: any): CalendarEvent => ({
  id: e.id,
  title: e.title,
  start: new Date(e.start_time),
  end: new Date(e.end_time),
  type: e.event_type,
  description: e.description,
  location: e.location,
  meeting_type: e.meeting_type,
  meeting_url: e.meeting_url,
  physical_address: e.physical_address,
  color: e.color,
  attendees: [], // TODO: map attendees if backend returns them in detailed view
});

export const fetchEvents = async (start: Date, end: Date): Promise<CalendarEvent[]> => {
  const startStr = start.toISOString();
  const endStr = end.toISOString();
  
  const events = await api(`/calendar/events?start=${startStr}&end=${endStr}`);
  return (events || []).map(mapBackendEvent);
};

export const createEvent = async (event: Omit<CalendarEvent, 'id'>): Promise<CalendarEvent> => {
  const payload = {
    title: event.title,
    description: event.description,
    start_time: event.start.toISOString(),
    end_time: event.end.toISOString(),
    event_type: event.type,
    location: event.location,
    meeting_type: event.meeting_type,
    meeting_url: event.meeting_url,
    physical_address: event.physical_address,
    color: event.color,
    attendees: event.attendees
  };

  const res = await api("/calendar/events", {
    method: "POST",
    body: JSON.stringify(payload),
  });
  
  return mapBackendEvent(res);
};

export const updateEvent = async (event: CalendarEvent): Promise<CalendarEvent> => {
  const payload = {
    title: event.title,
    description: event.description,
    start_time: event.start.toISOString(),
    end_time: event.end.toISOString(),
    event_type: event.type,
    location: event.location,
    meeting_type: event.meeting_type,
    meeting_url: event.meeting_url,
    physical_address: event.physical_address,
    color: event.color,
    attendees: event.attendees
  };

  await api(`/calendar/events/${event.id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
  
  return event;
};

export const deleteEvent = async (id: string): Promise<boolean> => {
  await api(`/calendar/events/${id}`, {
    method: "DELETE",
  });
  return true;
};
