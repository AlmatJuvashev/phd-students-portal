import React from 'react';
import { CalendarView } from './components/CalendarView';
import { ArrowLeft } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

export const CalendarPage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col">
      {/* Content */}
      <main className="flex-1 max-w-6xl mx-auto w-full px-4 py-6 sm:py-8">
        <CalendarView />
      </main>
    </div>
  );
};
