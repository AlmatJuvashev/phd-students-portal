import { api } from '@/api/client';
import { Enrollment, EnrollmentCreateRequest } from './types';

export const getEnrollments = (filters?: { program_id?: string; cohort_id?: string; search?: string }) => {
  const params = new URLSearchParams();
  if (filters?.program_id) params.append('program_id', filters.program_id);
  if (filters?.cohort_id) params.append('cohort_id', filters.cohort_id);
  if (filters?.search) params.append('search', filters.search);
  
  return api.get<Enrollment[]>(`/admin/enrollments?${params}`);
};

export const createEnrollment = (data: EnrollmentCreateRequest) => 
  api.post('/admin/enrollments', data);

export const bulkEnroll = (data: { program_id: string; cohort_id: string; student_ids: string[]; start_date: string }) =>
  api.post('/admin/enrollments/bulk', data);

export const updateEnrollmentStatus = (id: string, status: string) => 
  api.put(`/admin/enrollments/${id}/status`, { status });

export const dropEnrollment = (id: string) => 
  api.delete(`/admin/enrollments/${id}`);
