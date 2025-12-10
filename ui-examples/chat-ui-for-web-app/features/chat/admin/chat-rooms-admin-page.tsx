"use client"

import { useState } from "react"
import { Plus, MessageCircle, ArrowLeft } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ChatRoomTable } from "./chat-room-table"
import { ChatRoomFormModal } from "./chat-room-form-modal"
import { ChatRoomMembersModal } from "./chat-room-members-modal"
import { mockRooms } from "@/lib/data/chat-mock-data"
import Link from "next/link"
import type { ChatRoom, RoomFormValues } from "@/lib/types/chat"

export function ChatRoomsAdminPage() {
  const [rooms, setRooms] = useState<ChatRoom[]>(mockRooms)
  const [isFormModalOpen, setIsFormModalOpen] = useState(false)
  const [isMembersModalOpen, setIsMembersModalOpen] = useState(false)
  const [selectedRoom, setSelectedRoom] = useState<ChatRoom | null>(null)

  const handleCreateRoom = () => {
    setSelectedRoom(null)
    setIsFormModalOpen(true)
  }

  const handleEditRoom = (roomId: string) => {
    const room = rooms.find((r) => r.id === roomId)
    if (room) {
      setSelectedRoom(room)
      setIsFormModalOpen(true)
    }
  }

  const handleManageMembers = (roomId: string) => {
    const room = rooms.find((r) => r.id === roomId)
    if (room) {
      setSelectedRoom(room)
      setIsMembersModalOpen(true)
    }
  }

  const handleToggleArchive = (roomId: string) => {
    setRooms((prev) => prev.map((room) => (room.id === roomId ? { ...room, isArchived: !room.isArchived } : room)))
  }

  const handleFormSubmit = (data: RoomFormValues) => {
    if (selectedRoom) {
      setRooms((prev) => prev.map((room) => (room.id === selectedRoom.id ? { ...room, ...data } : room)))
    } else {
      const newRoom: ChatRoom = {
        id: `room-${Date.now()}`,
        ...data,
        membersCount: 0,
        createdAt: new Date().toISOString(),
      }
      setRooms((prev) => [...prev, newRoom])
    }
  }

  return (
    <div className="min-h-[100dvh] bg-muted/30">
      <div className="sticky top-0 z-10 border-b bg-card/80 backdrop-blur-sm">
        <div className="container mx-auto px-4 md:px-6 py-3 md:py-4">
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-2 md:gap-3 min-w-0">
              <Button variant="ghost" size="icon" asChild className="h-9 w-9 -ml-2">
                <Link href="/">
                  <ArrowLeft className="h-5 w-5" />
                  <span className="sr-only">Назад</span>
                </Link>
              </Button>
              <div className="p-2 rounded-xl bg-primary/10 hidden sm:flex">
                <MessageCircle className="h-5 w-5 text-primary" />
              </div>
              <div className="min-w-0">
                <h1 className="text-lg md:text-xl font-semibold text-foreground truncate">Группы чата</h1>
                <p className="text-xs md:text-sm text-muted-foreground hidden sm:block">
                  Управление группами и участниками
                </p>
              </div>
            </div>
            <Button onClick={handleCreateRoom} size="sm" className="shrink-0">
              <Plus className="h-4 w-4 md:mr-2" />
              <span className="hidden md:inline">Создать группу</span>
            </Button>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 md:px-6 py-4 md:py-6">
        <ChatRoomTable
          rooms={rooms}
          onEdit={handleEditRoom}
          onManageMembers={handleManageMembers}
          onToggleArchive={handleToggleArchive}
        />
      </div>

      <ChatRoomFormModal
        open={isFormModalOpen}
        onOpenChange={setIsFormModalOpen}
        room={selectedRoom}
        onSubmit={handleFormSubmit}
      />

      <ChatRoomMembersModal open={isMembersModalOpen} onOpenChange={setIsMembersModalOpen} room={selectedRoom} />
    </div>
  )
}
