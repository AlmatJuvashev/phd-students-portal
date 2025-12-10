"use client"

import { cn } from "@/lib/utils"
import type { ChatMessage, UserRole } from "@/lib/types/chat"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { CheckCheck } from "lucide-react"

interface ChatMessageBubbleProps {
  message: ChatMessage
  showSender: boolean
  showAvatar?: boolean
}

const roleLabels: Record<UserRole, string> = {
  admin: "Администратор",
  advisor: "Научный руководитель",
  student: "Докторант",
}

const roleColors: Record<UserRole, { text: string; bg: string; badge: string }> = {
  admin: { text: "text-rose-600", bg: "bg-rose-100", badge: "bg-rose-100 text-rose-700" },
  advisor: { text: "text-emerald-600", bg: "bg-emerald-100", badge: "bg-emerald-100 text-emerald-700" },
  student: { text: "text-blue-600", bg: "bg-blue-100", badge: "bg-blue-100 text-blue-700" },
}

function formatTime(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })
}

export function ChatMessageBubble({ message, showSender, showAvatar = true }: ChatMessageBubbleProps) {
  const isOwn = message.isOwn
  const formattedTime = formatTime(message.createdAt)
  const colors = roleColors[message.senderRole]

  const initials = message.senderName
    .split(" ")
    .slice(0, 2)
    .map((w) => w[0])
    .join("")
    .toUpperCase()

  return (
    <div className={cn("flex gap-2 max-w-[85%] md:max-w-[75%]", isOwn ? "ml-auto flex-row-reverse" : "mr-auto")}>
      {!isOwn && (
        <div className="flex-shrink-0 w-8">
          {showAvatar && (
            <Avatar className={cn("h-8 w-8", colors.bg)}>
              <AvatarFallback className={cn("text-xs font-medium", colors.text)}>{initials}</AvatarFallback>
            </Avatar>
          )}
        </div>
      )}

      <div className={cn("flex flex-col", isOwn ? "items-end" : "items-start")}>
        {/* Sender info - only show if showSender is true */}
        {showSender && !isOwn && (
          <div className="flex items-center gap-2 mb-1 px-1">
            <span className="text-sm font-medium text-foreground">{message.senderName}</span>
            <span className={cn("text-[10px] px-1.5 py-0.5 rounded-full font-medium", colors.badge)}>
              {roleLabels[message.senderRole]}
            </span>
          </div>
        )}

        <div
          className={cn(
            "rounded-2xl px-3.5 py-2 break-words shadow-sm",
            "transition-all duration-200",
            isOwn
              ? "bg-primary text-primary-foreground rounded-br-md"
              : message.senderRole === "admin"
                ? "bg-rose-50 border border-rose-100 text-foreground rounded-bl-md"
                : message.senderRole === "advisor"
                  ? "bg-emerald-50 border border-emerald-100 text-foreground rounded-bl-md"
                  : "bg-muted/80 text-foreground rounded-bl-md",
          )}
        >
          <p className="text-[15px] leading-relaxed whitespace-pre-wrap">{message.body}</p>
          <div className={cn("flex items-center gap-1 mt-1 -mb-0.5", isOwn ? "justify-end" : "justify-start")}>
            <span className={cn("text-[10px]", isOwn ? "text-primary-foreground/70" : "text-muted-foreground")}>
              {formattedTime}
            </span>
            {/* Read receipt indicator for own messages */}
            {isOwn && <CheckCheck className="h-3.5 w-3.5 text-primary-foreground/70" />}
          </div>
        </div>
      </div>
    </div>
  )
}
