"use client"

import { cn } from "@/lib/utils"
import type { ChatRoom, ChatRoomType, ChatMessage } from "@/lib/types/chat"
import { Users, GraduationCap, MessageSquare } from "lucide-react"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"

interface ChatRoomListItemProps {
  room: ChatRoom
  isSelected: boolean
  onClick: () => void
  lastMessage?: ChatMessage
}

const typeConfig: Record<ChatRoomType, { label: string; icon: typeof Users; color: string; bgColor: string }> = {
  cohort: { label: "Когорта", icon: Users, color: "text-blue-600", bgColor: "bg-blue-100" },
  advisory: { label: "Научрук", icon: GraduationCap, color: "text-emerald-600", bgColor: "bg-emerald-100" },
  other: { label: "Другое", icon: MessageSquare, color: "text-slate-600", bgColor: "bg-slate-100" },
}

export function ChatRoomListItem({ room, isSelected, onClick, lastMessage }: ChatRoomListItemProps) {
  const config = typeConfig[room.type]
  const Icon = config.icon

  const initials = room.name
    .split(" ")
    .slice(0, 2)
    .map((w) => w[0])
    .join("")
    .toUpperCase()

  const formatTime = (dateStr: string) => {
    const date = new Date(dateStr)
    const now = new Date()
    const diffDays = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24))

    if (diffDays === 0) {
      return date.toLocaleTimeString("ru", { hour: "2-digit", minute: "2-digit" })
    } else if (diffDays === 1) {
      return "Вчера"
    } else if (diffDays < 7) {
      return date.toLocaleDateString("ru", { weekday: "short" })
    }
    return date.toLocaleDateString("ru", { day: "numeric", month: "short" })
  }

  return (
    <button
      onClick={onClick}
      className={cn(
        "w-full text-left px-3 py-3 rounded-xl transition-all duration-200",
        "hover:bg-muted/60 focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "active:scale-[0.98]",
        isSelected && "bg-primary/10 hover:bg-primary/15",
        room.isArchived && "opacity-60",
      )}
    >
      <div className="flex items-center gap-3">
        <div className="relative flex-shrink-0">
          <Avatar className={cn("h-11 w-11", config.bgColor)}>
            <AvatarFallback className={cn("text-sm font-medium", config.color)}>{initials}</AvatarFallback>
          </Avatar>
          {/* Type icon badge */}
          <div
            className={cn(
              "absolute -bottom-0.5 -right-0.5 p-1 rounded-full bg-background border-2 border-background",
              config.bgColor,
            )}
          >
            <Icon className={cn("h-2.5 w-2.5", config.color)} />
          </div>
        </div>

        {/* Room info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <span
              className={cn(
                "font-medium text-sm truncate",
                isSelected ? "text-primary" : "text-foreground",
                room.unreadCount && room.unreadCount > 0 && "font-semibold",
              )}
            >
              {room.name}
            </span>
            {room.isArchived ? (
              <span className="text-[10px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded-full flex-shrink-0">
                Архив
              </span>
            ) : lastMessage ? (
              <span className="text-xs text-muted-foreground flex-shrink-0">{formatTime(lastMessage.createdAt)}</span>
            ) : null}
          </div>
          <div className="flex items-center justify-between gap-2 mt-0.5">
            <p
              className={cn(
                "text-sm truncate",
                room.unreadCount && room.unreadCount > 0 ? "text-foreground/80" : "text-muted-foreground",
              )}
            >
              {lastMessage ? (
                <>
                  {lastMessage.isOwn && <span className="text-muted-foreground">Вы: </span>}
                  {lastMessage.body}
                </>
              ) : (
                <span className="text-muted-foreground italic">Нет сообщений</span>
              )}
            </p>
            {/* Unread badge */}
            {room.unreadCount && room.unreadCount > 0 && (
              <span className="flex-shrink-0 min-w-5 h-5 flex items-center justify-center bg-primary text-primary-foreground text-xs font-semibold rounded-full px-1.5">
                {room.unreadCount > 99 ? "99+" : room.unreadCount}
              </span>
            )}
          </div>
        </div>
      </div>
    </button>
  )
}
