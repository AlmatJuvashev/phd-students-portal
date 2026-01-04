import { api } from '@/api/client';
import type { AssessmentAttempt, AttemptDetailsResponse, AssessmentForTakingResponse } from './types';

export const startAttempt = (assessmentId: string) =>
  api.post<AssessmentAttempt>(`/assessments/${assessmentId}/attempts`);

export const getAssessmentForTaking = (assessmentId: string) =>
  api.get<AssessmentForTakingResponse>(`/assessments/${assessmentId}`);

export const getAttemptDetails = (attemptId: string) =>
  api.get<AttemptDetailsResponse>(`/attempts/${attemptId}`);

export const submitResponse = (
  attemptId: string,
  payload: { question_id: string; option_id?: string; text_response?: string }
) => api.post<void>(`/attempts/${attemptId}/response`, payload);

export const completeAttempt = (attemptId: string) =>
  api.post<AssessmentAttempt>(`/attempts/${attemptId}/complete`);

export const listMyAttempts = (assessmentId: string) =>
  api.get<AssessmentAttempt[]>(`/assessments/${assessmentId}/my-attempts`);

