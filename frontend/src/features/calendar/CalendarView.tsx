import React, { useState, useEffect } from "react";
import { Calendar, dateFnsLocalizer } from "react-big-calendar";
import { format, parse, startOfWeek, getDay } from "date-fns";
import { enUS } from "date-fns/locale/en-US";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { EventModal } from "./EventModal";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "react-i18next";

const locales = {
  'en-US': enUS,
};

const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek,
  getDay,
  locales,
});

interface Event {
  id: string;
  title: string;
  start_time: string;
  end_time: string;
  event_type: string;
  description?: string;
  location?: string;
  creator_id: string;
}

interface CalendarEvent {
  id: string;
  title: string;
  start: Date;
  end: Date;
  resource?: Event;
}

export const CalendarView = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<Event | null>(null);
  const { token } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation("common");

  const fetchEvents = async () => {
    const start = new Date(new Date().getFullYear(), 0, 1).toISOString();
    const end = new Date(new Date().getFullYear() + 1, 0, 1).toISOString();
    const res = await fetch(`${import.meta.env.VITE_API_BASE}/api/events?start=${start}&end=${end}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) throw new Error(t("calendar.errors.fetch", { defaultValue: "Failed to fetch events" }));
    return res.json();
  };

  const { data: events = [] } = useQuery<Event[]>({
    queryKey: ['events'],
    queryFn: fetchEvents,
    enabled: !!token,
  });

  const calendarEvents: CalendarEvent[] = events.map((e) => ({
    id: e.id,
    title: e.title,
    start: new Date(e.start_time),
    end: new Date(e.end_time),
    resource: e,
  }));

  const handleSelectEvent = (event: CalendarEvent) => {
    setSelectedEvent(event.resource || null);
    setIsModalOpen(true);
  };

  const handleSelectSlot = ({ start, end }: { start: Date; end: Date }) => {
    setSelectedEvent(null);
    // Pre-fill dates if needed
    setIsModalOpen(true);
  };

  return (
    <div className="h-screen p-4 flex flex-col gap-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">
          {t("calendar.title", { defaultValue: "Calendar" })}
        </h1>
        <Button onClick={() => { setSelectedEvent(null); setIsModalOpen(true); }}>
          <Plus className="w-4 h-4 mr-2" />
          {t("calendar.new_event", { defaultValue: "New Event" })}
        </Button>
      </div>
      <div className="flex-1 bg-white p-4 rounded-lg shadow">
        <Calendar
          localizer={localizer}
          events={calendarEvents}
          startAccessor="start"
          endAccessor="end"
          style={{ height: '100%' }}
          onSelectEvent={handleSelectEvent}
          onSelectSlot={handleSelectSlot}
          selectable
        />
      </div>
      <EventModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        event={selectedEvent}
        onSuccess={() => queryClient.invalidateQueries({ queryKey: ['events'] })}
      />
    </div>
  );
};
