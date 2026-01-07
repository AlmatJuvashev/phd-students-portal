import { api } from '@/api/client';
import { CourseContent, Module, Lesson, Activity } from './types';

export const getCourseContent = (courseId: string) => 
  api.get<CourseContent>(`/admin/studio/courses/${courseId}/content`);

export const updateCourseContent = (courseId: string, content: CourseContent) => 
  api.put(`/admin/studio/courses/${courseId}/content`, content);

export const addModule = (courseId: string, data: Partial<Module>) => 
  api.post(`/admin/studio/courses/${courseId}/modules`, data);

export const updateModule = (courseId: string, moduleId: string, data: Partial<Module>) => 
  api.put(`/admin/studio/courses/${courseId}/modules/${moduleId}`, data);

export const deleteModule = (courseId: string, moduleId: string) => 
  api.delete(`/admin/studio/courses/${courseId}/modules/${moduleId}`);

export const addLesson = (moduleId: string, data: Partial<Lesson>) => 
  api.post(`/admin/studio/modules/${moduleId}/lessons`, data);

export const addActivity = (lessonId: string, data: Partial<Activity>) => 
  api.post(`/admin/studio/lessons/${lessonId}/activities`, data);

export const updateActivity = (activityId: string, data: Partial<Activity>) => 
  api.put(`/admin/studio/activities/${activityId}`, data);

export const generateCourseStructure = (syllabus: string) =>
  api.post<{ modules: Module[] }>('/admin/ai/generate-course', { syllabus_text: syllabus });

export const generateQuiz = (topic: string, difficulty: string, count: number) =>
  api.post<any>('/admin/ai/generate-quiz', { topic, difficulty, count });

export const generateSurvey = (topic: string, count: number) =>
  api.post<any>('/admin/ai/generate-survey', { topic, count });

export const generateAssessmentItems = (topic: string, type: string, count: number) =>
  api.post<{ items: any[] }>('/admin/ai/generate-assessment-items', { topic, type, count });
