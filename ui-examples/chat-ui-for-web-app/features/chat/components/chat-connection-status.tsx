"use client"

import { cn } from "@/lib/utils"
import { WifiOff, Loader2 } from "lucide-react"

type ConnectionStatus = "connected" | "connecting" | "disconnected"

interface ChatConnectionStatusProps {
  status: ConnectionStatus
  className?: string
}

export function ChatConnectionStatus({ status, className }: ChatConnectionStatusProps) {
  return (
    <div
      className={cn(
        "flex items-center gap-1.5 text-[11px] font-medium rounded-full px-2 py-0.5 transition-colors",
        status === "connected" && "text-emerald-600 bg-emerald-50",
        status === "connecting" && "text-amber-600 bg-amber-50",
        status === "disconnected" && "text-red-600 bg-red-50",
        className,
      )}
    >
      {status === "connected" && (
        <>
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75" />
            <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500" />
          </span>
          <span className="hidden sm:inline">Подключено</span>
        </>
      )}
      {status === "connecting" && (
        <>
          <Loader2 className="h-3 w-3 animate-spin" />
          <span className="hidden sm:inline">Подключение...</span>
        </>
      )}
      {status === "disconnected" && (
        <>
          <WifiOff className="h-3 w-3" />
          <span className="hidden sm:inline">Отключено</span>
        </>
      )}
    </div>
  )
}
