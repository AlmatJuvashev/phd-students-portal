"use client"

import { useRef, useEffect } from "react"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ChatMessageBubble } from "./chat-message-bubble"
import { Skeleton } from "@/components/ui/skeleton"
import { AlertCircle, MessageSquareOff } from "lucide-react"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { cn } from "@/lib/utils"
import type { ChatMessage } from "@/lib/types/chat"

interface ChatMessageListProps {
  messages: ChatMessage[]
  isLoading?: boolean
  hasError?: boolean
}

function isToday(date: Date): boolean {
  const today = new Date()
  return date.toDateString() === today.toDateString()
}

function isYesterday(date: Date): boolean {
  const yesterday = new Date()
  yesterday.setDate(yesterday.getDate() - 1)
  return date.toDateString() === yesterday.toDateString()
}

function isSameDay(date1: Date, date2: Date): boolean {
  return date1.toDateString() === date2.toDateString()
}

function formatDateLabel(date: Date): string {
  if (isToday(date)) return "Сегодня"
  if (isYesterday(date)) return "Вчера"
  return date.toLocaleDateString("ru-RU", { day: "numeric", month: "long", year: "numeric" })
}

function DateSeparator({ date }: { date: Date }) {
  return (
    <div className="flex items-center justify-center my-4">
      <span className="text-xs text-muted-foreground bg-muted/60 px-3 py-1 rounded-full">{formatDateLabel(date)}</span>
    </div>
  )
}

export function ChatMessageList({ messages, isLoading, hasError }: ChatMessageListProps) {
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [messages])

  // Loading state
  if (isLoading) {
    return (
      <div className="flex-1 p-4 space-y-4 bg-muted/10">
        {[1, 2, 3, 4, 5].map((i) => (
          <div key={i} className={cn("flex gap-2", i % 2 === 0 ? "justify-end" : "")}>
            {i % 2 !== 0 && <Skeleton className="h-8 w-8 rounded-full flex-shrink-0" />}
            <div className="space-y-2 max-w-[75%]">
              {i % 2 !== 0 && <Skeleton className="h-3 w-24" />}
              <Skeleton className={cn("h-16 rounded-2xl", i % 2 === 0 ? "w-48" : "w-64")} />
            </div>
          </div>
        ))}
      </div>
    )
  }

  // Error state
  if (hasError) {
    return (
      <div className="flex-1 flex items-center justify-center p-8 bg-muted/10">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Ошибка загрузки</AlertTitle>
          <AlertDescription>Не удалось загрузить сообщения. Попробуйте обновить страницу.</AlertDescription>
        </Alert>
      </div>
    )
  }

  // Empty state
  if (messages.length === 0) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center p-8 text-center bg-gradient-to-b from-transparent to-muted/20">
        <div className="p-4 rounded-full bg-gradient-to-br from-muted to-muted/50 mb-4">
          <MessageSquareOff className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="font-medium text-foreground mb-1">Сообщений пока нет</h3>
        <p className="text-sm text-muted-foreground max-w-xs">Будьте первым, кто напишет в этот чат!</p>
      </div>
    )
  }

  const groupedMessages = messages.reduce<
    {
      message: ChatMessage
      showSender: boolean
      showAvatar: boolean
      showDate: Date | null
    }[]
  >((acc, message, index) => {
    const prevMessage = messages[index - 1]
    const currentDate = new Date(message.createdAt)
    const prevDate = prevMessage ? new Date(prevMessage.createdAt) : null

    // Show date separator if it's the first message or day changed
    const showDate = !prevDate || !isSameDay(currentDate, prevDate) ? currentDate : null

    // Show sender name if sender changed or 5+ minutes passed
    const showSender =
      !prevMessage ||
      prevMessage.senderId !== message.senderId ||
      currentDate.getTime() - new Date(prevMessage.createdAt).getTime() > 5 * 60 * 1000

    // Show avatar only for the last message in a group from the same sender
    const nextMessage = messages[index + 1]
    const showAvatar =
      !nextMessage ||
      nextMessage.senderId !== message.senderId ||
      new Date(nextMessage.createdAt).getTime() - currentDate.getTime() > 5 * 60 * 1000

    acc.push({ message, showSender, showAvatar, showDate })
    return acc
  }, [])

  return (
    <ScrollArea className="flex-1 bg-gradient-to-b from-muted/5 to-muted/20">
      <div className="p-4 space-y-1">
        {groupedMessages.map(({ message, showSender, showAvatar, showDate }) => (
          <div key={message.id}>
            {showDate && <DateSeparator date={showDate} />}
            <div className={cn(!showSender && !message.isOwn && "ml-10")}>
              <ChatMessageBubble message={message} showSender={showSender} showAvatar={showAvatar} />
            </div>
          </div>
        ))}
        <div ref={scrollRef} />
      </div>
    </ScrollArea>
  )
}
