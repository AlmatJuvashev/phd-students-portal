import React from 'react';
import { CalendarView } from './CalendarView';
import { ArrowLeft } from 'lucide-react';

interface CalendarPageProps {
  onBack: () => void;
}

export const CalendarPage: React.FC<CalendarPageProps> = ({ onBack }) => {
  return (
    <div className="min-h-screen bg-slate-50 flex flex-col">
      {/* Header */}
      <div className="bg-white border-b border-slate-200 sticky top-0 z-20">
        <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <button 
                onClick={onBack}
                className="p-2 -ml-2 text-slate-500 hover:bg-slate-100 rounded-full transition-colors"
            >
                <ArrowLeft size={24} />
            </button>
            <div>
              <h1 className="text-xl sm:text-2xl font-bold text-slate-900">Schedule</h1>
              <p className="text-xs sm:text-sm text-slate-500">Manage your deadlines and events</p>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <main className="flex-1 max-w-6xl mx-auto w-full px-4 py-6 sm:py-8">
        <CalendarView />
      </main>
    </div>
  );
};