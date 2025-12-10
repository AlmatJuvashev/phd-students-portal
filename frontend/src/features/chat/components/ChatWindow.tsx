import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowLeft, Send, Paperclip, MoreHorizontal, Image as ImageIcon, FileText, Check, Clock, CloudDownload } from 'lucide-react';
import { cn } from "@/lib/utils";
import { useTranslation } from 'react-i18next';
import { Button } from "@/components/ui/button";

// Interfaces to match backend or generic use
export interface ChatMessageDisplay {
    id: string;
    senderId: string;
    senderName?: string;
    senderRole?: string;
    content: string;
    timestamp: string;
    status?: 'sending' | 'sent' | 'read' | 'error';
    attachments?: {
        id?: string;
        name: string;
        url: string;
        type: string;
        size?: string;
    }[];
}

interface ChatWindowProps {
  room: {
      id: string;
      name: string;
      type: string;
      participants?: any[];
  };
  currentUser: {
      id: string;
  };
  messages: ChatMessageDisplay[];
  onBack: () => void; // Mobile only
  onSendMessage: (text: string, attachments: File[]) => void;
  className?: string;
  isLoading?: boolean;
}

export const ChatWindow: React.FC<ChatWindowProps> = ({ 
    room, 
    currentUser, 
    messages,
    onBack, 
    onSendMessage,
    className,
    isLoading 
}) => {
  const { t } = useTranslation("common");
  const [inputValue, setInputValue] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, room.id]);

  const handleSendMessage = () => {
    if (!inputValue.trim() && selectedFiles.length === 0) return;
    onSendMessage(inputValue, selectedFiles);
    setInputValue('');
    setSelectedFiles([]);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files) {
          setSelectedFiles(Array.from(e.target.files));
      }
  };

  const formatMessageDate = (isoString: string) => {
    return new Date(isoString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const renderBubble = (msg: ChatMessageDisplay, isMe: boolean) => {
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
            ? "bg-primary text-primary-foreground rounded-tr-none origin-bottom-right" 
            : "bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 text-slate-800 dark:text-slate-100 rounded-tl-none origin-bottom-left"
        )}>
          {/* Sender Name for Group Chats (Not me) */}
          {!isMe && room.type !== 'private' && (
             <div className="text-xs font-bold text-primary mb-1 opacity-90">
               {msg.senderName || 'Unknown'} <span className="text-[10px] font-normal opacity-70">({msg.senderRole})</span>
             </div>
          )}

          {/* Attachments */}
          {msg.attachments && msg.attachments.length > 0 && (
            <div className="mb-2 space-y-2">
              {msg.attachments.map((att, idx) => (
                <div key={idx} className="flex items-center gap-3 p-2 bg-black/5 dark:bg-white/10 rounded-lg">
                  <div className="p-1.5 bg-white/20 rounded">
                     {att.type?.includes('image') ? <ImageIcon size={16} /> : <FileText size={16} />}
                  </div>
                  <div className="flex-1 min-w-0">
                    <a href={att.url} target="_blank" rel="noopener noreferrer" className="font-medium truncate text-xs hover:underline block">
                        {att.name}
                    </a>
                  </div>
                  <a href={att.url} download className="p-1 hover:bg-black/10 rounded-full">
                      <CloudDownload size={14} />
                  </a>
                </div>
              ))}
            </div>
          )}

          {/* Content */}
          <div className="whitespace-pre-wrap leading-relaxed">{msg.content}</div>

          {/* Metadata */}
          <div className={cn(
            "flex items-center justify-end gap-1 mt-1.5 text-[10px]",
            isMe ? "text-primary-foreground/70" : "text-slate-400"
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
    <div className={cn("flex flex-col h-full bg-[#f0f2f5] dark:bg-slate-950", className)}>
      {/* Header */}
      <div className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-md p-3 sm:p-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-3">
          <button onClick={onBack} className="md:hidden p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full">
            <ArrowLeft size={20} />
          </button>
          
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-slate-200 dark:bg-slate-800 flex items-center justify-center font-bold text-slate-500 dark:text-slate-400">
               {room.name.charAt(0)}
            </div>
            <div>
               <h3 className="font-bold text-slate-900 dark:text-slate-100 leading-tight">{room.name}</h3>
               <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">
                 {room.type === 'private' 
                   ? 'Private' 
                   : `${room.participants?.length || 0} participants`
                 }
               </p>
            </div>
          </div>
        </div>

        <button className="p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors">
           <MoreHorizontal size={20} />
        </button>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
        {isLoading ? (
             <div className="h-full flex items-center justify-center text-slate-400">Loading messages...</div>
        ) : messages.length === 0 ? (
          <div className="h-full flex flex-col items-center justify-center text-slate-400 opacity-60">
             <div className="w-20 h-20 bg-slate-200 dark:bg-slate-800 rounded-full mb-4 flex items-center justify-center">
               <Send size={32} />
             </div>
             <p>{t("chat.no_messages", "No messages yet. Start the conversation!")}</p>
          </div>
        ) : (
          messages.map(msg => renderBubble(msg, msg.senderId === currentUser.id))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Area */}
      <div className="p-3 sm:p-4 bg-white dark:bg-slate-900 border-t border-slate-200 dark:border-slate-800">
         
         {/* Selected files preview */}
         {selectedFiles.length > 0 && (
             <div className="flex gap-2 mb-2 overflow-x-auto pb-2">
                 {selectedFiles.map((f, i) => (
                     <div key={i} className="bg-slate-100 dark:bg-slate-800 border rounded-lg p-2 text-xs flex items-center gap-2 whitespace-nowrap">
                         <Paperclip size={12} /> {f.name}
                         <button onClick={() => setSelectedFiles(prev => prev.filter((_, idx) => idx !== i))} className="hover:text-red-500"><ArrowLeft size={12} className="rotate-45" /></button>
                     </div>
                 ))}
             </div>
         )}

         <div className="max-w-4xl mx-auto flex items-end gap-2 bg-slate-100 dark:bg-slate-950 p-2 rounded-3xl border border-transparent focus-within:bg-white dark:focus-within:bg-slate-900 focus-within:border-primary/50 focus-within:shadow-md transition-all">
            <button 
                className="p-2 text-slate-400 hover:text-primary transition-colors rounded-full hover:bg-slate-200 dark:hover:bg-slate-800 flex-shrink-0"
                onClick={() => fileInputRef.current?.click()}
            >
               <Paperclip size={20} />
            </button>
            <input 
                type="file" 
                multiple 
                className="hidden" 
                ref={fileInputRef}
                onChange={handleFileSelect}
            />
            
            <textarea 
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={t("chat.type_message", "Type a message...")}
              rows={1}
              className="flex-1 bg-transparent border-none focus:ring-0 resize-none py-2 px-1 text-sm sm:text-base max-h-32 focus:outline-none dark:text-slate-200"
              style={{ minHeight: '40px' }}
            />
            <button 
              onClick={handleSendMessage}
              disabled={!inputValue.trim() && selectedFiles.length === 0}
              className={cn(
                "p-2.5 rounded-full transition-all duration-200 flex-shrink-0",
                (inputValue.trim() || selectedFiles.length > 0)
                  ? "bg-primary text-primary-foreground shadow-lg hover:bg-primary/90 hover:scale-105 active:scale-95" 
                  : "bg-slate-200 dark:bg-slate-800 text-slate-400 cursor-not-allowed"
              )}
            >
               <Send size={18} fill="currentColor" />
            </button>
         </div>
      </div>
    </div>
  );
};
