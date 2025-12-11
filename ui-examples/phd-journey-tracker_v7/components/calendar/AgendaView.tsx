import React, { useEffect, useRef } from 'react';
import { format, isSameDay, isToday, startOfMonth, endOfMonth, eachDayOfInterval } from 'date-fns';
import { CalendarEvent, EventType } from '../../types';
import { cn } from '../../lib/utils';
import { MapPin, Calendar as CalendarIcon } from 'lucide-react';

interface AgendaViewProps {
  currentDate: Date;
  events: CalendarEvent[];
  onSelectEvent: (event: CalendarEvent) => void;
}

const TYPE_STYLES: Record<EventType, string> = {
  academic: 'bg-blue-50 border-blue-200 text-blue-700',
  exam: 'bg-red-50 border-red-200 text-red-700',
  personal: 'bg-emerald-50 border-emerald-200 text-emerald-700',
  holiday: 'bg-purple-50 border-purple-200 text-purple-700',
};

export const AgendaView: React.FC<AgendaViewProps> = ({ currentDate, events, onSelectEvent }) => {
  const scrollRef = useRef<HTMLDivElement>(null);
  
  // Show events for the current month range
  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(currentDate);
  const days = eachDayOfInterval({ start: monthStart, end: monthEnd });

  // Determine which days to show: 
  // 1. All days in the month? Or just days with events? 
  // For a "Schedule" feel, showing consecutive days is better context.
  // We'll filter to days that have events OR are today, to keep the list compact but useful.
  // If the list is too short, we can show the whole month.
  // Let's show the whole month to allow users to see empty days too.
  
  // Scroll to today on mount if it's in the list
  useEffect(() => {
    const todayElement = document.getElementById('agenda-today');
    if (todayElement && scrollRef.current) {
        const top = todayElement.offsetTop;
        scrollRef.current.scrollTop = top - 20; // some padding
    }
  }, [currentDate]);

  return (
    <div ref={scrollRef} className="flex flex-col h-full bg-white rounded-xl shadow-sm border border-slate-200 overflow-y-auto custom-scrollbar p-3 sm:p-4 space-y-4">
      {days.map(day => {
        const dayEvents = events.filter(e => isSameDay(new Date(e.start), day));
        dayEvents.sort((a, b) => new Date(a.start).getTime() - new Date(b.start).getTime());
        
        const isCurrentDay = isToday(day);
        const hasEvents = dayEvents.length > 0;

        // Skip empty days in the past to reduce clutter? 
        // No, let's keep it consistent like a planner.

        return (
          <div 
            key={day.toISOString()} 
            id={isCurrentDay ? 'agenda-today' : undefined}
            className={cn("flex gap-3 sm:gap-4", !hasEvents && !isCurrentDay && "opacity-50")}
          >
            {/* Date Column */}
            <div className="flex flex-col items-center w-12 sm:w-14 flex-shrink-0 pt-1">
               <span className={cn(
                   "text-[10px] sm:text-xs font-bold uppercase", 
                   isCurrentDay ? "text-primary-600" : "text-slate-400"
                )}>
                   {format(day, 'EEE')}
               </span>
               <div className={cn(
                 "w-8 h-8 sm:w-10 sm:h-10 flex items-center justify-center rounded-full text-sm sm:text-lg font-bold mt-1 transition-all",
                 isCurrentDay 
                    ? "bg-primary-600 text-white shadow-md shadow-primary-500/30 scale-110" 
                    : "bg-slate-100 text-slate-700"
               )}>
                 {format(day, 'd')}
               </div>
            </div>

            {/* Events Column */}
            <div className={cn(
                "flex-1 space-y-2 pb-4 border-b border-slate-100 min-h-[3rem]",
                isCurrentDay && "border-primary-100"
            )}>
               {hasEvents ? (
                 dayEvents.map(event => (
                   <button
                     key={event.id}
                     onClick={() => onSelectEvent(event)}
                     className={cn(
                       "w-full text-left p-3 rounded-xl border flex flex-col gap-1.5 transition-all active:scale-[0.98] shadow-sm hover:shadow-md",
                       TYPE_STYLES[event.type as EventType]
                     )}
                   >
                     <div className="flex justify-between items-start gap-2">
                        <span className="font-bold text-sm leading-tight">{event.title}</span>
                        {event.type !== 'holiday' && (
                            <span className="text-[10px] font-bold opacity-80 whitespace-nowrap bg-white/60 px-2 py-0.5 rounded-full">
                                {format(new Date(event.start), 'h:mm a')}
                            </span>
                        )}
                     </div>
                     
                     {(event.location || event.description) && (
                        <div className="flex items-center gap-3 text-xs opacity-80">
                            {event.location && (
                                <span className="flex items-center gap-1 font-medium">
                                    <MapPin size={12} /> {event.location}
                                </span>
                            )}
                        </div>
                     )}
                   </button>
                 ))
               ) : (
                   isCurrentDay ? (
                       <div className="text-sm text-slate-400 italic py-2 flex items-center gap-2">
                           <CalendarIcon size={14} /> No events scheduled for today
                       </div>
                   ) : (
                       <div className="h-full" /> // Spacer
                   )
               )}
            </div>
          </div>
        )
      })}
      
      {days.length === 0 && (
         <div className="text-center py-10 text-slate-400">
             No dates to display
         </div>
      )}
    </div>
  );
};