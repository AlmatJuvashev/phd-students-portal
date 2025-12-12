import React from 'react';
import { startOfMonth, endOfMonth, startOfWeek, endOfWeek, eachDayOfInterval, isSameMonth, isSameDay, format, isToday } from 'date-fns';
import { CalendarEvent, EventType } from '../types';
import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';

interface MonthViewProps {
  currentDate: Date;
  events: CalendarEvent[];
  onSelectEvent: (event: CalendarEvent) => void;
  onSelectSlot: (date: Date) => void;
  onEventDrop: (event: CalendarEvent, newStart: Date) => void;
}

const TYPE_STYLES: Record<EventType, string> = {
  academic: 'bg-blue-100 text-blue-700 border-l-2 border-blue-500',
  exam: 'bg-red-100 text-red-700 border-l-2 border-red-500',
  personal: 'bg-emerald-100 text-emerald-700 border-l-2 border-emerald-500',
  holiday: 'bg-purple-100 text-purple-700 border-l-2 border-purple-500',
  meeting: 'bg-amber-100 text-amber-700 border-l-2 border-amber-500',
  deadline: 'bg-orange-100 text-orange-700 border-l-2 border-orange-500',
};

export const MonthView: React.FC<MonthViewProps> = ({ currentDate, events, onSelectEvent, onSelectSlot, onEventDrop }) => {
  const { t } = useTranslation();
  
  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(monthStart);
  const startDate = startOfWeek(monthStart, { weekStartsOn: 1 }); // Monday start
  const endDate = endOfWeek(monthEnd, { weekStartsOn: 1 });

  const calendarDays = eachDayOfInterval({ start: startDate, end: endDate });
  
  const weekDays = [
    t('calendar.days.mon', 'Mon'),
    t('calendar.days.tue', 'Tue'),
    t('calendar.days.wed', 'Wed'),
    t('calendar.days.thu', 'Thu'),
    t('calendar.days.fri', 'Fri'),
    t('calendar.days.sat', 'Sat'),
    t('calendar.days.sun', 'Sun'),
  ];

  const [dragTarget, setDragTarget] = React.useState<string | null>(null);

  const handleDragStart = (e: React.DragEvent, event: CalendarEvent) => {
    e.dataTransfer.setData('eventId', event.id);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };
  
  const handleDragEnter = (e: React.DragEvent, dateStr: string) => {
      e.preventDefault();
      setDragTarget(dateStr);
  };
  
  const handleDragLeave = (e: React.DragEvent) => {
      e.preventDefault();
      // Only clear if we are actually leaving the container, not entering a child
      // But simpler for now: layout handles it or we rely on Enter of next sibling
  };

  const handleDrop = (e: React.DragEvent, date: Date) => {
    e.preventDefault();
    setDragTarget(null);
    const eventId = e.dataTransfer.getData('eventId');
    const event = events.find(e => e.id === eventId);
    if (event) {
        // Preserve original time when dropping on a day in month view
        const newStart = new Date(date);
        const originalStart = new Date(event.start);
        newStart.setHours(originalStart.getHours(), originalStart.getMinutes());
        onEventDrop(event, newStart);
    }
  };

  return (
    <div className="flex flex-col h-full bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden" onMouseLeave={() => setDragTarget(null)}>
      {/* Header Row */}
      <div className="grid grid-cols-7 border-b border-slate-200 bg-slate-50">
        {weekDays.map(day => (
          <div key={day} className="py-2 text-center text-xs font-bold text-slate-500 uppercase tracking-wider">
            {day} 
          </div>
        ))}
      </div>

      {/* Days Grid */}
      <div className="flex-1 grid grid-cols-7 auto-rows-fr bg-slate-200 gap-px border-b border-slate-200">
        {calendarDays.map((day) => {
          const dayStr = day.toISOString();
          const dayEvents = events.filter(e => isSameDay(new Date(e.start), day));
          const isCurrentMonth = isSameMonth(day, monthStart);
          const isDragTarget = dragTarget === dayStr;
          
          return (
            <div 
              key={dayStr}
              onClick={() => onSelectSlot(day)}
              onDragOver={handleDragOver}
              onDragEnter={(e) => handleDragEnter(e, dayStr)}
              onDrop={(e) => handleDrop(e, day)}
              className={cn(
                "min-h-[100px] bg-white p-1 transition-colors relative group cursor-pointer",
                isDragTarget ? "bg-primary/10 ring-2 ring-primary/20 z-10" : "hover:bg-slate-50",
                !isCurrentMonth && !isDragTarget && "bg-slate-50/50 text-slate-400"
              )}
            >
              {/* Day Number */}
              <div className="flex justify-between items-start px-1">
                <span className={cn(
                  "text-xs font-semibold w-6 h-6 flex items-center justify-center rounded-full mb-1 transition-all",
                  isToday(day) 
                    ? "bg-primary text-white shadow-md shadow-primary/30 scale-110" 
                    : "text-slate-700",
                  isDragTarget && !isToday(day) && "bg-white text-primary shadow-sm"
                )}>
                  {format(day, 'd')}
                </span>
              </div>

              {/* Events List */}
              <div className="flex flex-col gap-1.5 overflow-hidden max-h-[calc(100%-2rem)]">
                {dayEvents.slice(0, 3).map(event => (
                  <div
                    key={event.id}
                    draggable
                    onDragStart={(e) => handleDragStart(e, event)}
                    onClick={(e) => {
                      e.stopPropagation();
                      onSelectEvent(event);
                    }}
                    className={cn(
                      "w-full text-left px-2 py-1 text-[10px] font-semibold rounded-md truncate transition-all cursor-move shadow-sm border border-transparent hover:border-black/5 hover:shadow-md",
                      TYPE_STYLES[event.type as EventType]
                    )}
                  >
                    {format(new Date(event.start), 'HH:mm')} {event.title}
                  </div>
                ))}
                {dayEvents.length > 3 && (
                  <span className="text-[10px] text-slate-400 font-bold pl-1 hover:text-primary transition-colors">
                    +{dayEvents.length - 3} more
                  </span>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
