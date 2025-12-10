import { useState } from 'react';
import { motion } from 'framer-motion';
import { Search, MoreVertical, MessageSquarePlus, User, Users, Megaphone, Inbox } from 'lucide-react';
import { cn } from "@/lib/utils";
import { useTranslation } from 'react-i18next';
import { Badge } from "@/components/ui/badge";

// Interface definitions adjusted to match existing types or generic enough
export interface ChatRoomDisplay {
  id: string;
  name: string;
  type: string;
  lastMessage?: {
      content: string;
      timestamp: string;
      senderId: string;
  };
  unreadCount: number;
  participants?: any[];
}

interface ChatSidebarProps {
  rooms: ChatRoomDisplay[];
  selectedRoomId: string | null;
  onSelectRoom: (roomId: string) => void;
  className?: string;
  isLoading?: boolean;
}

export const ChatSidebar: React.FC<ChatSidebarProps> = ({ 
    rooms, 
    selectedRoomId, 
    onSelectRoom, 
    className,
    isLoading = false 
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const { t } = useTranslation("common");

  const filteredRooms = rooms.filter(room => 
    room.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getRoomIcon = (type: string) => {
    switch(type) {
      case 'group': return <Users size={18} />;
      case 'channel': return <Megaphone size={18} />;
      default: return <User size={18} />;
    }
  };

  const getInitials = (name: string) => {
    return name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase();
  };

  const formatTime = (isoString?: string) => {
    if (!isoString) return '';
    const date = new Date(isoString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className={cn("flex flex-col h-full bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800", className)}>
      {/* Sidebar Header */}
      <div className="p-4 border-b border-slate-100 dark:border-slate-800 flex-shrink-0">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-slate-800 dark:text-slate-100">{t("chat.title", "Messages")}</h2>
          <div className="flex gap-2">
            <button className="p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full transition-colors">
              <MessageSquarePlus size={20} />
            </button>
            <button className="p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full transition-colors">
              <MoreVertical size={20} />
            </button>
          </div>
        </div>
        
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={16} />
          <input 
            type="text"
            placeholder={t("chat.search_placeholder", "Search conversations...")}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-9 pr-4 py-2.5 bg-slate-100 dark:bg-slate-800 border-transparent focus:bg-white dark:focus:bg-slate-900 focus:border-primary focus:ring-4 focus:ring-primary/10 rounded-xl text-sm transition-all outline-none"
          />
        </div>
      </div>

      {/* Room List */}
      <div className="flex-1 overflow-y-auto custom-scrollbar p-2 space-y-1">
        {isLoading ? (
            <div className="p-4 text-center text-slate-500 text-sm">{t("chat.loading_rooms", "Loading...")}</div>
        ) : filteredRooms.length === 0 ? (
            <div className="p-8 text-center text-slate-400 flex flex-col items-center">
                <Inbox size={32} className="mb-2 opacity-50" />
                <span className="text-sm">{t("chat.no_rooms", "No conversations found")}</span>
            </div>
        ) : (
            filteredRooms.map(room => (
            <motion.button
                key={room.id}
                onClick={() => onSelectRoom(room.id)}
                whileTap={{ scale: 0.98 }}
                className={cn(
                "w-full flex items-start gap-3 p-3 rounded-xl transition-all text-left group",
                selectedRoomId === room.id 
                    ? "bg-primary/10 shadow-sm" 
                    : "hover:bg-slate-50 dark:hover:bg-slate-800"
                )}
            >
                {/* Avatar */}
                <div className={cn(
                "w-12 h-12 rounded-full flex items-center justify-center text-white font-bold flex-shrink-0 shadow-sm transition-colors",
                room.type === 'channel' ? "bg-purple-500" : 
                room.type === 'group' ? "bg-indigo-500" : "bg-emerald-500"
                )}>
                {room.type === 'private' ? getInitials(room.name) : getRoomIcon(room.type)}
                </div>

                {/* Content */}
                <div className="flex-1 min-w-0">
                <div className="flex justify-between items-baseline mb-0.5">
                    <h3 className={cn(
                    "font-semibold truncate pr-2 transiiton-colors",
                    selectedRoomId === room.id ? "text-primary" : "text-slate-900 dark:text-slate-100"
                    )}>
                    {room.name}
                    </h3>
                    <span className="text-[10px] font-medium text-slate-400 flex-shrink-0">
                    {formatTime(room.lastMessage?.timestamp)}
                    </span>
                </div>
                <p className={cn(
                    "text-sm truncate leading-snug transition-colors",
                    room.unreadCount > 0 ? "font-semibold text-slate-700 dark:text-slate-300" : "text-slate-500 dark:text-slate-400",
                    selectedRoomId === room.id && "text-primary/80"
                )}>
                    {room.lastMessage?.senderId === 'me' && <span className="opacity-70">{t("chat.you", "You")}: </span>}
                    {room.lastMessage?.content || t("chat.no_messages", "No messages yet")}
                </p>
                </div>

                {/* Badges */}
                {room.unreadCount > 0 && (
                <div className="flex flex-col items-center justify-center h-12">
                    <span className="bg-primary text-primary-foreground text-[10px] font-bold px-1.5 py-0.5 min-w-[1.25rem] text-center rounded-full shadow-sm">
                    {room.unreadCount}
                    </span>
                </div>
                )}
            </motion.button>
            ))
        )}
      </div>
    </div>
  );
};
