"use client"

import { useState, useMemo } from "react"
import { MessageSquare, ArrowLeft } from "lucide-react"
import { ChatSidebar } from "./components/chat-sidebar"
import { ChatRoomHeader } from "./components/chat-room-header"
import { ChatMessageList } from "./components/chat-message-list"
import { ChatMessageInput } from "./components/chat-message-input"
import { mockRooms, mockMessages } from "@/lib/data/chat-mock-data"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import type { ChatMessage } from "@/lib/types/chat"

interface ChatPageProps {
  simulateLoading?: boolean
  simulateError?: boolean
}

export function ChatPage({ simulateLoading = false, simulateError = false }: ChatPageProps) {
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null)
  const [messages, setMessages] = useState<ChatMessage[]>(mockMessages)
  const [isSending, setIsSending] = useState(false)
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(true)

  const selectedRoom = useMemo(() => mockRooms.find((room) => room.id === selectedRoomId) || null, [selectedRoomId])

  const roomMessages = useMemo(
    () => messages.filter((msg) => msg.roomId === selectedRoomId),
    [messages, selectedRoomId],
  )

  const handleSelectRoom = (roomId: string) => {
    setSelectedRoomId(roomId)
    setIsMobileSidebarOpen(false)
  }

  const handleBackToList = () => {
    setIsMobileSidebarOpen(true)
  }

  const handleSendMessage = (text: string) => {
    if (!selectedRoomId) return

    setIsSending(true)

    setTimeout(() => {
      const newMessage: ChatMessage = {
        id: `msg-${Date.now()}`,
        roomId: selectedRoomId,
        senderId: "current",
        senderName: "Вы",
        senderRole: "student",
        body: text,
        createdAt: new Date().toISOString(),
        isOwn: true,
      }
      setMessages((prev) => [...prev, newMessage])
      setIsSending(false)
    }, 500)
  }

  return (
    <div className="flex h-[100dvh] bg-background overflow-hidden">
      <div
        className={cn(
          "absolute inset-0 z-30 bg-background transition-transform duration-300 ease-out",
          "md:relative md:inset-auto md:z-auto md:w-80 md:flex-shrink-0 md:translate-x-0",
          isMobileSidebarOpen ? "translate-x-0" : "-translate-x-full",
        )}
      >
        <ChatSidebar
          rooms={mockRooms}
          selectedRoomId={selectedRoomId}
          onSelectRoom={handleSelectRoom}
          isLoading={simulateLoading}
          messages={messages}
        />
      </div>

      {/* Main content area */}
      <div
        className={cn(
          "flex-1 flex flex-col min-w-0 transition-opacity duration-200",
          isMobileSidebarOpen ? "opacity-0 md:opacity-100" : "opacity-100",
        )}
      >
        {selectedRoom ? (
          <>
            <ChatRoomHeader room={selectedRoom} connectionStatus="connected" onBack={handleBackToList} />
            <ChatMessageList messages={roomMessages} isLoading={simulateLoading} hasError={simulateError} />
            <ChatMessageInput onSend={handleSendMessage} isSending={isSending} isArchived={selectedRoom.isArchived} />
          </>
        ) : (
          /* No room selected state - improved styling */
          <div className="flex-1 flex flex-col items-center justify-center p-8 text-center bg-muted/20">
            <Button variant="ghost" size="sm" onClick={handleBackToList} className="md:hidden absolute top-4 left-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Чаты
            </Button>
            <div className="p-6 rounded-full bg-gradient-to-br from-primary/10 to-primary/5 mb-6">
              <MessageSquare className="h-12 w-12 text-primary/60" />
            </div>
            <h2 className="text-xl font-semibold text-foreground mb-2">Выберите чат</h2>
            <p className="text-muted-foreground max-w-md text-sm leading-relaxed">
              Слева отображаются доступные вам группы. Выберите одну, чтобы просматривать сообщения и задавать вопросы.
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
