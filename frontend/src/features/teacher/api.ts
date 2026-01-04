import { api } from '@/api/client';
import { 
  TeacherDashboardStats,
  CourseOffering,
  CourseEnrollment,
  ActivitySubmission,
  GradebookEntry,
  GradeSubmissionRequest
} from './types';

export const getTeacherDashboard = () =>
  api.get<TeacherDashboardStats>('/teacher/dashboard');

export const getTeacherCourses = () =>
  api.get<CourseOffering[]>('/teacher/courses');

export const getCourseRoster = (id: string) => 
  api.get<CourseEnrollment[]>(`/teacher/courses/${id}/roster`);

export const getCourseGradebook = (id: string) =>
  api.get<GradebookEntry[]>(`/teacher/courses/${id}/gradebook`);

export const getTeacherSubmissions = () =>
  api.get<ActivitySubmission[]>('/teacher/submissions');

export const submitGradeForSubmission = (data: GradeSubmissionRequest) =>
  api.post<GradebookEntry>('/grading/entries', data);
