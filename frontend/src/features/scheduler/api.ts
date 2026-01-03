import { api } from "@/lib/admin/api";
import { 
  AcademicTerm, 
  Building, 
  Room, 
  ClassSession, 
  CourseOffering, 
  SolverConfig 
} from "./types";

const BASE_URL = '/scheduling';
const RESOURCES_URL = '/resources'; // Now registered

export async function fetchTerms(tenantId?: string): Promise<AcademicTerm[]> {
  const query = tenantId ? `?tenant_id=${tenantId}` : '';
  return api(`${BASE_URL}/terms${query}`);
}

export async function createTerm(term: Partial<AcademicTerm>): Promise<AcademicTerm> {
  return api(`${BASE_URL}/terms`, {
    method: 'POST',
    body: JSON.stringify(term),
  });
}

export async function fetchBuildings(tenantId?: string): Promise<Building[]> {
  return api(`${RESOURCES_URL}/buildings`);
}

export async function fetchRooms(buildingId?: string): Promise<Room[]> {
  const query = buildingId ? `?building_id=${buildingId}` : '';
  return api(`${RESOURCES_URL}/rooms${query}`);
}

export async function fetchOfferings(termId?: string, tenantId?: string): Promise<CourseOffering[]> {
  const params = new URLSearchParams();
  if (termId) params.append('term_id', termId);
  if (tenantId) params.append('tenant_id', tenantId);
  return api(`${BASE_URL}/offerings?${params.toString()}`);
}

export async function fetchSessions(offeringId?: string, start?: Date, end?: Date): Promise<ClassSession[]> {
  const params = new URLSearchParams();
  if (offeringId) params.append('offering_id', offeringId);
  if (start) params.append('start', start.toISOString());
  if (end) params.append('end', end.toISOString());
  return api(`${BASE_URL}/sessions?${params.toString()}`);
}

export async function createSession(session: Partial<ClassSession>): Promise<ClassSession> {
  return api(`${BASE_URL}/sessions`, {
    method: 'POST',
    body: JSON.stringify(session),
  });
}

export async function optimizeSchedule(termId: string, config?: SolverConfig): Promise<any> {
    return api(`${BASE_URL}/optimize`, {
      method: 'POST',
      body: JSON.stringify({ term_id: termId, config }),
    });
}
