"use client"

import { cn } from "@/lib/utils"
import type { ChatRoom, ChatRoomType } from "@/lib/types/chat"
import { Users, GraduationCap, MessageSquare, MoreVertical, Archive, ArrowLeft, Phone, Video } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { ChatConnectionStatus } from "./chat-connection-status"

interface ChatRoomHeaderProps {
  room: ChatRoom
  connectionStatus?: "connected" | "connecting" | "disconnected"
  onBack?: () => void
}

const typeConfig: Record<ChatRoomType, { label: string; icon: typeof Users; color: string; bgColor: string }> = {
  cohort: { label: "Когорта", icon: Users, color: "text-blue-600", bgColor: "bg-blue-100" },
  advisory: { label: "Научное руководство", icon: GraduationCap, color: "text-emerald-600", bgColor: "bg-emerald-100" },
  other: { label: "Общий чат", icon: MessageSquare, color: "text-slate-600", bgColor: "bg-slate-100" },
}

export function ChatRoomHeader({ room, connectionStatus = "connected", onBack }: ChatRoomHeaderProps) {
  const config = typeConfig[room.type]
  const Icon = config.icon

  const initials = room.name
    .split(" ")
    .slice(0, 2)
    .map((w) => w[0])
    .join("")
    .toUpperCase()

  return (
    <div className="flex items-center justify-between px-2 md:px-4 py-3 border-b bg-card/80 backdrop-blur-sm sticky top-0 z-10">
      <div className="flex items-center gap-2 md:gap-3 min-w-0">
        <Button variant="ghost" size="icon" onClick={onBack} className="md:hidden h-9 w-9 -ml-1">
          <ArrowLeft className="h-5 w-5" />
          <span className="sr-only">Назад</span>
        </Button>

        <Avatar className={cn("h-10 w-10", config.bgColor)}>
          <AvatarFallback className={cn("text-sm font-medium", config.color)}>{initials}</AvatarFallback>
        </Avatar>

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <h1 className="font-semibold text-foreground truncate text-sm md:text-base">{room.name}</h1>
            {room.isArchived && (
              <span className="hidden sm:inline-flex items-center gap-1 text-xs text-muted-foreground bg-muted px-2 py-0.5 rounded-full">
                <Archive className="h-3 w-3" />
                Архив
              </span>
            )}
          </div>
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <Icon className="h-3 w-3" />
            <span className="hidden sm:inline">{config.label}</span>
            <span className="hidden sm:inline">·</span>
            <span>{room.membersCount} участников</span>
            <ChatConnectionStatus status={connectionStatus} className="sm:hidden" />
          </div>
        </div>
      </div>

      <div className="flex items-center gap-1">
        <ChatConnectionStatus status={connectionStatus} className="hidden sm:flex mr-2" />

        <Button
          variant="ghost"
          size="icon"
          className="hidden md:flex h-9 w-9 text-muted-foreground hover:text-foreground"
        >
          <Phone className="h-4 w-4" />
          <span className="sr-only">Позвонить</span>
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="hidden md:flex h-9 w-9 text-muted-foreground hover:text-foreground"
        >
          <Video className="h-4 w-4" />
          <span className="sr-only">Видеозвонок</span>
        </Button>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-9 w-9">
              <MoreVertical className="h-4 w-4" />
              <span className="sr-only">Меню</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <DropdownMenuItem>Информация о группе</DropdownMenuItem>
            <DropdownMenuItem>Участники ({room.membersCount})</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>Настройки уведомлений</DropdownMenuItem>
            <DropdownMenuItem>Поиск в чате</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}
