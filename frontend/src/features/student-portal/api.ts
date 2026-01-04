import { api } from '@/api/client';
import type { StudentAssignment, StudentCourse, StudentDashboard, StudentGradeEntry } from './types';

export const getStudentDashboard = () => api.get<StudentDashboard>('/student/dashboard');
export const getStudentCourses = () => api.get<StudentCourse[]>('/student/courses');
export const getStudentAssignments = () => api.get<StudentAssignment[]>('/student/assignments');
export const getStudentGrades = () => api.get<StudentGradeEntry[]>('/student/grades');

