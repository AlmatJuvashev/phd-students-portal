import React, { useState, useEffect, useMemo } from "react";
import { Calendar, dateFnsLocalizer, Views } from "react-big-calendar";
import { format, parse, startOfWeek, getDay } from "date-fns";
import enUS from "date-fns/locale/en-US";
import ru from "date-fns/locale/ru";
import kk from "date-fns/locale/kk";
import "react-big-calendar/lib/css/react-big-calendar.css";
import "./calendar-mobile.css";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { API_URL } from "@/api/client";
import { Plus } from "lucide-react";
import { EventModal } from "./EventModal";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "react-i18next";
import { BackButton } from "@/components/ui/back-button";

type View = "month" | "week" | "work_week" | "day" | "agenda";

const locales = {
  en: enUS,
  "en-US": enUS,
  ru: ru,
  kk: kk,
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
  const { t, i18n } = useTranslation("common");

  // Ensure culture is supported
  const culture = useMemo(() => {
    return locales[i18n.language as keyof typeof locales]
      ? i18n.language
      : "en-US";
  }, [i18n.language]);

  // Detect mobile and set default view
  const [isMobile, setIsMobile] = useState(false);
  const [defaultView, setDefaultView] = useState<View>("month");

  useEffect(() => {
    const checkMobile = () => {
      const mobile = window.innerWidth < 768; // md breakpoint
      setIsMobile(mobile);
      setDefaultView(mobile ? "day" : "month"); // Use day view for mobile instead of agenda
    };
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const fetchEvents = async () => {
    const start = new Date(new Date().getFullYear(), 0, 1).toISOString();
    const end = new Date(new Date().getFullYear() + 1, 0, 1).toISOString();
    const res = await fetch(`${API_URL}/events?start=${start}&end=${end}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok)
      throw new Error(
        t("calendar.errors.fetch", { defaultValue: "Failed to fetch events" })
      );
    return res.json();
  };

  const { data: events = [], isError } = useQuery<Event[]>({
    queryKey: ["events"],
    queryFn: fetchEvents,
    enabled: !!token,
    retry: 1,
  });

  const calendarEvents: CalendarEvent[] = useMemo(() => {
    return (events || [])
      .filter((e) => {
        if (!e || !e.start_time || !e.end_time) return false;
        const start = new Date(e.start_time);
        const end = new Date(e.end_time);
        return !isNaN(start.getTime()) && !isNaN(end.getTime());
      })
      .map((e) => ({
        id: e.id || crypto.randomUUID(),
        title: e.title || t("calendar.untitled", { defaultValue: "Untitled" }),
        start: new Date(e.start_time),
        end: new Date(e.end_time),
        resource: e,
      }));
  }, [events, t]);

  const handleSelectEvent = (event: CalendarEvent) => {
    setSelectedEvent(event.resource || null);
    setIsModalOpen(true);
  };

  const handleSelectSlot = ({ start, end }: { start: Date; end: Date }) => {
    setSelectedEvent(null);
    // Pre-fill dates if needed
    setIsModalOpen(true);
  };

  // Show error message if events fail to load
  if (isError) {
    return (
      <div className="min-h-screen p-4 md:p-6 max-w-6xl mx-auto">
        <div className="space-y-4">
          <div className="flex items-center gap-4">
            <BackButton to="/" />
            <h1 className="text-2xl md:text-3xl font-bold flex-1 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              {t("calendar.title", { defaultValue: "Calendar" })}
            </h1>
          </div>
          <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-6 text-center">
            <p className="text-destructive font-medium">
              {t("calendar.errors.fetch", {
                defaultValue: "Failed to load calendar events",
              })}
            </p>
            <p className="text-sm text-muted-foreground mt-2">
              {t("calendar.errors.check_connection", {
                defaultValue: "Please check your connection or try again later",
              })}
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-4 md:p-6 max-w-6xl mx-auto">
      <div className="space-y-4">
        {/* Header with Back Button and Title */}
        <div className="flex items-center gap-4">
          <BackButton to="/" />
          <h1 className="text-2xl md:text-3xl font-bold flex-1 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            {t("calendar.title", { defaultValue: "Calendar" })}
          </h1>
          <Button
            onClick={() => {
              setSelectedEvent(null);
              setIsModalOpen(true);
            }}
            className="shadow-lg hover:shadow-xl transition-all"
          >
            <Plus className="w-4 h-4 mr-2" />
            <span className="hidden sm:inline">
              {t("calendar.new_event", { defaultValue: "New Event" })}
            </span>
            <span className="sm:hidden">
              {t("calendar.new", { defaultValue: "New" })}
            </span>
          </Button>
        </div>

        {/* Calendar Container with improved styling */}
        <div className="bg-card border-2 border-border/50 rounded-2xl shadow-xl overflow-hidden">
          <div
            className="p-4 md:p-6"
            style={{ height: "calc(100vh - 200px)", minHeight: "600px" }}
          >
            <Calendar
              localizer={localizer}
              events={calendarEvents}
              startAccessor="start"
              endAccessor="end"
              style={{ height: "100%" }}
              onSelectEvent={handleSelectEvent}
              onSelectSlot={handleSelectSlot}
              selectable
              culture={culture}
              defaultView={defaultView}
              views={["month", "week", "day"]} // Agenda view is broken, so we exclude it
              messages={{
                // Navigation
                today: t("calendar.today", { defaultValue: "Today" }),
                previous: t("calendar.previous", { defaultValue: "Back" }),
                next: t("calendar.next", { defaultValue: "Next" }),

                // View names (toolbar buttons)
                month: t("calendar.month", { defaultValue: "Month" }),
                week: t("calendar.week", { defaultValue: "Week" }),
                day: t("calendar.day", { defaultValue: "Day" }),
                agenda: t("calendar.agenda", { defaultValue: "Agenda" }),

                // Agenda view columns
                date: t("calendar.date", { defaultValue: "Date" }),
                time: t("calendar.time", { defaultValue: "Time" }),
                event: t("calendar.event", { defaultValue: "Event" }),

                // Additional labels
                allDay: t("calendar.all_day", { defaultValue: "All Day" }),
                yesterday: t("calendar.yesterday", {
                  defaultValue: "Yesterday",
                }),
                tomorrow: t("calendar.tomorrow", { defaultValue: "Tomorrow" }),

                // Messages
                noEventsInRange: t("calendar.no_events", {
                  defaultValue: "No events in this range",
                }),
                showMore: (total) =>
                  t("calendar.show_more", {
                    defaultValue: `+${total} more`,
                    count: total,
                  }),
              }}
              className="custom-calendar"
            />
          </div>
        </div>
      </div>
      <EventModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        event={selectedEvent}
        onSuccess={() =>
          queryClient.invalidateQueries({ queryKey: ["events"] })
        }
      />
    </div>
  );
};
