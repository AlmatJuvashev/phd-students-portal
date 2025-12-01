import { api } from "@/api/client";
import { ChatRoom } from "@/features/chat/api";

export async function listAdminRooms(): Promise<ChatRoom[]> {
  const res = await api("/chat/rooms");
  return res.rooms ?? [];
}

export async function listRoomMembers(roomId: string) {
  const res = await api(`/chat/rooms/${roomId}/members`);
  return res.members as Array<{
    user_id: string;
    first_name?: string;
    last_name?: string;
    email?: string;
    role_in_room: string;
  }>;
}

export async function createRoom(input: {
  name: string;
  type: "cohort" | "advisory" | "other";
  meta?: unknown;
}) {
  const res = await api("/chat/rooms", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return res.room as ChatRoom;
}

export async function updateRoom(roomId: string, input: { name?: string; is_archived?: boolean }) {
  const res = await api(`/chat/rooms/${roomId}`, {
    method: "PATCH",
    body: JSON.stringify(input),
  });
  return res.room as ChatRoom;
}

export async function addRoomMember(roomId: string, payload: { user_id: string; role_in_room?: string }) {
  return api(`/chat/rooms/${roomId}/members`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function removeRoomMember(roomId: string, userId: string) {
  return api(`/chat/rooms/${roomId}/members/${userId}`, {
    method: "DELETE",
  });
}

export type UserSearchResult = {
  id: string;
  name: string;
  email: string;
  role: string;
};

export async function searchUsers(query: string, filters?: {
  program?: string;
  department?: string;
  cohort?: string;
  specialty?: string;
}) {
  const params = new URLSearchParams({ q: query, limit: "50" });
  if (filters?.program) params.append("program", filters.program);
  if (filters?.department) params.append("department", filters.department);
  if (filters?.cohort) params.append("cohort", filters.cohort);
  if (filters?.specialty) params.append("specialty", filters.specialty);

  const res = await api(`/admin/users?${params.toString()}`);
  return (res.data as UserSearchResult[]) ?? [];
}

export async function addRoomMembersBatch(roomId: string, payload: {
  user_ids?: string[];
  program?: string;
  department?: string;
  cohort?: string;
  specialty?: string;
  role?: string;
}) {
  return api(`/chat/rooms/${roomId}/members/batch`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function removeRoomMembersBatch(roomId: string, payload: {
  user_ids?: string[];
  program?: string;
  department?: string;
  cohort?: string;
  specialty?: string;
  role?: string;
}) {
  return api(`/chat/rooms/${roomId}/members/batch`, {
    method: "DELETE",
    body: JSON.stringify(payload),
  });
}
