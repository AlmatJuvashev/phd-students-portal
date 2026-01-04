import { api } from '@/api/client';
import type { GradeSubmissionRequest, GradebookEntry, GradingSchema } from './types';

export const listGradingSchemas = () => api.get<GradingSchema[]>('/grading/schemas');
export const createGradingSchema = (data: Partial<GradingSchema>) => api.post<GradingSchema>('/grading/schemas', data);

export const submitGrade = (data: GradeSubmissionRequest) => api.post<GradebookEntry>('/grading/entries', data);
export const listStudentGrades = (studentId: string) => api.get<GradebookEntry[]>(`/grading/student/${studentId}`);

