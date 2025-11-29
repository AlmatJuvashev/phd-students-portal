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

export type ChatMessage = {
  id: string;
  room_id: string;
  sender_id: string;
  sender_name?: string;
  sender_role?: ChatRole;
  body: string;
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

export async function sendMessage(roomId: string, body: string) {
  const res = await api(`/chat/rooms/${roomId}/messages`, {
    method: "POST",
    body: JSON.stringify({ body }),
  });
  return res.message as ChatMessage;
}
