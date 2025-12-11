import React, { useState, useEffect } from 'react';
import { format, addMonths, subMonths, addWeeks, subWeeks, addDays, subDays } from 'date-fns';
import { CalendarEvent } from '../types';
import { ChevronLeft, ChevronRight, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';
import { EventDialog } from './EventDialog';
import { fetchEvents, createEvent, updateEvent, deleteEvent } from '../api';
import { MonthView } from './MonthView';
import { WeekDayView } from './WeekDayView';
import { useTranslation } from 'react-i18next';

type ViewType = 'month' | 'week' | 'day';

export const CalendarView: React.FC = () => {
  const { t } = useTranslation();
  const [currentDate, setCurrentDate] = useState(new Date());
  const [view, setView] = useState<ViewType>('month');
  const [events, setEvents] = useState<CalendarEvent[]>([]);
  const [selectedEvent, setSelectedEvent] = useState<Partial<CalendarEvent> | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadEvents();
  }, []);

  const loadEvents = async () => {
    setLoading(true);
    // In real app, pass start/end based on view
    const data = await fetchEvents(new Date(), new Date());
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

  const handleSaveEvent = async (eventData: CalendarEvent | Omit<CalendarEvent, 'id'>) => {
    if ('id' in eventData) {
      await updateEvent(eventData as CalendarEvent);
      setEvents(prev => prev.map(e => e.id === eventData.id ? eventData as CalendarEvent : e));
    } else {
      const newEvent = await createEvent(eventData);
      setEvents(prev => [...prev, newEvent]);
    }
    setIsDialogOpen(false);
    setSelectedEvent(null);
  };

  const handleDeleteEvent = async (id: string) => {
    await deleteEvent(id);
    setEvents(prev => prev.filter(e => e.id !== id));
    setIsDialogOpen(false);
    setSelectedEvent(null);
  };

  return (
    <div className="h-[700px] flex flex-col gap-4">
      {/* Calendar Toolbar */}
      <div className="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white p-3 rounded-2xl shadow-sm border border-slate-200">
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

        <h2 className="text-xl font-bold text-slate-800 tabular-nums capitalize">
          {format(currentDate, 'MMMM yyyy')}
        </h2>

        <div className="flex items-center gap-2">
            <div className="flex bg-slate-100 rounded-xl p-1">
                {(['month', 'week', 'day'] as ViewType[]).map((v) => (
                    <button
                        key={v}
                        onClick={() => setView(v)}
                        className={cn(
                            "px-4 py-2 text-sm font-medium rounded-lg transition-all capitalize",
                            view === v 
                            ? "bg-white text-primary shadow-sm font-bold" 
                            : "text-slate-500 hover:text-slate-700"
                        )}
                    >
                        {t(`calendar.${v}`, v)}
                    </button>
                ))}
            </div>
            <button 
                onClick={() => handleSelectSlot(new Date())}
                className="p-2.5 bg-primary hover:bg-primary/90 text-white rounded-xl shadow-lg shadow-primary/30 transition-all active:scale-95"
            >
                <Plus size={20} />
            </button>
        </div>
      </div>

      {/* Main View Area */}
      <div className="flex-1 relative min-h-0">
         {loading && (
             <div className="absolute inset-0 bg-white/60 backdrop-blur-[1px] z-50 flex items-center justify-center rounded-2xl">
                 <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
             </div>
         )}
         
         {view === 'month' ? (
             <MonthView 
                currentDate={currentDate} 
                events={events}
                onSelectEvent={handleSelectEvent}
                onSelectSlot={handleSelectSlot}
             />
         ) : (
             <WeekDayView 
                currentDate={currentDate} 
                view={view} 
                events={events}
                onSelectEvent={handleSelectEvent}
                onSelectSlot={handleSelectSlot}
             />
         )}
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
