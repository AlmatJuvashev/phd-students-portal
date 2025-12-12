import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowLeft, Send, Paperclip, MoreHorizontal, Image as ImageIcon, FileText, Check, Clock, CloudDownload, Archive, Reply, Trash2, Edit2, Info, Share, ExternalLink, X, CornerUpRight } from 'lucide-react';
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { archiveRoom, updateMessage, deleteMessage, createMessage, getRoomMembers } from "../api";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import { useTranslation } from 'react-i18next';
import { Button } from "@/components/ui/button";
import { AddMemberDialog } from "../AddMemberDialog";
import {
  DropdownMenu,
  DropdownItem,
} from "@/components/ui/dropdown-menu";
import { UserPlus } from "lucide-react";
import { ChatRoomDetails } from "./ChatRoomDetails";

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
    meta?: any;
    edited_at?: string | null;
    isForwarded?: boolean;
}

interface ChatWindowProps {
  room: {
      id: string;
      name: string;
      type: string;
      participants?: any[];
      is_archived?: boolean;
  };
  currentUser: {
      id: string;
      role?: string;
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
  const [isAddMemberOpen, setIsAddMemberOpen] = useState(false);
  const [editingMessage, setEditingMessage] = useState<ChatMessageDisplay | null>(null);
  const [forwardDialogOpen, setForwardDialogOpen] = useState(false);
  const [messageToForward, setMessageToForward] = useState<ChatMessageDisplay | null>(null);
  const [infoOpen, setInfoOpen] = useState(false);
  const [selectedMessageInfo, setSelectedMessageInfo] = useState<ChatMessageDisplay | null>(null);
  const [readers, setReaders] = useState<any[]>([]);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [showTyping, setShowTyping] = useState(false);
  const queryClient = useQueryClient();

  // Simulated Typing Effect
  useEffect(() => {
    // Only simulate typing if we are not the one typing and there are no messages or periodically
    if (messages.length > 0 && Math.random() > 0.7) {
       const timeout = setTimeout(() => {
          setShowTyping(true);
          setTimeout(() => setShowTyping(false), 3000);
       }, 5000);
       return () => clearTimeout(timeout);
    }
  }, [messages]);

  const archiveMutation = useMutation({
    mutationFn: (id: string) => archiveRoom(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chat", "rooms"] });
      alert(t("chat.archive_success", "Room archived."));
      // onBack(); // Optional: go back to list
    },
    onError: (err) => {
        alert(t("chat.archive_error", "Failed to archive room"));
        console.error(err);
    }
  });

  const handleArchiveRoom = () => {
      if (confirm(t("chat.confirm_archive", "Are you sure you want to archive this group?"))) {
          archiveMutation.mutate(room.id);
      }
  };

  const isAdmin = currentUser.role === 'admin' || currentUser.role === 'superadmin';

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, room.id, showTyping]); // Scroll when typing appears too

  const handleSendMessage = () => {
    if (!inputValue.trim() && selectedFiles.length === 0) return;
    onSendMessage(inputValue, selectedFiles);
    setInputValue('');
    setSelectedFiles([]);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      if (editingMessage) {
        handleUpdate();
      } else {
        handleSendMessage();
      }
    }
  };

  const cancelEdit = () => {
    setEditingMessage(null);
    setInputValue("");
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files) {
          setSelectedFiles(prev => [...prev, ...Array.from(e.target.files || [])]);
      }
  };

  const handleEditMessage = (msg: ChatMessageDisplay) => {
      setEditingMessage(msg);
      setInputValue(msg.content);
      fileInputRef.current?.focus();
  };

  const updateMutation = useMutation({
      mutationFn: (vars: {id: string, body: string}) => updateMessage(vars.id, vars.body),
      onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["chat", "messages"] });
          setEditingMessage(null);
          setInputValue("");
      }
  });

  const deleteMutation = useMutation({
      mutationFn: (id: string) => deleteMessage(id),
      onSuccess: () => queryClient.invalidateQueries({ queryKey: ["chat", "messages"] })
  });

  const handleUpdate = () => {
      if (editingMessage && inputValue.trim()) {
          updateMutation.mutate({ id: editingMessage.id, body: inputValue });
      }
  };

  const handleForward = async (roomTargetId: string) => {
      if (messageToForward) {
          try {
            await createMessage(roomTargetId, messageToForward.content, [], undefined, { forwarded_from: messageToForward.senderName });
            setForwardDialogOpen(false);
            setMessageToForward(null);
            alert(t("chat.forward_success", "Forwarded!"));
          } catch (e) {
            console.error("Failed to forward:", e);
            alert(t("chat.forward_error", "Failed to forward message."));
          }
      }
  };

  const handleInfo = async (msg: ChatMessageDisplay) => {
      setSelectedMessageInfo(msg);
      // Fetch members and check read status
      // Note: This is a simplified "Who could have read it" based on LastReadAt
      const members = await getRoomMembers(room.id);
      const msgTime = new Date(msg.timestamp).getTime();
      const readBy = members.filter((m: any) => m.last_read_at && new Date(m.last_read_at).getTime() >= msgTime && m.user_id !== currentUser.id);
      setReaders(readBy);
      setInfoOpen(true);
  };

  const renderUrlPreview = (text: string) => {
      const urlRegex = /(https?:\/\/[^\s]+)/g;
      const urls = text.match(urlRegex);
      if (!urls) return null;

      return (
          <div className="mt-2 space-y-2">
              {urls.map((url, i) => {
                  // YouTube
                  const ytMatch = url.match(/(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/);
                  if (ytMatch) {
                      return (
                          <div key={i} className="rounded-lg overflow-hidden relative pb-[56.25%] h-0 bg-black/10">
                              <iframe 
                                  src={`https://www.youtube.com/embed/${ytMatch[1]}`}
                                  className="absolute top-0 left-0 w-full h-full border-0"
                                  allowFullScreen
                                  title="YouTube preview"
                              />
                          </div>
                      );
                  }
                  // Image
                  if (url.match(/\.(jpeg|jpg|gif|png|webp)$/i)) {
                      return (
                          <img key={i} src={url} alt="Preview" className="rounded-lg max-h-48 object-cover mt-1" />
                      );
                  }
                  return null;
              })}
          </div>
      );
  };

  const formatMessageDate = (isoString: string) => {
    return new Date(isoString).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  const renderBubble = (msg: ChatMessageDisplay, isMe: boolean) => {
    return (
      <motion.div
        key={msg.id} 
        layout
        initial={{ opacity: 0, scale: 0.95, y: 10 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        transition={{ 
          opacity: { duration: 0.15 },
          layout: { type: "spring", bounce: 0.3 },
          type: "spring", stiffness: 500, damping: 30 
        }}
        className={cn(
          "flex w-full mb-3 px-2",
          isMe ? "justify-end" : "justify-start"
        )}
      >
        {/* Avatar for others */}
        {!isMe && room.type !== 'private' && (
          <div className="w-8 h-8 rounded-full border border-slate-100 overflow-hidden flex-shrink-0 mr-2 mt-1 shadow-sm">
             <img 
               src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${msg.senderName || 'Unknown'}`} 
               alt={msg.senderName}
               className="w-full h-full object-cover"
             />
          </div>
        )}
        
        <div className={cn(
          "max-w-[85%] sm:max-w-[70%] min-w-[60px] rounded-2xl relative group",
          isMe 
            ? "bg-primary text-primary-foreground rounded-br-sm shadow-sm" 
            : "bg-surface-0 dark:bg-slate-800 text-foreground border border-border/40 rounded-bl-sm shadow-sm"
        )}>
          {/* Elegant Context Menu Trigger */}
          <DropdownMenu position={isMe ? "right" : "left"} trigger={
            <button 
              className={cn(
                "absolute -right-1 top-1/2 -translate-y-1/2 p-1 rounded-full opacity-0 group-hover:opacity-100 transition-all duration-200 hover:scale-110",
                isMe 
                  ? "bg-white/20 hover:bg-white/30 text-white/80 hover:text-white" 
                  : "bg-slate-100 dark:bg-slate-700 hover:bg-slate-200 dark:hover:bg-slate-600 text-slate-500 dark:text-slate-300 shadow-sm"
              )}
            >
              <MoreHorizontal size={14} />
            </button>
          }>
              <DropdownItem onClick={() => { setMessageToForward(msg); setForwardDialogOpen(true); }}>
                <div className="flex items-center gap-2 py-0.5">
                  <CornerUpRight className="h-4 w-4 text-slate-500" />
                  <span>{t('chat.forward', 'Forward')}</span>
                </div>
              </DropdownItem>
            {isMe && (
              <>
                <DropdownItem onClick={() => handleEditMessage(msg)}>
                  <div className="flex items-center gap-2 py-0.5">
                    <Edit2 className="h-4 w-4 text-slate-500" />
                    <span>{t('chat.edit', 'Edit')}</span>
                  </div>
                </DropdownItem>
                <DropdownItem onClick={() => { if(confirm(t("chat.delete_confirm", "Delete this message?"))) deleteMutation.mutate(msg.id) }}>
                  <div className="flex items-center gap-2 py-0.5 text-red-500">
                    <Trash2 className="h-4 w-4" />
                    <span>{t('chat.delete', 'Delete')}</span>
                  </div>
                </DropdownItem>
              </>
            )}
            <DropdownItem onClick={() => handleInfo(msg)}>
              <div className="flex items-center gap-2 py-0.5">
                <Info className="h-4 w-4 text-slate-500" />
                <span>{t('chat.info', 'Info')}</span>
              </div>
            </DropdownItem>
          </DropdownMenu>

          {/* Message Content Container */}
          <div className={cn(
            "px-2.5 pt-1.5 pb-2 sm:px-3 sm:pb-2.5",
            !isMe && room.type !== 'private' && "pt-2"
          )}>
            {/* Forwarded indicator */}
            {msg.meta?.forwarded_from && (
              <div className={cn(
                "text-[11px] italic mb-1.5 flex items-center gap-1.5 pb-1.5 border-b",
                isMe ? "border-white/20 text-white/70" : "border-slate-200 dark:border-slate-600 text-slate-500"
              )}>
                <Reply size={12} className="-scale-x-100" />
                <span>{t("chat.forwarded_from", "Forwarded from {{name}}", { name: msg.meta.forwarded_from })}</span>
              </div>
            )}

            {/* Sender Name for Group Chats */}
            {!isMe && room.type !== 'private' && (
              <div className="flex items-center gap-1.5 mb-1">
                <span className="text-xs font-semibold text-primary">
                  {msg.senderName || t("chat.unknown_sender", "Unknown")}
                </span>
                <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400">
                  {msg.senderRole}
                </span>
              </div>
            )}

            {/* Attachments */}
            {msg.attachments && msg.attachments.length > 0 && (
              <div className="mb-2 space-y-1.5">
                {msg.attachments.map((att, idx) => (
                  <a 
                    key={idx} 
                    href={att.url} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className={cn(
                      "flex items-center gap-3 p-3 rounded-xl transition-all border",
                      isMe 
                        ? "bg-primary-foreground/10 border-primary-foreground/20 hover:bg-primary-foreground/20" 
                        : "bg-muted/50 border-border hover:bg-muted"
                    )}
                  >
                    <div className="w-9 h-9 rounded-lg flex items-center justify-center">
                      {att.type?.includes('image') ? <ImageIcon size={18} className={isMe ? "text-white" : "text-primary"} /> : <FileText size={18} className={isMe ? "text-white" : "text-primary"} />}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-xs truncate">{att.name}</p>
                      <p className={cn("text-[10px]", isMe ? "text-white/60" : "text-slate-400")}>{att.size}</p>
                    </div>
                    <CloudDownload size={16} className={cn(isMe ? "text-white/60" : "text-slate-400")} />
                  </a>
                ))}
              </div>
            )}

            {/* Message Content & Timestamp Combined */}
            <div className="text-[14px] sm:text-[15px] leading-snug whitespace-pre-wrap break-words relative">
              {msg.content}
              {renderUrlPreview(msg.content)}
              
              {/* Timestamp & Status - Floating */}
              <span className={cn(
                "float-right flex items-center gap-1 ml-2 mt-2 select-none text-[10px]",
                isMe ? "text-white/70" : "text-slate-400 dark:text-slate-500"
              )}>
                {msg.edited_at && <span className="italic mr-0.5">edited</span>}
                <span>{formatMessageDate(msg.timestamp)}</span>
                {isMe && (
                  <span className="flex items-center ml-0.5 inline-flex">
                    {msg.status === 'sending' && <Clock size={10} className="opacity-70" />}
                    {msg.status === 'sent' && <Check size={12} className="opacity-80" />}
                    {msg.status === 'read' && (
                      <div className="flex -space-x-1">
                        <Check size={12} className="text-blue-200" />
                        <Check size={12} className="text-blue-200" />
                      </div>
                    )}
                  </span>
                )}
              </span>
            </div>
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
             <div className="w-10 h-10 rounded-full bg-slate-100 overflow-hidden border border-slate-200">
               {room.type === 'private' ? (
                  <img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${room.name}`} alt={room.name} className="w-full h-full object-cover" />
               ) : (
                  <div className="w-full h-full flex items-center justify-center bg-indigo-100 text-indigo-600">
                     {room.name.charAt(0)}
                  </div>
               )}
            </div>
            <div>
               <h3 className="font-bold text-slate-900 dark:text-slate-100 leading-tight">{room.name}</h3>
               <p className="text-xs text-slate-500 dark:text-slate-400 font-medium">
                 {room.type === 'private' 
                   ? t('chat.private_room', 'Private') 
                   : t('chat.participants_count', '{{count}} participants', { count: room.participants?.length || 0 })
                 }
               </p>
            </div>
          </div>
        </div>

        <DropdownMenu
          trigger={
            <button className="p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors">
               <MoreHorizontal size={20} />
            </button>
          }
        >
            {isAdmin && room.type !== 'private' && (
              <DropdownItem onClick={() => setIsAddMemberOpen(true)}>
                <div className="flex items-center">
                  <UserPlus className="mr-2 h-4 w-4" />
                  {t("chat.add_member", "Add Member")}
                </div>
              </DropdownItem>
            )}
             {isAdmin && !room.is_archived && (
              <DropdownItem onClick={handleArchiveRoom}>
                 <div className="flex items-center text-red-500">
                   <Archive className="mr-2 h-4 w-4" />
                   {t("chat.archive_group", "Archive Group")}
                 </div>
              </DropdownItem>
            )}
             <DropdownItem onClick={() => setDetailsOpen(true)}>
                <div className="flex items-center">
                  {t("chat.room_details", "Group Info")}
                </div>
             </DropdownItem>
        </DropdownMenu>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
        {isLoading ? (
             <div className="h-full flex items-center justify-center text-slate-400">{t("chat.loading_messages", "Loading messages…")}</div>
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
        
        {/* Simulated Typing Indicator */}
        <AnimatePresence>
          {showTyping && (
             <motion.div 
               initial={{ opacity: 0, y: 10 }}
               animate={{ opacity: 1, y: 0 }}
               exit={{ opacity: 0, scale: 0.9 }}
               className="flex items-center gap-2 mb-4 px-2"
             >
                <div className="w-8 h-8 rounded-full bg-slate-100 border border-slate-200 overflow-hidden">
                   <img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=SimulatedUser`} alt="Typing..." />
                </div>
                <div className="bg-slate-200 dark:bg-slate-800 rounded-2xl p-3 rounded-bl-sm">
                   <div className="flex gap-1">
                      <motion.div 
                        className="w-1.5 h-1.5 bg-slate-400 rounded-full"
                        animate={{ y: [0, -3, 0] }}
                        transition={{ duration: 0.6, repeat: Infinity, delay: 0 }}
                      />
                      <motion.div 
                        className="w-1.5 h-1.5 bg-slate-400 rounded-full"
                        animate={{ y: [0, -3, 0] }}
                        transition={{ duration: 0.6, repeat: Infinity, delay: 0.2 }}
                      />
                      <motion.div 
                        className="w-1.5 h-1.5 bg-slate-400 rounded-full"
                        animate={{ y: [0, -3, 0] }}
                        transition={{ duration: 0.6, repeat: Infinity, delay: 0.4 }}
                      />
                   </div>
                </div>
             </motion.div>
          )}
        </AnimatePresence>
        <div ref={messagesEndRef} />
      </div>

      {/* Input Area */}
      <div className="p-3 sm:p-4 bg-white dark:bg-slate-900 border-t border-slate-200 dark:border-slate-800">
         
         {/* Editing Indicator */}
         {editingMessage && (
             <div className="flex items-center gap-3 bg-gradient-to-r from-primary/5 to-primary/10 border-l-4 border-primary px-3 py-2.5 mb-3 rounded-r-xl">
                 <div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center flex-shrink-0">
                   <Edit2 className="h-4 w-4 text-primary" />
                 </div>
                 <div className="flex-1 min-w-0">
                     <span className="font-semibold text-primary text-xs block">{t("chat.editing_message", "Editing message")}</span>
                     <span className="text-sm text-slate-600 dark:text-slate-300 truncate block">{editingMessage.content}</span>
                 </div>
                 <button onClick={cancelEdit} className="p-1.5 hover:bg-slate-200 dark:hover:bg-slate-700 rounded-full transition-colors flex-shrink-0">
                     <X className="h-4 w-4 text-slate-500" />
                 </button>
             </div>
         )}

         {/* Selected files preview */}
         {selectedFiles.length > 0 && (
             <div className="flex gap-2 mb-3 overflow-x-auto pb-1 scrollbar-thin">
                 {selectedFiles.map((f, i) => (
                     <div key={i} className="bg-gradient-to-r from-slate-50 to-slate-100 dark:from-slate-800 dark:to-slate-700 border border-slate-200 dark:border-slate-600 rounded-xl px-3 py-2 text-xs flex items-center gap-2 whitespace-nowrap shadow-sm hover:shadow transition-shadow">
                         <div className="w-6 h-6 rounded-lg bg-primary/10 flex items-center justify-center">
                           <Paperclip size={12} className="text-primary" />
                         </div>
                         <span className="font-medium text-slate-700 dark:text-slate-200 max-w-[120px] truncate">{f.name}</span>
                         <button 
                           onClick={() => setSelectedFiles(prev => prev.filter((_, idx) => idx !== i))} 
                           className="w-5 h-5 rounded-full bg-slate-200 dark:bg-slate-600 hover:bg-red-100 dark:hover:bg-red-900/50 flex items-center justify-center transition-colors group"
                         >
                           <X size={10} className="text-slate-500 group-hover:text-red-500" />
                         </button>
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
              placeholder={editingMessage ? t("chat.edit_placeholder", "Do changes...") : t("chat.type_message", "Type a message...")}
              rows={1}
              className="flex-1 bg-transparent border-none focus:ring-0 resize-none py-2 px-1 text-sm sm:text-base max-h-32 focus:outline-none dark:text-slate-200"
              style={{ minHeight: '40px' }}
            />
            <button 
              onClick={editingMessage ? handleUpdate : handleSendMessage}
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

      <AddMemberDialog 
        open={isAddMemberOpen} 
        onOpenChange={setIsAddMemberOpen} 
        roomId={room.id}
      />
      
      {/* Forward Dialog */}
      <Dialog open={forwardDialogOpen} onOpenChange={setForwardDialogOpen}>
          <DialogContent>
              <DialogHeader><DialogTitle>{t("chat.forward_message", "Forward Message")}</DialogTitle></DialogHeader>
              <div className="py-4">
                  <p className="mb-4 text-sm text-slate-500">{t("chat.select_forward", "Select group to forward to")}</p>
                  {/* Implementing room select adds complexity. For now, prompt ID or just alert */}
                   <div className="space-y-2">
                       {/* In real app, list rooms here. For demo, we just forward to SAME room to test? Or user input? */}
                      <p className="text-xs text-yellow-600">{t('chat.demo_forward_note', "Note: For this demo, clicking 'Forward' simulates forwarding to the same room.")}</p>
                      <Button onClick={() => handleForward(room.id)}>{t("chat.forward_here", "Forward here")}</Button>
                  </div>
              </div>
          </DialogContent>
      </Dialog>

      {/* Info Dialog */}
      <Dialog open={infoOpen} onOpenChange={setInfoOpen}>
          <DialogContent className="sm:max-w-md">
                         <DialogHeader>
                 <DialogTitle className="flex items-center gap-2 text-slate-900 dark:text-slate-100">
                   <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                     <Info className="h-4 w-4 text-primary" />
                   </div>
                   {t("chat.message_details", "Message Details")}
                 </DialogTitle>
               </DialogHeader>
               
               {/* Message Preview */}
               {selectedMessageInfo && (
                 <div className="bg-slate-50 dark:bg-slate-800/50 rounded-xl p-3 mb-4 border border-slate-200 dark:border-slate-700">
                   <p className="text-sm text-slate-600 dark:text-slate-300 line-clamp-2">{selectedMessageInfo.content}</p>
                   <p className="text-xs text-slate-400 mt-1">
                     {selectedMessageInfo.timestamp ? format(new Date(selectedMessageInfo.timestamp), 'MMM d, yyyy at HH:mm') : ''}
                   </p>
                 </div>
               )}

               {/* Read By Section */}
               <div>
                 <div className="flex items-center gap-2 mb-3">
                   <Check className="h-4 w-4 text-primary" />
                   <h4 className="text-sm font-semibold text-slate-700 dark:text-slate-200">{t("chat.read_by", "Read by")}</h4>
                   <span className="text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary font-medium">
                     {readers.length}
                   </span>
                 </div>
                 
                 <div className="max-h-48 overflow-y-auto space-y-1">
                   {readers.length === 0 ? (
                     <div className="flex flex-col items-center justify-center py-6 text-slate-400">
                       <Clock className="h-8 w-8 mb-2 opacity-50" />
                       <p className="text-sm">{t("chat.no_read", "No one has read this message yet")}</p>
                     </div>
                   ) : (
                     readers.map((r: any) => (
                       <div key={r.user_id} className="flex items-center gap-3 p-2.5 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                         <div className="w-8 h-8 rounded-full bg-gradient-to-br from-primary/20 to-primary/40 flex items-center justify-center text-xs font-bold text-primary flex-shrink-0">
                           {(r.first_name?.charAt(0) || r.email?.charAt(0) || '?').toUpperCase()}
                         </div>
                         <div className="flex-1 min-w-0">
                           <p className="text-sm font-medium text-slate-700 dark:text-slate-200 truncate">
                             {r.first_name ? `${r.first_name} ${r.last_name || ''}`.trim() : r.email}
                           </p>
                         </div>
                         <span className="text-xs text-slate-400 flex-shrink-0">
                           {r.last_read_at ? format(new Date(r.last_read_at), 'HH:mm') : ''}
                         </span>
                       </div>
                     ))
                   )}
                 </div>
               </div>
          </DialogContent>
      </Dialog>
      <ChatRoomDetails 
        open={detailsOpen} 
        onOpenChange={setDetailsOpen} 
        roomId={room.id}
        roomName={room.name}
        currentUser={currentUser}
        onArchive={isAdmin && !room.is_archived ? handleArchiveRoom : undefined}
      />
    </div>
  );
};
