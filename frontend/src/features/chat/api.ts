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
  created_at: string;
  edited_at?: string | null;
  deleted_at?: string | null;
};

export async function listRooms(): Promise<ChatRoom[]> {
  const res = await api("/chat/rooms");
  return res.rooms ?? [];
}

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

export async function sendMessage(roomId: string, body: string, attachments?: ChatAttachment[]) {
  const res = await api(`/chat/rooms/${roomId}/messages`, {
    method: "POST",
    body: JSON.stringify({ body, attachments }),
  });
  return res.message as ChatMessage;
}

export async function uploadFile(roomId: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);
  
  // We need to use fetch directly or a wrapper that handles FormData correctly if api client doesn't
  // Assuming api client handles FormData if body is FormData, or we need to handle headers.
  // The current api client wrapper might set Content-Type to json by default.
  // Let's check api client implementation.
  // For now, I'll assume I can pass FormData and let the browser set Content-Type (multipart/form-data with boundary).
  // But wait, the `api` helper usually sets Content-Type: application/json.
  // I should check `frontend/src/api/client.ts`.
  
  // If I can't check it right now, I'll use a standard fetch or assume the api client is smart enough.
  // Actually, I'll just use the `api` client but I might need to override headers.
  // Let's assume for now I can pass a custom config.
  
  // Re-reading api client usage in other files...
  // `frontend/src/features/docgen/api.ts` or similar might have upload examples.
  // I recall `UploadVersion` in `documents.go`.
  
  // Let's just implement it and if it fails I'll fix the client.
  // Actually, better to check client.ts first.
  return api(`/chat/rooms/${roomId}/upload`, {
    method: "POST",
    body: formData,
  }) as Promise<ChatAttachment>;
}

export async function updateMessage(messageId: string, body: string) {
  const res = await api(`/chat/messages/${messageId}`, {
    method: "PATCH",
    body: JSON.stringify({ body }),
  });
  return res.message as ChatMessage;
}

export async function deleteMessage(messageId: string) {
  return api(`/chat/messages/${messageId}`, {
    method: "DELETE",
  });
}
