import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowLeft, Send, Paperclip, MoreHorizontal, Image as ImageIcon, FileText, Check, Clock } from 'lucide-react';
import { ChatRoom, Message, ChatUser } from '../../types';
import { cn } from '../../lib/utils';

interface ChatWindowProps {
  room: ChatRoom;
  currentUser: ChatUser;
  onBack: () => void; // Mobile only
  className?: string;
}

export const ChatWindow: React.FC<ChatWindowProps> = ({ room, currentUser, onBack, className }) => {
  const [messages, setMessages] = useState<Message[]>(room.messages || []);
  const [inputValue, setInputValue] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    setMessages(room.messages);
    scrollToBottom();
  }, [room]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = () => {
    if (!inputValue.trim()) return;

    const newMessage: Message = {
      id: Date.now().toString(),
      senderId: currentUser.id,
      content: inputValue,
      timestamp: new Date().toISOString(),
      status: 'sending'
    };

    // Optimistic Update
    setMessages(prev => [...prev, newMessage]);
    setInputValue('');

    // Simulate API delay
    setTimeout(() => {
      setMessages(prev => prev.map(m => 
        m.id === newMessage.id ? { ...m, status: 'sent' } : m
      ));
    }, 1000);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const formatMessageDate = (isoString: string) => {
    return new Date(isoString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const renderBubble = (msg: Message, isMe: boolean) => {
    return (
      <motion.div 
        layout
        initial={{ opacity: 0, scale: 0.8, y: 20 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        transition={{ 
          opacity: { duration: 0.2 },
          layout: { type: "spring", bounce: 0.4 },
          type: "spring", stiffness: 400, damping: 25 
        }}
        className={cn(
          "flex w-full mb-4",
          isMe ? "justify-end" : "justify-start"
        )}
      >
        <div className={cn(
          "max-w-[85%] sm:max-w-[70%] rounded-2xl p-3 sm:p-4 relative shadow-sm text-sm sm:text-base",
          isMe 
            ? "bg-primary-600 text-white rounded-tr-none origin-bottom-right" 
            : "bg-white border border-slate-100 text-slate-800 rounded-tl-none origin-bottom-left"
        )}>
          {/* Sender Name for Group Chats (Not me) */}
          {!isMe && room.type !== 'private' && (
             <div className="text-xs font-bold text-primary-600 mb-1">
               {room.participants.find(p => p.id === msg.senderId)?.name || 'Unknown'}
             </div>
          )}

          {/* Attachments */}
          {msg.attachments && msg.attachments.length > 0 && (
            <div className="mb-2 space-y-2">
              {msg.attachments.map(att => (
                <div key={att.id} className="flex items-center gap-3 p-2 bg-black/5 rounded-lg">
                  <div className="p-1.5 bg-white/20 rounded">
                     {att.type === 'image' ? <ImageIcon size={16} /> : <FileText size={16} />}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium truncate text-xs">{att.name}</div>
                    <div className="text-[10px] opacity-70">{att.size}</div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Content */}
          <div className="whitespace-pre-wrap leading-relaxed">{msg.content}</div>

          {/* Metadata */}
          <div className={cn(
            "flex items-center justify-end gap-1 mt-1.5 text-[10px]",
            isMe ? "text-primary-100" : "text-slate-400"
          )}>
            <span>{formatMessageDate(msg.timestamp)}</span>
            {isMe && (
              <span>
                {msg.status === 'sending' && <Clock size={10} />}
                {msg.status === 'sent' && <Check size={10} />}
                {msg.status === 'read' && <div className="flex"><Check size={10}/><Check size={10} className="-ml-1" /></div>}
              </span>
            )}
          </div>
        </div>
      </motion.div>
    );
  };

  return (
    <div className={cn("flex flex-col h-full bg-[#f0f2f5]", className)}>
      {/* Header */}
      <div className="bg-white/90 backdrop-blur-md p-3 sm:p-4 border-b border-slate-200 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-3">
          <button onClick={onBack} className="md:hidden p-2 text-slate-500 hover:bg-slate-100 rounded-full">
            <ArrowLeft size={20} />
          </button>
          
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-slate-200 flex items-center justify-center font-bold text-slate-500">
               {room.name.charAt(0)}
            </div>
            <div>
               <h3 className="font-bold text-slate-900 leading-tight">{room.name}</h3>
               <p className="text-xs text-slate-500 font-medium">
                 {room.type === 'private' 
                   ? 'Online' 
                   : `${room.participants.length} participants`
                 }
               </p>
            </div>
          </div>
        </div>

        <button className="p-2 text-slate-400 hover:text-slate-600 transition-colors">
           <MoreHorizontal size={20} />
        </button>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
        {messages.length === 0 ? (
          <div className="h-full flex flex-col items-center justify-center text-slate-400 opacity-60">
             <div className="w-20 h-20 bg-slate-200 rounded-full mb-4 flex items-center justify-center">
               <Send size={32} />
             </div>
             <p>No messages yet. Start the conversation!</p>
          </div>
        ) : (
          messages.map(msg => renderBubble(msg, msg.senderId === currentUser.id))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Area */}
      <div className="p-3 sm:p-4 bg-white border-t border-slate-200">
         <div className="max-w-4xl mx-auto flex items-end gap-2 bg-slate-100 p-2 rounded-3xl border border-transparent focus-within:bg-white focus-within:border-primary-200 focus-within:shadow-md transition-all">
            <button className="p-2 text-slate-400 hover:text-primary-600 transition-colors rounded-full hover:bg-slate-100 flex-shrink-0">
               <Paperclip size={20} />
            </button>
            <textarea 
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Type a message..."
              rows={1}
              className="flex-1 bg-transparent border-none focus:ring-0 resize-none py-2 px-1 text-sm sm:text-base max-h-32"
              style={{ minHeight: '40px' }}
            />
            <button 
              onClick={handleSendMessage}
              disabled={!inputValue.trim()}
              className={cn(
                "p-2.5 rounded-full transition-all duration-200 flex-shrink-0",
                inputValue.trim() 
                  ? "bg-primary-600 text-white shadow-lg hover:bg-primary-700 hover:scale-105 active:scale-95" 
                  : "bg-slate-200 text-slate-400 cursor-not-allowed"
              )}
            >
               <Send size={18} fill="currentColor" />
            </button>
         </div>
      </div>
    </div>
  );
};