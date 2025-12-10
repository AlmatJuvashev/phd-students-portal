"use client"

import { useState, useRef, useEffect, type KeyboardEvent } from "react"
import { Send, Lock, Paperclip, Smile, Mic } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Textarea } from "@/components/ui/textarea"
import { cn } from "@/lib/utils"

interface ChatMessageInputProps {
  onSend: (text: string) => void
  isSending?: boolean
  isArchived?: boolean
  placeholder?: string
}

export function ChatMessageInput({
  onSend,
  isSending = false,
  isArchived = false,
  placeholder = "Введите сообщение...",
}: ChatMessageInputProps) {
  const [text, setText] = useState("")
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // Auto-resize textarea
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto"
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 120)}px`
    }
  }, [text])

  const handleSend = () => {
    if (text.trim() && !isSending && !isArchived) {
      onSend(text.trim())
      setText("")
    }
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  // Archived state
  if (isArchived) {
    return (
      <div className="px-4 py-4 border-t bg-muted/30">
        <div className="flex items-center justify-center gap-2 text-muted-foreground py-2 bg-muted/50 rounded-xl">
          <Lock className="h-4 w-4" />
          <span className="text-sm">Этот чат заархивирован</span>
        </div>
      </div>
    )
  }

  return (
    <div className="p-3 md:p-4 border-t bg-card/80 backdrop-blur-sm">
      <div className="flex items-end gap-2 bg-muted/50 rounded-2xl p-2 border border-border/50 focus-within:border-primary/30 focus-within:bg-background transition-all">
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9 flex-shrink-0 text-muted-foreground hover:text-foreground rounded-xl"
        >
          <Paperclip className="h-5 w-5" />
          <span className="sr-only">Прикрепить файл</span>
        </Button>

        {/* Text input */}
        <div className="flex-1">
          <Textarea
            ref={textareaRef}
            value={text}
            onChange={(e) => setText(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            disabled={isSending}
            rows={1}
            className={cn(
              "min-h-9 max-h-[120px] resize-none border-0 bg-transparent p-1",
              "focus-visible:ring-0 focus-visible:ring-offset-0",
              "placeholder:text-muted-foreground/60",
            )}
          />
        </div>

        <div className="flex items-center gap-1 flex-shrink-0">
          {!text.trim() && (
            <>
              {/* Emoji button (visual only) */}
              <Button
                variant="ghost"
                size="icon"
                className="h-9 w-9 text-muted-foreground hover:text-foreground rounded-xl hidden sm:flex"
              >
                <Smile className="h-5 w-5" />
                <span className="sr-only">Эмодзи</span>
              </Button>
              {/* Voice message button (visual only) */}
              <Button
                variant="ghost"
                size="icon"
                className="h-9 w-9 text-muted-foreground hover:text-foreground rounded-xl"
              >
                <Mic className="h-5 w-5" />
                <span className="sr-only">Голосовое сообщение</span>
              </Button>
            </>
          )}
          {/* Send button - shows when there's text */}
          {text.trim() && (
            <Button onClick={handleSend} disabled={isSending} size="icon" className="h-9 w-9 rounded-xl transition-all">
              <Send className="h-4 w-4" />
              <span className="sr-only">Отправить</span>
            </Button>
          )}
        </div>
      </div>
      <p className="hidden md:block text-xs text-muted-foreground mt-2 text-center">
        Enter — отправить · Shift + Enter — новая строка
      </p>
    </div>
  )
}
