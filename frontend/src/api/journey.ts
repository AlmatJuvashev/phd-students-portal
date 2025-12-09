import { api, API_URL } from "./client";

export type NodeSubmissionDTO = {
  node_id: string;
  playbook_version_id: string;
  state: string;
  locale?: string;
  form?: { rev: number; data: any };
  slots: Array<{
    key: string;
    required: boolean;
    multiplicity: "single" | "multi";
    mime: string[];
    attachments: Array<{
      version_id: string;
      filename: string;
      size_bytes: number;
      is_active: boolean;
      attached_at?: string;
      download_url: string;
      status?: "submitted" | "approved" | "rejected";
      review_note?: string;
      approved_at?: string;
      approved_by?: string;
      reviewed_document?: {
        version_id: string;
        filename?: string;
        size_bytes?: number;
        mime_type?: string;
        download_url: string;
        reviewed_at?: string;
        reviewed_by?: string;
      };
    }>;
  }>;
  outcomes?: Array<{
    value: string;
    decided_by: string;
    note?: string;
    created_at: string;
  }>;
};

export async function getNodeSubmission(nodeId: string) {
  return api(`/journey/nodes/${nodeId}/submission`);
}

export async function saveNodeSubmission(
  nodeId: string,
  payload: { form_data?: any; state?: string }
) {
  return api(`/journey/nodes/${nodeId}/submission`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function presignNodeUpload(
  nodeId: string,
  payload: {
    slot_key: string;
    filename: string;
    content_type: string;
    size_bytes: number;
  }
) {
  return api(`/journey/nodes/${nodeId}/uploads/presign`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function attachNodeUpload(
  nodeId: string,
  payload: {
    slot_key: string;
    filename: string;
    object_key: string;
    content_type: string;
    size_bytes: number;
    etag?: string;
  }
) {
  return api(`/journey/nodes/${nodeId}/uploads/attach`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function patchNodeState(
  nodeId: string,
  payload: { state: string }
) {
  return api(`/journey/nodes/${nodeId}/state`, {
    method: "PATCH",
    body: JSON.stringify(payload),
  });
}

export async function getProfileSnapshot() {
  const token = localStorage.getItem("token");
  const res = await fetch(`${API_URL}/journey/profile`, {
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });
  if (res.status === 404) {
    return null;
  }
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  return res.json();
}
