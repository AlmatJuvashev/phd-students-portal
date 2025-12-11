import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ArrowLeft, Send, Paperclip, MoreHorizontal, Image as ImageIcon, FileText, Check, Clock, CloudDownload, Archive, Reply, Trash2, Edit2, Info, Share, ExternalLink } from 'lucide-react';
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
  const queryClient = useQueryClient();

  const archiveMutation = useMutation({
    mutationFn: (id: string) => archiveRoom(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chat", "rooms"] });
      alert("Room archived.");
      // onBack(); // Optional: go back to list
    },
    onError: (err) => {
        alert("Failed to archive room");
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

  const handleForward = (roomTargetId: string) => {
      if (messageToForward) {
          createMessage(roomTargetId, messageToForward.content, [], undefined, { forwarded_from: messageToForward.senderName });
          setForwardDialogOpen(false);
          setMessageToForward(null);
          alert("Forwarded!");
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
          "max-w-[85%] sm:max-w-[70%] rounded-2xl p-3 sm:p-4 relative shadow-sm text-sm sm:text-base group",
          isMe 
            ? "bg-primary text-primary-foreground rounded-tr-none origin-bottom-right" 
            : "bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 text-slate-800 dark:text-slate-100 rounded-tl-none origin-bottom-left"
        )}>
          {/* Context Menu Hook (using Dropdown for simplicity) */}
          <DropdownMenu trigger={
              <button className="absolute top-2 right-2 p-1 text-slate-400 opacity-0 group-hover:opacity-100 transition-opacity bg-white/80 dark:bg-black/50 rounded-full">
                  <MoreHorizontal size={14} />
              </button>
          }>
               <DropdownItem onClick={() => { setMessageToForward(msg); setForwardDialogOpen(true); }}>
                   <div className="flex items-center"><Share className="mr-2 h-3 w-3"/> {t('Forward')}</div>
               </DropdownItem>
               {isMe && (
                   <>
                       <DropdownItem onClick={() => handleEditMessage(msg)}>
                           <div className="flex items-center"><Edit2 className="mr-2 h-3 w-3"/> {t('Edit')}</div>
                       </DropdownItem>
                       <DropdownItem onClick={() => { if(confirm("Delete?")) deleteMutation.mutate(msg.id) }}>
                           <div className="flex items-center text-red-500"><Trash2 className="mr-2 h-3 w-3"/> {t('Delete')}</div>
                       </DropdownItem>
                   </>
               )}
               <DropdownItem onClick={() => handleInfo(msg)}>
                   <div className="flex items-center"><Info className="mr-2 h-3 w-3"/> {t('Info')}</div>
               </DropdownItem>
          </DropdownMenu>

          {msg.meta?.forwarded_from && (
             <div className="text-[10px] italic opacity-70 mb-1 flex items-center gap-1">
                 <Reply size={10} className="-scale-x-100" /> Forwarded from {msg.meta.forwarded_from}
             </div>
          )}
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
          <div className="whitespace-pre-wrap leading-relaxed">
              {msg.content}
              {renderUrlPreview(msg.content)}
          </div>

          {/* Metadata */}
          <div className={cn(
            "flex items-center justify-end gap-1 mt-1.5 text-[10px] min-h-[14px]",
            isMe ? "text-primary-foreground/70" : "text-slate-400"
          )}>
            <span>{formatMessageDate(msg.timestamp)}</span>
            {msg.edited_at && <span>(edited)</span>}
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
                   {t("chat.archive_room", "Archive Group")}
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
              <DialogHeader><DialogTitle>Forward Message</DialogTitle></DialogHeader>
              <div className="py-4">
                  <p className="mb-4 text-sm text-slate-500">Select group to forward to (Feature simplified: auto-forward to current for demo, or implement room select)</p>
                  {/* Implementing room select adds complexity. For now, prompt ID or just alert */}
                  <div className="space-y-2">
                       {/* In real app, list rooms here. For demo, we just forward to SAME room to test? Or user input? */}
                      <p className="text-xs text-yellow-600">Note: For this demo, clicking 'Forward' simulates forwarding to the same room.</p>
                      <Button onClick={() => handleForward(room.id)}>Forward here</Button>
                  </div>
              </div>
          </DialogContent>
      </Dialog>

      {/* Info Dialog */}
      <Dialog open={infoOpen} onOpenChange={setInfoOpen}>
          <DialogContent>
               <DialogHeader><DialogTitle>Message Info</DialogTitle></DialogHeader>
               <div className="max-h-60 overflow-y-auto">
                   <h4 className="text-sm font-bold mb-2">Read by:</h4>
                   {readers.length === 0 ? <p className="text-sm text-slate-500">No one yet</p> : (
                       readers.map((r: any) => (
                           <div key={r.user_id} className="flex items-center justify-between py-2 border-b">
                               <span className="text-sm">{r.first_name || r.email}</span>
                               <span className="text-xs text-slate-400">{r.last_read_at ? format(new Date(r.last_read_at), 'HH:mm') : ''}</span>
                           </div>
                       ))
                   )}
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
