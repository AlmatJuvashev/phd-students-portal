import { api } from "./client";

export type Contact = {
  id: string;
  name: Record<string, string>;
  title?: Record<string, string>;
  email?: string | null;
  phone?: string | null;
  sort_order?: number;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
};

export type ContactPayload = {
  name: Record<string, string>;
  title?: Record<string, string>;
  email?: string;
  phone?: string;
  sort_order?: number;
  is_active?: boolean;
};

export async function fetchContacts(): Promise<Contact[]> {
  return api("/contacts");
}

export async function fetchAdminContacts(all = true): Promise<Contact[]> {
  const query = all ? "?all=true" : "";
  return api(`/admin/contacts${query}`);
}

export async function createContact(payload: ContactPayload) {
  return api("/admin/contacts", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function updateContact(id: string, payload: ContactPayload) {
  return api(`/admin/contacts/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function deleteContact(id: string) {
  return api(`/admin/contacts/${id}`, { method: "DELETE" });
}
