import { api } from '@/api/client';
import { Program, Course } from './types';

// Programs
export const getPrograms = () => api.get<Program[]>('/curriculum/programs');
// Program Versions (Builder)
export const getProgramVersionMap = (programId: string) =>
  api.get(`/curriculum/programs/${programId}/builder/map`);

export const getProgramVersionNodes = (programId: string) =>
  api.get(`/curriculum/programs/${programId}/builder/nodes`);

export const createProgramVersionNode = (programId: string, data: any) =>
  api.post(`/curriculum/programs/${programId}/builder/nodes`, data);

export const updateProgramVersionNode = (programId: string, nodeId: string, data: any) =>
  api.put(`/curriculum/programs/${programId}/builder/nodes/${nodeId}`, data);

export const updateProgramVersionMap = (programId: string, data: any) =>
  api.put(`/curriculum/programs/${programId}/builder/map`, data);

// Backward-compatible aliases (deprecated)
export const getProgramJourneyMap = getProgramVersionMap;
export const getProgramNodes = getProgramVersionNodes;
export const createNode = createProgramVersionNode;
export const getProgram = (id: string) => api.get<Program>(`/curriculum/programs/${id}`);
export const createProgram = (data: Partial<Program>) => api.post('/curriculum/programs', data);
export const updateProgram = (id: string, data: Partial<Program>) => 
  api.put(`/curriculum/programs/${id}`, data);
export const deleteProgram = (id: string) => api.delete(`/curriculum/programs/${id}`);

// Courses  
export const getCourses = (programId?: string) => 
  api.get<Course[]>(`/curriculum/courses${programId ? `?program_id=${programId}` : ''}`);
export const getCourse = (id: string) => api.get<Course>(`/curriculum/courses/${id}`);
export const createCourse = (data: Partial<Course>) => api.post('/curriculum/courses', data);
export const updateCourse = (id: string, data: Partial<Course>) => 
  api.put(`/curriculum/courses/${id}`, data);
export const deleteCourse = (id: string) => api.delete(`/curriculum/courses/${id}`);
