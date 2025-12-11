import React, { useEffect, useRef } from 'react';
import { format, addDays, startOfWeek, isSameDay, isToday, startOfDay, getHours, getMinutes, setHours, setMinutes } from 'date-fns';
import { CalendarEvent, EventType } from '../types';
import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';

interface WeekDayViewProps {
  currentDate: Date;
  view: 'week' | 'day';
  events: CalendarEvent[];
  onSelectEvent: (event: CalendarEvent) => void;
  onSelectSlot: (date: Date) => void;
}

const EVENT_COLORS: Record<EventType, string> = {
  academic: 'bg-blue-500/10 text-blue-700 border-l-4 border-blue-500',
  exam: 'bg-red-500/10 text-red-700 border-l-4 border-red-500',
  personal: 'bg-emerald-500/10 text-emerald-700 border-l-4 border-emerald-500',
  holiday: 'bg-purple-500/10 text-purple-700 border-l-4 border-purple-500',
};

export const WeekDayView: React.FC<WeekDayViewProps> = ({ currentDate, view, events, onSelectEvent, onSelectSlot }) => {
  const scrollRef = useRef<HTMLDivElement>(null);
  
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

  const getEventStyle = (event: CalendarEvent) => {
    const start = new Date(event.start);
    const end = new Date(event.end);
    
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

  return (
    <div className="flex flex-col h-full bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
      {/* Header */}
      <div className="flex border-b border-slate-200 bg-slate-50">
        <div className="w-16 flex-shrink-0 border-r border-slate-200" /> {/* Time gutter */}
        <div className={cn("flex-1 grid", view === 'week' ? "grid-cols-7" : "grid-cols-1")}>
          {days.map(day => (
            <div key={day.toISOString()} className={cn("py-2 text-center border-r border-slate-200 last:border-0", isToday(day) && "bg-primary/5")}>
              <div className="text-xs font-bold text-slate-500 uppercase">{format(day, 'EEE')}</div>
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
                  {hour}:00
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
            {days.map(day => (
              <div key={day.toISOString()} className="relative border-r border-slate-100 last:border-0 h-full group">
                {/* Click Listeners for slots */}
                {hours.map(hour => (
                   <div 
                      key={hour} 
                      className="absolute w-full hover:bg-primary/10 transition-colors z-0"
                      style={{ top: `${hour * CELL_HEIGHT}px`, height: `${CELL_HEIGHT}px` }}
                      onClick={() => handleSlotClick(day, hour)}
                   />
                ))}

                {/* Events */}
                {events
                  .filter(e => isSameDay(new Date(e.start), day))
                  .map(event => (
                    <div
                      key={event.id}
                      onClick={(e) => {
                        e.stopPropagation();
                        onSelectEvent(event);
                      }}
                      className={cn(
                        "absolute left-1 right-1 rounded-md px-2 py-1 text-xs cursor-pointer hover:brightness-95 transition-all shadow-sm z-10 overflow-hidden",
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
                    className="absolute left-0 right-0 h-0.5 bg-red-500 z-20 pointer-events-none flex items-center"
                    style={{ top: `${(getHours(new Date()) * 60 + getMinutes(new Date())) / 60 * CELL_HEIGHT}px` }}
                  >
                    <div className="w-2 h-2 rounded-full bg-red-500 -ml-1" />
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
