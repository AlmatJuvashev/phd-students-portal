import { api } from '@/api/client';
import { Assessment, AssessmentListFilters, Question, AttemptDetailsResponse } from './types';

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

export const getAssessment = async (id: string): Promise<{ assessment: Assessment; questions: Question[] }> => {
  const response = await api.get(`/assessments/${id}`);
  return response.data;
};

export const startAttempt = async (assessmentId: string): Promise<{ id: string }> => {
  const response = await api.post(`/assessments/${assessmentId}/attempts`);
  return response.data;
};

export const getAttemptDetails = async (attemptId: string): Promise<AttemptDetailsResponse> => {
  const response = await api.get(`/attempts/${attemptId}`);
  return response.data;
};

export const submitResponse = async (
  attemptId: string,
  payload: { question_id: string; option_id?: string; text_response?: string }
): Promise<void> => {
  await api.post(`/attempts/${attemptId}/response`, payload);
};

export const completeAttempt = async (attemptId: string): Promise<void> => {
  await api.post(`/attempts/${attemptId}/complete`);
};
