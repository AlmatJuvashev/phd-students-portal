import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  MessageCircle,
  ShieldCheck,
  Sparkles,
  Clock4,
  Send,
  Inbox,
  Search,
  ArrowLeft,
  MessageSquareOff,
  Paperclip,
  X,
  File as FileIcon,
  MoreVertical,
  Pencil,
  Trash2,
  UserPlus,
} from "lucide-react";
import { AddMemberDialog } from "@/features/chat/AddMemberDialog";
import { api } from "@/api/client";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu-radix";
import { cn } from "@/lib/utils";
import {
  ChatAttachment,
  ChatMessage,
  ChatRoom,
  ChatMember,
  listMessages,
  listRooms,
  listRoomMembers,
  sendMessage,
  uploadFile,
  updateMessage,
  deleteMessage,
  markAsRead,
} from "@/features/chat/api";

function formatTime(iso?: string) {
  if (!iso) return "";
  const date = new Date(iso);
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

function formatDateLabel(date: Date) {
  const today = new Date();
  const yesterday = new Date();
  yesterday.setDate(today.getDate() - 1);
  if (date.toDateString() === today.toDateString()) return "Сегодня";
  if (date.toDateString() === yesterday.toDateString()) return "Вчера";
  return date.toLocaleDateString("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

function initials(name?: string) {
  if (!name) return "?";
  const parts = name.split(" ").filter(Boolean);
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
  return name.slice(0, 2).toUpperCase();
}

export function ChatPage() {
  const { t } = useTranslation("common");
  const qc = useQueryClient();
  const [activeRoomId, setActiveRoomId] = useState("");
  const [draft, setDraft] = useState("");
  const [attachments, setAttachments] = useState<ChatAttachment[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [editingMessage, setEditingMessage] = useState<ChatMessage | null>(
    null
  );
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [search, setSearch] = useState("");
  const [showAddMember, setShowAddMember] = useState(false);

  const { data: me } = useQuery({
    queryKey: ["me"],
    queryFn: () => api("/me"),
  });
  const isAdmin = me?.role === "admin" || me?.role === "superadmin";

  // ... (existing queries)

  const updateMutation = useMutation({
    mutationFn: (vars: { id: string; body: string }) =>
      updateMessage(vars.id, vars.body),
    onSuccess: () => {
      setDraft("");
      setEditingMessage(null);
      qc.invalidateQueries({ queryKey: ["chat", "messages", activeRoomId] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteMessage(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["chat", "messages", activeRoomId] });
    },
  });

  const sendMutation = useMutation({
    mutationFn: (vars: { body: string; attachments: ChatAttachment[] }) =>
      sendMessage(activeRoomId, vars.body, vars.attachments),
    onSuccess: () => {
      setDraft("");
      setAttachments([]);
      qc.invalidateQueries({ queryKey: ["chat", "messages", activeRoomId] });
    },
  });

  const handleSend = () => {
    if (!activeRoomId) return;

    if (editingMessage) {
      if (!draft.trim()) return; // Cannot have empty message
      updateMutation.mutate({ id: editingMessage.id, body: draft.trim() });
      return;
    }

    if (!draft.trim() && attachments.length === 0) return;
    sendMutation.mutate({ body: draft.trim(), attachments });
  };

  const startEditing = (msg: ChatMessage) => {
    setEditingMessage(msg);
    setDraft(msg.body);
    // Focus textarea? (handled by auto-focus if we re-render, or ref)
  };

  const cancelEditing = () => {
    setEditingMessage(null);
    setDraft("");
  };

  const {
    data: rooms = [],
    isLoading: roomsLoading,
    isFetching: roomsFetching,
  } = useQuery<ChatRoom[]>({
    queryKey: ["chat", "rooms"],
    queryFn: listRooms,
  });

  const roomsCount = rooms?.length ?? 0;

  // Filter rooms by search and archive status
  const filteredRooms = useMemo(() => {
    if (!rooms) return [];
    const searchLower = search.toLowerCase();
    return rooms.filter((r) => r.name.toLowerCase().includes(searchLower));
  }, [rooms, search]);

  const activeRooms = useMemo(
    () => filteredRooms.filter((r) => !r.is_archived),
    [filteredRooms]
  );

  const archivedRooms = useMemo(
    () => filteredRooms.filter((r) => r.is_archived),
    [filteredRooms]
  );

  const activeRoom = useMemo(
    () => rooms?.find((r) => r.id === activeRoomId) ?? null,
    [rooms, activeRoomId]
  );

  // Messages query
  const { data: messages = [], isLoading: messagesLoading } = useQuery<
    ChatMessage[]
  >({
    queryKey: ["chat", "messages", activeRoomId],
    queryFn: () => listMessages({ roomId: activeRoomId }),
    enabled: !!activeRoomId,
  });

  // Room members query
  const { data: roomMembers = [] } = useQuery<ChatMember[]>({
    queryKey: ["chat", "members", activeRoomId],
    queryFn: () => listRoomMembers(activeRoomId),
    enabled: !!activeRoomId,
  });

  // Format members for display in header
  const membersDisplay = useMemo(() => {
    if (!roomMembers || roomMembers.length === 0) return "";
    const names = roomMembers
      .slice(0, 3)
      .map((m) => {
        const name = [m.first_name, m.last_name].filter(Boolean).join(" ");
        return name || m.email || "User";
      });
    if (roomMembers.length > 3) {
      return `${names.join(", ")} +${roomMembers.length - 3}`;
    }
    return names.join(", ");
  }, [roomMembers]);

  // Sort messages by created_at
  const sortedMessages = useMemo(() => {
    if (!messages) return [];
    return [...messages].sort(
      (a, b) =>
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    );
  }, [messages]);

  // Disable input if no active room or room is archived
  const disableInput = !activeRoomId || activeRoom?.is_archived;

  // Handle file selection for upload
  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    console.log("[handleFileSelect] Files selected:", files);
    console.log("[handleFileSelect] Active room:", activeRoomId);
    
    if (!files || files.length === 0 || !activeRoomId) {
      console.log("[handleFileSelect] Early return:", { files: !!files, length: files?.length, activeRoomId });
      return;
    }

    setIsUploading(true);
    try {
      const file = files[0];
      console.log("[handleFileSelect] Uploading file:", {
        name: file.name,
        size: file.size,
        type: file.type,
      });
      const uploaded = await uploadFile(activeRoomId, file);
      console.log("[handleFileSelect] Upload successful:", uploaded);
      setAttachments((prev) => [...prev, uploaded]);
    } catch (err) {
      console.error("[handleFileSelect] File upload failed:", err);
    } finally {
      setIsUploading(false);
      e.target.value = ""; // Reset input
    }
  };

  // Mobile-specific logic
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 768);
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  // On mobile, if we select a room, we want to show the chat (hide sidebar)
  // If we don't have a room, we show sidebar.
  // We reuse sidebarOpen to mean "Show Room List" on mobile.
  useEffect(() => {
    if (isMobile) {
      if (activeRoomId) {
        setSidebarOpen(false); // Show chat
      } else {
        setSidebarOpen(true); // Show list
      }
    } else {
      setSidebarOpen(true); // Always show sidebar on desktop
    }
  }, [activeRoomId, isMobile]);

  // Mark room as read when entering and refresh rooms to update unread counts
  useEffect(() => {
    if (activeRoomId) {
      markAsRead(activeRoomId)
        .then(() => {
          // Refresh room list to update unread counts
          qc.invalidateQueries({ queryKey: ["chat", "rooms"] });
        })
        .catch((err) => console.error("Failed to mark room as read:", err));
    }
  }, [activeRoomId, qc]);

  const handleBackToRooms = () => {
    setSidebarOpen(true);
    setActiveRoomId(""); // Clear active room to ensure we go back to "no selection" state conceptually, though UI just shows list
  };

  return (
    <div className="space-y-4 h-[calc(100vh-100px)] flex flex-col p-4 md:p-6">
      <div className="space-y-1 shrink-0">
        <div className="flex items-center gap-2">
          <MessageCircle className="h-6 w-6 text-primary" />
          <h1 className="text-2xl font-bold tracking-tight">
            {t("chat.title")}
          </h1>
          <Badge variant="outline" className="text-xs">
            {t("chat.preview")}
          </Badge>
        </div>
        <p className="text-sm text-muted-foreground">{t("chat.subtitle")}</p>
      </div>

      <div className="flex-1 flex min-h-0 rounded-xl border shadow-sm overflow-hidden bg-background relative">
        {/* Sidebar */}
        <div
          className={cn(
            "w-full md:w-80 lg:w-96 border-r bg-muted/30 flex flex-col absolute inset-0 z-20 md:relative md:z-0 bg-background md:bg-muted/30 transition-transform duration-300",
            isMobile
              ? sidebarOpen
                ? "translate-x-0"
                : "-translate-x-full"
              : "translate-x-0"
          )}
        >
          <div className="p-4 border-b bg-gradient-to-b from-primary/5 to-transparent space-y-3 shrink-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="p-2 rounded-lg bg-primary/10">
                  <MessageCircle className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <div className="font-semibold">{t("chat.rooms_heading")}</div>
                  <div className="text-xs text-muted-foreground">
                    {roomsFetching ? "…" : roomsCount}{" "}
                    {t("chat.room_members", { count: roomsCount })}
                  </div>
                </div>
              </div>
            </div>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder={t("chat.search_placeholder", {
                  defaultValue: "Поиск чатов...",
                })}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          <div className="p-3 space-y-4 overflow-y-auto flex-1">
            {roomsLoading ? (
              <div className="rounded-md border border-dashed p-3 text-sm text-muted-foreground">
                {t("chat.loading_rooms")}
              </div>
            ) : roomsCount === 0 ? (
              <div className="rounded-md border border-dashed p-3 text-sm text-muted-foreground flex items-center gap-2">
                <Inbox className="h-4 w-4" />
                <span>{t("chat.no_rooms")}</span>
              </div>
            ) : (
              <>
                <div>
                  <div className="text-[11px] uppercase tracking-wide text-muted-foreground px-2 mb-1">
                    {t("chat.active_rooms", { defaultValue: "Мои группы" })}
                  </div>
                  <div className="space-y-1">
                    {activeRooms.length ? (
                      activeRooms.map((room) => (
                        <button
                          key={room.id}
                          onClick={() => {
                            setActiveRoomId(room.id);
                            // Sidebar closing handled by effect
                          }}
                          className={cn(
                            "w-full rounded-lg border px-3 py-2 text-left hover:border-primary/60 hover:bg-primary/5 transition",
                            room.id === activeRoomId &&
                              "border-primary bg-primary/5 shadow-sm"
                          )}
                        >
                          <div className="flex items-center justify-between gap-2">
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center gap-2">
                                <span className="font-semibold leading-tight truncate">
                                  {room.name}
                                </span>
                                {/* Unread badge */}
                                {room.unread_count && room.unread_count > 0 && (
                                  <span className="inline-flex items-center justify-center h-5 min-w-5 px-1.5 text-[10px] font-bold rounded-full bg-primary text-primary-foreground shrink-0">
                                    {room.unread_count > 99 ? "99+" : room.unread_count}
                                  </span>
                                )}
                              </div>
                              <div className="flex flex-wrap items-center gap-2 text-[11px] text-muted-foreground mt-1">
                                <Badge
                                  variant="outline"
                                  className="text-[10px]"
                                >
                                  {t(`chat.types.${room.type}`)}
                                </Badge>
                                {room.created_by_role && (
                                  <Badge
                                    variant="outline"
                                    className="text-[10px]"
                                  >
                                    {room.created_by_role}
                                  </Badge>
                                )}
                              </div>
                            </div>
                            <span className="text-[11px] text-muted-foreground flex items-center gap-1 shrink-0">
                              <Clock4 className="h-3 w-3" />
                              {formatTime(room.last_message_at || room.created_at)}
                            </span>
                          </div>
                        </button>
                      ))
                    ) : (
                      <div className="text-xs text-muted-foreground px-2 py-4">
                        {t("chat.no_rooms")}
                      </div>
                    )}
                  </div>
                </div>

                {archivedRooms.length > 0 && (
                  <div className="pt-2 border-t">
                    <div className="text-[11px] uppercase tracking-wide text-muted-foreground px-2 mb-1">
                      {t("chat.archived")}
                    </div>
                    <div className="space-y-1">
                      {archivedRooms.map((room) => (
                        <button
                          key={room.id}
                          onClick={() => {
                            setActiveRoomId(room.id);
                            // Sidebar closing handled by effect
                          }}
                          className={cn(
                            "w-full rounded-lg border px-3 py-2 text-left hover:border-primary/40 hover:bg-primary/5 transition",
                            room.id === activeRoomId &&
                              "border-primary bg-primary/5 shadow-sm"
                          )}
                        >
                          <div className="flex items-center justify-between gap-2">
                            <div>
                              <div className="font-semibold leading-tight">
                                {room.name}
                              </div>
                              <div className="flex flex-wrap items-center gap-2 text-[11px] text-muted-foreground mt-1">
                                <Badge
                                  variant="outline"
                                  className="text-[10px]"
                                >
                                  {t(`chat.types.${room.type}`)}
                                </Badge>
                                <Badge
                                  variant="outline"
                                  className="text-[10px]"
                                >
                                  {t("chat.archived")}
                                </Badge>
                              </div>
                            </div>
                            <span className="text-[11px] text-muted-foreground flex items-center gap-1">
                              <Clock4 className="h-3 w-3" />
                              {formatTime(room.created_at)}
                            </span>
                          </div>
                        </button>
                      ))}
                    </div>
                  </div>
                )}
              </>
            )}
          </div>
        </div>

        {/* Messages */}
        <div className="flex-1 flex flex-col min-w-0 w-full h-full absolute inset-0 md:static bg-background z-10 md:z-0">
          {/* Enhanced Header */}
          <div className="relative shrink-0 z-20">
            {/* Glassmorphism background with gradient */}
            <div className="absolute inset-0 bg-background/80 backdrop-blur-md border-b z-0">
              <div className="absolute inset-0 bg-gradient-to-r from-primary/5 via-transparent to-primary/5 opacity-50" />
            </div>

            <div className="relative z-10 flex items-center justify-between px-4 py-3">
              <div className="flex items-center gap-3 min-w-0">
                <Button
                  variant="ghost"
                  size="icon"
                  className="md:hidden -ml-2 h-9 w-9 rounded-full hover:bg-primary/10"
                  onClick={handleBackToRooms}
                >
                  <ArrowLeft className="h-5 w-5 text-primary" />
                </Button>
                
                {/* Room Avatar/Icon */}
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-primary/10 to-primary/5 flex items-center justify-center shrink-0 border border-primary/10 shadow-sm">
                  {activeRoom ? (
                     <span className="text-sm font-bold text-primary">
                       {activeRoom.name.substring(0, 2).toUpperCase()}
                     </span>
                  ) : (
                    <Sparkles className="h-5 w-5 text-primary" />
                  )}
                </div>

                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="font-bold text-base tracking-tight truncate">
                      {activeRoom ? activeRoom.name : t("chat.messages_heading")}
                    </span>
                    {activeRoom?.is_archived && (
                      <Badge variant="secondary" className="text-[10px] px-1.5 h-5">
                        {t("chat.archived")}
                      </Badge>
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground truncate flex items-center gap-1.5">
                    {activeRoom ? (
                      <>
                        <span className="inline-block w-1.5 h-1.5 rounded-full bg-green-500/50 animate-pulse" />
                        {membersDisplay ? (
                          <span title={roomMembers.map(m => [m.first_name, m.last_name].filter(Boolean).join(" ") || m.email).join(", ")}>
                            {membersDisplay}
                          </span>
                        ) : (
                          t(`chat.types.${activeRoom.type}`)
                        )}
                      </>
                    ) : (
                      t("chat.empty_state")
                    )}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-1">
                {activeRoom && isAdmin && !activeRoom.is_archived && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-9 w-9 rounded-full hover:bg-primary/10 text-muted-foreground hover:text-primary transition-colors"
                    onClick={() => setShowAddMember(true)}
                    title={t("chat.add_member", { defaultValue: "Add Member" })}
                  >
                    <UserPlus className="h-5 w-5" />
                  </Button>
                )}
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-9 w-9 rounded-full hover:bg-primary/10 text-muted-foreground hover:text-primary transition-colors"
                >
                  <MoreVertical className="h-5 w-5" />
                </Button>
              </div>
            </div>
          </div>

          {activeRoomId && (
            <AddMemberDialog
              open={showAddMember}
              onOpenChange={setShowAddMember}
              roomId={activeRoomId}
            />
          )}

          <div className="flex-1 flex flex-col bg-muted/10 min-h-0">
            {messagesLoading ? (
              <div className="flex-1 flex items-center justify-center text-sm text-muted-foreground">
                {t("chat.loading_messages")}
              </div>
            ) : sortedMessages.length === 0 ? (
              <div className="flex-1 flex flex-col items-center justify-center gap-3 text-center px-6">
                <div className="p-4 rounded-full bg-muted">
                  <MessageSquareOff className="h-6 w-6 text-muted-foreground" />
                </div>
                <div className="text-sm text-muted-foreground max-w-sm">
                  {t("chat.empty_state")}
                </div>
              </div>
            ) : (
              <div className="flex-1 overflow-y-auto px-4 py-3 space-y-3">
                {sortedMessages.map((msg, idx) => {
                  const prev = sortedMessages[idx - 1];
                  const thisDate = new Date(msg.created_at);
                  const showDate =
                    !prev ||
                    new Date(prev.created_at).toDateString() !==
                      thisDate.toDateString();
                  const showSender =
                    !prev ||
                    prev.sender_id !== msg.sender_id ||
                    thisDate.getTime() - new Date(prev.created_at).getTime() >
                      5 * 60 * 1000;
                  
                  // WhatsApp-style: check if this is the current user's message
                  const isOwnMessage = msg.sender_id === me?.id;
                  
                  // Role-based styling for admin/advisor messages
                  const isAdminMessage = msg.sender_role === "admin" || msg.sender_role === "superadmin";
                  const isAdvisorMessage = msg.sender_role === "advisor" || msg.sender_role === "chair";
                  
                  return (
                    <div
                      key={msg.id}
                      className={cn(
                        "space-y-1 group",
                        editingMessage?.id === msg.id && "opacity-70"
                      )}
                    >
                      {showDate && (
                        <div className="flex items-center justify-center my-2">
                          <span className="text-[11px] text-muted-foreground bg-muted px-3 py-1 rounded-full">
                            {formatDateLabel(thisDate)}
                          </span>
                        </div>
                      )}
                      <div className={cn(
                        "flex items-start gap-2",
                        isOwnMessage && "flex-row-reverse"
                      )}>
                        {/* Avatar - hide for own messages or show on right side */}
                        <div className={cn(
                          "h-8 w-8 rounded-full flex items-center justify-center text-xs font-semibold shrink-0",
                          isOwnMessage 
                            ? "bg-primary text-primary-foreground" 
                            : isAdminMessage 
                              ? "bg-purple-500/20 text-purple-600 dark:text-purple-400"
                              : isAdvisorMessage
                                ? "bg-blue-500/20 text-blue-600 dark:text-blue-400"
                                : "bg-primary/10 text-primary"
                        )}>
                          {initials(msg.sender_name || msg.sender_id)}
                        </div>
                        <div className={cn(
                          "flex-1 min-w-0 max-w-[75%]",
                          isOwnMessage && "flex flex-col items-end"
                        )}>
                          {showSender && !isOwnMessage && (
                            <div className="flex items-center gap-2 text-[11px] text-muted-foreground">
                              <span className={cn(
                                "font-medium",
                                isAdminMessage && "text-purple-600 dark:text-purple-400",
                                isAdvisorMessage && "text-blue-600 dark:text-blue-400",
                                !isAdminMessage && !isAdvisorMessage && "text-foreground"
                              )}>
                                {msg.sender_name || msg.sender_id}
                              </span>
                              {msg.sender_role && (
                                <Badge
                                  variant="outline"
                                  className={cn(
                                    "text-[10px]",
                                    isAdminMessage && "border-purple-500/50 text-purple-600 dark:text-purple-400",
                                    isAdvisorMessage && "border-blue-500/50 text-blue-600 dark:text-blue-400"
                                  )}
                                >
                                  {msg.sender_role}
                                </Badge>
                              )}
                              <span>{formatTime(msg.created_at)}</span>
                            </div>
                          )}
                          {/* Show time for own messages */}
                          {showSender && isOwnMessage && (
                            <div className="text-[11px] text-muted-foreground mb-1">
                              {formatTime(msg.created_at)}
                            </div>
                          )}

                          {msg.deleted_at ? (
                            <div className="mt-1 rounded-2xl bg-muted/50 border px-3 py-2 shadow-sm italic text-muted-foreground text-sm flex items-center gap-2">
                              <MessageSquareOff className="h-3 w-3" />
                              {t("chat.message_deleted", {
                                defaultValue: "This message was deleted",
                              })}
                            </div>
                          ) : (
                            <div className="group/msg relative mt-1">
                              <div className={cn(
                                "rounded-2xl px-3 py-2 shadow-sm",
                                isOwnMessage 
                                  ? "bg-primary text-primary-foreground rounded-tr-sm" 
                                  : cn(
                                      "bg-background border rounded-tl-sm",
                                      isAdminMessage && "border-purple-500/30 bg-purple-50/50 dark:bg-purple-950/20",
                                      isAdvisorMessage && "border-blue-500/30 bg-blue-50/50 dark:bg-blue-950/20"
                                    )
                              )}>
                                {msg.attachments &&
                                  msg.attachments.length > 0 && (
                                    <div className="mb-2 space-y-2">
                                      {msg.attachments.map((att, i) => (
                                        <div key={i}>
                                          {att.type.startsWith("image/") ? (
                                            <a
                                              href={att.url}
                                              target="_blank"
                                              rel="noopener noreferrer"
                                              className="block"
                                            >
                                              <img
                                                src={att.url}
                                                alt={att.name}
                                                className="max-w-full rounded-md max-h-60 object-cover border"
                                              />
                                            </a>
                                          ) : (
                                            <a
                                              href={att.url}
                                              target="_blank"
                                              rel="noopener noreferrer"
                                              className={cn(
                                                "flex items-center gap-2 p-2 rounded-md transition border",
                                                isOwnMessage 
                                                  ? "bg-primary-foreground/10 hover:bg-primary-foreground/20 border-primary-foreground/20" 
                                                  : "bg-muted/50 hover:bg-muted"
                                              )}
                                            >
                                              <FileIcon className={cn(
                                                "h-4 w-4",
                                                isOwnMessage ? "text-primary-foreground" : "text-primary"
                                              )} />
                                              <div className="flex-1 min-w-0">
                                                <div className="text-sm font-medium truncate">
                                                  {att.name}
                                                </div>
                                                <div className={cn(
                                                  "text-xs",
                                                  isOwnMessage ? "text-primary-foreground/70" : "text-muted-foreground"
                                                )}>
                                                  {(att.size / 1024).toFixed(1)}{" "}
                                                  KB
                                                </div>
                                              </div>
                                            </a>
                                          )}
                                        </div>
                                      ))}
                                    </div>
                                  )}
                                <p className="text-sm whitespace-pre-line leading-relaxed break-words">
                                  {msg.body}
                                </p>
                                {msg.edited_at && (
                                  <span className={cn(
                                    "text-[10px] italic ml-1",
                                    isOwnMessage ? "text-primary-foreground/70" : "text-muted-foreground"
                                  )}>
                                    (
                                    {t("chat.edited", {
                                      defaultValue: "edited",
                                    })}
                                    )
                                  </span>
                                )}
                              </div>

                              {/* Only show edit/delete menu for own messages */}
                              {isOwnMessage && (
                                <div className={cn(
                                  "absolute top-1 opacity-0 group-hover/msg:opacity-100 transition-opacity",
                                  isOwnMessage ? "left-1" : "right-1"
                                )}>
                                  <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                      <Button
                                        variant="ghost"
                                        size="icon"
                                        className="h-6 w-6 rounded-full bg-background/80 hover:bg-muted shadow-sm"
                                      >
                                        <MoreVertical className="h-3 w-3" />
                                      </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align={isOwnMessage ? "start" : "end"}>
                                      <DropdownMenuItem
                                        onClick={() => startEditing(msg)}
                                      >
                                        <Pencil className="mr-2 h-3 w-3" />
                                        {t("common.edit", {
                                          defaultValue: "Edit",
                                        })}
                                      </DropdownMenuItem>
                                      <DropdownMenuItem
                                        className="text-destructive focus:text-destructive"
                                        onClick={() =>
                                          deleteMutation.mutate(msg.id)
                                        }
                                      >
                                        <Trash2 className="mr-2 h-3 w-3" />
                                        {t("common.delete", {
                                          defaultValue: "Delete",
                                        })}
                                      </DropdownMenuItem>
                                    </DropdownMenuContent>
                                  </DropdownMenu>
                                </div>
                              )}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          <div className="border-t p-3 bg-background space-y-3 shrink-0">
            {editingMessage && (
              <div className="flex items-center justify-between bg-muted/50 p-2 rounded-md border-l-4 border-primary text-sm">
                <div className="flex flex-col">
                  <span className="font-semibold text-primary">
                    {t("chat.editing", { defaultValue: "Editing message" })}
                  </span>
                  <span className="text-muted-foreground line-clamp-1">
                    {editingMessage.body}
                  </span>
                </div>
                <Button variant="ghost" size="icon" onClick={cancelEditing}>
                  <X className="h-4 w-4" />
                </Button>
              </div>
            )}
            {attachments.length > 0 && (
              <div className="flex gap-2 overflow-x-auto pb-2">
                {attachments.map((att, i) => (
                  <div key={i} className="relative group flex-shrink-0">
                    <div className="w-16 h-16 rounded-md border bg-muted flex items-center justify-center overflow-hidden">
                      {att.type.startsWith("image/") ? (
                        <img
                          src={att.url}
                          alt={att.name}
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <FileIcon className="h-6 w-6 text-muted-foreground" />
                      )}
                    </div>
                    <button
                      onClick={() =>
                        setAttachments((prev) =>
                          prev.filter((_, idx) => idx !== i)
                        )
                      }
                      className="absolute -top-1 -right-1 bg-destructive text-destructive-foreground rounded-full p-0.5 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <X className="h-3 w-3" />
                    </button>
                  </div>
                ))}
              </div>
            )}
            <div className="flex gap-2">
              <input
                type="file"
                id="file-upload"
                className="hidden"
                onChange={handleFileSelect}
                disabled={disableInput || isUploading || !!editingMessage}
              />
              <Button
                variant="ghost"
                size="icon"
                className="shrink-0"
                disabled={disableInput || isUploading || !!editingMessage}
                onClick={() => document.getElementById("file-upload")?.click()}
              >
                <Paperclip className="h-5 w-5 text-muted-foreground" />
              </Button>
              <Textarea
                placeholder={t("chat.input_placeholder")}
                className="resize-none min-h-[40px] max-h-[120px]"
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    handleSend();
                  }
                }}
                disabled={
                  disableInput ||
                  sendMutation.isPending ||
                  updateMutation.isPending ||
                  (!activeRoomId && !editingMessage)
                }
              />
              <Button
                size="icon"
                className="shrink-0"
                onClick={handleSend}
                disabled={
                  disableInput ||
                  sendMutation.isPending ||
                  updateMutation.isPending ||
                  isUploading ||
                  (!draft.trim() && attachments.length === 0) ||
                  (!activeRoomId && !editingMessage)
                }
              >
                {editingMessage ? (
                  <Pencil className="h-4 w-4" />
                ) : (
                  <Send className="h-4 w-4" />
                )}
              </Button>
            </div>
            <div className="flex justify-between items-center">
              <p className="text-xs text-muted-foreground">
                {disableInput
                  ? t("chat.send_disabled")
                  : isUploading
                  ? "Uploading..."
                  : t("chat.live_hint")}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default ChatPage;
