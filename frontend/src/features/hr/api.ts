import { api } from '@/api/client';

// Types
export interface Staff {
  id: string;
  name: string;
  email: string;
  role: 'student' | 'advisor' | 'chair' | 'admin';
  username: string;
  program: string;
  specialty: string;
  department: string;
  cohort: string;
  created_at: string;
  is_active: boolean;
}

export interface StaffListResponse {
  data: Staff[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface CreateStaffRequest {
  first_name: string;
  last_name: string;
  email?: string;
  role: string;
  phone?: string;
  program?: string;
  specialty?: string;
  department?: string;
  cohort?: string;
  advisor_ids?: string[];
}

export interface UpdateStaffRequest extends CreateStaffRequest {
  id: string;
}

// API Functions
export const getStaffList = async (params?: {
  page?: number;
  limit?: number;
  role?: string;
  department?: string;
  search?: string;
  active?: string;
}): Promise<StaffListResponse> => {
  const searchParams = new URLSearchParams();
  if (params?.page) searchParams.set('page', String(params.page));
  if (params?.limit) searchParams.set('limit', String(params.limit));
  if (params?.role) searchParams.set('role', params.role);
  if (params?.department) searchParams.set('department', params.department);
  if (params?.search) searchParams.set('q', params.search);
  if (params?.active) searchParams.set('active', params.active);
  
  const res = await api.get(`/users?${searchParams.toString()}`);
  return res.data;
};

export const createStaff = async (data: CreateStaffRequest): Promise<{ username: string; temp_password: string }> => {
  const res = await api.post('/users', data);
  return res.data;
};

export const updateStaff = async (id: string, data: Omit<UpdateStaffRequest, 'id'>): Promise<void> => {
  await api.put(`/users/${id}`, data);
};

export const setStaffActive = async (id: string, active: boolean): Promise<void> => {
  await api.put(`/users/${id}/active`, { active });
};

export const resetStaffPassword = async (id: string): Promise<{ username: string; temp_password: string }> => {
  const res = await api.post(`/users/${id}/reset-password`);
  return res.data;
};
