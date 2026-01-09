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

export const deleteProgramVersionNode = (programId: string, nodeId: string) =>
  api.delete(`/curriculum/programs/${programId}/builder/nodes/${nodeId}`);

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

// Course Content (Modules, Lessons, Activities)
export const getCourseModules = (courseId: string) => api.get(`/curriculum/courses/${courseId}/modules`);
export const createCourseModule = (courseId: string, data: any) => api.post(`/curriculum/courses/${courseId}/modules`, data);
export const updateCourseModule = (id: string, courseId: string, data: any) => api.put(`/curriculum/courses/${courseId}/modules/${id}`, data);
export const deleteCourseModule = (id: string, courseId: string) => api.delete(`/curriculum/courses/${courseId}/modules/${id}`);

export const createCourseLesson = (moduleId: string, courseId: string, data: any) => api.post(`/curriculum/courses/${courseId}/modules/${moduleId}/lessons`, data);
export const updateCourseLesson = (id: string, moduleId: string, courseId: string, data: any) => api.put(`/curriculum/courses/${courseId}/modules/${moduleId}/lessons/${id}`, data);
export const deleteCourseLesson = (id: string, moduleId: string, courseId: string) => api.delete(`/curriculum/courses/${courseId}/modules/${moduleId}/lessons/${id}`);

export const createCourseActivity = (lessonId: string, moduleId: string, courseId: string, data: any) => api.post(`/curriculum/courses/${courseId}/modules/${moduleId}/lessons/${lessonId}/activities`, data);
export const updateCourseActivity = (id: string, lessonId: string, moduleId: string, courseId: string, data: any) => api.put(`/curriculum/courses/${courseId}/modules/${moduleId}/lessons/${lessonId}/activities/${id}`, data);
export const deleteCourseActivity = (id: string, lessonId: string, moduleId: string, courseId: string) => api.delete(`/curriculum/courses/${courseId}/modules/${moduleId}/lessons/${lessonId}/activities/${id}`);
