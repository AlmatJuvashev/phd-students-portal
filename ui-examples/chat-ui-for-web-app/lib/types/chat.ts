// Core chat types for the WebApp Checklist chat module

export type ChatRoomType = "cohort" | "advisory" | "other"
export type UserRole = "admin" | "advisor" | "student"

export interface ChatRoom {
  id: string
  name: string
  type: ChatRoomType
  description?: string
  isArchived?: boolean
  unreadCount?: number
  membersCount?: number
  createdAt?: string
}

export interface ChatMessage {
  id: string
  roomId: string
  senderId: string
  senderName: string
  senderRole: UserRole
  body: string
  createdAt: string // ISO
  isOwn: boolean
}

export interface ChatMember {
  id: string
  name: string
  email: string
  role: UserRole
  avatarUrl?: string
}

export interface RoomFormValues {
  name: string
  type: ChatRoomType
  description?: string
}
