import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { api } from "@/api/client";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChatSidebar, ChatRoomDisplay } from "@/features/chat/components/ChatSidebar";
import { ChatWindow, ChatMessageDisplay } from "@/features/chat/components/ChatWindow";
import { cn } from "@/lib/utils";
import {
  ChatAttachment,
  ChatMessage,
  ChatRoom,
  ChatMember,
  listMessages,
  listRooms,
  listRoomMembers,
  createMessage,
  uploadFile,
  markAsRead,
} from "@/features/chat/api";

export function ChatPage() {
  const { t } = useTranslation("common");
  const qc = useQueryClient();
  const [activeRoomId, setActiveRoomId] = useState("");
  
  // Queries
  const { data: me } = useQuery({
    queryKey: ["me"],
    queryFn: () => api("/me"),
  });

  const {
    data: rooms = [],
    isLoading: roomsLoading,
  } = useQuery<ChatRoom[]>({
    queryKey: ["chat", "rooms"],
    queryFn: listRooms,
  });

  const { data: messages = [], isLoading: messagesLoading } = useQuery<ChatMessage[]>({
    queryKey: ["chat", "messages", activeRoomId],
    queryFn: () => listMessages({ roomId: activeRoomId }),
    enabled: !!activeRoomId,
    refetchInterval: 10000, // Poll every 10 seconds for new messages
  });

  const { data: roomMembers = [] } = useQuery<ChatMember[]>({
    queryKey: ["chat", "members", activeRoomId],
    queryFn: () => listRoomMembers(activeRoomId),
    enabled: !!activeRoomId,
  });

  // Derived State
  const activeRoom = useMemo(
    () => rooms?.find((r) => r.id === activeRoomId) ?? null,
    [rooms, activeRoomId]
  );

  // Mappers
  const displayRooms: ChatRoomDisplay[] = useMemo(() => {
    return rooms.map(room => ({
        id: room.id,
        name: room.name,
        type: room.type,
        unreadCount: room.unread_count || 0,
        lastMessage: room.last_message_at ? {
            content: t("chat.message_placeholder", "Message"), // Backend doesn't provide preview yet
            timestamp: room.last_message_at,
            senderId: "unknown"
        } : undefined
    })).sort((a, b) => {
        const dateA = a.lastMessage?.timestamp || "";
        const dateB = b.lastMessage?.timestamp || "";
        return new Date(dateB).getTime() - new Date(dateA).getTime();
    });
  }, [rooms]);

  const displayMessages: ChatMessageDisplay[] = useMemo(() => {
      // Sort by time
      const sorted = [...messages].sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
      
      return sorted.map(msg => ({
          id: msg.id,
          senderId: msg.sender_id,
          senderName: msg.sender_name,
          senderRole: msg.sender_role,
          content: msg.body,
          timestamp: msg.created_at,
          status: 'read', // Simplified status for now
          attachments: msg.attachments?.map(att => ({
              name: att.name,
              url: att.url,
              type: att.type,
              size: att.size + ' B'
          }))
      }));
  }, [messages]);


  // Mutations
  const sendMutation = useMutation({
    mutationFn: (vars: { body: string; attachments: ChatAttachment[] }) =>
      createMessage(activeRoomId, vars.body, vars.attachments), // Changed from sendMessage to createMessage
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["chat", "messages", activeRoomId] });
      // Invalidate rooms to update last message time/preview if we had it
      qc.invalidateQueries({ queryKey: ["chat", "rooms"] }); 
    },
  });

  // Handlers
  const handleSendMessage = async (text: string, files: File[]) => {
      if (!activeRoomId) return;
      
      let uploadedAttachments: ChatAttachment[] = [];
      
      // Upload files first
      if (files.length > 0) {
          try {
             const uploadPromises = files.map(file => uploadFile(activeRoomId, file));
             uploadedAttachments = await Promise.all(uploadPromises);
          } catch (e) {
              console.error("Failed to upload files", e);
              // Handle error (toast)
              return;
          }
      }

      sendMutation.mutate({
          body: text,
          attachments: uploadedAttachments
      });
  };

  const handleSelectRoom = (roomId: string) => {
      setActiveRoomId(roomId);
  };

  const handleBack = () => {
      setActiveRoomId("");
  };

  // Mark read effect
  useEffect(() => {
    if (activeRoomId) {
      markAsRead(activeRoomId)
        .then(() => qc.invalidateQueries({ queryKey: ["chat", "rooms"] }));
    }
  }, [activeRoomId, qc]);

  return (
    <div className="fixed inset-0 pt-[57px] md:pt-[65px] bg-background">
      <div className="flex h-full border-t">
        {/* Sidebar - Hidden on mobile if room selected */}
        <div className={cn(
            "w-full md:w-80 lg:w-96 flex-shrink-0 h-full border-r",
            activeRoomId ? "hidden md:flex" : "flex"
        )}>
           <ChatSidebar 
              rooms={displayRooms}
              selectedRoomId={activeRoomId}
              onSelectRoom={handleSelectRoom}
              isLoading={roomsLoading}
              className="w-full"
              currentUser={me}
           />
        </div>

        {/* Chat Window - Hidden on mobile if no room selected */}
        <div className={cn(
            "flex-1 h-full flex flex-col min-w-0 bg-[#f0f2f5] dark:bg-slate-950",
            !activeRoomId ? "hidden md:flex" : "flex"
        )}>
            {activeRoom && activeRoomId ? (
                <ChatWindow 
                    room={{
                        id: activeRoom.id,
                        name: activeRoom.name,
                        type: activeRoom.type,
                        participants: roomMembers
                    }}
                    currentUser={me || { id: "", role: "" }}
                    messages={displayMessages}
                    onBack={handleBack}
                    onSendMessage={handleSendMessage}
                    isLoading={messagesLoading}
                />
            ) : (
                <div className="hidden md:flex flex-col items-center justify-center h-full text-slate-400">
                    <div className="w-24 h-24 bg-slate-100 dark:bg-slate-900 rounded-full mb-6 flex items-center justify-center animate-pulse">
                        <span className="text-4xl text-slate-300">💬</span>
                    </div>
                    <h2 className="text-xl font-bold text-slate-700 dark:text-slate-200 mb-2">{t("chat.welcome", "Welcome to Chat")}</h2>
                    <p className="text-sm max-w-xs text-center px-4">
                        {t("chat.select_room_prompt", "Select a conversation from the sidebar to start messaging.")}
                    </p>
                </div>
            )}
        </div>
      </div>
    </div>
  );
}
