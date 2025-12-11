import React from 'react';
import { startOfMonth, endOfMonth, startOfWeek, endOfWeek, eachDayOfInterval, isSameMonth, isSameDay, format, isToday } from 'date-fns';
import { CalendarEvent, EventType } from '../../types';
import { cn } from '../../lib/utils';
import { Clock } from 'lucide-react';

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
};

export const MonthView: React.FC<MonthViewProps> = ({ currentDate, events, onSelectEvent, onSelectSlot, onEventDrop }) => {
  const monthStart = startOfMonth(currentDate);
  const monthEnd = endOfMonth(monthStart);
  const startDate = startOfWeek(monthStart, { weekStartsOn: 0 }); // Sunday start
  const endDate = endOfWeek(monthEnd, { weekStartsOn: 0 });

  const calendarDays = eachDayOfInterval({ start: startDate, end: endDate });
  const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  const handleDragStart = (e: React.DragEvent, event: CalendarEvent) => {
    e.dataTransfer.setData('eventId', event.id);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent, date: Date) => {
    e.preventDefault();
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
    <div className="flex flex-col h-full bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
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
          const dayEvents = events.filter(e => isSameDay(new Date(e.start), day));
          const isCurrentMonth = isSameMonth(day, monthStart);
          
          return (
            <div 
              key={day.toISOString()}
              onClick={() => onSelectSlot(day)}
              onDragOver={handleDragOver}
              onDrop={(e) => handleDrop(e, day)}
              className={cn(
                "min-h-[100px] bg-white p-1 transition-colors relative group hover:bg-slate-50 cursor-pointer",
                !isCurrentMonth && "bg-slate-50/50 text-slate-400"
              )}
            >
              {/* Day Number */}
              <div className="flex justify-between items-start px-1">
                <span className={cn(
                  "text-xs font-semibold w-6 h-6 flex items-center justify-center rounded-full mb-1",
                  isToday(day) 
                    ? "bg-primary-600 text-white shadow-md shadow-primary-500/30" 
                    : "text-slate-700"
                )}>
                  {format(day, 'd')}
                </span>
              </div>

              {/* Events List */}
              <div className="flex flex-col gap-1 overflow-hidden max-h-[calc(100%-2rem)]">
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
                      "w-full text-left px-1.5 py-0.5 text-[10px] font-medium rounded-sm truncate transition-transform hover:scale-[1.02] cursor-move",
                      TYPE_STYLES[event.type as EventType]
                    )}
                  >
                    {format(new Date(event.start), 'HH:mm')} {event.title}
                  </div>
                ))}
                {dayEvents.length > 3 && (
                  <span className="text-[10px] text-slate-400 font-medium pl-1">
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