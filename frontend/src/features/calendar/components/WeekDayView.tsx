
import React, { useEffect, useRef } from 'react';
import { format, addDays, startOfWeek, isSameDay, isToday, startOfDay, getHours, getMinutes, setHours, setMinutes } from 'date-fns';
import { enUS, ru, kk } from 'date-fns/locale';
import { CalendarEvent, EventType } from '../types';
import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';

const LOCALES: Record<string, any> = {
  en: enUS,
  ru: ru,
  kz: kk,
};

interface WeekDayViewProps {
  currentDate: Date;
  view: 'week' | 'day';
  events: CalendarEvent[];
  onSelectEvent: (event: CalendarEvent) => void;
  onSelectSlot: (date: Date) => void;
  onEventDrop: (event: CalendarEvent, newStart: Date) => void;
}

const EVENT_COLORS: Record<EventType, string> = {
  academic: 'bg-blue-500/10 text-blue-700 border-l-4 border-blue-500',
  exam: 'bg-red-500/10 text-red-700 border-l-4 border-red-500',
  personal: 'bg-emerald-500/10 text-emerald-700 border-l-4 border-emerald-500',
  holiday: 'bg-purple-500/10 text-purple-700 border-l-4 border-purple-500',
  meeting: 'bg-amber-500/10 text-amber-700 border-l-4 border-amber-500',
  deadline: 'bg-orange-500/10 text-orange-500 border-l-4 border-orange-500',
};

export const WeekDayView: React.FC<WeekDayViewProps> = ({ currentDate, events, view, onSelectEvent, onSelectSlot, onEventDrop }) => {
  const scrollRef = useRef<HTMLDivElement>(null);
  const { t, i18n } = useTranslation();
  const currentLocale = i18n.language === 'kz' ? kk : i18n.language === 'ru' ? ru : enUS;
  
  // Generate columns
  const days = view === 'week' 
    ? Array.from({ length: 7 }, (_, i) => addDays(startOfWeek(currentDate, { weekStartsOn: 1 }), i))
    : [currentDate];

  // Hours 0-23
  const hours = Array.from({ length: 24 }, (_, i) => i);
  const CELL_HEIGHT = 60; // px per hour

  // Scroll to 8 AM on mount
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = 8 * CELL_HEIGHT;
    }
  }, [view]);

  const [dragTarget, setDragTarget] = React.useState<string | null>(null);

  const getEventStyle = (event: CalendarEvent) => {
    const start = new Date(event.start);
    const end = new Date(event.end);
    
    // Normalize dates to current day for positioning if spanning multiple days (simplified for now)
    const startHour = getHours(start);
    const startMin = getMinutes(start);
    const endHour = getHours(end);
    const endMin = getMinutes(end);
    
    // Calculate top position in minutes from start of day
    const topMinutes = startHour * 60 + startMin;
    const durationMinutes = (endHour * 60 + endMin) - topMinutes;
    
    return {
      top: `${(topMinutes / 60) * CELL_HEIGHT}px`,
      height: `${Math.max((durationMinutes / 60) * CELL_HEIGHT, 20)}px`, // Min height 20px
    };
  };

  const handleSlotClick = (day: Date, hour: number) => {
    const clickedDate = setMinutes(setHours(day, hour), 0);
    onSelectSlot(clickedDate);
  };

  const handleDragStart = (e: React.DragEvent, event: CalendarEvent) => {
    e.dataTransfer.setData('eventId', event.id);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };
  
  const handleDragEnter = (e: React.DragEvent, dayIso: string, hour: number) => {
      e.preventDefault();
      setDragTarget(`${dayIso}-${hour}`);
  };

  const handleDrop = (e: React.DragEvent, day: Date, hour: number) => {
    e.preventDefault();
    setDragTarget(null);
    const eventId = e.dataTransfer.getData('eventId');
    const event = events.find(e => e.id === eventId);
    if (event) {
        const newStart = new Date(day);
        newStart.setHours(hour);
        newStart.setMinutes(0); // Snap to top of the hour for now
        onEventDrop(event, newStart);
    }
  };

  return (
    <div className="flex flex-col h-full bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden" onMouseLeave={() => setDragTarget(null)}>
      {/* Header */}
      <div className="flex border-b border-slate-200 bg-slate-50">
        <div className="w-16 flex-shrink-0 border-r border-slate-200" /> {/* Time gutter */}
        <div className={cn("flex-1 grid", view === 'week' ? "grid-cols-7" : "grid-cols-1")}>
          {days.map(day => (
            <div key={day.toISOString()} className={cn("py-2 text-center border-r border-slate-200 last:border-0", isToday(day) && "bg-primary/5")}>
              <div className="text-xs font-bold text-slate-500 uppercase">{format(day, 'EEE', { locale: currentLocale })}</div>
              <div className={cn(
                "w-8 h-8 mx-auto mt-1 flex items-center justify-center rounded-full text-sm font-bold",
                isToday(day) ? "bg-primary text-white shadow-md shadow-primary/30" : "text-slate-800"
              )}>
                {format(day, 'd')}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Scrollable Grid */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto relative custom-scrollbar">
        <div className="flex relative min-h-[1440px]"> {/* 24 * 60 */}
          
          {/* Time Labels */}
          <div className="w-16 flex-shrink-0 border-r border-slate-200 bg-white sticky left-0 z-10">
            {hours.map(hour => (
              <div key={hour} className="relative" style={{ height: `${CELL_HEIGHT}px` }}>
                <span className="absolute -top-3 right-2 text-xs text-slate-400 font-medium bg-white px-1">
                  {hour < 10 ? `0${hour}:00` : `${hour}:00`}
                </span>
                {/* Tick mark */}
                <div className="absolute top-0 right-0 w-2 h-px bg-slate-200" />
              </div>
            ))}
          </div>

          {/* Columns */}
          <div className={cn("flex-1 grid relative", view === 'week' ? "grid-cols-7" : "grid-cols-1")}>
            {/* Horizontal Grid Lines */}
            <div className="absolute inset-0 z-0 pointer-events-none">
              {hours.map(hour => (
                <div key={hour} className="border-t border-slate-100 w-full" style={{ height: `${CELL_HEIGHT}px` }} />
              ))}
            </div>

            {/* Day Columns */}
            {days.map(day => {
                const dayIso = day.toISOString();
                return (
                  <div key={dayIso} className="relative border-r border-slate-100 last:border-0 h-full group">
                    {/* Click Listeners for slots */}
                    {hours.map(hour => {
                       const isTarget = dragTarget === `${dayIso}-${hour}`;
                       return (
                           <div 
                              key={hour} 
                              className={cn(
                                "absolute w-full transition-colors z-0",
                                isTarget ? "bg-primary/10 ring-inset ring-2 ring-primary/20" : "hover:bg-primary/5"
                              )}
                              style={{ top: `${hour * CELL_HEIGHT}px`, height: `${CELL_HEIGHT}px` }}
                              onClick={() => handleSlotClick(day, hour)}
                              onDragOver={handleDragOver}
                              onDragEnter={(e) => handleDragEnter(e, dayIso, hour)}
                              onDrop={(e) => handleDrop(e, day, hour)}
                           />
                       );
                    })}

                {/* Events */}
                {events
                  .filter(e => isSameDay(new Date(e.start), day))
                  .map(event => (
                    <div
                      key={event.id}
                      draggable
                      onDragStart={(e) => handleDragStart(e, event)}
                      onClick={(e) => {
                        e.stopPropagation();
                        onSelectEvent(event);
                      }}
                      className={cn(
                        "absolute left-1 right-1 rounded-md px-2 py-1 text-xs hover:brightness-95 transition-all shadow-sm z-10 overflow-hidden cursor-move",
                        EVENT_COLORS[event.type as EventType]
                      )}
                      style={getEventStyle(event)}
                    >
                      <div className="font-bold truncate leading-tight">{event.title}</div>
                      <div className="opacity-80 truncate text-[10px]">
                        {format(new Date(event.start), 'HH:mm')} - {format(new Date(event.end), 'HH:mm')}
                      </div>
                    </div>
                  ))}

                {/* Current Time Indicator */}
                {isToday(day) && (
                  <div 
                    className="absolute left-0 right-0 h-0.5 bg-primary z-20 pointer-events-none flex items-center"
                    style={{ top: `${(getHours(new Date()) * 60 + getMinutes(new Date())) / 60 * CELL_HEIGHT}px` }}
                  >
                    <div className="w-2 h-2 rounded-full bg-primary -ml-1" />
                  </div>
                )}
              </div>
            );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};
