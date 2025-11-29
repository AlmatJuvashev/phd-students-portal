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
} from "lucide-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import {
  ChatMessage,
  ChatRoom,
  listMessages,
  listRooms,
  sendMessage,
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
  return date.toLocaleDateString("ru-RU", { day: "numeric", month: "long", year: "numeric" });
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
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [search, setSearch] = useState("");

  const {
    data: rooms = [],
    isLoading: roomsLoading,
    isFetching: roomsFetching,
  } = useQuery<ChatRoom[]>({
    queryKey: ["chat", "rooms"],
    queryFn: listRooms,
  });

  useEffect(() => {
    if (!rooms.length) return;
    const exists = rooms.find((r) => r.id === activeRoomId);
    if (!activeRoomId || !exists) {
      setActiveRoomId(rooms[0].id);
      setSidebarOpen(false);
    }
  }, [rooms, activeRoomId]);

  const activeRoom = useMemo(
    () => rooms.find((room) => room.id === activeRoomId),
    [rooms, activeRoomId]
  );

  const {
    data: messages = [],
    isLoading: messagesLoading,
    isFetching: messagesFetching,
  } = useQuery<ChatMessage[]>({
    queryKey: ["chat", "messages", activeRoomId],
    queryFn: () => listMessages({ roomId: activeRoomId, limit: 50 }),
    enabled: !!activeRoomId,
  });

  const sortedMessages = useMemo(
    () =>
      [...messages].sort(
        (a, b) =>
          new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
      ),
    [messages]
  );

  const sendMutation = useMutation({
    mutationFn: (body: string) => sendMessage(activeRoomId, body),
    onSuccess: () => {
      setDraft("");
      qc.invalidateQueries({ queryKey: ["chat", "messages", activeRoomId] });
    },
  });

  const handleSend = () => {
    if (!draft.trim() || !activeRoomId) return;
    sendMutation.mutate(draft.trim());
  };

  const roomsCount = rooms.length;
  const disableInput = !activeRoom || activeRoom.is_archived;
  const filteredRooms = rooms.filter((r) =>
    r.name.toLowerCase().includes(search.toLowerCase())
  );
  const activeRooms = filteredRooms.filter((r) => !r.is_archived);
  const archivedRooms = filteredRooms.filter((r) => r.is_archived);

  return (
    <div className="space-y-4">
      <div className="space-y-1">
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
        <div className="inline-flex items-center gap-2 rounded-md bg-muted px-3 py-2 text-xs text-muted-foreground">
          <ShieldCheck className="h-4 w-4 text-primary" />
          <span>{t("chat.demo_notice")}</span>
        </div>
      </div>

      <div className="flex min-h-[60vh] rounded-xl border shadow-sm overflow-hidden bg-background">
        {/* Sidebar */}
        <div
          className={cn(
            "w-full max-w-sm border-r bg-muted/30 transition-transform duration-300 md:translate-x-0",
            sidebarOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"
          )}
        >
          <div className="p-4 border-b bg-gradient-to-b from-primary/5 to-transparent space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="p-2 rounded-lg bg-primary/10">
                  <MessageCircle className="h-5 w-5 text-primary" />
                </div>
                <div>
                  <div className="font-semibold">{t("chat.rooms_heading")}</div>
                  <div className="text-xs text-muted-foreground">
                    {roomsFetching ? "…" : roomsCount} {t("chat.room_members", { count: roomsCount })}
                  </div>
                </div>
              </div>
            </div>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder={t("chat.search_placeholder", { defaultValue: "Поиск чатов..." })}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>

          <div className="p-3 space-y-4 overflow-y-auto max-h-[70vh]">
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
                            setSidebarOpen(false);
                          }}
                          className={cn(
                            "w-full rounded-lg border px-3 py-2 text-left hover:border-primary/60 hover:bg-primary/5 transition",
                            room.id === activeRoomId && "border-primary bg-primary/5 shadow-sm"
                          )}
                        >
                          <div className="flex items-center justify-between gap-2">
                            <div>
                              <div className="font-semibold leading-tight">{room.name}</div>
                              <div className="flex flex-wrap items-center gap-2 text-[11px] text-muted-foreground mt-1">
                                <Badge variant="outline" className="text-[10px]">
                                  {t(`chat.types.${room.type}`)}
                                </Badge>
                                {room.created_by_role && (
                                  <Badge variant="outline" className="text-[10px]">
                                    {room.created_by_role}
                                  </Badge>
                                )}
                              </div>
                            </div>
                            <span className="text-[11px] text-muted-foreground flex items-center gap-1">
                              <Clock4 className="h-3 w-3" />
                              {formatTime(room.created_at)}
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
                            setSidebarOpen(false);
                          }}
                          className={cn(
                            "w-full rounded-lg border px-3 py-2 text-left hover:border-primary/40 hover:bg-primary/5 transition",
                            room.id === activeRoomId && "border-primary bg-primary/5 shadow-sm"
                          )}
                        >
                          <div className="flex items-center justify-between gap-2">
                            <div>
                              <div className="font-semibold leading-tight">{room.name}</div>
                              <div className="flex flex-wrap items-center gap-2 text-[11px] text-muted-foreground mt-1">
                                <Badge variant="outline" className="text-[10px]">
                                  {t(`chat.types.${room.type}`)}
                                </Badge>
                                <Badge variant="outline" className="text-[10px]">
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
        <div className="flex-1 flex flex-col min-w-0">
          <div className="flex items-center justify-between border-b px-4 py-3">
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="icon"
                className="md:hidden"
                onClick={() => setSidebarOpen((v) => !v)}
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <div>
                <div className="flex items-center gap-2">
                  <Sparkles className="h-4 w-4 text-primary" />
                  <span className="font-semibold text-sm">
                    {activeRoom ? activeRoom.name : t("chat.messages_heading")}
                  </span>
                </div>
                <p className="text-xs text-muted-foreground">
                  {activeRoom
                    ? t(`chat.types.${activeRoom.type}`)
                    : t("chat.empty_state")}
                </p>
              </div>
            </div>
            {activeRoom?.is_archived && (
              <Badge variant="outline" className="text-[10px]">
                {t("chat.archived")}
              </Badge>
            )}
          </div>

          <div className="flex-1 flex flex-col bg-muted/10">
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
                    new Date(prev.created_at).toDateString() !== thisDate.toDateString();
                  const showSender =
                    !prev ||
                    prev.sender_id !== msg.sender_id ||
                    thisDate.getTime() - new Date(prev.created_at).getTime() > 5 * 60 * 1000;
                  return (
                    <div key={msg.id} className="space-y-1">
                      {showDate && (
                        <div className="flex items-center justify-center my-2">
                          <span className="text-[11px] text-muted-foreground bg-muted px-3 py-1 rounded-full">
                            {formatDateLabel(thisDate)}
                          </span>
                        </div>
                      )}
                      <div className="flex items-start gap-2">
                        <div className="h-8 w-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-semibold">
                          {initials(msg.sender_name || msg.sender_id)}
                        </div>
                        <div className="flex-1">
                          {showSender && (
                            <div className="flex items-center gap-2 text-[11px] text-muted-foreground">
                              <span className="font-medium text-foreground">
                                {msg.sender_name || msg.sender_id}
                              </span>
                              {msg.sender_role && (
                                <Badge variant="outline" className="text-[10px]">
                                  {msg.sender_role}
                                </Badge>
                              )}
                              <span>{formatTime(msg.created_at)}</span>
                            </div>
                          )}
                          <div className="mt-1 rounded-2xl bg-background border px-3 py-2 shadow-sm">
                            <p className="text-sm whitespace-pre-line leading-relaxed">
                              {msg.body}
                            </p>
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          <div className="border-t p-3 bg-background">
            <Textarea
              placeholder={t("chat.input_placeholder")}
              className="resize-none"
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              disabled={disableInput || sendMutation.isPending || !activeRoomId}
            />
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between mt-2">
              <p className="text-xs text-muted-foreground">
                {disableInput ? t("chat.send_disabled") : t("chat.live_hint")}
              </p>
              <Button
                size="sm"
                className="sm:w-auto"
                onClick={handleSend}
                disabled={
                  disableInput ||
                  sendMutation.isPending ||
                  !draft.trim() ||
                  !activeRoomId
                }
              >
                <Send className="mr-2 h-4 w-4" />
                {sendMutation.isPending ? t("common.loading") : t("chat.send")}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default ChatPage;
