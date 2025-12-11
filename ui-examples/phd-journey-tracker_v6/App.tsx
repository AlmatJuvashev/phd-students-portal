import React, { useState } from 'react';
import { JourneyMap } from './components/JourneyMap';
import { ChatLayout } from './components/chat/ChatLayout';
import { CalendarPage } from './components/calendar/CalendarPage';
import { PHD_PLAYBOOK } from './data';
import { Locale } from './types';

export default function App() {
  const [locale, setLocale] = useState<Locale>('ru');
  const [currentView, setCurrentView] = useState<'map' | 'chat' | 'calendar'>('map');

  const renderView = () => {
    switch (currentView) {
      case 'chat':
        return <ChatLayout onBack={() => setCurrentView('map')} />;
      case 'calendar':
        return <CalendarPage onBack={() => setCurrentView('map')} />;
      case 'map':
      default:
        return (
          <div className="max-w-4xl mx-auto px-4 py-8">
            {/* Locale Switcher for Demo */}
            <div className="flex justify-end mb-4 gap-2">
              {(['ru', 'kz', 'en'] as Locale[]).map((l) => (
                <button
                  key={l}
                  onClick={() => setLocale(l)}
                  className={`px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wide transition-colors ${
                    locale === l 
                      ? 'bg-primary-600 text-white' 
                      : 'bg-white text-slate-400 hover:text-slate-600'
                  }`}
                >
                  {l}
                </button>
              ))}
            </div>

            {/* The Map Component */}
            <JourneyMap 
              playbook={PHD_PLAYBOOK} 
              currentActiveNodeId="VI_attestation_file"
              locale={locale}
              onOpenChat={() => setCurrentView('chat')}
              onOpenCalendar={() => setCurrentView('calendar')}
            />
            
            <div className="text-center text-slate-400 text-sm mt-12 pb-8">
              PhD Portal © 2025 KazNMU
            </div>
          </div>
        );
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 font-sans text-slate-900 selection:bg-primary-200 selection:text-primary-900">
      {renderView()}
    </div>
  );
}