import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, MapPin, Trash2, Save } from 'lucide-react';
import { CalendarEvent, EventType } from '../types';
import { cn } from '@/lib/utils';
import { useTranslation } from 'react-i18next';

interface EventDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (event: Omit<CalendarEvent, 'id'> | CalendarEvent) => void;
  onDelete?: (id: string) => void;
  initialEvent?: Partial<CalendarEvent>;
}

export const EventDialog: React.FC<EventDialogProps> = ({ 
  isOpen, 
  onClose, 
  onSave, 
  onDelete,
  initialEvent 
}) => {
  const { t } = useTranslation();
  const [title, setTitle] = useState('');
  const [type, setType] = useState<EventType>('academic');
  const [start, setStart] = useState('');
  const [end, setEnd] = useState('');
  const [location, setLocation] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');

  const EVENT_TYPES: { value: EventType; label: string; color: string }[] = [
    { value: 'academic', label: t('calendar.types.academic', 'Academic'), color: 'bg-blue-100 text-blue-700 border-blue-200' },
    { value: 'exam', label: t('calendar.types.exam', 'Exam'), color: 'bg-red-100 text-red-700 border-red-200' },
    { value: 'personal', label: t('calendar.types.personal', 'Personal'), color: 'bg-emerald-100 text-emerald-700 border-emerald-200' },
    { value: 'holiday', label: t('calendar.types.holiday', 'Holiday'), color: 'bg-purple-100 text-purple-700 border-purple-200' },
  ];

  useEffect(() => {
    if (isOpen) {
      if (initialEvent) {
        setTitle(initialEvent.title || '');
        setType(initialEvent.type || 'academic');
        setLocation(initialEvent.location || '');
        setDescription(initialEvent.description || '');
        
        // Format dates for datetime-local input
        const formatDate = (d?: Date) => {
            if (!d) return '';
            const offset = d.getTimezoneOffset() * 60000;
            return new Date(d.getTime() - offset).toISOString().slice(0, 16);
        };
        
        setStart(formatDate(initialEvent.start));
        setEnd(formatDate(initialEvent.end));
      } else {
        // Reset form
        setTitle('');
        setType('academic');
        setStart('');
        setEnd('');
        setLocation('');
        setDescription('');
      }
      setError('');
    }
  }, [isOpen, initialEvent]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      setError(t('calendar.validation.title_required', 'Title is required'));
      return;
    }
    if (!start || !end) {
      setError(t('calendar.validation.dates_required', 'Start and End times are required'));
      return;
    }
    if (new Date(start) >= new Date(end)) {
      setError(t('calendar.validation.end_after_start', 'End time must be after start time'));
      return;
    }

    onSave({
      ...(initialEvent?.id ? { id: initialEvent.id } : {}),
      title,
      type,
      start: new Date(start),
      end: new Date(end),
      location,
      description
    } as CalendarEvent);
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="bg-white w-full max-w-md rounded-2xl shadow-2xl overflow-hidden"
          >
            {/* Header */}
            <div className="flex justify-between items-center p-4 border-b border-slate-100">
              <h2 className="text-lg font-bold text-slate-800">
                {initialEvent?.id ? t('calendar.edit_event', 'Edit Event') : t('calendar.new_event', 'New Event')}
              </h2>
              <button 
                onClick={onClose}
                className="p-2 hover:bg-slate-100 rounded-full text-slate-400 hover:text-slate-600 transition-colors"
              >
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              {error && (
                <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg">
                  {error}
                </div>
              )}

              {/* Title */}
              <div>
                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">
                  {t('calendar.form.title', 'Title')}
                </label>
                <input
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all font-medium text-slate-800"
                  placeholder={t('calendar.placeholders.title', 'Enter event title')}
                />
              </div>

              {/* Type Selection */}
              <div>
                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">
                  {t('calendar.form.type', 'Type')}
                </label>
                <div className="grid grid-cols-2 gap-2">
                  {EVENT_TYPES.map(t => (
                    <button
                      key={t.value}
                      type="button"
                      onClick={() => setType(t.value)}
                      className={cn(
                        "px-3 py-2 rounded-lg text-sm font-bold border transition-all text-left flex items-center gap-2",
                        type === t.value 
                          ? t.color + " ring-2 ring-offset-1 ring-primary/40" 
                          : "bg-white border-slate-200 text-slate-500 hover:border-slate-300"
                      )}
                    >
                      <div className={cn("w-2 h-2 rounded-full", t.color.split(' ')[0].replace('bg-', 'bg-current opacity-50'))} />
                      {t.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* Date & Time */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">
                    {t('calendar.form.start', 'Start')}
                  </label>
                  <input
                    type="datetime-local"
                    value={start}
                    onChange={(e) => setStart(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all text-sm"
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">
                    {t('calendar.form.end', 'End')}
                  </label>
                  <input
                    type="datetime-local"
                    value={end}
                    onChange={(e) => setEnd(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all text-sm"
                  />
                </div>
              </div>

              {/* Location */}
              <div>
                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">
                  {t('calendar.form.location', 'Location')}
                </label>
                <div className="relative">
                  <MapPin size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input
                    type="text"
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                    className="w-full pl-9 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all text-sm"
                    placeholder={t('calendar.placeholders.location', 'e.g. Conference Hall B')}
                  />
                </div>
              </div>

              {/* Description */}
              <div>
                <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1.5">
                 {t('calendar.form.description', 'Description')}
                </label>
                <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  rows={3}
                  className="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all text-sm resize-none"
                  placeholder={t('calendar.placeholders.description', 'Add details about the event...')}
                />
              </div>

              {/* Actions */}
              <div className="flex items-center justify-between pt-4 mt-2 border-t border-slate-100">
                {initialEvent?.id && onDelete ? (
                  <button
                    type="button"
                    onClick={() => onDelete(initialEvent.id!)}
                    className="flex items-center gap-2 px-4 py-2 text-red-600 hover:bg-red-50 rounded-lg font-medium transition-colors text-sm"
                  >
                    <Trash2 size={16} /> {t('calendar.form.delete', 'Delete')}
                  </button>
                ) : <div />} {/* Spacer */}

                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={onClose}
                    className="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-lg font-medium transition-colors text-sm"
                  >
                    {t('calendar.form.cancel', 'Cancel')}
                  </button>
                  <button
                    type="submit"
                    className="flex items-center gap-2 px-6 py-2 bg-primary hover:bg-primary/90 text-white rounded-lg font-bold shadow-lg shadow-primary/30 transition-all active:scale-95 text-sm"
                  >
                    <Save size={16} /> {t('calendar.form.save', 'Save Event')}
                  </button>
                </div>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
};
