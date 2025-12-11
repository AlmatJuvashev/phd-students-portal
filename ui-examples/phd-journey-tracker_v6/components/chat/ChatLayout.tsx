import React, { useState } from 'react';
import { ChatSidebar } from './ChatSidebar';
import { ChatWindow } from './ChatWindow';
import { MOCK_CHAT_ROOMS } from '../../data';
import { ChatRoom, ChatUser } from '../../types';

interface ChatLayoutProps {
  onBack: () => void;
}

// Mock current user
const CURRENT_USER: ChatUser = {
  id: 'me',
  name: 'Alikhan',
  role: 'student'
};

export const ChatLayout: React.FC<ChatLayoutProps> = ({ onBack }) => {
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null);

  const activeRoom = MOCK_CHAT_ROOMS.find(r => r.id === selectedRoomId);

  return (
    <div className="fixed inset-0 bg-white z-[60] flex flex-col md:flex-row h-[100dvh]">
      {/* Mobile: Sidebar shown if no room selected. Desktop: Always shown */}
      <div className={`w-full md:w-80 lg:w-96 flex-shrink-0 h-full ${selectedRoomId ? 'hidden md:flex' : 'flex'}`}>
        <ChatSidebar 
          rooms={MOCK_CHAT_ROOMS} 
          selectedRoomId={selectedRoomId}
          onSelectRoom={setSelectedRoomId}
          className="w-full"
        />
      </div>

      {/* Mobile: Window shown if room selected. Desktop: Always shown (or placeholder) */}
      <div className={`flex-1 h-full flex flex-col relative ${!selectedRoomId ? 'hidden md:flex' : 'flex'}`}>
        {/* Mobile-only back to map button (visible only when sidebar is active) */}
        {!selectedRoomId && (
          <div className="absolute top-4 right-4 md:hidden z-50">
             <button onClick={onBack} className="bg-slate-900 text-white px-4 py-2 rounded-full text-sm font-bold shadow-lg">
               Close Chat
             </button>
          </div>
        )}

        {activeRoom ? (
          <ChatWindow 
            room={activeRoom} 
            currentUser={CURRENT_USER}
            onBack={() => setSelectedRoomId(null)}
          />
        ) : (
          <div className="hidden md:flex flex-col items-center justify-center h-full bg-slate-50 text-slate-400 p-8 text-center border-l border-slate-200">
            <div className="w-64 h-64 bg-slate-100 rounded-full mb-8 flex items-center justify-center animate-pulse">
               <span className="text-6xl">💬</span>
            </div>
            <h2 className="text-2xl font-bold text-slate-800 mb-2">Welcome to Student Chat</h2>
            <p>Select a conversation from the sidebar to start messaging your advisors or support.</p>
            
            <button 
              onClick={onBack} 
              className="mt-8 px-6 py-3 bg-white border border-slate-300 rounded-xl font-bold text-slate-600 hover:bg-slate-50 transition-colors shadow-sm"
            >
              Back to Journey Map
            </button>
          </div>
        )}
      </div>
    </div>
  );
};