import { api } from '@/api/client';
import { 
  TeacherDashboardStats,
  CourseOffering,
  CourseEnrollment,
  ActivitySubmission,
  GradebookEntry,
  GradeSubmissionRequest,
  StudentRiskProfile,
  TeacherStudentActivityEvent
} from './types';
import { submitGrade } from '@/features/grading/api';

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
  submitGrade(data);

export const getCourseStudents = (id: string) =>
  api.get<StudentRiskProfile[]>(`/teacher/courses/${id}/students`);

export const getAtRiskStudents = (id: string) =>
  api.get<StudentRiskProfile[]>(`/teacher/courses/${id}/at-risk`);

export const getStudentActivityLog = (studentId: string, courseOfferingId: string, limit = 50) =>
  api.get<TeacherStudentActivityEvent[]>(
    `/teacher/students/${studentId}/activity?course_offering_id=${encodeURIComponent(courseOfferingId)}&limit=${limit}`
  );
