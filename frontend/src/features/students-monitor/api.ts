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
};

export async function fetchMonitorStudents(params: Record<string, any> = {}) {
  const sp = new URLSearchParams();
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && String(v) !== "") sp.set(k, String(v));
  });
  return api(`/admin/monitor/students${sp.toString() ? `?${sp.toString()}` : ""}`);
}

export type JourneyNode = { node_id: string; state: string; updated_at?: string; attachments?: number };
export async function fetchStudentJourney(id: string): Promise<{ nodes: JourneyNode[] }> {
  return api(`/admin/students/${id}/journey`);
}

export async function patchStudentNodeState(id: string, nodeId: string, state: string) {
  return api(`/admin/students/${id}/nodes/${nodeId}/state`, { method: "PATCH", body: JSON.stringify({ state }) });
}

