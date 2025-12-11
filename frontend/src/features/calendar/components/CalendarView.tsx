import React, { useState, useEffect, useMemo } from 'react';
import { format, addMonths, subMonths, addWeeks, subWeeks, addDays, subDays, startOfMonth, endOfMonth } from 'date-fns';
import { enUS, ru, kk } from 'date-fns/locale';
import { motion, AnimatePresence } from 'framer-motion';
import { CalendarEvent, EventType } from '../types';
import { ChevronLeft, ChevronRight, Calendar as CalendarIcon, Clock, Plus, List, Filter, Bell } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';
import { EventDialog } from './EventDialog';
import { fetchEvents, createEvent, updateEvent, deleteEvent } from '../api';
import { MonthView } from './MonthView';
import { WeekDayView } from './WeekDayView';
import { AgendaView } from './AgendaView';

const LOCALES: Record<string, any> = {
  en: enUS,
  ru: ru,
  kz: kk,
};

type ViewType = 'month' | 'week' | 'day' | 'agenda';

export const CalendarView: React.FC = () => {
  const { t, i18n } = useTranslation();
  const currentLocale = LOCALES[i18n.language] || enUS;

  const [currentDate, setCurrentDate] = useState(new Date());
  const [view, setView] = useState<ViewType>('month');
  const [events, setEvents] = useState<CalendarEvent[]>([]);
  const [selectedEvent, setSelectedEvent] = useState<Partial<CalendarEvent> | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  
  // New State for Features
  const [filterType, setFilterType] = useState<EventType | 'all'>('all');
  const [notifiedEvents, setNotifiedEvents] = useState<Set<string>>(new Set());
  const [permission, setPermission] = useState(
    typeof Notification !== 'undefined' ? Notification.permission : 'default'
  );

  // Responsive: Default to Agenda view on mobile
  useEffect(() => {
    if (window.innerWidth < 640) {
      setView('agenda');
    }
  }, []);

  // Fetch events when current date or view changes
  useEffect(() => {
    loadEvents();
  }, [currentDate, view]);

  // Notification Check Interval
  useEffect(() => {
    const checkUpcomingEvents = () => {
      if (!('Notification' in window) || Notification.permission !== 'granted') return;

      const now = new Date();
      events.forEach(event => {
        if (notifiedEvents.has(event.id)) return;

        const startTime = new Date(event.start);
        const diffMs = startTime.getTime() - now.getTime();
        const diffMins = diffMs / (1000 * 60);

        // Notify if the event starts in 15 minutes or less (but hasn't started yet)
        if (diffMins > 0 && diffMins <= 15) {
          try {
            new Notification(t('calendar.notifications.upcoming', { title: event.title }) as string, {
              body: `${format(startTime, 'HH:mm')} - ${event.location || ''}`,
              icon: '/favicon.ico', // Use a default icon
              tag: event.id // Prevent duplicate notifications for same event
            });
            
            setNotifiedEvents(prev => {
              const next = new Set(prev);
              next.add(event.id);
              return next;
            });
          } catch (error) {
            console.error("Notification error:", error);
          }
        }
      });
    };

    // Check immediately and then every minute
    checkUpcomingEvents();
    const intervalId = setInterval(checkUpcomingEvents, 60000);

    return () => clearInterval(intervalId);
  }, [events, notifiedEvents, t]);

  const requestNotificationPermission = async () => {
    if (!('Notification' in window)) return;
    try {
      const result = await Notification.requestPermission();
      setPermission(result);
      if (result === 'granted') {
         new Notification(t('calendar.notifications.enabled_title', 'Notifications Enabled') as string, {
           body: t('calendar.notifications.enabled_msg', 'You will now be notified 15 minutes before your events.') as string,
           icon: '/favicon.ico'
         });
      }
    } catch (e) {
      console.error(e);
    }
  };

  const loadEvents = async () => {
    setLoading(true);
    // Determine fetch range based on view + buffer to handle transitions smoothly
    const rangeStart = subMonths(currentDate, 1);
    const rangeEnd = addMonths(currentDate, 2);
    
    const data = await fetchEvents(rangeStart, rangeEnd);
    setEvents(data);
    setLoading(false);
  };

  const handleNavigate = (direction: 'prev' | 'next' | 'today') => {
    if (direction === 'today') {
      setCurrentDate(new Date());
      return;
    }

    const modifier = direction === 'next' ? 1 : -1;

    switch (view) {
      case 'month':
      case 'agenda':
        setCurrentDate(prev => addMonths(prev, modifier));
        break;
      case 'week':
        setCurrentDate(prev => addWeeks(prev, modifier));
        break;
      case 'day':
        setCurrentDate(prev => addDays(prev, modifier));
        break;
    }
  };

  const handleSelectSlot = (date: Date) => {
    // Default to 1 hour duration
    const end = new Date(date.getTime() + 60 * 60 * 1000);
    setSelectedEvent({ start: date, end, type: 'academic' });
    setIsDialogOpen(true);
  };

  const handleSelectEvent = (event: CalendarEvent) => {
    setSelectedEvent(event);
    setIsDialogOpen(true);
  };

  const handleEventDrop = async (event: CalendarEvent, newStart: Date) => {
    // Calculate duration to preserve it
    const duration = event.end.getTime() - event.start.getTime();
    const newEnd = new Date(newStart.getTime() + duration);

    const updatedEvent = { ...event, start: newStart, end: newEnd };
    
    // Optimistic update
    setEvents(prev => prev.map(e => e.id === event.id ? updatedEvent : e));
    
    await updateEvent(updatedEvent);
    // Reload to ensure consistency (especially for recurring series logic in backend)
    await loadEvents();
  };

  const handleSaveEvent = async (eventData: CalendarEvent | Omit<CalendarEvent, 'id'>) => {
    if ('id' in eventData) {
      await updateEvent(eventData as CalendarEvent);
    } else {
      await createEvent(eventData);
    }
    // Reload to reflect updates/expansions
    await loadEvents();
    setIsDialogOpen(false);
    setSelectedEvent(null);
  };

  const handleDeleteEvent = async (id: string) => {
    await deleteEvent(id);
    await loadEvents();
    setIsDialogOpen(false);
    setSelectedEvent(null);
  };

  // Filtered Events Logic
  const filteredEvents = useMemo(() => {
    if (filterType === 'all') return events;
    return events.filter(e => e.type === filterType);
  }, [events, filterType]);

  return (
    <div className="h-[calc(100vh-140px)] sm:h-[700px] flex flex-col gap-4">
      {/* Calendar Toolbar */}
      <div className="flex flex-col xl:flex-row items-center justify-between gap-4 bg-white p-3 rounded-2xl shadow-sm border border-slate-200 sticky top-0 z-20">
        
        {/* Navigation & Title Group */}
        <div className="flex items-center justify-between w-full xl:w-auto gap-4">
          <div className="flex items-center gap-2 bg-slate-100 rounded-xl p-1">
            <button onClick={() => handleNavigate('prev')} className="p-2 hover:bg-white hover:shadow-sm rounded-lg transition-all text-slate-500">
              <ChevronLeft size={20} />
            </button>
            <button onClick={() => handleNavigate('today')} className="px-4 py-2 text-xs font-bold text-slate-600 hover:text-slate-900 transition-colors uppercase tracking-wider">
              {t('calendar.today', 'Today')}
            </button>
            <button onClick={() => handleNavigate('next')} className="p-2 hover:bg-white hover:shadow-sm rounded-lg transition-all text-slate-500">
              <ChevronRight size={20} />
            </button>
          </div>

          <h2 className="text-lg sm:text-xl font-bold text-slate-800 tabular-nums text-center xl:hidden capitalize">
            {format(currentDate, 'MMMM yyyy', { locale: currentLocale })}
          </h2>
        </div>

        {/* Center Title (Desktop) */}
        <h2 className="hidden xl:block text-xl font-bold text-slate-800 tabular-nums text-center absolute left-1/2 -translate-x-1/2 capitalize">
          {format(currentDate, 'MMMM yyyy', { locale: currentLocale })}
        </h2>

        {/* Controls Group */}
        <div className="flex flex-col sm:flex-row items-center gap-3 w-full xl:w-auto">
            
            {/* Filter Dropdown */}
            <div className="relative w-full sm:w-auto">
              <Filter className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" size={16} />
              <select
                value={filterType}
                onChange={(e) => setFilterType(e.target.value as EventType | 'all')}
                className="w-full sm:w-40 pl-9 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm font-medium text-slate-700 outline-none focus:ring-2 focus:ring-primary/20 appearance-none cursor-pointer hover:bg-slate-100 transition-colors"
              >
                <option value="all">{t('calendar.filters.all', 'All Events')}</option>
                <option value="academic">{t('calendar.types.academic', 'Academic')}</option>
                <option value="exam">{t('calendar.types.exam', 'Exams')}</option>
                <option value="personal">{t('calendar.types.personal', 'Personal')}</option>
                <option value="holiday">{t('calendar.types.holiday', 'Holidays')}</option>
                <option value="meeting">{t('calendar.types.meeting', 'Meeting')}</option>
                <option value="deadline">{t('calendar.types.deadline', 'Deadline')}</option>
              </select>
            </div>

            {/* Notification Button */}
            <button
              onClick={requestNotificationPermission}
              disabled={permission === 'granted' || permission === 'denied'}
              className={cn(
                "p-2.5 rounded-xl border transition-all flex items-center justify-center gap-2",
                permission === 'granted'
                  ? "bg-primary/5 text-primary border-primary/20 shadow-sm opacity-100"
                  : permission === 'denied'
                  ? "bg-slate-100 text-slate-400 border-slate-200 opacity-50 cursor-not-allowed"
                  : "bg-white text-slate-500 border-slate-200 hover:text-slate-700 hover:border-slate-300 hover:shadow-sm"
              )}
              title={
                permission === 'granted' ? t('calendar.notifications.active', "Notifications Active") as string : 
                permission === 'denied' ? t('calendar.notifications.blocked', "Notifications Blocked") as string : 
                t('calendar.notifications.enable', "Enable Notifications") as string
              }
            >
               <Bell size={20} className={cn(permission === 'granted' && "fill-current")} />
            </button>

            <div className="flex items-center gap-2 w-full sm:w-auto">
                <div className="flex bg-slate-100 rounded-xl p-1 flex-1 sm:flex-none">
                    {(['month', 'week', 'day', 'agenda'] as ViewType[]).map((v) => {
                        const label = v === 'agenda' ? t('calendar.agenda.title', 'Agenda') : t(`calendar.${v}`, v);
                        return (
                        <button
                            key={v}
                            onClick={() => setView(v)}
                            className={cn(
                                "flex-1 sm:flex-none px-3 py-2 text-sm font-medium rounded-lg transition-all capitalize flex items-center justify-center gap-1",
                                view === v 
                                ? "bg-white text-primary shadow-sm font-bold" 
                                : "text-slate-500 hover:text-slate-700"
                            )}
                            title={label as string}
                        >
                            <span className="sm:hidden">
                                {v === 'month' && <CalendarIcon size={16} />}
                                {v === 'week' && <span className="text-xs font-bold">Wk</span>}
                                {v === 'day' && <span className="text-xs font-bold">Dy</span>}
                                {v === 'agenda' && <List size={16} />}
                            </span>
                            <span className="hidden sm:inline">{label}</span>
                        </button>
                    )})}
                </div>
                <button 
                    onClick={() => handleSelectSlot(new Date())}
                    className="p-2.5 bg-primary hover:bg-primary/90 text-white rounded-xl shadow-lg shadow-primary/30 transition-all active:scale-95"
                    title={t('calendar.new_event', 'Add Event')}
                >
                    <Plus size={20} />
                </button>
            </div>
        </div>
      </div>

      {/* Main View Area */}
      <div className="flex-1 relative min-h-0 overflow-hidden">
         {loading && (
             <div className="absolute inset-0 bg-white/60 backdrop-blur-[1px] z-50 flex items-center justify-center rounded-2xl">
                 <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
             </div>
         )}
         
         <AnimatePresence mode="wait">
            <motion.div
                key={view + currentDate.toString()}
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.2 }}
                className="h-full"
            >
                {view === 'agenda' ? (
                     <AgendaView 
                        currentDate={currentDate}
                        events={filteredEvents}
                        onSelectEvent={handleSelectEvent}
                     />
                ) : view === 'month' ? (
                     <MonthView 
                        currentDate={currentDate} 
                        events={filteredEvents}
                        onSelectEvent={handleSelectEvent}
                        onSelectSlot={handleSelectSlot}
                     />
                 ) : (
                     <WeekDayView 
                        currentDate={currentDate} 
                        view={view} 
                        events={filteredEvents}
                        onSelectEvent={handleSelectEvent}
                        onSelectSlot={handleSelectSlot}
                     />
                 )}
            </motion.div>
         </AnimatePresence>
      </div>

      <EventDialog 
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
        onSave={handleSaveEvent}
        onDelete={handleDeleteEvent}
        initialEvent={selectedEvent || undefined}
      />
    </div>
  );
};
