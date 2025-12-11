import { api } from "@/api/client";

export type ChatRoomType = "cohort" | "advisory" | "other";
export type ChatRole = "student" | "advisor" | "chair" | "admin" | "superadmin";

export type ChatRoom = {
  id: string;
  name: string;
  type: ChatRoomType;
  created_by: string;
  created_by_role?: ChatRole;
  is_archived: boolean;
  meta?: unknown;
  created_at: string;
  unread_count?: number;
  last_message_at?: string;
};

export type ChatAttachment = {
  url: string;
  type: string;
  name: string;
  size: number;
};

export type ChatMessage = {
  id: string;
  room_id: string;
  sender_id: string;
  sender_name?: string;
  sender_role?: ChatRole;
  body: string;
  attachments?: ChatAttachment[];
  importance?: "alert" | "warning" | null;
  created_at: string;
  edited_at?: string | null;
  deleted_at?: string | null;
};

export type ChatMember = {
  user_id: string;
  first_name?: string;
  last_name?: string;
  email?: string;
  role_in_room: string;
  last_read_at?: string;
  username?: string;
};

export async function listRooms(): Promise<ChatRoom[]> {
  const res = await api("/chat/rooms");
  return res.rooms ?? [];
}

export async function listRoomMembers(roomId: string): Promise<ChatMember[]> {
  const res = await api(`/chat/rooms/${roomId}/members`);
  return res.members ?? [];
}

export const createRoom = async (name: string, type: "group" | "channel", meta?: any): Promise<ChatRoom> => {
  const res = await api("/chat/rooms", {
    method: "POST",
    body: JSON.stringify({ name, type, meta: meta || {} }),
  });
  return res.room;
};

export const archiveRoom = async (roomId: string) => {
  return api(`/chat/rooms/${roomId}`, {
    method: "PUT",
    body: JSON.stringify({ is_archived: true }),
  });
};

export async function listMessages(params: {
  roomId: string;
  limit?: number;
  before?: string;
  after?: string;
}): Promise<ChatMessage[]> {
  const search = new URLSearchParams();
  if (params.limit) search.set("limit", String(params.limit));
  if (params.before) search.set("before", params.before);
  if (params.after) search.set("after", params.after);
  const qs = search.toString();
  const res = await api(
    `/chat/rooms/${params.roomId}/messages${qs ? `?${qs}` : ""}`
  );
  return res.messages ?? [];
}

export async function createMessage(roomId: string, body: string, attachments: any[] = [], importance?: string, meta?: any): Promise<ChatMessage> {
  const res = await api(`/chat/rooms/${roomId}/messages`, {
    method: "POST",
    body: JSON.stringify({ body, attachments, importance, meta }),
  });
  return res.message;
}

export async function getRoomMembers(roomId: string) {
  const res = await api(`/chat/rooms/${roomId}/members`);
  return res.members;
}

export async function updateMessage(messageId: string, body: string) {
  const res = await api(`/chat/messages/${messageId}`, {
    method: "PUT",
    body: JSON.stringify({ body }),
  });
  return res.message;
}

export async function deleteMessage(messageId: string) {
  return api(`/chat/messages/${messageId}`, {
    method: "DELETE",
  });
}

export async function uploadFile(roomId: string, file: File): Promise<ChatAttachment> {
  console.log("[uploadFile] Starting upload:", {
    roomId,
    fileName: file.name,
    fileSize: file.size,
    fileType: file.type,
  });
  
  const formData = new FormData();
  formData.append("file", file);
  
  console.log("[uploadFile] FormData created, entries:");
  for (const [key, value] of formData.entries()) {
    console.log(`  ${key}:`, value);
  }
  
  try {
    const res = await api(`/chat/rooms/${roomId}/upload`, {
      method: "POST",
      body: formData,
    });
    console.log("[uploadFile] Upload response:", res);
    return res as ChatAttachment;
  } catch (error) {
    console.error("[uploadFile] Upload failed:", error);
    throw error;
  }
}



export async function addMember(roomId: string, userId: string, role: string = "member") {
  return api(`/chat/rooms/${roomId}/members`, {
    method: "POST",
    body: JSON.stringify({ user_id: userId, role_in_room: role }),
  });
}

export async function markAsRead(roomId: string) {
  return api(`/chat/rooms/${roomId}/read`, {
    method: "POST",
  });
}

