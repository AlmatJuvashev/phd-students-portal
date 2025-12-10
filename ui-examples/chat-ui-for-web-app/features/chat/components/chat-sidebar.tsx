"use client"

import { useState } from "react"
import { Search, MessageCircle } from "lucide-react"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ChatRoomListItem } from "./chat-room-list-item"
import { Skeleton } from "@/components/ui/skeleton"
import type { ChatRoom, ChatMessage } from "@/lib/types/chat"

interface ChatSidebarProps {
  rooms: ChatRoom[]
  selectedRoomId: string | null
  onSelectRoom: (roomId: string) => void
  isLoading?: boolean
  messages?: ChatMessage[]
}

export function ChatSidebar({ rooms, selectedRoomId, onSelectRoom, isLoading, messages = [] }: ChatSidebarProps) {
  const [searchQuery, setSearchQuery] = useState("")

  const filteredRooms = rooms.filter((room) => room.name.toLowerCase().includes(searchQuery.toLowerCase()))

  const activeRooms = filteredRooms.filter((room) => !room.isArchived)
  const archivedRooms = filteredRooms.filter((room) => room.isArchived)

  const getLastMessage = (roomId: string) => {
    const roomMessages = messages.filter((m) => m.roomId === roomId)
    return roomMessages[roomMessages.length - 1]
  }

  return (
    <div className="flex flex-col h-full bg-card border-r">
      {/* Header with gradient accent */}
      <div className="p-4 border-b bg-gradient-to-b from-primary/5 to-transparent">
        <div className="flex items-center gap-3 mb-4">
          <div className="p-2 rounded-xl bg-primary/10">
            <MessageCircle className="h-5 w-5 text-primary" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground">Сообщения</h2>
            <p className="text-xs text-muted-foreground">{activeRooms.length} активных чатов</p>
          </div>
        </div>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Поиск чатов..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 bg-background/50 border-muted focus-visible:bg-background transition-colors"
          />
        </div>
      </div>

      {/* Room list */}
      <ScrollArea className="flex-1">
        <div className="p-2">
          {isLoading ? (
            <div className="space-y-2 p-2">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="flex items-center gap-3 p-3">
                  <Skeleton className="h-10 w-10 rounded-full" />
                  <div className="flex-1 space-y-2">
                    <Skeleton className="h-4 w-3/4" />
                    <Skeleton className="h-3 w-1/2" />
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <>
              {/* Active rooms section */}
              <div className="mb-2">
                <h3 className="text-xs font-medium text-muted-foreground uppercase tracking-wider px-3 py-2">
                  Мои группы
                </h3>
                <div className="space-y-0.5">
                  {activeRooms.length > 0 ? (
                    activeRooms.map((room) => (
                      <ChatRoomListItem
                        key={room.id}
                        room={room}
                        isSelected={selectedRoomId === room.id}
                        onClick={() => onSelectRoom(room.id)}
                        lastMessage={getLastMessage(room.id)}
                      />
                    ))
                  ) : (
                    <p className="text-sm text-muted-foreground px-3 py-4 text-center">
                      {searchQuery ? "Чаты не найдены" : "Нет активных чатов"}
                    </p>
                  )}
                </div>
              </div>

              {/* Archived rooms section */}
              {archivedRooms.length > 0 && (
                <div className="mt-4 pt-4 border-t">
                  <h3 className="text-xs font-medium text-muted-foreground uppercase tracking-wider px-3 py-2">
                    Архив
                  </h3>
                  <div className="space-y-0.5">
                    {archivedRooms.map((room) => (
                      <ChatRoomListItem
                        key={room.id}
                        room={room}
                        isSelected={selectedRoomId === room.id}
                        onClick={() => onSelectRoom(room.id)}
                        lastMessage={getLastMessage(room.id)}
                      />
                    ))}
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </ScrollArea>
    </div>
  )
}
