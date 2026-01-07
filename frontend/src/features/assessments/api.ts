import { api } from '@/api/client';
import { Assessment, AssessmentListFilters } from './types';

export const getAssessments = async (filters?: AssessmentListFilters): Promise<Assessment[]> => {
  const params = new URLSearchParams();
  if (filters?.course_id) params.append('course_id', filters.course_id);
  if (filters?.status) params.append('status', filters.status);
  if (filters?.search) params.append('search', filters.search);

  // Assuming endpoint is /api/assessments or /api/admin/assessments
  // Based on other features, likely /api/assessments for now
  const response = await api.get(`/assessments?${params.toString()}`);
  return response.data;
};

export const deleteAssessment = async (id: string): Promise<void> => {
  await api.delete(`/assessments/${id}`);
};

export const getAssessment = async (id: string): Promise<{ assessment: Assessment; questions: any[] }> => {
  const response = await api.get(`/assessments/${id}`);
  return response.data;
};
