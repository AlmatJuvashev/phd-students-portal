import { api } from "@/api/client";

export type MonitorStudent = {
  id: string;
  name: string;
  email: string;
  phone?: string;
  program?: string;
  department?: string;
  cohort?: string;
  advisors?: { id: string; name: string; email?: string }[];
  rp_required?: boolean;
  overall_progress_pct?: number;
  last_update?: string;
  current_stage?: string;
  stage_done?: number;
  stage_total?: number;
  risk_score?: number;
  risk_level?: "HIGH" | "MEDIUM" | "LOW";
  last_risk_analysis?: string;
};

export async function runBatchAnalysis() {
  return api('/analytics/batch-analysis', { method: 'POST' });
}

export async function fetchMonitorStudents(params: Record<string, any> = {}) {
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && String(v) !== "") sp.set(k, String(v));
  });
  return api(
    `/admin/monitor/students${sp.toString() ? `?${sp.toString()}` : ""}`
  );
}

export type JourneyNode = {
  node_id: string;
  state: string;
  updated_at?: string;
  attachments?: number;
  files?: Array<{
    filename: string;
    download_url: string;
    size_bytes: number;
    attached_at?: string;
  }>;
};
export async function fetchStudentJourney(
  id: string
): Promise<{ nodes: JourneyNode[] }> {
  return api(`/admin/students/${id}/journey`);
}

export async function patchStudentNodeState(
  id: string,
  nodeId: string,
  state: string
) {
  return api(`/admin/students/${id}/nodes/${nodeId}/state`, {
    method: "PATCH",
    body: JSON.stringify({ state }),
  });
}

export async function fetchDeadlines(
  id: string
): Promise<{ node_id: string; due_at: string }[]> {
  return api(`/admin/students/${id}/deadlines`);
}

export async function putDeadline(
  id: string,
  nodeId: string,
  due_at: string,
  note?: string
) {
  return api(`/admin/students/${id}/nodes/${nodeId}/deadline`, {
    method: "PUT",
    body: JSON.stringify({ due_at, note }),
  });
}

export async function postReminders(payload: {
  student_ids: string[];
  title: string;
  message?: string;
  due_at?: string;
}) {
  return api(`/admin/reminders`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function fetchMonitorAnalytics(params: Record<string, any> = {}) {
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && String(v) !== "") sp.set(k, String(v));
  });
  return api(
    `/admin/monitor/analytics${sp.toString() ? `?${sp.toString()}` : ""}`
  );
}

export async function fetchStudentDetails(id: string): Promise<MonitorStudent> {
  return api(`/admin/students/${id}`);
}

export type NodeFileRow = {
  slot_key: string;
  attachment_id: string;
  filename: string;
  size_bytes: number;
  status: "submitted" | "approved" | "rejected";
  review_note?: string;
  attached_at?: string;
  approved_at?: string;
  approved_by?: string;
  uploaded_by?: string;
  version_id?: string;
  mime_type?: string;
  download_url: string;
  is_active: boolean;
};

export async function fetchStudentNodeFiles(
  studentId: string,
  nodeId: string
): Promise<NodeFileRow[]> {
  return api(`/admin/students/${studentId}/nodes/${nodeId}/files`);
}

export async function reviewAttachment(
  attachmentId: string,
  payload: { status: "approved" | "rejected" | "submitted"; note?: string }
) {
  return api(`/admin/attachments/${attachmentId}/review`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
