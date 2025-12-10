"use client"

import { format } from "date-fns"
import { ru } from "date-fns/locale"
import { MoreHorizontal, Users, Archive, ArchiveRestore, Pencil, GraduationCap, MessageSquare } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Card, CardContent } from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn } from "@/lib/utils"
import type { ChatRoom, ChatRoomType } from "@/lib/types/chat"

interface ChatRoomTableProps {
  rooms: ChatRoom[]
  onEdit: (roomId: string) => void
  onManageMembers: (roomId: string) => void
  onToggleArchive: (roomId: string) => void
}

const typeConfig: Record<ChatRoomType, { label: string; icon: typeof Users; color: string; bgColor: string }> = {
  cohort: { label: "Когорта", icon: Users, color: "text-blue-600", bgColor: "bg-blue-100" },
  advisory: { label: "Научрук", icon: GraduationCap, color: "text-emerald-600", bgColor: "bg-emerald-100" },
  other: { label: "Другое", icon: MessageSquare, color: "text-slate-600", bgColor: "bg-slate-100" },
}

export function ChatRoomTable({ rooms, onEdit, onManageMembers, onToggleArchive }: ChatRoomTableProps) {
  return (
    <div className="space-y-3">
      {rooms.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <div className="p-4 rounded-full bg-muted mb-4">
              <MessageSquare className="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 className="font-medium text-foreground mb-1">Групп пока нет</h3>
            <p className="text-sm text-muted-foreground">Создайте первую группу чата</p>
          </CardContent>
        </Card>
      ) : (
        rooms.map((room) => {
          const config = typeConfig[room.type]
          const Icon = config.icon
          const initials = room.name
            .split(" ")
            .slice(0, 2)
            .map((w) => w[0])
            .join("")
            .toUpperCase()

          return (
            <Card key={room.id} className={cn("transition-all hover:shadow-md", room.isArchived && "opacity-70")}>
              <CardContent className="p-4">
                <div className="flex items-start gap-3">
                  {/* Avatar */}
                  <Avatar className={cn("h-12 w-12 flex-shrink-0", config.bgColor)}>
                    <AvatarFallback className={cn("text-sm font-medium", config.color)}>{initials}</AvatarFallback>
                  </Avatar>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div className="min-w-0">
                        <h3 className="font-medium text-foreground truncate">{room.name}</h3>
                        {room.description && (
                          <p className="text-sm text-muted-foreground truncate mt-0.5">{room.description}</p>
                        )}
                      </div>
                      {/* Actions */}
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8 flex-shrink-0 -mr-2">
                            <MoreHorizontal className="h-4 w-4" />
                            <span className="sr-only">Действия</span>
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-48">
                          <DropdownMenuItem onClick={() => onEdit(room.id)}>
                            <Pencil className="h-4 w-4 mr-2" />
                            Редактировать
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => onManageMembers(room.id)}>
                            <Users className="h-4 w-4 mr-2" />
                            Участники
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem onClick={() => onToggleArchive(room.id)}>
                            {room.isArchived ? (
                              <>
                                <ArchiveRestore className="h-4 w-4 mr-2" />
                                Разархивировать
                              </>
                            ) : (
                              <>
                                <Archive className="h-4 w-4 mr-2" />
                                Архивировать
                              </>
                            )}
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>

                    {/* Meta info */}
                    <div className="flex flex-wrap items-center gap-2 mt-3">
                      <Badge variant="secondary" className={cn("text-xs", config.bgColor, config.color)}>
                        <Icon className="h-3 w-3 mr-1" />
                        {config.label}
                      </Badge>
                      <span className="flex items-center gap-1 text-xs text-muted-foreground">
                        <Users className="h-3 w-3" />
                        {room.membersCount || 0}
                      </span>
                      {room.isArchived ? (
                        <Badge variant="outline" className="text-xs text-muted-foreground">
                          Архив
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="text-xs text-emerald-600 border-emerald-200 bg-emerald-50">
                          Активен
                        </Badge>
                      )}
                      {room.createdAt && (
                        <span className="text-xs text-muted-foreground ml-auto">
                          {format(new Date(room.createdAt), "d MMM yyyy", { locale: ru })}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )
        })
      )}
    </div>
  )
}
