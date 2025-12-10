"use client"

import { useState, useMemo } from "react"
import { Search, X, UserPlus, UserMinus } from "lucide-react"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { cn } from "@/lib/utils"
import type { ChatMember, ChatRoom, UserRole } from "@/lib/types/chat"
import { mockMembers, mockAllUsers } from "@/lib/data/chat-mock-data"

interface ChatRoomMembersModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  room: ChatRoom | null
}

const roleLabels: Record<UserRole, string> = {
  admin: "Админ",
  advisor: "Научрук",
  student: "Студент",
}

const roleColors: Record<UserRole, { text: string; bg: string }> = {
  admin: { text: "text-rose-600", bg: "bg-rose-100" },
  advisor: { text: "text-emerald-600", bg: "bg-emerald-100" },
  student: { text: "text-blue-600", bg: "bg-blue-100" },
}

export function ChatRoomMembersModal({ open, onOpenChange, room }: ChatRoomMembersModalProps) {
  const [members, setMembers] = useState<ChatMember[]>(mockMembers)
  const [searchQuery, setSearchQuery] = useState("")
  const [addUserQuery, setAddUserQuery] = useState("")

  const filteredMembers = useMemo(() => {
    return members.filter(
      (member) =>
        member.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        member.email.toLowerCase().includes(searchQuery.toLowerCase()),
    )
  }, [members, searchQuery])

  const availableUsers = useMemo(() => {
    const memberIds = new Set(members.map((m) => m.id))
    return mockAllUsers.filter(
      (user) =>
        !memberIds.has(user.id) &&
        (user.name.toLowerCase().includes(addUserQuery.toLowerCase()) ||
          user.email.toLowerCase().includes(addUserQuery.toLowerCase())),
    )
  }, [members, addUserQuery])

  const handleRemoveMember = (memberId: string) => {
    setMembers((prev) => prev.filter((m) => m.id !== memberId))
  }

  const handleAddMember = (user: ChatMember) => {
    setMembers((prev) => [...prev, user])
    setAddUserQuery("")
  }

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[550px] max-h-[85dvh] flex flex-col p-0 gap-0">
        <DialogHeader className="p-4 pb-2 sm:p-6 sm:pb-4">
          <DialogTitle>Управление участниками</DialogTitle>
          <DialogDescription className="truncate">{room?.name}</DialogDescription>
        </DialogHeader>

        <Tabs defaultValue="members" className="flex-1 flex flex-col min-h-0">
          <TabsList className="mx-4 sm:mx-6 grid grid-cols-2">
            <TabsTrigger value="members" className="gap-2">
              <UserMinus className="h-4 w-4" />
              <span className="hidden sm:inline">Участники</span>
              <span className="sm:hidden">Список</span>
              <Badge variant="secondary" className="ml-1 h-5 px-1.5 text-xs">
                {members.length}
              </Badge>
            </TabsTrigger>
            <TabsTrigger value="add" className="gap-2">
              <UserPlus className="h-4 w-4" />
              Добавить
            </TabsTrigger>
          </TabsList>

          {/* Members Tab */}
          <TabsContent value="members" className="flex-1 flex flex-col min-h-0 mt-0 p-4 sm:p-6 pt-4">
            <div className="relative mb-3">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Поиск участников..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9"
              />
            </div>
            <ScrollArea className="flex-1 -mx-4 sm:-mx-6 px-4 sm:px-6">
              <div className="space-y-1">
                {filteredMembers.length > 0 ? (
                  filteredMembers.map((member) => {
                    const colors = roleColors[member.role]
                    return (
                      <div
                        key={member.id}
                        className="flex items-center gap-3 p-2.5 rounded-xl hover:bg-muted/60 transition-colors group"
                      >
                        <Avatar className={cn("h-10 w-10", colors.bg)}>
                          <AvatarFallback className={cn("text-sm", colors.text)}>
                            {getInitials(member.name)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="flex-1 min-w-0">
                          <div className="text-sm font-medium truncate">{member.name}</div>
                          <div className="text-xs text-muted-foreground truncate">{member.email}</div>
                        </div>
                        <Badge variant="secondary" className={cn("text-xs shrink-0", colors.bg, colors.text)}>
                          {roleLabels[member.role]}
                        </Badge>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 opacity-0 group-hover:opacity-100 text-muted-foreground hover:text-destructive transition-opacity"
                          onClick={() => handleRemoveMember(member.id)}
                        >
                          <X className="h-4 w-4" />
                          <span className="sr-only">Удалить</span>
                        </Button>
                      </div>
                    )
                  })
                ) : (
                  <p className="text-sm text-muted-foreground text-center py-8">
                    {searchQuery ? "Участники не найдены" : "Нет участников"}
                  </p>
                )}
              </div>
            </ScrollArea>
          </TabsContent>

          {/* Add Member Tab */}
          <TabsContent value="add" className="flex-1 flex flex-col min-h-0 mt-0 p-4 sm:p-6 pt-4">
            <div className="relative mb-3">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Найти пользователя..."
                value={addUserQuery}
                onChange={(e) => setAddUserQuery(e.target.value)}
                className="pl-9"
              />
            </div>
            <ScrollArea className="flex-1 -mx-4 sm:-mx-6 px-4 sm:px-6">
              <div className="space-y-1">
                {availableUsers.length > 0 ? (
                  availableUsers.map((user) => {
                    const colors = roleColors[user.role]
                    return (
                      <button
                        key={user.id}
                        onClick={() => handleAddMember(user)}
                        className="w-full flex items-center gap-3 p-2.5 rounded-xl hover:bg-muted/60 transition-colors text-left group"
                      >
                        <Avatar className={cn("h-10 w-10", colors.bg)}>
                          <AvatarFallback className={cn("text-sm", colors.text)}>
                            {getInitials(user.name)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="flex-1 min-w-0">
                          <div className="text-sm font-medium truncate">{user.name}</div>
                          <div className="text-xs text-muted-foreground truncate">{user.email}</div>
                        </div>
                        <Badge variant="secondary" className={cn("text-xs shrink-0", colors.bg, colors.text)}>
                          {roleLabels[user.role]}
                        </Badge>
                        <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                          <UserPlus className="h-4 w-4 text-primary" />
                        </div>
                      </button>
                    )
                  })
                ) : (
                  <p className="text-sm text-muted-foreground text-center py-8">
                    {addUserQuery ? "Пользователи не найдены" : "Введите имя для поиска"}
                  </p>
                )}
              </div>
            </ScrollArea>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
